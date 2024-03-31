package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"url-shortener/internal/config"
	"url-shortener/internal/http-server/handlers/url/save"
	"url-shortener/internal/http-server/middleware"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage/sqllite"
)

func main() {
	cfg := config.MustLoad()
	logger := initLogger(cfg.Env)

	logger.Info("starting application", slog.String("env", cfg.Env))

	storage, err := sqllite.New(cfg.StoragePath)
	if err != nil {
		logger.Error("failed to initialize storage", sl.Err(err))
		os.Exit(1)
	}
	defer func(storage *sqllite.Storage) {
		err := storage.Close()
		if err != nil {
			logger.Error("failed to close storage", sl.Err(err))
		}
	}(storage)

	logger.Info("storage initialized", slog.String("storage", "sqllite"))

	router := http.NewServeMux()
	router.HandleFunc("POST /api/url", save.New(logger, storage))

	mw := middleware.Chain(
		middleware.RequestIDMiddleware,
		middleware.NewLoggingMiddleware(logger),
		middleware.RecovererMiddleware,
		middleware.ContentTypeJsonMiddleware,
	)

	server := &http.Server{
		Addr:         cfg.HttpServer.Address,
		Handler:      mw(router),
		WriteTimeout: cfg.HttpServer.Timeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
		ReadTimeout:  cfg.HttpServer.Timeout,
	}

	startServer(server, logger)
}

func initLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case config.LocalEnv:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case config.ProdEnv:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return logger
}

func startServer(server *http.Server, logger *slog.Logger) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("listen and serve returned err", sl.Err(err))
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	shutdownGracefully(server, logger)
}

func shutdownGracefully(server *http.Server, logger *slog.Logger) {
	logger.Info("gracefully shutting down")
	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(c); err != nil {
		logger.Error("server shutdown returned an err: %v\n", err)
	}

	logger.Info("server stopped")
}

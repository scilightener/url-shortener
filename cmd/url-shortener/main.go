package main

import (
	"log/slog"
	"net/http"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/http-server/handlers/url/save"
	"url-shortener/internal/http-server/middleware"
	"url-shortener/internal/storage/sqllite"
)

func main() {
	cfg := config.MustLoad()
	logger := initLogger(cfg.Env)

	logger.Info("starting application", slog.String("env", cfg.Env))

	storage, err := sqllite.New(cfg.StoragePath)
	if err != nil {
		logger.Error("failed to initialize storage", slog.String("error", err.Error()))
		os.Exit(1)
	}

	logger.Info("storage initialized", slog.String("storage", "sqllite"))

	router := http.NewServeMux()
	router.HandleFunc("POST /api/url", save.New(logger, storage))

	mw := middleware.Chain(
		middleware.RequestIDMiddleware,
		middleware.NewLoggingMiddleware(logger),
		middleware.RecovererMiddleware,
	)

	server := http.Server{
		Addr:         cfg.HttpServer.Address,
		Handler:      mw(router),
		WriteTimeout: cfg.HttpServer.Timeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
		ReadTimeout:  cfg.HttpServer.Timeout,
	}
	err = server.ListenAndServe()
	if err != nil {
		logger.Error("server error", slog.String("error", err.Error()))
		os.Exit(1)
	}
	logger.Info("server stopped")
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

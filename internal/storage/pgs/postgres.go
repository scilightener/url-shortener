package pgs

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
	"url-shortener/internal/storage"
)

type Storage struct {
	dbPool *pgxpool.Pool
}

func New(connectionString string) (*Storage, error) {
	const eo = "storage.pgs.New"

	dbPool, err := pgxpool.New(context.Background(), connectionString)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", eo, err)
	}

	batch := new(pgx.Batch)
	batch.Queue(`
		CREATE TABLE IF NOT EXISTS urls (
    		id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    		alias TEXT NOT NULL UNIQUE,
    		valid_until TIMESTAMP NOT NULL,
    		url TEXT NOT NULL);
	`)
	batch.Queue("CREATE INDEX IF NOT EXISTS idx_alias ON urls(alias);")
	br := dbPool.SendBatch(context.Background(), batch)
	_, err = br.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", eo, err)
	}

	return &Storage{dbPool: dbPool}, nil
}

func (s *Storage) SaveUrl(ctx context.Context, alias, url string, validUntilUTC time.Time) (int64, error) {
	const eo = "storage.pgs.SaveUrl"

	row := s.dbPool.QueryRow(ctx,
		"INSERT INTO urls (alias, valid_until, url) VALUES ($1, $2, $3) RETURNING id", alias, validUntilUTC, url)
	var id int64
	err := row.Scan(&id)
	if err != nil {
		pgErr := new(pgconn.PgError)
		// 23505: unique_violation (the only unique constraint is on the alias column)
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, fmt.Errorf("%s: %w", eo, storage.ResourceAlreadyExists)
		}

		return 0, fmt.Errorf("%s: %w", eo, err)
	}

	return id, nil
}

func (s *Storage) GetUrl(ctx context.Context, alias string) (string, error) {
	const eo = "storage.pgs.GetUrl"

	row := s.dbPool.QueryRow(ctx,
		"SELECT url FROM urls WHERE alias = $1 AND valid_until > $2", alias, time.Now().UTC())

	var url string
	err := row.Scan(&url)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", eo, storage.ResourceNotFound)
		}

		return "", fmt.Errorf("%s: %w", eo, err)
	}

	return url, nil
}

func (s *Storage) DeleteUrl(ctx context.Context, alias string) error {
	const eo = "storage.pgs.DeleteUrl"

	_, err := s.dbPool.Exec(ctx, "DELETE FROM urls WHERE alias = $1", alias)
	if err != nil {
		return fmt.Errorf("%s: %w", eo, err)
	}

	return nil
}

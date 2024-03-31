package sqllite

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
	"url-shortener/internal/storage"

	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const eo = "storage.sqlLite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", eo, err)
	}

	query, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS urls (
    		id INTEGER PRIMARY KEY,
    		alias TEXT NOT NULL UNIQUE,
    		valid_until TIMESTAMP NOT NULL,
    		url TEXT NOT NULL);
		CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", eo, err)
	}

	_, err = query.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", eo, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Close() error {
	const eo = "storage.sqlLite.Close"

	if err := s.db.Close(); err != nil {
		return fmt.Errorf("%s: %w", eo, err)
	}

	return nil
}

func (s *Storage) SaveUrl(alias, url string, validUntilUTC time.Time) (int64, error) {
	const eo = "storage.sqlLite.SaveUrl"

	query, err := s.db.Prepare("INSERT INTO urls (alias, valid_until, url) VALUES (?, ?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", eo, err)
	}

	res, err := query.Exec(alias, validUntilUTC, url)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && errors.Is(sqliteErr.Code, sqlite3.ErrConstraint) {
			return 0, fmt.Errorf("%s: %w", eo, storage.ResourceAlreadyExists)
		}

		return 0, fmt.Errorf("%s: %w", eo, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", eo, err)
	}

	return id, nil
}

func (s *Storage) GetUrl(alias string) (string, error) {
	const eo = "storage.sqlLite.GetUrl"

	query, err := s.db.Prepare("SELECT url FROM urls WHERE alias = ? AND valid_until > ?")
	if err != nil {
		return "", fmt.Errorf("%s: %w", eo, err)
	}

	var url string
	err = query.QueryRow(alias, 0).Scan(&url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", eo, storage.ResourceNotFound)
		}

		return "", fmt.Errorf("%s: %w", eo, err)
	}

	return url, nil
}

func (s *Storage) GetUrlById(id int64) (string, error) {
	const eo = "storage.sqlLite.GetUrlById"

	query, err := s.db.Prepare("SELECT url FROM urls WHERE id = ? AND valid_until > ?")
	if err != nil {
		return "", fmt.Errorf("%s: %w", eo, err)
	}

	var url string
	err = query.QueryRow(id, 0).Scan(&url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", eo, storage.ResourceNotFound)
		}

		return "", fmt.Errorf("%s: %w", eo, err)
	}

	return url, nil
}

func (s *Storage) DeleteUrl(alias string) error {
	const eo = "storage.sqlLite.DeleteUrl"

	query, err := s.db.Prepare("DELETE FROM urls WHERE alias = ?")
	if err != nil {
		return fmt.Errorf("%s: %w", eo, err)
	}

	_, err = query.Exec(alias)
	if err != nil {
		return fmt.Errorf("%s: %w", eo, err)
	}

	return nil
}

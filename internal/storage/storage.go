package storage

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

var (
	ErrURLNotFound = errors.New("Url not found")
	ErrURLExists   = errors.New("URL already exists")
)

type Storage struct {
	db *sql.DB
}

func New(db *sql.DB) (*Storage, error) {
	const op = "storage.New"

	if db == nil {
		return nil, fmt.Errorf("%s: db is nil", op)
	}

	// Create table and index if they don't exist
	if _, err := db.Exec(`
CREATE TABLE IF NOT EXISTS public.url (
    id SERIAL PRIMARY KEY,
    alias TEXT NOT NULL UNIQUE,
    url TEXT NOT NULL
);
`); err != nil {
		return nil, fmt.Errorf("%s: create table: %w", op, err)
	}

	if _, err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_alias ON public.url(alias);`); err != nil {
		return nil, fmt.Errorf("%s: create index: %w", op, err)
	}

	return &Storage{db: db}, nil
}

// SaveURL inserts a new URL and alias, returning the generated id.
func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	const op = "storage.SaveURL"

	if s == nil || s.db == nil {
		return 0, fmt.Errorf("%s: db is nil", op)
	}

	var id int64
	// Use PostgreSQL-style placeholders and RETURNING to get the id
	err := s.db.QueryRow(
		`INSERT INTO public.url (url, alias) VALUES ($1, $2) RETURNING id`,
		urlToSave, alias,
	).Scan(&id)
	if err != nil {
		// Unique violation code for Postgres is 23505
		if pgErr, ok := err.(*pq.Error); ok && string(pgErr.Code) == "23505" {
			return 0, fmt.Errorf("%s: unique violation: %w", op, err)
		}
		return 0, fmt.Errorf("%s: insert: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.GetURL"

	if s == nil || s.db == nil {
		return "", fmt.Errorf("%s: db is nil", op)
	}

	stmt, err := s.db.Prepare("SELECT url FROM public.url WHERE alias = $1")
	if err != nil {
		return "", fmt.Errorf("%s: prepare: %w", op, err)
	}
	defer stmt.Close()

	var ResUrl string
	err = stmt.QueryRow(alias).Scan(&ResUrl)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("%s: %w", op, ErrURLNotFound)
	}
	if err != nil {
		return "", fmt.Errorf("%s: scan: %w", op, err)
	}
	if ResUrl == "" {
		return "", fmt.Errorf("%s: %w", op, ErrURLNotFound)
	}

	return ResUrl, nil
}

func (s *Storage) DeleteURL(alias string) error {
	const op = "storage.DeleteURL"

	if s == nil || s.db == nil {
		return fmt.Errorf("%s: db is nil", op)
	}

	stmt, err := s.db.Prepare("DELETE FROM public.url WHERE alias = $1")
	if err != nil {
		return fmt.Errorf("%s: prepare delete: %w", op, err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(alias)
	if err != nil {
		return fmt.Errorf("%s: delete: %w", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: get rows affected: %w", op, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%s: %w", op, ErrURLNotFound)
	}

	return nil
}

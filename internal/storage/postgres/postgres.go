package postgres

import (
	"database/sql"
	"fmt"
	"url-shortener/internal/storage"

	"github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(pgCfg string) (*Storage, error) {
	const op = "storage.postgres.New"

	db, err := sql.Open("postgres", pgCfg)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(url string, alias string) (int64, error) {
	const op = "storage.postgres.SaveURL"

	var id int64
	err := s.db.QueryRow("INSERT INTO urls(url, alias) VALUES($1, $2) RETURNING id", url, alias).Scan(&id)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			return 0, fmt.Errorf("%s: %s", op, err.Code.Name())
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.postgres.GetURL"

	rows, err := s.db.Query("SELECT url FROM urls WHERE alias = $1", alias)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			return "", fmt.Errorf("%s: %s", op, err.Code.Name())
		}

		return "", fmt.Errorf("%s: %w", op, err)
	}

	if urlNotFount := !rows.Next(); urlNotFount {
		return "", storage.ErrURLNotFound
	}

	var url string
	err = rows.Scan(&url)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			return "", fmt.Errorf("%s: %s", op, err.Code.Name())
		}

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return url, nil
}

func (s *Storage) DeleteURL(url string) error {
	const op = "storage.postgres.DeleteURL"

	stmt, err := s.db.Prepare("DELETE FROM urls WHERE url = $1")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(url)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			return fmt.Errorf("%s: %s", op, err.Code.Name())
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

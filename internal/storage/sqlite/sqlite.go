package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/akamaaru/url-shortener/internal/storage"
	"modernc.org/sqlite"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS url(
		id INTEGER PRIMARY KEY,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_alias ON url(alias)
	`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) (error) {
	const op = "storage.sqlite.SaveURL"

	stmt, err := s.db.Prepare("INSERT INTO url(url, alias) VALUES(?, ?)")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(urlToSave, alias)
	if err != nil {
		if sqliteErr, ok := err.(*sqlite.Error); ok && sqliteErr.Code() == 19 { // sqlite3.SQLITE_CONSTRAINT = 19
			return fmt.Errorf("%s: %w", op, storage.ErrURLExists)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	// TODO remove if needed
	_, err = res.LastInsertId()
	if err != nil {
		return fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}

	return nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.sqlite.GetURL"

	// exists, err := s.ExistsURL(alias)
	// if err != nil {
	// 	return "", fmt.Errorf("%s: %w", op, err)
	// }

	// if !exists {
	// 	return "", storage.ErrURLNotFound
	// }

	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias = ?")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRow(alias)
	if row.Err() != nil {
		return "", storage.ErrURLNotFound
	}

	var item string
	err = row.Scan(&item)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return item, nil
}

func (s *Storage) DeleteURL(alias string) error {
	const op = "storage.sqlite.DeleteURL"

	// exists, err := s.ExistsURL(alias)
	// if err != nil {
	// 	return fmt.Errorf("%s: %w", op, err)
	// }

	// if !exists {
	// 	return storage.ErrURLNotFound
	// }

	stmt, err := s.db.Prepare("DELETE FROM url WHERE alias = ?")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(alias)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) ExistsURL(alias string) (bool, error) {
	const op = "storage.sqlite.ExistsURL"

	stmt, err := s.db.Prepare("SELECT COUNT(1) FROM url WHERE alias = ?")
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(alias)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return rows != 0, nil
}
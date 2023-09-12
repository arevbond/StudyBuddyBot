package sqlite

import (
	"context"
	"database/sql"
	"tg_ics_useful_bot/lib/e"
	"tg_ics_useful_bot/storage"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

// New creates new SQLite storage.
func New(path string) (storage *Storage, err error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, e.Wrap("can't open db (probably wrong path): ", err)
	}

	if err := db.Ping(); err != nil {
		return nil, e.Wrap("can't open db: cant't ping db: ", err)
	}
	return &Storage{db: db}, nil
}

// Save saves page to storage.
func (s *Storage) Save(ctx context.Context, p *storage.Page) error {
	q := `INSERT INTO pages (url, user_name) VALUES (?, ?)`

	if _, err := s.db.ExecContext(ctx, q, p.URL, p.UserName); err != nil {
		return e.Wrap("can't save page:", err)
	}
	return nil
}

// PickRandom picks random page from storage.
func (s *Storage) PickRandom(ctx context.Context, userName string) (*storage.Page, error) {
	q := `SELECT url FROM pages WHERE user_name = ? ORDER BY RANDOM() LIMIT 1`

	var url string

	err := s.db.QueryRowContext(ctx, q, userName).Scan(&url)
	if err == sql.ErrNoRows {
		return nil, storage.ErrNoSavedPages
	}
	if err != nil {
		return nil, e.Wrap("can't pick random page", err)
	}
	return &storage.Page{URL: url, UserName: userName}, nil
}

// Remove removes page from storage.
func (s *Storage) Remove(ctx context.Context, p *storage.Page) error {
	q := `DELETE FROM pages WHERE url = ? AND user_name = ?`
	if _, err := s.db.ExecContext(ctx, q, p.URL, p.UserName); err != nil {
		return e.Wrap("can't remove page", err)
	}
	return nil
}

// IsExists checks if page exists in storage.
func (s *Storage) IsExists(ctx context.Context, p *storage.Page) (bool, error) {
	q := `SELECT COUNT(*) FROM pages WHERE url = ? AND user_name = ?`

	var count int

	if err := s.db.QueryRowContext(ctx, q, p.URL, p.UserName).Scan(&count); err != nil {
		return false, e.Wrap("can't check page is exist", err)
	}
	return count > 0, nil
}

// Init creates tables to storage.
func (s *Storage) Init(ctx context.Context) error {
	q := `CREATE TABLE IF NOT EXISTS pages (url TEXT, user_name TEXT)`

	_, err := s.db.ExecContext(ctx, q)
	if err != nil {
		return e.Wrap("can't create tables", err)
	}
	return nil
}

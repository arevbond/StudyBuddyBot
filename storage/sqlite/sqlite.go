package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"tg_ics_useful_bot/lib/e"
	"tg_ics_useful_bot/storage"
	"time"
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

// Init creates tables to storage.
func (s *Storage) Init(ctx context.Context) error {
	q1 := `CREATE TABLE IF NOT EXISTS pages (url TEXT, user_name TEXT)`
	q2 := `CREATE TABLE IF NOT EXISTS users (tg_id int, chat_id int, is_bot BIT, first_name TEXT, last_name TEXT, 
			username TEXT, is_premium BIT, dick_size INT, last_try_change_dick DATE)`
	_, err := s.db.ExecContext(ctx, q1)
	if err != nil {
		return e.Wrap("can't create table pages", err)
	}

	_, err = s.db.ExecContext(ctx, q2)
	if err != nil {
		return e.Wrap("can't create table users", err)
	}
	return nil
}

// CreateUser new user by chatID and telegramID.
func (s *Storage) CreateUser(ctx context.Context, u *storage.DBUser) error {
	q := `INSERT INTO users (tg_id, chat_id, is_bot, first_name, last_name, username, is_premium, dick_size, last_try_change_dick) 
							VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	log.Printf("create user %d %s %s %s, chat_id = %d, dick size = %d", u.TgID, u.Username, u.FirstName,
		u.LastName, u.ChatID, u.DickSize)

	if _, err := s.db.ExecContext(ctx, q, u.TgID, u.ChatID, u.IsBot, u.FirstName,
		u.LastName, u.Username, u.IsPremium, u.DickSize, u.LastTryChangeDick); err != nil {
		return e.Wrap(fmt.Sprintf("can't create user %d %s: ", u.TgID, u.Username), err)
	}
	return nil
}

// User got user by chatID and telegram ID.
func (s *Storage) User(ctx context.Context, tgID, chatID int) (*storage.DBUser, error) {
	q := `SELECT * FROM users WHERE tg_id = ? AND chat_id = ?`

	var id, cID, dickSize int
	var isBot, isPremium bool
	var firstName, lastName, username string
	var lastTryChangeDickStr time.Time

	err := s.db.QueryRowContext(ctx, q, tgID, chatID).Scan(&id, &cID, &isBot, &firstName, &lastName,
		&username, &isPremium, &dickSize, &lastTryChangeDickStr)

	user := &storage.DBUser{
		TgID:              id,
		ChatID:            cID,
		IsBot:             isBot,
		FirstName:         firstName,
		LastName:          lastName,
		Username:          username,
		IsPremium:         isPremium,
		DickSize:          dickSize,
		LastTryChangeDick: lastTryChangeDickStr,
	}

	if err == sql.ErrNoRows {
		return nil, storage.ErrUserNotExist
	}

	if err != nil {
		return nil, e.Wrap(fmt.Sprintf("can't get user from storage tg id: %d, chat id: %d", tgID, chatID), err)
	}

	// log.Printf("from storage get user: tg id = %d, chat id = %d, dick size = %d", user.TgID, user.ChatID, user.DickSize)

	return user, nil
}

func (s *Storage) UpdateUserDickSize(ctx context.Context, u *storage.DBUser, dickSize int) error {
	q := `UPDATE users SET dick_size = ?, last_try_change_dick = ? WHERE tg_id = ? AND chat_id = ?`
	oldDickSize := u.DickSize
	if _, err := s.db.ExecContext(ctx, q, dickSize, time.Now(), u.TgID, u.ChatID); err != nil {
		return e.Wrap(fmt.Sprintf("can't update dick size user %d chat id %d from %d to %d",
			u.TgID, u.ChatID, u.DickSize, dickSize), err)
	}
	u.DickSize = dickSize
	log.Printf("user %d change his dick from %d to %d", u.TgID, oldDickSize, u.DickSize)
	return nil
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

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
	q2 := `CREATE TABLE IF NOT EXISTS users (tg_id int, chat_id int, is_bot BIT, first_name TEXT, last_name TEXT, 
			username TEXT, is_premium BIT, dick_size INT, last_try_change_dick DATE)`

	_, err := s.db.ExecContext(ctx, q2)
	if err != nil {
		return e.Wrap("can't create table users", err)
	}
	return nil
}

// CreateUser new user by chatID and telegramID.
func (s *Storage) CreateUser(ctx context.Context, u *storage.DBUser) error {
	q := `INSERT INTO users (tg_id, chat_id, is_bot, first_name, last_name, username, is_premium, dick_size, last_try_change_dick) 
							VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	log.Printf("create user #%d '%s' '%s' '%s', chat_id = %d, dick size = %d", u.TgID, u.Username, u.FirstName,
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

	user := &storage.DBUser{}

	err := s.db.QueryRowContext(ctx, q, tgID, chatID).Scan(&user.TgID, &user.ChatID, &user.IsBot, &user.FirstName, &user.LastName,
		&user.Username, &user.IsPremium, &user.DickSize, &user.LastTryChangeDick)

	if err == sql.ErrNoRows {
		return nil, storage.ErrUserNotExist
	}

	if err != nil {
		return nil, e.Wrap(fmt.Sprintf("can't get user from storage tg id: %d, chat id: %d", tgID, chatID), err)
	}

	// log.Printf("from storage get user: tg id = %d, chat id = %d, dick size = %d", user.TgID, user.ChatID, user.DickSize)

	return user, nil
}

func (s *Storage) UserByUsername(ctx context.Context, username string, chatID int) (*storage.DBUser, error) {
	q := `SELECT * FROM users WHERE username = ? AND chat_id = ?`
	log.Printf(username)
	user := &storage.DBUser{}

	err := s.db.QueryRowContext(ctx, q, username, chatID).Scan(&user.TgID, &user.ChatID, &user.IsBot, &user.FirstName, &user.LastName,
		&user.Username, &user.IsPremium, &user.DickSize, &user.LastTryChangeDick)

	if err == sql.ErrNoRows {
		return nil, storage.ErrUserNotExist
	}

	if err != nil {
		return nil, e.Wrap(fmt.Sprintf("can't get user from storage username: %s, chat id: %d", username, chatID), err)
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

func (s *Storage) UsersByChat(ctx context.Context, chatID int) ([]*storage.DBUser, error) {
	q := `SELECT * FROM users WHERE chat_id = ? ORDER BY -dick_size`
	rows, err := s.db.QueryContext(ctx, q, chatID)
	if err != nil {
		return nil, e.Wrap(fmt.Sprintf("can't get users by chat id: %s", chatID), err)
	}
	defer rows.Close()

	var users []*storage.DBUser

	for rows.Next() {
		user := &storage.DBUser{}
		if err := rows.Scan(&user.TgID, &user.ChatID, &user.IsBot, &user.FirstName, &user.LastName,
			&user.Username, &user.IsPremium, &user.DickSize, &user.LastTryChangeDick); err != nil {
			return users, e.Wrap(fmt.Sprintf("can't get users by chat id: %s", chatID), err)
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return users, err
	}
	return users, nil
}

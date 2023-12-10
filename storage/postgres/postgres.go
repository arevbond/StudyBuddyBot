package postgres

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"log"
	"tg_ics_useful_bot/config"
	"tg_ics_useful_bot/lib/e"
	"tg_ics_useful_bot/storage"
	"time"
)

type Storage struct {
	db *sqlx.DB
}

// New создаёт подключение к PostgreSQL базе данных.
func New(cfg *config.Config) (*Storage, error) {
	dbSource := fmt.Sprintf("postgres://%s:%s@localhost:5430/%s", cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDBName)
	conn, err := sqlx.Connect("pgx", dbSource)
	if err != nil {
		return nil, e.Wrap("connect to pgx failed", err)
	}

	err = conn.Ping()
	if err != nil {
		return nil, errors.Wrap(err, "ping failed")
	}
	return &Storage{db: conn}, nil
}

// CreateUser создаёт нового пользователя из телеграмма в базе данных.
func (s *Storage) CreateUser(ctx context.Context, u *storage.DBUser) error {
	q := `INSERT INTO users (tg_id, chat_id, is_bot, is_premium, first_name, last_name, username,
        dick_size, change_dick_at, user_stat_id, health_points, hp_taked_at, is_gay, gay_at, points, cur_dick_change_count, max_dick_change_count)
                   VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)`

	if _, err := s.db.ExecContext(ctx, q, u.TgID, u.ChatID, u.IsBot, u.IsPremium, u.FirstName, u.LastName, u.Username, u.DickSize, u.ChangeDickAt,
		u.UserStatId, u.HealthPoints, u.HpTakedAt, u.IsGay, u.GayAt, u.Points, u.CurDickChangeCount, u.MaxDickChangeCount); err != nil {
		return e.Wrap(fmt.Sprintf("can't create user %d %s: ", u.TgID, u.Username), err)
	}
	log.Printf("[INFO] create user #%d '%s' '%s' '%s', chat_id = %d, dick size = %d", u.TgID, u.Username, u.FirstName, u.LastName, u.ChatID, u.DickSize)
	return nil
}

// UpdateUser обновляет всю информацию о пользователей в базе данных.
func (s *Storage) UpdateUser(ctx context.Context, u *storage.DBUser) error {
	q := `UPDATE users SET is_premium = $1, first_name = $2, last_name = $3, username = $4, dick_size = $5, change_dick_at = $6, 
    		health_points = $7, hp_taked_at = $8, is_gay = $9, gay_at = $10, points = $11, cur_dick_change_count = $12, max_dick_change_count = $13 
             WHERE id = $14`
	_, err := s.db.ExecContext(ctx, q, u.IsPremium, u.FirstName, u.LastName, u.Username, u.DickSize, u.ChangeDickAt,
		u.HealthPoints, u.HpTakedAt, u.IsGay, u.GayAt, u.Points, u.CurDickChangeCount, u.MaxDickChangeCount, u.ID)

	if err != nil {
		return e.Wrap("can't update user", err)
	}

	return nil
}

// GetUser возвращает пользователя из базы данных по его телеграм id и чат id.
func (s *Storage) GetUser(ctx context.Context, tgID, chatID int) (*storage.DBUser, error) {
	q := `SELECT * FROM users WHERE tg_id = $1 AND chat_id = $2`

	var user storage.DBUser

	err := s.db.GetContext(ctx, &user, q, tgID, chatID)

	if err == sql.ErrNoRows {
		return nil, storage.ErrUserNotExist
	}

	if err != nil {
		return nil, e.Wrap(fmt.Sprintf("can't get user from storage tg id: %d, chat id: %d", tgID, chatID), err)
	}

	return &user, nil
}

func (s *Storage) UserByUsername(ctx context.Context, username string, chatID int) (*storage.DBUser, error) {
	q := `SELECT * FROM users WHERE username = $1 AND chat_id = $2`

	user := &storage.DBUser{}

	err := s.db.Get(&user, q)

	if err == sql.ErrNoRows {
		return nil, storage.ErrUserNotExist
	}

	if err != nil {
		return nil, e.Wrap(fmt.Sprintf("[ERROR] can't get user from storage username: %s, chat id: %d", username, chatID), err)
	}

	// log.Printf("from storage get user: tg id = %d, chat id = %d, dick size = %d", user.TgID, user.ChatID, user.DickSize)

	return user, nil
}

// UsersByChat возвращает всех пользователей из базы данных, которые находятся в одном телеграм чате.
func (s *Storage) UsersByChat(ctx context.Context, chatID int) ([]*storage.DBUser, error) {
	q := `SELECT * FROM users WHERE chat_id = $1 ORDER BY -dick_size`

	users := []*storage.DBUser{}

	err := s.db.Select(&users, q, chatID)
	if err != nil {
		return nil, errors.Wrap(err, "can't get users by chatID")
	}

	return users, nil
}

// GetGayOfDay возвращает запись о пидоре дне из базы данных.
func (s *Storage) GetGayOfDay(ctx context.Context, chatID int) (*storage.DBGay, error) {
	q := `SELECT * FROM gays WHERE chat_id = $1`

	gay := &storage.DBGay{}

	err := s.db.QueryRowContext(ctx, q, chatID).Scan(&gay.ID, &gay.ChatID, &gay.TgID, &gay.Username, &gay.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, storage.ErrUserNotExist
	}

	if err != nil {
		return nil, e.Wrap(fmt.Sprintf("[ERROR] can't get gay from table gays chat id: %d", chatID), err)
	}

	return gay, nil
}

// CreateGayOfDay Создаёт запись о пидоре дне в базе данных.
func (s *Storage) CreateGayOfDay(ctx context.Context, gay *storage.DBGay) error {
	q := `INSERT INTO gays (chat_id, tg_id, username, created_at) 
							VALUES ($1, $2, $3, $4)`

	log.Printf("[INFO] create gay of day #%d '%s', chat_id = %d", gay.TgID, gay.Username, gay.ChatID)

	if _, err := s.db.ExecContext(ctx, q, gay.ChatID, gay.TgID, gay.Username, gay.CreatedAt); err != nil {
		return e.Wrap(fmt.Sprintf("can't create gay %d %s: ", gay.TgID, gay.Username), err)
	}
	return nil
}

// RemoveGayOfDay удаляет запись о пидоре дня из базы данных.
func (s *Storage) RemoveGayOfDay(ctx context.Context, chatID int) error {
	q := `DELETE FROM gays WHERE chat_id = $1`

	if _, err := s.db.ExecContext(ctx, q, chatID); err != nil {
		return e.Wrap(fmt.Sprintf("[ERROR] can't remove gay %d %s: ", chatID), err)
	}
	return nil
}

// GetCalendarID возвращает Google Calendar ID из базы данных.
func (s *Storage) GetCalendarID(ctx context.Context, chatID int) (string, error) {
	q := `SELECT * from calendars WHERE chat_id = $1`

	var id int
	var calendarID string
	err := s.db.QueryRowContext(ctx, q, chatID).Scan(&id, &calendarID)
	if err == sql.ErrNoRows {
		return "", storage.ErrUserNotExist
	}

	if err != nil {
		return "", e.Wrap(fmt.Sprintf("[ERROR] can't get calendar_id from table calendars chat id: %d", chatID), err)
	}
	return calendarID, nil
}

// AddCalendarID добавляет Google Calendar ID в базу данных.
func (s *Storage) AddCalendarID(ctx context.Context, chatID int, calendarID string) error {
	q := `INSERT INTO calendars (chat_id, calendar_id) VALUES ($1, $2) ON CONFLICT (chat_id) DO UPDATE SET calendar_id = $3`
	if _, err := s.db.ExecContext(ctx, q, chatID, calendarID, calendarID); err != nil {
		return e.Wrap(fmt.Sprintf("can't update or create calendar_id in chat #%d: ", chatID), err)
	}
	return nil
}

// AddHomework добавляет запись домашнего задания в таблицу базы данных.
func (s *Storage) AddHomework(ctx context.Context, chatID int, subject string, task string) error {
	q := `INSERT INTO homeworks (chat_id, subject, task, created_at) VALUES ($1, $2, $3, $4)`
	if _, err := s.db.ExecContext(ctx, q, chatID, subject, task, time.Now()); err != nil {
		return e.Wrap("can't add homework:", err)
	}
	return nil
}

// GetHomeworkByChatID возвращает запись домашнего задания по id в таблице.
func (s *Storage) GetHomeworkByChatID(ctx context.Context, chatID int, limit int) ([]*storage.DBHomework, error) {
	q := `SELECT *  from homeworks WHERE chat_id = $1 ORDER BY created_at DESC LIMIT $2`

	homeworks := []*storage.DBHomework{}
	err := s.db.Select(&homeworks, q, chatID, limit)
	if err != nil {
		return nil, e.Wrap("can't get all homeworks", err)
	}
	return homeworks, nil
}

// GetHomeworkBySubject возвращает запись домашнего задания по названию предмета.
func (s *Storage) GetHomeworkBySubject(ctx context.Context, chatID int, subject string) ([]*storage.DBHomework, error) {
	q := `SELECT * from homeworks WHERE chat_id = $1 AND subject = $2 ORDER BY created_at DESC `

	homeworks := []*storage.DBHomework{}
	err := s.db.Select(&homeworks, q, chatID, subject)
	if err != nil {
		return nil, e.Wrap("can't get all homeworks", err)
	}
	return homeworks, nil
}

// DeleteHomework удаляет домашнее задание из базы данных.
func (s *Storage) DeleteHomework(ctx context.Context, id int) error {
	q := `DELETE FROM homeworks WHERE id = $1`
	_, err := s.db.ExecContext(ctx, q, id)
	if err != nil {
		return e.Wrap("can't delete row:", err)
	}
	return nil
}

// CreateUserStats создаёт статистику пользователя в базе данных.
func (s *Storage) CreateUserStats(ctx context.Context, u *storage.DBUserStat) (int, error) {
	q := `INSERT INTO user_stats (message_count, dick_plus_count, dick_minus_count, yes_count, no_count, duels_count, 
                        duels_win_count, duels_lose_count, kill_count, die_count, gay_count) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id`
	var id int
	row := s.db.QueryRow(q, u.MessageCount, u.DickPlusCount, u.DickMinusCount, u.YesCount, u.NoCount, u.DuelsCount, u.DuelsWinCount,
		u.DuelsLoseCount, u.KillCount, u.DieCount, u.GayCount)
	if err := row.Scan(&id); err != nil {
		return 0, e.Wrap("can't create user stats", err)
	}
	return id, nil
}

// GetUserStats возвращает ститистку пользователя из базы данных.
func (s *Storage) GetUserStats(ctx context.Context, u *storage.DBUser) (*storage.DBUserStat, error) {
	q := `SELECT * from user_stats WHERE id = $1`

	userStats := storage.DBUserStat{}
	err := s.db.Get(&userStats, q, u.UserStatId)
	if err != nil {
		return nil, e.Wrap("can't get user stats", err)
	}
	return &userStats, nil
}

// UpdateUserStats обновляет статистику пользователя в чате.
func (s *Storage) UpdateUserStats(ctx context.Context, u *storage.DBUserStat) error {
	q := `UPDATE user_stats SET message_count = $1, dick_plus_count = $2, dick_minus_count = $3, yes_count = $4, no_count = $5, 
                      duels_count = $6, duels_win_count = $7, duels_lose_count = $8, gay_count = $9 WHERE id = $10`
	_, err := s.db.ExecContext(ctx, q, u.MessageCount, u.DickPlusCount, u.DickMinusCount, u.YesCount, u.NoCount, u.DuelsCount, u.DuelsWinCount,
		u.DuelsLoseCount, u.GayCount, u.ID)

	if err != nil {
		return e.Wrap("can't update user stats", err)
	}

	return nil
}

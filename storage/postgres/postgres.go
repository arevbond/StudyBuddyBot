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

// CreateUser new user by chatID and telegramID.
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

// User get user by chatID and telegram ID.
func (s *Storage) UserByTelegramID(ctx context.Context, tgID, chatID int) (*storage.DBUser, error) {
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

func (s *Storage) UpdateUserDickSize(ctx context.Context, u *storage.DBUser, dickSize int) error {
	q := `UPDATE users SET dick_size = $1 WHERE tg_id = $2 AND chat_id = $3`
	oldDickSize := u.DickSize
	if _, err := s.db.ExecContext(ctx, q, dickSize, u.TgID, u.ChatID); err != nil {
		return e.Wrap(fmt.Sprintf("[ERROR] can't update dick size user %d chat id %d from %d to %d",
			u.TgID, u.ChatID, u.DickSize, dickSize), err)
	}
	u.DickSize = dickSize
	log.Printf("[INFO] user %d %s change his dick from %d to %d", u.TgID, u.Username, oldDickSize, u.DickSize)
	return nil
}

func (s *Storage) UpdateDateLastTryChangeDickToNow(ctx context.Context, u *storage.DBUser) error {
	q := `UPDATE users SET last_try_change_dick = $1 WHERE tg_id = $2 AND chat_id = $3`
	currentTime := time.Now()
	if _, err := s.db.ExecContext(ctx, q, currentTime, u.TgID, u.ChatID); err != nil {
		return e.Wrap(fmt.Sprintf("[ERROR] can't update date last try change dick to now user %d chat id",
			u.TgID, u.ChatID), err)
	}
	log.Printf("[INFO] user #%d %s change his date last try change dick to %s", u.TgID, u.Username, currentTime.Format("02-Jan-2006"))
	return nil
}

func (s *Storage) UsersByChat(ctx context.Context, chatID int) ([]*storage.DBUser, error) {
	q := `SELECT * FROM users WHERE chat_id = $1 ORDER BY -dick_size`

	users := []*storage.DBUser{}

	err := s.db.Select(&users, q, chatID)
	if err != nil {
		return nil, errors.Wrap(err, "can't get users by chatID")
	}

	return users, nil
}

func (s *Storage) GayOfDay(ctx context.Context, chatID int) (*storage.DBGay, error) {
	q := `SELECT * FROM gays WHERE chat_id = $1`

	gay := &storage.DBGay{}

	err := s.db.QueryRowContext(ctx, q, chatID).Scan(&gay.ChatID, &gay.TgID, &gay.Username, &gay.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, storage.ErrUserNotExist
	}

	if err != nil {
		return nil, e.Wrap(fmt.Sprintf("[ERROR] can't get gay from table gays chat id: %d", chatID), err)
	}

	// log.Printf("from storage get user: tg id = %d, chat id = %d, dick size = %d", user.TgID, user.ChatID, user.DickSize)

	return gay, nil
}

func (s *Storage) CreateGayOfDay(ctx context.Context, gay *storage.DBGay) error {
	q := `INSERT INTO gays (chat_id, tg_id, username, date_last_used) 
							VALUES ($1, $2, $3, $4)`

	log.Printf("[INFO] create gay of day #%d '%s', chat_id = %d", gay.TgID, gay.Username, gay.ChatID)

	if _, err := s.db.ExecContext(ctx, q, gay.ChatID, gay.TgID, gay.Username, gay.CreatedAt); err != nil {
		return e.Wrap(fmt.Sprintf("can't create gay %d %s: ", gay.TgID, gay.Username), err)
	}
	return nil
}

func (s *Storage) RemoveGayOfDay(ctx context.Context, chatID int) error {
	q := `DELETE FROM gays WHERE chat_id = $1`

	if _, err := s.db.ExecContext(ctx, q, chatID); err != nil {
		return e.Wrap(fmt.Sprintf("[ERROR] can't remove gay %d %s: ", chatID), err)
	}
	return nil
}

func (s *Storage) IncreaseCountOfGay(ctx context.Context, u *storage.DBUser) error {
	q1 := `SELECT * FROM user_stats WHERE id = $1`
	userStats := storage.DBUserStat{}
	s.db.Get(&userStats, q1, u.UserStatId)

	q2 := `UPDATE user_stats SET gay_count = $1 WHERE id = $2`

	if _, err := s.db.ExecContext(ctx, q2, userStats.GayCount+1, userStats.ID); err != nil {
		return errors.Wrap(err, "can't increase count of gay in user_stats")
	}
	return nil
}

func (s *Storage) CalendarID(ctx context.Context, chatID int) (string, error) {
	q := `SELECT * from calendars WHERE chat_id = $1`
	var calendarID string
	err := s.db.QueryRowContext(ctx, q, chatID).Scan(&calendarID)
	if err == sql.ErrNoRows {
		return "", storage.ErrUserNotExist
	}

	if err != nil {
		return "", e.Wrap(fmt.Sprintf("[ERROR] can't get calendar_id from table calendars chat id: %d", chatID), err)
	}
	return calendarID, nil
}

func (s *Storage) AddCalendarID(ctx context.Context, chatID int, calendarID string) error {
	q := `INSERT INTO calendars (chat_id, calendar_id) VALUES ($1, $2) ON CONFLICT (chat_id) DO UPDATE SET calendar_id = $3`
	if _, err := s.db.ExecContext(ctx, q, chatID, calendarID, calendarID); err != nil {
		return e.Wrap(fmt.Sprintf("can't update or create calendar_id in chat #%d: ", chatID), err)
	}
	return nil
}

func (s *Storage) AddHomework(ctx context.Context, chatID int, subject string, task string) error {
	q := `INSERT INTO homeworks (chat_id, subject, task, created_at) VALUES ($1, $2, $3, $4)`
	if _, err := s.db.ExecContext(ctx, q, chatID, subject, task, time.Now()); err != nil {
		return e.Wrap("can't add homework:", err)
	}
	return nil
}

func (s *Storage) GetHomeworkByChatID(ctx context.Context, chatID int, limit int) ([]*storage.DBHomework, error) {
	q := `SELECT rowid, *  from homeworks WHERE chat_id = $1 ORDER BY -created_at LIMIT $2`

	rows, err := s.db.QueryContext(ctx, q, chatID, limit)
	if err != nil {
		return nil, e.Wrap(fmt.Sprintf("can't get homeworks by chat id: %s", chatID), err)
	}
	defer rows.Close()

	var homeworks []*storage.DBHomework

	for rows.Next() {
		homework := &storage.DBHomework{}
		if err := rows.Scan(&homework.ID, &homework.ChatID, &homework.Subject, &homework.Task, &homework.CreatedAT); err != nil {
			return nil, e.Wrap(fmt.Sprintf("can't get homeworks by chat id: %s", chatID), err)
		}
		homeworks = append(homeworks, homework)
	}
	if err = rows.Err(); err != nil {
		return homeworks, err
	}
	return homeworks, nil
}

func (s *Storage) GetHomeworkBySubject(ctx context.Context, chatID int, subject string) ([]*storage.DBHomework, error) {
	q := `SELECT rowid, * from homeworks WHERE subject = $1 ORDER BY -created_at`

	rows, err := s.db.QueryContext(ctx, q, subject)
	if err != nil {
		return nil, e.Wrap(fmt.Sprintf("can't get homeworks by chat id: %s", chatID), err)
	}
	defer rows.Close()

	var homeworks []*storage.DBHomework

	for rows.Next() {
		homework := &storage.DBHomework{}
		if err := rows.Scan(&homework.ID, &homework.ChatID, &homework.Subject, &homework.Task, &homework.CreatedAT); err != nil {
			return nil, e.Wrap(fmt.Sprintf("can't get homeworks by subject: %d", chatID), err)
		}
		homeworks = append(homeworks, homework)
	}
	if err = rows.Err(); err != nil {
		return homeworks, err
	}
	return homeworks, nil
}

func (s *Storage) DeleteHomeworkByRowID(ctx context.Context, rowID int) error {
	q := `DELETE FROM homeworks WHERE rowid = $1`
	_, err := s.db.ExecContext(ctx, q, rowID)
	if err != nil {
		return e.Wrap("can't delete row:", err)
	}
	return nil
}

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

func (s *Storage) GetUserStats(ctx context.Context, u *storage.DBUser) (*storage.DBUserStat, error) {
	q := `SELECT * from user_stats WHERE id = $1`

	userStats := storage.DBUserStat{}
	err := s.db.Get(&userStats, q, u.UserStatId)
	if err != nil {
		return nil, e.Wrap("can't get user stats", err)
	}
	return &userStats, nil
}

func (s *Storage) UserStatsByTelegramIDAndChatID(ctx context.Context, tgID, chatID int) (*storage.DBUserStat, error) {
	return nil, nil
}

func (s *Storage) IncreaseMessageCount(ctx context.Context, u *storage.DBUserStat) error {
	return nil
}

func (s *Storage) IncreaseDickPlusCount(ctx context.Context, u *storage.DBUserStat) error {
	return nil
}

func (s *Storage) IncreaseDickMinusCount(ctx context.Context, u *storage.DBUserStat) error {
	return nil
}

func (s *Storage) IncreaseYesCount(ctx context.Context, u *storage.DBUserStat) error {
	return nil
}

func (s *Storage) IncreaseNoCount(ctx context.Context, u *storage.DBUserStat) error {
	return nil
}

func (s *Storage) UsersStatsByChatID(ctx context.Context, chatID int) ([]*storage.DBUserStat, error) {
	return nil, nil
}

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

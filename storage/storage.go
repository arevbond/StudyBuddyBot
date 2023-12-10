package storage

import (
	"context"
	"errors"
	"time"
)

type Storage interface {
	GetUser(ctx context.Context, tgID, chatID int) (*DBUser, error)
	CreateUser(ctx context.Context, u *DBUser) error
	UpdateUser(ctx context.Context, u *DBUser) error
	UsersByChat(ctx context.Context, chatID int) ([]*DBUser, error)

	UserByUsername(ctx context.Context, username string, chatID int) (*DBUser, error)

	GayOfDay(ctx context.Context, chatID int) (*DBGay, error)
	CreateGayOfDay(ctx context.Context, gay *DBGay) error
	RemoveGayOfDay(ctx context.Context, chatID int) error

	CalendarID(ctx context.Context, chatID int) (string, error)
	AddCalendarID(ctx context.Context, chatID int, calendarID string) error

	AddHomework(ctx context.Context, chatID int, subject string, task string) error
	GetHomeworkByChatID(ctx context.Context, chatID int, limit int) ([]*DBHomework, error)
	GetHomeworkBySubject(ctx context.Context, chatID int, subject string) ([]*DBHomework, error)
	DeleteHomeworkByRowID(ctx context.Context, rowID int) error

	CreateUserStats(ctx context.Context, u *DBUserStat) (int, error)
	GetUserStats(ctx context.Context, u *DBUser) (*DBUserStat, error)
	UpdateUserStats(ctx context.Context, u *DBUserStat) error
}

var ErrUserNotExist = errors.New("user not exists")

type DBUser struct {
	ID                 int       `json:"-" db:"id"`
	TgID               int       `db:"tg_id"`
	ChatID             int       `db:"chat_id"`
	IsBot              bool      `db:"is_bot"`
	IsPremium          bool      `db:"is_premium"`
	FirstName          string    `db:"first_name"`
	LastName           string    `db:"last_name"`
	Username           string    `db:"username"`
	DickSize           int       `db:"dick_size"`
	ChangeDickAt       time.Time `db:"change_dick_at"`
	UserStatId         int       `db:"user_stat_id"`
	HealthPoints       int       `db:"health_points"`
	HpTakedAt          time.Time `db:"hp_taked_at"`
	IsGay              bool      `db:"is_gay"`
	GayAt              time.Time `db:"gay_at"`
	Points             int       `db:"points"`
	CurDickChangeCount int       `db:"cur_dick_change_count"`
	MaxDickChangeCount int       `db:"max_dick_change_count"`
}

type DBGay struct {
	ID        int       `db:"id"`
	ChatID    int       `db:"chat_id"`
	TgID      int       `db:"tg_id"`
	Username  string    `db:"username"`
	CreatedAt time.Time `db:"created_at"`
}

type DBHomework struct {
	ID        int       `db:"id"`
	ChatID    int       `db:"chat_id"`
	Subject   string    `db:"subject"`
	Task      string    `db:"task"`
	CreatedAT time.Time `db:"created_at"`
}

type DBUserStat struct {
	ID             int `db:"id"`
	MessageCount   int `db:"message_count"`
	DickPlusCount  int `db:"dick_plus_count"`
	DickMinusCount int `db:"dick_minus_count"`
	YesCount       int `db:"yes_count"`
	NoCount        int `db:"no_count"`
	DuelsCount     int `db:"duels_count"`
	DuelsWinCount  int `db:"duels_win_count"`
	DuelsLoseCount int `db:"duels_lose_count"`
	KillCount      int `db:"kill_count"`
	DieCount       int `db:"die_count"`
	GayCount       int `db:"gay_count"`
}

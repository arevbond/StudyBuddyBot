package storage

import (
	"context"
	"errors"
	"time"
)

type Storage interface {
	CreateUser(ctx context.Context, u *DBUser) error
	UserByTelegramID(ctx context.Context, tgID, chatID int) (*DBUser, error)
	UserByUsername(ctx context.Context, username string, chatID int) (*DBUser, error)
	UpdateUserDickSize(ctx context.Context, u *DBUser, dickSize int) error
	UpdateDateLastTryChangeDickToNow(ctx context.Context, u *DBUser) error
	UsersByChat(ctx context.Context, chatID int) ([]*DBUser, error)
	IncreaseCountOfGay(ctx context.Context, u *DBUser) error

	GayOfDay(ctx context.Context, chatID int) (*DBGayOfDay, error)
	RemoveGayOfDay(ctx context.Context, chatID int) error
	CreateGayOfDay(ctx context.Context, gay *DBGayOfDay) error

	CalendarID(ctx context.Context, chatID int) (string, error)
	AddCalendarID(ctx context.Context, chatID int, calendarID string) error

	AddHomework(ctx context.Context, chatID int, subject string, task string) error
	GetHomeworkByChatID(ctx context.Context, chatID int, limit int) ([]*DBHomework, error)
	GetHomeworkBySubject(ctx context.Context, chatID int, subject string) ([]*DBHomework, error)
	DeleteHomeworkByRowID(ctx context.Context, rowID int) error

	CreateUserStats(ctx context.Context, u *DBUserStats) error
	UsersStatsByChatID(ctx context.Context, chatID int) ([]*DBUserStats, error)
	UserStatsByTelegramIDAndChatID(ctx context.Context, tgID, chatID int) (*DBUserStats, error)
	IncreaseMessageCount(ctx context.Context, u *DBUserStats) error
	IncreaseDickPlusCount(ctx context.Context, u *DBUserStats) error
	IncreaseDickMinusCount(ctx context.Context, u *DBUserStats) error
	IncreaseYesCount(ctx context.Context, u *DBUserStats) error
	IncreaseNoCount(ctx context.Context, u *DBUserStats) error
}

var ErrUserNotExist = errors.New("user not exists")

type DBUser struct {
	TgID           int
	ChatID         int
	IsBot          bool
	FirstName      string
	LastName       string
	Username       string
	IsPremium      bool
	DickSize       int
	CountGayOfDay  int
	DateChangeDick time.Time
}

type DBGayOfDay struct {
	ChatID       int
	TgID         int
	Username     string
	DateLastUsed time.Time
}

type DBHomework struct {
	ID            int
	ChatID        int
	Subject, Task string
	CreatedAT     time.Time
}

type DBUserStats struct {
	TelegramID     int
	ChatID         int
	UserName       string
	FirstName      string
	LastName       string
	MessageCount   int
	DickPlusCount  int
	DickMinusCount int
	YesCount       int
	NoCount        int
}

func NewDBUserStats(telegramID, chatID int, username string, firstName, lastName string) *DBUserStats {
	return &DBUserStats{telegramID, chatID, username, firstName, lastName, 0, 0, 0, 0, 0}
}

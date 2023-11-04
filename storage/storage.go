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

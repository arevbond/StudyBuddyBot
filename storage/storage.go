package storage

import (
	"context"
	"errors"
	"time"
)

type Storage interface {
	User(ctx context.Context, tgID, chatID int) (*DBUser, error)
	UserByUsername(ctx context.Context, username string, chatID int) (*DBUser, error)
	CreateUser(ctx context.Context, u *DBUser) error
	UpdateUserDickSize(ctx context.Context, u *DBUser, dickSize int) error
	UsersByChat(ctx context.Context, chatID int) ([]*DBUser, error)

	GayOfDay(ctx context.Context, chatID int) (*DBGayOfDay, error)
	CreateGayOfDay(ctx context.Context, gay *DBGayOfDay) error
}

var ErrUserNotExist = errors.New("user not exists")

type DBUser struct {
	TgID              int
	ChatID            int
	IsBot             bool
	FirstName         string
	LastName          string
	Username          string
	IsPremium         bool
	DickSize          int
	CountGayOfDay     int
	LastTryChangeDick time.Time
}

type DBGayOfDay struct {
	ChatID       int
	TgID         int
	Username     string
	DateLastUsed time.Time
}

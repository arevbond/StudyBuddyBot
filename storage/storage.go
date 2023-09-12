package storage

import (
	"context"
	"errors"
	"time"
)

type Storage interface {
	User(ctx context.Context, tgID, chatID int) (*DBUser, error)
	CreateUser(ctx context.Context, u *DBUser) error
	UpdateUserDickSize(ctx context.Context, u *DBUser, dickSize int) error
	UsersByChat(ctx context.Context, chatID int) ([]*DBUser, error)
}

var ErrUserNotExist = errors.New("user not exists")

type DBUser struct {
	TgID              int
	ChatID            int
	IsBot             bool   `json:"is_bot"`
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	Username          string `json:"username"`
	IsPremium         bool   `json:"is_premium"`
	DickSize          int
	LastTryChangeDick time.Time
}

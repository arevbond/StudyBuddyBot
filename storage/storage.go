package storage

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"tg_ics_useful_bot/lib/e"
	"time"
)

type Storage interface {
	Save(ctx context.Context, p *Page) error
	PickRandom(ctx context.Context, userName string) (*Page, error)
	Remove(ctx context.Context, p *Page) error
	IsExists(ctx context.Context, p *Page) (bool, error)

	User(ctx context.Context, tgID, chatID int) (*DBUser, error)
	CreateUser(ctx context.Context, u *DBUser) error
	UpdateUserDickSize(ctx context.Context, u *DBUser, dickSize int) error
}

var ErrNoSavedPages = errors.New("no saved pages")
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

type Page struct {
	URL      string
	UserName string
}

func (p Page) Hash() (string, error) {
	h := sha1.New()

	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", e.Wrap("can't calculate hash", err)
	}

	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", e.Wrap("can't calculate hash", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

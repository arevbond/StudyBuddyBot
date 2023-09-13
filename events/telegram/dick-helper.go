package telegram

import (
	"context"
	"fmt"
	"tg_ics_useful_bot/clients/telegram"
	"tg_ics_useful_bot/lib/e"
	"tg_ics_useful_bot/lib/game"
	"tg_ics_useful_bot/storage"
	"time"
)

func (p *Processor) createNewPlayer(chatID int, user *telegram.User) (*storage.DBUser, error) {
	dbUser := &storage.DBUser{
		TgID:              user.ID,
		ChatID:            chatID,
		IsBot:             user.IsBot,
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		Username:          user.Username,
		IsPremium:         user.IsPremium,
		DickSize:          game.PositiveRandomValue(),
		LastTryChangeDick: time.Now(),
	}
	err := p.storage.CreateUser(context.Background(), dbUser)
	if err != nil {
		return nil, err
	}
	return dbUser, err
}

func (p *Processor) changeDickSize(user *storage.DBUser, value int) (int, error) {
	oldDickSize := user.DickSize

	err := p.storage.UpdateUserDickSize(context.Background(), user, user.DickSize+value)
	if err != nil {
		return 0, e.Wrap(fmt.Sprintf("chat id %d, user %s can't change dick size: ", user.ChatID, user.Username), err)
	}
	return oldDickSize, nil
}

func (p *Processor) changeRandomDickSize(user *storage.DBUser) (bool, int, error) {
	value := game.RandomValue()
	oldDickSize := user.DickSize

	err := p.storage.UpdateUserDickSize(context.Background(), user, user.DickSize+value)
	if err != nil {
		return false, 0, e.Wrap(fmt.Sprintf("chat id %d, user %s can't change dick size: ", user.ChatID, user.Username), err)
	}
	return value >= 0, oldDickSize, nil
}

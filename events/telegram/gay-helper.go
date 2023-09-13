package telegram

import (
	"context"
	"math/rand"
	"tg_ics_useful_bot/clients/telegram"
	"tg_ics_useful_bot/storage"
	"time"
)

func (p *Processor) createNewGayOfDay(chatID int, admins []telegram.User) (*storage.DBGayOfDay, error) {
	rand.Seed(time.Now().Unix())
	n := rand.Intn(len(admins))
	u := admins[n]
	dbUser, err := p.storage.User(context.Background(), u.ID, chatID)
	if err == storage.ErrUserNotExist {
		dbUser, err = p.createNewPlayer(chatID, &u)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	gay := &storage.DBGayOfDay{
		ChatID:       chatID,
		TgID:         dbUser.TgID,
		Username:     dbUser.Username,
		DateLastUsed: time.Now(),
	}
	err = p.storage.CreateGayOfDay(context.Background(), gay)
	if err != nil {
		return nil, err
	}
	err = p.storage.IncreaseCountOfGay(context.Background(), dbUser)
	if err != nil {
		return nil, err
	}
	return gay, nil
}

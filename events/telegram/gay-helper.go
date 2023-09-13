package telegram

import (
	"context"
	"fmt"
	"math/rand"
	"tg_ics_useful_bot/clients/telegram"
	"tg_ics_useful_bot/storage"
	"time"
)

func (p *Processor) createNewGayOfDay(chatID int, admins []telegram.User) error {
	rand.Seed(time.Now().Unix())
	n := rand.Intn(len(admins))
	user := admins[n]
	gay := &storage.DBGayOfDay{
		ChatID:       chatID,
		TgID:         user.ID,
		Username:     user.Username,
		DateLastUsed: time.Now(),
	}
	err := p.storage.CreateGayOfDay(context.Background(), gay)
	if err != nil {
		return err
	}
	return p.tg.SendMessage(chatID, fmt.Sprintf(msgNewGayOfDay, gay.Username))
}

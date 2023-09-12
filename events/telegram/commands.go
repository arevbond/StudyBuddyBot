package telegram

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"tg_ics_useful_bot/clients/telegram"
	"tg_ics_useful_bot/lib/e"
	"tg_ics_useful_bot/storage"
	"time"
)

const (
	DickCmd  = "/dick"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(text string, chat *telegram.Chat, user *telegram.User) error {
	text = strings.TrimSpace(text)

	switch {
	case strings.HasPrefix(text, DickCmd):
		log.Printf("got new command '%s' from '%s", text, user.Username)
		return p.gameDick(chat, user)
	case strings.HasPrefix(text, HelpCmd):
		log.Printf("got new command '%s' from '%s", text, user.Username)
		return p.sendHelp(chat.ID)
	case strings.HasPrefix(text, StartCmd):
		log.Printf("got new command '%s' from '%s", text, user.Username)
		return p.sendHello(chat.ID)
	default:
		return nil
	}

}

func (p *Processor) gameDick(chat *telegram.Chat, user *telegram.User) (err error) {
	defer func() { err = e.WrapIfErr("can't change dick size: ", err) }()

	dbUser, err := p.storage.User(context.Background(), user.ID, chat.ID)

	if err == storage.ErrUserNotExist {
		dbUser = &storage.DBUser{
			TgID:              user.ID,
			ChatID:            chat.ID,
			IsBot:             user.IsBot,
			FirstName:         user.FirstName,
			LastName:          user.LastName,
			Username:          user.Username,
			IsPremium:         user.IsPremium,
			DickSize:          rand.Intn(25),
			LastTryChangeDick: time.Now(),
		}
		err = p.storage.CreateUser(context.Background(), dbUser)
		if err != nil {
			return err
		}
		return p.tg.SendMessage(chat.ID, fmt.Sprintf(msgCreateUser, dbUser.Username)+fmt.Sprintf(msgDickSize, dbUser.DickSize))
	} else if err != nil {
		return err
	}

	// TODO: if canChangeDickSize()
	isPlus, value, err := p.changeDickSize(dbUser)
	if err != nil {
		return err
	}
	if isPlus {
		return p.tg.SendMessage(chat.ID, fmt.Sprintf(msgDickIncrease, dbUser.Username, value)+fmt.Sprintf(msgDickSize, dbUser.DickSize))
	} else {
		return p.tg.SendMessage(chat.ID, fmt.Sprintf(msgDickDecrease, dbUser.Username, value)+fmt.Sprintf(msgDickSize, dbUser.DickSize))
	}
}

func (p *Processor) changeDickSize(user *storage.DBUser) (bool, int, error) {
	value := randomValue()

	log.Printf("%d user old dick size = %d, new dick size = %d", user.TgID, user.DickSize, user.DickSize+value)

	err := p.storage.UpdateUserDickSize(context.Background(), user, user.DickSize+value)
	if err != nil {
		return false, 0, e.Wrap(fmt.Sprintf("chat id %d, user %s can't change dick size: ", user.ChatID, user.Username), err)
	}
	return value >= 0, value, nil
}

func randomValue() int {
	sign := rand.Intn(5)
	value := rand.Intn(10)
	if sign > 0 {
		return value
	}
	return -1 * value
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

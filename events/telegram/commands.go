package telegram

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"strings"
	"tg_ics_useful_bot/lib/e"
	"tg_ics_useful_bot/storage"
	"time"
)

const (
	DickCmd = "/dick"

	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(text string, chatID int, userID int, username string, firstName, lastName string,
	isBot, isPremium bool) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s", text, username)

	if isAddCmd(text) {
		return p.savePage(chatID, text, username)
	}
	switch text {
	case DickCmd:
		return p.dick(chatID, userID, username, firstName, lastName, isBot, isPremium)
	case RndCmd:
		return p.sendRandom(chatID, username)
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHello(chatID)
	default:
		return nil
	}

}

func (p *Processor) dick(chatID int, userID int, username string, firstName, LastName string,
	isBot, isPremium bool) (err error) {
	defer func() { err = e.WrapIfErr("can't change dick size: ", err) }()

	user, err := p.storage.User(context.Background(), userID, chatID)
	if err == storage.ErrUserNotExist {
		user = &storage.User{
			TgID:              userID,
			ChatID:            chatID,
			IsBot:             isBot,
			FirstName:         firstName,
			LastName:          LastName,
			Username:          username,
			IsPremium:         isPremium,
			DickSize:          rand.Intn(25),
			LastTryChangeDick: time.Now(),
		}
		err = p.storage.CreateUser(context.Background(), user)
		if err != nil {
			return err
		}
		return p.tg.SendMessage(chatID,
			fmt.Sprintf(msgCreateUser+msgDickSize, username, user.DickSize))
	} else if err != nil {
		return err
	}
	isPlus, value, err := p.changeDickSize(user)
	if err != nil {
		return err
	}
	if isPlus {
		return p.tg.SendMessage(chatID, fmt.Sprintf(msgDickIncrease, username, value)+
			fmt.Sprintf(msgDickSize, user.DickSize))
	} else {
		return p.tg.SendMessage(chatID, fmt.Sprintf(msgDickDecrease, username, value)+
			fmt.Sprintf(msgDickSize, user.DickSize))
	}

	// if userCanChangeDick -> change his dick size -> sendMessage his dickSize success change
	// else -> sendMessage your dick can't change
}
func (p *Processor) changeDickSize(user *storage.User) (bool, int, error) {
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

func (p *Processor) savePage(chatID int, pageURL string, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: save page", err) }()

	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}

	isExists, err := p.storage.IsExists(context.Background(), page)
	if err != nil {
		return err
	}
	if isExists {
		return p.tg.SendMessage(chatID, msgAlreadyExists)
	}

	if err := p.storage.Save(context.Background(), page); err != nil {
		return err
	}

	if err := p.tg.SendMessage(chatID, msgSaved); err != nil {
		return err
	}

	return nil
}

func (p *Processor) sendRandom(chatID int, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't send random", err) }()

	page, err := p.storage.PickRandom(context.Background(), username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}
	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.tg.SendMessage(chatID, msgNoSavedPages)
	}

	if err := p.tg.SendMessage(chatID, page.URL); err != nil {
		return err
	}

	return p.storage.Remove(context.Background(), page)
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}

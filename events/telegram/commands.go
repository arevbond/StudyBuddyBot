package telegram

import (
	"context"
	"fmt"
	"log"
	"strings"
	"tg_ics_useful_bot/clients/telegram"
	"tg_ics_useful_bot/game"
	"tg_ics_useful_bot/lessons"
	"tg_ics_useful_bot/lib/e"
	"tg_ics_useful_bot/storage"
	"time"
)

const (
	GayStartCmd = "/gay"
	GayTopCmd   = "/gay_top"

	DicStartCmd = "/dick"
	DickTopCmd  = "/top_dick"

	TodayLessonsCmd    = "/today"
	LessonsCmd         = "/lessons"
	TomorrowLessonsCmd = "/tomorrow"

	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(text string, chat *telegram.Chat, user *telegram.User) error {
	text = strings.TrimSpace(text)

	switch {
	case strings.HasPrefix(text, TomorrowLessonsCmd):
		return p.tomorrowLessons(chat.ID)
	case strings.HasPrefix(text, LessonsCmd):
		return p.allLessons(chat.ID)
	case strings.HasPrefix(text, DicStartCmd):
		log.Printf("got new command '%s' from '%s", text, user.Username)
		return p.gameDick(chat, user)
	case strings.HasPrefix(text, DickTopCmd):
		return p.topDick(chat)
	case strings.HasPrefix(text, TodayLessonsCmd):
		return p.lessonsToday(chat.ID)
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

func (p *Processor) tomorrowLessons(chatID int) error {
	result := "Расписание на завтра:\n\n"
	result += lessons.StringTomorrowLessons(time.Now().Weekday())
	return p.tg.SendMessage(chatID, result)
}

func (p *Processor) allLessons(chatID int) error {
	result := "Расписание на неделю:\n\n"
	result += lessons.StringAllLessons()
	return p.tg.SendMessage(chatID, result)
}

func (p *Processor) lessonsToday(chatID int) error {
	result := "Расписание на сегодня:\n\n"
	result += lessons.StringLessonsByDay(time.Now().Weekday())
	return p.tg.SendMessage(chatID, result)
}

func (p *Processor) topDick(chat *telegram.Chat) (err error) {
	users, err := p.storage.UsersByChat(context.Background(), chat.ID)
	if err != nil {
		return err
	}
	result := ""
	for i, u := range users {
		result += fmt.Sprintf("%d. %s — %d см\n", i+1, u.FirstName+" "+u.LastName, u.DickSize)
	}
	return p.tg.SendMessage(chat.ID, result)
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
			DickSize:          game.RandomValue(),
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

	if game.CanChangeDickSize(dbUser) {
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
	return p.tg.SendMessage(chat.ID, fmt.Sprintf(msgAlreadyPlays, dbUser.Username))
}

func (p *Processor) changeDickSize(user *storage.DBUser) (bool, int, error) {
	value := game.RandomValue()

	log.Printf("%d user old dick size = %d, new dick size = %d", user.TgID, user.DickSize, user.DickSize+value)

	err := p.storage.UpdateUserDickSize(context.Background(), user, user.DickSize+value)
	if err != nil {
		return false, 0, e.Wrap(fmt.Sprintf("chat id %d, user %s can't change dick size: ", user.ChatID, user.Username), err)
	}
	return value >= 0, value, nil
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

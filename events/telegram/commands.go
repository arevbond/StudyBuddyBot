package telegram

import (
	"context"
	"fmt"
	"log"
	"strings"
	"tg_ics_useful_bot/clients/telegram"
	"tg_ics_useful_bot/clients/xkcd"
	"tg_ics_useful_bot/lessons"
	"tg_ics_useful_bot/lib/e"
	"tg_ics_useful_bot/lib/game"
	"tg_ics_useful_bot/lib/lib"
	"tg_ics_useful_bot/storage"
	"time"
)

const (
	GayStartCmd = "/gay"
	GayTopCmd   = "/gay_top"

	XkcdCmd = "/xkcd"

	DicStartCmd = "/dick"
	DickTopCmd  = "/top_dick"
	DickDuelCmd = "/duel"

	TodayLessonsCmd    = "/today"
	LessonsCmd         = "/lessons"
	TomorrowLessonsCmd = "/tomorrow"
)

func (p *Processor) doCmd(text string, chat *telegram.Chat, user *telegram.User, message *telegram.IncomingMessage) error {
	text = strings.TrimSpace(text)
	if strings.HasPrefix(text, "/") {
		log.Printf("got new command '%s' from '%s", text, user.Username)
	}

	switch {

	case strings.HasPrefix(text, DicStartCmd):
		return p.gameDick(chat, user, message)
	case strings.HasPrefix(text, DickTopCmd):
		return p.topDick(chat)
	case strings.HasPrefix(text, DickDuelCmd):
		if lib.Contains("@", text) {
			target := strings.Split(text, "@")[1]
			log.Printf("@%s вызывает на дуель @%s", user.Username, target)
			return p.duelDick(chat, user, target)
		}
		return nil

	case strings.HasPrefix(text, TodayLessonsCmd):
		return p.lessonsToday(chat.ID)
	case strings.HasPrefix(text, TomorrowLessonsCmd):
		return p.tomorrowLessons(chat.ID)
	case strings.HasPrefix(text, LessonsCmd):
		return p.allLessons(chat.ID)

	case strings.HasPrefix(text, XkcdCmd):
		return p.sendXkcd(chat.ID)
	default:
		return nil
	}

}

func (p *Processor) sendXkcd(chatID int) error {
	comics, err := xkcd.RandomComics()
	if err != nil {
		return err
	}
	return p.tg.SendPhoto(chatID, comics.Img)
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
	if err != nil {
		return err
	}
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

func (p *Processor) gameDick(chat *telegram.Chat, user *telegram.User, message *telegram.IncomingMessage) (err error) {
	defer func() { err = e.WrapIfErr("can't change dick size: ", err) }()

	err = p.tg.DeleteMessage(chat.ID, message.ID)
	if err != nil {
		return err
	}

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
			DickSize:          game.PositiveRandomValue(),
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
		_, oldDickSize, err := p.changeDickSize(dbUser)
		if err != nil {
			return err
		}
		return p.tg.SendMessage(chat.ID, fmt.Sprintf(msgChangeDickSize, dbUser.Username, oldDickSize, dbUser.DickSize))
	}
	return p.tg.SendMessage(chat.ID, fmt.Sprintf(msgAlreadyPlays, dbUser.Username))
}

func (p *Processor) duelDick(chat *telegram.Chat, user *telegram.User, targetUsername string) error {
	u1, err := p.storage.User(context.Background(), user.ID, chat.ID)
	if err != nil {
		return err
	}
	u2, err := p.storage.UserByUsername(context.Background(), targetUsername, chat.ID)
	if err == storage.ErrUserNotExist {
		return p.tg.SendMessage(chat.ID, fmt.Sprintf(msgTargetNotFound, targetUsername))
	}

	isUser1Win, ch1, ch2 := game.Duel(u1.DickSize, u2.DickSize)
	if isUser1Win {
		return p.tg.SendMessage(chat.ID, fmt.Sprintf(msgChanceDuel, u1.Username, u1.DickSize, ch1, targetUsername, u2.DickSize, ch2)+fmt.Sprintf(msgVictoryInDuel, u1.Username, u2.Username))
	} else {
		return p.tg.SendMessage(chat.ID, fmt.Sprintf(msgChanceDuel, u1.Username, u1.DickSize, ch1, targetUsername, u2.DickSize, ch2)+fmt.Sprintf(msgVictoryInDuel, u2.Username, u1.Username))
	}

}

func (p *Processor) changeDickSize(user *storage.DBUser) (bool, int, error) {
	value := game.RandomValue()
	oldDickSize := user.DickSize

	err := p.storage.UpdateUserDickSize(context.Background(), user, user.DickSize+value)
	if err != nil {
		return false, 0, e.Wrap(fmt.Sprintf("chat id %d, user %s can't change dick size: ", user.ChatID, user.Username), err)
	}
	return value >= 0, oldDickSize, nil
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

package telegram

import (
	"context"
	"fmt"
	"log"
	"sort"
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
	GayTopCmd   = "/top_gay"

	XkcdCmd = "/xkcd"

	DicStartCmd = "/dick"
	DickTopCmd  = "/top_dick"
	DickDuelCmd = "/duel"

	TodayLessonsCmd    = "/today"
	LessonsCmd         = "/lessons"
	TomorrowLessonsCmd = "/tomorrow"
)

func (p *Processor) doCmd(text string, chat *telegram.Chat, user *telegram.User, messageID int) error {
	text = strings.TrimSpace(text)
	if strings.HasPrefix(text, "/") {
		log.Printf("got new command '%s' from '%s", text, user.Username)
	}
	if chat.Type == "group" || chat.Type == "supergroup" {
		switch {
		case strings.HasPrefix(text, GayTopCmd):
			return p.gameGayTop(chat.ID)
		case strings.HasPrefix(text, GayStartCmd):
			return p.gameGay(chat.ID)

		case strings.HasPrefix(text, DickTopCmd):
			return p.topDick(chat)
		case strings.HasPrefix(text, DicStartCmd):
			return p.gameDick(chat, user, messageID)
		case strings.HasPrefix(text, DickDuelCmd):
			if lib.Contains("@", text) {
				textSplited := strings.Split(text, "@")
				target := textSplited[len(textSplited)-1]
				log.Printf("@%s вызывает на дуель @%s", user.Username, target)
				return p.duelDick(chat, user, target)
			}
			return p.duelDick(chat, user, user.Username)

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
	return nil
}

func (p *Processor) gameGayTop(chatID int) (err error) {
	admins, err := p.tg.ChatAdministrators(chatID)
	if err != nil {
		return e.Wrap("can't get chat administrators: ", err)
	}
	dbUsers := []*storage.DBUser{}
	for _, u := range admins {
		dbUser, err := p.storage.User(context.Background(), u.ID, chatID)
		if err == storage.ErrUserNotExist {
			dbUser = &storage.DBUser{
				TgID:      u.ID,
				ChatID:    chatID,
				IsBot:     u.IsBot,
				FirstName: u.FirstName,
				LastName:  u.LastName,
				Username:  u.Username,
			}
			err = p.storage.CreateUser(context.Background(), dbUser)
			if err != nil {
				return err
			}
		}
		dbUsers = append(dbUsers, dbUser)
	}
	sort.Slice(dbUsers, func(i, j int) bool {
		return dbUsers[i].CountGayOfDay >= dbUsers[j].CountGayOfDay
	})
	result := "Рейтинг пидоров: \n\n"

	for i, dbU := range dbUsers {
		result += fmt.Sprintf("%d. %s %s - %d раз \n", i+1, dbU.FirstName, dbU.LastName, dbU.CountGayOfDay)
	}
	return p.tg.SendMessage(chatID, result)
}

func (p *Processor) gameGay(chatID int) error {
	admins, err := p.tg.ChatAdministrators(chatID)
	if err != nil {
		return e.Wrap("can't get chat administrators: ", err)
	}

	gay, err := p.storage.GayOfDay(context.Background(), chatID)
	if err == storage.ErrUserNotExist {
		gay, err = p.createNewGayOfDay(chatID, admins)
		return p.tg.SendMessage(chatID, fmt.Sprintf(msgNewGayOfDay, gay.Username))
	} else if err != nil {
		return e.Wrap("can't get gay of day: ", err)
	}
	if gay.DateLastUsed.Month() >= time.Now().Month() && gay.DateLastUsed.Day() < time.Now().Day() {
		err = p.storage.RemoveGayOfDay(context.Background(), chatID)
		if err != nil {
			return err
		}
		gay, err = p.createNewGayOfDay(chatID, admins)
		return p.tg.SendMessage(chatID, fmt.Sprintf(msgNewGayOfDay, gay.Username))
	}
	return p.tg.SendMessage(chatID, fmt.Sprintf(msgCurrentGayOfDay, gay.Username))
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

func (p *Processor) gameDick(chat *telegram.Chat, user *telegram.User, messageID int) (err error) {
	defer func() { err = e.WrapIfErr("can't change dick size: ", err) }()

	err = p.tg.DeleteMessage(chat.ID, messageID)
	if err != nil {
		return err
	}

	dbUser, err := p.storage.User(context.Background(), user.ID, chat.ID)

	if err == storage.ErrUserNotExist {
		u, err := p.createNewPlayer(chat.ID, user)
		if err != nil {
			return err
		}
		return p.tg.SendMessage(chat.ID, fmt.Sprintf(msgCreateUser, u.Username)+fmt.Sprintf(msgDickSize, u.DickSize))
	} else if err != nil {
		return err
	}

	if game.CanChangeDickSize(dbUser) {
		_, oldDickSize, err := p.changeRandomDickSize(dbUser)
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
	} else if err != nil {
		return err
	}

	User1Win, ch1, ch2 := game.Duel(u1.DickSize, u2.DickSize)
	if User1Win {
		oldDickSize, err := p.changeDickSize(u1, game.PositiveRandomValue())
		if err != nil {
			return err
		}
		return p.tg.SendMessage(chat.ID,
			fmt.Sprintf(msgChanceDuel, u1.Username, oldDickSize, ch1, targetUsername, u2.DickSize, ch2)+
				fmt.Sprintf(msgVictoryInDuel, u1.Username, u2.Username)+
				fmt.Sprintf(msgDickSize, u1.DickSize))
	} else {
		_, err := p.changeDickSize(u1, -1*game.PositiveRandomValue())
		if err != nil {
			return err
		}
		return p.tg.SendMessage(chat.ID,
			fmt.Sprintf(msgChanceDuel, u1.Username, u1.DickSize, ch1, targetUsername, u2.DickSize, ch2)+
				fmt.Sprintf(msgVictoryInDuel, u2.Username, u1.Username))
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

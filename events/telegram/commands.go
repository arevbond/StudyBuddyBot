package telegram

import (
	"log"
	"strings"
	"tg_ics_useful_bot/clients/jokesrv"
	"tg_ics_useful_bot/clients/telegram"
	"tg_ics_useful_bot/clients/xkcd"
	"tg_ics_useful_bot/lessons"
	"tg_ics_useful_bot/lib/e"
	"tg_ics_useful_bot/lib/lib"
	"time"
)

const (
	AnecdotCmd = "/joke"

	FlipCmd = "/flip"

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
	if strings.HasPrefix(strings.ToLower(text), "да") {
		return p.tg.SendMessage(chat.ID, "Пизда")
	}

	if strings.HasPrefix(text, "/") {
		log.Printf("[INFO] got new command '%s' from '%s", text, user.Username)
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
				log.Printf("[INFO] @%s вызывает на дуель @%s", user.Username, target)
				return p.gameDuelDick(chat, messageID, user, target)
			}
			return p.gameDuelDick(chat, messageID, user, user.Username)

		case strings.HasPrefix(text, TodayLessonsCmd):
			return p.lessonsToday(chat.ID)
		case strings.HasPrefix(text, TomorrowLessonsCmd):
			return p.tomorrowLessons(chat.ID)
		case strings.HasPrefix(text, LessonsCmd):
			return p.allLessons(chat.ID)

		case strings.HasPrefix(text, XkcdCmd):
			return p.sendXkcd(chat.ID)

		case strings.HasPrefix(text, AnecdotCmd):
			anecdot, err := jokesrv.Anecdot()
			if err != nil {
				return e.Wrap("can't get anecdot: ", err)
			}
			return p.tg.SendMessage(chat.ID, anecdot)

		case strings.HasPrefix(text, FlipCmd):
			return p.tg.SendPhoto(chat.ID, RandomPhotoHinkOrRoom())

		default:
			return nil
		}
	}
	return nil
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

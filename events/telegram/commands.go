package telegram

import (
	"context"
	"fmt"
	"log"
	"strings"
	"tg_ics_useful_bot/clients/jokesrv"
	"tg_ics_useful_bot/clients/telegram"
	"tg_ics_useful_bot/clients/xkcd"
	"tg_ics_useful_bot/events/schedule"
	"tg_ics_useful_bot/lessons"
	"tg_ics_useful_bot/lib/e"
	"tg_ics_useful_bot/lib/utils"
)

type method int

const (
	UnsupportedMethod method = iota
	sendMessageMethod
	sendPhotoMethod
)

var (
	AllCmd             = []string{"/all", "/all@ics_useful_bot"}
	AnecdotCmd         = []string{"/joke", "/joke@ics_useful_bot"}
	FlipCmd            = []string{"/flip", "/flip@ics_useful_bot"}
	GayStartCmd        = []string{"/gay", "/gay@ics_useful_bot"}
	GayTopCmd          = []string{"/top_gay", "/top_gay@ics_useful_bot"}
	XkcdCmd            = []string{"/xkcd", "/xkcd@ics_useful_bot"}
	DicStartCmd        = []string{"/dick", "/dick@ics_useful_bot"}
	DickTopCmd         = []string{"/top_dick", "/top_dick@ics_useful_bot"}
	DickDuelCmd        = []string{"/duel", "/duel@ics_useful_bot"}
	TodayLessonsCmd    = []string{"/today", "/today@ics_useful_bot"}
	LessonsCmd         = []string{"/lessons", "/lessons@ics_useful_bot"}
	TomorrowLessonsCmd = []string{"/tomorrow", "/tomorrow@ics_useful_bot"}

	// test
	ScheduleCmd = []string{"/schedule"}
)

// selectCommand select one of available commands.
func (p *Processor) selectCommand(text string, chat *telegram.Chat, user *telegram.User, messageID int) (string, method, error) {
	var message string
	var mthd method
	var err error
	switch {

	case utils.Contains(text, AllCmd):
		message = p.allUsernames(chat.ID)
		mthd = sendMessageMethod

	case utils.Contains(text, GayTopCmd):
		message, err = p.gameGayTop(chat.ID)
		if err != nil {
			return "", UnsupportedMethod, e.Wrap("can't do GayTopCmd: ", err)
		}
		mthd = sendMessageMethod

	case utils.Contains(text, GayStartCmd):
		message, err = p.gameGay(chat.ID)
		if err != nil {
			return "", UnsupportedMethod, e.Wrap("can't get message from gameGay: ", err)
		}
		mthd = sendMessageMethod

	case utils.Contains(text, DickTopCmd):
		message, err = p.topDick(chat.ID)
		if err != nil {
			return "", UnsupportedMethod, e.Wrap(fmt.Sprintf("can't get top dics from chat %d: ", chat.ID), err)
		}
		mthd = sendMessageMethod

	case utils.Contains(text, DicStartCmd):
		message, err = p.gameDick(chat, user, messageID)
		if err != nil {
			return "", UnsupportedMethod, e.Wrap("can't get message from gameDick: ", err)
		}
		mthd = sendMessageMethod

	case strings.HasPrefix(text, DickDuelCmd[0]):
		message, err = p.gameDuelDick(chat, messageID, user, user.Username)
		if err != nil {
			return "", UnsupportedMethod, e.Wrap("can't do gameDuelDick: ", err)
		}
		if utils.StringContains("@", text) {
			textSplited := strings.Split(text, "@")
			target := textSplited[len(textSplited)-1]
			log.Printf("[INFO] @%s вызывает на дуель @%s", user.Username, target)
			message, err = p.gameDuelDick(chat, messageID, user, target)
			if err != nil {
				return "", UnsupportedMethod, e.Wrap("can't do gameDuelDick: ", err)
			}
		}
		mthd = sendMessageMethod

	case utils.Contains(text, TodayLessonsCmd):
		message = lessons.LessonsToday()
		mthd = sendMessageMethod
	case utils.Contains(text, TomorrowLessonsCmd):
		message = lessons.TomorrowLessons()
		mthd = sendMessageMethod
	case utils.Contains(text, LessonsCmd):
		message = lessons.AllLessons()
		mthd = sendMessageMethod

	case utils.Contains(text, XkcdCmd):
		var comics xkcd.Comics
		comics, err = xkcd.RandomComics()
		if err != nil {
			return "", UnsupportedMethod, e.Wrap("can't get comics from xkcd: ", err)
		}
		message = comics.Img
		mthd = sendPhotoMethod

	case utils.Contains(text, AnecdotCmd):
		message, err = jokesrv.Anecdot()
		if err != nil {
			return "", UnsupportedMethod, e.Wrap("can't get anecdot: ", err)
		}
		mthd = sendMessageMethod
	case utils.Contains(text, FlipCmd):
		message = RandomPhotoHinkOrRoom()
		mthd = sendPhotoMethod
	case utils.Contains(text, ScheduleCmd):
		calendarID, err := p.storage.CalendarID(context.Background(), chat.ID)
		if err != nil || calendarID == "" {
			return "", UnsupportedMethod, e.Wrap("can't get calendarID: ", err)
		}

		message = schedule.Schedule(calendarID)
		mthd = sendMessageMethod
	}
	return message, mthd, nil
}

func (p *Processor) allUsernames(chatID int) string {
	admins, err := p.tg.ChatAdministrators(chatID)
	if err != nil {
		log.Printf("can't get admins in chat #%d: ", chatID, err)
	}
	result := ""
	for _, a := range admins {
		result += "@" + a.Username + " "
	}
	return result[:len(result)-1]
}

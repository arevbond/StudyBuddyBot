package telegram

import (
	"context"
	"fmt"
	"log"
	"strings"
	"tg_ics_useful_bot/clients/jokesrv"
	"tg_ics_useful_bot/clients/telegram"
	"tg_ics_useful_bot/clients/xkcd"
	"tg_ics_useful_bot/lib/e"
	"tg_ics_useful_bot/lib/schedule"
	"tg_ics_useful_bot/lib/utils"
)

type method int

const (
	UnsupportedMethod method = iota
	sendMessageMethod
	sendPhotoMethod
)

var (
	AllCmd      = []string{"/all", "/all@ics_useful_bot"}
	AnecdotCmd  = []string{"/joke", "/joke@ics_useful_bot"}
	FlipCmd     = []string{"/flip", "/flip@ics_useful_bot"}
	GayStartCmd = []string{"/gay", "/gay@ics_useful_bot"}
	GayTopCmd   = []string{"/top_gay", "/top_gay@ics_useful_bot"}
	XkcdCmd     = []string{"/xkcd", "/xkcd@ics_useful_bot"}
	DicStartCmd = []string{"/dick", "/dick@ics_useful_bot"}
	DickTopCmd  = []string{"/top_dick", "/top_dick@ics_useful_bot"}
	DickDuelCmd = []string{"/duel", "/duel@ics_useful_bot"}
	ScheduleCmd = []string{"/schedule", "/schedule@ics_useful_bot"}

	AddCalendarIDCmd = []string{"/add_calendar", "/add_calendar@ics_useful_bot"}
)

// selectCommand select one of available commands.
func (p *Processor) selectCommand(text string, chat *telegram.Chat, user *telegram.User, messageID int) (string, method, error) {
	var message string
	var mthd method
	var err error
	switch {

	case utils.Equal(text, AllCmd):
		message = p.allUsernames(chat.ID)
		mthd = sendMessageMethod

	case utils.Equal(text, GayTopCmd):
		message, err = p.gameGayTop(chat.ID)
		if err != nil {
			return "", UnsupportedMethod, e.Wrap("can't do GayTopCmd: ", err)
		}
		mthd = sendMessageMethod

	case utils.Equal(text, GayStartCmd):
		message, err = p.gameGay(chat.ID)
		if err != nil {
			return "", UnsupportedMethod, e.Wrap("can't get message from gameGay: ", err)
		}
		mthd = sendMessageMethod

	case utils.Equal(text, DickTopCmd):
		message, err = p.topDick(chat.ID)
		if err != nil {
			return "", UnsupportedMethod, e.Wrap(fmt.Sprintf("can't get top dics from chat %d: ", chat.ID), err)
		}
		mthd = sendMessageMethod

	case utils.Equal(text, DicStartCmd):
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

	case utils.Equal(text, XkcdCmd):
		var comics xkcd.Comics
		comics, err = xkcd.RandomComics()
		if err != nil {
			return "", UnsupportedMethod, e.Wrap("can't get comics from xkcd: ", err)
		}
		message = comics.Img
		mthd = sendPhotoMethod

	case utils.Equal(text, AnecdotCmd):
		message, err = jokesrv.Anecdot()
		if err != nil {
			return "", UnsupportedMethod, e.Wrap("can't get anecdot: ", err)
		}
		mthd = sendMessageMethod
	case utils.Equal(text, FlipCmd):
		message = RandomPhotoHinkOrRoom()
		mthd = sendPhotoMethod
	case utils.Equal(text, ScheduleCmd):
		calendarID, err := p.storage.CalendarID(context.Background(), chat.ID)
		if err != nil || calendarID == "" {
			return "", UnsupportedMethod, e.Wrap("can't get calendarID: ", err)
		}

		message, err = schedule.Schedule(calendarID)
		if err != nil {
			log.Printf("[ERROR] can't send schedule: %v", err)
			message = fmt.Sprintf(msgErrorSendMessage, calendarID)
		}
		mthd = sendMessageMethod
	case utils.Equal(strings.Split(text, " ")[0], AddCalendarIDCmd):
		if !p.isAdmin(user, chat.ID) {
			return msgForbiddenCalendarUpdate, sendMessageMethod, nil
		}
		strs := strings.Split(text, " ")
		calendarID := ""
		for _, str := range strs {
			if len(str) > 0 {
				calendarID = str
			}
		}
		err = p.storage.AddCalendarID(context.Background(), chat.ID, calendarID)
		if err != nil {
			message = fmt.Sprintf(msgErrorUpdateCalendarID, calendarID)
			log.Printf("[ERROR] can't update calender_id: %v", err)
		} else {
			message = msgSuccessUpdateCalendarID
		}
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

func (p *Processor) isAdmin(user *telegram.User, chatID int) bool {
	admins, err := p.tg.ChatAdministrators(chatID)
	if err != nil {
		log.Printf("can't get admins in chat #%d: ", chatID, err)
	}
	for _, admin := range admins {
		if user.ID == admin.ID {
			return true
		}
	}
	return false
}

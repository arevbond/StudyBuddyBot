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

const (
	suffix = "@ics_useful_bot"
)

var (
	HelpCmd = "/help"

	DicStartCmd = "/dick"
	DickTopCmd  = "/top_dick"
	DickDuelCmd = "/duel"

	GayStartCmd = "/gay"
	GayTopCmd   = "/top_gay"

	AddCalendarIDCmd = "/add_calendar"
	ScheduleCmd      = "/schedule"

	AnecdotCmd = "/joke"
	XkcdCmd    = "/xkcd"
	FlipCmd    = "/flip"

	AllCmd = "/all"
)

// selectCommand select one of available commands.
func (p *Processor) selectCommand(cmd string, chat *telegram.Chat, user *telegram.User, messageID int) (string, method, error) {
	var message string
	var mthd method
	var err error
	switch {

	case isCommand(cmd, AllCmd):
		message = p.allUsernames(chat.ID)
		mthd = sendMessageMethod

	case isCommand(cmd, GayTopCmd):
		message, err = p.gameGayTop(chat.ID)
		if err != nil {
			return "", UnsupportedMethod, e.Wrap("can't do GayTopCmd: ", err)
		}
		mthd = sendMessageMethod

	case isCommand(cmd, GayStartCmd):
		message, err = p.gameGay(chat.ID)
		if err != nil {
			return "", UnsupportedMethod, e.Wrap("can't get message from gameGay: ", err)
		}
		mthd = sendMessageMethod

	case isCommand(cmd, DickTopCmd):
		message, err = p.topDick(chat.ID)
		if err != nil {
			return "", UnsupportedMethod, e.Wrap(fmt.Sprintf("can't get top dics from chat %d: ", chat.ID), err)
		}
		mthd = sendMessageMethod

	case isCommand(cmd, DicStartCmd):
		message, err = p.gameDick(chat, user, messageID)
		if err != nil {
			return "", UnsupportedMethod, e.Wrap("can't get message from gameDick: ", err)
		}
		mthd = sendMessageMethod

	case isCommand(strings.Split(cmd, " ")[0], DickDuelCmd) || isCommand(cmd, DickDuelCmd):
		message, err = p.gameDuelDick(chat, messageID, user, user.Username)
		if err != nil {
			return "", UnsupportedMethod, e.Wrap("can't do gameDuelDick: ", err)
		}
		if utils.StringContains("@", cmd) {
			textSplited := strings.Split(cmd, "@")
			target := textSplited[len(textSplited)-1]
			log.Printf("[INFO] @%s вызывает на дуель @%s", user.Username, target)
			message, err = p.gameDuelDick(chat, messageID, user, target)
			if err != nil {
				return "", UnsupportedMethod, e.Wrap("can't do gameDuelDick: ", err)
			}
		}
		mthd = sendMessageMethod

	case isCommand(cmd, XkcdCmd):
		var comics xkcd.Comics
		comics, err = xkcd.RandomComics()
		if err != nil {
			return "", UnsupportedMethod, e.Wrap("can't get comics from xkcd: ", err)
		}
		message = comics.Img
		mthd = sendPhotoMethod

	case isCommand(cmd, AnecdotCmd):
		message, err = jokesrv.Anecdot()
		if err != nil {
			return "", UnsupportedMethod, e.Wrap("can't get anecdot: ", err)
		}
		mthd = sendMessageMethod
	case isCommand(cmd, FlipCmd):
		message = RandomPhotoHinkOrRoom()
		mthd = sendPhotoMethod
	case isCommand(cmd, ScheduleCmd):
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
	case isCommand(strings.Split(cmd, " ")[0], AddCalendarIDCmd):
		if !p.isAdmin(user, chat.ID) {
			return msgForbiddenCalendarUpdate, sendMessageMethod, nil
		}
		strs := strings.Split(cmd, " ")
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
	case isCommand(cmd, HelpCmd):
		message = msgHelp
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

func isCommand(cmd string, correctCmd string) bool {
	if cmd == correctCmd || cmd == correctCmd+suffix {
		return true
	}
	return false
}

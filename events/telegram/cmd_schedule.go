package telegram

import (
	"context"
	"fmt"
	"log"
	"strings"
	"tg_ics_useful_bot/clients/telegram"
	"tg_ics_useful_bot/lib/schedule"
	"tg_ics_useful_bot/storage"
)

// addCalendarExec предоставляет Exec метод для выполнения /add_calendar.
type addCalendarExec struct {
	command string
}

// Exec: /add_calendar {calendar_id}
func (a *addCalendarExec) Exec(p *Processor, inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {

	if !p.isChatAdmin(user, chat.ID) {
		return &Response{message: msgForbiddenCalendarUpdate, method: sendMessageMethod}, nil
	}
	strs := strings.Split(inMessage, " ")
	calendarID := ""
	for _, str := range strs {
		if len(str) > 0 {
			calendarID = str
		}
	}
	err := p.storage.AddCalendarID(context.Background(), chat.ID, calendarID)
	var message string
	if err != nil {
		message = fmt.Sprintf(msgErrorUpdateCalendarID, calendarID)
		log.Printf("can't update calender_id: %v", err)
	} else {
		message = msgSuccessUpdateCalendarID
	}
	mthd := sendMessageMethod
	return &Response{message: message, method: mthd, replyMessageId: -1}, nil
}

// scheduleExec предоставляет Exec метод для выполнения /schedule.
type scheduleExec struct {
	command string
}

// Exec: /schedule - возвращает расписание из Google Calender.
func (a *scheduleExec) Exec(p *Processor, inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {
	var message string
	var parseMode telegram.ParseMode
	calendarID, err := p.storage.GetCalendarID(context.Background(), chat.ID)
	if err != nil || calendarID == "" {
		message = msgCalendarNotExists
		log.Print("can't get calendarID: ", err)
	} else {
		message, err = schedule.ScheduleCmd(calendarID)
		parseMode = telegram.Markdown
		if err != nil {
			log.Printf("[ERROR] can't send schedule: %v", err)
			message = fmt.Sprintf(msgErrorSendMessage, calendarID)
			parseMode = ""
		}
	}
	mthd := sendMessageMethod
	return &Response{message: message, method: mthd, replyMessageId: -1, parseMode: parseMode}, nil
}

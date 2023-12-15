package telegram

import (
	"fmt"
	"tg_ics_useful_bot/clients/telegram"
	"tg_ics_useful_bot/lib/e"
	"tg_ics_useful_bot/storage"
)

// dickTopExec предоставляет метод Exec для выполнения /top_dick.
type dickTopExec struct {
	command string
}

// Exec: /top_dick - пишет топ всех пенисов в чат.
func (a *dickTopExec) Exec(p *Processor, inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {
	message, err := p.topDicksCmd(chat.ID)
	if err != nil {
		return nil, e.Wrap(fmt.Sprintf("can't get top dics from chat %d: ", chat.ID), err)
	}
	mthd := sendMessageMethod
	return &Response{message: message, method: mthd, replyMessageId: -1}, nil
}

// dickStartExec предоставляет метод Exec для выполнения /dick.
type dickStartExec struct {
	command string
}

// Exec: /dick - игра в пенис.
func (a *dickStartExec) Exec(p *Processor, inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {
	message, err := p.gameDickCmd(chat, user, userStats)
	if err != nil {
		return nil, e.Wrap("can't get message from gameDickCmd: ", err)
	}
	mthd := sendMessageMethod
	return &Response{message: message, method: mthd, replyMessageId: -1}, nil
}

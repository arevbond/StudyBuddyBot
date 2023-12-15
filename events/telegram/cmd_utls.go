package telegram

import (
	"tg_ics_useful_bot/clients/telegram"
	"tg_ics_useful_bot/lib/e"
	"tg_ics_useful_bot/storage"
)

type helpCmd struct {
	command string
	p       *Processor
}

func (a *helpCmd) Exec(inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {
	message, err := a.p.getHp(user, chat)
	if err != nil {
		return nil, e.Wrap("can't get hp in 'selectCommand':", err)
	}
	mthd := sendMessageMethod
	return &Response{message: message, method: mthd, replyMessageId: messageID}, nil
}

func (a *helpCmd) SetProcessor(p *Processor) {
	a.p = p
}

type getChatIDcmd struct {
	command string
	p       *Processor
}

func (a *getChatIDcmd) Exec(inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {
	message, err := a.p.getHp(user, chat)
	if err != nil {
		return nil, e.Wrap("can't get hp in 'selectCommand':", err)
	}
	mthd := sendMessageMethod
	return &Response{message: message, method: mthd, replyMessageId: messageID}, nil
}

func (a *getChatIDcmd) SetProcessor(p *Processor) {
	a.p = p
}

package telegram

import (
	"tg_ics_useful_bot/clients/telegram"
	"tg_ics_useful_bot/lib/e"
	"tg_ics_useful_bot/storage"
)

// gayExec предоставляет метод Exec для выполнения /gay.
type gayExec struct {
	command string
}

// Exec: /gay - определяет случайного пидора в чате среди админов чата.
func (a *gayExec) Exec(p *Processor, inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {
	message, err := p.gameGay(chat.ID)
	if err != nil {
		return nil, e.Wrap("can't get message from gameGay: ", err)
	}
	mthd := sendMessageMethod
	return &Response{message: message, method: mthd, replyMessageId: messageID}, nil
}

// topGaysExec предоставляет метод Exec для вывода топа пидоров.
type topGaysExec struct {
	command string
}

// Exec: /top_gay - выводит список участников чата и их кол-во становления пидором дня.
func (a *topGaysExec) Exec(p *Processor, inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {
	message, err := p.topGays(chat.ID)
	if err != nil {
		return nil, e.Wrap("can't do GayTop: ", err)
	}
	mthd := sendMessageMethod
	return &Response{message: message, method: mthd, replyMessageId: messageID}, nil
}

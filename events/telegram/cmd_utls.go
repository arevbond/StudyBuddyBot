package telegram

import (
	"strconv"
	"tg_ics_useful_bot/clients/telegram"
	"tg_ics_useful_bot/lib/e"
	"tg_ics_useful_bot/storage"
)

// helpExec предоставляет метод Exec для выполнения /help.
type helpExec struct {
	command string
}

// Exec: /help - возвращает help сообщениею
func (a *helpExec) Exec(p *Processor, inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {
	message, err := p.getHp(user, chat)
	if err != nil {
		return nil, e.Wrap("can't get hp in 'selectCommand':", err)
	}
	mthd := sendMessageMethod
	return &Response{message: message, method: mthd, replyMessageId: messageID}, nil
}

// chatIDExec предоставляет метод Exec для выполнения /chat_id.
type chatIDExec struct {
	command string
}

// Exec: /chat_id - возвращает chat id.
func (a *chatIDExec) Exec(p *Processor, inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {

	message := strconv.Itoa(chat.ID)
	mthd := sendMessageMethod
	return &Response{message: message, method: mthd, replyMessageId: messageID}, nil
}

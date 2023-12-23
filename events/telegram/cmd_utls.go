package telegram

import (
	"strconv"
	"tg_ics_useful_bot/clients/telegram"
	"tg_ics_useful_bot/storage"
)

// helpExec предоставляет метод Exec для выполнения /help.
type helpExec string

// Exec: /help - возвращает help сообщениею
func (a helpExec) Exec(p *Processor, inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {

	message := msgHelp
	mthd := sendMessageMethod
	return &Response{message: message, method: mthd, replyMessageId: -1, parseMode: telegram.Markdown}, nil
}

// chatIDExec предоставляет метод Exec для выполнения /chat_id.
type chatIDExec string

// Exec: /chat_id - возвращает chat id.
func (a chatIDExec) Exec(p *Processor, inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {

	message := strconv.Itoa(chat.ID)
	mthd := sendMessageMethod
	return &Response{message: message, method: mthd, replyMessageId: -1}, nil
}

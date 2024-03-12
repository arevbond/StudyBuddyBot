package telegram

import (
	"log"
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

// allUsernamesExec предоставляет метод Exec для вызова всех админов в чате.
type allUsernamesExec string

// Exec: /all - тэгает всех админов в чате.
func (a allUsernamesExec) Exec(p *Processor, inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {

	message := a.allUsernames(chat.ID, p.tg)
	mthd := sendMessageMethod
	return &Response{message: message, method: mthd, replyMessageId: -1}, nil
}

// allUsernames возвращает строку "@username1, @username2...".
func (a allUsernamesExec) allUsernames(chatID int, tgClient *telegram.Client) string {
	admins, err := tgClient.ChatAdministrators(chatID)
	if err != nil {
		log.Printf("can't get admins in chat #%d: %v", chatID, err)
	}
	result := ""
	for _, a := range admins {
		result += "@" + a.Username + " "
	}
	return result[:len(result)-1]
}

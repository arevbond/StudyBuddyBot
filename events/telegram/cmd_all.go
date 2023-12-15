package telegram

import (
	"log"
	"tg_ics_useful_bot/clients/telegram"
	"tg_ics_useful_bot/storage"
)

// allUsernamesExec предоставляет метод Exec для вызова всех админов в чате.
type allUsernamesExec struct {
	command string
}

// Exec: /all - тэгает всех админов в чате.
func (a *allUsernamesExec) Exec(p *Processor, inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {

	message := p.allUsernames(chat.ID)
	mthd := sendMessageMethod
	return &Response{message: message, method: mthd, replyMessageId: -1}, nil
}

// allUsernames возвращает строку "@username1, @username2...".
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

// isChatAdmin определяет является ли пользователь админов в чате.
func (p *Processor) isChatAdmin(user *telegram.User, chatID int) bool {
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

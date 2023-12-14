package telegram

import (
	"log"
	"tg_ics_useful_bot/clients/telegram"
	"tg_ics_useful_bot/storage"
)

type allUsernamesCmd struct {
	command string
	p       *Processor
}

func (a *allUsernamesCmd) Exec(inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {
	message := a.p.allUsernames(chat.ID)
	mthd := sendMessageMethod
	return &Response{message: message, method: mthd, replyMessageId: messageID}, nil
}

func (a *allUsernamesCmd) SetProcessor(p *Processor) {
	a.p = p
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

func isCommand(cmd string, correctCmd string) bool {
	if cmd == correctCmd || cmd == correctCmd+suffix {
		return true
	}
	return false
}

package telegram

import (
	"fmt"
	"log"
	"strings"
	"tg_ics_useful_bot/clients/telegram"
	"tg_ics_useful_bot/lib/e"
	"tg_ics_useful_bot/lib/utils"
)

func (p *Processor) doCmd(text string, chat *telegram.Chat, user *telegram.User, messageID int) error {
	text = strings.TrimSpace(text)

	if utils.IsYesCommand(text) {
		return p.tg.SendMessage(chat.ID, "Пизда")
	} else if utils.IsNoCommand(text) {
		return p.tg.SendMessage(chat.ID, "Пидора ответ")
	}

	if utils.IsCommand(text) {
		log.Printf("[INFO] got new command '%s' from '%s' in '%s'", text, user.Username, chat.Title)
		msg, mthd, err := p.selectCommand(text, chat, user, messageID)
		if err != nil {
			return e.Wrap(fmt.Sprintf("can't select command from message: %s", text), err)
		}

		switch mthd {
		case UnsupportedMethod:
			return e.Wrap(fmt.Sprintf("unsupported method from message: %s", text), nil)
		case sendMessageMethod:
			return p.tg.SendMessage(chat.ID, msg)
		case sendPhotoMethod:
			return p.tg.SendPhoto(chat.ID, msg)
		}
	}

	return nil
}

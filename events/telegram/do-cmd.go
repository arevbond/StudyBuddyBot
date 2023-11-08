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

	switch utils.CheckYesOrNo(text) {
	case utils.IsYesCommand:
		return p.tg.SendMessage(chat.ID, "Пизда")
	case utils.IsNoCommand:
		return p.tg.SendMessage(chat.ID, "Пид*ра ответ")
	}

	if utils.IsCommand(text) {
		log.Printf("[INFO] got new command '%s' from '%s' in '%s'", text, user.Username, chat.Title)
		msg, mthd, buttons, err := p.selectCommand(text, chat, user, messageID)
		if err != nil {
			return e.Wrap(fmt.Sprintf("can't select command from message: %s", text), err)
		}

		switch mthd {
		case UnsupportedMethod:
			return e.Wrap(fmt.Sprintf("unsupported method from message %s: %v", text), nil)
		case sendMessageMethod:
			return p.tg.SendMessage(chat.ID, msg)
		case sendPhotoMethod:
			return p.tg.SendPhoto(chat.ID, msg)
		case sendMessageWithButtonsMethod:
			return p.tg.SendMessageWithMarkup(chat.ID, msg, buttons)
		}
	}

	return nil
}

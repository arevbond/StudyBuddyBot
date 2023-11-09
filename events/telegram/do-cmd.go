package telegram

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"tg_ics_useful_bot/clients/telegram"
	"tg_ics_useful_bot/lib/e"
	"tg_ics_useful_bot/lib/utils"
)

func (p *Processor) doCmd(text string, chat *telegram.Chat, user *telegram.User, messageID int) error {
	text = strings.TrimSpace(text)

	parseMode := ""

	switch utils.CheckYesOrNo(text) {
	case utils.IsYesCommand:
		return p.tg.SendMessage(chat.ID, "Пизда", parseMode, messageID)
	case utils.IsNoCommand:
		return p.tg.SendMessage(chat.ID, "Пидора ответ", parseMode, messageID)
	}

	if utils.IsCommand(text) || len(stateHomework) > 0 {
		log.Printf("[INFO] got new command '%s' from '%s' in '%s'", text, user.Username, chat.Title)
		msg, mthd, parseMode, replyToMessageID, err := p.selectCommand(text, chat, user, messageID)

		if err != nil {
			return e.Wrap(fmt.Sprintf("can't select command from message: %s", text), err)
		}

		switch mthd {
		case UnsupportedMethod:
			return e.Wrap("unsupported method:", errors.New("unknown method"))
		case sendMessageMethod:
			return p.tg.SendMessage(chat.ID, msg, parseMode, replyToMessageID)
		case sendPhotoMethod:
			return p.tg.SendPhoto(chat.ID, msg)
		case sendMessageWithButtonsMethod:
			return p.tg.SendMessage(chat.ID, msg, parseMode, replyToMessageID)
		}
	}

	return nil
}

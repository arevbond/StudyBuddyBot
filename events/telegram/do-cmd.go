package telegram

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"tg_ics_useful_bot/clients/telegram"
	"tg_ics_useful_bot/lib/e"
	"tg_ics_useful_bot/lib/utils"
	"tg_ics_useful_bot/storage"
)

func (p *Processor) doCmd(text string, chat *telegram.Chat, user *telegram.User, messageID int) error {
	userStats, err := p.storage.UserStatsByTelegramIDAndChatID(context.Background(), user.ID, chat.ID)
	if err == storage.ErrUserNotExist {
		userStats = storage.NewDBUserStats(user.ID, chat.ID, user.Username, user.FirstName, user.LastName)
		err = p.storage.CreateUserStats(context.Background(), userStats)
		if err != nil {
			return err
		}
	}
	err = p.storage.IncreaseMessageCount(context.Background(), userStats)
	if err != nil {
		log.Print(err)
	}

	text = strings.TrimSpace(text)

	parseMode := ""

	switch utils.CheckYesOrNo(text) {
	case utils.IsYesCommand:
		err = p.storage.IncreaseYesCount(context.Background(), userStats)
		if err != nil {
			log.Print(err)
		}
		return p.tg.SendMessage(chat.ID, "Пизда", parseMode, messageID)
	case utils.IsNoCommand:
		err = p.storage.IncreaseNoCount(context.Background(), userStats)
		if err != nil {
			log.Print(err)
		}
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

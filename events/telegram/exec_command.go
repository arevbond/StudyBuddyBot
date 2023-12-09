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

const (
	MAX_DICK_CHANGE_COUNT = 3
)

func (p *Processor) doCmd(text string, chat *telegram.Chat, user *telegram.User, messageID int) error {
	dbUser, err := p.storage.UserByTelegramID(context.Background(), user.ID, chat.ID)
	if err == storage.ErrUserNotExist {
		dbUser, err = p.createNewUserInDB(chat.ID, user)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	userStats, err := p.storage.GetUserStats(context.Background(), dbUser)
	if err == storage.ErrUserNotExist {
		return e.Wrap("not find user stats", err)
	} else if err != nil {
		return e.Wrap("not user stats", err)
	}

	userStats.MessageCount++
	err = p.storage.UpdateUserStats(context.Background(), userStats)
	if err != nil {
		log.Print(err)
	}

	text, parseMode := strings.TrimSpace(text), ""

	switch utils.CheckYesOrNo(text) {
	case utils.IsYesCommand:

		userStats.YesCount++
		err = p.storage.UpdateUserStats(context.Background(), userStats)
		if err != nil {
			log.Print(err)
		}
		return p.tg.SendMessage(chat.ID, "Пизда", parseMode, messageID)
	case utils.IsNoCommand:

		userStats.NoCount++
		err = p.storage.UpdateUserStats(context.Background(), userStats)
		if err != nil {
			log.Print(err)
		}
		return p.tg.SendMessage(chat.ID, "Пидора ответ", parseMode, messageID)
	}

	if utils.IsCommand(text) || len(stateHomework) > 0 {
		log.Printf("[INFO] got new command '%s' from '%s' in '%s'", text, user.Username, chat.Title)
		msg, mthd, parseMode, replyToMessageID, err := p.selectCommand(text, chat, user, userStats, messageID)

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
		case doNothingMethod:
			log.Printf("Message: \"%s\" - do nothing", text)
		}
	}

	return nil
}

func (p *Processor) createNewUserInDB(chatID int, user *telegram.User) (*storage.DBUser, error) {
	dbUserStatID, err := p.storage.CreateUserStats(context.Background(), &storage.DBUserStat{})

	if err != nil {
		return nil, e.Wrap("can't create user stats in 'createNewUserInDB'", err)
	}
	dbUser := &storage.DBUser{
		TgID:               user.ID,
		ChatID:             chatID,
		IsBot:              user.IsBot,
		IsPremium:          user.IsPremium,
		FirstName:          user.FirstName,
		LastName:           user.LastName,
		Username:           user.Username,
		UserStatId:         dbUserStatID,
		MaxDickChangeCount: MAX_DICK_CHANGE_COUNT,
	}
	err = p.storage.CreateUser(context.Background(), dbUser)

	if err != nil {
		return nil, err
	}
	return dbUser, nil
}

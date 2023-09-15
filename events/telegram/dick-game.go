package telegram

import (
	"context"
	"fmt"
	"tg_ics_useful_bot/clients/telegram"
	"tg_ics_useful_bot/lib/e"
	"tg_ics_useful_bot/lib/game"
	"tg_ics_useful_bot/storage"
	"time"
)

func (p *Processor) topDick(chat *telegram.Chat) (err error) {
	users, err := p.storage.UsersByChat(context.Background(), chat.ID)
	if err != nil {
		return e.Wrap("[ERROR] can't get users: ", err)
	}
	result := ""
	for i, u := range users {
		result += fmt.Sprintf("%d. %s — %d см\n", i+1, u.FirstName+" "+u.LastName, u.DickSize)
	}
	return p.tg.SendMessage(chat.ID, result)
}

func (p *Processor) gameDick(chat *telegram.Chat, user *telegram.User, messageID int) (err error) {
	defer func() { err = e.WrapIfErr("[ERROR] can't change dick size: ", err) }()

	err = p.tg.DeleteMessage(chat.ID, messageID)
	if err != nil {
		return e.Wrap(fmt.Sprintf("[ERROR] can't delete message: user #%d, chat id #%d", user.ID, chat.ID), err)
	}

	dbUser, err := p.storage.UserByTelegramID(context.Background(), user.ID, chat.ID)

	if err == storage.ErrUserNotExist {
		u, err2 := p.createNewPlayer(chat.ID, user)
		if err2 != nil {
			return e.Wrap(fmt.Sprintf("[ERROR] can't create new player telegram id = #%d", user.ID), err)
		}
		return p.tg.SendMessage(chat.ID, fmt.Sprintf(msgCreateUser, u.Username)+fmt.Sprintf(msgDickSize, u.DickSize))
	} else if err != nil {
		return err
	}

	if game.CanChangeDickSize(dbUser) {
		_, oldDickSize, err := p.changeRandomDickSize(dbUser)
		if err != nil {
			return err
		}
		return p.tg.SendMessage(chat.ID, fmt.Sprintf(msgChangeDickSize, dbUser.Username, oldDickSize, dbUser.DickSize))
	}
	return p.tg.SendMessage(chat.ID, fmt.Sprintf(msgAlreadyPlays, dbUser.Username))
}

func (p *Processor) duelDick(chat *telegram.Chat, user *telegram.User, targetUsername string) error {
	u1, err := p.storage.UserByTelegramID(context.Background(), user.ID, chat.ID)
	if err != nil {
		return err
	}
	u2, err := p.storage.UserByUsername(context.Background(), targetUsername, chat.ID)
	if err == storage.ErrUserNotExist {
		return p.tg.SendMessage(chat.ID, fmt.Sprintf(msgTargetNotFound, targetUsername))
	} else if err != nil {
		return err
	}

	User1Win, ch1, ch2 := game.Duel(u1.DickSize, u2.DickSize)
	if User1Win {
		oldDickSize, err := p.changeDickSize(u1, game.PositiveRandomValue())
		if err != nil {
			return err
		}
		return p.tg.SendMessage(chat.ID,
			fmt.Sprintf(msgChanceDuel, u1.Username, oldDickSize, ch1, targetUsername, u2.DickSize, ch2)+
				fmt.Sprintf(msgVictoryInDuel, u1.Username, u2.Username)+
				fmt.Sprintf(msgDickSize, u1.DickSize))
	} else {
		oldDickSize, err := p.changeDickSize(u1, -1*game.PositiveRandomValue())
		if err != nil {
			return err
		}
		return p.tg.SendMessage(chat.ID,
			fmt.Sprintf(msgChanceDuel, u1.Username, oldDickSize, ch1, targetUsername, u2.DickSize, ch2)+
				fmt.Sprintf(msgVictoryInDuel, u2.Username, u1.Username))
	}

}

func (p *Processor) createNewPlayer(chatID int, user *telegram.User) (*storage.DBUser, error) {
	dbUser := &storage.DBUser{
		TgID:           user.ID,
		ChatID:         chatID,
		IsBot:          user.IsBot,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		Username:       user.Username,
		IsPremium:      user.IsPremium,
		DickSize:       game.PositiveRandomValue(),
		DateChangeDick: time.Now(),
	}
	err := p.storage.CreateUser(context.Background(), dbUser)
	if err != nil {
		return nil, err
	}
	return dbUser, err
}

func (p *Processor) changeDickSize(user *storage.DBUser, value int) (int, error) {
	oldDickSize := user.DickSize

	err := p.storage.UpdateUserDickSize(context.Background(), user, user.DickSize+value)
	if err != nil {
		return 0, e.Wrap(fmt.Sprintf("[ERROR] chat id %d, user %s can't change dick size: ", user.ChatID, user.Username), err)
	}
	err = p.storage.UpdateDateLastTryChangeDickToNow(context.Background(), user)
	if err != nil {
		return 0, e.Wrap("[ERROR] can't update time to now: ", err)
	}
	return oldDickSize, nil
}

func (p *Processor) changeRandomDickSize(user *storage.DBUser) (bool, int, error) {
	value := game.RandomValue()
	oldDickSize, err := p.changeDickSize(user, value)
	if err != nil {
		return false, 0, e.Wrap(fmt.Sprintf("[ERROR] chat id %d, user %s can't change random dick size: ",
			user.ChatID, user.Username), err)
	}
	return value >= 0, oldDickSize, nil
}

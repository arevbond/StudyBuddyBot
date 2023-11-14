package telegram

import (
	"context"
	"fmt"
	"log"
	"tg_ics_useful_bot/clients/telegram"
	"tg_ics_useful_bot/lib/dick"
	"tg_ics_useful_bot/lib/e"
	"tg_ics_useful_bot/storage"
	"time"
)

var duels = make(map[string]*storage.DBUser)

var reward = 15

func (p *Processor) topDick(chatID int) (msg string, err error) {
	users, err := p.storage.UsersByChat(context.Background(), chatID)
	if err != nil {
		return "", e.Wrap("[ERROR] can't get users: ", err)
	}
	result := ""
	for i, u := range users {
		result += fmt.Sprintf("%d. %s — %d см\n", i+1, u.FirstName+" "+u.LastName, u.DickSize)
	}
	return result, nil
}

func (p *Processor) gameDick(chat *telegram.Chat, user *telegram.User, messageID int) (msg string, err error) {
	defer func() { err = e.WrapIfErr("[ERROR] can't change dick size: ", err) }()

	err = p.tg.DeleteMessage(chat.ID, messageID)
	if err != nil {
		return "", e.Wrap(fmt.Sprintf("[ERROR] can't delete message: user #%d, chat id #%d", user.ID, chat.ID), err)
	}

	dbUser, err := p.storage.UserByTelegramID(context.Background(), user.ID, chat.ID)

	if err == storage.ErrUserNotExist {
		u, err2 := p.createNewPlayer(chat.ID, user)
		if err2 != nil {
			return "", e.Wrap(fmt.Sprintf("[ERROR] can't create new player telegram id = #%d", user.ID), err)
		}
		return fmt.Sprintf(msgCreateUser, u.Username) + fmt.Sprintf(msgDickSize, u.DickSize), nil
	} else if err != nil {
		return "", err
	}

	if dick.CanChangeDickSize(dbUser) {
		_, oldDickSize, err := p.changeRandomDickSize(dbUser)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf(msgChangeDickSize, dbUser.Username, oldDickSize, dbUser.DickSize), nil
	}
	return fmt.Sprintf(msgAlreadyPlays, dbUser.Username), nil
}

// TODO: изменить формулу для награды
func (p *Processor) gameDuelDick(chat *telegram.Chat, messageID int, user *telegram.User, targetUsername string) (string, error) {
	err := p.tg.DeleteMessage(chat.ID, messageID)
	if err != nil {
		return "", e.Wrap(fmt.Sprintf("[ERROR] can't delete message: user #%d, chat id #%d", user.ID, chat.ID), err)
	}

	u1, err := p.storage.UserByTelegramID(context.Background(), user.ID, chat.ID)
	if err != nil {
		return "", err
	}
	u2, err := p.storage.UserByUsername(context.Background(), targetUsername, chat.ID)
	if err == storage.ErrUserNotExist {
		return fmt.Sprintf(msgTargetNotFound, targetUsername), nil
	} else if err != nil {
		return "", err
	}

	if u1.TgID == u2.TgID || u2.IsBot {
		return fmt.Sprintf(msgDuelWithYourself, u1.Username), nil
	}

	if enemy, ok := duels[u1.Username]; ok && enemy.TgID == u2.TgID {
		delete(duels, u1.Username)
		User1Win, ch1, ch2 := dick.Duel(u1.DickSize, u2.DickSize)
		if User1Win {
			if ch1 > 65 {
				reward = 5
			}
			oldDickSize1, err2 := p.changeDickSize(u1, reward)
			if err2 != nil {
				return "", err
			}
			oldDickSize2, err3 := p.changeDickSize(u2, -1*reward)
			if err3 != nil {
				return "", err3
			}
			return fmt.Sprintf(msgAcceptDuel, u1.Username, oldDickSize1, ch1, u2.Username, oldDickSize2, ch2) +
				fmt.Sprintf(msgUser1Wins, u1.Username, u1.DickSize, u2.Username, u2.DickSize), nil
		} else {
			if ch1 <= 35 {
				reward = 5
			}
			oldDickSize1, err2 := p.changeDickSize(u1, -1*reward)
			if err2 != nil {
				return "", err
			}
			oldDickSize2, err3 := p.changeDickSize(u2, reward)
			if err3 != nil {
				return "", err3
			}
			return fmt.Sprintf(msgAcceptDuel, u1.Username, oldDickSize1, ch1, u2.Username, oldDickSize2, ch2) +
				fmt.Sprintf(msgUser1Lost, u1.Username, u1.DickSize, u2.Username, u2.DickSize), nil
		}
	} else {
		duels[targetUsername] = u1
		return fmt.Sprintf(msgChallengeToDuel, u1.Username, targetUsername), nil
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
		DickSize:       dick.PositiveRandomValue(),
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
	value := dick.RandomValue()

	userStats, err := p.storage.UserStatsByTelegramIDAndChatID(context.Background(), user.TgID, user.ChatID)
	if err != nil {
		log.Print(err)
	}
	if value > 0 {
		err = p.storage.IncreaseDickPlusCount(context.Background(), userStats)
		if err != nil {
			log.Print(err)
		}
	} else {
		err = p.storage.IncreaseDickMinusCount(context.Background(), userStats)
		if err != nil {
			log.Print(err)
		}
	}

	oldDickSize, err := p.changeDickSize(user, value)
	if err != nil {
		return false, 0, e.Wrap(fmt.Sprintf("[ERROR] chat id %d, user %s can't change random dick size: ",
			user.ChatID, user.Username), err)
	}
	return value >= 0, oldDickSize, nil
}

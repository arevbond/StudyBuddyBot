package telegram

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"tg_ics_useful_bot/clients/telegram"
	"tg_ics_useful_bot/lib/dick"
	"tg_ics_useful_bot/lib/e"
	"tg_ics_useful_bot/storage"
	"time"
)

func (p *Processor) topDicks(chatID int) (msg string, err error) {
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

func (p *Processor) gameDick(chat *telegram.Chat, user *telegram.User, userStats *storage.DBUserStat, messageID int) (msg string, err error) {
	defer func() { err = e.WrapIfErr("can't change dick size: ", err) }()

	err = p.tg.DeleteMessage(chat.ID, messageID)
	if err != nil {
		return "", e.Wrap(fmt.Sprintf("can't delete message: user #%d, chat id #%d", user.ID, chat.ID), err)
	}

	dbUser, err := p.storage.UserByTelegramID(context.Background(), user.ID, chat.ID)

	if err != nil {
		return "", err
	}

	canChange, err := p.CanChangeDickSize(dbUser)
	if err != nil {
		return "", err
	}
	if canChange {
		_, oldDickSize, err := p.changeRandomDickSize(dbUser, userStats)
		if err != nil {
			return "", err
		}
		if oldDickSize == 0 {
			return fmt.Sprintf(msgCreateUser, dbUser.Username) + fmt.Sprintf(msgDickSize, dbUser.DickSize), nil
		}
		return fmt.Sprintf(msgChangeDickSize, dbUser.Username, oldDickSize, dbUser.DickSize), nil
	}
	return fmt.Sprintf(msgAlreadyPlays, dbUser.Username), nil
}

func (p *Processor) changeDickSizeAndTime(user *storage.DBUser, value int) (int, error) {
	oldDickSize := user.DickSize

	user.DickSize += value
	user.ChangeDickAt = time.Now()
	err := p.storage.UpdateUser(context.Background(), user)
	if err != nil {
		return 0, e.Wrap(fmt.Sprintf("chat id %d, user %s can't change dick size or change dick at: ", user.ChatID, user.Username), err)
	}
	return oldDickSize, nil
}

// TODO: переделать
func (p *Processor) changeRandomDickSize(user *storage.DBUser, userStats *storage.DBUserStat) (bool, int, error) {
	var value int
	for {
		value = dick.RandomValue()
		if user.DickSize+value > 0 {
			break
		}
	}

	if rand.Intn(101) == 77 {
		value = 100
	}

	if value > 0 {
		userStats.DickPlusCount++
		err := p.storage.UpdateUserStats(context.Background(), userStats)
		if err != nil {
			log.Print(err)
		}
	} else {
		userStats.DickMinusCount++
		err := p.storage.UpdateUserStats(context.Background(), userStats)
		if err != nil {
			log.Print(err)
		}
	}

	oldDickSize, err := p.changeDickSizeAndTime(user, value)
	if err != nil {
		return false, 0, e.Wrap(fmt.Sprintf("[ERROR] chat id %d, user %s can't change random dick size: ",
			user.ChatID, user.Username), err)
	}
	return value >= 0, oldDickSize, nil
}

func (p *Processor) CanChangeDickSize(user *storage.DBUser) (bool, error) {
	yearLastTry, monthLastTry, dayLastTry := user.ChangeDickAt.Date()
	year, month, today := time.Now().Date()
	if (month == monthLastTry && today > dayLastTry) || month > monthLastTry || year > yearLastTry {
		user.CurDickChangeCount = 0
		err := p.storage.UpdateUser(context.Background(), user)
		if err != nil {
			return false, e.Wrap("can't update user in 'CanChangeDickSize'", err)
		}
	}
	if user.CurDickChangeCount+1 <= user.MaxDickChangeCount {
		user.CurDickChangeCount++
		err := p.storage.UpdateUser(context.Background(), user)
		if err != nil {
			return false, e.Wrap("can't update user in 'CanChangeDickSize'", err)
		}
		return true, nil
	}

	return false, nil
}

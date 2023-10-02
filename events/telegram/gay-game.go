package telegram

import (
	"context"
	"fmt"
	"math/rand"
	"sort"
	"tg_ics_useful_bot/clients/telegram"
	"tg_ics_useful_bot/lib/e"
	"tg_ics_useful_bot/storage"
	"time"
)

func (p *Processor) gameGayTop(chatID int) (message string, err error) {
	admins, err := p.tg.ChatAdministrators(chatID)
	if err != nil {
		return "", e.Wrap("[ERROR] can't get chat administrators: ", err)
	}
	dbUsers := []*storage.DBUser{}
	for _, u := range admins {
		dbUser, err := p.storage.UserByTelegramID(context.Background(), u.ID, chatID)
		if err == storage.ErrUserNotExist {
			dbUser = &storage.DBUser{
				TgID:      u.ID,
				ChatID:    chatID,
				IsBot:     u.IsBot,
				FirstName: u.FirstName,
				LastName:  u.LastName,
				Username:  u.Username,
			}
			err = p.storage.CreateUser(context.Background(), dbUser)
			if err != nil {
				return "", err
			}
		}
		dbUsers = append(dbUsers, dbUser)
	}
	sort.Slice(dbUsers, func(i, j int) bool {
		return dbUsers[i].CountGayOfDay >= dbUsers[j].CountGayOfDay
	})
	result := "Рейтинг пидоров: \n\n"

	for i, dbU := range dbUsers {
		result += fmt.Sprintf("%d. %s %s - %d раз \n", i+1, dbU.FirstName, dbU.LastName, dbU.CountGayOfDay)
	}
	return result, nil
}

func (p *Processor) gameGay(chatID int) (string, error) {
	admins, err := p.tg.ChatAdministrators(chatID)
	if err != nil {
		return "", e.Wrap("can't get chat administrators: ", err)
	}

	gay, err := p.storage.GayOfDay(context.Background(), chatID)
	if err == storage.ErrUserNotExist {
		gay, err = p.createNewGayOfDay(chatID, admins)
		return fmt.Sprintf(msgNewGayOfDay, gay.Username), nil
	} else if err != nil {
		return "", e.Wrap("can't get gay of day: ", err)
	}
	if (gay.DateLastUsed.Month() == time.Now().Month() && gay.DateLastUsed.Day() < time.Now().Day()) || gay.DateLastUsed.Month() < time.Now().Month() {
		err = p.storage.RemoveGayOfDay(context.Background(), chatID)
		if err != nil {
			return "", err
		}
		gay, err = p.createNewGayOfDay(chatID, admins)
		return fmt.Sprintf(msgNewGayOfDay, gay.Username), nil
	}
	return fmt.Sprintf(msgCurrentGayOfDay, gay.Username), nil
}

func (p *Processor) createNewGayOfDay(chatID int, admins []telegram.User) (*storage.DBGayOfDay, error) {
	rand.Seed(time.Now().Unix())
	n := rand.Intn(len(admins))
	u := admins[n]
	dbUser, err := p.storage.UserByTelegramID(context.Background(), u.ID, chatID)
	if err == storage.ErrUserNotExist {
		dbUser = &storage.DBUser{
			TgID:      u.ID,
			ChatID:    chatID,
			IsBot:     u.IsBot,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Username:  u.Username,
		}
		err = p.storage.CreateUser(context.Background(), dbUser)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	gay := &storage.DBGayOfDay{
		ChatID:       chatID,
		TgID:         dbUser.TgID,
		Username:     dbUser.Username,
		DateLastUsed: time.Now(),
	}
	err = p.storage.CreateGayOfDay(context.Background(), gay)
	if err != nil {
		return nil, err
	}
	err = p.storage.IncreaseCountOfGay(context.Background(), dbUser)
	if err != nil {
		return nil, err
	}
	return gay, nil
}

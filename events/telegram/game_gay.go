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

// gameGay определяет пидора дня среди администратора и возвращает сообщение для чата.
func (p *Processor) gameGay(chatID int) (string, error) {
	admins, err := p.tg.ChatAdministrators(chatID)
	if err != nil {
		return "", e.Wrap("can't get chat administrators: ", err)
	}

	gay, err := p.storage.GetGayOfDay(context.Background(), chatID)
	if err == storage.ErrUserNotExist {
		gay, err = p.createNewGayOfDay(chatID, admins)
		return fmt.Sprintf(msgNewGayOfDay, gay.Username), nil
	} else if err != nil {
		return "", e.Wrap("can't get gay of day: ", err)
	}
	if (gay.CreatedAt.Month() == time.Now().Month() && gay.CreatedAt.Day() < time.Now().Day()) || gay.CreatedAt.Month() < time.Now().Month() {
		err = p.storage.RemoveGayOfDay(context.Background(), chatID)
		if err != nil {
			return "", err
		}
		gay, err = p.createNewGayOfDay(chatID, admins)
		return fmt.Sprintf(msgNewGayOfDay, gay.Username), nil
	}
	return fmt.Sprintf(msgCurrentGayOfDay, gay.Username), nil
}

// createNewGayOfDay создаёт пидора дня.
func (p *Processor) createNewGayOfDay(chatID int, admins []telegram.User) (*storage.DBGay, error) {
	n := rand.Intn(len(admins))
	u := &admins[n]
	dbUser, err := p.storage.GetUser(context.Background(), u.ID, chatID)
	if err == storage.ErrUserNotExist {
		dbUser, err = p.createNewUserInDB(chatID, u)
		if err != nil {
			return nil, e.Wrap("can't create new user in db; in 'createNewGayOfDay'", err)
		}
	} else if err != nil {
		return nil, err
	}
	gay := &storage.DBGay{
		ChatID:    chatID,
		TgID:      dbUser.TgID,
		Username:  dbUser.Username,
		CreatedAt: time.Now(),
	}
	err = p.storage.CreateGayOfDay(context.Background(), gay)
	if err != nil {
		return nil, err
	}

	userStats, err := p.storage.GetUserStats(context.Background(), dbUser)
	if err != nil {
		return nil, e.Wrap("can't get user stats in 'createNewGayOfDay'", err)
	}
	userStats.GayCount++
	err = p.storage.UpdateUserStats(context.Background(), userStats)
	if err != nil {
		return nil, e.Wrap("can't update userStats in 'createNewGayOfDay'", err)
	}

	return gay, nil
}

// topGaysCmd возвращает список всех админов и сколько раз они были пидорами.
func (p *Processor) topGays(chatID int) (message string, err error) {
	admins, err := p.tg.ChatAdministrators(chatID)
	if err != nil {
		return "", e.Wrap("[ERROR] can't get chat administrators: ", err)
	}
	dbUsers := []*storage.DBUser{}
	dbUsersStats := []*storage.DBUserStat{}
	for _, u := range admins {
		dbUser, err := p.storage.GetUser(context.Background(), u.ID, chatID)
		if err == storage.ErrUserNotExist {
			dbUser, err = p.createNewUserInDB(chatID, &u)
			if err != nil {
				return "", e.Wrap("can't create new user in db; in 'createNewGayOfDay'", err)
			}
		}
		dbUserStat, err := p.storage.GetUserStats(context.Background(), dbUser)
		if err != nil {
			return "", err
		}
		dbUsersStats = append(dbUsersStats, dbUserStat)
		dbUsers = append(dbUsers, dbUser)
	}
	sort.Slice(dbUsers, func(i, j int) bool {
		return dbUsersStats[i].GayCount > dbUsersStats[j].GayCount
	})
	sort.Slice(dbUsersStats, func(i, j int) bool {
		return dbUsersStats[i].GayCount > dbUsersStats[j].GayCount
	})
	result := "Рейтинг пидоров: \n\n"

	for i, dbU := range dbUsers {
		result += fmt.Sprintf("%d. %s %s %d раз \n", i+1, dbU.FirstName, dbU.LastName, dbUsersStats[i].GayCount)
	}
	return result, nil
}

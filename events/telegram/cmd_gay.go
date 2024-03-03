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

// gayExec предоставляет метод Exec для выполнения /gay.
type gayExec string

// Exec: /gay - определяет случайного пидора в чате среди админов чата.
func (a gayExec) Exec(p *Processor, inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {

	message, err := p.gameGay(chat.ID)
	if err != nil {
		return nil, e.Wrap("can't get message from gameGay: ", err)
	}
	mthd := sendMessageMethod
	return &Response{message: message, method: mthd, replyMessageId: -1}, nil
}

// topGaysExec предоставляет метод Exec для вывода топа пидоров.
type topGaysExec string

// Exec: /top_gay - выводит список участников чата и их кол-во становления пидором дня.
func (a topGaysExec) Exec(p *Processor, inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {
	message, err := p.topGays(chat.ID)
	if err != nil {
		return nil, e.Wrap("can't do GayTop: ", err)
	}
	mthd := sendMessageMethod
	return &Response{message: message, method: mthd, replyMessageId: -1}, nil
}

// gameGay определяет пидора дня среди администратора и возвращает сообщение для чата.
func (p *Processor) gameGay(chatID int) (string, error) {
	admins, err := p.tg.ChatAdministrators(chatID)
	if err != nil {
		return "", e.Wrap("can't get chat administrators: ", err)
	}

	gay, isCreated, err := p.gayOfDay(chatID, admins)
	if err != nil {
		return "", e.Wrap("can't initial gameGay", err)
	}
	name, hasUsername, err := p.getGayName(gay)
	if err != nil {
		return "", e.Wrap("can't get gay name", err)
	}
	if isCreated {
		if hasUsername {
			return fmt.Sprintf(msgNewGayOfDayUsername, name), nil
		} else {
			return fmt.Sprintf(msgNewGayOfDayFullName, name), nil
		}
	}
	if hasUsername {
		return fmt.Sprintf(msgCurrentGayOfDayUsername, name), nil
	}
	return fmt.Sprintf(msgCurrentGayOfDayFullName, name), nil
}

func (p *Processor) gayOfDay(chatID int, admins []telegram.User) (*storage.DBGay, bool, error) {
	gay, err := p.storage.GetGayOfDay(context.Background(), chatID)
	if err == storage.ErrUserNotExist {
		gay, err = p.createNewGayOfDay(chatID, admins)
		return gay, true, nil
	} else if err != nil {
		return nil, false, e.Wrap("can't get gay of day: ", err)
	}
	if gayIsOld(gay) {
		err = p.storage.RemoveGayOfDay(context.Background(), chatID)
		if err != nil {
			return nil, false, err
		}
		gay, err = p.createNewGayOfDay(chatID, admins)
		return gay, true, nil
	}
	return gay, false, nil
}

func gayIsOld(gay *storage.DBGay) bool {
	return (gay.CreatedAt.Month() == time.Now().Month() && gay.CreatedAt.Day() < time.Now().Day()) || gay.CreatedAt.Month() < time.Now().Month() || gay.CreatedAt.Year() < time.Now().Year()
}

func (p *Processor) getGayName(gay *storage.DBGay) (string, bool, error) {
	dbUser, err := p.storage.GetUser(context.Background(), gay.TgID, gay.ChatID)
	if dbUser.Username != "" {
		return dbUser.Username, true, nil
	}
	if err != nil {
		return "", false, e.Wrap("can't gay user from storage", err)
	}
	return dbUser.FirstName + " " + dbUser.LastName, false, nil
}

// createNewGayOfDay создаёт пидора дня.
func (p *Processor) createNewGayOfDay(chatID int, admins []telegram.User) (*storage.DBGay, error) {
	var user *telegram.User
	for {
		n := rand.Intn(len(admins))
		u := &admins[n]
		if !u.IsBot {
			user = u
			break
		}
	}
	dbUser, err := p.storage.GetUser(context.Background(), user.ID, chatID)
	if err == storage.ErrUserNotExist {
		dbUser, err = p.createNewUserInDB(chatID, user)
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

// topGaysExec возвращает список всех админов и сколько раз они были пидорами.
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
		result += fmt.Sprintf("%d. %s %s: %d раз(а)\n", i+1, dbU.FirstName, dbU.LastName, dbUsersStats[i].GayCount)
	}
	return result, nil
}

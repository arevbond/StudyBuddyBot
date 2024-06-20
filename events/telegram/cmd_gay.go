package telegram

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"tg_ics_useful_bot/clients/telegram"
	"tg_ics_useful_bot/lib/e"
	"tg_ics_useful_bot/storage"
	"time"
)

// topGaysExec предоставляет метод Exec для вывода топа пидоров.
type topGaysExec string

// Exec: /top_gay - выводит список участников чата и их кол-во становления пидором дня.
func (a topGaysExec) Exec(p *Processor, inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {
	message, err := a.topGays(chat.ID, p)
	if err != nil {
		return nil, e.Wrap("can't do GayTop: ", err)
	}
	mthd := sendMessageMethod
	return &Response{message: message, method: mthd, replyMessageId: -1}, nil
}

// topGaysExec возвращает список всех админов и сколько раз они были пидорами.
func (a topGaysExec) topGays(chatID int, p *Processor) (message string, err error) {
	admins, err := p.tg.ChatAdministrators(chatID)

	if err != nil {
		return "", e.Wrap("can't get chat administrators: ", err)
	}

	type userWithCount struct {
		name  string
		count int
	}
	users := []userWithCount{}

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
		users = append(users, userWithCount{
			name:  dbUser.FirstName + " " + dbUser.LastName,
			count: dbUserStat.GayCount,
		})
	}

	sort.Slice(users, func(i, j int) bool {
		return users[i].count > users[j].count
	})

	result := "Рейтинг пидоров: \n\n"
	for i, user := range users {
		result += fmt.Sprintf("%d. %s %d раз(а)\n", i+1, user.name, user.count)
	}
	return result, nil
}

// gayExec предоставляет метод Exec для выполнения /gay.
type gayExec string

// Exec: /gay - определяет случайного пидора в чате среди админов чата.
func (g gayExec) Exec(p *Processor, inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {

	message, err := g.gameGay(chat.ID, p)
	if err != nil {
		return nil, e.Wrap("can't get message from gameGay: ", err)
	}
	mthd := sendMessageMethod
	return &Response{message: message, method: mthd, replyMessageId: -1}, nil
}

// gameGay определяет пидора дня среди администратора и возвращает сообщение для чата.
func (g gayExec) gameGay(chatID int, p *Processor) (string, error) {
	admins, err := p.tg.ChatAdministrators(chatID)
	if err != nil {
		return "", e.Wrap("can't get chat administrators: ", err)
	}

	gay, isCreated, err := g.gayOfDay(chatID, admins, p)
	if err != nil {
		return "", e.Wrap("can't initial gameGay", err)
	}
	name, hasUsername, err := g.getGayName(gay, p.storage)
	if err != nil {
		return "", e.Wrap("can't get gay name", err)
	}
	return g.formatOutputGay(hasUsername, isCreated, name), nil
}

func (g gayExec) formatOutputGay(hasUsername bool, isCreated bool, name string) string {
	if isCreated {
		if hasUsername {
			return fmt.Sprintf(msgNewGayOfDayUsername, name)
		} else {
			return fmt.Sprintf(msgNewGayOfDayFullName, name)
		}
	}
	if hasUsername {
		return fmt.Sprintf(msgCurrentGayOfDayUsername, name)
	}
	return fmt.Sprintf(msgCurrentGayOfDayFullName, name)
}

func (g gayExec) gayOfDay(chatID int, admins []telegram.User, p *Processor) (*storage.DBGay, bool, error) {
	gay, err := p.storage.GetGayOfDay(context.Background(), chatID)
	if err == storage.ErrUserNotExist {
		gay, err = g.createNewGayOfDay(chatID, admins, p)
		return gay, true, nil
	} else if err != nil {
		return nil, false, e.Wrap("can't get gay of day: ", err)
	}
	if g.gayIsOld(gay) {
		err = p.storage.RemoveGayOfDay(context.Background(), chatID)
		if err != nil {
			return nil, false, err
		}
		gay, err = g.createNewGayOfDay(chatID, admins, p)
		return gay, true, nil
	}
	return gay, false, nil
}

func (g gayExec) gayIsOld(gay *storage.DBGay) bool {
	return (gay.CreatedAt.Month() == time.Now().Month() && gay.CreatedAt.Day() < time.Now().Day()) || gay.CreatedAt.Month() < time.Now().Month() || gay.CreatedAt.Year() < time.Now().Year()
}

func (g gayExec) getGayName(gay *storage.DBGay, db storage.Storage) (string, bool, error) {
	dbUser, err := db.GetUser(context.Background(), gay.TgID, gay.ChatID)
	if dbUser.Username != "" {
		return dbUser.Username, true, nil
	}
	if err != nil {
		return "", false, e.Wrap("can't gay user from storage", err)
	}
	return dbUser.FirstName + " " + dbUser.LastName, false, nil
}

// createNewGayOfDay создаёт пидора дня.
func (g gayExec) createNewGayOfDay(chatID int, admins []telegram.User, p *Processor) (*storage.DBGay, error) {
	if len(admins) < 1 {
		return nil, e.Wrap("can't create gay of day", errors.New("zero admins"))
	}
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

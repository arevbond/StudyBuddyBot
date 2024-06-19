package telegram

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"strconv"
	"strings"
	"tg_ics_useful_bot/clients/telegram"
	"tg_ics_useful_bot/lib/e"
	"tg_ics_useful_bot/storage"
	"time"
)

const (
	jackpotValue = 100
	maxReward    = 20
	minReward    = 5
)

// dickTopExec предоставляет метод Exec для выполнения /top_dick.
type dickTopExec string

// Exec: /top_dick - пишет топ всех пенисов в чат.
func (d dickTopExec) Exec(p *Processor, inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {

	message, err := getTopDicks(chat.ID, p)
	if err != nil {
		return nil, e.Wrap(fmt.Sprintf("can't get top dics from chat %d: ", chat.ID), err)
	}
	mthd := sendMessageMethod
	return &Response{message: message, method: mthd, replyMessageId: -1, parseMode: telegram.Markdown}, nil
}

// getTopDicks возвращает string сообщение со списком всех dick > 0 в чате.
func getTopDicks(chatID int, p *Processor) (msg string, err error) {
	users, err := p.storage.UsersByChat(context.Background(), chatID)
	if err != nil {
		return "", e.Wrap("can't get users: ", err)
	}

	result := ""
	for i, u := range users {
		if u.DickSize > 0 && !u.IsBot {
			if i == 0 {
				result += fmt.Sprintf("👑 *%s* — _%d см_\n", u.FirstName+" "+u.LastName, u.DickSize)
			} else {
				result += fmt.Sprintf("%d. %s — %d см\n", i+1, u.FirstName+" "+u.LastName, u.DickSize)
			}
		}
	}
	return result, nil
}

// dickStartExec предоставляет метод Exec для выполнения /dick.
type dickStartExec string

// Exec: /dick - игра в пенис.
func (d dickStartExec) Exec(p *Processor, inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {

	message, err := d.gameDick(chat, user, userStats, p.storage, p.logger)
	if err != nil {
		return nil, e.Wrap("can't get message from gameDick: ", err)
	}
	mthd := sendMessageMethod
	return &Response{message: message, method: mthd, replyMessageId: -1}, nil
}

// gameDick это функция изменяющая размер пениса на случайное число и время изменения пениса.
// /dick - command
// Возвращает сообщение, отправляемое в чат.
func (d dickStartExec) gameDick(chat *telegram.Chat, user *telegram.User, userStats *storage.DBUserStat, db storage.Storage, logger *slog.Logger) (msg string, err error) {
	defer func() { err = e.WrapIfErr("error in gameDick: ", err) }()

	dbUser, err := db.GetUser(context.Background(), user.ID, chat.ID)
	if err != nil {
		return "", err
	}

	message, err := d.proccessDickGame(dbUser, userStats, db, logger)
	if err != nil {
		return "", e.Wrap("can't work game dick cmd", err)
	}
	return message, nil
}

func (d dickStartExec) proccessDickGame(dbUser *storage.DBUser, userStats *storage.DBUserStat, db storage.Storage, logger *slog.Logger) (string, error) {
	canChange, err := d.canChangeDickSize(dbUser, db)
	if err != nil {
		return "", err
	}

	if !canChange {
		return d.formatAlreadyPlaying(dbUser), nil
	}

	oldDickSize := dbUser.DickSize
	err = d.updateRandomDickAndChangeTime(dbUser, userStats, db, logger)
	if err != nil {
		return "", err
	}
	return d.formatOutputGameDick(dbUser, oldDickSize), nil
}

func (d dickStartExec) formatOutputGameDick(dbUser *storage.DBUser, oldDickSize int) string {
	name, hasUsername := d.getName(dbUser)
	if oldDickSize == 0 {
		if hasUsername {
			return fmt.Sprintf(msgCreateUserWithUsername, name) + fmt.Sprintf(msgDickSize, dbUser.DickSize)
		}
		return fmt.Sprintf(msgCreateUserWithFullName, name) + fmt.Sprintf(msgDickSize, dbUser.DickSize)
	}
	if hasUsername {
		return fmt.Sprintf(msgChangeDickSizeWithUsername, name, oldDickSize, dbUser.DickSize)
	}
	return fmt.Sprintf(msgChangeDickSizeWithFullName, name, oldDickSize, dbUser.DickSize)
}

func (d dickStartExec) formatAlreadyPlaying(dbUser *storage.DBUser) string {
	name, hasUsername := d.getName(dbUser)
	if hasUsername {
		return fmt.Sprintf(msgAlreadyPlaysWithUsername, name)
	}
	return fmt.Sprintf(msgAlreadyPlaysWithFullName, name)
}

func (d dickStartExec) getName(dbUser *storage.DBUser) (string, bool) {
	if dbUser.Username != "" {
		return dbUser.Username, true
	}
	return dbUser.FirstName + " " + dbUser.LastName, false
}

// updateRandomDickAndChangeTime изменяет значение пениса на слуайное число и время его изменения в базе данных.
func (d dickStartExec) updateRandomDickAndChangeTime(user *storage.DBUser, userStats *storage.DBUserStat, db storage.Storage, logger *slog.Logger) error {
	reward := d.calculateReward(user)

	if reward > 0 {
		userStats.DickPlusCount++
	} else {
		userStats.DickMinusCount++
	}
	err := db.UpdateUserStats(context.Background(), userStats)
	if err != nil {
		logger.Error("can't update users stats", slog.Any("error", err))
	}

	user.DickSize += reward
	user.ChangeDickAt = time.Now()
	err = db.UpdateUser(context.Background(), user)
	if err != nil {
		return e.Wrap(fmt.Sprintf("chat id %d, user %s can't change dick size or change dick at: ", user.ChatID, user.Username), err)
	}
	return nil
}

func (d dickStartExec) calculateReward(user *storage.DBUser) int {
	var reward int
	for {
		reward = d.reward()
		if reward > 0 || (reward < 0 && user.DickSize+reward > 0) {
			break
		}
	}
	if d.isJackpot() {
		reward = jackpotValue
	}
	return reward
}

// canChangeDickSize - может ли пользователь изменить пенис сегодня. (остались ли у него попытки)
// Обновляет попытки каждый день до 0.
func (d dickStartExec) canChangeDickSize(user *storage.DBUser, db storage.Storage) (bool, error) {
	if d.isNewDay(user.ChangeDickAt) {
		user.CurDickChangeCount = 0
		if err := db.UpdateUser(context.Background(), user); err != nil {
			return false, fmt.Errorf("can't update user in 'canChangeDickSize': %w", err)
		}
	}
	if user.CurDickChangeCount < user.MaxDickChangeCount {
		user.CurDickChangeCount++
		if err := db.UpdateUser(context.Background(), user); err != nil {
			return false, fmt.Errorf("can't update user in 'canChangeDickSize': %w", err)
		}
		return true, nil
	}
	return false, nil
}

func (d dickStartExec) isNewDay(lastChange time.Time) bool {
	yearLastTry, monthLastTry, dayLastTry := lastChange.Date()
	year, month, today := time.Now().Date()
	return (month == monthLastTry && today > dayLastTry) || month > monthLastTry || year > yearLastTry
}

// reward возвращает случайное положительное или отрицательное число в конкретном диапозоне.
func (d dickStartExec) reward() int {
	reward := minReward + rand.Intn(maxReward)
	if !d.isPlus() {
		reward = -reward
	}
	return reward
}

func (d dickStartExec) isPlus() bool {
	return rand.Intn(100) > 35
}

// IsJackpot показывает выиграл ли пользователь джекпот.
func (d dickStartExec) isJackpot() bool {
	return rand.Intn(100) == 69
}

type finishSeasonExec string

func (f finishSeasonExec) Exec(p *Processor, inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {

	if !p.isAdmin(user.ID) {
		return nil, e.Wrap("no admin can't do this cmd (/send_message)", ErrNotAdmin)
	}

	strs := strings.Split(inMessage, " ")
	if len(strs) < 2 {
		return nil, errors.New("invalid input message")
	}
	chatIDStr := strs[1]
	chatID, err := strconv.Atoi(chatIDStr)
	if err != nil {
		p.logger.Error("invalid type of chat id", slog.Any("error", err))
		return nil, e.Wrap("invalid chat ID", err)
	}

	users, err := p.storage.UsersByChat(context.Background(), chatID)
	if err != nil {
		return nil, e.Wrap("[ERROR] can't get users: ", err)
	}

	topDicksMessage, err := getTopDicks(chatID, p)
	if err != nil {
		return nil, e.Wrap("can't get top dics", err)
	}

	if topDicksMessage == "" {
		return &Response{message: msgErrorZeroUsersInSeason, method: sendMessageMethod, parseMode: telegram.Markdown, replyMessageId: messageID}, nil
	}

	winner := users[0]

	err = f.processFinishSeason(p, users, winner)
	if err != nil {
		return nil, e.Wrap("can't proccess finish season", err)
	}

	resultMessages := []string{msgEndSeason, fmt.Sprintf(msgSeasonResult, topDicksMessage), msgStartSeason}

	resultMessage := strings.Join(resultMessages, "\n")

	err = p.tg.SendMessage(chatID, resultMessage, telegram.Markdown, -1)
	if err != nil {
		p.logger.Error("can't send message in finish season command", slog.Any("error", err))
		return &Response{message: msgError, method: sendMessageMethod, parseMode: telegram.Markdown, replyMessageId: messageID}, nil
	}
	return &Response{message: msgSuccess, method: sendMessageMethod, parseMode: telegram.Markdown, replyMessageId: messageID}, nil
}

func (f finishSeasonExec) processFinishSeason(p *Processor, users []*storage.DBUser, winner *storage.DBUser) error {
	winner.MaxDickChangeCount++

	for _, user := range users {
		user.DickSize = 0
		user.ChangeDickAt = time.Now().Add(-24 * time.Hour)
		err := p.storage.UpdateUser(context.Background(), user)
		if err != nil {
			return e.Wrap(fmt.Sprintf("can't update user %s", user.Username), err)
		}
	}

	return nil
}

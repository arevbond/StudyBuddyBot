package telegram

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"tg_ics_useful_bot/clients/telegram"
	"tg_ics_useful_bot/lib/e"
	"tg_ics_useful_bot/storage"
	"time"
)

const (
	jackpotValue = 50000
	minValue     = 500
)

// dickTopExec предоставляет метод Exec для выполнения /top_dick.
type dickTopExec string

// Exec: /top_dick - пишет топ всех пенисов в чат.
func (d dickTopExec) Exec(p *Processor, inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {

	message, err := d.getTopDicks(chat.ID, p)
	if err != nil {
		return nil, e.Wrap(fmt.Sprintf("can't get top dics from chat %d: ", chat.ID), err)
	}
	mthd := sendMessageMethod
	return &Response{message: message, method: mthd, replyMessageId: -1, parseMode: telegram.Markdown}, nil
}

// getTopDicks возвращает string сообщение со списком всех dick > 0 в чате.
func (d dickTopExec) getTopDicks(chatID int, p *Processor) (msg string, err error) {
	users, err := p.storage.UsersByChat(context.Background(), chatID)
	if err != nil {
		return "", e.Wrap("[ERROR] can't get users: ", err)
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

	message, err := d.gameDick(chat, user, userStats, p.storage)
	if err != nil {
		return nil, e.Wrap("can't get message from gameDick: ", err)
	}
	mthd := sendMessageMethod
	return &Response{message: message, method: mthd, replyMessageId: -1}, nil
}

// gameDick это функция изменяющая размер пениса на случайное число и время изменения пениса.
// /dick - command
// Возвращает сообщение, отправляемое в чат.
func (d dickStartExec) gameDick(chat *telegram.Chat, user *telegram.User, userStats *storage.DBUserStat, db storage.Storage) (msg string, err error) {
	defer func() { err = e.WrapIfErr("error in gameDick: ", err) }()

	dbUser, err := db.GetUser(context.Background(), user.ID, chat.ID)
	if err != nil {
		return "", err
	}

	message, err := d.proccessDickGame(dbUser, userStats, db)
	if err != nil {
		return "", e.Wrap("can't work game dick cmd", err)
	}
	return message, nil
}

func (d dickStartExec) proccessDickGame(dbUser *storage.DBUser, userStats *storage.DBUserStat, db storage.Storage) (string, error) {
	canChange, err := d.canChangeDickSize(dbUser, db)
	if err != nil {
		return "", err
	}

	if !canChange {
		return d.formatAlreadyPlaying(dbUser), nil
	}

	oldDickSize := dbUser.DickSize
	reward, err := d.updateRandomDickAndChangeTime(dbUser, userStats, db)
	if err != nil {
		return "", err
	}
	return d.formatOutputGameDick(dbUser, oldDickSize, reward), nil
}

func (d dickStartExec) formatOutputGameDick(dbUser *storage.DBUser, oldDickSize int, reward *dickReward) string {
	name, hasUsername := d.getName(dbUser)
	if oldDickSize == 0 {
		if hasUsername {
			return fmt.Sprintf(msgCreateUserWithUsername, name) + fmt.Sprintf(msgDickSize, dbUser.DickSize)
		}
		return fmt.Sprintf(msgCreateUserWithFullName, name) + fmt.Sprintf(msgDickSize, dbUser.DickSize)
	}
	if hasUsername {
		return fmt.Sprintf(msgChangeDickSizeWithUsername, name, oldDickSize, dbUser.DickSize, reward.basic, reward.basic+reward.maxAdditional, reward.chance)
	}
	return fmt.Sprintf(msgChangeDickSizeWithFullName, name, oldDickSize, dbUser.DickSize, reward.basic, reward.basic+reward.maxAdditional, reward.chance)
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
func (d dickStartExec) updateRandomDickAndChangeTime(user *storage.DBUser, userStats *storage.DBUserStat, db storage.Storage) (*dickReward, error) {
	var value int
	var reward *dickReward
	for {
		reward = d.reward(user.DickSize)
		value = reward.getReward()
		if user.DickSize+value > 0 {
			break
		}
	}

	if d.IsJackpot() {
		value = jackpotValue
	}

	if value > 0 {
		userStats.DickPlusCount++
	} else {
		userStats.DickMinusCount++
	}
	err := db.UpdateUserStats(context.Background(), userStats)
	if err != nil {
		log.Print(err)
	}

	user.DickSize += value
	user.ChangeDickAt = time.Now()
	err = db.UpdateUser(context.Background(), user)
	if err != nil {
		return nil, e.Wrap(fmt.Sprintf("chat id %d, user %s can't change dick size or change dick at: ", user.ChatID, user.Username), err)
	}
	return reward, nil
}

// canChangeDickSize - может ли пользователь изменить пенис сегодня. (остались ли у него попытки)
// Обновляет попытки каждый день до 0.
func (d dickStartExec) canChangeDickSize(user *storage.DBUser, db storage.Storage) (bool, error) {
	yearLastTry, monthLastTry, dayLastTry := user.ChangeDickAt.Date()
	year, month, today := time.Now().Date()
	if (month == monthLastTry && today > dayLastTry) || month > monthLastTry || year > yearLastTry {
		user.CurDickChangeCount = 0
		err := db.UpdateUser(context.Background(), user)
		if err != nil {
			return false, e.Wrap("can't update user in 'canChangeDickSize'", err)
		}
	}
	if user.CurDickChangeCount+1 <= user.MaxDickChangeCount {
		user.CurDickChangeCount++
		err := db.UpdateUser(context.Background(), user)
		if err != nil {
			return false, e.Wrap("can't update user in 'canChangeDickSize'", err)
		}
		return true, nil
	}

	return false, nil
}

// reward возвращает случайное положительное или отрицательное число в конкретном диапозоне.
func (d dickStartExec) reward(dickSize int) *dickReward {
	reward := d.getRewardByDickSize(dickSize)

	reward.isPlus = d.isPlus()
	return reward
}

type dickLevel int

const (
	lowDick dickLevel = iota
	highDick
)

type dickReward struct {
	basic         int
	maxAdditional int
	chance        int
	isPlus        bool
}

func (d dickReward) getReward() int {
	result := d.basic + rand.Intn(d.maxAdditional)
	if d.isPlus {
		return result
	}
	return -1 * result
}

func (d dickStartExec) getDickReward(dickSize int, level dickLevel) *dickReward {
	if level == lowDick {
		return &dickReward{basic: 10, maxAdditional: 100, chance: 100}
	}
	n := rand.Intn(101)
	switch {
	case n < 80:
		return &dickReward{basic: 100, maxAdditional: dickSize / 10, chance: 80}
	case n >= 80 && n < 95:
		return &dickReward{basic: 100, maxAdditional: dickSize / 5, chance: 14}
	case n >= 95 && n < 100:
		return &dickReward{basic: 100, maxAdditional: dickSize, chance: 5}
	case n == 100:
		return &dickReward{basic: 100, maxAdditional: 2 * dickSize, chance: 1}
	default:
		return &dickReward{1, 1, 1, true}
	}
}

func (d dickStartExec) getRewardByDickSize(dickSize int) *dickReward {
	level := d.DickLevel(dickSize)
	return d.getDickReward(dickSize, level)
}

func (d dickStartExec) DickLevel(dickSize int) dickLevel {
	switch {
	case dickSize < 100:
		return lowDick
	default:
		return highDick
	}
}

func (d dickStartExec) isPlus() bool {
	sign := rand.Intn(20)
	if sign <= 7 {
		return false
	}
	return true
}

// IsJackpot показывает выиграл ли пользователь джекпот.
func (d dickStartExec) IsJackpot() bool {
	if value := rand.Intn(100); value == 77 {
		return true
	}
	return false
}

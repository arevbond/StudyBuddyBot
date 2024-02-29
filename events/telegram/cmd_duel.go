package telegram

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"tg_ics_useful_bot/clients/telegram"
	"tg_ics_useful_bot/lib/e"
	"tg_ics_useful_bot/lib/utils"
	"tg_ics_useful_bot/storage"
	"time"
)

var duels = make(map[string]*storage.DBUser)

const (
	RewardForKill = 25
	HealthPoints  = 5

	HeartEmoji = "❤️"
	DeathEmoji = "☠️"
)

// getHpExec предоставляет Exec метод для выполнения /hp.
type getHpExec string

// Exec: /hp - один раз в день пополняет здоровье пользователя.
func (a getHpExec) Exec(p *Processor, inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {

	message, err := p.getHp(user, chat)
	if err != nil {
		return nil, err
	}
	mthd := sendMessageMethod
	return &Response{message: message, method: mthd}, nil
}

// duelExec предоставляет Exec метод для выполнения /duel.
type duelExec string

// Exec: /duel {@username} - игра дуели.
func (a duelExec) Exec(p *Processor, inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {

	message, err := p.gameDuel(chat, user, user.Username)
	if err != nil {
		return nil, e.Wrap("can't do gameDuel: ", err)
	}
	if utils.StringContains("@", inMessage) {
		textSplited := strings.Split(inMessage, "@")
		target := textSplited[len(textSplited)-1]
		log.Printf("[INFO] @%s вызывает на дуель @%s", user.Username, target)
		message, err = p.gameDuel(chat, user, target)
		if err != nil {
			return nil, e.Wrap("can't do gameDuel: ", err)
		}
	}
	mthd := sendMessageMethod
	return &Response{message: message, method: mthd, replyMessageId: -1}, nil
}

// getHp пополняет HP пользователя раз в день.
func (p *Processor) getHp(user *telegram.User, chat *telegram.Chat) (string, error) {
	dbUser, err := p.storage.GetUser(context.Background(), user.ID, chat.ID)
	if err != nil {
		return "", err
	}

	if !p.canGetHp(dbUser) {
		return fmt.Sprintf(msgCantGetHP, user.Username, p.hpString(dbUser)), nil
	}

	dbUser.HpTakedAt = time.Now()

	if dbUser.HealthPoints < defaultHpUser {
		dbUser.HealthPoints += defaultHpUser
	} else {
		dbUser.HealthPoints += 1
	}
	err = p.storage.UpdateUser(context.Background(), dbUser)
	if err != nil {
		return "", e.Wrap("can't update hp in 'canChangeDickSize'", err)
	}
	return fmt.Sprintf(msgGetHp, dbUser.Username, p.hpString(dbUser)), nil
}

// gameDuel проводит дуель между двумя участиками чата на оснвое их DickSize и HP.
func (p *Processor) gameDuel(chat *telegram.Chat, user *telegram.User, targetUsername string) (string, error) {
	u1, err := p.storage.GetUser(context.Background(), user.ID, chat.ID)
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

	if !p.canDuel(u1, u2) {
		return fmt.Sprintf(msgCantCreateDuel, u1.Username, p.hpString(u1), u2.Username, p.hpString(u2)), nil
	}

	stats1, err := p.storage.GetUserStats(context.Background(), u1)
	if err != nil {
		return "", e.Wrap("can't get user stats in 'gameDuel'", err)
	}
	stats2, err := p.storage.GetUserStats(context.Background(), u2)
	if err != nil {
		return "", e.Wrap("can't get user stats in 'gameDuel'", err)
	}

	oldDickSize1 := u1.DickSize
	oldDickSize2 := u2.DickSize

	oldHP1 := p.hpString(u1)
	oldHP2 := p.hpString(u2)

	finishMessage := msgFinishDuel
	if enemy, ok := duels[u1.Username]; ok && enemy.TgID == u2.TgID {
		delete(duels, u1.Username)
		stats1.DuelsCount++
		stats2.DuelsCount++

		isUser1Win, ch1, ch2 := duel(u1.DickSize, u2.DickSize)
		if isUser1Win {
			stats1.DuelsWinCount++
			stats2.DuelsLoseCount++

			err = p.changeHP(u2, -1)
			if err != nil {
				return "", err
			}

			reward := getReward(u2.DickSize, ch1)

			if p.isDie(u2) {
				stats1.KillCount++
				stats2.DieCount++

				reward += RewardForKill

				finishMessage = msgPlayerDie
			}

			err = p.changeDickSize(u1, reward)
			if err != nil {
				return "", err
			}
			err = p.changeDickSize(u2, -1*reward)
			if err != nil {
				return "", err
			}

			err1 := p.storage.UpdateUserStats(context.Background(), stats1)
			err2 := p.storage.UpdateUserStats(context.Background(), stats2)
			if err1 != nil || err2 != nil {
				log.Println("[ERROR] can't update user stats in 'gameDuel'")
			}

			return fmt.Sprintf(msgAcceptDuel, u1.Username, oldHP1, oldDickSize1, ch1, u2.Username, oldHP2, oldDickSize2, ch2) +
				fmt.Sprintf(finishMessage, u1.Username, p.hpString(u1), u1.DickSize, reward, u2.Username, p.hpString(u2),
					u2.DickSize, reward), nil
		} else {
			stats2.DuelsWinCount++
			stats1.DuelsLoseCount++

			err = p.changeHP(u1, -1)
			if err != nil {
				return "", err
			}

			reward := getReward(u1.DickSize, ch2)
			if p.isDie(u1) {
				stats2.KillCount++
				stats1.DieCount++

				reward += RewardForKill

				finishMessage = msgPlayerDie
			}

			err = p.changeDickSize(u1, -1*reward)
			if err != nil {
				return "", err
			}
			err = p.changeDickSize(u2, reward)
			if err != nil {
				return "", err
			}

			err1 := p.storage.UpdateUserStats(context.Background(), stats1)
			err2 := p.storage.UpdateUserStats(context.Background(), stats2)
			if err1 != nil || err2 != nil {
				log.Println("[ERROR] can't update user stats in 'gameDuel'")
			}

			return fmt.Sprintf(msgAcceptDuel, u1.Username, oldHP1, oldDickSize1, ch1, u2.Username, oldHP2, oldDickSize2, ch2) +
				fmt.Sprintf(finishMessage, u2.Username, p.hpString(u2), u2.DickSize, reward, u1.Username, p.hpString(u1),
					u1.DickSize, reward), nil
		}
	} else {
		duels[targetUsername] = u1
		return fmt.Sprintf(msgChallengeToDuel, u1.Username, targetUsername), nil
	}
}

// hpString возвращает unicode строку, в которой кол-во hp пользователя
// конвертируется в строку с сердечками.
func (p *Processor) hpString(user *storage.DBUser) string {
	heart := HeartEmoji
	result := ""
	for i := 1; i <= user.HealthPoints; i++ {
		result += heart
	}
	if len(result) == 0 {
		return DeathEmoji
	}
	return result
}

// changeDickSize изменяет размер пениса после дуели.
// Не позволяет размеру пениса быть меньше 1.
func (p *Processor) changeDickSize(user *storage.DBUser, value int) error {
	user.DickSize += value
	if user.DickSize <= 0 {
		user.DickSize = 1
	}
	err := p.storage.UpdateUser(context.Background(), user)
	if err != nil {
		return e.Wrap(fmt.Sprintf("chat id %d, user %s can't change dick size :", user.ChatID, user.Username), err)
	}
	return nil
}

// changeHP изменяет значение health_points пользователя в базе данных.
func (p *Processor) changeHP(user *storage.DBUser, value int) error {
	user.HealthPoints += value
	err := p.storage.UpdateUser(context.Background(), user)
	if err != nil {
		return e.Wrap(fmt.Sprintf("chat id %d, user %s can't change health points :", user.ChatID, user.Username), err)
	}
	return nil
}

// isDie возвращает равно ли 0 хп пользователя.
func (p *Processor) isDie(user *storage.DBUser) bool {
	return user.HealthPoints == 0
}

// canDuel возвращает имеют ли два пользователя больше 0 хп или имеют ли они писюны.
func (p *Processor) canDuel(user1 *storage.DBUser, user2 *storage.DBUser) bool {
	return (user1.HealthPoints > 0 && user2.HealthPoints > 0) && (user1.DickSize > 0 && user2.DickSize > 0)
}

// canGetHp возвращает может ли пользватель сегодня пополнить хп.
func (p *Processor) canGetHp(user *storage.DBUser) bool {
	yearLastTry, monthLastTry, dayLastTry := user.HpTakedAt.Date()
	year, month, today := time.Now().Date()
	return ((month == monthLastTry && today > dayLastTry) || (month > monthLastTry || year > yearLastTry)) && (user.HealthPoints < HealthPoints)

}

// duel return true if dick1 wins.
func duel(dick1 int, dick2 int) (bool, float64, float64) {
	allChance := dick1 + dick2
	chance1 := float64(dick1) / float64(allChance) * 100
	chance2 := float64(dick2) / float64(allChance) * 100

	result := float64(rand.Intn(100))
	log.Printf("[INFO] duel between dick1 = %d and dick2 = %d. chance1 = %.2f and chance2 = %.2f", dick1, dick2, chance1, chance2)
	return result <= chance1, chance1, chance2
}

// getReward считает награду при дуели.
// (enemyDick * 1/10) * (1 - chance %)
func getReward(enemyDick int, chance float64) int {
	reward := enemyDick / 10
	return int(float64(reward) * (1 - (chance / 100)))
}

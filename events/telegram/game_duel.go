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

var duels = make(map[string]*storage.DBUser)

const (
	REWARD_FOR_KILL = 25
)

// getHpCmd пополняет HP пользователя раз в день.
func (p *Processor) getHpCmd(user *telegram.User, chat *telegram.Chat) (string, error) {
	dbUser, err := p.storage.GetUser(context.Background(), user.ID, chat.ID)
	if err != nil {
		return "", err
	}

	if !p.canGetHp(dbUser) {
		return fmt.Sprintf(msgCantGetHP, dbUser.Username), nil
	}

	dbUser.HpTakedAt = time.Now()

	if dbUser.HealthPoints < DEFAULT_HP_USER {
		dbUser.HealthPoints += DEFAULT_HP_USER
	} else {
		dbUser.HealthPoints += 1
	}
	err = p.storage.UpdateUser(context.Background(), dbUser)
	if err != nil {
		return "", e.Wrap("can't update hp in 'canChangeDickSize'", err)
	}
	return fmt.Sprintf(msgGetHp, dbUser.Username, p.hpString(dbUser)), nil
}

// gameDuelCmd проводит дуель между двумя участиками чата на оснвое их DickSize и HP.
func (p *Processor) gameDuelCmd(chat *telegram.Chat, user *telegram.User, targetUsername string) (string, error) {
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
		return fmt.Sprintf(msgCantCreateDuel, u1.Username, u2.Username), nil
	}

	stats1, err := p.storage.GetUserStats(context.Background(), u1)
	if err != nil {
		return "", e.Wrap("can't get user stats in 'gameDuelCmd'", err)
	}
	stats2, err := p.storage.GetUserStats(context.Background(), u2)
	if err != nil {
		return "", e.Wrap("can't get user stats in 'gameDuelCmd'", err)
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

				reward += REWARD_FOR_KILL

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
				log.Println("[ERROR] can't update user stats in 'gameDuelCmd'")
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

				reward += REWARD_FOR_KILL

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
				log.Println("[ERROR] can't update user stats in 'gameDuelCmd'")
			}

			return fmt.Sprintf(msgAcceptDuel, u1.Username, oldHP1, oldDickSize1, ch1, u2.Username, oldHP2, oldDickSize2, ch2) +
				fmt.Sprintf(finishMessage, u2.Username, p.hpString(u1), u2.DickSize, reward, u1.Username, p.hpString(u2),
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
	heart := "❤️"
	result := ""
	for i := 1; i <= user.HealthPoints; i++ {
		result += heart
	}
	if len(result) == 0 {
		return "☠️"
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

// canDuel возвращает имеют ли два пользователя больше 0 хп.
func (p *Processor) canDuel(user1 *storage.DBUser, user2 *storage.DBUser) bool {
	return user1.HealthPoints > 0 && user2.HealthPoints > 0
}

// canGetHp возвращает может ли пользватель сегодня пополнить хп.
func (p *Processor) canGetHp(user *storage.DBUser) bool {
	yearLastTry, monthLastTry, dayLastTry := user.HpTakedAt.Date()
	year, month, today := time.Now().Date()
	return (month == monthLastTry && today > dayLastTry) || month > monthLastTry || year > yearLastTry

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

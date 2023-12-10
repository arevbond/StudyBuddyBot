package telegram

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"tg_ics_useful_bot/clients/telegram"
	"tg_ics_useful_bot/lib/e"
	"tg_ics_useful_bot/storage"
)

var duels = make(map[string]*storage.DBUser)

// TODO: добавить уменьшение HP
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

	oldDickSize1 := u1.DickSize
	oldDickSize2 := u2.DickSize

	if enemy, ok := duels[u1.Username]; ok && enemy.TgID == u2.TgID {
		delete(duels, u1.Username)
		User1Win, ch1, ch2 := duel(u1.DickSize, u2.DickSize)
		if User1Win {
			reward := getReward(u2.DickSize, ch1)
			err2 := p.changeDickSize(u1, reward)
			if err2 != nil {
				return "", err
			}
			err3 := p.changeDickSize(u2, -1*reward)
			if err3 != nil {
				return "", err3
			}
			return fmt.Sprintf(msgAcceptDuel, u1.Username, p.hpString(u1), oldDickSize1, ch1, u2.Username, p.hpString(u2), oldDickSize2, ch2) +
				fmt.Sprintf(msgFinishDuel, u1.Username, u1.DickSize, reward, u2.Username, u2.DickSize, reward), nil
		} else {
			reward := getReward(u1.DickSize, ch2)
			err2 := p.changeDickSize(u1, -1*reward)
			if err2 != nil {
				return "", err
			}
			err3 := p.changeDickSize(u2, reward)
			if err3 != nil {
				return "", err3
			}
			return fmt.Sprintf(msgAcceptDuel, u1.Username, p.hpString(u1), oldDickSize1, ch1, u2.Username, p.hpString(u2), oldDickSize2, ch2) +
				fmt.Sprintf(msgFinishDuel, u2.Username, u2.DickSize, reward, u1.Username, u1.DickSize, reward), nil
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

// duel return true if dick1 wins.
func duel(dick1 int, dick2 int) (bool, float64, float64) {
	allChance := dick1 + dick2
	chance1 := float64(dick1) / float64(allChance) * 100
	chance2 := float64(dick2) / float64(allChance) * 100

	result := float64(rand.Intn(100))
	log.Printf("[INFO] duel between dick1 = %d and dick2 = %d. chance1 = %f and chance2 = %f", dick1, dick2, chance1, chance2)
	return result <= chance1, chance1, chance2
}

// getReward считает награду при дуели.
// (enemyDick * 1/10) * (1 - chance %)
func getReward(enemyDick int, chance float64) int {
	reward := enemyDick / 10
	return int(float64(reward) * (1 - (chance / 100)))
}

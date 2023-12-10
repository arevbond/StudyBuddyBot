package telegram

import (
	"context"
	"fmt"
	"tg_ics_useful_bot/clients/telegram"
	"tg_ics_useful_bot/lib/duel"
	"tg_ics_useful_bot/lib/e"
	"tg_ics_useful_bot/storage"
)

var duels = make(map[string]*storage.DBUser)

var reward = 15

// TODO: изменить формулу для награды
func (p *Processor) gameDuelDick(chat *telegram.Chat, messageID int, user *telegram.User, targetUsername string) (string, error) {
	err := p.tg.DeleteMessage(chat.ID, messageID)
	if err != nil {
		return "", e.Wrap(fmt.Sprintf("[ERROR] can't delete message: user #%d, chat id #%d", user.ID, chat.ID), err)
	}

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
		User1Win, ch1, ch2 := duel.Duel(u1.DickSize, u2.DickSize)
		if User1Win {
			if ch1 > 65 {
				reward = 5
			}
			err2 := p.changeDickSizeAndTime(u1, reward)
			if err2 != nil {
				return "", err
			}
			err3 := p.changeDickSizeAndTime(u2, -1*reward)
			if err3 != nil {
				return "", err3
			}
			return fmt.Sprintf(msgAcceptDuel, u1.Username, oldDickSize1, ch1, u2.Username, oldDickSize2, ch2) +
				fmt.Sprintf(msgFinishDuel, u1.Username, u1.DickSize, u2.Username, u2.DickSize), nil
		} else {
			if ch1 <= 35 {
				reward = 5
			}
			err2 := p.changeDickSizeAndTime(u1, -1*reward)
			if err2 != nil {
				return "", err
			}
			err3 := p.changeDickSizeAndTime(u2, reward)
			if err3 != nil {
				return "", err3
			}
			return fmt.Sprintf(msgAcceptDuel, u1.Username, oldDickSize1, ch1, u2.Username, oldDickSize2, ch2) +
				fmt.Sprintf(msgFinishDuel, u2.Username, u2.DickSize, u1.Username, u1.DickSize), nil
		}
	} else {
		duels[targetUsername] = u1
		return fmt.Sprintf(msgChallengeToDuel, u1.Username, targetUsername), nil
	}
}

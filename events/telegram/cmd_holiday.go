package telegram

import (
	"context"
	"fmt"
	"math/rand"
	"tg_ics_useful_bot/clients/telegram"
	"tg_ics_useful_bot/lib/e"
	"tg_ics_useful_bot/storage"
	"time"
)

var holidays = make(map[int]struct{})

type holidayExec string

func (a holidayExec) Exec(p *Processor, inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {

	mthd := sendMessageMethod
	if _, ok := holidays[user.ID]; ok {
		return &Response{message: msgAlreadyUsedHoliday, method: mthd}, nil
	}

	var message string

	dbUser, err := p.storage.GetUser(context.Background(), user.ID, chat.ID)
	if err != nil {
		return nil, err
	}
	holidays[user.ID] = struct{}{}

	for i := 3; i > 0; i-- {
		p.tg.SendMessage(chat.ID, fmt.Sprintf("До результата Новогоднего Розыгрыша: %d", i), "", -1)
		time.Sleep(1 * time.Second)
	}

	oldDickSize := dbUser.DickSize
	if isWin() {
		err = p.changeDickSize(dbUser, oldDickSize)
		message = msgWinHoliday
		if err != nil {
			return nil, e.Wrap("can't change dick size", err)
		}
	} else {
		err = p.changeDickSize(dbUser, -oldDickSize/2)
		message = msgLoseHoliday
		if err != nil {
			return nil, e.Wrap("can't change dick size", err)
		}
	}

	return &Response{message: message, method: mthd, replyMessageId: messageID}, nil
}

func isWin() bool {
	return rand.Intn(5) > 1
}

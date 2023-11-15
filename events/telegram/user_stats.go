package telegram

import (
	"context"
	"log"
)

func (p *Processor) chatStats(chatID int) (int, int, int, int, int) {
	users, err := p.storage.UsersStatsByChatID(context.Background(), chatID)
	if err != nil {
		log.Print(err)
		return -1, -1, -1, -1, -1
	}
	var msgCnt, dickPlusCnt, dickMinusCnt, yesCnt, noCnt int
	for _, u := range users {
		msgCnt += u.MessageCount
		dickPlusCnt += u.DickPlusCount
		dickMinusCnt += u.DickMinusCount
		yesCnt += u.YesCount
		noCnt += u.NoCount
	}
	return msgCnt, dickMinusCnt, dickPlusCnt, yesCnt, noCnt
}

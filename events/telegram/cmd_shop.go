package telegram

import (
	"tg_ics_useful_bot/clients/telegram"
	"tg_ics_useful_bot/storage"
)

const (
	priceHp = 1500

	priceDickSpin  = 1500
	limitDickSpins = 3
)

type shopExec string

func (s shopExec) Exec(p *Processor, inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {

	return nil, nil
}

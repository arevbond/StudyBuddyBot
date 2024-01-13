package telegram

import (
	"tg_ics_useful_bot/clients/telegram"
	"tg_ics_useful_bot/storage"
)

type shopExec string

func (s shopExec) Exec(p *Processor, inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {

	return nil, nil
}
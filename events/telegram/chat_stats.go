package telegram

import (
	"context"
	"tg_ics_useful_bot/storage"
)

func (p *Processor) chatStats(chatID int) (*storage.DBUserStat, error) {
	users, err := p.storage.UsersByChat(context.Background(), chatID)
	if err != nil {
		return nil, err
	}
	allStats := &storage.DBUserStat{}
	for _, u := range users {
		userStats, err := p.storage.GetUserStats(context.Background(), u)
		if err != nil {
			return nil, err
		}
		allStats.MessageCount += userStats.MessageCount
		allStats.DickPlusCount += userStats.DickMinusCount
		allStats.DickMinusCount += userStats.DickMinusCount
		allStats.YesCount += userStats.YesCount
		allStats.NoCount += userStats.NoCount
		allStats.DuelsCount += userStats.DuelsCount
		allStats.DuelsWinCount += userStats.DuelsWinCount
		allStats.DuelsLoseCount += userStats.DuelsLoseCount
		allStats.KillCount += userStats.KillCount
		allStats.DieCount += userStats.KillCount
	}
	return allStats, nil
}

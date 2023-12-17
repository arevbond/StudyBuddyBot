package telegram

import (
	"context"
	"fmt"
	"tg_ics_useful_bot/clients/telegram"
	"tg_ics_useful_bot/lib/e"
	"tg_ics_useful_bot/storage"
)

// myStatsExec предоставляет Exec метод для выполнения /my_stats.
type myStatsExec string

// Exec: /my_stats - возвращает статистику пользователя в данном чате.
func (a myStatsExec) Exec(p *Processor, inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {
	message := fmt.Sprintf(msgUserStats, userStats.MessageCount, userStats.DickPlusCount,
		userStats.DickMinusCount, userStats.YesCount, userStats.NoCount, userStats.DuelsCount,
		userStats.DuelsWinCount, userStats.DuelsLoseCount, userStats.KillCount, userStats.DieCount)
	mthd := sendMessageMethod
	replyMessageId := messageID
	return &Response{message: message, method: mthd, replyMessageId: replyMessageId}, nil
}

// chatStatsExec предоставляет Exec метод для выполнения /chat_stats.
type chatStatsExec string

// Exec: /chat_stats - возвращает всю статистику данного чата.
func (a chatStatsExec) Exec(p *Processor, inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {
	userStats, err := p.chatStats(chat.ID)
	if err != nil {
		return nil, e.Wrap("can't get chat stats: ", err)
	}
	message := fmt.Sprintf(msgUserStats, userStats.MessageCount, userStats.DickPlusCount,
		userStats.DickMinusCount, userStats.YesCount, userStats.NoCount, userStats.DuelsCount,
		userStats.DuelsWinCount, userStats.DuelsLoseCount, userStats.KillCount, userStats.DieCount)
	replyMessageId := messageID
	mthd := sendMessageMethod
	return &Response{message: message, method: mthd, replyMessageId: replyMessageId}, nil
}

// chatStats формирует статистику чата, суммирая все статистики пользователей данного чата.
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

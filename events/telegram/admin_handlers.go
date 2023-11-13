package telegram

import (
	"context"
	"log"
	"strconv"
)

func (p *Processor) changeAnyDickSize(chatIDStr, userIDStr, valueStr string) error {
	chatID, err := strconv.Atoi(chatIDStr)
	if err != nil {
		return err
	}
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return err
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return err
	}
	dbUser, err := p.storage.UserByTelegramID(context.Background(), userID, chatID)
	err = p.storage.UpdateUserDickSize(context.Background(), dbUser, value)
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

package telegram

import (
	"context"
	"log"
	"strconv"
)

// changeDickByAdminCmd админская ручка, позволяющая изменить пенис любому пользователю.
func (p *Processor) changeDickByAdminCmd(chatIDStr, userIDStr, valueStr string) error {
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
	dbUser, err := p.storage.GetUser(context.Background(), userID, chatID)
	if err != nil {
		return err
	}
	dbUser.DickSize += value
	err = p.storage.UpdateUser(context.Background(), dbUser)
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

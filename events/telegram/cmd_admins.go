package telegram

import (
	"context"
	"log"
	"strconv"
	"strings"
	"tg_ics_useful_bot/clients/telegram"
	"tg_ics_useful_bot/storage"
)

// adminSendMessageExec предоставляет метод Exec для отправки сообщения от имени бота в любой чат
type adminSendMessageExec struct {
	command string
}

// Exec: /send_message {chat_id} {message}
func (a *adminSendMessageExec) Exec(p *Processor, inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {

	strs := strings.Split(inMessage, " ")
	chatIDStr, message := strs[1], strings.Join(strs[2:], " ")
	chatID, err := strconv.Atoi(chatIDStr)
	if err != nil {
		log.Print(err)
	}
	err = p.tg.SendMessage(chatID, message, "", -1)
	if err != nil {
		log.Println("can't send message by admin:", err)
	}
	mthd := doNothingMethod
	return &Response{message: message, method: mthd, replyMessageId: -1}, nil
}

// adminChangeDickExec предоставляет метод Exec для выполнения команды /change_dick
type adminChangeDickExec struct {
	command string
}

// Exec: /change_dick {chat_id} {user_id} {value}
func (a *adminChangeDickExec) Exec(p *Processor, inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {

	strs := strings.Split(inMessage, " ")
	chatIDStr, userIDStr, valueStr := strs[1], strs[2], strs[3]
	err := p.changeDickByAdminCmd(chatIDStr, userIDStr, valueStr)
	if err != nil {
		return nil, err
	}
	message := msgSuccessAdminChangeDickSize
	mthd := sendMessageMethod
	return &Response{message: message, method: mthd, replyMessageId: -1}, nil
}

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

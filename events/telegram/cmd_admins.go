package telegram

import (
	"context"
	"errors"
	"log/slog"
	"strconv"
	"strings"
	"tg_ics_useful_bot/clients/telegram"
	"tg_ics_useful_bot/lib/e"
	"tg_ics_useful_bot/storage"
)

// adminSendMessageExec предоставляет метод Exec для отправки сообщения от имени бота в любой чат
// Только для админов бота.
type adminSendMessageExec string

// Exec: /send_message {chat_id} {message}
func (a adminSendMessageExec) Exec(p *Processor, inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {

	if !p.isAdmin(user.ID) {
		return nil, e.Wrap("no admin can't do this cmd (/send_message)", errors.New("can't do this cmd"))
	}

	strs := strings.Split(inMessage, " ")
	chatIDStr, message := strs[1], strings.Join(strs[2:], " ")
	chatID, err := strconv.Atoi(chatIDStr)
	if err != nil {
		p.logger.Error("invalid type chat id", slog.Any("error", err))
	}
	err = p.tg.SendMessage(chatID, message, "", -1)
	if err != nil {
		p.logger.Error("can't send message by admin", slog.Any("error", err))
	}

	message = msgSuccess
	mthd := sendMessageMethod
	return &Response{message: message, method: mthd, replyMessageId: messageID}, nil
}

// adminChangeDickExec предоставляет метод Exec для выполнения команды /change_dick.
// Только для админов бота.
type adminChangeDickExec string

// Exec: /change_dick {chat_id} {user_id} {value}
func (a adminChangeDickExec) Exec(p *Processor, inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {

	if !p.isAdmin(user.ID) {
		return nil, e.Wrap("no admin can't do this cmd (/send_message)", errors.New("can't do this cmd"))
	}

	strs := strings.Split(inMessage, " ")
	chatIDStr, userIDStr, valueStr := strs[1], strs[2], strs[3]
	err := a.changeDickByAdmin(chatIDStr, userIDStr, valueStr, p.storage)
	if err != nil {
		return nil, err
	}
	message := msgSuccess
	mthd := sendMessageMethod
	return &Response{message: message, method: mthd, replyMessageId: messageID}, nil
}

// changeDickByAdmin админская ручка, позволяющая изменить пенис любому пользователю.
func (a adminChangeDickExec) changeDickByAdmin(chatIDStr, userIDStr, valueStr string, db storage.Storage) error {
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
	dbUser, err := db.GetUser(context.Background(), userID, chatID)
	if err != nil {
		return err
	}
	dbUser.DickSize += value
	err = db.UpdateUser(context.Background(), dbUser)
	if err != nil {
		return e.Wrap("can't update user", err)
	}
	return nil
}

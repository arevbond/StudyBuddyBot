package telegram

import (
	"log/slog"
	"strconv"
	"tg_ics_useful_bot/clients/telegram"
	"tg_ics_useful_bot/storage"
)

// helpExec предоставляет метод Exec для выполнения /help.
type helpExec string

// Exec: /help - возвращает help сообщениею
func (h helpExec) Exec(p *Processor, inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {

	message := msgHelp
	mthd := sendMessageMethod
	return &Response{message: message, method: mthd, replyMessageId: -1, parseMode: telegram.Markdown}, nil
}

// chatIDExec предоставляет метод Exec для выполнения /chat_id.
type chatIDExec string

// Exec: /chat_id - возвращает chat id.
func (c chatIDExec) Exec(p *Processor, inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {

	message := strconv.Itoa(chat.ID)
	mthd := sendMessageMethod
	return &Response{message: message, method: mthd, replyMessageId: -1}, nil
}

// allUsernamesExec предоставляет метод Exec для вызова всех админов в чате.
type allUsernamesExec string

// Exec: /all - тэгает всех админов в чате.
func (a allUsernamesExec) Exec(p *Processor, inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {

	message := a.allUsernames(chat.ID, p.tg, p.logger)
	mthd := sendMessageMethod
	return &Response{message: message, method: mthd, replyMessageId: -1}, nil
}

// allUsernames возвращает строку "@username1, @username2...".
func (a allUsernamesExec) allUsernames(chatID int, tgClient *telegram.Client, logger *slog.Logger) string {
	admins, err := tgClient.ChatAdministrators(chatID)
	if err != nil {
		logger.Error("can't get admins", slog.Any("error", err), slog.Int("chat id", chatID))
	}
	result := ""
	for _, a := range admins {
		result += "@" + a.Username + " "
	}
	return result[:len(result)-1]
}

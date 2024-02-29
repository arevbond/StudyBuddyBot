package telegram

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"tg_ics_useful_bot/clients/telegram"
	"tg_ics_useful_bot/lib/e"
	"tg_ics_useful_bot/lib/utils"
	"tg_ics_useful_bot/storage"
)

type method int

const (
	UnsupportedMethod method = iota
	sendMessageMethod
	sendPhotoMethod
	sendMessageWithButtonsMethod
	sendPoll
	doNothingMethod

	maxDickChangeCount = 1
	defaultHpUser      = 3

	answerOnYes = "Пизда"
	answerOnNo  = "Пидора ответ"
)

// CmdExecutor предоставляет интерфейс с методом Exec
// для процедуры выполнения команды пользователя.
type CmdExecutor interface {
	Exec(p *Processor, inMessage string, user *telegram.User, chat *telegram.Chat,
		userStats *storage.DBUserStat, messageID int) (*Response, error)
}

// Response структура содержащая поля для отправки сообщения пользователю.
type Response struct {
	message        string
	method         method
	parseMode      telegram.ParseMode
	replyMessageId int
	poll           telegram.SendPoll
}

// doCmd выбирает необходимую логику для выолнения команды.
func (p *Processor) doCmd(text string, chat *telegram.Chat, user *telegram.User, messageID int) error {
	ctx := context.Background()
	dbUser, err := p.getUser(ctx, user, chat.ID)
	userStats, err := p.storage.GetUserStats(context.Background(), dbUser)
	if err == storage.ErrUserNotExist {
		return e.Wrap("not find user stats", err)
	} else if err != nil {
		return e.Wrap("not user stats", err)
	}

	userStats.MessageCount++
	err = p.storage.UpdateUserStats(context.Background(), userStats)
	if err != nil {
		log.Print(err)
	}

	text, parseMode := strings.TrimSpace(text), telegram.WithoutParseMode

	switch utils.CheckYesOrNo(text) {
	case utils.IsYesCommand:

		userStats.YesCount++
		err = p.storage.UpdateUserStats(context.Background(), userStats)
		if err != nil {
			log.Print(err)
		}
		return p.tg.SendMessage(chat.ID, answerOnYes, parseMode, messageID)
	case utils.IsNoCommand:

		userStats.NoCount++
		err = p.storage.UpdateUserStats(context.Background(), userStats)
		if err != nil {
			log.Print(err)
		}
		return p.tg.SendMessage(chat.ID, answerOnNo, parseMode, messageID)
	}

	userWithChat := UserWithChat{chat.ID, user.ID}

	if _, ok := stateHomework[userWithChat]; ok {
		msg := p.addHomeworkCmd(text, userWithChat)
		replyToMessageID := messageID
		return p.tg.SendMessage(chat.ID, msg, parseMode, replyToMessageID)
	}

	if utils.IsCommand(text) {
		return p.handleCommand(text, chat, user, messageID, userStats)
	}
	return nil
}

func (p *Processor) getUser(ctx context.Context, user *telegram.User, chatID int) (*storage.DBUser, error) {
	dbUser, err := p.userCache.GetUser(user.ID, chatID)
	if err == storage.ErrUserNotExist {
		dbUser, err = p.storage.GetUser(ctx, user.ID, chatID)
		if err == storage.ErrUserNotExist {
			dbUser, err = p.createNewUserInDB(chatID, user)
			if err != nil {
				return nil, e.Wrap("can't create new user in 'getUser'", err)
			}
		} else if err != nil {
			return nil, e.Wrap("unknown error in 'doCmd'", err)
		}
		p.userCache.AddUser(dbUser)
		log.Printf("[INFO] Get user %d from storage\n", user.ID)
	}

	dbUser, err = p.userChangeInfo(user, dbUser)
	if err != nil {
		return nil, e.Wrap("can't update user info in 'doCmd'", err)
	}
	return dbUser, nil
}

func (p *Processor) handleCommand(text string, chat *telegram.Chat, user *telegram.User, messageID int, userStats *storage.DBUserStat) error {
	log.Printf("[INFO] got new command '%s' from '%s' in '%s'", text, user.Username, chat.Title)

	strCmd := strings.Split(text, " ")[0]
	cmd := p.getCmd(strCmd)
	if cmd == nil {
		return e.Wrap(fmt.Sprintf("can't get command from %s", strCmd), errors.New("can't get cmd"))
	}

	response, err := cmd.Exec(p, text, user, chat, userStats, messageID)
	if err != nil {
		return e.Wrap(fmt.Sprintf("can't select command from message: %s", text), err)
	}

	msg, mthd, parseMode, replyToMessageID := response.message, response.method, response.parseMode, response.replyMessageId
	if replyToMessageID <= 0 {
		err = p.tg.DeleteMessage(chat.ID, messageID)
		if err != nil {
			return err
		}
	}

	switch mthd {
	case UnsupportedMethod:
		return e.Wrap("unsupported method:", errors.New("unknown method"))
	case sendMessageMethod:
		return p.tg.SendMessage(chat.ID, msg, parseMode, replyToMessageID)
	case sendPhotoMethod:
		return p.tg.SendPhoto(chat.ID, msg)
	case sendMessageWithButtonsMethod:
		return p.tg.SendMessage(chat.ID, msg, parseMode, replyToMessageID)
	case sendPoll:
		return p.tg.SendPoll(response.poll)
	case doNothingMethod:
		log.Printf("Message: \"%s\" - do nothing", text)
	}
	return nil
}

// getCmd возвращает executor для команды, если она существует.
func (p *Processor) getCmd(strCmd string) CmdExecutor {
	if !strings.Contains(strCmd, "@") {
		strCmd += suffix
	}
	cmd, ok := p.commands[strCmd]
	if !ok {
		return nil
	}
	return cmd
}

// createNewUserInDB создаёт пользователя в базе данных, если он там ещё не существует.
func (p *Processor) createNewUserInDB(chatID int, user *telegram.User) (*storage.DBUser, error) {
	dbUserStatID, err := p.storage.CreateUserStats(context.Background(), &storage.DBUserStat{})

	if err != nil {
		return nil, e.Wrap("can't create user stats in 'createNewUserInDB'", err)
	}
	dbUser := &storage.DBUser{
		TgID:               user.ID,
		ChatID:             chatID,
		IsBot:              user.IsBot,
		IsPremium:          user.IsPremium,
		FirstName:          user.FirstName,
		LastName:           user.LastName,
		Username:           user.Username,
		UserStatId:         dbUserStatID,
		MaxDickChangeCount: maxDickChangeCount,
		HealthPoints:       defaultHpUser,
	}
	err = p.storage.CreateUser(context.Background(), dbUser)

	if err != nil {
		return nil, err
	}
	return dbUser, nil
}

// userChangeInfo изменяет данные пользователя в базе данных, если он поменял данные в телеграмме.
func (p *Processor) userChangeInfo(user *telegram.User, dbUser *storage.DBUser) (*storage.DBUser, error) {
	if user.FirstName != dbUser.FirstName || user.LastName != dbUser.LastName ||
		user.Username != dbUser.Username || user.IsPremium != dbUser.IsPremium {
		newDbUser := &storage.DBUser{
			ID:                 dbUser.ID,
			TgID:               dbUser.TgID,
			ChatID:             dbUser.ChatID,
			IsBot:              user.IsBot,
			IsPremium:          user.IsPremium,
			FirstName:          user.FirstName,
			LastName:           user.LastName,
			Username:           user.Username,
			DickSize:           dbUser.DickSize,
			ChangeDickAt:       dbUser.ChangeDickAt,
			UserStatId:         dbUser.UserStatId,
			HealthPoints:       dbUser.HealthPoints,
			HpTakedAt:          dbUser.HpTakedAt,
			IsGay:              dbUser.IsGay,
			GayAt:              dbUser.GayAt,
			Points:             dbUser.Points,
			CurDickChangeCount: dbUser.CurDickChangeCount,
			MaxDickChangeCount: dbUser.MaxDickChangeCount,
		}
		err := p.storage.UpdateUser(context.Background(), newDbUser)
		if err != nil {
			return newDbUser, err
		}
		p.userCache.AddUser(newDbUser)
	}
	return dbUser, nil
}

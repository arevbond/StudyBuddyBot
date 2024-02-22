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
)

var (
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

// allCommands список всех возможных команд бота.
var allCommands = map[string]CmdExecutor{
	AllCmd + suffix:                allUsernamesExec(AllCmd + suffix),
	GayTopCmd + suffix:             topGaysExec(GayTopCmd + suffix),
	GayStartCmd + suffix:           gayExec(GayStartCmd + suffix),
	DickTopCmd + suffix:            dickTopExec(DickTopCmd + suffix),
	DicStartCmd + suffix:           dickStartExec(DicStartCmd + suffix),
	GetHPCmd + suffix:              getHpExec(GetHPCmd + suffix),
	DickDuelCmd + suffix:           duelExec(DickDuelCmd + suffix),
	HelpCmd + suffix:               helpExec(HelpCmd + suffix),
	GetMyStatsCmd + suffix:         myStatsExec(GetMyStatsCmd + suffix),
	GetChatStatsCmd + suffix:       chatStatsExec(GetChatStatsCmd + suffix),
	ChangeDickCmd + suffix:         adminChangeDickExec(ChangeDickCmd + suffix),
	SendMessageByAdminCmd + suffix: adminSendMessageExec(SendMessageByAdminCmd + suffix),
	AddCalendarIDCmd + suffix:      addCalendarExec(AddCalendarIDCmd + suffix),
	ScheduleCmd + suffix:           scheduleExec(ScheduleCmd + suffix),
	XkcdCmd + suffix:               xkcdExec(XkcdCmd + suffix),
	AnecdotCmd + suffix:            anekdotExec(AnecdotCmd + suffix),
	AufCmd + suffix:                aufExec(AufCmd + suffix),
	FlipCmd + suffix:               flipExec(FlipCmd + suffix),
	GetChatIDCmd + suffix:          chatIDExec(GetChatIDCmd + suffix),

	AddHomeworkCmd + suffix:    addHomeworkExec(AddHomeworkCmd + suffix),
	GetHomeworkCmd + suffix:    getHomeworkExec(GetHomeworkCmd + suffix),
	DeleteHomeworkCmd + suffix: deleteHomeworkExec(DeleteHomeworkCmd + suffix),

	StartAuctionCmd + suffix: startAuctionExec(StartAuctionCmd + suffix),
	AddDepositCmd + suffix:   addDepositExec(AddDepositCmd + suffix),
	AuctionCmd + suffix:      auctionExec(AuctionCmd + suffix),

	StartQuizCmd + suffix: startQuizExec(StartQuizCmd + suffix),
}

const (
	MAX_DICK_CHANGE_COUNT = 1
	DEFAULT_HP_USER       = 3
)

// doCmd выбирает необходимую логику для выолнения команды.
func (p *Processor) doCmd(text string, chat *telegram.Chat, user *telegram.User, messageID int) error {
	dbUser, err := p.userCache.GetUser(user.ID, chat.ID)
	if err == storage.ErrUserNotExist {
		dbUser, err = p.storage.GetUser(context.Background(), user.ID, chat.ID)
		if err == storage.ErrUserNotExist {
			dbUser, err = p.createNewUserInDB(chat.ID, user)
			if err != nil {
				return e.Wrap("can't create new user in 'doCmd'", err)
			}
		} else if err != nil {
			return e.Wrap("unknown error in 'doCmd'", err)
		}
		p.userCache.AddUser(dbUser)
	}

	dbUser, err = p.userChangeInfo(user, dbUser)
	if err != nil {
		return e.Wrap("can't update user info in 'doCmd'", err)
	}

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

	text, parseMode := strings.TrimSpace(text), telegram.ParseMode("")

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
		log.Printf("[INFO] got new command '%s' from '%s' in '%s'", text, user.Username, chat.Title)

		strCmd := strings.Split(text, " ")[0]
		cmd := p.getCmd(strCmd)
		if cmd == nil {
			return e.Wrap(fmt.Sprintf("can't get command from %s", strCmd), err)
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
	}

	return nil
}

// getCmd возвращает executor для команды, если она существует.
func (p *Processor) getCmd(strCmd string) CmdExecutor {
	if !strings.Contains(strCmd, "@") {
		strCmd += suffix
	}
	cmd, ok := allCommands[strCmd]
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
		MaxDickChangeCount: MAX_DICK_CHANGE_COUNT,
		HealthPoints:       DEFAULT_HP_USER,
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

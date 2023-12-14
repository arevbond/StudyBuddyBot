package telegram

import (
	"tg_ics_useful_bot/clients/telegram"
	"tg_ics_useful_bot/lib/e"
	"tg_ics_useful_bot/storage"
)

type method int

const (
	UnsupportedMethod method = iota
	sendMessageMethod
	sendPhotoMethod
	sendMessageWithButtonsMethod
	doNothingMethod
)

const (
	suffix = "@ics_useful_bot"
)

const (
	HelpCmd = "/help"

	DicStartCmd = "/dick"
	DickTopCmd  = "/top_dick"
	DickDuelCmd = "/duel"

	GayStartCmd = "/gay"
	GayTopCmd   = "/top_gay"

	AddCalendarIDCmd = "/add_calendar"

	ScheduleCmd = "/schedule"
	AnecdotCmd  = "/joke"
	XkcdCmd     = "/xkcd"
	FlipCmd     = "/flip"

	AllCmd = "/all"

	AddHomeworkCmd    = "/add"
	GetHomeworkCmd    = "/get"
	DeleteHomeworkCmd = "/delete"
	CancelHomeworkCmd = "/cancel"

	GetMyStatsCmd   = "/my_stats"
	GetChatStatsCmd = "/chat_stats"

	GetChatIDCmd = "/chat_id"

	GetHPCmd = "/hp"

	// admins commands
	ChangeDickCmd         = "/change_dick"
	SendMessageByAdminCmd = "/send_message"
)

type gameGayCmd struct {
	command string
	p       *Processor
}

func (a *gameGayCmd) Exec(inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {
	message, err := a.p.gameGay(chat.ID)
	if err != nil {
		return nil, e.Wrap("can't get message from gameGay: ", err)
	}
	mthd := sendMessageMethod
	return &Response{message: message, method: mthd, replyMessageId: messageID}, nil
}

func (a *gameGayCmd) SetProcessor(p *Processor) {
	a.p = p
}

type topGaysCmd struct {
	command string
	p       *Processor
}

func (a *topGaysCmd) Exec(inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {
	message, err := a.p.topGays(chat.ID)
	if err != nil {
		return nil, e.Wrap("can't do GayTop: ", err)
	}
	mthd := sendMessageMethod
	return &Response{message: message, method: mthd, replyMessageId: messageID}, nil
}

func (a *topGaysCmd) SetProcessor(p *Processor) {
	a.p = p
}

package telegram

const (
	suffix = "@ics_useful_bot"
)

// dick
const (
	DicStartCmd    = "/dick"
	dickTopCommand = "/top_dick"
	DickDuelCmd    = "/duel"
	GetHPCmd       = "/hp"

	// gay
	GayStartCmd = "/gay"
	GayTopCmd   = "/top_gay"

	// calendar
	AddCalendarIDCmd = "/add_calendar"
	ScheduleCmd      = "/schedule"

	// homework
	AddHomeworkCmd    = "/add"
	GetHomeworkCmd    = "/get"
	DeleteHomeworkCmd = "/delete"
	CancelHomeworkCmd = "/cancel"

	// quit
	StartQuizCmd = "/quiz"
	StopQuizCmd  = "/stop"

	// auction
	StartAuctionCmd  = "/start_auction"
	FinishAuctionCmd = "/finish_auction"
	AddDepositCmd    = "/deposit"
	AuctionCmd       = "/auction"

	// stats
	GetMyStatsCmd   = "/my_stats"
	GetChatStatsCmd = "/chat_stats"

	// admins
	ChangeDickCmd         = "/change_dick"
	SendMessageByAdminCmd = "/send_message"

	// utils
	AllCmd       = "/all"
	HelpCmd      = "/help"
	AnecdotCmd   = "/joke"
	AufCmd       = "/auf"
	XkcdCmd      = "/xkcd"
	FlipCmd      = "/flip"
	GetChatIDCmd = "/chat_id"

	// other
	// HOLIDAY
	HolidayCmd = "/holiday"
)

// TODO: отрефакторить функцию; сделать админские команды более явными с помощью фабрик
func getAllCommands() map[string]CmdExecutor {
	return map[string]CmdExecutor{
		AllCmd + suffix:                allUsernamesExec(AllCmd + suffix),
		GayTopCmd + suffix:             topGaysExec(GayTopCmd + suffix),
		GayStartCmd + suffix:           gayExec(GayStartCmd + suffix),
		dickTopCommand + suffix:        dickTopExec(dickTopCommand + suffix),
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
		StopQuizCmd + suffix:  stopQuizExec(StopQuizCmd + suffix),
	}
}

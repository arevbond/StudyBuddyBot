package telegram

type method int

const (
	UnsupportedMethod method = iota
	sendMessageMethod
	sendPhotoMethod
	sendMessageWithButtonsMethod
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

	ChangeDickCmd = "/change_dick"
)

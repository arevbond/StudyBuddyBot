package telegram

type method int

const (
	sendMessageMethod method = iota
	sendPhotoMethod
)

var (
	AnecdotCmd         = "/joke"
	FlipCmd            = "/flip"
	GayStartCmd        = "/gay"
	GayTopCmd          = "/top_gay"
	XkcdCmd            = "/xkcd"
	DicStartCmd        = "/dick"
	DickTopCmd         = "/top_dick"
	DickDuelCmd        = "/duel"
	TodayLessonsCmd    = "/today"
	LessonsCmd         = "/lessons"
	TomorrowLessonsCmd = "/tomorrow"
)

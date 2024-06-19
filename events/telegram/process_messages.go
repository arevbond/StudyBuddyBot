package telegram

import (
	"errors"
	"log/slog"
	"tg_ics_useful_bot/clients/telegram"
	"tg_ics_useful_bot/events"
	"tg_ics_useful_bot/lib/e"
	"tg_ics_useful_bot/lib/quiz"
	"tg_ics_useful_bot/storage"
	"tg_ics_useful_bot/storage/cache"
)

type Processor struct {
	logger    *slog.Logger
	tg        *telegram.Client
	offset    int
	storage   storage.Storage
	userCache cache.UserCache
	quiz      *quizState
	commands  map[string]CmdExecutor
	auctions  map[int][]*AuctionPlayer
}

type quizState struct {
	currentQuestion *quiz.Question
	currentPlayers  map[int]int
	quit            chan bool
	allCommands     map[string]CmdExecutor
}

func newQuizState() *quizState {
	return &quizState{
		currentQuestion: &quiz.Question{},
		currentPlayers:  make(map[int]int),
		quit:            make(chan bool),
	}
}

type Meta struct {
	MessageID int

	TgID      int
	Username  string
	FirstName string
	LastName  string
	IsBot     bool
	IsPremium bool

	ChatID              int
	ChatType            string
	ChatTitle           string
	ChatActiveUsernames []string

	PollID    string
	OptionIds []int
}

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")

	ErrNotAdmin = errors.New("not admin")
)

func New(client *telegram.Client, storage storage.Storage, userCache cache.UserCache, logger *slog.Logger) *Processor {
	return &Processor{
		logger:    logger,
		tg:        client,
		storage:   storage,
		userCache: userCache,
		quiz:      newQuizState(),
		commands:  getAllCommands(),
		auctions:  make(map[int][]*AuctionPlayer),
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, e.Wrap("can't get events", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}

	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	case events.PollAnswer:
		return p.processPollAnswer(event)
	default:
		return e.Wrap("can't process message", ErrUnknownEventType)
	}
}

func (p *Processor) processPollAnswer(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.Wrap("can't process message", err)
	}
	userID := meta.TgID
	optionIds := meta.OptionIds
	if p.quiz.currentQuestion != nil {
		for _, id := range optionIds {
			if p.quiz.currentQuestion.CorrectOptionID == id {
				p.quiz.currentPlayers[userID]++
			}
		}
	}

	return nil
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.Wrap("can't process message", err)
	}

	messageID := meta.MessageID

	user := &telegram.User{
		ID:        meta.TgID,
		IsBot:     meta.IsBot,
		FirstName: meta.FirstName,
		LastName:  meta.LastName,
		Username:  meta.Username,
		IsPremium: meta.IsPremium,
	}

	chat := &telegram.Chat{
		ID:              meta.ChatID,
		Type:            meta.ChatType,
		Title:           meta.ChatTitle,
		ActiveUsernames: meta.ChatActiveUsernames,
	}

	if chat.Type == "private" && !p.isAdmin(user.ID) {
		p.logger.Info("user send private message", slog.Int("tg id", user.ID))
		return e.Wrap("can't process private not admin message", ErrNotAdmin)
	}

	if err = p.doCmd(event.Text, chat, user, messageID); err != nil {
		return e.Wrap("can't process message", err)
	}

	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("can't get meta", ErrUnknownMetaType)
	}

	return res, nil
}

func event(upd telegram.Update) events.Event {
	updType := fetchType(upd)

	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}

	if updType == events.Message {
		res.Meta = Meta{
			MessageID: upd.Message.ID,

			TgID:      upd.Message.From.ID,
			FirstName: upd.Message.From.FirstName,
			LastName:  upd.Message.From.LastName,
			Username:  upd.Message.From.Username,
			IsBot:     upd.Message.From.IsBot,
			IsPremium: upd.Message.From.IsPremium,

			ChatID:              upd.Message.Chat.ID,
			ChatType:            upd.Message.Chat.Type,
			ChatTitle:           upd.Message.Chat.Title,
			ChatActiveUsernames: upd.Message.Chat.ActiveUsernames,
		}
	} else if updType == events.PollAnswer {
		res.Meta = Meta{
			PollID:    upd.PollAnswer.PollID,
			TgID:      upd.PollAnswer.User.ID,
			OptionIds: upd.PollAnswer.OptionIds,
		}
	}
	return res
}

func fetchText(upd telegram.Update) string {
	if upd.Message == nil {
		return ""
	}

	return upd.Message.Text
}

func fetchType(upd telegram.Update) events.Type {
	if upd.Message != nil {
		return events.Message
	} else if upd.PollAnswer != nil {
		return events.PollAnswer
	}
	return events.Unknown
}

// isAdmin проверяет является ли телеграм пользователь админом бота.
func (p *Processor) isAdmin(userID int) bool {
	for _, id := range p.tg.AdminsID {
		if userID == id {
			return true
		}
	}
	return false
}

// isChatAdmin определяет является ли пользователь админов в чате.
func (p *Processor) isChatAdmin(user *telegram.User, chatID int) bool {
	admins, err := p.tg.ChatAdministrators(chatID)
	if err != nil {
		p.logger.Error("can't get admins", slog.Any("error", err), slog.Int("chat id", chatID))
	}
	for _, admin := range admins {
		if user.ID == admin.ID {
			return true
		}
	}
	return false
}

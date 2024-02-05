package telegram

import (
	"errors"
	"tg_ics_useful_bot/clients/telegram"
	"tg_ics_useful_bot/events"
	"tg_ics_useful_bot/lib/e"
	"tg_ics_useful_bot/storage"
	"tg_ics_useful_bot/storage/cache"
)

type Processor struct {
	tg        *telegram.Client
	offset    int
	storage   storage.Storage
	userCache cache.UserCache
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
}

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

func New(client *telegram.Client, storage storage.Storage, userCache cache.UserCache) *Processor {
	return &Processor{
		tg:        client,
		storage:   storage,
		userCache: userCache,
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
	default:
		return e.Wrap("can't process message", ErrUnknownEventType)
	}
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
		return nil
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

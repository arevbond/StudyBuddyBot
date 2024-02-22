package telegram

import "tg_ics_useful_bot/lib/quiz"

type ChatMemberAdministratorResponse struct {
	Ok     bool                      `json:"ok"`
	Result []ChatMemberAdministrator `json:"result"`
}

type ChatMemberAdministrator struct {
	User User `json:"user"`
}

type UpdatesResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	ID            int              `json:"update_id"`
	Message       *IncomingMessage `json:"message"`
	CallbackQuery *CallbackQuery   `json:"callback_query"`
	PollAnswer    *PollAnswer      `json:"poll_answer"`
}

type IncomingMessage struct {
	ID             int     `json:"message_id"`
	Text           string  `json:"text"`
	From           User    `json:"from"`
	Date           int     `json:"date"` // Date the message was sent in Unix time
	Chat           Chat    `json:"chat"`
	ReplyToMessage Message `json:"reply_to_message"`
}

type CallbackQuery struct {
	ID      string          `json:"id"`
	From    User            `json:"from"`
	Message IncomingMessage `json:"message"`
	Data    string          `json:"data"`
}

type User struct {
	ID        int    `json:"id"`
	IsBot     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	IsPremium bool   `json:"is_premium"`
}

type Chat struct {
	ID              int      `json:"id"`
	Type            string   `json:"type"`
	Title           string   `json:"title"`
	ActiveUsernames []string `json:"active_usernames"`
}

type Message struct {
	ChatID           int    `json:"chat_id"`
	Text             string `json:"text"`
	ParseMode        string `json:"parse_mode"`
	ReplyToMessageID int    `json:"reply_to_message_id"`
}

type ForceReply struct {
	ForceReply       bool   `json:"force_reply"`
	InputPlaceHolder string `json:"input_place_holder"`
	Selective        bool   `json:"selective"`
}

type InlineKeyboardMarkup struct {
	Keyboard        [][]InlineKeyboardButton `json:"inline_keyboard"`
	OneTimeKeyboard bool                     `json:"one_time_keyboard"`
}

type InlineKeyboardButton struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data"`
}

type Poll struct {
	ID                    string   `json:"id"`
	Type                  string   `json:"type"`
	IsAnonymous           bool     `json:"is_anonymous"`
	Question              string   `json:"question"`
	Options               []string `json:"options"`
	AllowsMultipleAnswers bool     `json:"allows_multiple_answers"`
	CorrectOptionID       int      `json:"correct_option_id"`
	Explanation           string   `json:"explanation"`
	OpenPeriod            int      `json:"open_period"`
}

type PollAnswer struct {
	PollID    string `json:"poll_id"`
	VoterChat Chat   `json:"voter_chat"`
	User      User   `json:"user"`
	OptionIds []int  `json:"option_ids"`
}

type SendPoll struct {
	ChatID                int      `json:"chat_id"`
	Type                  string   `json:"type"`
	IsAnonymous           bool     `json:"is_anonymous"`
	Question              string   `json:"question"`
	Options               []string `json:"options"`
	AllowsMultipleAnswers bool     `json:"allows_multiple_answers"`
	CorrectOptionID       int      `json:"correct_option_id"`
	Explanation           string   `json:"explanation"`
	OpenPeriod            int      `json:"open_period"`
}

func NewSendPoll(chatID int, question *quiz.Question) SendPoll {
	return SendPoll{
		ChatID:                chatID,
		Type:                  "quiz",
		IsAnonymous:           false,
		Question:              question.Question,
		Options:               question.Options,
		AllowsMultipleAnswers: question.AllowsMultipleAnswers,
		CorrectOptionID:       question.CorrectOptionID,
		Explanation:           question.Explanation,
		OpenPeriod:            question.OpenPeriod,
	}
}

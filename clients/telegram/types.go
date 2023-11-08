package telegram

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
}

type IncomingMessage struct {
	ID   int    `json:"message_id"`
	Text string `json:"text"`
	From User   `json:"from"`
	Date int    `json:"date"` // Date the message was sent in Unix time
	Chat Chat   `json:"chat"`
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

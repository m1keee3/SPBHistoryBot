package telegram

type UpdatesService interface {
	Updates(offset int, limit int) ([]Update, error)
}

type MessageService interface {
	SendMessage(chatID int, text string, keyboard InlineKeyboardMarkup) error
	SendPhoto(chatID int, text string, photoURL string, keyboard InlineKeyboardMarkup) error
	EditMessage(chatID int, messageID int, text string, markup InlineKeyboardMarkup) error
	EditPhoto(chatID int, messageID int, text string, photoURL string, markup InlineKeyboardMarkup) error
	SendNoButtonsMessage(chatID int, text string) error
	DeleteMessage(chatID int, messageID int) error
}

type UpdatesResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	ID            int              `json:"update_id"`
	Message       *IncomingMessage `json:"message,omitempty"`
	CallbackQuery *CallbackQuery   `json:"callback_query,omitempty"`
}

type IncomingMessage struct {
	MessageID int    `json:"message_id"`
	Text      string `json:"text"`
	From      User   `json:"from"`
	Chat      Chat   `json:"chat"`
}

type CallbackQuery struct {
	ID      string           `json:"id"`
	From    User             `json:"from"`
	Message *IncomingMessage `json:"message"`
	Data    string           `json:"data"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type Chat struct {
	ID int `json:"id"`
}
type InlineKeyboardButton struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data"`
}

type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

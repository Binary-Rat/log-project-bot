package tg

type UpdatesResponse struct {
	OK     bool      `json:"ok"`
	Result []Updates `json:"result"`
}

type Updates struct {
	ID            int                    `json:"update_id"`
	Message       *IncomingMessage       `json:"message"`
	CallbackQuery *IncomingCallbackQuery `json:"callback_query"`
}

// Refactor message structures
type IncomingMessage struct {
	Text string `json:"text"`
	From From   `json:"from"`
	Chat Chat   `json:"chat"`
}

type OutcomingMessage struct {
	Text        string       `json:"text"`
	ChatID      int          `json:"chat_id"`
	ReplyMarkup *ReplyMarkup `json:"reply_markup"`
}

type From struct {
	Username string `json:"username"`
}

type Chat struct {
	ID int `json:"id"`
}

type ReplyMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

type IncomingCallbackQuery struct {
	ID      string          `json:"id"`
	From    From            `json:"from"`
	Message IncomingMessage `json:"message"`
	Data    string          `json:"data"`
}

type InlineKeyboardButton struct {
	Text         string `json:"text"`
	URL          string `json:"url"`
	CallbackData string `json:"callback_data"`
}

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
	Text        string   `json:"text"`
	ChatID      int      `json:"chat_id"`
	ReplyMarkup keyboard `json:"reply_markup,omitempty"`
}

type From struct {
	Username string `json:"username"`
}

type Chat struct {
	ID int `json:"id"`
}

type ReplyMarkup struct {
	Keyboard       [][]KeyboardButton `json:"keyboard,omitempty"`
	InlineKeyboard [][]KeyboardButton `json:"inline_keyboard,omitempty"`
	OneTime        bool               `json:"one_time_keyboard,omitempty"`
}

type IncomingCallbackQuery struct {
	ID      string          `json:"id"`
	From    From            `json:"from"`
	Message IncomingMessage `json:"message"`
	Data    string          `json:"data"`
}

type KeyboardButton struct {
	Text         string `json:"text"`
	URL          string `json:"url,omitempty"`
	CallbackData string `json:"callback_data,omitempty"`
}

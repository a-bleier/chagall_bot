package comm

/*
This file contains types which may be send to the api.
*/
type SendMessage struct {
	ChatID                string      `json:"chat_id,omitempty"`
	Text                  string      `json:"text,omitempty"`
	ParseMode             string      `json:"parse_mode,omitempty"`
	DisableWebPagePreview bool        `json:"disable_web_page_preview,omitempty"`
	ReplyMarkup           interface{} `json:"reply_markup,omitempty"`
}
type InlineKeyboardMarkup struct {
	Keyboard [][]KeyboardButton `json:"inline_keyboard"`
}

type ReplyKeyboardMarkup struct {
	Keyboard [][]KeyboardButton `json:"keyboard"`
}

func NewInlinekeyboardMarkup(names [][]string) InlineKeyboardMarkup {
	buttonField := make([][]KeyboardButton, 0)
	for _, row := range names {
		buttonRow := make([]KeyboardButton, 0)
		for _, name := range row {
			buttonRow = append(buttonRow, KeyboardButton{name, name})
		}
		buttonField = append(buttonField, buttonRow)
	}

	return InlineKeyboardMarkup{buttonField}
}

func NewReplykeyboardMarkup(names [][]string) ReplyKeyboardMarkup {
	buttonField := make([][]KeyboardButton, 0)
	for _, row := range names {
		buttonRow := make([]KeyboardButton, 0)
		for _, name := range row {
			buttonRow = append(buttonRow, KeyboardButton{Text: name,
				CallbackData: ""})
		}
		buttonField = append(buttonField, buttonRow)
	}

	return ReplyKeyboardMarkup{buttonField}
}

type KeyboardButton struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data,omitempty"`
	//Some optional stuff can be inserted here later
}

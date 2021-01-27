package comm

import "encoding/json"

//Update contains unmarshalled JSON Data of an Update From Telegram
type APIResponse struct {
	Ok          bool               `json:"ok"`
	Result      json.RawMessage    `json:"result"`
	ErrorCode   int                `json:"error_code"`
	Description string             `json:"description"`
	Paramters   ResponseParameters `json:"parameters"`
}

type ResponseParameters struct {
	MigrateToChatID int64 `json:"migrate_to_chat_id"`
	RetryAfter      int   `json:"retry_after"`
}
type Update struct {
	Id            int           `json:"update_id"`
	Message       Message       `json:"message"`
	InlineQuery   InlineQuery   `json:"inline_query"`
	CallbackQuery CallbackQuery `json:"callback_query"`
}

func (u *Update) GetUserId() uint64 {

	if u.Message.From.Id != 0 {
		return u.Message.From.Id
	} else {
		return u.CallbackQuery.From.Id
	}
	return 0
}

//Message is Message Type
type Message struct {
	Id   int    `json:"message_id"`
	From User   `json:"from"`
	Chat Chat   `json:"chat"`
	Date int    `json:"date"`
	Text string `json:"text"`
}

//User is User type
type User struct {
	Id           uint64 `json:"id"`
	IsBot        bool   `json:"is_bot"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	LanguageCode string `json:"language_code"`
}

//Chat is Chat type
type Chat struct {
	Id        uint64 `json:"id"`
	UserName  string `json:"username"`
	FirstName string `json:"first_name"`
	Type      string `json:"type"`
}

type InlineQuery struct {
	Id     string `json:"id"`
	From   User   `json:"from"`
	Query  string `json:"query"`
	Offset string `json:"offset"`
}

type CallbackQuery struct {
	Id              string  `json:"id"`
	From            User    `json:"from"`
	Message         Message `json:"message,omitempty"`
	InlineMessageId string  `json:"inline_message_id.omitempty"`
	Data            string  `json:"data"`
}

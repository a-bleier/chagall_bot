package comm

//Update contains unmarshalled JSON Data of an Update From Telegram
type Update struct {
	Id      int
	Message Message
}

//Message is Message Type
type Message struct {
	ID    int
	From  User
	Mchat Chat
	Date  int
	Text  string
}

func NewUpdateFromJSON(rawUpdate interface{}) Update {
	updateMap := rawUpdate.(map[string]interface{})
	updateId := int(updateMap["update_id"].(float64))
	messageMap := updateMap["message"].(map[string]interface{})
	messageId := int(messageMap["message_id"].(float64))
	fromMap := messageMap["from"].(map[string]interface{})
	chatMap := messageMap["chat"].(map[string]interface{})

	chat := Chat{
		ID:        int(chatMap["id"].(float64)),
		firstName: chatMap["first_name"].(string),
		cType:     chatMap["type"].(string),
	}
	from := User{
		id:           int(fromMap["id"].(float64)),
		isBot:        fromMap["is_bot"].(bool),
		firstName:    fromMap["first_name"].(string),
		languageCode: fromMap["language_code"].(string),
	}
	msg := Message{
		ID:    messageId,
		From:  from,
		Mchat: chat,
		Date:  int(messageMap["date"].(float64)),
		Text:  messageMap["text"].(string),
	}
	return Update{
		Id:      updateId,
		Message: msg,
	}
}

//User is User type
type User struct {
	id           int
	isBot        bool
	firstName    string
	languageCode string
}

//Chat is Chat type
type Chat struct {
	ID        int
	firstName string
	cType     string
}

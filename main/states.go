package main

import (
	"encoding/json"
	"fmt"
	"github.com/a-bleier/chagall_bot/comm"
)

type state int

var userStateLookup map[uint64]state

const (
	START_STATE state = iota
	CHOOSING_SERVICE_STATE
	BIRTHDAYS_STATE
	ADD_BIRTHDAY_STATE
	REMOVE_BIRTHDAYY_STATE
)

func initStates() {
	userStateLookup = make(map[uint64]state)
}

func transitStates(update comm.Update) bool {

	var userId uint64
	var chatId uint64
	if update.Message.Id != 0 { //This means, the update contains a message
		userId = update.Message.From.Id
		chatId = update.Message.Chat.Id
	} else if update.InlineQuery.Id != "" { // This means, the update contains an callback from an inline
		userId = update.InlineQuery.From.Id
	} else if update.CallbackQuery.Id != "" { //This means, the update contains a callback query
		userId = update.CallbackQuery.From.Id
	}
	currentState := userStateLookup[userId]
	switch currentState {
	case START_STATE:
		sendServiceOffer(chatId)
		userStateLookup[userId] = CHOOSING_SERVICE_STATE
		break
	case CHOOSING_SERVICE_STATE:
		break
	}
	return false
}

func processServiceCallback(cbQuery comm.CallbackQuery) state {

	if cbQuery.Id == "" {
		return CHOOSING_SERVICE_STATE
	} else {

	}
	return START_STATE
}

func sendServiceOffer(chatId uint64) {
	row := make([]string, 2)
	var field [][]string
	row[0] = "Birthdays"
	row[1] = "Quit"
	field = append(field, row)

	inlineKeyboard := comm.NewInlinekeyboardMarkup(field)
	text := "Hi !\nI'm Chagall, a friendly bot.\nHow can I help you ?"
	sMessage := comm.SendMessage{Text: text,
		ReplyMarkup: inlineKeyboard,
		ChatID:      fmt.Sprintf("%d", chatId),
	}
	data, err := json.Marshal(sMessage)
	if err != nil {
		panic(err)
	}
	item := comm.QueueItem{data, "sendMessage"}
	txQueue.EnQueue(item)

}

func sendMessageInlineKeyboard(message string, keyboardTexts, callbackMessages [][]string) {

}

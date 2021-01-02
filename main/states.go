package main

//TODO maybe move in own package
import (
	"encoding/json"
	"fmt"
	"github.com/a-bleier/chagall_bot/comm"
)

type state int

const (
	START_STATE state = iota
	CHOOSING_SERVICE_STATE
	BIRTHDAYS_STATE
	ADD_BIRTHDAY_STATE
	REMOVE_BIRTHDAYY_STATE
)

type StateMachine struct {
	userStateLookup map[uint64]state
	textFacility    TextFacility
}

func NewStateMachine() StateMachine {
	return StateMachine{userStateLookup: make(map[uint64]state),
		textFacility: NewTextFacility()}
}

func (s *StateMachine) transitStates(update comm.Update) bool {

	var userId uint64
	var chatId uint64
	if update.Message.Id != 0 { //This means, the update contains a message
		userId = update.Message.From.Id
		chatId = update.Message.Chat.Id
	} else if update.InlineQuery.Id != "" { // This means, the update contains an callback from an inline
		userId = update.InlineQuery.From.Id
	} else if update.CallbackQuery.Id != "" { //This means, the update contains a callback query
		userId = update.CallbackQuery.From.Id
		answerCallbackQuery(update.CallbackQuery.Id)
	}
	currentState := s.userStateLookup[userId]
	switch currentState {
	case START_STATE:
		sendTextInlineKeyboard(
			fmt.Sprintf("%d", chatId),
			"introduction",
			"serviceOffer",
			&s.textFacility)
		s.userStateLookup[userId] = CHOOSING_SERVICE_STATE
		break
	case CHOOSING_SERVICE_STATE:
		s.userStateLookup[userId] = processServiceCallback(update.CallbackQuery, &s.textFacility)
		break
	case BIRTHDAYS_STATE: //NOTE: When adding new services, thi from here encapsule in new function
		s.userStateLookup[userId] = processBirthdaysCallback(update.CallbackQuery, &s.textFacility)
		break
	}
	return false
}

func processBirthdaysCallback(cbQuery comm.CallbackQuery, facility *TextFacility) state {
	var retState state

	if cbQuery.Id == "" {
		retState = BIRTHDAYS_STATE

		if cbQuery.Data == "Back" {
			sendTextInlineKeyboard(fmt.Sprintf("%d", cbQuery.Message.Chat.Id),
				"offerServiceAgain",
				"serviceOffer",
				facility)
			retState = CHOOSING_SERVICE_STATE
		} else if cbQuery.Data == "List" {

		} else if cbQuery.Data == "Add" {

		} else if cbQuery.Data == "Remove" {

		} else if cbQuery.Data == "Edit" {

		}
	}
	return retState
}

//TODO better error processing
func processServiceCallback(cbQuery comm.CallbackQuery, facility *TextFacility) state {

	retState := START_STATE

	if cbQuery.Id == "" { //No callback, oopsie
		retState = CHOOSING_SERVICE_STATE
	} else {
		if cbQuery.Data == "Birthdays" {

			sendTextInlineKeyboard(fmt.Sprintf("%d", cbQuery.Message.Chat.Id),
				"birthdayService",
				"birthdayService",
				facility)
			retState = BIRTHDAYS_STATE
		} else {
			sendTextInlineKeyboard(fmt.Sprintf("%d", cbQuery.Message.Chat.Id),
				"goodbye",
				"",
				facility)
			retState = START_STATE
		}

	}
	return retState
}

func answerCallbackQuery(id string) {
	answer := comm.AnswerCallbackQuery{
		CallbackQueryId: id,
	}
	data, err := json.Marshal(answer)
	if err != nil {
		panic(err)
	}
	item := comm.QueueItem{data, "answerCallbackQuery"}
	txQueue.EnQueue(item)
}

//This function shall abstract the messages to the user which have the form  [text + inline buttons]
func sendTextInlineKeyboard(chatId string, messageKey string, inlineButtonGroupKey string, facility *TextFacility) {
	//give the key to textFacility, receive a inlinekeyboardTemplate [][]string
	//build a inlineKeyboard

	var sMessage comm.SendMessage
	if len(inlineButtonGroupKey) != 0 { //Inline keyboard needed when the key is not ""
		field := facility.getKeyboardTemplate(inlineButtonGroupKey)
		inlineKeyboard := comm.NewInlinekeyboardMarkup(field)
		sMessage = comm.SendMessage{Text: facility.getMessageText(messageKey),
			ReplyMarkup: inlineKeyboard,
			ChatID:      chatId,
		}
	} else { // No inline keyboard needed
		sMessage = comm.SendMessage{
			Text:   facility.getMessageText(messageKey),
			ChatID: chatId,
		}
	}

	data, err := json.Marshal(sMessage)
	if err != nil {
		panic(err)
	}
	item := comm.QueueItem{data, "sendMessage"}
	txQueue.EnQueue(item)

}

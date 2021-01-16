package main

//TODO maybe move in own package
import (
	"encoding/json"
	"fmt"
	"github.com/a-bleier/chagall_bot/comm"
	"github.com/a-bleier/chagall_bot/db"
	"strings"
)

type state int

const (
	START_STATE state = iota
	CHOOSING_SERVICE_STATE
	BIRTHDAYS_STATE
	ASK_BIRTHDAY_NAME
	ASK_BIRTHDAY_DATE
	ASK_BIRTHDAY_CONTACT
	ASK_BIRTHDAY_CONFIRMATION
	ADD_BIRTHDAY_CONFIRMATION
	REMOVE_BIRTHDAYY_STATE
)

type StateMachine struct {
	userStateLookup map[uint64]state
	textFacility    TextFacility
	bdStateMachine  birthdayStateMachine
}

func NewStateMachine() StateMachine {
	textFacility := NewTextFacility()
	bdStateMachine := birthdayStateMachine{&textFacility,
		make(map[uint64]state),
		birthdayEntry{}}
	return StateMachine{userStateLookup: make(map[uint64]state),
		textFacility:   textFacility,
		bdStateMachine: bdStateMachine}
}

//TODO: Encapsule each service into its own type with own states i.o. to shrink the main state machine
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
		sendTextInlineKeyboard("",
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
		s.userStateLookup[userId] = s.bdStateMachine.transitStates(update)
		break
	}
	return false
}

//TODO better error processing
func processServiceCallback(cbQuery comm.CallbackQuery, facility *TextFacility) state {

	retState := START_STATE

	if cbQuery.Id == "" { //No callback, oopsie
		retState = CHOOSING_SERVICE_STATE
	} else {
		if cbQuery.Data == "Birthdays" {

			sendTextInlineKeyboard("",
				fmt.Sprintf("%d", cbQuery.Message.Chat.Id),
				"birthdayService",
				"birthdayService",
				facility)
			retState = BIRTHDAYS_STATE
		} else {
			sendTextInlineKeyboard("",
				fmt.Sprintf("%d", cbQuery.Message.Chat.Id),
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
//TODO: Remove user id and database access from this function

func sendSimpleMessage(chatId, messageText string) {

	var sMessage comm.SendMessage
	sMessage = comm.SendMessage{
		Text:   messageText,
		ChatID: chatId,
	}
	data, err := json.Marshal(sMessage)
	if err != nil {
		panic(err)
	}
	item := comm.QueueItem{data, "sendMessage"}
	txQueue.EnQueue(item)
}
func sendTextInlineKeyboard(userId string, chatId string, messageKey string, inlineButtonGroupKey string, facility *TextFacility) {
	//give the key to textFacility, receive a inlinekeyboardTemplate [][]string
	//build a inlineKeyboard

	var sMessage comm.SendMessage
	var messageText string
	if messageKey == "listBirthdays" {
		messageText = strings.Join(db.ListAllBirthdays(userId), "\n") //TODO throw out db query from here
	} else if messageKey == "" {
		messageText = ""
	} else {
		messageText = facility.getMessageText(messageKey)
	}
	if len(inlineButtonGroupKey) != 0 { //Inline keyboard needed when the key is not ""
		field := facility.getKeyboardTemplate(inlineButtonGroupKey)
		inlineKeyboard := comm.NewInlinekeyboardMarkup(field)
		sMessage = comm.SendMessage{
			Text:        messageText,
			ReplyMarkup: inlineKeyboard,
			ChatID:      chatId,
		}
	} else { // No inline keyboard needed
		sMessage = comm.SendMessage{
			Text:   messageText,
			ChatID: chatId,
		}
	}
	//
	//data, err := json.Marshal(sMessage)
	//if err != nil {
	//	panic(err)
	//}
	//item := comm.QueueItem{data, "sendMessage"}
	//txQueue.EnQueue(item)
	sender.AddMessageToTx(sMessage, "sendMessage")

}

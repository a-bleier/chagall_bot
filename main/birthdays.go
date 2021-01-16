package main

import (
	"fmt"
	"github.com/a-bleier/chagall_bot/comm"
	"github.com/a-bleier/chagall_bot/db"
	"strings"
)

type birthdayStateMachine struct {
	textFacility    *TextFacility
	userStateLookup map[uint64]state
	newBirthday     birthdayEntry //FIXME: Only works when one user is adding a new entry; ---> add map[userId]entry
}

type birthdayEntry struct {
	name    string
	date    string
	contact string
}

func (b *birthdayStateMachine) transitStates(update comm.Update) state {
	var userId uint64
	var chatId uint64
	if update.Message.Id != 0 { //This means, the update contains a message
		userId = update.Message.From.Id
		chatId = update.Message.Chat.Id
	} else if update.InlineQuery.Id != "" { // This means, the update contains an callback from an inline
		userId = update.InlineQuery.From.Id
	} else if update.CallbackQuery.Id != "" { //This means, the update contains a callback query
		userId = update.CallbackQuery.From.Id
		chatId = update.CallbackQuery.Message.Chat.Id
	}

	switch b.userStateLookup[userId] {
	case 0, BIRTHDAYS_STATE:
		b.userStateLookup[userId] = b.processBirthdaysCallback(update, b.textFacility)
		if b.userStateLookup[userId] == CHOOSING_SERVICE_STATE {
			b.userStateLookup[userId] = BIRTHDAYS_STATE
			return CHOOSING_SERVICE_STATE
		}
		break
	case ASK_BIRTHDAY_NAME:
		fallthrough
	case ASK_BIRTHDAY_DATE:
		fallthrough
	case ASK_BIRTHDAY_CONTACT:
		fallthrough
	case ASK_BIRTHDAY_CONFIRMATION:
		b.userStateLookup[userId] = b.addRoutine(b.userStateLookup[userId], chatId, update.Message.Text)
	case ADD_BIRTHDAY_CONFIRMATION:
		confValue := update.CallbackQuery.Data
		if confValue == "Yes" {
			sendSimpleMessage(fmt.Sprintf("%d", chatId), "new entry confirmed")
			err := db.AddBirthday(fmt.Sprintf("%d", userId), b.newBirthday.date, b.newBirthday.name, b.newBirthday.contact)
			if err != nil {
				panic(err)
			}
			b.newBirthday = birthdayEntry{}
		} else if confValue == "No" {
			b.userStateLookup[userId] = BIRTHDAYS_STATE
			sendSimpleMessage(fmt.Sprintf("%d", chatId), "new entry discarded")
		}
		b.userStateLookup[userId] = BIRTHDAYS_STATE
		sendTextInlineKeyboard("",
			fmt.Sprintf("%d", chatId),
			"birthdayService",
			"birthdayService",
			b.textFacility)
		break
	}

	return BIRTHDAYS_STATE
}

func (b *birthdayStateMachine) processBirthdaysCallback(update comm.Update, facility *TextFacility) state {
	var retState state

	cbQuery := update.CallbackQuery

	if cbQuery.Id == "" {
		retState = BIRTHDAYS_STATE
	} else {

		if cbQuery.Data == "Back" {
			sendTextInlineKeyboard("",
				fmt.Sprintf("%d", cbQuery.Message.Chat.Id),
				"offerServiceAgain",
				"serviceOffer",
				facility)
			retState = CHOOSING_SERVICE_STATE
		} else if cbQuery.Data == "List" {
			sendSimpleMessage(fmt.Sprintf("%d", cbQuery.Message.Chat.Id),
				strings.Join(db.ListAllBirthdays(fmt.Sprintf("%d", cbQuery.From.Id)), "\n"))
			sendTextInlineKeyboard(fmt.Sprintf("%d", cbQuery.From.Id),
				fmt.Sprintf("%d", cbQuery.Message.Chat.Id),
				"birthdayService",
				"birthdayService",
				facility)
			retState = BIRTHDAYS_STATE
		} else if cbQuery.Data == "Add" {
			retState = b.addRoutine(ASK_BIRTHDAY_NAME, cbQuery.Message.Chat.Id, update.Message.Text)
		} else if cbQuery.Data == "Remove" {
			//tbd
		} else if cbQuery.Data == "Edit" {
			//tbd
		}
	}
	return retState
}

//TODO: Write a  function which puts the birthday reminders from db in the cron jobs

func (b *birthdayStateMachine) addRoutine(currentState state, chatId uint64, messageText string) state {

	var retState state = ASK_BIRTHDAY_NAME
	switch currentState {
	case ASK_BIRTHDAY_NAME:
		sendSimpleMessage(fmt.Sprintf("%d", chatId),
			b.textFacility.getMessageText("askBirthdayName"))
		retState = ASK_BIRTHDAY_DATE
	case ASK_BIRTHDAY_DATE:
		b.newBirthday.name = messageText
		sendSimpleMessage(fmt.Sprintf("%d", chatId),
			b.textFacility.getMessageText("askBirthdayDate"))
		retState = ASK_BIRTHDAY_CONTACT
	case ASK_BIRTHDAY_CONTACT:
		//Check if the received date is right
		//When wrong -> ASK_BIRTHDAY_DATE
		//When true -> ASK for contact
		b.newBirthday.date = messageText
		sendSimpleMessage(fmt.Sprintf("%d", chatId),
			b.textFacility.getMessageText("askBirthdayContact"))
		retState = ASK_BIRTHDAY_CONFIRMATION
	case ASK_BIRTHDAY_CONFIRMATION:
		//receive contact
		//ask for confirmation
		b.newBirthday.contact = messageText

		sendSimpleMessage(fmt.Sprintf("%d", chatId), fmt.Sprintf("%s %s %s %s",
			b.textFacility.getMessageText("askBirthdayConfirmation"),
			b.newBirthday.name,
			b.newBirthday.date,
			b.newBirthday.contact,
		),
		)
		sendTextInlineKeyboard("",
			fmt.Sprintf("%d", chatId),
			"confirmMessage",
			"yesNo",
			b.textFacility)
		retState = ADD_BIRTHDAY_CONFIRMATION //when birthday confirmed
	}
	return retState
}

func removeRoutine() {

	//give the users choices

	//ask the user which one to delete

	//delete from db

	//delete from cron
}

//Not so important here
func editRoutine() {

}

package main

import (
	"fmt"
	"github.com/a-bleier/chagall_bot/comm"
	"github.com/a-bleier/chagall_bot/db"
	"github.com/robfig/cron"
)

//TODO: Write a cron job which will check everyday at 00:00 if someone needs to be remembered because of an birthday

func initCronJobs() {
	cronObject := cron.New()
	fmt.Println("Setting up a cron job")
	cronObject.AddFunc("*/10 */1 * * * *", func() {
		fmt.Println("Job done")
	})
	cronObject.AddFunc("*/30 */1 * * * *", checkBirthdays)
	//For production
	//cronObject.AddFunc("* * * */31 * *", checkBirthdays)
	cronObject.Start()
}

func checkBirthdays() {
	/*
		TODO: query the birthdays which are today
		write a message to the users
	*/
	fmt.Println("checking for birthdays")
	ebrList, err := db.GetAllEntryBirthdayReminders()
	if err != nil {
		fmt.Println("error in the birthday cron job")
	}

	for _, entry := range ebrList {
		responseString := fmt.Sprintf("Hi %s\n", entry.UserName)
		responseString += fmt.Sprintf("%s has his / her birthday today.\n He / She was born on %s. You can reach this person in this way:\n %s", entry.Name, entry.Date, entry.Contact)
		sMessage := comm.SendMessage{
			Text:   responseString,
			ChatID: entry.ChatId,
		}
		sender.AddMessageToTx(sMessage, "sendMessage")
	}
}

package main

import (
	"fmt"
	"github.com/a-bleier/chagall_bot/comm"
	"github.com/a-bleier/chagall_bot/db"
	"github.com/a-bleier/chagall_bot/logging"
	"github.com/robfig/cron"
	"regexp"
	"time"
)

func initCronJobs() {
	cronObject := cron.New()
	logging.LogInfo("Setting up cron jobs")
	//TODO: Check if this works
	cronObject.AddFunc("0 23 * * *", checkBirthdays)
	cronObject.Start()
}

//checkBirthdays goes through entries in Birthdays table and check is a person has its birthday today. if it found one,
//a message will be put in the queue
func checkBirthdays() {
	logging.LogInfo("Checking for birthdays")
	ebrList, err := db.GetAllEntryBirthdayReminders()
	if err != nil {
		logging.LogWarning("An error occured while querying the database for birthdays")
	}

	for _, entry := range ebrList {

		t := time.Now()
		_, month, day := t.Date()
		pattern := fmt.Sprintf("(-|\\s)%02d(-|\\s)%02d", int(month), day)
		r, _ := regexp.Compile(pattern)
		if r.MatchString(entry.Date) == false {
			continue
		}
		responseString := fmt.Sprintf("Hi %s\n", entry.UserName)
		responseString += fmt.Sprintf("%s has his / her birthday today.\n He / She was born on %s. You can reach this person in this way:\n %s", entry.Name, entry.Date, entry.Contact)
		sMessage := comm.SendMessage{
			Text:   responseString,
			ChatID: entry.ChatId,
		}
		sender.AddMessageToTx(sMessage, "sendMessage")
	}
}

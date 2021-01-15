package main

import (
	"fmt"
	"github.com/robfig/cron"
)

//TODO: Write a cron job which will check everyday at 00:00 if someone needs to be remembered because of an birthday

func initCronJobs() {
	cronObject := cron.New()
	fmt.Println("Setting up a cron job")
	cronObject.AddFunc("*/10 */1 * * * *", func() {
		fmt.Println("Job done")
	})
	cronObject.Start()
}

func checkBirthdays() {
	/*
		TODO: query the birthdays which are today
		write a message to the users
	*/
}

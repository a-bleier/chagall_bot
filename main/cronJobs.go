package main

import (
	"fmt"
	"github.com/robfig/cron"
)

func initCronJobs() {
	cronObject := cron.New()
	fmt.Println("Setting up a cron job")
	cronObject.AddFunc("*/10 */1 * * * *", func() {
		fmt.Println("Job done")
	})
	cronObject.Start()
}

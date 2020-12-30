package main

import (
	"encoding/json"
	"fmt"
	"github.com/a-bleier/chagall_bot/comm"
	"github.com/a-bleier/chagall_bot/db"
	"io/ioutil"
	"strconv"
	"sync"
)

type globalConfig struct {
	telegramKey string
}

func (gConf *globalConfig) setupGlobalConfig() {
	byteConfig, _ := ioutil.ReadFile("chagallconfig.json")

	var temp map[string]interface{}
	json.Unmarshal(byteConfig, &temp)
	gConf.telegramKey = temp["telegramKey"].(string)
}

var txQueue comm.SafeQueue

func main() {
	//1264160269

	//Init db
	db.InitChagDB("simple.sqlite")
	fmt.Println(db.CheckUserIsRegistered(strconv.Itoa(1264160269)))

	//Init cron
	initCronJobs()

	//init state machine
	initStates()

	//TODO Store last update in global config
	var gConfig globalConfig
	gConfig.setupGlobalConfig()

	fmt.Println("Chagall says hello")

	mR := sync.Mutex{}
	condR := sync.NewCond(&mR)

	rxQueue := comm.NewSafeQueue()
	listener := comm.NewStub(&rxQueue, 698207968, condR, gConfig.telegramKey, true)

	mT := sync.Mutex{}
	condT := sync.NewCond(&mT)

	txQueue = comm.NewSafeQueue()
	sender := comm.NewStub(&txQueue, 698207968, condT, gConfig.telegramKey, false)

	go listener.Listen()
	go sender.Send()

	//Main loop here
	for {
		condR.L.Lock()
		for rxQueue.IsEmpty() {
			condR.Wait()
		}
		item := rxQueue.DeQueue()
		condR.L.Unlock()

		update := item.Data.(comm.Update)
		fmt.Println(update.Id)
		fmt.Println(update.Message.Text)

		if db.CheckUserIsRegistered(strconv.Itoa(1264160269)) == false {
			fmt.Println("No registered user tries to log in")
			//responseText = "Sorry, I don't know you. Please contact Adrian if you want to use this bot"
			//txQueue.EnQueue(fmt.Sprintf(`{"chat_id" : %d,"text" : "%s"}`, responseChatId, responseText))
			continue
		}
		//responseText = update.Message.Text

		fmt.Println("From: ", update.Message.From)

		transitStates(update)
		//ReplyKeyboardMarkup
		//TODO check user state
		/*
			When fresh conversation -> send list of services [Birthdays | Quit]
			When Quit -> goodbye message
			when birthdays checked -> send options to modify [list | add | remove | exit]
			when add -> ask for name and birthday
			when list or remove -> birthdays checked
			when exit -> fresh conversation
			after 5 minutes all states will timeout to fresh conversation
		*/

		//txQueue.EnQueue(fmt.Sprintf(`{"chat_id" : %d,"text" : "%s"}`, responseChatId, responseText))

		condT.Broadcast()
	}

}

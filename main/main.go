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
	telegramKey string `json:"telegramKey"`
}

func (gConf *globalConfig) setupGlobalConfig() {
	byteConfig, _ := ioutil.ReadFile("chagallconfig.json")

	var temp map[string]interface{}
	json.Unmarshal(byteConfig, &temp)
	gConf.telegramKey = temp["telegramKey"].(string)
}

//TODO: Remove this one
var sender comm.Stub

func main() {

	//Init db
	db.InitChagDB("simple.sqlite")

	//Init cron
	initCronJobs()

	//init state machine
	stateMachine := NewStateMachine()

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

	txQueue := comm.NewSafeQueue()
	sender = comm.NewStub(&txQueue, 698207968, condT, gConfig.telegramKey, false)

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
		fmt.Println("chatId", update.Message.Chat.Id)
		fmt.Println("userId ", update.Message.From.Id)
		fmt.Println(update.Message.Text)

		if db.CheckUserIsRegistered(strconv.Itoa(1264160269)) == false {
			fmt.Println("A unregistered user tries to log in")
			//TODO Send info text to unregistered user
			continue
		}

		stateMachine.transitStates(update)
	}

}

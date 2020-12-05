package main

import (
	"encoding/json"
	"fmt"
	"github.com/a-bleier/chagall_bot/comm"
	"io/ioutil"
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

func main() {

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
	sender := comm.NewStub(&txQueue, 698207968, condT, gConfig.telegramKey, false)

	go listener.Listen()
	go sender.Send()

	//Main loop here
	for {
		condR.L.Lock()
		for rxQueue.IsEmpty() {
			condR.Wait()
		}

		update := rxQueue.DeQueue().(comm.Update)
		fmt.Println(update.Id)
		fmt.Println(update.Message.Text)
		condR.L.Unlock()

		txQueue.EnQueue(fmt.Sprintf(`{"chat_id" : %d,"text" : "%s"}`, update.Message.Mchat.ID, update.Message.Text))

		condT.Broadcast()
	}

}

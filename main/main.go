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

	var gConfig globalConfig
	gConfig.setupGlobalConfig()

	fmt.Println("Chagall says hello")

	m := sync.Mutex{}
	cond := sync.NewCond(&m)

	queue := comm.NewSafeRxQueue()
	server := comm.NewServer(&queue, 698207968, cond, gConfig.telegramKey)

	go server.Run()

	//Main loop here
	for {
		cond.L.Lock()
		for queue.IsEmpty() {
			cond.Wait()
		}

		update := queue.DeQueue()
		fmt.Println(update.Id)
		fmt.Println(update.Message.Text)
		cond.L.Unlock()
	}

}

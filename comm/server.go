package comm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
)

//Server receives updates and sends messages.
type Server struct {
	rx         *SafeRxQueue
	cond       *sync.Cond
	lastUpdate int
	apiKey     string
}

//Run will listen and serve in its own routine. Call this function in its own routine
func (s *Server) Run() {
	for true {
		s.tick()
	}

}

func (s *Server) tick() {
	//receive something

	//TODO Read out the the api key From json config file

	url := "https://api.telegram.org/bot" + s.apiKey + "/getUpdates"
	url += "?offset=" + strconv.Itoa(s.lastUpdate+1) + "&timeout=60"
	//fmt.Print(url)
	resp, err := http.Get(url)

	if err != nil {
		fmt.Print("Error")
	}

	response, _ := ioutil.ReadAll(resp.Body)

	var formattedJSONBuffer bytes.Buffer
	json.Indent(&formattedJSONBuffer, response, "", "\t")

	//fmt.Printf("%s", formattedJSONBuffer.Bytes())

	defer resp.Body.Close()

	//put in rxQueue
	var obj map[string]interface{}
	err = json.Unmarshal(response, &obj)
	if err != nil {
		panic(err)
	}
	result := obj["result"].([]interface{})
	fmt.Println("me here")
	for _, res := range result {
		update := NewUpdateFromJSON(res)
		s.lastUpdate = update.Id
		s.rx.EnQueue(update)
	}

	//wake up main routine

	s.cond.Signal()

}

//NewServer returns a new Server
func NewServer(rx *SafeRxQueue, lastUpdateID int, cond *sync.Cond, apiKey string) Server {
	return Server{
		rx:         rx,
		cond:       cond,
		lastUpdate: lastUpdateID,
		apiKey:     apiKey,
	}
}

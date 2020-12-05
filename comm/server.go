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

//Stub receives updates and sends messages.
type Stub struct {
	queue      *SafeQueue
	cond       *sync.Cond
	lastUpdate int
	apiKey     string
	isListener bool
}

//Run will listen and serve in its own routine. Call this function in its own routine
func (s *Stub) Listen() {

	if s.isListener == false {
		fmt.Println("This stub is a Sender ! Don't listen here")
		return
	}
	for true {
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
			s.queue.EnQueue(update)
		}

		//wake up main routine

		s.cond.Signal()

	}

}

func (s *Stub) Send() {
	if s.isListener == true {
		fmt.Println("This stub is a Listener ! Don't send here")
		return
	}

	for {
		s.cond.L.Lock()
		for s.queue.IsEmpty() {
			s.cond.Wait()
		}
		text := s.queue.DeQueue().(string)
		s.cond.L.Unlock()

		fmt.Println("Gonna send", text)

		url := "https://api.telegram.org/bot" + s.apiKey + "/sendMessage"
		jsonStr := []byte(text)

		resp, _ := http.Post(url, "application/json", bytes.NewBuffer(jsonStr))

		resp.Body.Close()

		fmt.Println("response Status:", resp.Status)
		fmt.Println("response Headers:", resp.Header)

	}

	//Do some sending here
}

//NewServer returns a new Stub
func NewStub(rx *SafeQueue, lastUpdateID int, cond *sync.Cond, apiKey string, isListener bool) Stub {
	return Stub{
		queue:      rx,
		cond:       cond,
		lastUpdate: lastUpdateID,
		apiKey:     apiKey,
		isListener: isListener,
	}
}

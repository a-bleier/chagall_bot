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

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("error")
			return
		}
		defer resp.Body.Close()

		var apiResp APIResponse
		err = json.Unmarshal(data, &apiResp)
		if err != nil {
			panic(err)
		}
		var updates []Update
		json.Unmarshal(apiResp.Result, &updates)

		for _, u := range updates {
			item := QueueItem{u, "update"}
			s.lastUpdate = u.Id
			s.queue.EnQueue(item)
		}
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
		//text := s.queue.DeQueue().(string)
		item := s.queue.DeQueue()
		s.cond.L.Unlock()

		jsonStr := item.Data.([]byte)
		method := item.Info
		//fmt.Println("Gonna send", text)
		fmt.Println(method)
		url := "https://api.telegram.org/bot" + s.apiKey + "/" + method
		//jsonStr := []byte(text)

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

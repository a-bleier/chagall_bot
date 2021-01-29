package comm

import (
	"bytes"
	"encoding/json"
	"github.com/a-bleier/chagall_bot/logging"
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
		logging.LogWarning("There was a try to use this sending stub as a listener.")
		return
	}
	for true {
		url := "https://api.telegram.org/bot" + s.apiKey + "/getUpdates"
		url += "?offset=" + strconv.Itoa(s.lastUpdate+1) + "&timeout=60"
		//fmt.Print(url)
		resp, err := http.Get(url)

		if err != nil {
			logging.LogFatalError("Couldn't fetch " + url)
		}

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logging.LogFatalError("Couldn't read the body of the response")
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
		logging.LogWarning("There was a try to use this listening stub as a sender.")
		return
	}

	for {
		s.cond.L.Lock()
		for s.queue.IsEmpty() {
			s.cond.Wait()
		}
		item := s.queue.DeQueue()
		s.cond.L.Unlock()

		jsonStr := item.Data.([]byte)
		method := item.Info
		logging.LogInfo("Using method " + method)

		url := "https://api.telegram.org/bot" + s.apiKey + "/" + method
		resp, _ := http.Post(url, "application/json", bytes.NewBuffer(jsonStr))

		resp.Body.Close()

		//fmt.Println("response Status: ", resp.Status)
		logging.LogInfo("response Status" + resp.Status)

	}
}

//Could turn out wonky due to race conditions
func (s *Stub) AddMessageToTx(v interface{}, dataType string) {
	data, err := json.Marshal(v)
	if err != nil {
		logging.LogFatalError("Couldn't queue a message into tx queue")
		panic(err)
	}
	item := QueueItem{data, dataType}
	s.queue.EnQueue(item)
	s.cond.Broadcast()
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

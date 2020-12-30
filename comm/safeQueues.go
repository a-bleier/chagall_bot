package comm

import (
	"sync"
)

//SafeQueue : All operations shall be thread safe to use

//TODO make a queue item

type QueueItem struct {
	Data interface{}
	Info string
}
type SafeQueue struct {
	Mu      sync.Mutex
	UpdateQ []QueueItem
}

func (s *SafeQueue) IsEmpty() bool {
	s.Mu.Lock()
	retVal := len(s.UpdateQ) == 0
	s.Mu.Unlock()
	return retVal
}

func (s *SafeQueue) EnQueue(item QueueItem) {
	s.Mu.Lock()
	s.UpdateQ = append(s.UpdateQ, item)
	s.Mu.Unlock()

}

func (s *SafeQueue) DeQueue() QueueItem {
	s.Mu.Lock()
	qi := s.UpdateQ[0]
	s.UpdateQ = s.UpdateQ[1:]
	s.Mu.Unlock()
	return qi
}

func NewSafeQueue() SafeQueue {
	return SafeQueue{}
}

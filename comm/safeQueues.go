package comm

import (
	"sync"
)

//SafeRxQueue : All operations shall be thread safe to use
type SafeRxQueue struct {
	Mu      sync.Mutex
	UpdateQ []Update
}

func (s *SafeRxQueue) IsEmpty() bool {
	s.Mu.Lock()
	retVal := len(s.UpdateQ) == 0
	s.Mu.Unlock()
	return retVal
}

func (s *SafeRxQueue) EnQueue(u Update) {
	s.Mu.Lock()
	s.UpdateQ = append(s.UpdateQ, u)
	s.Mu.Unlock()

}

func (s *SafeRxQueue) DeQueue() Update {
	s.Mu.Lock()
	u := s.UpdateQ[0]
	s.UpdateQ = s.UpdateQ[1:]
	s.Mu.Unlock()
	return u
}

func NewSafeRxQueue() SafeRxQueue {
	return SafeRxQueue{}
}

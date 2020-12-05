package comm

import (
	"sync"
)

//SafeQueue : All operations shall be thread safe to use
type SafeQueue struct {
	Mu      sync.Mutex
	UpdateQ []interface{}
}

func (s *SafeQueue) IsEmpty() bool {
	s.Mu.Lock()
	retVal := len(s.UpdateQ) == 0
	s.Mu.Unlock()
	return retVal
}

func (s *SafeQueue) EnQueue(u interface{}) {
	s.Mu.Lock()
	s.UpdateQ = append(s.UpdateQ, u)
	s.Mu.Unlock()

}

func (s *SafeQueue) DeQueue() interface{} {
	s.Mu.Lock()
	u := s.UpdateQ[0]
	s.UpdateQ = s.UpdateQ[1:]
	s.Mu.Unlock()
	return u
}

func NewSafeQueue() SafeQueue {
	return SafeQueue{}
}

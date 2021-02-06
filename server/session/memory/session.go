package memory

import (
	"sync"
	"time"
)

type MemorySession struct {
	mtx       sync.Mutex
	id        string
	startTime time.Time
	datas     map[string]interface{}
}

func (s *MemorySession) Set(k string, v interface{}) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.datas[k] = v
}

func (s *MemorySession) Get(k string) interface{} {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	d, ok := s.datas[k]
	if ok {
		return d
	}

	return nil
}

func (s *MemorySession) Del(k string) bool {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if _, ok := s.datas[k]; ok {
		delete(s.datas, k)
		return true
	}

	return false
}

func (s *MemorySession) ID() string {
	return s.id
}

func (s *MemorySession) StartTime() time.Time {
	return s.startTime
}

func (s *MemorySession) Update() bool {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.startTime = time.Now()
	return true
}

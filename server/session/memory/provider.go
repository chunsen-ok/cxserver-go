package memory

import (
	"cxfw/session"
	"sync"
	"time"
)

type MemoryProvider struct {
	mtx      sync.Mutex
	sessions map[string]*MemorySession
}

func NewProvider() *MemoryProvider {
	return &MemoryProvider{
		mtx:      sync.Mutex{},
		sessions: make(map[string]*MemorySession),
	}
}

func (s *MemoryProvider) NewSession(sessionID string) session.ISession {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	se := &MemorySession{
		mtx:       sync.Mutex{},
		id:        sessionID,
		startTime: time.Now(),
		datas:     make(map[string]interface{}),
	}

	s.sessions[sessionID] = se

	return se
}

func (s *MemoryProvider) DelSession(sessionID string) bool {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if _, ok := s.sessions[sessionID]; ok {
		delete(s.sessions, sessionID)
		return true
	}

	return false
}

func (s *MemoryProvider) GetSession(sessionID string) session.ISession {
	if se, ok := s.sessions[sessionID]; ok {
		return se
	}

	return nil
}

// 删除过期 session
func (s *MemoryProvider) GC(maxLifeTime int) {
	now := time.Now().Unix()

	s.mtx.Lock()
	defer s.mtx.Unlock()

	expired := make([]string, 0)
	for k, se := range s.sessions {
		if se.startTime.Unix()+int64(maxLifeTime) < now {
			expired = append(expired, k)
		}
	}

	for _, k := range expired {
		delete(s.sessions, k)
	}
}

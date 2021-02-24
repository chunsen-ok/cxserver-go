package memory

import (
	"cxfw/session/ses"
	"log"
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

func (s *MemoryProvider) NewSession(sessionID string) ses.ISession {
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

func (s *MemoryProvider) GetSession(sessionID string) ses.ISession {
	if se, ok := s.sessions[sessionID]; ok {
		return se
	}

	return nil
}

func (s *MemoryProvider) UpdateSession(sessionID string) ses.ISession {
	if se, ok := s.sessions[sessionID]; ok {
		se.Update()
		return se
	}

	return nil
}

// 删除过期 session
func (s *MemoryProvider) GC(maxLifeTime int64) {
	now := time.Now().Unix()
	log.Println("provider GC: ", now)

	s.mtx.Lock()
	defer s.mtx.Unlock()

	expired := make([]string, 0)
	for k, se := range s.sessions {
		log.Println("session :", k, " time: ", se.startTime.Unix(), "max life time: ", maxLifeTime)
		if se.startTime.Unix()+maxLifeTime < now {
			expired = append(expired, k)
		}
	}

	for _, k := range expired {
		log.Println("expired:", k)
		delete(s.sessions, k)
	}
}

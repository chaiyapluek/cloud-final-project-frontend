package service

import (
	"math/rand"
	"sync"

	"dev.chaiyapluek.cloud.final.frontend/template/pages/order"
)

type SessionService interface {
	GetSessionDetail(sessionId string) *SessionDetail
	NewSessionDetail() (string, *SessionDetail)
	UpdateSessionDetail(sessionId string, sessionDetail *SessionDetail)
	DeleteSessionDetail(sessionId string)
}

type SessionDetail struct {
	CurrentMenu     *order.OrderProps
	Preferences     map[string]interface{}
	CartId          string
	CurrentLocation string
}

type inMemorySessionService struct {
	mutex   sync.RWMutex
	session map[string]*SessionDetail
}

func NewInMemorySessionService() SessionService {
	return &inMemorySessionService{
		session: map[string]*SessionDetail{},
	}
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randomSessionId() string {
	b := make([]byte, 16)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func (s *inMemorySessionService) GetSessionDetail(sessionId string) *SessionDetail {
	s.mutex.RLock()
	detail, ok := s.session[sessionId]
	s.mutex.RUnlock()
	if !ok {
		return nil
	}
	return detail
}

func (s *inMemorySessionService) NewSessionDetail() (string, *SessionDetail) {
	var id string = randomSessionId()
	var ok bool
	s.mutex.Lock()
	_, ok = s.session[id]
	for ok {
		id = randomSessionId()
		_, ok = s.session[id]
	}
	s.session[id] = &SessionDetail{
		CurrentMenu: &order.OrderProps{},
		Preferences: map[string]interface{}{},
		CartId:      "",
	}
	defer s.mutex.Unlock()
	return id, s.session[id]
}

func (s *inMemorySessionService) UpdateSessionDetail(sessionId string, sessionDetail *SessionDetail) {
	s.mutex.Lock()
	s.session[sessionId] = sessionDetail
	s.mutex.Unlock()
}

func (s *inMemorySessionService) DeleteSessionDetail(sessionId string) {
	s.mutex.Lock()
	delete(s.session, sessionId)
	s.mutex.Unlock()
}

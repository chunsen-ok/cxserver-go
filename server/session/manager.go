package session

import (
	"cxfw/session/memory"
	"cxfw/session/ses"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const gcInterval = 30 * time.Minute

var instance *CookieSessionManager

type CookieSessionManager struct {
	maxLifeTime int64
	cookieName  string
	provider    ses.IProvider
}

func Init(cookieName string, maxLifeTime int64) *CookieSessionManager {
	instance = &CookieSessionManager{
		maxLifeTime: maxLifeTime,
		cookieName:  cookieName,
		provider:    memory.NewProvider(),
	}

	return instance
}

func S() *CookieSessionManager {
	return instance
}

func (s *CookieSessionManager) StartSession(c *gin.Context) ses.ISession {
	id := ses.SessionID()
	se := s.provider.NewSession(id)

	cookie := http.Cookie{
		Name:       s.cookieName,
		Value:      id,
		Path:       "/",
		Domain:     "",
		Expires:    time.Now().Add(time.Duration(s.maxLifeTime)),
		RawExpires: "",
		MaxAge:     int(s.maxLifeTime),
		Secure:     true,
		HttpOnly:   true,
		SameSite:   http.SameSiteLaxMode,
		Raw:        "",
		Unparsed:   []string{},
	}

	http.SetCookie(c.Writer, &cookie)

	return se
}

func (s *CookieSessionManager) StopSession(c *gin.Context) bool {
	cookieVal, err := c.Cookie(s.cookieName)
	if len(cookieVal) == 0 || err != nil {
		return false
	}

	ok := s.provider.DelSession(cookieVal)
	if !ok {
		return false
	}

	cookie := http.Cookie{
		Name:     s.cookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now(),
		MaxAge:   -1,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(c.Writer, &cookie)

	return true
}

func (s *CookieSessionManager) GetSession(c *gin.Context) ses.ISession {
	sessionID, err := c.Cookie(s.cookieName)
	if err != nil {
		return nil
	}

	se := s.provider.GetSession(sessionID)
	return se
}

func (s *CookieSessionManager) GC() {
	s.provider.GC(s.maxLifeTime)
	time.AfterFunc(gcInterval, func() { s.GC() })
}

package sessionMgr

import (
	"cxfw/session"
	"cxfw/session/memory"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const gcInterval = 30 * time.Minute

type CookieSessionManager struct {
	maxLifeTime int
	cookieName  string
	provider    session.IProvider
}

func Init(cookieName string, maxLifeTime int) *CookieSessionManager {
	s := &CookieSessionManager{
		maxLifeTime: maxLifeTime,
		cookieName:  cookieName,
		provider:    memory.NewProvider(),
	}

	return s
}

func (s *CookieSessionManager) StartSession(c *gin.Context) session.ISession {
	id := sessionID()
	se := s.provider.NewSession(id)

	cookie := http.Cookie{
		Name:     s.cookieName,
		Value:    id,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(time.Duration(s.maxLifeTime)),
		MaxAge:   s.maxLifeTime,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
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
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(c.Writer, &cookie)

	return true
}

func (s *CookieSessionManager) GetSession(c *gin.Context) session.ISession {
	sessionID, err := c.Cookie(s.cookieName)
	if err != nil {
		return nil
	}

	return s.provider.GetSession(sessionID)
}

func (s *CookieSessionManager) GC() {
	s.provider.GC(s.maxLifeTime)
	time.AfterFunc(gcInterval, func() { s.GC() })
}

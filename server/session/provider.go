package session

type IProvider interface {
	NewSession(sessionID string) ISession
	DelSession(sessionID string) bool
	GetSession(sessionID string) ISession
	GC(maxLifeTime int)
}

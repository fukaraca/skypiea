package session

import (
	"crypto/rand"
	"time"

	"github.com/dgraph-io/ristretto/v2"
)

const DefaultCookieName = "ss_skypiea"

type storage struct {
	*ristretto.Cache[string, *Session]
}

type Manager struct {
	cache      *storage
	defaultTTL time.Duration
	cookieName string
}

type Session struct {
	id        string
	data      map[string]any
	createdAt time.Time
	updatedAt time.Time
}

func New() *Session {
	return &Session{
		id:        rand.Text(),
		data:      nil,
		createdAt: time.Now(),
	}
}

func NewSessionManager(maxSize int64, ttl time.Duration) *Manager {
	c, err := ristretto.NewCache(&ristretto.Config[string, *Session]{
		NumCounters: 10 * maxSize,
		MaxCost:     maxSize,
		BufferItems: 64,
	})
	if err != nil {
		return nil
	}
	return &Manager{
		cache:      &storage{c},
		defaultTTL: ttl,
		cookieName: DefaultCookieName,
	}
}

func (sm *Manager) SetSession(userID string, sess *Session, ttl time.Duration) {
	sm.cache.SetWithTTL(userID, sess, 0, ttl)
	sm.cache.Wait()
}

func (sm *Manager) ValidateSession(userID string) bool {
	v, found := sm.cache.Get(userID)
	return found && v != nil
}

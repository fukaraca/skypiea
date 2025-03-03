package session

import (
	"time"

	"github.com/dgraph-io/ristretto/v2"
)

type Manager struct {
	cache      *ristretto.Cache[string, *Session]
	defaultTTL time.Duration
}

type Session struct{}

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
		cache:      c,
		defaultTTL: ttl,
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

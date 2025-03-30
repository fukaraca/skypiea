package session

import (
	"fmt"
	"time"

	"github.com/fukaraca/skypiea/pkg/cache"
	"github.com/fukaraca/skypiea/pkg/gwt"
	"github.com/google/uuid"
)

const DefaultCookieName = "ss_skypiea"

var Cache *Manager

type Manager struct {
	cache      *cache.Storage
	defaultTTL time.Duration
	jwtManager gwt.Manager
}

type Session struct {
	ID        string
	UserID    uuid.UUID
	token     string
	createdAt time.Time
	updatedAt time.Time
	EOL       time.Time
}

func (sm *Manager) New(userID uuid.UUID) *Session {
	if uuid.Validate(userID.String()) != nil {
		userID = uuid.New()
	}
	t0 := time.Now()
	id := fmt.Sprintf("%s.%d", userID.String(), t0.UnixNano())
	// Get user details from DB
	tkn, err := sm.jwtManager.GenerateToken(userID.String(), "admin")
	if err != nil {
		return nil
	}
	return &Session{
		ID:        id,
		UserID:    userID,
		token:     tkn,
		createdAt: t0,
		updatedAt: t0,
	}
}

func (s *Session) Valid() bool {
	return s.EOL.After(time.Now().Add(-time.Second * 5))
}

func (s *Session) Token() string {
	return s.token
}

func NewManager(jwtConfig *gwt.Config, ttl time.Duration) *Manager {
	return &Manager{
		cache:      cache.New(),
		defaultTTL: ttl,
		jwtManager: gwt.NewJWTService(jwtConfig),
	}
}

func (sm *Manager) Set(sess *Session) {
	sess.EOL = sess.updatedAt.Add(sm.defaultTTL)
	sm.cache.Set(sess.ID, sess)
}

func (sm *Manager) Get(sessionID string) *Session {
	sess, _ := sm.cache.Get(sessionID).(*Session)
	return sess
}

func (sm *Manager) RevokeAllSessions(userID string) {
	sm.cache.DeleteByPrefix(userID)
}

func (sm *Manager) RefreshSession(sess *Session) {
	sess.updatedAt = time.Now()
	sm.Set(sess)
}

func (sm *Manager) ValidateSession(sessionID string) bool {
	sess := sm.Get(sessionID)
	if sess == nil || !sess.Valid() {
		return false
	}
	sm.RefreshSession(sess)
	return true
}

func (sm *Manager) ValidateToken(tkn string) bool {
	_, err := sm.jwtManager.ValidateToken(tkn)
	return err == nil
}

func (sm *Manager) GetJWTBySessionID(sessionID string) *gwt.Token {
	v := sm.Get(sessionID)
	if v != nil && v.Valid() {
		tkn, _ := sm.jwtManager.ValidateToken(v.token)
		return tkn
	}
	return nil
}

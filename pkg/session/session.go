package session

import (
	"context"
	"fmt"
	"time"

	"github.com/fukaraca/skypiea/internal/model"
	"github.com/fukaraca/skypiea/internal/storage"
	"github.com/fukaraca/skypiea/pkg/cache"
	"github.com/fukaraca/skypiea/pkg/gwt"
	"github.com/google/uuid"
)

const (
	DefaultCookieName   = "ss_skypiea"
	DefaultCookieMaxAge = 1000
	DefaultCookieDomain = "localhost"

	DefaultSessionEarlyTimeout = -time.Second * 5

	CtxLoggedIn = "logged_in"
)

type UserReader interface {
	GetUserByUUID(context.Context, uuid.UUID) (*storage.User, error)
}

var Cache *Manager

type Manager struct {
	cache      *cache.Storage
	repo       UserReader
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

type Cookie struct {
	Name, Value, Path, Domain string
	MaxAge                    int
	Secure, HTTPOnly          bool
}

func NewCookie(id string) *Cookie {
	return &Cookie{
		Name:     DefaultCookieName,
		Value:    id,
		Path:     model.PathMain,
		Domain:   DefaultCookieDomain,
		MaxAge:   DefaultCookieMaxAge,
		Secure:   false,
		HTTPOnly: true,
	}
}

func (sm *Manager) NewSession(ctx context.Context, userID uuid.UUID) *Session {
	if uuid.Validate(userID.String()) != nil {
		userID = uuid.New()
	}
	t0 := time.Now()
	id := fmt.Sprintf("%s.%d", userID.String(), t0.UnixNano())
	// Get user details from DB
	user, err := sm.repo.GetUserByUUID(ctx, userID)
	if err != nil {
		return nil
	}
	tkn, err := sm.jwtManager.GenerateToken(userID.String(), user.Role)
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
	return s.EOL.After(time.Now().Add(DefaultSessionEarlyTimeout))
}

func (s *Session) Token() string {
	return s.token
}

func NewManager(jwtConfig *gwt.Config, repo UserReader, ttl time.Duration) *Manager {
	return &Manager{
		cache:      cache.New(),
		repo:       repo,
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

func (sm *Manager) ValidateSession(sessionID string) (*Session, bool) {
	sess := sm.Get(sessionID)
	if sess == nil || !sess.Valid() {
		return nil, false
	}
	sm.RefreshSession(sess)
	return sess, true
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

func (sm *Manager) GetUserUUIDByToken(tkn string) *uuid.UUID {
	t, err := sm.jwtManager.ValidateToken(tkn)
	if err != nil {
		return nil
	}
	uid := uuid.MustParse(t.UserID)
	return &uid
}

func (sm *Manager) Delete(sessionID string) {
	sm.cache.Del(sessionID)
}

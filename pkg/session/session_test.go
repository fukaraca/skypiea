package session_test

import (
	"context"
	"testing"
	"time"

	"github.com/fukaraca/skypiea/internal/storage"
	"github.com/fukaraca/skypiea/pkg/gwt"
	"github.com/fukaraca/skypiea/pkg/session"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockDBRegistery struct{}

func (r *mockDBRegistery) GetUserByUUID(ctx context.Context, userUUID uuid.UUID) (*storage.User, error) {
	return &storage.User{UserUUID: userUUID.String()}, nil
}

func TestSessionManager(t *testing.T) {
	jwtConfig := &gwt.Config{Secret: "secret", Domain: "localhost"}

	ttl := time.Minute * 30
	manager := session.NewManager(jwtConfig, &mockDBRegistery{}, ttl)

	userID := uuid.New()

	t.Run("Create and Retrieve Session", func(t *testing.T) {
		sess := manager.NewSession(t.Context(), userID)
		assert.NotNil(t, sess)
		assert.Equal(t, userID, sess.UserID)

		manager.Set(sess)

		retrieved := manager.Get(sess.ID)
		assert.NotNil(t, retrieved)
		assert.Equal(t, sess.ID, retrieved.ID)
	})

	t.Run("Validate Session", func(t *testing.T) {
		sess := manager.NewSession(t.Context(), userID)
		manager.Set(sess)

		_, valid := manager.ValidateSession(sess.ID)
		assert.True(t, valid)
	})

	t.Run("Refresh Session", func(t *testing.T) {
		sess := manager.NewSession(t.Context(), userID)
		manager.Set(sess)

		oldUpdatedAt := sess.EOL
		time.Sleep(time.Millisecond * 10)

		manager.RefreshSession(sess)
		retrieved := manager.Get(sess.ID)

		assert.True(t, retrieved.EOL.After(oldUpdatedAt))
	})

	t.Run("Revoke All Sessions", func(t *testing.T) {
		sess1 := manager.NewSession(t.Context(), userID)
		sess2 := manager.NewSession(t.Context(), userID)
		manager.Set(sess1)
		manager.Set(sess2)

		manager.RevokeAllSessions(userID.String())

		assert.Nil(t, manager.Get(sess1.ID))
		assert.Nil(t, manager.Get(sess2.ID))
	})

	t.Run("Validate JWT Token", func(t *testing.T) {
		sess := manager.NewSession(t.Context(), userID)
		manager.Set(sess)

		valid := manager.ValidateToken(sess.Token())
		assert.True(t, valid)
	})

	t.Run("Get JWT by Session ID", func(t *testing.T) {
		sess := manager.NewSession(t.Context(), userID)
		manager.Set(sess)

		jwt := manager.GetJWTBySessionID(sess.ID)
		assert.NotNil(t, jwt)
		assert.Equal(t, userID.String(), jwt.Subject)
	})
}

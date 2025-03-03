package session

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSession(t *testing.T) {
	testSM := NewSessionManager(1000, time.Minute)
	testSM.SetSession("1", &Session{}, time.Minute)
	assert.True(t, testSM.ValidateSession("1"))
}

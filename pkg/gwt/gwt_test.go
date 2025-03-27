package gwt

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWT(t *testing.T) {
	testService := NewJWTService(&Config{Secret: []byte("secret")})
	token, err := testService.GenerateToken("1", "admin")
	require.NoError(t, err)
	assert.NotNil(t, token)
	claims, err := testService.ValidateToken(token)
	require.NoError(t, err)
	require.Equal(t, "1", claims.UserID)
}

package jwt_test

import (
	"testing"
	"time"

	"github.com/minhhoccode111/go-clean-template-gin/pkg/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testSecret   = "test-secret"
	testUserID   = int64(42)
	testUsername = "testuser"
)

var testUserRoles = []string{"user"} //nolint:gochecknoglobals // intended

func TestGenerateAndValidateToken(t *testing.T) {
	t.Parallel()

	token, err := jwt.GenerateToken(testUserID, testUsername, testSecret, testUserRoles, time.Hour)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	claims, err := jwt.ValidateToken(token, testSecret)
	require.NoError(t, err)
	assert.Equal(t, testUserID, claims.UserID)
	assert.Equal(t, testUsername, claims.Username)
}

func TestValidateToken_WrongSecret(t *testing.T) {
	t.Parallel()

	token, err := jwt.GenerateToken(testUserID, testUsername, testSecret, testUserRoles, time.Hour)
	require.NoError(t, err)

	_, err = jwt.ValidateToken(token, "wrong-secret")
	assert.Error(t, err)
}

func TestValidateToken_Expired(t *testing.T) {
	t.Parallel()

	token, err := jwt.GenerateToken(
		testUserID,
		testUsername,
		testSecret,
		testUserRoles,
		-time.Second,
	)
	require.NoError(t, err)

	_, err = jwt.ValidateToken(token, testSecret)
	assert.Error(t, err)
}

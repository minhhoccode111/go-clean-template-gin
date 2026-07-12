package jwt_test

import (
	"testing"
	"time"

	"github.com/minhhoccode111/go-clean-template-gin/pkg/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testSecret = "test-secret"
	testUserID = "42"
)

func TestGenerateAndParseToken(t *testing.T) {
	t.Parallel()

	manager := jwt.New(testSecret, time.Hour)

	token, err := manager.GenerateToken(testUserID)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	parsedUserID, err := manager.ParseToken(token)
	require.NoError(t, err)
	assert.Equal(t, testUserID, parsedUserID)
}

func TestParseToken_WrongSecret(t *testing.T) {
	t.Parallel()

	manager := jwt.New(testSecret, time.Hour)

	token, err := manager.GenerateToken(testUserID)
	require.NoError(t, err)

	differentManager := jwt.New("wrong-secret", time.Hour)
	_, err = differentManager.ParseToken(token)
	assert.Error(t, err)
}

func TestParseToken_Expired(t *testing.T) {
	t.Parallel()

	manager := jwt.New(testSecret, -time.Second)

	token, err := manager.GenerateToken(testUserID)
	require.NoError(t, err)

	_, err = manager.ParseToken(token)
	assert.Error(t, err)
}

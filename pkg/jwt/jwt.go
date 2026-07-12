// Package jwt provides JWT generation and validation.
package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Manager -.
type Manager struct {
	secret      string
	tokenExpiry time.Duration
}

// New -.
func New(secret string, tokenExpiry time.Duration) *Manager {
	return &Manager{
		secret:      secret,
		tokenExpiry: tokenExpiry,
	}
}

// GenerateToken creates a signed HS256 JWT containing the userID.
func (m *Manager) GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(m.tokenExpiry).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(m.secret))
	if err != nil {
		return "", fmt.Errorf("JWT - GenerateToken - token.SignedString: %w", err)
	}

	return signed, nil
}

// ParseToken validates and returns the userID from a JWT string.
func (m *Manager) ParseToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenSignatureInvalid
		}

		return []byte(m.secret), nil
	})
	if err != nil {
		return "", fmt.Errorf("JWT - ParseToken - jwt.Parse: %w", err)
	}

	if !token.Valid {
		return "", jwt.ErrTokenInvalidClaims
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", jwt.ErrTokenInvalidClaims
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return "", jwt.ErrTokenInvalidClaims
	}

	return userID, nil
}

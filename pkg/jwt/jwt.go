// Package jwt provides JWT generation and validation for the application.
package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ClaimsJWT holds the custom JWT claims.
type ClaimsJWT struct {
	UserID   int64    `json:"user_id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	jwt.RegisteredClaims
}

// GenerateToken creates a signed HS256 JWT for the given user.
func GenerateToken(
	userID int64,
	username, secret string,
	roles []string,
	ttl time.Duration,
) (string, error) {
	claims := ClaimsJWT{
		UserID:   userID,
		Username: username,
		Roles:    roles,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprintf("%d", userID),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
}

// ValidateToken parses and validates a JWT string, returning the typed claims.
func ValidateToken(tokenStr, secret string) (*ClaimsJWT, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &ClaimsJWT{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenSignatureInvalid
		}

		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*ClaimsJWT)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}

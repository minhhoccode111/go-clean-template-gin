package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/minhhoccode111/go-clean-template-gin/pkg/jwt"
)

const (
	_bearerParts = 2
	errorKey     = "error"
)

// Auth returns a JWT authentication middleware for Gin.
func Auth(jwtManager *jwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{errorKey: "missing authorization header"})

			return
		}

		parts := strings.SplitN(header, " ", _bearerParts)
		if len(parts) != _bearerParts || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{errorKey: "invalid authorization header format"})

			return
		}

		userID, err := jwtManager.ParseToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{errorKey: "invalid or expired token"})

			return
		}

		c.Set("userID", userID)

		c.Next()
	}
}

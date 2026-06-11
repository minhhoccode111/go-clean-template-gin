package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	pkgjwt "github.com/minhhoccode111/go-clean-template-gin/pkg/jwt"
)

type ctxKey string

const (
	CtxUserIDKey ctxKey = "userID"
	CtxUsername  ctxKey = "username"
	CtxUserRoles ctxKey = "userRoles"
)

func extractToken(c *gin.Context) string {
	const scheme = "Bearer "
	if h := c.GetHeader("Authorization"); len(h) > len(scheme) && h[:len(scheme)] == scheme {
		return h[len(scheme):]
	}

	return ""
}

// Auth validates the JWT and stores user claims in the Gin context.
func Auth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token == "" {
			messageResponse(c, http.StatusUnauthorized, "Unauthorized")

			return
		}

		claims, err := pkgjwt.ValidateToken(token, secret)
		if err != nil {
			messageResponse(c, http.StatusUnauthorized, "Unauthorized")

			return
		}

		c.Set(string(CtxUserIDKey), claims.UserID)
		c.Set(string(CtxUsername), claims.Username)
		c.Set(string(CtxUserRoles), claims.Roles)
		c.Next()
	}
}

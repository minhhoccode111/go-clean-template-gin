package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type MatchMode string

const (
	MatchModeAny MatchMode = "any"
	MatchModeAll MatchMode = "all"
)

// Authz validates that the authenticated user has the required roles.
//
//nolint:gocognit,cyclop // role-checking logic is inherently branchy
func Authz(required []string, mode MatchMode) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRolesRaw, exists := c.Get(string(CtxUserRoles))
		if !exists {
			messageResponse(c, http.StatusForbidden, "Forbidden")

			return
		}

		userRoles, ok := userRolesRaw.([]string)
		if !ok {
			messageResponse(c, http.StatusForbidden, "Forbidden")

			return
		}

		userRoleSet := make(map[string]bool)
		for _, role := range userRoles {
			userRoleSet[role] = true
		}

		switch mode {
		case MatchModeAny:
			for _, reqRole := range required {
				if userRoleSet[reqRole] {
					c.Next()

					return
				}
			}

			messageResponse(c, http.StatusForbidden, "Forbidden")
		case MatchModeAll:
			for _, reqRole := range required {
				if !userRoleSet[reqRole] {
					messageResponse(c, http.StatusForbidden, "Forbidden")

					return
				}
			}

			c.Next()
		default:
			messageResponse(c, http.StatusForbidden, "Forbidden")
		}
	}
}

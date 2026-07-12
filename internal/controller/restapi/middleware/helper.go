package middleware

import "github.com/gin-gonic/gin"

type contextKey string

const (
	CtxUserIDKey contextKey = "userID"
	CtxUserRoles contextKey = "userRoles"
)

type msg struct {
	Message string `json:"message"`
}

func messageResponse(c *gin.Context, code int, s string) {
	c.AbortWithStatusJSON(code, msg{Message: s})
}

package middleware

import "github.com/gin-gonic/gin"

type msg struct {
	Message string `json:"message"`
}

func messageResponse(c *gin.Context, code int, s string) {
	c.AbortWithStatusJSON(code, msg{Message: s})
}

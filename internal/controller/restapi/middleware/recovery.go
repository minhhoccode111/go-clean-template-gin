package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/minhhoccode111/go-clean-template-gin/pkg/logger"
)

func Recovery(l logger.Interface) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				l.Error("%s - %s %s PANIC DETECTED: %v\n%s\n",
					c.ClientIP(),
					c.Request.Method,
					c.Request.URL.Path,
					err,
					debug.Stack(),
				)

				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()

		c.Next()
	}
}

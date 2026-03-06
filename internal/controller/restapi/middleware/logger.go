package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/minhhoccode111/go-clean-template-gin/pkg/logger"
)

func Logger(l logger.Interface) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		l.Info("%s - %s %s - %d %d",
			c.ClientIP(),
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			c.Writer.Size(),
		)
	}
}

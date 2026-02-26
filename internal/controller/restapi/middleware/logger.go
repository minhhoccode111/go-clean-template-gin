package middleware

import (
	"github.com/evrone/go-clean-template/pkg/logger"
	"github.com/gin-gonic/gin"
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

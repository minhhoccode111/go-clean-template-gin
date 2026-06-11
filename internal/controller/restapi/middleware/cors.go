package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/minhhoccode111/go-clean-template-gin/config"
)

// CORS returns a Gin middleware that sets the appropriate CORS headers.
// It supports multiple comma-separated origins in cfg.AllowOrigins and
// reflects the matched origin back to the client, which is required when
// credentials are enabled. A Vary: Origin header is always added so that
// caches do not serve one origin's response to another.
func CORS(cfg *config.Config) gin.HandlerFunc {
	allowedOrigins := make(map[string]struct{})

	for o := range strings.SplitSeq(cfg.CORS.AllowOrigins, ",") {
		if trimmed := strings.TrimSpace(o); trimmed != "" {
			allowedOrigins[trimmed] = struct{}{}
		}
	}

	wildcardAll := cfg.CORS.AllowOrigins == "*"

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		c.Writer.Header().Add("Vary", "Origin")

		if wildcardAll && !cfg.CORS.AllowCredentials {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		} else if _, ok := allowedOrigins[origin]; ok {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}

		c.Writer.Header().Set("Access-Control-Allow-Methods", cfg.CORS.AllowMethods)
		c.Writer.Header().Set("Access-Control-Allow-Headers", cfg.CORS.AllowHeaders)
		c.Writer.Header().
			Set("Access-Control-Allow-Credentials", strconv.FormatBool(cfg.CORS.AllowCredentials))

		if c.Request.Method == http.MethodOptions {
			c.Writer.Header().Set("Access-Control-Max-Age", "86400")
			c.AbortWithStatus(http.StatusNoContent)

			return
		}

		c.Next()
	}
}

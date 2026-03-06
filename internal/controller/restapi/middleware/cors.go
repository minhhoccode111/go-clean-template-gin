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
func CORS(cfg config.CORS) gin.HandlerFunc {
	// Pre-compute the set of allowed origins once at startup.
	allowedOrigins := make(map[string]struct{})

	for o := range strings.SplitSeq(cfg.AllowOrigins, ",") {
		if trimmed := strings.TrimSpace(o); trimmed != "" {
			allowedOrigins[trimmed] = struct{}{}
		}
	}

	wildcardAll := cfg.AllowOrigins == "*"

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Vary must be set on every response so caches key on origin.
		c.Writer.Header().Add("Vary", "Origin")

		if wildcardAll && !cfg.AllowCredentials {
			// Wildcard is only valid without credentials.
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		} else if _, ok := allowedOrigins[origin]; ok {
			// Reflect the matched origin — required for credentials and
			// the only correct behavior for a multi-origin allow-list.
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}
		// If the origin is not in the list we simply omit the header,
		// which causes the browser to block the cross-origin request.

		c.Writer.Header().Set("Access-Control-Allow-Methods", cfg.AllowMethods)
		c.Writer.Header().Set("Access-Control-Allow-Headers", cfg.AllowHeaders)
		c.Writer.Header().
			Set("Access-Control-Allow-Credentials", strconv.FormatBool(cfg.AllowCredentials))

		if c.Request.Method == http.MethodOptions {
			c.Writer.Header().Set("Access-Control-Max-Age", "86400")
			c.AbortWithStatus(http.StatusNoContent)

			return
		}

		c.Next()
	}
}

package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/minhhoccode111/go-clean-template-gin/config"
	"github.com/minhhoccode111/go-clean-template-gin/internal/controller/restapi/v1/response"
	"golang.org/x/time/rate"
)

// RateLimit creates a Gin middleware that enforces rate limiting.
func RateLimit(cfg config.RateLimit) gin.HandlerFunc {
	// Initialize a global rate limiter.
	limiter := rate.NewLimiter(rate.Limit(cfg.RequestsPerSecond), cfg.Burst)

	return func(c *gin.Context) {
		// Attempt to take a token from the bucket.
		if !limiter.Allow() {
			// If not allowed, abort with 429 Too Many Requests
			c.AbortWithStatusJSON(http.StatusTooManyRequests, response.Error{
				Error: "Too many requests. Please try again later.",
			})

			return
		}

		// Proceed to the next handler if allowed.
		c.Next()
	}
}

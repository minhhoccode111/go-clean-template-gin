package middleware

import (
	"math/rand/v2"
	"time"

	"github.com/gin-gonic/gin"
)

// Sleep adds a random latency between minMs and maxMs milliseconds to each request.
//
//nolint:gosec // math/rand is fine for non-security-sensitive latency jitter
func Sleep(minMs, maxMs int) gin.HandlerFunc {
	return func(c *gin.Context) {
		time.Sleep(time.Duration(minMs+rand.IntN(maxMs-minMs+1)) * time.Millisecond)
		c.Next()
	}
}

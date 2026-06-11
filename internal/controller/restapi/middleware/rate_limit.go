package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/minhhoccode111/go-clean-template-gin/config"
	pkgcache "github.com/minhhoccode111/go-clean-template-gin/pkg/cache"
	"golang.org/x/time/rate"
)

// RateLimit creates a Gin middleware that enforces per-IP rate limiting.
func RateLimit(cfg *config.Config) gin.HandlerFunc {
	c, err := pkgcache.New[string, *rate.Limiter](
		pkgcache.MaxCost(cfg.RateLimit.MaxCost),
		pkgcache.TTL(cfg.RateLimit.TTL),
	)
	if err != nil {
		panic(fmt.Sprintf("middleware: RateLimit: failed to create cache: %v", err))
	}

	rps := rate.Limit(cfg.RateLimit.RequestsPerSecond)
	burst := cfg.RateLimit.Burst

	return func(ctx *gin.Context) {
		limiter, err := c.GetOrLoad(ctx.Request.Context(), ctx.ClientIP(),
			func(_ context.Context, _ string) (*rate.Limiter, error) {
				return rate.NewLimiter(rps, burst), nil
			},
		)
		if err != nil || !limiter.Allow() {
			messageResponse(
				ctx,
				http.StatusTooManyRequests,
				"Too many requests. Please try again later.",
			)

			return
		}

		ctx.Next()
	}
}

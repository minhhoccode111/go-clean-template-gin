package v1

import (
	v1 "github.com/minhhoccode111/go-clean-template-gin/internal/controller/nats_rpc/v1"
	"github.com/minhhoccode111/go-clean-template-gin/internal/usecase"
	"github.com/minhhoccode111/go-clean-template-gin/pkg/logger"
	"github.com/minhhoccode111/go-clean-template-gin/pkg/nats/nats_rpc/server"
)

// NewRouter -.
func NewRouter(t usecase.Translation, l logger.Interface) map[string]server.CallHandler {
	routes := make(map[string]server.CallHandler)

	{
		v1.NewTranslationRoutes(routes, t, l)
	}

	return routes
}

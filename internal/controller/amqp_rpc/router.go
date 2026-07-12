// Package amqprpc defines the RabbitMQ RPC router.
package amqprpc

import (
	v1 "github.com/minhhoccode111/go-clean-template-gin/internal/controller/amqp_rpc/v1"
	"github.com/minhhoccode111/go-clean-template-gin/internal/usecase"
	"github.com/minhhoccode111/go-clean-template-gin/pkg/jwt"
	"github.com/minhhoccode111/go-clean-template-gin/pkg/logger"
	"github.com/minhhoccode111/go-clean-template-gin/pkg/rabbitmq/rmq_rpc/server"
)

// NewRouter -.
func NewRouter(t usecase.Translation, u usecase.User, tk usecase.Task, jwtManager *jwt.Manager, l logger.Interface) map[string]server.CallHandler {
	routes := make(map[string]server.CallHandler)

	{
		v1.NewRoutes(routes, t, u, tk, jwtManager, l)
	}

	return routes
}

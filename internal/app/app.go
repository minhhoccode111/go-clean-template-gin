// Package app configures and runs application.
package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/minhhoccode111/go-clean-template-gin/config"
	amqprpc "github.com/minhhoccode111/go-clean-template-gin/internal/controller/amqp_rpc"
	"github.com/minhhoccode111/go-clean-template-gin/internal/controller/grpc"
	natsrpc "github.com/minhhoccode111/go-clean-template-gin/internal/controller/nats_rpc"
	"github.com/minhhoccode111/go-clean-template-gin/internal/controller/restapi"
	"github.com/minhhoccode111/go-clean-template-gin/internal/entity"
	repocache "github.com/minhhoccode111/go-clean-template-gin/internal/repo/cache"
	"github.com/minhhoccode111/go-clean-template-gin/internal/repo/persistent"
	"github.com/minhhoccode111/go-clean-template-gin/internal/repo/webapi"
	"github.com/minhhoccode111/go-clean-template-gin/internal/usecase/translation"
	"github.com/minhhoccode111/go-clean-template-gin/pkg/cache"
	"github.com/minhhoccode111/go-clean-template-gin/pkg/grpcserver"
	"github.com/minhhoccode111/go-clean-template-gin/pkg/httpserver"
	"github.com/minhhoccode111/go-clean-template-gin/pkg/logger"
	natsRPCServer "github.com/minhhoccode111/go-clean-template-gin/pkg/nats/nats_rpc/server"
	"github.com/minhhoccode111/go-clean-template-gin/pkg/postgres"
	rmqRPCServer "github.com/minhhoccode111/go-clean-template-gin/pkg/rabbitmq/rmq_rpc/server"
	"github.com/minhhoccode111/go-clean-template-gin/pkg/validatorx"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) { //nolint: gocyclo,cyclop,funlen,gocritic,nolintlint
	l := logger.New(cfg.Log.Level)
	v := validatorx.New()

	// Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	// Cache
	otterCache, err := cache.New[string, []entity.Translation](
		cache.MaxCost(cfg.Cache.MaxCost),
		cache.TTL(cfg.Cache.TTL),
	)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - ottercache.New: %w", err))
	}

	// Use-Case
	translationUseCase := translation.New(
		persistent.New(pg),
		webapi.New(),
		repocache.New(otterCache),
	)

	// RabbitMQ RPC Server
	rmqRouter := amqprpc.NewRouter(translationUseCase, l)

	rmqServer, err := rmqRPCServer.New(cfg.RMQ.URL, cfg.RMQ.ServerExchange, rmqRouter, l)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - rmqServer - server.New: %w", err))
	}

	// NATS RPC Server
	natsRouter := natsrpc.NewRouter(translationUseCase, l)

	natsServer, err := natsRPCServer.New(cfg.NATS.URL, cfg.NATS.ServerExchange, natsRouter, l)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - natsServer - server.New: %w", err))
	}

	// gRPC Server
	grpcServer := grpcserver.New(l, grpcserver.Port(cfg.GRPC.Port))
	grpc.NewRouter(grpcServer.App, translationUseCase, l)

	// HTTP Server
	httpServer := httpserver.New(l, httpserver.Port(cfg.HTTP.Port))
	restapi.NewRouter(httpServer.Engine, cfg, translationUseCase, l, v)

	// Start servers
	rmqServer.Start()
	natsServer.Start()
	grpcServer.Start()
	httpServer.Start()

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: %s", s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	case err = <-grpcServer.Notify():
		l.Error(fmt.Errorf("app - Run - grpcServer.Notify: %w", err))
	case err = <-rmqServer.Notify():
		l.Error(fmt.Errorf("app - Run - rmqServer.Notify: %w", err))
	case err = <-natsServer.Notify():
		l.Error(fmt.Errorf("app - Run - natsServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}

	err = grpcServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - grpcServer.Shutdown: %w", err))
	}

	err = rmqServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - rmqServer.Shutdown: %w", err))
	}

	err = natsServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - natsServer.Shutdown: %w", err))
	}
}

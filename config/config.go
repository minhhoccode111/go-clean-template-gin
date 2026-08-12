package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

type (
	// Config -.
	Config struct {
		App       App
		HTTP      HTTP
		Log       Log
		PG        PG
		GRPC      GRPC
		RMQ       RMQ
		NATS      NATS
		Metrics   Metrics
		Swagger   Swagger
		CORS      CORS
		Cache     Cache
		JWT       JWT
		RateLimit RateLimit
		Tracing   Tracing
	}

	// App -.
	App struct {
		Name    string `env:"APP_NAME,required"`
		Version string `env:"APP_VERSION,required"`
	}

	// HTTP -.
	// Port, MaxMultipartMemory (max memory for multipart form parsing, default 32MB),
	// MaxBodySize (max request body size for all endpoints, default 10MB).
	HTTP struct {
		Port               string `env:"HTTP_PORT,required"`
		MaxMultipartMemory int64  `env:"HTTP_MAX_MULTIPART_MEMORY" envDefault:"33554432"`
		MaxBodySize        int64  `env:"HTTP_MAX_BODY_SIZE"        envDefault:"10485760"`
	}

	// Log -.
	Log struct {
		Level string `env:"LOG_LEVEL,required"`
	}

	// PG -.
	PG struct {
		PoolMax int    `env:"PG_POOL_MAX,required"`
		URL     string `env:"PG_URL,required"`
	}

	// GRPC -.
	GRPC struct {
		Port string `env:"GRPC_PORT,required"`
	}

	// RMQ -.
	RMQ struct {
		ServerExchange string `env:"RMQ_RPC_SERVER,required"`
		ClientExchange string `env:"RMQ_RPC_CLIENT,required"`
		URL            string `env:"RMQ_URL,required"`
	}

	// NATS -.
	NATS struct {
		ServerExchange string `env:"NATS_RPC_SERVER,required"`
		URL            string `env:"NATS_URL,required"`
	}

	// Metrics -.
	Metrics struct {
		Enabled bool `env:"METRICS_ENABLED" envDefault:"true"`
	}

	// Swagger -.
	Swagger struct {
		Enabled bool `env:"SWAGGER_ENABLED" envDefault:"false"`
	}

	// CORS -.
	CORS struct {
		AllowOrigins     string `env:"CORS_ALLOW_ORIGINS,required"`
		AllowCredentials bool   `env:"CORS_ALLOW_CREDENTIALS,required"`
		AllowHeaders     string `env:"CORS_ALLOW_HEADERS,required"`
		AllowMethods     string `env:"CORS_ALLOW_METHODS,required"`
	}

	// Cache -.
	Cache struct {
		MaxCost int           `env:"CACHE_MAX_COST" envDefault:"10000"`
		TTL     time.Duration `env:"CACHE_TTL"      envDefault:"5m"`
	}

	// JWT -.
	JWT struct {
		Secret string        `env:"JWT_SECRET,required"`
		TTL    time.Duration `env:"JWT_TTL"             envDefault:"24h"`
	}

	// RateLimit -.
	RateLimit struct {
		RequestsPerSecond float64       `env:"RATE_LIMIT_RPS"       envDefault:"10"`
		Burst             int           `env:"RATE_LIMIT_BURST"     envDefault:"20"`
		MaxCost           int           `env:"RATE_LIMIT_CACHE_MAX" envDefault:"10000"`
		TTL               time.Duration `env:"RATE_LIMIT_CACHE_TTL" envDefault:"10m"`
	}

	// Tracing -.
	Tracing struct {
		Enabled      bool    `env:"TRACING_ENABLED" envDefault:"false"`
		OTLPEndpoint string  `env:"TRACING_OTLP_ENDPOINT" envDefault:"localhost:4317"`
		OTLPInsecure bool    `env:"TRACING_OTLP_INSECURE" envDefault:"true"`
		SampleRate   float64 `env:"TRACING_SAMPLE_RATE" envDefault:"0.1"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}

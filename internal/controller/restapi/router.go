package restapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/minhhoccode111/go-clean-template-gin/config"
	_ "github.com/minhhoccode111/go-clean-template-gin/docs" // Swagger docs.
	"github.com/minhhoccode111/go-clean-template-gin/internal/controller/restapi/middleware"
	v1 "github.com/minhhoccode111/go-clean-template-gin/internal/controller/restapi/v1"
	"github.com/minhhoccode111/go-clean-template-gin/internal/usecase"
	"github.com/minhhoccode111/go-clean-template-gin/pkg/logger"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

// NewRouter -.
// Swagger spec:
// @title       Go Clean Template API
// @description Using a translation service as an example
// @version     1.0
// @host        localhost:8080
// @BasePath    /v1
func NewRouter(
	handler *gin.Engine,
	cfg *config.Config,
	t usecase.Translation,
	l logger.Interface,
	v *validator.Validate,
) {
	// Options
	handler.Use(middleware.Logger(l))
	handler.Use(middleware.Recovery(l))
	handler.Use(middleware.CORS(cfg.CORS))
	handler.Use(middleware.RateLimit(cfg.RateLimit))

	// Prometheus metrics
	if cfg.Metrics.Enabled {
		p := ginprometheus.NewPrometheus("go-clean-template-gin-gin")
		p.Use(handler)
	}

	// Swagger
	if cfg.Swagger.Enabled {
		handler.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// K8s probe
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Routers
	apiV1Group := handler.Group("/v1")
	{
		v1.NewTranslationRoutes(apiV1Group, t, l, v)
	}
}

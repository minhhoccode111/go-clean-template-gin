package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/minhhoccode111/go-clean-template-gin/config"
	"github.com/minhhoccode111/go-clean-template-gin/internal/usecase"
	"github.com/minhhoccode111/go-clean-template-gin/pkg/logger"
)

// NewTranslationRoutes -.
func NewTranslationRoutes(
	apiV1Group *gin.RouterGroup,
	cfg *config.Config,
	t usecase.Translation,
	l logger.Interface,
	v *validator.Validate,
) {
	r := NewV1(cfg, t, l, v)

	translationGroup := apiV1Group.Group("/translation")

	{
		translationGroup.GET("/history", r.history)
		translationGroup.POST("/do-translate", r.doTranslate)
	}
}

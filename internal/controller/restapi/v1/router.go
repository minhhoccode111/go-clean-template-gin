package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/minhhoccode111/go-clean-template-gin/internal/usecase"
	"github.com/minhhoccode111/go-clean-template-gin/pkg/logger"
)

// NewTranslationRoutes -.
func NewTranslationRoutes(
	apiV1Group *gin.RouterGroup,
	t usecase.Translation,
	l logger.Interface,
	v *validator.Validate,
) {
	r := &V1{t: t, l: l, v: v}

	translationGroup := apiV1Group.Group("/translation")

	{
		translationGroup.GET("/history", r.history)
		translationGroup.POST("/do-translate", r.doTranslate)
	}
}

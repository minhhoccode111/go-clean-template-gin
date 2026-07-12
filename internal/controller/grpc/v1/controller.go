package v1

import (
	"github.com/go-playground/validator/v10"
	v1 "github.com/minhhoccode111/go-clean-template-gin/docs/proto/v1"
	"github.com/minhhoccode111/go-clean-template-gin/internal/usecase"
	"github.com/minhhoccode111/go-clean-template-gin/pkg/logger"
)

// TranslationController -.
type TranslationController struct {
	v1.UnimplementedTranslationServer

	t usecase.Translation
	l logger.Interface
	v *validator.Validate
}

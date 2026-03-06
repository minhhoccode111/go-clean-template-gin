package v1

import (
	"github.com/go-playground/validator/v10"
	"github.com/minhhoccode111/go-clean-template-gin/internal/usecase"
	"github.com/minhhoccode111/go-clean-template-gin/pkg/logger"
)

// V1 -.
type V1 struct {
	t usecase.Translation
	l logger.Interface
	v *validator.Validate
}

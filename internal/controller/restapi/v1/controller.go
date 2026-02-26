package v1

import (
	"github.com/minhhoccode111/go-clean-template-gin/internal/usecase"
	"github.com/minhhoccode111/go-clean-template-gin/pkg/logger"
	"github.com/go-playground/validator/v10"
)

// V1 -.
type V1 struct {
	t usecase.Translation
	l logger.Interface
	v *validator.Validate
}

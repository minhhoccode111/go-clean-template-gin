package v1

import (
	"github.com/go-playground/validator/v10"
	"github.com/minhhoccode111/go-clean-template-gin/config"
	"github.com/minhhoccode111/go-clean-template-gin/internal/usecase"
	"github.com/minhhoccode111/go-clean-template-gin/pkg/logger"
)

// V1 -.
type V1 struct {
	cfg *config.Config
	t   usecase.Translation
	l   logger.Interface
	v   *validator.Validate
}

func NewV1(cfg *config.Config, t usecase.Translation, l logger.Interface, v *validator.Validate) *V1 {
	return &V1{
		cfg: cfg,
		t:   t,
		l:   l,
		v:   v,
	}
}

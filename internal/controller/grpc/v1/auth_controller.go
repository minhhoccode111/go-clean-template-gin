package v1

import (
	v1 "github.com/minhhoccode111/go-clean-template-gin/docs/proto/v1"
	"github.com/minhhoccode111/go-clean-template-gin/internal/usecase"
	"github.com/minhhoccode111/go-clean-template-gin/pkg/logger"
	"github.com/go-playground/validator/v10"
)

// AuthController -.
type AuthController struct {
	v1.UnimplementedAuthServiceServer

	u usecase.User
	l logger.Interface
	v *validator.Validate
}

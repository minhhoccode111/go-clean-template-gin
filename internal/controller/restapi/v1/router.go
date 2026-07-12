package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/minhhoccode111/go-clean-template-gin/config"
	"github.com/minhhoccode111/go-clean-template-gin/internal/controller/restapi/middleware"
	"github.com/minhhoccode111/go-clean-template-gin/internal/usecase"
	"github.com/minhhoccode111/go-clean-template-gin/pkg/jwt"
	"github.com/minhhoccode111/go-clean-template-gin/pkg/logger"
)

// NewRoutes -.
func NewRoutes(
	apiV1Group *gin.RouterGroup,
	cfg *config.Config,
	t usecase.Translation,
	u usecase.User,
	tk usecase.Task,
	jwtManager *jwt.Manager,
	l logger.Interface,
	v *validator.Validate,
) {
	r := NewV1(cfg, t, u, tk, l, v)

	// Public routes
	authGroup := apiV1Group.Group("/auth")
	{
		authGroup.POST("/register", r.register)
		authGroup.POST("/login", r.login)
	}

	// Protected routes
	protected := apiV1Group.Group("", middleware.Auth(jwtManager))

	userGroup := protected.Group("/user")
	{
		userGroup.GET("/profile", r.profile)
	}

	taskGroup := protected.Group("/tasks")
	{
		taskGroup.POST("/", r.createTask)
		taskGroup.GET("/", r.listTasks)
		taskGroup.GET("/:id", r.getTask)
		taskGroup.PUT("/:id", r.updateTask)
		taskGroup.PATCH("/:id/status", r.transitionTask)
		taskGroup.DELETE("/:id", r.deleteTask)
	}

	translationGroup := protected.Group("/translation")
	{
		translationGroup.GET("/history", r.history)
		translationGroup.POST("/do-translate", r.doTranslate)
	}
}

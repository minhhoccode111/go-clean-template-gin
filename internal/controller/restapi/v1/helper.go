package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/minhhoccode111/go-clean-template-gin/internal/controller/restapi/middleware"
	"github.com/minhhoccode111/go-clean-template-gin/internal/controller/restapi/v1/response"
)

//nolint:unused // template utility, used when auth helpers are wired
func messageResponse(c *gin.Context, code int, msg string) {
	c.JSON(code, response.Message{Message: msg})
}

// extractActorID extracts the authenticated user ID from gin context.
//
//nolint:unused // template utility, wired when auth middleware is used
func extractActorID(c *gin.Context) (int64, bool) {
	userIDRaw, exists := c.Get(string(middleware.CtxUserIDKey))
	if !exists {
		errorResponse(c, http.StatusUnauthorized, "unauthorized")

		return 0, false
	}

	id, ok := userIDRaw.(int64)
	if !ok {
		errorResponse(c, http.StatusUnauthorized, "invalid user ID")

		return 0, false
	}

	return id, true
}

//nolint:unused // template utility, used when extending with route handlers
func parseParamInt64(c *gin.Context, paramName, errMsg string) (int64, bool) {
	id, err := strconv.ParseInt(c.Param(paramName), 10, 64)
	if err != nil || id <= 0 {
		errorResponse(c, http.StatusBadRequest, errMsg)

		return 0, false
	}

	return id, true
}

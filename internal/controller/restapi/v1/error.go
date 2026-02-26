package v1

import (
	"github.com/minhhoccode111/go-clean-template-gin/internal/controller/restapi/v1/response"
	"github.com/gin-gonic/gin"
)

func errorResponse(c *gin.Context, code int, msg string) {
	c.AbortWithStatusJSON(code, response.Error{Error: msg})
}

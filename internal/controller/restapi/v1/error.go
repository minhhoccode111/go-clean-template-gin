package v1

import (
	"github.com/evrone/go-clean-template/internal/controller/restapi/v1/response"
	"github.com/gin-gonic/gin"
)

func errorResponse(c *gin.Context, code int, msg string) {
	c.AbortWithStatusJSON(code, response.Error{Error: msg})
}

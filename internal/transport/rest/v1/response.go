package v1

import (
	"strings"

	"github.com/gin-gonic/gin"
)

type response struct {
	Message string `json:"message"`
}

func newResponse(c *gin.Context, statusCode int, message string) {
	c.AbortWithStatusJSON(statusCode, response{Message: message})
}

func getPrefixFromParam(c *gin.Context, prefix string) string {
	param := c.Param(prefix)
	return strings.Replace(param, "_", " ", -1)
}

package httputil

import (
	"github.com/gin-gonic/gin"
	"github.com/primosec/pulse/pkg/apperror"
)

type Response struct {
	Data  any    `json:"data,omitempty"`
	Error string `json:"error,omitempty"`
}

func JSON(c *gin.Context, status int, data any) {
	c.JSON(status, Response{Data: data})
}

func Error(c *gin.Context, err error) {
	if appErr, ok := err.(*apperror.AppError); ok {
		c.JSON(appErr.Code, Response{Error: appErr.Message})
		return
	}
	c.JSON(500, Response{Error: "internal server error"})
}

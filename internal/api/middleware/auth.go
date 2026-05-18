package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/primosec/pulse/pkg/apperror"
	"github.com/primosec/pulse/pkg/httputil"
)

type TokenValidator interface {
	ValidateToken(tokenStr string) (string, error)
}

func Auth(validator TokenValidator) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			httputil.Error(c, apperror.ErrUnauthorized)
			c.Abort()
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			httputil.Error(c, apperror.ErrUnauthorized)
			c.Abort()
			return
		}

		userID, err := validator.ValidateToken(parts[1])
		if err != nil {
			httputil.Error(c, apperror.ErrUnauthorized)
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}

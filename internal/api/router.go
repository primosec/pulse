package api

import (
	"github.com/gin-gonic/gin"
	"github.com/primosec/pulse/internal/api/handler"
	"github.com/primosec/pulse/internal/service"
)

func NewRouter(authService *service.AuthService) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	authHandler := handler.NewAuthHandler(authService)

	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}
	}

	return r
}

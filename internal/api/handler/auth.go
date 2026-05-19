package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/primosec/pulse/internal/dto"
	"github.com/primosec/pulse/internal/service"
	"github.com/primosec/pulse/pkg/apperror"
	"github.com/primosec/pulse/pkg/httputil"
)

type AuthHandler struct {
	auth *service.AuthService
}

func NewAuthHandler(auth *service.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.Error(c, apperror.New(http.StatusBadRequest, err.Error()))
		return
	}

	user, err := h.auth.Register(c.Request.Context(), req.Email, req.Password, req.Name)
	if err != nil {
		httputil.Error(c, err)
		return
	}

	httputil.JSON(c, http.StatusCreated, dto.UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBind(&req); err != nil {
		httputil.Error(c, apperror.New(http.StatusBadRequest, err.Error()))
		return
	}

	token, err := h.auth.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		httputil.Error(c, err)
		return
	}

	httputil.JSON(c, http.StatusOK, dto.TokenResponse{Token: token})
}

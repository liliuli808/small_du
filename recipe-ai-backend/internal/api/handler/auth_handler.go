package handler

import (
	"net/http"
	"recipe-ai-backend/internal/model"
	"recipe-ai-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) WxLogin(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Code:    "INVALID_PARAMS",
			Message: "缺少 code 参数",
		})
		return
	}

	resp, err := h.authService.WxLogin(c.Request.Context(), req.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Code:    "LOGIN_FAILED",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

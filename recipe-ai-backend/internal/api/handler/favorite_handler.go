package handler

import (
	"net/http"
	"recipe-ai-backend/internal/model"
	"recipe-ai-backend/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// FavoriteHandler 收藏Handler
type FavoriteHandler struct {
	favoriteService *service.FavoriteService
}

// NewFavoriteHandler 创建收藏Handler
func NewFavoriteHandler(favoriteService *service.FavoriteService) *FavoriteHandler {
	return &FavoriteHandler{
		favoriteService: favoriteService,
	}
}

// ToggleFavorite 切换收藏状态
func (h *FavoriteHandler) ToggleFavorite(c *gin.Context) {
	userOpenID := c.GetString("user_openid")
	if userOpenID == "" {
		c.JSON(http.StatusUnauthorized, model.APIResponse{
			Code:    "UNAUTHORIZED",
			Message: "请先登录",
		})
		return
	}

	recipeIDStr := c.Param("recipe_id")
	recipeID, err := strconv.ParseInt(recipeIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Code:    "INVALID_RECIPE_ID",
			Message: "菜谱ID格式错误",
		})
		return
	}

	isFavorited, err := h.favoriteService.ToggleFavorite(c.Request.Context(), userOpenID, recipeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Code:    "TOGGLE_FAILED",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.FavoriteToggleResponse{IsFavorited: isFavorited})
}

// GetFavoriteStatus 获取收藏状态
func (h *FavoriteHandler) GetFavoriteStatus(c *gin.Context) {
	userOpenID := c.GetString("user_openid")
	if userOpenID == "" {
		c.JSON(http.StatusOK, model.FavoriteToggleResponse{IsFavorited: false})
		return
	}

	recipeIDStr := c.Param("recipe_id")
	recipeID, err := strconv.ParseInt(recipeIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Code:    "INVALID_RECIPE_ID",
			Message: "菜谱ID格式错误",
		})
		return
	}

	isFavorited, err := h.favoriteService.IsFavorited(c.Request.Context(), userOpenID, recipeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Code:    "CHECK_FAILED",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.FavoriteToggleResponse{IsFavorited: isFavorited})
}

// ListFavorites 列出用户收藏
func (h *FavoriteHandler) ListFavorites(c *gin.Context) {
	userOpenID := c.GetString("user_openid")
	if userOpenID == "" {
		c.JSON(http.StatusUnauthorized, model.APIResponse{
			Code:    "UNAUTHORIZED",
			Message: "请先登录",
		})
		return
	}

	limit := 20
	offset := 0
	if l, err := strconv.Atoi(c.Query("limit")); err == nil && l > 0 {
		limit = l
	}
	if o, err := strconv.Atoi(c.Query("offset")); err == nil && o >= 0 {
		offset = o
	}

	resp, err := h.favoriteService.ListUserFavorites(c.Request.Context(), userOpenID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Code:    "LIST_FAILED",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

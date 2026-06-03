package handler

import (
	"net/http"
	"recipe-ai-backend/internal/model"
	"recipe-ai-backend/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserRecipeHandler 用户菜谱Handler
type UserRecipeHandler struct {
	userRecipeService *service.UserRecipeService
}

// NewUserRecipeHandler 创建用户菜谱Handler
func NewUserRecipeHandler(userRecipeService *service.UserRecipeService) *UserRecipeHandler {
	return &UserRecipeHandler{
		userRecipeService: userRecipeService,
	}
}

// CreateUserRecipe 创建用户菜谱
func (h *UserRecipeHandler) CreateUserRecipe(c *gin.Context) {
	userOpenID := c.GetString("user_openid")
	if userOpenID == "" {
		c.JSON(http.StatusUnauthorized, model.APIResponse{
			Code:    "UNAUTHORIZED",
			Message: "请先登录",
		})
		return
	}

	var req model.CreateUserRecipeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Code:    "INVALID_PARAMS",
			Message: "请求参数错误",
		})
		return
	}

	recipe, err := h.userRecipeService.Create(c.Request.Context(), userOpenID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Code:    "CREATE_FAILED",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": recipe.ID})
}

// GetUserRecipe 获取用户菜谱详情
func (h *UserRecipeHandler) GetUserRecipe(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Code:    "INVALID_ID",
			Message: "ID格式错误",
		})
		return
	}

	// 获取用户 openid（允许未登录用户返回空字符串，由 service 层校验归属）
	userOpenID := c.GetString("user_openid")

	resp, err := h.userRecipeService.GetByID(c.Request.Context(), userOpenID, id)
	if err != nil {
		c.JSON(http.StatusNotFound, model.APIResponse{
			Code:    "NOT_FOUND",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListUserRecipes 列出用户菜谱
func (h *UserRecipeHandler) ListUserRecipes(c *gin.Context) {
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

	resp, err := h.userRecipeService.ListByUser(c.Request.Context(), userOpenID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Code:    "LIST_FAILED",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateUserRecipe 更新用户菜谱
func (h *UserRecipeHandler) UpdateUserRecipe(c *gin.Context) {
	userOpenID := c.GetString("user_openid")
	if userOpenID == "" {
		c.JSON(http.StatusUnauthorized, model.APIResponse{
			Code:    "UNAUTHORIZED",
			Message: "请先登录",
		})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Code:    "INVALID_ID",
			Message: "ID格式错误",
		})
		return
	}

	var req model.UpdateUserRecipeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Code:    "INVALID_PARAMS",
			Message: "请求参数错误",
		})
		return
	}

	if err := h.userRecipeService.Update(c.Request.Context(), userOpenID, id, &req); err != nil {
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Code:    "UPDATE_FAILED",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// DeleteUserRecipe 删除用户菜谱
func (h *UserRecipeHandler) DeleteUserRecipe(c *gin.Context) {
	userOpenID := c.GetString("user_openid")
	if userOpenID == "" {
		c.JSON(http.StatusUnauthorized, model.APIResponse{
			Code:    "UNAUTHORIZED",
			Message: "请先登录",
		})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Code:    "INVALID_ID",
			Message: "ID格式错误",
		})
		return
	}

	if err := h.userRecipeService.Delete(c.Request.Context(), userOpenID, id); err != nil {
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Code:    "DELETE_FAILED",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

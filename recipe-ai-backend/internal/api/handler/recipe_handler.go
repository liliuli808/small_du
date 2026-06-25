package handler

import (
	"context"
	"net/http"
	"recipe-ai-backend/internal/model"
	"recipe-ai-backend/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// RecipeHandler 菜谱Handler
type RecipeHandler struct {
	recipeService    *service.RecipeService
	nutritionService *service.NutritionService
}

// NewRecipeHandler 创建菜谱Handler
func NewRecipeHandler(recipeService *service.RecipeService, nutritionService *service.NutritionService) *RecipeHandler {
	return &RecipeHandler{
		recipeService:    recipeService,
		nutritionService: nutritionService,
	}
}

// GetRecipe 获取菜谱结果
func (h *RecipeHandler) GetRecipe(c *gin.Context) {
	recipeIDStr := c.Param("recipe_id")
	recipeID, err := strconv.ParseInt(recipeIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Code:    "INVALID_RECIPE_ID",
			Message: "菜谱ID格式错误",
		})
		return
	}

	// 异步增加浏览量（不阻塞响应，使用独立 context）
	go h.recipeService.IncrementViewCount(context.Background(), recipeID)

	resp, err := h.recipeService.GetRecipeResponse(c.Request.Context(), recipeID)
	if err != nil {
		c.JSON(http.StatusNotFound, model.APIResponse{
			Code:    "RECIPE_NOT_FOUND",
			Message: "菜谱不存在",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// SearchRecipes 搜索菜谱
func (h *RecipeHandler) SearchRecipes(c *gin.Context) {
	keyword := c.Query("q")
	if keyword == "" {
		c.JSON(http.StatusOK, model.PopularRecipesResponse{Recipes: []model.PopularRecipeItem{}})
		return
	}

	limit := 20
	offset := 0
	if l, err := strconv.Atoi(c.Query("limit")); err == nil && l > 0 && l <= 50 {
		limit = l
	}
	if o, err := strconv.Atoi(c.Query("offset")); err == nil && o >= 0 {
		offset = o
	}

	resp, err := h.recipeService.SearchRecipes(c.Request.Context(), keyword, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Code:    "SEARCH_FAILED",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListPopularRecipes 获取热门菜谱
func (h *RecipeHandler) ListPopularRecipes(c *gin.Context) {
	limit := 10
	if l, err := strconv.Atoi(c.Query("limit")); err == nil && l > 0 && l <= 50 {
		limit = l
	}

	resp, err := h.recipeService.ListPopularRecipes(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Code:    "LIST_FAILED",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Recalculate 重新计算热量
func (h *RecipeHandler) Recalculate(c *gin.Context) {
	recipeIDStr := c.Param("recipe_id")
	_, err := strconv.ParseInt(recipeIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Code:    "INVALID_RECIPE_ID",
			Message: "菜谱ID格式错误",
		})
		return
	}

	var req model.RecalculateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Code:    "INVALID_PARAMS",
			Message: "请求参数错误",
		})
		return
	}

	// 转换请求参数
	adjustments := make([]model.IngredientAdjustment, len(req.Ingredients))
	for i, ing := range req.Ingredients {
		adjustments[i] = model.IngredientAdjustment{
			Name:  ing.Name,
			Grams: ing.Grams,
		}
	}

	nutrition, err := h.nutritionService.Recalculate(c.Request.Context(), req.Servings, adjustments)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Code:    "NUTRITION_FAILED",
			Message: "热量计算失败",
		})
		return
	}

	c.JSON(http.StatusOK, model.RecalculateResponse{
		Nutrition: *nutrition,
	})
}

// DeriveAsUserRecipe 将AI菜谱转为用户可编辑的格式
func (h *RecipeHandler) DeriveAsUserRecipe(c *gin.Context) {
	recipeIDStr := c.Param("recipe_id")
	recipeID, err := strconv.ParseInt(recipeIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Code:    "INVALID_RECIPE_ID",
			Message: "菜谱ID格式错误",
		})
		return
	}

	data, err := h.recipeService.DeriveAsUserRecipeData(c.Request.Context(), recipeID)
	if err != nil {
		c.JSON(http.StatusNotFound, model.APIResponse{
			Code:    "NOT_FOUND",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, data)
}

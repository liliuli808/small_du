package service

import (
	"context"
	"encoding/json"
	"fmt"
	"recipe-ai-backend/internal/model"
	"recipe-ai-backend/internal/repository"
)

// UserRecipeService 用户菜谱服务
type UserRecipeService struct {
	userRecipeRepo repository.UserRecipeRepository
}

// NewUserRecipeService 创建用户菜谱服务
func NewUserRecipeService(userRecipeRepo repository.UserRecipeRepository) *UserRecipeService {
	return &UserRecipeService{
		userRecipeRepo: userRecipeRepo,
	}
}

// Create 创建用户菜谱
func (s *UserRecipeService) Create(ctx context.Context, userOpenID string, req *model.CreateUserRecipeRequest) (*model.UserRecipe, error) {
	recipeData := model.UserRecipeData{
		DishName:    req.DishName,
		Servings:    req.Servings,
		Ingredients: req.Ingredients,
		Steps:       req.Steps,
		Tips:        req.Tips,
	}

	recipeJSON, err := json.Marshal(recipeData)
	if err != nil {
		return nil, fmt.Errorf("序列化菜谱数据失败: %w", err)
	}

	recipe := &model.UserRecipe{
		UserOpenID: userOpenID,
		DishName:   req.DishName,
		Servings:   req.Servings,
		RecipeJSON: model.JSONB{},
	}
	json.Unmarshal(recipeJSON, &recipe.RecipeJSON)

	if err := s.userRecipeRepo.Create(ctx, recipe); err != nil {
		return nil, fmt.Errorf("创建菜谱失败: %w", err)
	}

	return recipe, nil
}

// GetByID 获取用户菜谱详情
func (s *UserRecipeService) GetByID(ctx context.Context, userOpenID string, id int64) (*model.UserRecipeDetailResponse, error) {
	recipe, err := s.userRecipeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取菜谱失败: %w", err)
	}

	// 校验归属
	if recipe.UserOpenID != userOpenID {
		return nil, fmt.Errorf("无权限访问该菜谱")
	}

	recipeData, err := repository.ParseUserRecipeJSON(recipe)
	if err != nil {
		return nil, fmt.Errorf("解析菜谱数据失败: %w", err)
	}

	return &model.UserRecipeDetailResponse{
		ID:        recipe.ID,
		DishName:  recipe.DishName,
		Servings:  recipe.Servings,
		Recipe:    *recipeData,
		CreatedAt: recipe.CreatedAt.Format("2006-01-02 15:04"),
	}, nil
}

// ListByUser 列出用户菜谱
func (s *UserRecipeService) ListByUser(ctx context.Context, userOpenID string, limit, offset int) (*model.UserRecipesResponse, error) {
	recipes, err := s.userRecipeRepo.ListByUser(ctx, userOpenID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("获取菜谱列表失败: %w", err)
	}

	items := make([]model.UserRecipeListItem, len(recipes))
	for i, r := range recipes {
		items[i] = model.UserRecipeListItem{
			ID:        r.ID,
			DishName:  r.DishName,
			Servings:  r.Servings,
			CreatedAt: r.CreatedAt.Format("2006-01-02"),
		}
	}

	return &model.UserRecipesResponse{Recipes: items}, nil
}

// Update 更新用户菜谱
func (s *UserRecipeService) Update(ctx context.Context, userOpenID string, id int64, req *model.UpdateUserRecipeRequest) error {
	recipeData := model.UserRecipeData{
		DishName:    req.DishName,
		Servings:    req.Servings,
		Ingredients: req.Ingredients,
		Steps:       req.Steps,
		Tips:        req.Tips,
	}

	recipeJSON, err := json.Marshal(recipeData)
	if err != nil {
		return fmt.Errorf("序列化菜谱数据失败: %w", err)
	}

	recipe := &model.UserRecipe{
		ID:         id,
		UserOpenID: userOpenID,
		DishName:   req.DishName,
		Servings:   req.Servings,
		RecipeJSON: model.JSONB{},
	}
	json.Unmarshal(recipeJSON, &recipe.RecipeJSON)

	return s.userRecipeRepo.Update(ctx, recipe)
}

// Delete 删除用户菜谱
func (s *UserRecipeService) Delete(ctx context.Context, userOpenID string, id int64) error {
	return s.userRecipeRepo.Delete(ctx, id, userOpenID)
}

// CountByUser 统计用户菜谱数量
func (s *UserRecipeService) CountByUser(ctx context.Context, userOpenID string) (int64, error) {
	return s.userRecipeRepo.CountByUser(ctx, userOpenID)
}

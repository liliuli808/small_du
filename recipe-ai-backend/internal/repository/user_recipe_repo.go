package repository

import (
	"context"
	"encoding/json"
	"recipe-ai-backend/internal/model"

	"gorm.io/gorm"
)

// UserRecipeRepository 用户菜谱存储接口
type UserRecipeRepository interface {
	Create(ctx context.Context, recipe *model.UserRecipe) error
	GetByID(ctx context.Context, id int64) (*model.UserRecipe, error)
	ListByUser(ctx context.Context, userOpenID string, limit, offset int) ([]model.UserRecipe, error)
	CountByUser(ctx context.Context, userOpenID string) (int64, error)
	Update(ctx context.Context, recipe *model.UserRecipe) error
	Delete(ctx context.Context, id int64, userOpenID string) error
}

// userRecipeRepository 用户菜谱存储实现
type userRecipeRepository struct {
	db *gorm.DB
}

// NewUserRecipeRepository 创建用户菜谱存储
func NewUserRecipeRepository(db *gorm.DB) UserRecipeRepository {
	return &userRecipeRepository{db: db}
}

// Create 创建用户菜谱
func (r *userRecipeRepository) Create(ctx context.Context, recipe *model.UserRecipe) error {
	return r.db.WithContext(ctx).Create(recipe).Error
}

// GetByID 根据ID查询
func (r *userRecipeRepository) GetByID(ctx context.Context, id int64) (*model.UserRecipe, error) {
	var recipe model.UserRecipe
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&recipe).Error
	if err != nil {
		return nil, err
	}
	return &recipe, nil
}

// ListByUser 列出用户的菜谱
func (r *userRecipeRepository) ListByUser(ctx context.Context, userOpenID string, limit, offset int) ([]model.UserRecipe, error) {
	var recipes []model.UserRecipe
	err := r.db.WithContext(ctx).
		Where("user_openid = ?", userOpenID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&recipes).Error
	return recipes, err
}

// CountByUser 统计用户菜谱数量
func (r *userRecipeRepository) CountByUser(ctx context.Context, userOpenID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.UserRecipe{}).
		Where("user_openid = ?", userOpenID).
		Count(&count).Error
	return count, err
}

// Update 更新用户菜谱
func (r *userRecipeRepository) Update(ctx context.Context, recipe *model.UserRecipe) error {
	return r.db.WithContext(ctx).
		Model(&model.UserRecipe{}).
		Where("id = ? AND user_openid = ?", recipe.ID, recipe.UserOpenID).
		Updates(map[string]interface{}{
			"dish_name":   recipe.DishName,
			"servings":    recipe.Servings,
			"recipe_json": recipe.RecipeJSON,
			"updated_at":  gorm.Expr("CURRENT_TIMESTAMP"),
		}).Error
}

// Delete 删除用户菜谱
func (r *userRecipeRepository) Delete(ctx context.Context, id int64, userOpenID string) error {
	return r.db.WithContext(ctx).
		Where("id = ? AND user_openid = ?", id, userOpenID).
		Delete(&model.UserRecipe{}).Error
}

// ParseUserRecipeJSON 解析用户菜谱JSON
func ParseUserRecipeJSON(recipe *model.UserRecipe) (*model.UserRecipeData, error) {
	if recipe == nil {
		return nil, nil
	}
	data, err := json.Marshal(recipe.RecipeJSON)
	if err != nil {
		return nil, err
	}
	var result model.UserRecipeData
	err = json.Unmarshal(data, &result)
	return &result, err
}

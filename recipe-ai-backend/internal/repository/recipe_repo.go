package repository

import (
	"context"
	"encoding/json"
	"recipe-ai-backend/internal/model"

	"gorm.io/gorm"
)

// RecipeRepository 菜谱存储接口
type RecipeRepository interface {
	Create(ctx context.Context, recipe *model.Recipe) error
	GetByID(ctx context.Context, id int64) (*model.Recipe, error)
	GetByVideoID(ctx context.Context, videoID int64) (*model.Recipe, error)
	GetByBVID(ctx context.Context, bvid string) (*model.Recipe, error)
}

// recipeRepository 菜谱存储实现
type recipeRepository struct {
	db *gorm.DB
}

// NewRecipeRepository 创建菜谱存储
func NewRecipeRepository(db *gorm.DB) RecipeRepository {
	return &recipeRepository{db: db}
}

// Create 创建菜谱
func (r *recipeRepository) Create(ctx context.Context, recipe *model.Recipe) error {
	return r.db.WithContext(ctx).Create(recipe).Error
}

// GetByID 根据ID查询
func (r *recipeRepository) GetByID(ctx context.Context, id int64) (*model.Recipe, error) {
	var recipe model.Recipe
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&recipe).Error
	if err != nil {
		return nil, err
	}
	return &recipe, nil
}

// GetByVideoID 根据视频ID查询
func (r *recipeRepository) GetByVideoID(ctx context.Context, videoID int64) (*model.Recipe, error) {
	var recipe model.Recipe
	err := r.db.WithContext(ctx).Where("video_id = ?", videoID).First(&recipe).Error
	if err != nil {
		return nil, err
	}
	return &recipe, nil
}

// GetByBVID 根据BVID查询已有菜谱（通过video表JOIN）
func (r *recipeRepository) GetByBVID(ctx context.Context, bvid string) (*model.Recipe, error) {
	var recipe model.Recipe
	err := r.db.WithContext(ctx).
		Joins("JOIN bilibili_videos ON bilibili_videos.id = recipes.video_id").
		Where("bilibili_videos.bvid = ?", bvid).
		First(&recipe).Error
	if err != nil {
		return nil, err
	}
	return &recipe, nil
}

// ParseRecipeJSON 解析菜谱JSON
func ParseRecipeJSON(recipe *model.Recipe) (*model.RecipeData, error) {
	if recipe == nil {
		return nil, nil
	}
	data, err := json.Marshal(recipe.RecipeJSON)
	if err != nil {
		return nil, err
	}
	var result model.RecipeData
	err = json.Unmarshal(data, &result)
	return &result, err
}

package repository

import (
	"context"
	"encoding/json"
	"recipe-ai-backend/internal/model"

	"gorm.io/gorm"
)

// NutritionRepository 营养存储接口
type NutritionRepository interface {
	CreateResult(ctx context.Context, result *model.NutritionResult) error
	GetResultByRecipeID(ctx context.Context, recipeID int64) (*model.NutritionResult, error)
	GetAllFoods(ctx context.Context) ([]model.NutritionFood, error)
	SearchFoods(ctx context.Context, keyword string) ([]model.NutritionFood, error)
}

// nutritionRepository 营养存储实现
type nutritionRepository struct {
	db *gorm.DB
}

// NewNutritionRepository 创建营养存储
func NewNutritionRepository(db *gorm.DB) NutritionRepository {
	return &nutritionRepository{db: db}
}

// CreateResult 创建营养结果
func (r *nutritionRepository) CreateResult(ctx context.Context, result *model.NutritionResult) error {
	return r.db.WithContext(ctx).Create(result).Error
}

// GetResultByRecipeID 根据菜谱ID查询营养结果
func (r *nutritionRepository) GetResultByRecipeID(ctx context.Context, recipeID int64) (*model.NutritionResult, error) {
	var result model.NutritionResult
	err := r.db.WithContext(ctx).Where("recipe_id = ?", recipeID).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetAllFoods 获取所有食材
func (r *nutritionRepository) GetAllFoods(ctx context.Context) ([]model.NutritionFood, error) {
	var foods []model.NutritionFood
	err := r.db.WithContext(ctx).Find(&foods).Error
	return foods, err
}

// SearchFoods 搜索食材
func (r *nutritionRepository) SearchFoods(ctx context.Context, keyword string) ([]model.NutritionFood, error) {
	var foods []model.NutritionFood
	err := r.db.WithContext(ctx).
		Where("canonical_name LIKE ? OR aliases::text LIKE ?", "%"+keyword+"%", "%"+keyword+"%").
		Find(&foods).Error
	return foods, err
}

// ParseNutritionResultJSON 解析营养结果JSON
func ParseNutritionResultJSON(result *model.NutritionResult) (*model.NutritionResultData, error) {
	if result == nil {
		return nil, nil
	}
	data, err := json.Marshal(result.ResultJSON)
	if err != nil {
		return nil, err
	}
	var nutrition model.NutritionResultData
	err = json.Unmarshal(data, &nutrition)
	return &nutrition, err
}

package repository

import (
	"context"
	"recipe-ai-backend/internal/model"

	"gorm.io/gorm"
)

// FavoriteRepository 收藏存储接口
type FavoriteRepository interface {
	Create(ctx context.Context, favorite *model.Favorite) error
	Delete(ctx context.Context, userOpenID string, recipeID int64) error
	Exists(ctx context.Context, userOpenID string, recipeID int64) (bool, error)
	ListByUser(ctx context.Context, userOpenID string, limit, offset int) ([]model.Favorite, error)
	CountByUser(ctx context.Context, userOpenID string) (int64, error)
	CountByRecipe(ctx context.Context, recipeID int64) (int64, error)
}

// favoriteRepository 收藏存储实现
type favoriteRepository struct {
	db *gorm.DB
}

// NewFavoriteRepository 创建收藏存储
func NewFavoriteRepository(db *gorm.DB) FavoriteRepository {
	return &favoriteRepository{db: db}
}

// Create 创建收藏
func (r *favoriteRepository) Create(ctx context.Context, favorite *model.Favorite) error {
	return r.db.WithContext(ctx).Create(favorite).Error
}

// Delete 取消收藏
func (r *favoriteRepository) Delete(ctx context.Context, userOpenID string, recipeID int64) error {
	return r.db.WithContext(ctx).
		Where("user_openid = ? AND recipe_id = ?", userOpenID, recipeID).
		Delete(&model.Favorite{}).Error
}

// Exists 检查是否已收藏
func (r *favoriteRepository) Exists(ctx context.Context, userOpenID string, recipeID int64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.Favorite{}).
		Where("user_openid = ? AND recipe_id = ?", userOpenID, recipeID).
		Count(&count).Error
	return count > 0, err
}

// ListByUser 列出用户的收藏
func (r *favoriteRepository) ListByUser(ctx context.Context, userOpenID string, limit, offset int) ([]model.Favorite, error) {
	var favorites []model.Favorite
	err := r.db.WithContext(ctx).
		Where("user_openid = ?", userOpenID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&favorites).Error
	return favorites, err
}

// CountByUser 统计用户收藏数量
func (r *favoriteRepository) CountByUser(ctx context.Context, userOpenID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.Favorite{}).
		Where("user_openid = ?", userOpenID).
		Count(&count).Error
	return count, err
}

// CountByRecipe 统计菜谱被收藏数量
func (r *favoriteRepository) CountByRecipe(ctx context.Context, recipeID int64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.Favorite{}).
		Where("recipe_id = ?", recipeID).
		Count(&count).Error
	return count, err
}

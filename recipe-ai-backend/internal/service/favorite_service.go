package service

import (
	"context"
	"fmt"
	"recipe-ai-backend/internal/model"
	"recipe-ai-backend/internal/repository"
)

// FavoriteService 收藏服务
type FavoriteService struct {
	favoriteRepo repository.FavoriteRepository
	recipeRepo   repository.RecipeRepository
	videoRepo    repository.VideoRepository
}

// NewFavoriteService 创建收藏服务
func NewFavoriteService(
	favoriteRepo repository.FavoriteRepository,
	recipeRepo repository.RecipeRepository,
	videoRepo repository.VideoRepository,
) *FavoriteService {
	return &FavoriteService{
		favoriteRepo: favoriteRepo,
		recipeRepo:   recipeRepo,
		videoRepo:    videoRepo,
	}
}

// ToggleFavorite 切换收藏状态
func (s *FavoriteService) ToggleFavorite(ctx context.Context, userOpenID string, recipeID int64) (bool, error) {
	exists, err := s.favoriteRepo.Exists(ctx, userOpenID, recipeID)
	if err != nil {
		return false, fmt.Errorf("检查收藏状态失败: %w", err)
	}

	if exists {
		// 取消收藏
		if err := s.favoriteRepo.Delete(ctx, userOpenID, recipeID); err != nil {
			return false, fmt.Errorf("取消收藏失败: %w", err)
		}
		// 减少菜谱收藏数（忽略错误，不影响主流程）
		_ = s.recipeRepo.DecrementFavoriteCount(ctx, recipeID)
		return false, nil
	}

	// 添加收藏
	favorite := &model.Favorite{
		UserOpenID: userOpenID,
		RecipeID:   recipeID,
	}
	if err := s.favoriteRepo.Create(ctx, favorite); err != nil {
		return false, fmt.Errorf("添加收藏失败: %w", err)
	}
	// 增加菜谱收藏数（忽略错误，不影响主流程）
	_ = s.recipeRepo.IncrementFavoriteCount(ctx, recipeID)
	return true, nil
}

// IsFavorited 检查是否已收藏
func (s *FavoriteService) IsFavorited(ctx context.Context, userOpenID string, recipeID int64) (bool, error) {
	return s.favoriteRepo.Exists(ctx, userOpenID, recipeID)
}

// ListUserFavorites 列出用户收藏
func (s *FavoriteService) ListUserFavorites(ctx context.Context, userOpenID string, limit, offset int) (*model.FavoritesResponse, error) {
	favorites, err := s.favoriteRepo.ListByUser(ctx, userOpenID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("获取收藏列表失败: %w", err)
	}

	items := make([]model.FavoriteItem, 0, len(favorites))
	for _, f := range favorites {
		recipe, err := s.recipeRepo.GetByID(ctx, f.RecipeID)
		if err != nil {
			continue // 跳过不存在的菜谱
		}

		video, _ := s.videoRepo.GetByID(ctx, recipe.VideoID)
		videoTitle := ""
		ownerName := ""
		if video != nil {
			videoTitle = video.Title
			ownerName = video.OwnerName
		}

		items = append(items, model.FavoriteItem{
			ID:            f.ID,
			RecipeID:      f.RecipeID,
			DishName:      recipe.DishName,
			VideoTitle:    videoTitle,
			OwnerName:     ownerName,
			ViewCount:     recipe.ViewCount,
			FavoriteCount: recipe.FavoriteCount,
			CreatedAt:     f.CreatedAt.Format("2006-01-02"),
		})
	}

	return &model.FavoritesResponse{Favorites: items}, nil
}

// CountByUser 统计用户收藏数量
func (s *FavoriteService) CountByUser(ctx context.Context, userOpenID string) (int64, error) {
	return s.favoriteRepo.CountByUser(ctx, userOpenID)
}

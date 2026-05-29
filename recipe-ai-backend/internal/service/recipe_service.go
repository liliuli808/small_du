package service

import (
	"context"
	"encoding/json"
	"fmt"
	"recipe-ai-backend/internal/model"
	"recipe-ai-backend/internal/repository"

	"gorm.io/gorm"
)

// RecipeService 菜谱服务
type RecipeService struct {
	db            *gorm.DB
	recipeRepo    repository.RecipeRepository
	nutritionRepo repository.NutritionRepository
	textRepo      repository.TextSourceRepository
	videoRepo     repository.VideoRepository
}

// NewRecipeService 创建菜谱服务
func NewRecipeService(
	db *gorm.DB,
	recipeRepo repository.RecipeRepository,
	nutritionRepo repository.NutritionRepository,
	textRepo repository.TextSourceRepository,
	videoRepo repository.VideoRepository,
) *RecipeService {
	return &RecipeService{
		db:            db,
		recipeRepo:    recipeRepo,
		nutritionRepo: nutritionRepo,
		textRepo:      textRepo,
		videoRepo:     videoRepo,
	}
}

// FindExistingRecipeByBVID 查询该BVID是否已有解析结果
func (s *RecipeService) FindExistingRecipeByBVID(ctx context.Context, bvid string) (*model.Recipe, error) {
	return s.recipeRepo.GetByBVID(ctx, bvid)
}

// SaveResult 保存解析结果
func (s *RecipeService) SaveResult(ctx context.Context, videoInfo *model.VideoInfo, bundle TextSourceBundle, recipeData *model.RecipeData, nutrition *model.NutritionResultData) (int64, error) {
	// 保存视频信息
	video := &model.BilibiliVideo{
		BVID:            videoInfo.BVID,
		AID:             videoInfo.AID,
		CID:             videoInfo.CID,
		Title:           videoInfo.Title,
		Description:     videoInfo.Description,
		OwnerName:       videoInfo.OwnerName,
		DurationSeconds: videoInfo.Duration,
	}

	video, err := s.videoRepo.CreateOrGet(ctx, video)
	if err != nil {
		return 0, fmt.Errorf("保存视频信息失败: %w", err)
	}

	// 保存文本来源
	if err := repository.SaveSources(ctx, s.db, video.ID, bundle.Subtitle, bundle.Description, bundle.TopComments); err != nil {
		// 非关键错误，继续
	}

	// 保存菜谱
	recipeJSON, _ := json.Marshal(recipeData)
	recipe := &model.Recipe{
		VideoID:    video.ID,
		DishName:   recipeData.DishName,
		Servings:   recipeData.Servings,
		RecipeJSON: model.JSONB{},
		Confidence: recipeData.Confidence,
	}
	json.Unmarshal(recipeJSON, &recipe.RecipeJSON)

	if err := s.recipeRepo.Create(ctx, recipe); err != nil {
		return 0, fmt.Errorf("保存菜谱失败: %w", err)
	}

	// 保存营养结果
	nutritionJSON, _ := json.Marshal(nutrition)
	nutritionResult := &model.NutritionResult{
		RecipeID:       recipe.ID,
		TotalKcal:      nutrition.TotalKcal,
		KcalPerServing: nutrition.KcalPerServing,
		ResultJSON:     model.JSONB{},
	}
	json.Unmarshal(nutritionJSON, &nutritionResult.ResultJSON)

	if err := s.nutritionRepo.CreateResult(ctx, nutritionResult); err != nil {
		return 0, fmt.Errorf("保存营养结果失败: %w", err)
	}

	return recipe.ID, nil
}

// GetRecipeResponse 获取完整菜谱响应
func (s *RecipeService) GetRecipeResponse(ctx context.Context, recipeID int64) (*model.RecipeResponse, error) {
	// 获取菜谱
	recipe, err := s.recipeRepo.GetByID(ctx, recipeID)
	if err != nil {
		return nil, fmt.Errorf("获取菜谱失败: %w", err)
	}

	// 解析菜谱JSON
	recipeData, err := repository.ParseRecipeJSON(recipe)
	if err != nil {
		return nil, fmt.Errorf("解析菜谱数据失败: %w", err)
	}

	// 获取营养结果
	nutritionResult, err := s.nutritionRepo.GetResultByRecipeID(ctx, recipeID)
	if err != nil {
		return nil, fmt.Errorf("获取营养结果失败: %w", err)
	}

	nutritionData, err := repository.ParseNutritionResultJSON(nutritionResult)
	if err != nil {
		return nil, fmt.Errorf("解析营养数据失败: %w", err)
	}

	// 获取视频信息 - 通过recipe关联
	video, _ := s.videoRepo.GetByBVID(ctx, "")
	if video == nil {
		video = &model.BilibiliVideo{}
	}

	// 获取文本来源
	sources, _ := s.textRepo.GetByVideoID(ctx, recipe.VideoID)

	// 构建来源类型列表
	sourceTypes := make([]string, 0)
	hasSubtitle := false
	hasDescription := false
	commentCount := 0
	for _, src := range sources {
		switch src.SourceType {
		case model.SourceTypeSubtitle:
			hasSubtitle = true
			sourceTypes = append(sourceTypes, "subtitle")
		case model.SourceTypeDescription:
			hasDescription = true
			sourceTypes = append(sourceTypes, "description")
		case model.SourceTypeTopComment:
			commentCount++
			sourceTypes = append(sourceTypes, "top_comment")
		case model.SourceTypeHotComment:
			commentCount++
			sourceTypes = append(sourceTypes, "hot_comment")
		}
	}

	return &model.RecipeResponse{
		Video: model.VideoInfoResponse{
			BVID:            video.BVID,
			Title:           video.Title,
			OwnerName:       video.OwnerName,
			DurationSeconds: video.DurationSeconds,
		},
		TextSources: model.TextSourcesResponse{
			HasSubtitle:    hasSubtitle,
			HasDescription: hasDescription,
			CommentCount:   commentCount,
			SourceTypes:    uniqueStrings(sourceTypes),
		},
		Recipe:    *recipeData,
		Nutrition: *nutritionData,
	}, nil
}

func uniqueStrings(arr []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, s := range arr {
		if !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}
	return result
}

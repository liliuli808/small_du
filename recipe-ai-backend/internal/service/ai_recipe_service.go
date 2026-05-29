package service

import (
	"context"
	"fmt"
	"recipe-ai-backend/internal/client"
	"recipe-ai-backend/internal/model"
)

// AIRecipeService AI菜谱服务
type AIRecipeService struct {
	client *client.AIClient
}

// NewAIRecipeService 创建AI菜谱服务
func NewAIRecipeService(client *client.AIClient) *AIRecipeService {
	return &AIRecipeService{client: client}
}

// Extract 调用AI解析菜谱
func (s *AIRecipeService) Extract(ctx context.Context, input string) (*model.RecipeData, error) {
	if input == "" {
		return nil, fmt.Errorf("输入文本为空")
	}

	recipe, err := s.client.ExtractRecipe(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("AI解析失败: %w", err)
	}

	return recipe, nil
}

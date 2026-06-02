package service

import (
	"context"
	"fmt"
	"math"
	"recipe-ai-backend/internal/model"
	"recipe-ai-backend/internal/pkg/logger"
	"recipe-ai-backend/internal/repository"
	"strings"
)

// NutritionService 营养计算服务
type NutritionService struct {
	repo repository.NutritionRepository
}

// NewNutritionService 创建营养计算服务
func NewNutritionService(repo repository.NutritionRepository) *NutritionService {
	return &NutritionService{repo: repo}
}

// Calculate 计算菜谱营养
func (s *NutritionService) Calculate(ctx context.Context, recipe *model.RecipeData) (*model.NutritionResultData, error) {
	foods, err := s.repo.GetAllFoods(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取营养数据库失败: %w", err)
	}

	result := &model.NutritionResultData{}

	for _, ing := range recipe.Ingredients {
		food, matched := s.matchNutritionFood(ing.Name, foods)
		if !matched {
			result.Details = append(result.Details, model.IngredientNutritionItem{
				Name:       ing.Name,
				Grams:      ing.Grams,
				Matched:    false,
				Confidence: ing.Confidence,
			})
			continue
		}

		ratio := ing.Grams / 100.0

		kcal := food.KcalPer100g * ratio
		protein := food.ProteinPer100g * ratio
		fat := food.FatPer100g * ratio
		carbs := food.CarbsPer100g * ratio

		result.TotalKcal += kcal
		result.TotalProtein += protein
		result.TotalFat += fat
		result.TotalCarbs += carbs

		result.Details = append(result.Details, model.IngredientNutritionItem{
			Name:        ing.Name,
			MatchedName: food.CanonicalName,
			Grams:       ing.Grams,
			Matched:     true,
			Kcal:        round(kcal, 1),
			Protein:     round(protein, 1),
			Fat:         round(fat, 1),
			Carbs:       round(carbs, 1),
			Source:      food.Source,
			Confidence:  ing.Confidence,
		})
	}

	servings := recipe.Servings
	if servings <= 0 {
		servings = 1
	}

	result.TotalKcal = round(result.TotalKcal, 1)
	result.TotalProtein = round(result.TotalProtein, 1)
	result.TotalFat = round(result.TotalFat, 1)
	result.TotalCarbs = round(result.TotalCarbs, 1)

	result.KcalPerServing = round(result.TotalKcal/float64(servings), 1)
	result.ProteinPerServing = round(result.TotalProtein/float64(servings), 1)
	result.FatPerServing = round(result.TotalFat/float64(servings), 1)
	result.CarbsPerServing = round(result.TotalCarbs/float64(servings), 1)

	logger.Info("营养计算完成",
		logger.String("dish", recipe.DishName),
		logger.Float64("total_kcal", result.TotalKcal))

	return result, nil
}

// Recalculate 重新计算营养（用户修改用量后）
func (s *NutritionService) Recalculate(ctx context.Context, servings int, ingredients []model.IngredientAdjustment) (*model.NutritionResultData, error) {
	foods, err := s.repo.GetAllFoods(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取营养数据库失败: %w", err)
	}

	result := &model.NutritionResultData{}

	for _, ing := range ingredients {
		food, matched := s.matchNutritionFood(ing.Name, foods)
		if !matched {
			result.Details = append(result.Details, model.IngredientNutritionItem{
				Name:    ing.Name,
				Grams:   ing.Grams,
				Matched: false,
			})
			continue
		}

		ratio := ing.Grams / 100.0

		result.TotalKcal += food.KcalPer100g * ratio
		result.TotalProtein += food.ProteinPer100g * ratio
		result.TotalFat += food.FatPer100g * ratio
		result.TotalCarbs += food.CarbsPer100g * ratio

		result.Details = append(result.Details, model.IngredientNutritionItem{
			Name:        ing.Name,
			MatchedName: food.CanonicalName,
			Grams:       ing.Grams,
			Matched:     true,
			Kcal:        round(food.KcalPer100g*ratio, 1),
			Protein:     round(food.ProteinPer100g*ratio, 1),
			Fat:         round(food.FatPer100g*ratio, 1),
			Carbs:       round(food.CarbsPer100g*ratio, 1),
		})
	}

	if servings <= 0 {
		servings = 1
	}

	result.TotalKcal = round(result.TotalKcal, 1)
	result.TotalProtein = round(result.TotalProtein, 1)
	result.TotalFat = round(result.TotalFat, 1)
	result.TotalCarbs = round(result.TotalCarbs, 1)

	result.KcalPerServing = round(result.TotalKcal/float64(servings), 1)
	result.ProteinPerServing = round(result.TotalProtein/float64(servings), 1)
	result.FatPerServing = round(result.TotalFat/float64(servings), 1)
	result.CarbsPerServing = round(result.TotalCarbs/float64(servings), 1)

	return result, nil
}

// matchNutritionFood 匹配营养食材
func (s *NutritionService) matchNutritionFood(name string, foods []model.NutritionFood) (*model.NutritionFood, bool) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, false
	}

	// 1. 精确匹配
	for _, f := range foods {
		if strings.EqualFold(f.CanonicalName, name) {
			return &f, true
		}
	}

	// 2. 别名匹配
	for _, f := range foods {
		for _, alias := range f.Aliases {
			if strings.EqualFold(alias, name) {
				return &f, true
			}
		}
	}

	// 3. 包含匹配（模糊匹配）
	for _, f := range foods {
		if strings.Contains(name, f.CanonicalName) || strings.Contains(f.CanonicalName, name) {
			return &f, true
		}
		for _, alias := range f.Aliases {
			if strings.Contains(name, alias) || strings.Contains(alias, name) {
				return &f, true
			}
		}
	}

	return nil, false
}

// round 四舍五入
func round(val float64, precision int) float64 {
	p := math.Pow(10, float64(precision))
	return math.Round(val*p) / p
}

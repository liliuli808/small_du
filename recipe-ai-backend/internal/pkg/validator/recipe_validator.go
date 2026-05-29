package validator

import (
	"recipe-ai-backend/internal/model"
	"strings"
)

// RecipeValidator 菜谱校验器
type RecipeValidator struct{}

// NewRecipeValidator 创建校验器
func NewRecipeValidator() *RecipeValidator {
	return &RecipeValidator{}
}

// Normalize 规范化菜谱数据
func (v *RecipeValidator) Normalize(recipe *model.RecipeData) *model.RecipeData {
	if recipe == nil {
		return nil
	}

	// 菜名不能为空
	recipe.DishName = strings.TrimSpace(recipe.DishName)
	if recipe.DishName == "" {
		recipe.DishName = "未知菜品"
		recipe.Confidence = minFloat(recipe.Confidence, 0.3)
	}

	// 份数处理
	if recipe.Servings <= 0 || recipe.Servings > 20 {
		recipe.Servings = 1
		recipe.UncertainItems = append(recipe.UncertainItems, "份数无法确定，默认按1份计算")
	}

	// 校验材料
	for i := range recipe.Ingredients {
		v.validateIngredient(&recipe.Ingredients[i])
	}

	// 校验步骤
	for i := range recipe.Steps {
		v.validateStep(&recipe.Steps[i])
	}

	// 校验置信度
	if recipe.Confidence < 0 {
		recipe.Confidence = 0
	}
	if recipe.Confidence > 1 {
		recipe.Confidence = 1
	}

	return recipe
}

// validateIngredient 校验单个材料
func (v *RecipeValidator) validateIngredient(ing *model.Ingredient) {
	// 名称清洗
	ing.Name = strings.TrimSpace(ing.Name)

	// 克重不能为负
	if ing.Grams < 0 {
		ing.Grams = 0
	}

	// 置信度范围
	if ing.Confidence < 0 {
		ing.Confidence = 0
	}
	if ing.Confidence > 1 {
		ing.Confidence = 1
	}

	// 异常用量标记
	switch ing.Name {
	case "食用油", "植物油", "油":
		if ing.Grams > 100 {
			ing.Confidence = minFloat(ing.Confidence, 0.3)
		}
	case "盐":
		if ing.Grams > 30 {
			ing.Confidence = minFloat(ing.Confidence, 0.3)
		}
	case "白糖", "糖", "冰糖":
		if ing.Grams > 200 {
			ing.Confidence = minFloat(ing.Confidence, 0.3)
		}
	}

	// 单个食材超过5kg标记异常
	if ing.Grams > 5000 {
		ing.Confidence = minFloat(ing.Confidence, 0.2)
	}
}

// validateStep 校验单个步骤
func (v *RecipeValidator) validateStep(step *model.Step) {
	// 确保步骤序号正确
	if step.Order <= 0 {
		step.Order = 1
	}

	// 时间处理
	if step.StartTime < 0 {
		step.StartTime = 0
	}
	if step.EndTime < 0 {
		step.EndTime = 0
	}
	if step.EndTime < step.StartTime {
		step.EndTime = step.StartTime
	}

	// 确保描述不为空
	step.Description = strings.TrimSpace(step.Description)
	if step.Description == "" {
		step.Description = step.Title
	}
}

// ValidateForSave 校验是否可保存
func (v *RecipeValidator) ValidateForSave(recipe *model.RecipeData) error {
	if recipe.DishName == "" {
		return ErrEmptyDishName
	}
	if len(recipe.Ingredients) == 0 {
		return ErrEmptyIngredients
	}
	if len(recipe.Steps) == 0 {
		return ErrEmptySteps
	}
	return nil
}

// HasMinimumContent 检查是否有最低限度的内容
func (v *RecipeValidator) HasMinimumContent(recipe *model.RecipeData) bool {
	if recipe == nil {
		return false
	}
	return recipe.DishName != "" && len(recipe.Ingredients) > 0
}

func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

var (
	ErrEmptyDishName    = NewValidationError("菜名不能为空")
	ErrEmptyIngredients = NewValidationError("材料列表不能为空")
	ErrEmptySteps       = NewValidationError("步骤列表不能为空")
)

type ValidationError struct {
	Message string
}

func NewValidationError(msg string) *ValidationError {
	return &ValidationError{Message: msg}
}

func (e *ValidationError) Error() string {
	return e.Message
}

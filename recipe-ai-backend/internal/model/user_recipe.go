package model

import "time"

// UserRecipe 用户自创菜谱表
type UserRecipe struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserOpenID string    `gorm:"type:varchar(64);not null" json:"-"`
	DishName   string    `gorm:"type:varchar(255)" json:"dish_name"`
	Servings   int       `json:"servings"`
	RecipeJSON JSONB     `gorm:"type:jsonb" json:"recipe_json"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (UserRecipe) TableName() string {
	return "user_recipes"
}

// UserRecipeData 用户菜谱数据结构（与 RecipeData 兼容）
type UserRecipeData struct {
	DishName    string       `json:"dish_name"`
	Servings    int          `json:"servings"`
	Ingredients []Ingredient `json:"ingredients"`
	Steps       []Step       `json:"steps"`
	Tips        []string     `json:"tips"`
}

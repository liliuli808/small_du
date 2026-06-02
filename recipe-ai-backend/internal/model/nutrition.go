package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// NutritionFood 营养食材表
type NutritionFood struct {
	ID             int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	CanonicalName  string    `gorm:"type:varchar(255);not null;index" json:"canonical_name"`
	Aliases        StringArray `gorm:"type:jsonb" json:"aliases"`
	KcalPer100g    float64   `json:"kcal_per_100g"`
	ProteinPer100g float64   `json:"protein_per_100g"`
	FatPer100g     float64   `json:"fat_per_100g"`
	CarbsPer100g   float64   `json:"carbs_per_100g"`
	Source         string    `gorm:"type:varchar(128)" json:"source"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (NutritionFood) TableName() string {
	return "nutrition_foods"
}

// NutritionResult 营养结果表
type NutritionResult struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	RecipeID     int64     `json:"recipe_id"`
	TotalKcal    float64   `json:"total_kcal"`
	KcalPerServing float64 `json:"kcal_per_serving"`
	ResultJSON   JSONB     `gorm:"type:jsonb" json:"result_json"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (NutritionResult) TableName() string {
	return "recipe_nutrition_results"
}

// NutritionResultData 营养计算结果
type NutritionResultData struct {
	TotalKcal         float64                   `json:"total_kcal"`
	KcalPerServing    float64                   `json:"kcal_per_serving"`
	TotalProtein      float64                   `json:"total_protein"`
	TotalFat          float64                   `json:"total_fat"`
	TotalCarbs        float64                   `json:"total_carbs"`
	ProteinPerServing float64                   `json:"protein_per_serving"`
	FatPerServing     float64                   `json:"fat_per_serving"`
	CarbsPerServing   float64                   `json:"carbs_per_serving"`
	Details           []IngredientNutritionItem `json:"details"`
}

// IngredientNutritionItem 单个食材营养信息
type IngredientNutritionItem struct {
	Name        string  `json:"name"`
	MatchedName string  `json:"matched_name,omitempty"`
	Grams       float64 `json:"grams"`
	Matched     bool    `json:"matched"`
	Kcal        float64 `json:"kcal,omitempty"`
	Protein     float64 `json:"protein,omitempty"`
	Fat         float64 `json:"fat,omitempty"`
	Carbs       float64 `json:"carbs,omitempty"`
	Source      string  `json:"source,omitempty"`
	Confidence  float64 `json:"confidence,omitempty"`
}

// StringArray 字符串数组类型
type StringArray []string

// Value 写入数据库
func (s StringArray) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	return json.Marshal(s)
}

// Scan 从数据库读取
func (s *StringArray) Scan(value interface{}) error {
	if value == nil {
		*s = nil
		return nil
	}
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, s)
	case string:
		return json.Unmarshal([]byte(v), s)
	default:
		return nil
	}
}

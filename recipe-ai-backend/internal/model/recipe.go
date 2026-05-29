package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// Recipe 菜谱表
type Recipe struct {
	ID         int64           `gorm:"primaryKey;autoIncrement" json:"id"`
	VideoID    int64           `json:"video_id"`
	DishName   string          `gorm:"type:varchar(255)" json:"dish_name"`
	Servings   int             `json:"servings"`
	RecipeJSON JSONB           `gorm:"type:jsonb" json:"recipe_json"`
	Confidence float64         `json:"confidence"`
	CreatedAt  time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Recipe) TableName() string {
	return "recipes"
}

// RecipeData AI解析的菜谱数据结构
type RecipeData struct {
	DishName       string       `json:"dish_name"`
	Servings       int          `json:"servings"`
	Ingredients    []Ingredient `json:"ingredients"`
	Steps          []Step       `json:"steps"`
	Tips           []string     `json:"tips"`
	UncertainItems []string     `json:"uncertain_items"`
	Confidence     float64      `json:"confidence"`
}

// Ingredient 材料
type Ingredient struct {
	Name       string  `json:"name"`
	Amount     float64 `json:"amount"`
	Unit       string  `json:"unit"`
	Grams      float64 `json:"grams"`
	Confidence float64 `json:"confidence"`
	Source     string  `json:"source"`
	SourceText string  `json:"source_text"`
}

// Step 做菜步骤
type Step struct {
	Order           int      `json:"order"`
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	StartTime       float64  `json:"start_time"`
	EndTime         float64  `json:"end_time"`
	Techniques      []string `json:"techniques"`
	DurationMinutes int      `json:"duration_minutes"`
	Source          string   `json:"source"`
}

// VideoTextSource 文本来源表
type VideoTextSource struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	VideoID    int64     `json:"video_id"`
	SourceType string    `gorm:"type:varchar(32);not null" json:"source_type"`
	Content    string    `gorm:"type:text" json:"content"`
	RawJSON    JSONB     `gorm:"type:jsonb" json:"raw_json"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (VideoTextSource) TableName() string {
	return "video_text_sources"
}

// TextSourceType 文本来源类型
const (
	SourceTypeSubtitle   = "subtitle"
	SourceTypeDescription = "description"
	SourceTypeTopComment = "top_comment"
	SourceTypeHotComment = "hot_comment"
)

// JSONB 自定义JSONB类型
type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, j)
	case string:
		return json.Unmarshal([]byte(v), j)
	default:
		return nil
	}
}

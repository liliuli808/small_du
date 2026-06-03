package model

import "time"

// Favorite 收藏表
type Favorite struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserOpenID string    `gorm:"type:varchar(64);not null" json:"-"`
	RecipeID   int64     `json:"recipe_id"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (Favorite) TableName() string {
	return "favorites"
}

package model

import "time"

type User struct {
	OpenID     string    `gorm:"type:varchar(64);primaryKey" json:"openid"`
	Nickname   string    `gorm:"type:varchar(64);default:''" json:"nickname"`
	AvatarURL  string    `gorm:"type:varchar(512);default:''" json:"avatar_url"`
	LastLoginAt time.Time `json:"last_login_at"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (User) TableName() string {
	return "users"
}

type LoginRequest struct {
	Code string `json:"code" binding:"required"`
}

type LoginResponse struct {
	OpenID string `json:"openid"`
	IsNew  bool   `json:"is_new"`
}

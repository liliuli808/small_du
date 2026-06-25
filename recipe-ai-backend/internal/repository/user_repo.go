package repository

import (
	"context"
	"recipe-ai-backend/internal/model"
	"time"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetByOpenID(ctx context.Context, openid string) (*model.User, error)
	Create(ctx context.Context, user *model.User) error
	UpdateLastLogin(ctx context.Context, openid string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByOpenID(ctx context.Context, openid string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("openid = ?", openid).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) UpdateLastLogin(ctx context.Context, openid string) error {
	return r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("openid = ?", openid).
		Update("last_login_at", time.Now()).Error
}

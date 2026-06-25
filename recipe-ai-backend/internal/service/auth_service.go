package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"recipe-ai-backend/internal/model"
	"recipe-ai-backend/internal/repository"
	"time"
)

type AuthService struct {
	userRepo repository.UserRepository
	secret   string
}

func NewAuthService(userRepo repository.UserRepository, secret string) *AuthService {
	return &AuthService{userRepo: userRepo, secret: secret}
}

func (s *AuthService) WxLogin(ctx context.Context, code string) (*model.LoginResponse, error) {
	// 开发环境：从 code 用 HMAC 生成 openid
	// 生产环境应替换为实际的 WeChat API 调用:
	//   GET https://api.weixin.qq.com/sns/jscode2session?appid=APPID&secret=SECRET&js_code=CODE&grant_type=authorization_code
	openid := s.generateOpenID(code)

	existing, err := s.userRepo.GetByOpenID(ctx, openid)
	isNew := false
	if err != nil {
		user := &model.User{
			OpenID:      openid,
			LastLoginAt: time.Now(),
		}
		if err := s.userRepo.Create(ctx, user); err != nil {
			return nil, fmt.Errorf("创建用户失败: %w", err)
		}
		isNew = true
	} else {
		s.userRepo.UpdateLastLogin(ctx, existing.OpenID)
	}

	return &model.LoginResponse{
		OpenID: openid,
		IsNew:  isNew,
	}, nil
}

func (s *AuthService) generateOpenID(code string) string {
	mac := hmac.New(sha256.New, []byte(s.secret))
	mac.Write([]byte(code))
	return fmt.Sprintf("o_%x", mac.Sum(nil)[:12])
}

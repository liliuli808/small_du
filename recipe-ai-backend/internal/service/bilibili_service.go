package service

import (
	"context"
	"recipe-ai-backend/internal/client"
	"recipe-ai-backend/internal/model"
)

// BilibiliService B站服务
type BilibiliService struct {
	client *client.BilibiliClient
}

// NewBilibiliService 创建B站服务
func NewBilibiliService(client *client.BilibiliClient) *BilibiliService {
	return &BilibiliService{client: client}
}

// ExtractBVID 提取BV号
func (s *BilibiliService) ExtractBVID(ctx context.Context, rawURL string) (string, error) {
	return s.client.ExtractBVID(rawURL)
}

// GetVideoInfo 获取视频信息
func (s *BilibiliService) GetVideoInfo(ctx context.Context, bvid string) (*model.VideoInfo, error) {
	return s.client.GetVideoInfo(ctx, bvid)
}

// TryGetSubtitle 尝试获取字幕
func (s *BilibiliService) TryGetSubtitle(ctx context.Context, videoInfo *model.VideoInfo) (string, bool) {
	return s.client.TryGetSubtitle(ctx, videoInfo)
}

// GetRecipeLikeComments 获取菜谱相关评论
func (s *BilibiliService) GetRecipeLikeComments(ctx context.Context, aid int64, limit int) ([]model.CommentItem, error) {
	return s.client.GetRecipeLikeComments(ctx, aid, limit)
}

// IsValidURL 检查URL是否有效
func (s *BilibiliService) IsValidURL(input string) bool {
	return s.client.IsValidURL(input)
}

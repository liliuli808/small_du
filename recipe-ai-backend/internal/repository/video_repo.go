package repository

import (
	"context"
	"recipe-ai-backend/internal/model"

	"gorm.io/gorm"
)

// VideoRepository 视频存储接口
type VideoRepository interface {
	CreateOrGet(ctx context.Context, video *model.BilibiliVideo) (*model.BilibiliVideo, error)
	GetByID(ctx context.Context, id int64) (*model.BilibiliVideo, error)
	GetByBVID(ctx context.Context, bvid string) (*model.BilibiliVideo, error)
}

// videoRepository 视频存储实现
type videoRepository struct {
	db *gorm.DB
}

// NewVideoRepository 创建视频存储
func NewVideoRepository(db *gorm.DB) VideoRepository {
	return &videoRepository{db: db}
}

// CreateOrGet 创建或获取视频记录
func (r *videoRepository) CreateOrGet(ctx context.Context, video *model.BilibiliVideo) (*model.BilibiliVideo, error) {
	var existing model.BilibiliVideo
	err := r.db.WithContext(ctx).Where("bvid = ?", video.BVID).First(&existing).Error
	if err == nil {
		// 已存在，更新信息
		existing.Title = video.Title
		existing.Description = video.Description
		existing.OwnerName = video.OwnerName
		existing.DurationSeconds = video.DurationSeconds
		existing.AID = video.AID
		existing.CID = video.CID
		r.db.WithContext(ctx).Save(&existing)
		return &existing, nil
	}

	// 不存在，创建新记录
	if err := r.db.WithContext(ctx).Create(video).Error; err != nil {
		return nil, err
	}
	return video, nil
}

// GetByID 根据ID查询
func (r *videoRepository) GetByID(ctx context.Context, id int64) (*model.BilibiliVideo, error) {
	var video model.BilibiliVideo
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&video).Error
	if err != nil {
		return nil, err
	}
	return &video, nil
}

// GetByBVID 根据BVID查询
func (r *videoRepository) GetByBVID(ctx context.Context, bvid string) (*model.BilibiliVideo, error) {
	var video model.BilibiliVideo
	err := r.db.WithContext(ctx).Where("bvid = ?", bvid).First(&video).Error
	if err != nil {
		return nil, err
	}
	return &video, nil
}

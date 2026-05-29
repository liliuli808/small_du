package repository

import (
	"context"
	"encoding/json"
	"recipe-ai-backend/internal/model"

	"gorm.io/gorm"
)

// TextSourceRepository 文本来源存储接口
type TextSourceRepository interface {
	Create(ctx context.Context, source *model.VideoTextSource) error
	GetByVideoID(ctx context.Context, videoID int64) ([]model.VideoTextSource, error)
}

// textSourceRepository 文本来源存储实现
type textSourceRepository struct {
	db *gorm.DB
}

// NewTextSourceRepository 创建文本来源存储
func NewTextSourceRepository(db *gorm.DB) TextSourceRepository {
	return &textSourceRepository{db: db}
}

// Create 创建文本来源记录
func (r *textSourceRepository) Create(ctx context.Context, source *model.VideoTextSource) error {
	return r.db.WithContext(ctx).Create(source).Error
}

// GetByVideoID 根据视频ID查询所有来源
func (r *textSourceRepository) GetByVideoID(ctx context.Context, videoID int64) ([]model.VideoTextSource, error) {
	var sources []model.VideoTextSource
	err := r.db.WithContext(ctx).Where("video_id = ?", videoID).Find(&sources).Error
	return sources, err
}

// SaveSources 批量保存文本来源
func SaveSources(ctx context.Context, db *gorm.DB, videoID int64, subtitle, description string, comments []model.CommentItem) error {
	// 保存字幕
	if subtitle != "" {
		subSource := &model.VideoTextSource{
			VideoID:    videoID,
			SourceType: model.SourceTypeSubtitle,
			Content:    subtitle,
		}
		if err := db.WithContext(ctx).Create(subSource).Error; err != nil {
			return err
		}
	}

	// 保存简介
	if description != "" {
		descSource := &model.VideoTextSource{
			VideoID:    videoID,
			SourceType: model.SourceTypeDescription,
			Content:    description,
		}
		if err := db.WithContext(ctx).Create(descSource).Error; err != nil {
			return err
		}
	}

	// 保存评论
	for _, comment := range comments {
		sourceType := model.SourceTypeHotComment
		if comment.IsTop {
			sourceType = model.SourceTypeTopComment
		}

		rawJSON, _ := json.Marshal(comment)
		commentSource := &model.VideoTextSource{
			VideoID:    videoID,
			SourceType: sourceType,
			Content:    comment.Message,
			RawJSON:    model.JSONB{"raw": string(rawJSON)},
		}
		if err := db.WithContext(ctx).Create(commentSource).Error; err != nil {
			return err
		}
	}

	return nil
}

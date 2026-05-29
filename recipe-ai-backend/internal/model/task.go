package model

import (
	"time"
)

// AnalyzeTask 解析任务表
type AnalyzeTask struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	TaskID       string    `gorm:"type:varchar(64);uniqueIndex;not null" json:"task_id"`
	Status       string    `gorm:"type:varchar(32);not null" json:"status"` // pending, processing, done, failed
	Message      string    `gorm:"type:text" json:"message"`
	RecipeID     *int64    `gorm:"" json:"recipe_id"`
	ErrorMessage string    `gorm:"type:text" json:"error_message"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (AnalyzeTask) TableName() string {
	return "analyze_tasks"
}

// TaskStatus 任务状态常量
const (
	TaskStatusPending    = "pending"
	TaskStatusProcessing = "processing"
	TaskStatusDone       = "done"
	TaskStatusFailed     = "failed"
)

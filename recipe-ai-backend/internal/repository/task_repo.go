package repository

import (
	"context"
	"recipe-ai-backend/internal/model"

	"gorm.io/gorm"
)

// TaskRepository 任务存储接口
type TaskRepository interface {
	Create(ctx context.Context, task *model.AnalyzeTask) error
	GetByTaskID(ctx context.Context, taskID string) (*model.AnalyzeTask, error)
	UpdateStatus(ctx context.Context, taskID, status, message string) error
	UpdateDone(ctx context.Context, taskID string, recipeID int64) error
	UpdateFailed(ctx context.Context, taskID, message, errorMsg string) error
}

// taskRepository 任务存储实现
type taskRepository struct {
	db *gorm.DB
}

// NewTaskRepository 创建任务存储
func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

// Create 创建任务
func (r *taskRepository) Create(ctx context.Context, task *model.AnalyzeTask) error {
	return r.db.WithContext(ctx).Create(task).Error
}

// GetByTaskID 根据taskID查询
func (r *taskRepository) GetByTaskID(ctx context.Context, taskID string) (*model.AnalyzeTask, error) {
	var task model.AnalyzeTask
	err := r.db.WithContext(ctx).Where("task_id = ?", taskID).First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// UpdateStatus 更新任务状态
func (r *taskRepository) UpdateStatus(ctx context.Context, taskID, status, message string) error {
	return r.db.WithContext(ctx).Model(&model.AnalyzeTask{}).
		Where("task_id = ?", taskID).
		Updates(map[string]interface{}{
			"status":  status,
			"message": message,
		}).Error
}

// UpdateDone 更新任务为完成状态
func (r *taskRepository) UpdateDone(ctx context.Context, taskID string, recipeID int64) error {
	return r.db.WithContext(ctx).Model(&model.AnalyzeTask{}).
		Where("task_id = ?", taskID).
		Updates(map[string]interface{}{
			"status":    model.TaskStatusDone,
			"message":   "解析完成",
			"recipe_id": recipeID,
		}).Error
}

// UpdateFailed 更新任务为失败状态
func (r *taskRepository) UpdateFailed(ctx context.Context, taskID, message, errorMsg string) error {
	return r.db.WithContext(ctx).Model(&model.AnalyzeTask{}).
		Where("task_id = ?", taskID).
		Updates(map[string]interface{}{
			"status":        model.TaskStatusFailed,
			"message":       message,
			"error_message": errorMsg,
		}).Error
}

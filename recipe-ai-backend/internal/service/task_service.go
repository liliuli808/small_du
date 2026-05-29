package service

import (
	"context"
	"fmt"
	"recipe-ai-backend/internal/model"
	"recipe-ai-backend/internal/repository"
	"time"
)

// TaskService 任务服务
type TaskService struct {
	repo repository.TaskRepository
}

// NewTaskService 创建任务服务
func NewTaskService(repo repository.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

// Create 创建任务
func (s *TaskService) Create(ctx context.Context) (string, error) {
	taskID := generateTaskID()
	task := &model.AnalyzeTask{
		TaskID:  taskID,
		Status:  model.TaskStatusPending,
		Message: "任务已创建",
	}
	if err := s.repo.Create(ctx, task); err != nil {
		return "", fmt.Errorf("创建任务失败: %w", err)
	}
	return taskID, nil
}

// GetStatus 获取任务状态
func (s *TaskService) GetStatus(ctx context.Context, taskID string) (*model.TaskResponse, error) {
	task, err := s.repo.GetByTaskID(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("查询任务失败: %w", err)
	}

	resp := &model.TaskResponse{
		TaskID:  task.TaskID,
		Status:  task.Status,
		Message: task.Message,
	}
	if task.RecipeID != nil {
		resp.RecipeID = task.RecipeID
	}
	return resp, nil
}

// Update 更新任务状态
func (s *TaskService) Update(ctx context.Context, taskID, status, message string) error {
	return s.repo.UpdateStatus(ctx, taskID, status, message)
}

// UpdateProgress 更新任务进度
func (s *TaskService) UpdateProgress(ctx context.Context, taskID, message string) error {
	return s.repo.UpdateStatus(ctx, taskID, model.TaskStatusProcessing, message)
}

// Done 标记任务完成
func (s *TaskService) Done(ctx context.Context, taskID string, recipeID int64) error {
	return s.repo.UpdateDone(ctx, taskID, recipeID)
}

// Fail 标记任务失败
func (s *TaskService) Fail(ctx context.Context, taskID, message, errorMsg string) error {
	return s.repo.UpdateFailed(ctx, taskID, message, errorMsg)
}

// generateTaskID 生成任务ID
func generateTaskID() string {
	return fmt.Sprintf("task_%s", time.Now().Format("20060102150405"))
}

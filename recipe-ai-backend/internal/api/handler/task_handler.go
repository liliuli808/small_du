package handler

import (
	"net/http"
	"recipe-ai-backend/internal/model"
	"recipe-ai-backend/internal/service"

	"github.com/gin-gonic/gin"
)

// TaskHandler 任务Handler
type TaskHandler struct {
	taskService *service.TaskService
}

// NewTaskHandler 创建任务Handler
func NewTaskHandler(taskService *service.TaskService) *TaskHandler {
	return &TaskHandler{taskService: taskService}
}

// GetTaskStatus 获取任务状态
func (h *TaskHandler) GetTaskStatus(c *gin.Context) {
	taskID := c.Param("task_id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Code:    "INVALID_TASK_ID",
			Message: "任务ID不能为空",
		})
		return
	}

	resp, err := h.taskService.GetStatus(c.Request.Context(), taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, model.APIResponse{
			Code:    "TASK_NOT_FOUND",
			Message: "任务不存在",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

package handler

import (
	"net/http"
	"recipe-ai-backend/internal/model"
	"recipe-ai-backend/internal/pkg/logger"
	"recipe-ai-backend/internal/pkg/parser"
	"recipe-ai-backend/internal/service"

	"github.com/gin-gonic/gin"
)

// AnalyzeHandler 解析Handler
type AnalyzeHandler struct {
	taskService     *service.TaskService
	bilibiliService *service.BilibiliService
	analyzeService  *service.AnalyzeService
	recipeService   *service.RecipeService
}

// NewAnalyzeHandler 创建解析Handler
func NewAnalyzeHandler(
	taskService *service.TaskService,
	bilibiliService *service.BilibiliService,
	analyzeService *service.AnalyzeService,
	recipeService *service.RecipeService,
) *AnalyzeHandler {
	return &AnalyzeHandler{
		taskService:     taskService,
		bilibiliService: bilibiliService,
		analyzeService:  analyzeService,
		recipeService:   recipeService,
	}
}

// CreateTask 创建解析任务
func (h *AnalyzeHandler) CreateTask(c *gin.Context) {
	var req model.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Code:    "INVALID_URL",
			Message: "请输入有效的B站视频链接",
		})
		return
	}

	req.URL = trimURL(req.URL)
	if req.URL == "" {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Code:    "INVALID_URL",
			Message: "链接不能为空",
		})
		return
	}

	if !parser.IsValidBilibiliURL(req.URL) {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Code:    "INVALID_URL",
			Message: "暂未识别到有效的B站视频链接",
		})
		return
	}

	ctx := c.Request.Context()

	// 同视频去重：提取 bvid 检查是否已有解析结果
	bvid, err := h.bilibiliService.ExtractBVID(ctx, req.URL)
	if err == nil && bvid != "" {
		existingRecipe, err := h.recipeService.FindExistingRecipeByBVID(ctx, bvid)
		if err == nil && existingRecipe != nil {
			// 已有结果，创建快捷任务并直接标记完成
			taskID, err := h.taskService.Create(ctx)
			if err != nil {
				logger.ErrorLog("创建去重任务失败", logger.Error(err))
				c.JSON(http.StatusInternalServerError, model.APIResponse{
					Code:    "INTERNAL_ERROR",
					Message: "任务创建失败，请稍后再试",
				})
				return
			}

			// 立即标记为完成
			if err := h.taskService.Done(ctx, taskID, existingRecipe.ID); err != nil {
				logger.ErrorLog("标记去重任务完成失败", logger.Error(err))
			}

			logger.Info("同视频去重命中，直接返回已有结果",
				logger.String("bvid", bvid),
				logger.Int64("recipe_id", existingRecipe.ID))

			c.JSON(http.StatusOK, model.CreateTaskResponse{
				TaskID:      taskID,
				Status:      model.TaskStatusDone,
				IsDuplicate: true,
				RecipeID:    &existingRecipe.ID,
			})
			return
		}
	}

	// 创建新任务
	taskID, err := h.taskService.Create(ctx)
	if err != nil {
		logger.ErrorLog("创建任务失败", logger.Error(err))
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Code:    "INTERNAL_ERROR",
			Message: "任务创建失败，请稍后再试",
		})
		return
	}

	// 异步投递任务
	go h.analyzeService.EnqueueTask(ctx, taskID, req.URL)

	c.JSON(http.StatusOK, model.CreateTaskResponse{
		TaskID:      taskID,
		Status:      model.TaskStatusPending,
		IsDuplicate: false,
	})
}

func trimURL(s string) string {
	// 简单去除前后空白和常见包裹字符
	for len(s) > 0 && (s[0] == ' ' || s[0] == '\n' || s[0] == '\t' || s[0] == '"' || s[0] == '\'') {
		s = s[1:]
	}
	for len(s) > 0 && (s[len(s)-1] == ' ' || s[len(s)-1] == '\n' || s[len(s)-1] == '\t' || s[len(s)-1] == '"' || s[len(s)-1] == '\'') {
		s = s[:len(s)-1]
	}
	return s
}

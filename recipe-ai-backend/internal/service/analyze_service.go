package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"recipe-ai-backend/internal/pkg/logger"
	"recipe-ai-backend/internal/pkg/validator"
	"time"

	"github.com/hibiken/asynq"
)

const (
	TypeAnalyzeBilibiliVideo = "bilibili:analyze"
)

// AnalyzeBilibiliPayload 分析任务Payload
type AnalyzeBilibiliPayload struct {
	TaskID string `json:"task_id"`
	URL    string `json:"url"`
}

// AnalyzeService 解析服务
type AnalyzeService struct {
	asynqClient     *asynq.Client
	taskService     *TaskService
	bilibiliService *BilibiliService
	textService     *TextSourceService
	aiService       *AIRecipeService
	nutritionService *NutritionService
	recipeService   *RecipeService
	recipeValidator *validator.RecipeValidator
}

// NewAnalyzeService 创建解析服务
func NewAnalyzeService(
	asynqClient *asynq.Client,
	taskService *TaskService,
	bilibiliService *BilibiliService,
	textService *TextSourceService,
	aiService *AIRecipeService,
	nutritionService *NutritionService,
	recipeService *RecipeService,
) *AnalyzeService {
	return &AnalyzeService{
		asynqClient:      asynqClient,
		taskService:      taskService,
		bilibiliService:  bilibiliService,
		textService:      textService,
		aiService:        aiService,
		nutritionService: nutritionService,
		recipeService:    recipeService,
		recipeValidator:  validator.NewRecipeValidator(),
	}
}

// EnqueueTask 投递任务到队列
func (s *AnalyzeService) EnqueueTask(ctx context.Context, taskID, url string) error {
	payload, err := json.Marshal(AnalyzeBilibiliPayload{
		TaskID: taskID,
		URL:    url,
	})
	if err != nil {
		return err
	}

	task := asynq.NewTask(TypeAnalyzeBilibiliVideo, payload)
	_, err = s.asynqClient.Enqueue(task, asynq.MaxRetry(2), asynq.Timeout(5*time.Minute))
	return err
}

// Process 处理解析任务（Worker调用）
func (s *AnalyzeService) Process(ctx context.Context, taskID, rawURL string) error {
	logger.Info("开始处理任务", logger.String("task_id", taskID))

	// Step 1: 解析视频链接
	s.taskService.UpdateProgress(ctx, taskID, "正在解析视频链接")
	bvid, err := s.bilibiliService.ExtractBVID(ctx, rawURL)
	if err != nil {
		s.taskService.Fail(ctx, taskID, "请输入有效的B站视频链接", err.Error())
		return err
	}

	// Step 2: 获取视频信息
	s.taskService.UpdateProgress(ctx, taskID, "正在获取视频信息")
	videoInfo, err := s.bilibiliService.GetVideoInfo(ctx, bvid)
	if err != nil {
		s.taskService.Fail(ctx, taskID, "获取视频信息失败", err.Error())
		return err
	}

	// Step 3: 获取字幕
	s.taskService.UpdateProgress(ctx, taskID, "正在获取字幕")
	subtitleText, hasSubtitle := s.bilibiliService.TryGetSubtitle(ctx, videoInfo)

	// Step 4: 获取评论
	s.taskService.UpdateProgress(ctx, taskID, "正在读取简介和评论")
	comments, _ := s.bilibiliService.GetRecipeLikeComments(ctx, videoInfo.AID, 5)

	bundle := TextSourceBundle{
		Title:       videoInfo.Title,
		Description: videoInfo.Description,
		Subtitle:    subtitleText,
		TopComments: comments,
		HasSubtitle: hasSubtitle,
	}

	// 检查是否有可用文本
	if !s.textService.HasUsefulText(bundle) {
		s.taskService.Fail(ctx, taskID, "没有找到字幕、简介或可用评论", "no useful text")
		return errors.New("no useful text")
	}

	// Step 5: 构建AI输入
	s.taskService.UpdateProgress(ctx, taskID, "正在整理文本内容")
	aiInput := s.textService.BuildAIInput(bundle)

	// Step 6: AI解析
	s.taskService.UpdateProgress(ctx, taskID, "正在用AI拆解菜谱")
	recipe, err := s.aiService.Extract(ctx, aiInput)
	if err != nil {
		s.taskService.Fail(ctx, taskID, "AI菜谱解析失败", err.Error())
		return err
	}

	// Step 7: 校验菜谱结果
	s.taskService.UpdateProgress(ctx, taskID, "正在校验菜谱结果")
	recipe = s.recipeValidator.Normalize(recipe)
	if !s.recipeValidator.HasMinimumContent(recipe) {
		s.taskService.Fail(ctx, taskID, "AI无法从该视频中提取有效菜谱", "insufficient content")
		return errors.New("insufficient content")
	}

	// Step 8: 计算热量
	s.taskService.UpdateProgress(ctx, taskID, "正在计算热量")
	nutrition, err := s.nutritionService.Calculate(ctx, recipe)
	if err != nil {
		s.taskService.Fail(ctx, taskID, "热量计算失败", err.Error())
		return err
	}

	// Step 9: 保存结果
	s.taskService.UpdateProgress(ctx, taskID, "正在保存结果")
	recipeID, err := s.recipeService.SaveResult(ctx, videoInfo, bundle, recipe, nutrition)
	if err != nil {
		s.taskService.Fail(ctx, taskID, "结果保存失败", err.Error())
		return err
	}

	// 完成任务
	s.taskService.Done(ctx, taskID, recipeID)
	logger.Info("任务处理完成",
		logger.String("task_id", taskID),
		logger.Int64("recipe_id", recipeID),
		logger.String("dish", recipe.DishName))

	return nil
}

// HandleAnalyzeTask Asynq任务处理器
func (s *AnalyzeService) HandleAnalyzeTask(ctx context.Context, t *asynq.Task) error {
	var payload AnalyzeBilibiliPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("解析任务payload失败: %w", err)
	}

	return s.Process(ctx, payload.TaskID, payload.URL)
}

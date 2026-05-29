package worker

import (
	"context"
	"recipe-ai-backend/internal/pkg/logger"
	"recipe-ai-backend/internal/service"

	"github.com/hibiken/asynq"
)

// AnalyzeVideoWorker 视频解析Worker
type AnalyzeVideoWorker struct {
	analyzeService *service.AnalyzeService
}

// NewAnalyzeVideoWorker 创建Worker
func NewAnalyzeVideoWorker(analyzeService *service.AnalyzeService) *AnalyzeVideoWorker {
	return &AnalyzeVideoWorker{analyzeService: analyzeService}
}

// RegisterHandlers 注册任务处理器
func (w *AnalyzeVideoWorker) RegisterHandlers(mux *asynq.ServeMux) {
	mux.HandleFunc(service.TypeAnalyzeBilibiliVideo, w.handleAnalyzeBilibiliVideo)
}

// handleAnalyzeBilibiliVideo 处理B站视频解析任务
func (w *AnalyzeVideoWorker) handleAnalyzeBilibiliVideo(ctx context.Context, t *asynq.Task) error {
	logger.Info("Worker收到任务", logger.String("type", t.Type()))
	return w.analyzeService.HandleAnalyzeTask(ctx, t)
}

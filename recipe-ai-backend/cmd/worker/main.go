package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"recipe-ai-backend/internal/client"
	"recipe-ai-backend/internal/pkg/config"
	"recipe-ai-backend/internal/pkg/database"
	"recipe-ai-backend/internal/pkg/logger"
	"recipe-ai-backend/internal/pkg/redis"
	"recipe-ai-backend/internal/repository"
	"recipe-ai-backend/internal/service"
	"recipe-ai-backend/internal/worker"

	"github.com/hibiken/asynq"
)

func main() {
	if err := logger.Init(); err != nil {
		fmt.Printf("初始化日志失败: %v\n", err)
		os.Exit(1)
	}

	// 加载配置
	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		cfgPath = "./config/config.yaml"
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		logger.Fatal("加载配置失败", logger.Error(err))
	}

	// 初始化数据库
	db, err := database.Init(cfg.Database)
	if err != nil {
		logger.Fatal("初始化数据库失败", logger.Error(err))
	}
	defer database.Close()

	// 初始化Redis
	redisClient := redis.NewClient(cfg.Redis)
	if err := redisClient.Ping(context.Background()); err != nil {
		logger.Fatal("Redis连接失败", logger.Error(err))
	}
	defer redisClient.Close()

	// 初始化HTTP客户端
	httpClient := client.NewHTTPClient(cfg.Bilibili.UserAgent, cfg.Bilibili.Cookie, cfg.Bilibili.Timeout)

	// 初始化B站客户端
	biliClient := client.NewBilibiliClient(httpClient)

	// 初始化AI客户端
	aiClient := client.NewAIClient(cfg.AI)

	// 初始化Repository
	taskRepo := repository.NewTaskRepository(db)
	videoRepo := repository.NewVideoRepository(db)
	textSourceRepo := repository.NewTextSourceRepository(db)
	recipeRepo := repository.NewRecipeRepository(db)
	nutritionRepo := repository.NewNutritionRepository(db)

	// 初始化Service
	taskService := service.NewTaskService(taskRepo)
	bilibiliService := service.NewBilibiliService(biliClient)
	textSourceService := service.NewTextSourceService()
	aiRecipeService := service.NewAIRecipeService(aiClient)
	nutritionService := service.NewNutritionService(nutritionRepo)
	recipeService := service.NewRecipeService(db, recipeRepo, nutritionRepo, textSourceRepo, videoRepo)

	// 初始化Asynq客户端（用于任务投递，如重试）
	asynqClient := asynq.NewClient(redisClient.AsynqRedisOpt())
	defer asynqClient.Close()

	analyzeService := service.NewAnalyzeService(
		asynqClient,
		taskService,
		bilibiliService,
		textSourceService,
		aiRecipeService,
		nutritionService,
		recipeService,
	)

	// 初始化Worker
	workerServer := worker.NewServer(redisClient.AsynqRedisOpt(), cfg.Worker)
	analyzeWorker := worker.NewAnalyzeVideoWorker(analyzeService)
	analyzeWorker.RegisterHandlers(workerServer.GetMux())

	// 启动Worker
	go func() {
		logger.Info("Worker服务启动", logger.Int("concurrency", cfg.Worker.Concurrency))
		if err := workerServer.Run(); err != nil {
			logger.Fatal("Worker服务启动失败", logger.Error(err))
		}
	}()

	// 等待退出信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("正在关闭Worker...")
	workerServer.Shutdown()
	logger.Info("Worker已关闭")
}

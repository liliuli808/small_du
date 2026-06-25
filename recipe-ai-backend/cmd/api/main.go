package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"recipe-ai-backend/internal/api"
	"recipe-ai-backend/internal/api/handler"
	"recipe-ai-backend/internal/client"
	"recipe-ai-backend/internal/pkg/config"
	"recipe-ai-backend/internal/pkg/database"
	"recipe-ai-backend/internal/pkg/logger"
	"recipe-ai-backend/internal/pkg/redis"
	"recipe-ai-backend/internal/repository"
	"recipe-ai-backend/internal/service"

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
		logger.Warn("Redis连接失败，继续启动", logger.Error(err))
	}
	defer redisClient.Close()

	// 初始化HTTP客户端
	httpClient := client.NewHTTPClient(cfg.Bilibili.UserAgent, cfg.Bilibili.Cookie, cfg.Bilibili.Timeout)

	// 初始化B站客户端
	biliClient := client.NewBilibiliClient(httpClient)

	// 初始化AI客户端
	aiClient := client.NewAIClient(cfg.AI)

	// 初始化Asynq客户端
	asynqClient := asynq.NewClient(redisClient.AsynqRedisOpt())
	defer asynqClient.Close()

	// 初始化Repository
	taskRepo := repository.NewTaskRepository(db)
	videoRepo := repository.NewVideoRepository(db)
	textSourceRepo := repository.NewTextSourceRepository(db)
	recipeRepo := repository.NewRecipeRepository(db)
	nutritionRepo := repository.NewNutritionRepository(db)
	userRecipeRepo := repository.NewUserRecipeRepository(db)
	favoriteRepo := repository.NewFavoriteRepository(db)
	userRepo := repository.NewUserRepository(db)

	// 初始化Service
	taskService := service.NewTaskService(taskRepo)
	bilibiliService := service.NewBilibiliService(biliClient)
	textSourceService := service.NewTextSourceService()
	aiRecipeService := service.NewAIRecipeService(aiClient)
	nutritionService := service.NewNutritionService(nutritionRepo)
	recipeService := service.NewRecipeService(db, recipeRepo, nutritionRepo, textSourceRepo, videoRepo)
	userRecipeService := service.NewUserRecipeService(userRecipeRepo)
	favoriteService := service.NewFavoriteService(favoriteRepo, recipeRepo, videoRepo)
	analyzeService := service.NewAnalyzeService(
		asynqClient,
		taskService,
		bilibiliService,
		textSourceService,
		aiRecipeService,
		nutritionService,
		recipeService,
	)

	// 初始化Service
	// (existing services above)
	authService := service.NewAuthService(userRepo, cfg.App.SecretKey)

	// 初始化Handler
	analyzeHandler := handler.NewAnalyzeHandler(taskService, bilibiliService, analyzeService, recipeService)
	taskHandler := handler.NewTaskHandler(taskService)
	recipeHandler := handler.NewRecipeHandler(recipeService, nutritionService)
	userRecipeHandler := handler.NewUserRecipeHandler(userRecipeService)
	favoriteHandler := handler.NewFavoriteHandler(favoriteService)
	authHandler := handler.NewAuthHandler(authService)

	// 设置路由
	router := api.NewRouter(analyzeHandler, taskHandler, recipeHandler, userRecipeHandler, favoriteHandler, authHandler)
	engine := router.Setup()

	// 启动HTTP服务
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.App.Port),
		Handler: engine,
	}

	go func() {
		logger.Info("API服务启动", logger.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("API服务启动失败", logger.Error(err))
		}
	}()

	// 等待退出信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("正在关闭服务...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.ErrorLog("服务关闭失败", logger.Error(err))
	}

	logger.Info("服务已关闭")
}

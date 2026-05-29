package api

import (
	"recipe-ai-backend/internal/api/handler"
	"recipe-ai-backend/internal/api/middleware"

	"github.com/gin-gonic/gin"
)

// Router API路由
type Router struct {
	engine          *gin.Engine
	analyzeHandler  *handler.AnalyzeHandler
	taskHandler     *handler.TaskHandler
	recipeHandler   *handler.RecipeHandler
}

// NewRouter 创建路由
func NewRouter(
	analyzeHandler *handler.AnalyzeHandler,
	taskHandler *handler.TaskHandler,
	recipeHandler *handler.RecipeHandler,
) *Router {
	engine := gin.New()
	engine.Use(middleware.Recover())
	engine.Use(middleware.CORS())

	return &Router{
		engine:         engine,
		analyzeHandler: analyzeHandler,
		taskHandler:    taskHandler,
		recipeHandler:  recipeHandler,
	}
}

// Setup 设置路由
func (r *Router) Setup() *gin.Engine {
	v1 := r.engine.Group("/api/v1")

	// 解析任务
	v1.POST("/analyze/bilibili", r.analyzeHandler.CreateTask)

	// 任务状态
	v1.GET("/tasks/:task_id", r.taskHandler.GetTaskStatus)

	// 菜谱结果
	v1.GET("/recipes/:recipe_id", r.recipeHandler.GetRecipe)

	// 重新计算热量
	v1.POST("/recipes/:recipe_id/recalculate", r.recipeHandler.Recalculate)

	// 健康检查
	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return r.engine
}

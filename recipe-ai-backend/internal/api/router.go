package api

import (
	"recipe-ai-backend/internal/api/handler"
	"recipe-ai-backend/internal/api/middleware"

	"github.com/gin-gonic/gin"
)

// Router API路由
type Router struct {
	engine             *gin.Engine
	analyzeHandler     *handler.AnalyzeHandler
	taskHandler        *handler.TaskHandler
	recipeHandler      *handler.RecipeHandler
	userRecipeHandler  *handler.UserRecipeHandler
	favoriteHandler    *handler.FavoriteHandler
	authHandler        *handler.AuthHandler
}

// NewRouter 创建路由
func NewRouter(
	analyzeHandler *handler.AnalyzeHandler,
	taskHandler *handler.TaskHandler,
	recipeHandler *handler.RecipeHandler,
	userRecipeHandler *handler.UserRecipeHandler,
	favoriteHandler *handler.FavoriteHandler,
	authHandler *handler.AuthHandler,
) *Router {
	engine := gin.New()
	engine.Use(middleware.Recover())
	engine.Use(middleware.CORS())
	engine.Use(middleware.UserOpenID())

	return &Router{
		engine:            engine,
		analyzeHandler:    analyzeHandler,
		taskHandler:       taskHandler,
		recipeHandler:     recipeHandler,
		userRecipeHandler: userRecipeHandler,
		favoriteHandler:   favoriteHandler,
		authHandler:       authHandler,
	}
}

// Setup 设置路由
func (r *Router) Setup() *gin.Engine {
	v1 := r.engine.Group("/api/v1")

	// 解析任务
	v1.POST("/analyze/bilibili", r.analyzeHandler.CreateTask)

	// 任务状态
	v1.GET("/tasks/:task_id", r.taskHandler.GetTaskStatus)

	// 热门菜谱
	v1.GET("/recipes/popular", r.recipeHandler.ListPopularRecipes)

	// 搜索菜谱
	v1.GET("/recipes/search", r.recipeHandler.SearchRecipes)

	// 菜谱结果
	v1.GET("/recipes/:recipe_id", r.recipeHandler.GetRecipe)

	// 菜谱收藏
	v1.POST("/recipes/:recipe_id/favorite", r.favoriteHandler.ToggleFavorite)
	v1.GET("/recipes/:recipe_id/favorite", r.favoriteHandler.GetFavoriteStatus)

	// 重新计算热量
	v1.POST("/recipes/:recipe_id/recalculate", r.recipeHandler.Recalculate)

	// AI菜谱派生（转为用户可编辑的数据）
	v1.GET("/recipes/:recipe_id/derive", r.recipeHandler.DeriveAsUserRecipe)

	// 用户菜谱
	v1.GET("/user/recipes", r.userRecipeHandler.ListUserRecipes)
	v1.POST("/user/recipes", r.userRecipeHandler.CreateUserRecipe)
	v1.GET("/user/recipes/:id", r.userRecipeHandler.GetUserRecipe)
	v1.PUT("/user/recipes/:id", r.userRecipeHandler.UpdateUserRecipe)
	v1.DELETE("/user/recipes/:id", r.userRecipeHandler.DeleteUserRecipe)

	// 用户收藏
	v1.GET("/user/favorites", r.favoriteHandler.ListFavorites)

	// 用户认证
	v1.POST("/auth/wx-login", r.authHandler.WxLogin)

	// 健康检查
	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return r.engine
}

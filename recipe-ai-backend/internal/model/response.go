package model

// APIResponse 通用API响应
type APIResponse struct {
	Code    string      `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// CreateTaskRequest 创建任务请求
type CreateTaskRequest struct {
	URL string `json:"url" binding:"required"`
}

// CreateTaskResponse 创建任务响应
type CreateTaskResponse struct {
	TaskID       string `json:"task_id"`
	Status       string `json:"status"`
	IsDuplicate  bool   `json:"is_duplicate"`
	RecipeID     *int64 `json:"recipe_id,omitempty"`
}

// TaskResponse 任务状态响应
type TaskResponse struct {
	TaskID   string `json:"task_id"`
	Status   string `json:"status"`
	Message  string `json:"message"`
	RecipeID *int64 `json:"recipe_id"`
}

// RecipeResponse 菜谱结果响应
type RecipeResponse struct {
	Video       VideoInfoResponse    `json:"video"`
	TextSources TextSourcesResponse  `json:"text_sources"`
	Recipe      RecipeData           `json:"recipe"`
	Nutrition   NutritionResultData  `json:"nutrition"`
}

// VideoInfoResponse 视频信息响应
type VideoInfoResponse struct {
	BVID           string `json:"bvid"`
	Title          string `json:"title"`
	OwnerName      string `json:"owner_name"`
	CoverURL       string `json:"cover_url"`
	DurationSeconds int   `json:"duration_seconds"`
}

// TextSourcesResponse 文本来源响应
type TextSourcesResponse struct {
	HasSubtitle   bool     `json:"has_subtitle"`
	HasDescription bool    `json:"has_description"`
	CommentCount  int      `json:"comment_count"`
	SourceTypes   []string `json:"source_types"`
}

// RecalculateRequest 重新计算热量请求
type RecalculateRequest struct {
	Servings    int                    `json:"servings" binding:"required,min=1"`
	Ingredients []IngredientAdjustment `json:"ingredients" binding:"required"`
}

// IngredientAdjustment 食材调整
type IngredientAdjustment struct {
	Name  string  `json:"name" binding:"required"`
	Grams float64 `json:"grams" binding:"required,min=0"`
}

// RecalculateResponse 重新计算热量响应
type RecalculateResponse struct {
	Nutrition NutritionResultData `json:"nutrition"`
}

// PopularRecipeItem 热门菜谱项
type PopularRecipeItem struct {
	ID            int64   `json:"id"`
	DishName      string  `json:"dish_name"`
	ViewCount     int     `json:"view_count"`
	FavoriteCount int     `json:"favorite_count"`
	VideoTitle    string  `json:"video_title"`
	OwnerName     string  `json:"owner_name"`
	IsFavorited   bool    `json:"is_favorited"`
}

// PopularRecipesResponse 热门菜谱响应
type PopularRecipesResponse struct {
	Recipes []PopularRecipeItem `json:"recipes"`
}

// FavoriteToggleResponse 收藏切换响应
type FavoriteToggleResponse struct {
	IsFavorited bool `json:"is_favorited"`
}

// FavoriteItem 收藏项
type FavoriteItem struct {
	ID            int64   `json:"id"`
	RecipeID      int64   `json:"recipe_id"`
	DishName      string  `json:"dish_name"`
	VideoTitle    string  `json:"video_title"`
	OwnerName     string  `json:"owner_name"`
	ViewCount     int     `json:"view_count"`
	FavoriteCount int     `json:"favorite_count"`
	CreatedAt     string  `json:"created_at"`
}

// FavoritesResponse 收藏列表响应
type FavoritesResponse struct {
	Favorites []FavoriteItem `json:"favorites"`
}

// UserRecipeListItem 用户菜谱列表项
type UserRecipeListItem struct {
	ID        int64  `json:"id"`
	DishName  string `json:"dish_name"`
	Servings  int    `json:"servings"`
	CreatedAt string `json:"created_at"`
}

// UserRecipesResponse 用户菜谱列表响应
type UserRecipesResponse struct {
	Recipes []UserRecipeListItem `json:"recipes"`
}

// UserRecipeDetailResponse 用户菜谱详情响应
type UserRecipeDetailResponse struct {
	ID       int64          `json:"id"`
	DishName string         `json:"dish_name"`
	Servings int            `json:"servings"`
	Recipe   UserRecipeData `json:"recipe"`
	CreatedAt string        `json:"created_at"`
}

// CreateUserRecipeRequest 创建用户菜谱请求
type CreateUserRecipeRequest struct {
	DishName    string       `json:"dish_name" binding:"required"`
	Servings    int          `json:"servings" binding:"required,min=1"`
	Ingredients []Ingredient `json:"ingredients" binding:"required"`
	Steps       []Step       `json:"steps" binding:"required"`
	Tips        []string     `json:"tips"`
}

// UpdateUserRecipeRequest 更新用户菜谱请求
type UpdateUserRecipeRequest struct {
	DishName    string       `json:"dish_name" binding:"required"`
	Servings    int          `json:"servings" binding:"required,min=1"`
	Ingredients []Ingredient `json:"ingredients" binding:"required"`
	Steps       []Step       `json:"steps" binding:"required"`
	Tips        []string     `json:"tips"`
}

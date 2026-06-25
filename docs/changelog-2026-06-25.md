# 改进记录 2026-06-25

## 1. 微信登录 (Item 1)

### 后端
- 新增 `users` 表（migrations/003_users.up.sql）
- 新增 `internal/model/user.go` — User 模型
- 新增 `internal/repository/user_repo.go` — UserRepository（GetByOpenID / Create / UpdateLastLogin）
- 新增 `internal/service/auth_service.go` — WxLogin 逻辑（开发环境用 HMAC 从 code 生成 openid）
- 新增 `internal/api/handler/auth_handler.go` — `POST /api/v1/auth/wx-login`
- 配置新增 `app.secret_key` 用于 HMAC 签名
- 路由注册: `POST /api/v1/auth/wx-login`

### 前端
- `app.js`: 新增 `doWxLogin()` 方法，onLaunch 时自动尝试微信登录；失败回退匿名ID
- `pages/mine/mine.wxml`: 未登录显示"微信一键登录"按钮，登录后显示已登录状态
- `pages/mine/mine.js`: 新增 `doLogin` 方法调用微信登录流程

---

## 2. 结果页编辑AI菜谱 (Item 2)

### 后端
- `recipe_service.go`: 新增 `DeriveAsUserRecipeData()` 方法，将 AI 菜谱数据转为用户可编辑格式
- `recipe_handler.go`: 新增 `DeriveAsUserRecipe` handler
- 路由: `GET /api/v1/recipes/:recipe_id/derive`

### 前端
- `pages/result/result.wxml`: 底部栏添加"编辑菜谱"按钮
- `pages/user-recipe-edit/user-recipe-edit.js`: 支持 `from_recipe_id` 参数，自动加载AI菜谱数据并预填表单

---

## 3. 分享功能 (Item 3)

### 前端
- `pages/result/result.js`: 新增 `onShareAppMessage` 生命周期方法
- `pages/result/result.wxml`: 底部栏添加 `button[open-type="share"]` 分享按钮
- `onLoad` 中调用 `wx.showShareMenu()` 启用系统分享菜单

---

## 4. 搜索功能 (Item 4)

### 后端
- `repository/recipe_repo.go`: 新增 `Search(ctx, keyword, limit, offset)` 方法（ILIKE 模糊搜索）
- `service/recipe_service.go`: 新增 `SearchRecipes()` 方法
- `handler/recipe_handler.go`: 新增 `SearchRecipes` handler
- 路由: `GET /api/v1/recipes/search?q=关键字&limit=20`

### 前端
- `pages/index/index.wxml`: 顶部添加搜索栏，输入关键字确认后显示搜索结果
- `pages/index/index.js`: 新增 `onSearchInput`、`doSearch`、`clearSearch` 方法
- `pages/index/index.wxss`: 搜索栏和搜索结果列表样式

---

## 5. 视频缩略图 (Item 5)

### 后端
- `model/bilibili.go`: `BilibiliVideo` 和 `VideoInfo` 新增 `CoverURL` 字段
- `model/response.go`: `VideoInfoResponse` 新增 `CoverURL` 字段
- `client/bilibili_client.go`: 从 B站 API 响应的 `pic` 字段提取封面 URL
- `service/recipe_service.go`: SaveResult 和 GetRecipeResponse 中传递 CoverURL
- 迁移: `migrations/004_cover_url.up.sql` — 添加 cover_url 列

### 前端
- `pages/result/result.wxml`: 有封面URL时显示 image，否则显示 emoji 占位
- `pages/result/result.wxss`: 缩略图样式

---

## 6. 设置页 (Item 6)

### 前端
- 新增 `pages/settings/settings` 页面（API地址修改、清除缓存、版本信息）
- `app.json`: 注册 settings 页面
- `pages/mine/mine.wxml`: 设置入口（应用设置 + 清除缓存）
- `pages/mine/mine.js`: 新增 `goToSettings`、`clearAllCache` 方法
- `app.js`: onLaunch 时从缓存恢复 API 地址

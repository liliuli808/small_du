# AI 菜谱解析助手 - Go 后端

B站做菜视频AI菜谱解析小程序的后端服务。

## 项目结构

```
recipe-ai-backend/
├── cmd/
│   ├── api/              # API 服务入口
│   └── worker/           # Worker 服务入口
├── internal/
│   ├── api/              # HTTP API 层
│   │   ├── handler/      # 请求处理器
│   │   └── middleware/   # 中间件
│   ├── client/           # 外部客户端
│   │   ├── bilibili_client.go  # B站 API 客户端
│   │   ├── ai_client.go        # AI 服务客户端
│   │   └── http_client.go      # HTTP 基础客户端
│   ├── model/            # 数据模型
│   ├── repository/       # 数据存储层
│   ├── service/          # 业务逻辑层
│   ├── worker/           # 异步 Worker
│   └── pkg/              # 公共包
│       ├── config/       # 配置管理
│       ├── database/     # 数据库连接
│       ├── redis/        # Redis 客户端
│       ├── logger/       # 日志工具
│       ├── parser/       # URL 解析
│       ├── unitconvert/  # 单位换算
│       ├── validator/    # 数据校验
│       └── xjson/        # JSON 修复
├── migrations/           # 数据库迁移
├── scripts/              # 脚本
├── config/               # 配置文件
├── Dockerfile
└── docker-compose.yml
```

## 技术栈

- **语言**: Go 1.22
- **Web 框架**: Gin
- **ORM**: GORM
- **数据库**: PostgreSQL
- **缓存/队列**: Redis + Asynq
- **日志**: Zap
- **配置**: Viper

## 快速开始

### 1. 环境要求

- Go 1.22+
- PostgreSQL 16
- Redis 7

### 2. 配置文件

复制配置文件模板：

```bash
cp config/config.example.yaml config/config.yaml
```

编辑 `config/config.yaml`，填入 AI API Key 等配置。

### 3. 数据库初始化

```bash
# 创建数据库
createdb recipe_ai

# 执行迁移
psql -U recipe -d recipe_ai -f migrations/001_init.up.sql

# 导入种子数据
psql -U recipe -d recipe_ai -f scripts/seed_nutrition_foods.sql
```

### 4. 运行服务

```bash
# 安装依赖
go mod download

# 运行 API 服务
go run cmd/api/main.go

# 运行 Worker 服务（另开终端）
go run cmd/worker/main.go
```

### 5. Docker 部署

```bash
docker-compose up -d
```

## API 接口

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/analyze/bilibili` | 创建 B 站解析任务 |
| GET | `/api/v1/tasks/{task_id}` | 查询任务状态 |
| GET | `/api/v1/recipes/{recipe_id}` | 获取菜谱结果 |
| POST | `/api/v1/recipes/{recipe_id}/recalculate` | 修改材料后重新计算热量 |

## 核心流程

1. 用户粘贴 B 站链接
2. API 创建异步任务，投递到 Redis 队列
3. Worker 消费任务：
   - 解析 BV 号
   - 获取视频信息
   - 获取字幕 / 简介 / 评论
   - 调用 AI 解析菜谱
   - 校验和标准化
   - 计算热量和营养素
   - 保存结果
4. 小程序轮询任务状态，完成后展示结果

## 配置说明

| 配置项 | 说明 | 必填 |
|--------|------|------|
| `ai.api_key` | AI 服务 API Key | 是 |
| `ai.base_url` | AI 服务 Base URL | 是 |
| `ai.model` | AI 模型名称 | 是 |
| `database.*` | PostgreSQL 连接配置 | 是 |
| `redis.*` | Redis 连接配置 | 是 |
| `bilibili.cookie` | B站 Cookie（可选） | 否 |

## 许可证

MIT

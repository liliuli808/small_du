package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"recipe-ai-backend/internal/model"
	"recipe-ai-backend/internal/pkg/config"
	"recipe-ai-backend/internal/pkg/logger"
	"recipe-ai-backend/internal/pkg/xjson"
	"time"
)

// AIClient AI服务客户端
type AIClient struct {
	apiKey      string
	baseURL     string
	model       string
	maxTokens   int
	temperature float64
	httpClient  *http.Client
}

// NewAIClient 创建AI客户端
func NewAIClient(cfg config.AIConfig) *AIClient {
	return &AIClient{
		apiKey:      cfg.APIKey,
		baseURL:     cfg.BaseURL,
		model:       cfg.Model,
		maxTokens:   cfg.MaxTokens,
		temperature: cfg.Temperature,
		httpClient:  &http.Client{Timeout: 120 * time.Second},
	}
}

// ChatMessage 对话消息
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatCompletionRequest 请求体
type ChatCompletionRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	Temperature float64       `json:"temperature,omitempty"`
}

// ChatCompletionResponse 响应体
type ChatCompletionResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error,omitempty"`
}

// ExtractRecipe 调用AI解析菜谱
func (c *AIClient) ExtractRecipe(ctx context.Context, input string) (*model.RecipeData, error) {
	prompt := c.buildSystemPrompt()

	reqBody := ChatCompletionRequest{
		Model:    c.model,
		Messages: []ChatMessage{
			{Role: "system", Content: prompt},
			{Role: "user", Content: input},
		},
		MaxTokens:   c.maxTokens,
		Temperature: c.temperature,
	}

	resp, err := c.callAPI(ctx, reqBody)
	if err != nil {
		return nil, fmt.Errorf("AI API调用失败: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("AI返回空结果")
	}

	content := resp.Choices[0].Message.Content

	var recipe model.RecipeData
	if err := xjson.TryParse(content, &recipe); err != nil {
		logger.ErrorLog("AI JSON解析失败", logger.String("content", content[:min(200, len(content))]), logger.Error(err))
		return nil, fmt.Errorf("AI输出JSON解析失败: %w", err)
	}

	logger.Info("AI解析成功", logger.String("dish", recipe.DishName), logger.Int("ingredients", len(recipe.Ingredients)))
	return &recipe, nil
}

// callAPI 调用AI API
func (c *AIClient) callAPI(ctx context.Context, reqBody ChatCompletionRequest) (*ChatCompletionResponse, error) {
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/chat/completions", bytes.NewReader(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("AI API返回错误状态码: %d, body: %s", resp.StatusCode, string(body))
	}

	var result ChatCompletionResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析AI响应失败: %w, body: %s", err, string(body))
	}

	if result.Error != nil {
		return nil, fmt.Errorf("AI API错误: %s", result.Error.Message)
	}

	return &result, nil
}

// buildSystemPrompt 构建系统提示词
func (c *AIClient) buildSystemPrompt() string {
	return `你是一个中餐菜谱结构化助手。请根据B站视频的字幕、简介和评论提取菜谱信息。

信息来源优先级（从高到低）：
1. 字幕内容最可信，包含操作步骤和顺序
2. 视频简介中的材料和配方很可信
3. 置顶评论可信，常包含文字版菜谱
4. 高赞评论只作为补充参考
5. 如果不同来源冲突，以字幕和简介为准

你必须严格遵守以下规则：
1. 只能基于字幕、简介和评论解析，不得编造关键材料
2. 评论里如果只是用户评价，不要当成菜谱步骤
3. "适量""少许"必须降低confidence到0.4以下
4. "一勺""两个""一把"需要按常见中餐标准换算成克或毫升
5. 如果无法判断菜名，可以从视频标题推断，但confidence降低
6. 如果无法判断份数，默认按1份，并在uncertain_items中标记
7. 如果无法提取有效材料和步骤，返回解析失败（ingredients和steps为空数组）
8. 食用油、盐、糖等调料如果原文未明确用量，估算后标记为低置信度
9. 输出source字段，说明来源：subtitle、description、top_comment、comment

常见单位换算参考：
- 1斤 = 500g
- 1两 = 50g
- 1勺生抽/酱油/醋/料酒/蚝油 ≈ 15ml
- 1勺油 ≈ 10g
- 1勺盐 ≈ 6g
- 1个鸡蛋 ≈ 50g
- 1个番茄 ≈ 150g
- 1个土豆 ≈ 150g
- 1瓣蒜 ≈ 5g
- 1块豆腐 ≈ 300g

请输出严格JSON，不要输出Markdown代码块或其他格式。

JSON格式如下：
{
  "dish_name": "菜名",
  "servings": 2,
  "ingredients": [
    {
      "name": "食材名称",
      "amount": 1,
      "unit": "斤",
      "grams": 500,
      "confidence": 0.95,
      "source": "subtitle",
      "source_text": "原文内容"
    }
  ],
  "steps": [
    {
      "order": 1,
      "title": "步骤标题",
      "description": "详细描述",
      "start_time": 0,
      "end_time": 0,
      "techniques": ["技法"],
      "duration_minutes": 5,
      "source": "subtitle"
    }
  ],
  "tips": ["烹饪技巧"],
  "uncertain_items": ["不确定项说明"],
  "confidence": 0.8
}`
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

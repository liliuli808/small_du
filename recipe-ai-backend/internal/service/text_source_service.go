package service

import (
	"fmt"
	"recipe-ai-backend/internal/model"
	"strings"
)

// TextSourceService 文本来源服务
type TextSourceService struct{}

// NewTextSourceService 创建文本来源服务
func NewTextSourceService() *TextSourceService {
	return &TextSourceService{}
}

// TextSourceBundle 文本来源包
type TextSourceBundle struct {
	Title       string
	Description string
	Subtitle    string
	TopComments []model.CommentItem
	HasSubtitle bool
}

// BuildAIInput 构建AI输入
func (s *TextSourceService) BuildAIInput(bundle TextSourceBundle) string {
	var b strings.Builder

	if bundle.Title != "" {
		b.WriteString("【视频标题】\n")
		b.WriteString(truncate(bundle.Title, 200))
		b.WriteString("\n\n")
	}

	if bundle.Description != "" {
		b.WriteString("【视频简介】\n")
		b.WriteString(truncate(bundle.Description, 2000))
		b.WriteString("\n\n")
	}

	if bundle.Subtitle != "" {
		b.WriteString("【字幕文本】\n")
		b.WriteString(truncate(bundle.Subtitle, 12000))
		b.WriteString("\n\n")
	}

	if len(bundle.TopComments) > 0 {
		b.WriteString("【评论补充】\n")
		for _, c := range bundle.TopComments {
			if c.IsTop {
				b.WriteString("置顶评论：")
			} else {
				b.WriteString(fmt.Sprintf("高赞评论，点赞%d：", c.Like))
			}
			b.WriteString(truncate(c.Message, 800))
			b.WriteString("\n")
		}
	}

	return b.String()
}

// HasUsefulText 检查是否有可用的文本内容
func (s *TextSourceService) HasUsefulText(bundle TextSourceBundle) bool {
	if bundle.Subtitle != "" {
		return true
	}
	if bundle.Description != "" {
		return true
	}
	if len(bundle.TopComments) > 0 {
		return true
	}
	return false
}

// GetSourceTypes 获取来源类型列表
func (s *TextSourceService) GetSourceTypes(bundle TextSourceBundle) []string {
	var types []string
	if bundle.HasSubtitle {
		types = append(types, "subtitle")
	}
	if bundle.Description != "" {
		types = append(types, "description")
	}
	for _, c := range bundle.TopComments {
		if c.IsTop {
			types = append(types, "top_comment")
		} else {
			types = append(types, "hot_comment")
		}
	}
	// 去重
	return uniqueStrings(types)
}

// truncate 截断文本
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen]) + "..."
}

func uniqueStrings(arr []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, s := range arr {
		if !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}
	return result
}

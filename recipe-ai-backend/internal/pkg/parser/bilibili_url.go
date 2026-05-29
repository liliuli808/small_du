package parser

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var bvidRegex = regexp.MustCompile(`BV[a-zA-Z0-9]+`)

// ParseBVID 从输入字符串中提取BV号
func ParseBVID(input string) (string, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", errors.New("输入为空")
	}

	// 先处理可能的短链
	if strings.Contains(input, "b23.tv") || strings.Contains(input, "bili.im") {
		resolved, err := ResolveShortURL(input)
		if err != nil {
			return "", fmt.Errorf("短链解析失败: %w", err)
		}
		input = resolved
	}

	bvid := bvidRegex.FindString(input)
	if bvid == "" {
		return "", errors.New("未识别到BV号")
	}

	return bvid, nil
}

// ResolveShortURL 解析B站短链
func ResolveShortURL(rawURL string) (string, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// 获取跳转后的URL
			return http.ErrUseLastResponse
		},
	}

	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 从Location头获取跳转URL
	if location := resp.Header.Get("Location"); location != "" {
		return location, nil
	}

	return "", errors.New("无法解析短链跳转")
}

// IsValidBilibiliURL 检查是否是有效的B站链接
func IsValidBilibiliURL(input string) bool {
	input = strings.ToLower(input)
	return strings.Contains(input, "bilibili.com") ||
		strings.Contains(input, "b23.tv") ||
		strings.Contains(input, "bili.im") ||
		strings.HasPrefix(input, "bv")
}

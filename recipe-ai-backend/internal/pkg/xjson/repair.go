package xjson

import (
	"encoding/json"
	"regexp"
	"strings"
)

// Repair 尝试修复常见JSON错误
func Repair(input string) string {
	// 去除markdown代码块标记
	input = removeMarkdownCodeBlock(input)

	// 去除前后空白
	input = strings.TrimSpace(input)

	// 确保以 { 开头
	if idx := strings.Index(input, "{"); idx > 0 {
		input = input[idx:]
	}

	// 确保以 } 结尾
	if idx := strings.LastIndex(input, "}"); idx > 0 && idx < len(input)-1 {
		input = input[:idx+1]
	}

	return input
}

// removeMarkdownCodeBlock 去除markdown代码块
func removeMarkdownCodeBlock(input string) string {
	// 匹配 ```json 或 ``` 开头
	re := regexp.MustCompile("```(?:json)?\\s*")
	input = re.ReplaceAllString(input, "")

	// 匹配 ``` 结尾
	input = strings.ReplaceAll(input, "```", "")

	return input
}

// TryParse 尝试解析JSON，失败则尝试修复后再解析
func TryParse(input string, v interface{}) error {
	// 先尝试直接解析
	if err := json.Unmarshal([]byte(input), v); err == nil {
		return nil
	}

	// 尝试修复后解析
	repaired := Repair(input)
	if err := json.Unmarshal([]byte(repaired), v); err == nil {
		return nil
	}

	// 尝试更激进的修复
	repaired = aggressiveRepair(repaired)
	return json.Unmarshal([]byte(repaired), v)
}

// aggressiveRepair 更激进的修复
func aggressiveRepair(input string) string {
	// 处理末尾多余的逗号
	re := regexp.MustCompile(`,(\s*[}\]])`)
	input = re.ReplaceAllString(input, "$1")

	// 处理单引号
	input = strings.ReplaceAll(input, "'", "\"")

	// 处理未转义的换行
	input = strings.ReplaceAll(input, "\n", "\\n")

	return input
}

// IsValidJSON 检查是否是有效的JSON
func IsValidJSON(input string) bool {
	var v interface{}
	return json.Unmarshal([]byte(input), &v) == nil
}

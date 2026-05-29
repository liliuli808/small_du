package unitconvert

import (
	"regexp"
	"strconv"
	"strings"
)

// UnitConversion 单位换算规则
type UnitConversion struct {
	Expression string
	Grams      float64
	Confidence float64 // 换算的置信度，用于标记估算值
}

// 常见中餐单位换算表
var unitConversions = []UnitConversion{
	{"1斤", 500, 0.95},
	{"1两", 50, 0.95},
	{"1克", 1, 1.0},
	{"1g", 1, 1.0},
	{"1kg", 1000, 1.0},
	{"1千克", 1000, 1.0},
	{"1毫升", 1, 0.95},
	{"1ml", 1, 0.95},
	{"1升", 1000, 0.95},
	{"1L", 1000, 0.95},
	{"1勺生抽", 15, 0.7},
	{"1勺酱油", 15, 0.7},
	{"1勺油", 10, 0.7},
	{"1勺醋", 15, 0.7},
	{"1勺料酒", 15, 0.7},
	{"1勺蚝油", 15, 0.7},
	{"1勺盐", 6, 0.7},
	{"1勺糖", 8, 0.7},
	{"1勺淀粉", 8, 0.7},
	{"1勺豆瓣酱", 15, 0.7},
	{"1个鸡蛋", 50, 0.9},
	{"1个番茄", 150, 0.7},
	{"1个土豆", 150, 0.7},
	{"1个洋葱", 200, 0.7},
	{"1个青椒", 80, 0.7},
	{"1个红椒", 80, 0.7},
	{"1瓣蒜", 5, 0.7},
	{"1根葱", 15, 0.7},
	{"1根胡萝卜", 100, 0.7},
	{"1块豆腐", 300, 0.7},
	{"1块姜", 10, 0.7},
	{"1片姜", 5, 0.7},
	{"1片香叶", 0.5, 0.5},
	{"1片桂皮", 2, 0.5},
	{"少许盐", 1, 0.4},
	{"少许糖", 2, 0.4},
	{"少许油", 5, 0.4},
	{"适量盐", 2, 0.4},
	{"适量糖", 5, 0.4},
	{"适量油", 10, 0.4},
	{"适量生抽", 10, 0.4},
	{"适量老抽", 5, 0.4},
	{"适量料酒", 10, 0.4},
	{"适量蚝油", 10, 0.4},
	{"适量淀粉", 10, 0.4},
	{"一把葱", 10, 0.4},
	{"一把香菜", 5, 0.4},
	{"一小撮盐", 1, 0.3},
	{"一小把", 10, 0.3},
	{"一撮", 2, 0.3},
	{"一小碗", 200, 0.6},
	{"一碗", 250, 0.6},
	{"一大碗", 400, 0.6},
	{"一小勺", 5, 0.7},
	{"一大勺", 20, 0.7},
	{"一茶匙", 5, 0.7},
	{"一汤匙", 15, 0.7},
	{"半勺", 5, 0.7},
}

// IngredientEstimate 食材估算结果
type IngredientEstimate struct {
	Grams      float64
	Confidence float64
}

// EstimateGrams 根据原始用量文本估算克重
func EstimateGrams(amount float64, unit, originalText string) IngredientEstimate {
	text := strings.ToLower(originalText)

	// 如果单位已经是克或毫升，直接返回
	lowerUnit := strings.ToLower(unit)
	if lowerUnit == "g" || lowerUnit == "克" || lowerUnit == "ml" || lowerUnit == "毫升" {
		return IngredientEstimate{Grams: amount, Confidence: 1.0}
	}

	// 如果是斤
	if lowerUnit == "斤" {
		return IngredientEstimate{Grams: amount * 500, Confidence: 0.95}
	}

	// 如果是两
	if lowerUnit == "两" {
		return IngredientEstimate{Grams: amount * 50, Confidence: 0.95}
	}

	// 如果是kg/千克
	if lowerUnit == "kg" || lowerUnit == "千克" {
		return IngredientEstimate{Grams: amount * 1000, Confidence: 1.0}
	}

	// 尝试从原始文本匹配单位换算
	for _, conv := range unitConversions {
		if strings.Contains(text, conv.Expression) || strings.Contains(text, strings.TrimPrefix(conv.Expression, "1")) {
			return IngredientEstimate{Grams: amount * conv.Grams, Confidence: conv.Confidence}
		}
	}

	// 默认情况，返回原始量值
	return IngredientEstimate{Grams: amount, Confidence: 0.5}
}

// ExtractAmountAndUnit 从文本中提取数量和单位
func ExtractAmountAndUnit(text string) (float64, string) {
	// 匹配数字+单位模式
	re := regexp.MustCompile(`(\d+(?:\.\d+)?)\s*([克g斤两kg千克毫升ml升L勺个瓣根块片碗茶匙汤匙撮把]+)`)
	matches := re.FindStringSubmatch(text)
	if len(matches) >= 3 {
		amount, _ := strconv.ParseFloat(matches[1], 64)
		return amount, matches[2]
	}

	// 匹配中文数字
	chineseNums := map[rune]float64{
		'一': 1, '二': 2, '两': 2, '三': 3, '四': 4,
		'五': 5, '六': 6, '七': 7, '八': 8, '九': 9,
		'十': 10, '半': 0.5, '几': 3,
	}

	for _, r := range text {
		if num, ok := chineseNums[r]; ok {
			// 尝试匹配单位
			re2 := regexp.MustCompile(string(r) + `\s*([克g斤两kg千克毫升ml升L勺个瓣根块片碗茶匙汤匙撮把]+)`)
			matches2 := re2.FindStringSubmatch(text)
			if len(matches2) >= 2 {
				return num, matches2[1]
			}
			return num, ""
		}
	}

	return 0, ""
}

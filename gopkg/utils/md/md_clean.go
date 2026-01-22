package md

import (
	"regexp"
	"strings"
)

// CleanString 去除字符串中的转义字符和特殊符号
func CleanString(input string) string {
	if input == "" {
		return ""
	}

	// 过滤掉 Markdown 图片链接，包括 s3 地址
	imageRegex := regexp.MustCompile(`!\[.*?\]\(.*?\)`)
	result := imageRegex.ReplaceAllString(input, "")

	// 去除常见转义字符
	result = strings.ReplaceAll(result, "\n", "")
	result = strings.ReplaceAll(result, "\r", "")
	result = strings.ReplaceAll(result, "\t", "")
	result = strings.ReplaceAll(result, "\v", "")
	result = strings.ReplaceAll(result, "\f", "")
	result = strings.ReplaceAll(result, "\b", "")
	result = strings.ReplaceAll(result, "\a", "")

	// 去除控制字符 (ASCII 0-31)
	controlCharsRegex := regexp.MustCompile("[\x00-\x1F]")
	result = controlCharsRegex.ReplaceAllString(result, "")

	// 去除特殊符号（可根据需要自定义）
	//specialCharsRegex := regexp.MustCompile("[\\\\~!@#$%^&*()+=|{}':;,.<>/?\\[\\]\"]+")
	//result = specialCharsRegex.ReplaceAllString(result, "")

	// 去除多余空格
	spaceRegex := regexp.MustCompile(`\s+`)
	result = spaceRegex.ReplaceAllString(result, " ")

	return strings.TrimSpace(result)
}

// CleanStringPreserveNewlines 去除字符串中的转义字符和特殊符号，但保留换行符
func CleanStringPreserveNewlines(input string) string {
	if input == "" {
		return ""
	}

	// 去除除换行符外的常见转义字符
	result := strings.ReplaceAll(input, "\r", "")
	result = strings.ReplaceAll(result, "\t", "")
	result = strings.ReplaceAll(result, "\v", "")
	result = strings.ReplaceAll(result, "\f", "")
	result = strings.ReplaceAll(result, "\b", "")
	result = strings.ReplaceAll(result, "\a", "")

	// 去除控制字符 (ASCII 0-31)，但保留换行符 (\n = ASCII 10)
	controlCharsRegex := regexp.MustCompile("[\x00-\x09\x0B-\x1F]")
	result = controlCharsRegex.ReplaceAllString(result, "")

	// 去除特殊符号（可根据需要自定义）
	specialCharsRegex := regexp.MustCompile("[\\\\~!@#$%^&*()+=|{}':;,.<>/?\\[\\]\"]+")
	result = specialCharsRegex.ReplaceAllString(result, "")

	// 去除多余空格，但保留换行符
	result = strings.ReplaceAll(result, "\n", "NEWLINE_PLACEHOLDER")
	spaceRegex := regexp.MustCompile(`\s+`)
	result = spaceRegex.ReplaceAllString(result, " ")
	result = strings.ReplaceAll(result, "NEWLINE_PLACEHOLDER", "\n")

	return strings.TrimSpace(result)
}

// CleanStringCustom 自定义去除字符串中的特定字符
func CleanStringCustom(input string, removeChars string) string {
	if input == "" {
		return ""
	}

	// 创建一个映射表，标记需要移除的字符
	removeMap := make(map[rune]bool)
	for _, char := range removeChars {
		removeMap[char] = true
	}

	// 过滤字符串
	var result strings.Builder
	for _, char := range input {
		if !removeMap[char] {
			result.WriteRune(char)
		}
	}

	return result.String()
}

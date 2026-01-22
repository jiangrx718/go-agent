package md

import (
	"regexp"
	"strings"
)

// Paragraph 表示 Markdown 中的一个段落
type Paragraph struct {
	Content string // 段落内容
	Line    int    // 段落在原文中的行号（从1开始）
}

// ExtractParagraphs 从 Markdown 文本中提取所有非标题的段落
// 会过滤掉图片和特定的 HTML 注释
func ExtractParagraphs(markdown string) []Paragraph {
	var paragraphs []Paragraph

	// 分割成行
	lines := strings.Split(markdown, "\n")

	// 定义正则表达式
	headingRegex := regexp.MustCompile(`^#{1,6}\s+.*$`)
	imageRegex := regexp.MustCompile(`!\[.*?\]\(.*?\)`)
	htmlCommentRegex := regexp.MustCompile(`<!--.*?-->`)
	htmlDivRegex := regexp.MustCompile(`<div.*?>.*?</div>`)

	var currentParagraph strings.Builder
	var paragraphStartLine int
	inParagraph := false

	for i, line := range lines {
		lineNum := i + 1
		trimmedLine := strings.TrimSpace(line)

		// 跳过空行
		if trimmedLine == "" {
			if inParagraph {
				// 结束当前段落
				paragraphs = append(paragraphs, Paragraph{
					Content: strings.TrimSpace(currentParagraph.String()),
					Line:    paragraphStartLine,
				})
				currentParagraph.Reset()
				inParagraph = false
			}
			continue
		}

		// 跳过标题行
		if headingRegex.MatchString(trimmedLine) {
			if inParagraph {
				// 结束当前段落
				paragraphs = append(paragraphs, Paragraph{
					Content: strings.TrimSpace(currentParagraph.String()),
					Line:    paragraphStartLine,
				})
				currentParagraph.Reset()
				inParagraph = false
			}
			continue
		}

		// 跳过图片行
		if imageRegex.MatchString(trimmedLine) {
			continue
		}

		// 跳过 HTML 注释和特定的 div
		if htmlCommentRegex.MatchString(trimmedLine) || htmlDivRegex.MatchString(trimmedLine) {
			continue
		}

		// 开始新段落或继续当前段落
		if !inParagraph {
			paragraphStartLine = lineNum
			inParagraph = true
		} else {
			currentParagraph.WriteString(" ")
		}

		currentParagraph.WriteString(trimmedLine)
	}

	// 处理最后一个段落
	if inParagraph {
		paragraphs = append(paragraphs, Paragraph{
			Content: strings.TrimSpace(currentParagraph.String()),
			Line:    paragraphStartLine,
		})
	}

	return paragraphs
}

// ExtractParagraphsAsText 从 Markdown 文本中提取所有非标题的段落，并以字符串数组形式返回
func ExtractParagraphsAsText(markdown string) []string {
	paragraphs := ExtractParagraphs(markdown)
	result := make([]string, len(paragraphs))

	for i, p := range paragraphs {
		result[i] = p.Content
	}

	return result
}

// ExtractParagraphsAsString 从 Markdown 文本中提取所有非标题的段落，并以单个字符串形式返回
// separator 参数指定段落之间的分隔符
func ExtractParagraphsAsString(markdown string, separator string) string {
	paragraphs := ExtractParagraphsAsText(markdown)
	return strings.Join(paragraphs, separator)
}

// ExtractParagraphsByMarker 根据内部分段标志提取段落
// 使用 <!--内部分段标志 --> 及其所在行的内容作为段落分隔标识
// 会过滤掉标题内容、内部使用注释和空段落
func ExtractParagraphsByMarker(markdown string) []Paragraph {
	var paragraphs []Paragraph

	// 定义分段标志的正则表达式 - 匹配包含<!--内部分段标志 -->的整行内容
	markerRegex := regexp.MustCompile(`(?m)^.*<!--内部分段标志\s*-->.*$`)

	// 定义标题的正则表达式
	headingRegex := regexp.MustCompile(`^#{1,6}\s+.*$`)

	// 定义内部使用注释的正则表达式
	internalCommentRegex := regexp.MustCompile(`<!--内部使用勿删\s*--><div.*?</div>`)

	// 将原文本中的分段标志替换为统一的分隔符，便于后续处理
	processedMarkdown := markerRegex.ReplaceAllString(markdown, "<<<PARAGRAPH_MARKER>>>")

	// 使用统一的分隔符分割文本
	segments := strings.Split(processedMarkdown, "<<<PARAGRAPH_MARKER>>>")

	// 处理每个分段
	lineOffset := 1
	for _, segment := range segments {
		if strings.TrimSpace(segment) == "" {
			// 跳过空段落
			// 计算段落中的行数以更新行号偏移
			lineOffset += len(strings.Split(segment, "\n"))
			continue
		}

		// 先过滤掉内部使用注释
		filteredSegment := internalCommentRegex.ReplaceAllString(segment, "")

		// 过滤掉标题行
		var filteredLines []string
		lines := strings.Split(filteredSegment, "\n")
		for _, line := range lines {
			trimmedLine := strings.TrimSpace(line)
			// 跳过标题行
			if !headingRegex.MatchString(trimmedLine) {
				filteredLines = append(filteredLines, line)
			}
		}

		// 如果过滤后没有内容，则跳过
		if len(filteredLines) == 0 {
			lineOffset += len(strings.Split(segment, "\n"))
			continue
		}

		// 获取段落内容并去除首尾空白
		paragraphContent := strings.TrimSpace(strings.Join(filteredLines, "\n"))

		// 如果段落内容为空，则跳过
		if paragraphContent == "" {
			lineOffset += len(strings.Split(segment, "\n"))
			continue
		}

		// 创建段落对象
		paragraphs = append(paragraphs, Paragraph{
			Content: paragraphContent,
			Line:    lineOffset,
		})

		// 更新行号偏移
		lineOffset += len(strings.Split(segment, "\n"))
	}

	return paragraphs
}

// ExtractParagraphsByMarkerAsText 根据内部分段标志提取段落，并以字符串数组形式返回
func ExtractParagraphsByMarkerAsText(markdown string) []string {
	paragraphs := ExtractParagraphsByMarker(markdown)
	result := make([]string, len(paragraphs))

	for i, p := range paragraphs {
		if p.Content == "" {
			continue
		}
		result[i] = p.Content
	}

	return result
}

// ExtractParagraphsByMarkerAsString 根据内部分段标志提取段落，并以单个字符串形式返回
func ExtractParagraphsByMarkerAsString(markdown string) string {
	paragraphs := ExtractParagraphsByMarkerAsText(markdown)
	return strings.Join(paragraphs, "\n\n")
}

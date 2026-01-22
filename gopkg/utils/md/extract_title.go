package md

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

// HeadingNode 表示标题节点
type HeadingNode struct {
	Title    string         `json:"title"`    // 标题文本
	Level    int            `json:"level"`    // 标题级别
	Children []*HeadingNode `json:"children"` // 子标题
}

// 预编译正则，避免每次调用重复编译
var (
	reAtxHeading    = regexp.MustCompile(`^(#{1,6})\s+(.+)$`)
	reRomanHeading  = regexp.MustCompile(`^([IVX]+)\.\s+(.+)$`)
	reLetterHeading = regexp.MustCompile(`^([A-Z])\.\s+(.+)$`)
	reAtxWithLetter = regexp.MustCompile(`^(#{1,6})\s+([A-Z])\.\s+(.+)$`)
	reAtxWithRoman  = regexp.MustCompile(`^(#{1,6})\s+([IVX]+)\.\s+(.+)$`)
)

// ExtractTOC 从Markdown内容中提取标题并构建树形结构
// maxDepth: 提取的最大标题深度，0表示不限制深度
func ExtractTOC(markdownContent string, maxDepth int) ([]*HeadingNode, error) {
	if markdownContent == "" {
		return nil, nil
	}

	// 提取所有标题
	headings, err := extractHeadings(markdownContent)
	if err != nil {
		return nil, err
	}

	// 构建树形结构
	return buildHeadingTree(headings, maxDepth), nil
}

// extractHeadings 提取所有标题
func extractHeadings(markdownContent string) ([]*HeadingNode, error) {
	var headings []*HeadingNode

	// 匹配Markdown标题的正则表达式
	// 支持以下格式:
	// 1. # 标题 (ATX风格)
	// 2. ## 标题 (ATX风格，多级)
	// 3. 标题 \n === (Setext风格，一级)
	// 4. 标题 \n --- (Setext风格，二级)
	// 5. 特殊格式: I. 标题, A. 标题 等
	// 6. 特殊格式: ### A. 标题 (ATX风格 + 字母标题)
	atxHeadingRegex := reAtxHeading
	romanHeadingRegex := reRomanHeading
	letterHeadingRegex := reLetterHeading
	atxWithLetterRegex := reAtxWithLetter
	atxWithRomanRegex := reAtxWithRoman

	scanner := bufio.NewScanner(strings.NewReader(markdownContent))
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 8*1024*1024)
	lineNum := 0
	var prevLine string

	for scanner.Scan() {
		line := scanner.Text()
		lineNum++

		// 处理ATX风格 + 字母标题 (### A. 标题)
		if matches := atxWithLetterRegex.FindStringSubmatch(line); len(matches) > 0 {
			level := len(matches[1]) // #的数量决定级别
			letter := matches[2]     // 字母部分
			title := strings.TrimSpace(matches[3])
			headings = append(headings, &HeadingNode{
				Title: fmt.Sprintf("%s. %s", letter, title),
				Level: level, // 使用#的数量作为级别
			})
			prevLine = line
			continue
		}

		// 处理ATX风格 + 罗马数字标题 (### I. 标题)
		if matches := atxWithRomanRegex.FindStringSubmatch(line); len(matches) > 0 {
			level := len(matches[1]) // #的数量决定级别
			roman := matches[2]      // 罗马数字部分
			title := strings.TrimSpace(matches[3])
			headings = append(headings, &HeadingNode{
				Title: fmt.Sprintf("%s. %s", roman, title),
				Level: level, // 使用#的数量作为级别
			})
			prevLine = line
			continue
		}

		// 处理ATX风格标题 (# 标题)
		if matches := atxHeadingRegex.FindStringSubmatch(line); len(matches) > 0 {
			level := len(matches[1]) // #的数量决定级别
			title := strings.TrimSpace(matches[2])
			headings = append(headings, &HeadingNode{
				Title: title,
				Level: level,
			})
			prevLine = line
			continue
		}

		// 处理罗马数字标题 (I. 标题, II. 标题等)
		if matches := romanHeadingRegex.FindStringSubmatch(line); len(matches) > 0 {
			title := strings.TrimSpace(matches[2])
			headings = append(headings, &HeadingNode{
				Title: fmt.Sprintf("%s. %s", matches[1], title),
				Level: 1, // 默认为一级标题
			})
			prevLine = line
			continue
		}

		// 处理字母标题 (A. 标题, B. 标题等)
		if matches := letterHeadingRegex.FindStringSubmatch(line); len(matches) > 0 {
			title := strings.TrimSpace(matches[2])
			headings = append(headings, &HeadingNode{
				Title: fmt.Sprintf("%s. %s", matches[1], title),
				Level: 2, // 默认为二级标题
			})
			prevLine = line
			continue
		}

		// 处理Setext风格标题 (标题\n===== 或 标题\n-----)
		if prevLine != "" {
			if strings.Trim(line, "=") == "" && strings.Contains(line, "=") {
				// 一级标题
				headings = append(headings, &HeadingNode{
					Title: strings.TrimSpace(prevLine),
					Level: 1,
				})
			} else if strings.Trim(line, "-") == "" && strings.Contains(line, "-") {
				// 二级标题
				headings = append(headings, &HeadingNode{
					Title: strings.TrimSpace(prevLine),
					Level: 2,
				})
			}
		}

		prevLine = line
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return headings, nil
}

// ExtractTOCReader 从 io.Reader 提取标题并构建树形结构，适用于大 Markdown 内容
func ExtractTOCReader(r io.Reader, maxDepth int) ([]*HeadingNode, error) {
	if r == nil {
		return nil, nil
	}

	headings, err := extractHeadingsFromReader(r)
	if err != nil {
		return nil, err
	}
	return buildHeadingTree(headings, maxDepth), nil
}

// extractHeadingsFromReader 按行从 Reader 提取所有标题
func extractHeadingsFromReader(r io.Reader) ([]*HeadingNode, error) {
	var headings []*HeadingNode

	atxHeadingRegex := reAtxHeading
	romanHeadingRegex := reRomanHeading
	letterHeadingRegex := reLetterHeading
	atxWithLetterRegex := reAtxWithLetter
	atxWithRomanRegex := reAtxWithRoman

	scanner := bufio.NewScanner(r)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 8*1024*1024)
	var prevLine string

	for scanner.Scan() {
		line := scanner.Text()

		if matches := atxWithLetterRegex.FindStringSubmatch(line); len(matches) > 0 {
			level := len(matches[1])
			letter := matches[2]
			title := strings.TrimSpace(matches[3])
			headings = append(headings, &HeadingNode{Title: fmt.Sprintf("%s. %s", letter, title), Level: level})
			prevLine = line
			continue
		}

		if matches := atxWithRomanRegex.FindStringSubmatch(line); len(matches) > 0 {
			level := len(matches[1])
			roman := matches[2]
			title := strings.TrimSpace(matches[3])
			headings = append(headings, &HeadingNode{Title: fmt.Sprintf("%s. %s", roman, title), Level: level})
			prevLine = line
			continue
		}

		if matches := atxHeadingRegex.FindStringSubmatch(line); len(matches) > 0 {
			level := len(matches[1])
			title := strings.TrimSpace(matches[2])
			headings = append(headings, &HeadingNode{Title: title, Level: level})
			prevLine = line
			continue
		}

		if matches := romanHeadingRegex.FindStringSubmatch(line); len(matches) > 0 {
			title := strings.TrimSpace(matches[2])
			headings = append(headings, &HeadingNode{Title: fmt.Sprintf("%s. %s", matches[1], title), Level: 1})
			prevLine = line
			continue
		}

		if matches := letterHeadingRegex.FindStringSubmatch(line); len(matches) > 0 {
			title := strings.TrimSpace(matches[2])
			headings = append(headings, &HeadingNode{Title: fmt.Sprintf("%s. %s", matches[1], title), Level: 2})
			prevLine = line
			continue
		}

		if prevLine != "" {
			if strings.Trim(line, "=") == "" && strings.Contains(line, "=") {
				headings = append(headings, &HeadingNode{Title: strings.TrimSpace(prevLine), Level: 1})
			} else if strings.Trim(line, "-") == "" && strings.Contains(line, "-") {
				headings = append(headings, &HeadingNode{Title: strings.TrimSpace(prevLine), Level: 2})
			}
		}

		prevLine = line
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return headings, nil
}

// buildHeadingTree 构建标题树
func buildHeadingTree(headings []*HeadingNode, maxDepth int) []*HeadingNode {
	if len(headings) == 0 {
		return nil
	}

	var root []*HeadingNode
	var stack []*HeadingNode

	for _, heading := range headings {
		// 如果设置了最大深度且当前标题级别超过最大深度，则跳过
		if maxDepth > 0 && heading.Level > maxDepth {
			continue
		}

		// 创建新节点
		node := &HeadingNode{
			Title: heading.Title,
			Level: heading.Level,
		}

		// 如果栈为空或当前标题级别小于等于栈顶标题级别，则作为根节点
		if len(stack) == 0 || node.Level <= stack[len(stack)-1].Level {
			// 清空栈并添加当前节点
			for len(stack) > 0 && node.Level <= stack[len(stack)-1].Level {
				stack = stack[:len(stack)-1]
			}

			if len(stack) == 0 {
				// 作为根节点
				root = append(root, node)
			} else {
				// 作为栈顶节点的子节点
				stack[len(stack)-1].Children = append(stack[len(stack)-1].Children, node)
			}
		} else {
			// 作为栈顶节点的子节点
			if len(stack) > 0 {
				stack[len(stack)-1].Children = append(stack[len(stack)-1].Children, node)
			} else {
				// 如果栈为空但级别大于1，则作为根节点
				root = append(root, node)
			}
		}

		// 将当前节点入栈
		stack = append(stack, node)
	}

	return root
}

// FormatTOC 将标题树格式化为字符串
func FormatTOC(headings []*HeadingNode) string {
	var sb strings.Builder
	formatTOCRecursive(&sb, headings, 0)
	return sb.String()
}

// formatTOCRecursive 递归格式化标题树
func formatTOCRecursive(sb *strings.Builder, headings []*HeadingNode, depth int) {
	for _, node := range headings {
		// 添加缩进
		sb.WriteString(strings.Repeat("- ", depth))
		sb.WriteString(node.Title)
		sb.WriteString("\n")

		if len(node.Children) > 0 {
			formatTOCRecursive(sb, node.Children, depth+1)
		}
	}
}

// FormatTOCWithPrefix 将标题树格式化为字符串，使用指定前缀
func FormatTOCWithPrefix(headings []*HeadingNode, prefix string) string {
	var sb strings.Builder
	formatTOCWithPrefixRecursive(&sb, headings, 0, prefix)
	return sb.String()
}

// formatTOCWithPrefixRecursive 递归格式化标题树，使用指定前缀
func formatTOCWithPrefixRecursive(sb *strings.Builder, headings []*HeadingNode, depth int, prefix string) {
	for _, node := range headings {
		// 添加缩进和前缀
		sb.WriteString(strings.Repeat(prefix, depth))
		sb.WriteString(" ")
		sb.WriteString(node.Title)
		sb.WriteString("\n")

		// 递归处理子节点
		if len(node.Children) > 0 {
			formatTOCWithPrefixRecursive(sb, node.Children, depth+1, prefix)
		}
	}
}

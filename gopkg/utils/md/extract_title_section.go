package md

import (
	"bufio"
	"regexp"
	"strings"
)

// Section 表示Markdown中的一个章节
type Section struct {
	Title   string     `json:"title"`   // 标题
	Content string     `json:"content"` // 内容
	Sub     []*Section `json:"sub"`     // 子章节
}

// ExtractSections 按标题层级提取Markdown内容
func ExtractSections(markdownContent string) ([]*Section, error) {
	// 解析所有标题及其位置信息
	headings, err := parseHeadingsWithPositions(markdownContent)
	if err != nil {
		return nil, err
	}

	// 构建标题树结构
	sections := buildSectionTree(headings, markdownContent)

	return sections, nil
}

// headingWithPosition 表示带位置信息的标题
type headingWithPosition struct {
	Title    string
	Level    int
	StartPos int // 标题开始位置
	EndPos   int // 标题结束位置
}

// parseHeadingsWithPositions 解析Markdown中的标题及其位置信息
func parseHeadingsWithPositions(markdownContent string) ([]*headingWithPosition, error) {
	var headings []*headingWithPosition

	// 定义匹配各种标题格式的正则表达式
	atxHeadingRegex := regexp.MustCompile(`^(#{1,6})\s+(.+)$`)                // # 标题
	setextHeading1Regex := regexp.MustCompile(`^(={3,})$`)                    // === (Setext一级标题线)
	setextHeading2Regex := regexp.MustCompile(`^(-{3,})$`)                    // --- (Setext二级标题线)
	romanHeadingRegex := regexp.MustCompile(`^([IVX]+)\.\s+(.+)$`)            // I. 标题
	letterHeadingRegex := regexp.MustCompile(`^([A-Z])\.\s+(.+)$`)            // A. 标题
	atxWithLetterRegex := regexp.MustCompile(`^(#{1,6})\s+([A-Z])\.\s+(.+)$`) // ### A. 标题
	atxWithRomanRegex := regexp.MustCompile(`^(#{1,6})\s+([IVX]+)\.\s+(.+)$`) // ### I. 标题

	scanner := bufio.NewScanner(strings.NewReader(markdownContent))
	lineNum := 0
	var prevLine string
	pos := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineNum++
		lineStartPos := pos

		// 处理ATX风格 + 字母标题 (### A. 标题)
		if matches := atxWithLetterRegex.FindStringSubmatch(line); len(matches) > 0 {
			level := len(matches[1]) // #的数量决定级别
			letter := matches[2]     // 字母部分
			title := strings.TrimSpace(matches[3])
			headings = append(headings, &headingWithPosition{
				Title:    strings.TrimSpace(letter + ". " + title),
				Level:    level,
				StartPos: lineStartPos,
				EndPos:   lineStartPos + len(line),
			})
			pos += len(line) + 1 // +1 for newline
			prevLine = line
			continue
		}

		// 处理ATX风格 + 罗马数字标题 (### I. 标题)
		if matches := atxWithRomanRegex.FindStringSubmatch(line); len(matches) > 0 {
			level := len(matches[1]) // #的数量决定级别
			roman := matches[2]      // 罗马数字部分
			title := strings.TrimSpace(matches[3])
			headings = append(headings, &headingWithPosition{
				Title:    strings.TrimSpace(roman + ". " + title),
				Level:    level,
				StartPos: lineStartPos,
				EndPos:   lineStartPos + len(line),
			})
			pos += len(line) + 1 // +1 for newline
			prevLine = line
			continue
		}

		// 处理ATX风格标题 (# 标题)
		if matches := atxHeadingRegex.FindStringSubmatch(line); len(matches) > 0 {
			level := len(matches[1]) // #的数量决定级别
			title := strings.TrimSpace(matches[2])
			headings = append(headings, &headingWithPosition{
				Title:    title,
				Level:    level,
				StartPos: lineStartPos,
				EndPos:   lineStartPos + len(line),
			})
			pos += len(line) + 1 // +1 for newline
			prevLine = line
			continue
		}

		// 处理罗马数字标题 (I. 标题, II. 标题等)
		if matches := romanHeadingRegex.FindStringSubmatch(line); len(matches) > 0 {
			title := strings.TrimSpace(matches[2])
			headings = append(headings, &headingWithPosition{
				Title:    strings.TrimSpace(matches[1] + ". " + title),
				Level:    1, // 默认为一级标题
				StartPos: lineStartPos,
				EndPos:   lineStartPos + len(line),
			})
			pos += len(line) + 1 // +1 for newline
			prevLine = line
			continue
		}

		// 处理字母标题 (A. 标题, B. 标题等)
		if matches := letterHeadingRegex.FindStringSubmatch(line); len(matches) > 0 {
			title := strings.TrimSpace(matches[2])
			headings = append(headings, &headingWithPosition{
				Title:    strings.TrimSpace(matches[1] + ". " + title),
				Level:    2, // 默认为二级标题
				StartPos: lineStartPos,
				EndPos:   lineStartPos + len(line),
			})
			pos += len(line) + 1 // +1 for newline
			prevLine = line
			continue
		}

		// 处理Setext风格标题 (标题\n===== 或 标题\n-----)
		if prevLine != "" {
			if setextHeading1Regex.MatchString(line) {
				// 一级标题
				headings = append(headings, &headingWithPosition{
					Title:    strings.TrimSpace(prevLine),
					Level:    1,
					StartPos: pos - len(prevLine) - 1, // -1 for newline
					EndPos:   lineStartPos + len(line),
				})
				pos += len(line) + 1 // +1 for newline
				prevLine = ""
				continue
			} else if setextHeading2Regex.MatchString(line) {
				// 二级标题
				headings = append(headings, &headingWithPosition{
					Title:    strings.TrimSpace(prevLine),
					Level:    2,
					StartPos: pos - len(prevLine) - 1, // -1 for newline
					EndPos:   lineStartPos + len(line),
				})
				pos += len(line) + 1 // +1 for newline
				prevLine = ""
				continue
			}
		}

		pos += len(line) + 1 // +1 for newline
		prevLine = line
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return headings, nil
}

// buildSectionTree 构建章节树
func buildSectionTree(headings []*headingWithPosition, markdownContent string) []*Section {
	if len(headings) == 0 {
		return nil
	}

	// 创建根节点列表
	var root []*Section

	// 创建节点列表
	nodes := make([]*Section, len(headings))
	for i, heading := range headings {
		nodes[i] = &Section{
			Title: heading.Title,
			Sub:   []*Section{},
		}
	}

	// 确定每个节点的内容
	for i, heading := range headings {
		// 内容开始位置是标题结束位置+1（换行符）
		startPos := heading.EndPos + 1

		// 内容结束位置是下一个标题的开始位置，如果没有下一个标题，则到文档末尾
		endPos := len(markdownContent)
		if i+1 < len(headings) {
			endPos = headings[i+1].StartPos
		}

		// 提取并清理内容
		if startPos < endPos && startPos < len(markdownContent) {
			content := markdownContent[startPos:endPos]
			// 移除首尾空白字符
			nodes[i].Content = strings.TrimSpace(content)
		}
	}

	// 构建树结构
	var stack []*headingWithPosition
	var nodeStack []*Section

	for i, heading := range headings {
		currentNode := nodes[i]

		// 维护栈，确保栈顶元素的级别小于当前元素
		for len(stack) > 0 && stack[len(stack)-1].Level >= heading.Level {
			stack = stack[:len(stack)-1]
			nodeStack = nodeStack[:len(nodeStack)-1]
		}

		if len(stack) == 0 {
			// 根节点
			root = append(root, currentNode)
		} else {
			// 子节点
			nodeStack[len(nodeStack)-1].Sub = append(nodeStack[len(nodeStack)-1].Sub, currentNode)
		}

		// 将当前节点压入栈
		stack = append(stack, heading)
		nodeStack = append(nodeStack, currentNode)
	}

	return root
}

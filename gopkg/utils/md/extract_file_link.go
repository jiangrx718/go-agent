package md

import (
	"go-agent/gopkg/utils/slices"
	"bufio"
	"io"
	"path"
	"regexp"
	"strings"
)

// 预编译的组合正则：一次扫描同时匹配 Markdown 链接/图片与 HTML img/a
// Markdown 语法支持可选的 title：![alt](url "title") 或 [text](url "title")
var reCombinedLinks = regexp.MustCompile(
	`(?:!\[[^\]]*\]|\[[^\]]*\])\(([^)\s]+)(?:\s+"[^"]*")?\)|<img[^>]+src=["']([^"']+)["']|<a[^>]+href=["']([^"']+)["']`,
)

// ExtractFileLinksByType 提取markdown中的所有文件链接，并按类型分组
func ExtractFileLinksByType(markdownContent string, isRemoveDuplicate bool) map[string][]string {
	var fileList []string

	// 空内容直接返回补齐的类型映射
	if markdownContent == "" {
		return ensureTypeKeys(map[string][]string{})
	}

	// 单次扫描提取所有链接（Markdown/HTML）
	matches := reCombinedLinks.FindAllStringSubmatch(markdownContent, -1)
	for _, m := range matches {
		// 三个捕获组之一会命中：1=markdown, 2=img src, 3=a href
		if len(m) >= 2 {
			url := firstNonEmpty(m[1], m[2], m[3])
			if url != "" {
				fileList = append(fileList, url)
			}
		}
	}

	if isRemoveDuplicate {
		fileList = slices.RemoveDuplicateElement(fileList)
	}

	// 按文件类型分组
	typeMap := map[string][]string{}
	// 分组
	for _, link := range fileList {
		// 计算扩展名前先去掉查询参数和锚点，避免 .png?x=1 被识别成错误扩展
		clean := stripQueryAndFragment(link)
		ext := strings.ToLower(path.Ext(clean))
		tp, ok := extTypeMap[ext]
		if !ok {
			tp = "other"
		}
		typeMap[tp] = append(typeMap[tp], link)
	}
	// 补齐所有类型键
	return ensureTypeKeys(typeMap)
}

// 辅助：返回第一个非空字符串
func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}

// 辅助：移除 URL 中的查询参数(?...)和片段(#...)
func stripQueryAndFragment(s string) string {
	if i := strings.IndexByte(s, '?'); i >= 0 {
		s = s[:i]
	}
	if i := strings.IndexByte(s, '#'); i >= 0 {
		s = s[:i]
	}
	return s
}

// 补齐类型键，保持返回结构稳定
func ensureTypeKeys(typeMap map[string][]string) map[string][]string {
	for _, tp := range []string{"image", "doc", "pdf", "excel", "ppt", "txt", "zip", "other"} {
		if _, ok := typeMap[tp]; !ok {
			typeMap[tp] = []string{}
		}
	}
	return typeMap
}

// ExtractFileLinksByTypeReader 从 io.Reader 流式提取文件链接并分组
// - 适用于大 Markdown 内容，避免一次性加载
// - 按行扫描，能覆盖常见的 Markdown/HTML 单行链接用法
// - isRemoveDuplicate：是否对提取到的链接去重
func ExtractFileLinksByTypeReader(r io.Reader, isRemoveDuplicate bool) map[string][]string {
	if r == nil {
		return ensureTypeKeys(map[string][]string{})
	}

	// 使用 Scanner 按行扫描，提升内存效率；加大缓冲以适配长行
	scanner := bufio.NewScanner(r)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 8*1024*1024) // 允许最长到 8MB 行

	var fileList []string
	for scanner.Scan() {
		line := scanner.Text()
		matches := reCombinedLinks.FindAllStringSubmatch(line, -1)
		for _, m := range matches {
			if len(m) >= 2 {
				url := firstNonEmpty(m[1], m[2], m[3])
				if url != "" {
					fileList = append(fileList, url)
				}

				if isRemoveDuplicate {
					fileList = slices.RemoveDuplicateElement(fileList)
				}
			}
		}
	}

	// 分组
	typeMap := map[string][]string{}
	for _, link := range fileList {
		clean := stripQueryAndFragment(link)
		ext := strings.ToLower(path.Ext(clean))
		tp, ok := extTypeMap[ext]
		if !ok {
			tp = "other"
		}
		typeMap[tp] = append(typeMap[tp], link)
	}
	return ensureTypeKeys(typeMap)
}

// extTypeMap 定义对应的类型
var extTypeMap = map[string]string{
	".jpg":      "image",
	".jpeg":     "image",
	".png":      "image",
	".gif":      "image",
	".bmp":      "image",
	".accountp": "image",
	".svg":      "image",
	".doc":      "doc",
	".docx":     "doc",
	".pdf":      "pdf",
	".xls":      "excel",
	".xlsx":     "excel",
	".csv":      "excel",
	".ppt":      "ppt",
	".pptx":     "ppt",
	".txt":      "txt",
	".zip":      "zip",
	".rar":      "zip",
}

package files

import (
	"archive/zip"
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// UrlExists 判断远程url是否存在
func UrlExists(url string) bool {
	resp, err := http.Head(url)
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true
	}
	return false
}

// DownloadFile 将远程文件下载到本地
func DownloadFile(url, path, filepath string) error {
	// 创建文件
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// 发送 GET 请求
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 检查 HTTP 响应状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// 将响应 Body 写入文件
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// ReadFileContent 使用 os 包中的 Open 函数打开文件，然后使用 io 包中的 Read 方法逐字节或指定大小读取文件内容。
func ReadFileContent(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer file.Close()

	content := make([]byte, 0)
	buf := make([]byte, 1024) // 指定读取的缓冲区大小

	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			log.Fatal(err)
			return nil, err
		}
		if n == 0 {
			break
		}
		content = append(content, buf[:n]...)
	}
	return content, nil
}

// ReadFileContentLineByLine 使用 os 包中的 Open 函数打开文件，然后使用 bufio 包中的 Scanner 对象逐行读取文件内容
func ReadFileContentLineByLine(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	content := make([]string, 0)
	for scanner.Scan() {
		content = append(content, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}
	return content, nil
}

// readFileContentOnce 使用 ioutil 包中的 ReadFile 函数一次性读取整个文件的内容
func readFileContentOnce(filename string) ([]byte, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return content, nil
}

// AppendToFile 将字符串内容追加到指定的文件
func AppendToFile(fileName, text string) error {
	// 打开文件，使用 os.O_APPEND 追加模式
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(text + "\n")
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	// 确保数据被刷新到文件中
	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("failed to flush writer: %w", err)
	}

	return nil
}

// IsDirectory 判断指定的目录是否存在
func IsDirectory(path string) (bool, error) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

func CreateFile(fileName string) {
	// 创建文件，使用 os.O_CREATE | os.O_WRONLY 标志
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("创建文件 %s 时出错: %v\n", fileName, err)
		return
	}
	defer file.Close()
}

func FileExists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func FilesExists(names []string) bool {
	for _, name := range names {
		ok, _ := FileExists(name)
		if ok {
			return ok
		}
	}
	return false
}

func DeleteFileIfExists(filePath string) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return
	}
	_ = os.Remove(filePath)
	return
}

func DeleteDirIfExists(dirPath string) error {
	// 检查目录是否存在
	_, err := os.Stat(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			// 目录不存在，无需操作
			return nil
		}
		// 其他错误（如权限问题）
		return err
	}

	// 删除目录及其内容
	err = os.RemoveAll(dirPath)
	if err != nil {
		return err
	}

	return nil
}

/**
 * 从路径中提取父目录名
 * @param path 文件路径
 * @eg:
 * 	/data/upload/year/20250415/123456_rain.xlsx -> 20250415
 * @return string
 */
func ExtractParentDir(path string) string {
	return filepath.Base(filepath.Dir(path))
}

/**
 * 判断指定文件是否为xlsx文件
 * @param path 文件路径
 * @return bool
 * @eg:
 * 	/data/upload/year/20250415/123456_rain.xlsx -> true
 * 	/data/upload/year/20250415/123456_rain.txt -> false
 */
func IsXLSXFile(path string) bool {
	// 获取文件扩展名并转换为小写
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".xlsx"
}

/**
 * 从路径中提取站点ID
 * @param path 文件路径
 * @return string
 */
func ExtractIdFromPath(path string) string {
	// 获取文件名
	filename := filepath.Base(path) // 例如：10810602.xlsx
	// 去掉扩展名
	name := strings.TrimSuffix(filename, filepath.Ext(filename)) // 结果为：10810602
	return name
}

/**
 * 从路径中提取站点ID
 * @param path 文件路径
 * @return string
 */
func ExtractIdFromPathId(path string) (name string) {
	// 获取文件名
	filename := ExtractIdFromPath(path) // 结果为：10800500_infer

	// 定义一个判断函数，当下划线出现时进行分割
	parts := strings.FieldsFunc(filename, func(r rune) bool {
		return r == '_' // 当字符是下划线时返回true，表示在此分割
	})

	if len(parts) > 0 {
		name = parts[0] //输出: 10800500
	}

	return name
}

/**
 * 解压zip文件到指定目录（不保留嵌套目录结构）
 * @param zipPath zip文件路径
 * @param destDir 目标目录
 * @return error
 */
func UnzipFlat(zipPath string, destDir string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("打开 zip 文件失败: %w", err)
	}
	defer r.Close()

	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("创建目标目录失败: %w", err)
	}

	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			continue // 忽略目录
		}

		// 始终将文件名放入目标目录根下（忽略嵌套路径）
		filename := filepath.Base(f.Name)
		targetPath := filepath.Join(destDir, filename)

		srcFile, err := f.Open()
		if err != nil {
			return fmt.Errorf("打开压缩文件项失败: %w", err)
		}
		defer srcFile.Close()

		destFile, err := os.Create(targetPath)
		if err != nil {
			return fmt.Errorf("创建目标文件失败: %w", err)
		}
		defer destFile.Close()

		if _, err := io.Copy(destFile, srcFile); err != nil {
			return fmt.Errorf("复制文件内容失败: %w", err)
		}
	}

	return nil
}

/**
 * 获取指定目录下的所有 .shp 文件
 * @param dirPath 目录路径
 * @return []string
 */
func GetShpFilesInDir(dir string, suffix string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("读取目录失败: %w", err)
	}

	var shpFiles []string
	for _, entry := range entries {
		name := entry.Name()
		if !entry.IsDir() &&
			strings.HasSuffix(name, suffix) &&
			!strings.HasPrefix(name, "._") {
			fullPath := filepath.Join(dir, name)
			shpFiles = append(shpFiles, fullPath)
		}
	}

	return shpFiles, nil
}

/**
 * 复制指定路径的文件到目标目录并重命名
 * srcPath 源文件路径（包含文件名）
 * dstDir 目标目录（只传目录名，不包含文件名）
 * newName：复制后的新文件名
 */
func CopyAndRenameFile(srcPath, dstDir, newName string) error {
	// 确保目标目录存在
	if err := os.MkdirAll(dstDir, os.ModePerm); err != nil {
		return fmt.Errorf("创建目标目录失败: %w", err)
	}

	// 打开源文件
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("打开源文件失败: %w", err)
	}
	defer srcFile.Close()

	// 构造目标文件路径
	dstPath := filepath.Join(dstDir, newName)

	// 创建目标文件
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("创建目标文件失败: %w", err)
	}
	defer dstFile.Close()

	// 复制内容
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("复制文件内容失败: %w", err)
	}

	return nil
}

// SplitFilePathTwoPart 将文件路径分割为目录和文件名
// 输入: 完整文件路径
// 输出: (目录路径, 文件名)
func SplitFilePathTwoPart(fullPath string) (dirPath, fileName string) {
	// 使用标准库的filepath包处理路径
	dirPath = filepath.Dir(fullPath)
	fileName = filepath.Base(fullPath)
	return
}

// FindLogDir 根据基准目录名查找实际可能延迟的日志目录
func FindLogDir(baseDir string, maxDelaySeconds int) string {
	// 解析基准目录名的时间
	baseTime, err := time.Parse("20060102_150405", filepath.Base(baseDir))
	if err != nil {
		return "" // 格式错误直接返回空
	}

	// 检查基准目录及其延迟版本
	for i := 0; i <= maxDelaySeconds; i++ {
		targetTime := baseTime.Add(time.Duration(i) * time.Second)
		targetDir := targetTime.Format("20060102_150405")
		path := filepath.Join(filepath.Dir(baseDir), targetDir)

		if _, err := os.Stat(path); err == nil {
			return path // 找到存在的目录
		}
	}

	return "" // 全部尝试后未找到
}

// FindLogFile 查找可能延迟的日志文件
// 文件格式示例: /Users/xxx/nc_path/out/jjmal1753068991.nc_20250721_151513.log
func FindLogFile(baseFile string, maxDelaySeconds int) string {
	// 1. 首先检查原始文件是否存在
	if _, err := os.Stat(baseFile); err == nil {
		return baseFile
	}

	// 2. 提取文件名中的时间部分
	fileName := filepath.Base(baseFile) // "jjmal1753068991.nc_20250721_151513.log"

	// 分割文件名获取时间部分
	parts := strings.Split(fileName, "_")
	if len(parts) < 3 {
		return "" // 文件名格式不符合预期
	}

	timePart := parts[len(parts)-2] + "_" + strings.TrimSuffix(parts[len(parts)-1], ".log") // "20250721_151513"

	// 3. 解析时间
	baseTime, err := time.Parse("20060102_150405", timePart)
	if err != nil {
		return "" // 时间格式错误
	}

	// 4. 构建基础文件名前缀
	prefix := strings.Join(parts[:len(parts)-2], "_") // "jjmal1753068991.nc"
	dir := filepath.Dir(baseFile)

	// 5. 检查延迟版本的文件
	for i := 1; i <= maxDelaySeconds; i++ {
		targetTime := baseTime.Add(time.Duration(i) * time.Second)
		targetFile := fmt.Sprintf("%s_%s.log", prefix, targetTime.Format("20060102_150405"))
		path := filepath.Join(dir, targetFile)

		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return "" // 未找到任何匹配文件
}

// SaveJSONToFile 将文件内容写入json文件
func SaveJSONToFile(data interface{}, filename string) error {
	// 检查目录是否存在，不存在则创建
	dir := filepath.Dir(filename)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return fmt.Errorf("创建目录失败: %v", err)
		}
	}

	// 序列化数据
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("JSON序列化失败: %v", err)
	}

	// 写入文件
	err = ioutil.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}

	return nil
}

package files

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

// ExtName 根据path文件路径获取无前缀(.)的文件名称
func ExtName(path string) string {
	ext := filepath.Ext(path)
	if ext != "" && ext[0] == '.' {
		ext = ext[1:]
	}
	return ext
}

// Split 切片
// return dir 路径
// return prefix 无后缀文件名称
// return ext 文件后缀
func Split(path string) (string, string, string) {
	if path == "" {
		return "", "", ""
	}

	dir, fileName := filepath.Split(path)
	extLastIndex := strings.LastIndex(fileName, ".")
	prefix := ""
	if extLastIndex == -1 {
		prefix = fileName
	} else {
		prefix = fileName[:extLastIndex]
	}

	return dir, prefix, ExtName(fileName)
}

// IsExist 判断path文件路径是否存在
func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// CopyFile 复制sourcePath路径文件到targetPath路径下
func CopyFile(sourcePath, targetPath string) error {
	srcFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("无法打开源文件: %w", err)
	}
	defer func() {
		_ = srcFile.Close()
	}()

	destFile, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("无法创建目标文件: %w", err)
	}
	defer func() {
		_ = destFile.Close()
	}()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return fmt.Errorf("文件拷贝失败: %w", err)
	}

	return nil
}

// ReplaceExt 替换文件路径path的后缀为ext
func ReplaceExt(path, ext string) string {
	return strings.ReplaceAll(path, filepath.Ext(path), ext)
}

// Write 将content写入指定文件(覆盖)
func Write(filePath string, content []byte) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}

	defer func() {
		_ = file.Close()
	}()

	if _, err = file.Write(content); err != nil {
		return err
	}

	return nil
}

func WriteReader(filePath string, reader io.Reader) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}

	defer func() {
		_ = file.Close()
	}()

	readerBytes, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	if _, err = file.Write(readerBytes); err != nil {
		return err
	}

	return nil
}

func WriteJSON(filePath string, a any) error {
	aJSON, err := jsoniter.Marshal(a)
	if err != nil {
		return err
	}

	return Write(filePath, aJSON)
}

// Read 读取指定路径文件
func Read(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = file.Close()
	}()

	return io.ReadAll(file)
}

func ReadJSON[T any](filePath string) (T, error) {
	var t T
	bytes, err := Read(filePath)
	if err != nil {
		return t, err
	}

	if err = jsoniter.Unmarshal(bytes, &t); err != nil {
		return t, err
	}

	return t, err
}

func Create(filePath string) error {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: error=%v, filepath=%s", err, filePath)
	}

	defer func() {
		_ = file.Close()
	}()

	return nil
}

func CreateFromBytes(filePath string, bytes []byte) error {
	if err := Create(filePath); err != nil {
		return err
	}

	return Write(filePath, bytes)
}

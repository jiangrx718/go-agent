package files

import (
	"go-agent/gopkg/utils"
	"fmt"
	"path/filepath"
	"strings"
)

type FileNames []FileName
type FileName string

func (f FileName) String() string {
	return string(f)
}

func (f FileName) Ext() string {
	ext := strings.ReplaceAll(filepath.Ext(f.String()), ".", "")
	return ext
}

func (f FileName) BaseExt() string {
	return filepath.Ext(f.String())
}

func (f FileName) LowerBaseExt() string {
	return strings.ToLower(f.BaseExt())
}

// Name 获取没有后缀的文件名称
func (f FileName) Name() string {
	base := filepath.Base(f.String())
	bases := strings.Split(base, ".")
	if len(bases) <= 1 {
		return base
	}

	bases = bases[:len(bases)-1]
	return strings.Join(bases, ".")
}

func (f FileName) Base() string {
	return filepath.Base(f.String())
}

func (f FileName) RemoveExt() string {
	return strings.ReplaceAll(f.String(), filepath.Ext(f.String()), "")
}

func (f FileName) GenSnowflakeFileName() FileName {
	return FileName(fmt.Sprintf("%s%s", utils.SnowflakeGenUUID(), f.Ext()))
}

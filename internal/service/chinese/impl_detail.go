package chinese

import (
	"context"
	"go-agent/gopkg/services"
	"strings"
	"unicode"

	"github.com/mozillazg/go-pinyin"
)

func (s *Service) Detail(ctx context.Context, chinese string) (services.Result, error) {
	var (
		result       strings.Builder
		args         = pinyin.NewArgs()
		chineseSlice []string
	)

	runes := []rune(chinese)
	for i, r := range runes {
		if !unicode.Is(unicode.Han, r) {
			result.WriteRune(r)
			continue
		}

		py := pinyin.SinglePinyin(r, args)
		if len(py) == 0 || len(py[0]) == 0 {
			continue
		}

		// 将当前汉字的拼音加入切片
		chineseSlice = append(chineseSlice, py[0])
		result.WriteString(py[0])
		if i < len(runes)-1 {
			result.WriteString("_")
		}
	}

	// 返回拼音切片结构
	return services.Success(ctx, chineseSlice)
}

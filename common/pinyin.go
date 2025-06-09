package common

import (
	"github.com/mozillazg/go-pinyin"
	"strings"
)

func GetFirstLetter(value string) string {
	// 默认
	py := pinyin.NewArgs()
	py.Style = pinyin.FirstLetter
	py.Fallback = func(r rune, a pinyin.Args) []string {
		return []string{string(r)}
	}
	retText := ""
	result := pinyin.Pinyin(value, py)
	for _, res := range result {
		for _, s := range res {
			retText += strings.ToLower(s)
			break
		}
	}
	return retText
}

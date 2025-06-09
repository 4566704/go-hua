// 通用函数

package common

import (
	"fmt"
)

func FormatFlow(value int64) string {
	if value < 1024 {
		return fmt.Sprintf("%d B", value)
	} else if value < 1024*1024 {
		return fmt.Sprintf("%.1f KB", float64(value)/float64(1024))
	} else if value < 1024*1024*1024 {
		return fmt.Sprintf("%.1f MB", float64(value)/float64(1024*1024))
	} else if value < 1024*1024*1024*1024 {
		return fmt.Sprintf("%.1f GB", float64(value)/float64(1024*1024*1024))
	}
	return fmt.Sprintf("%.2f TB", float64(value)/float64(1024*1024*1024*1024))
}

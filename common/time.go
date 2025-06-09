package common
import (
	"time"
)

// NtToUtc Windows NT timestamp to Unix UTC timestamp
func NtToUtc(ntTimestamp uint64) uint64 {
	unixTimestamp := (ntTimestamp - 11644473600000000) / 1000000
	return unixTimestamp
}

// UtcToNt Unix UTC timestamp to Windows NT timestamp
func UtcToNt(unixTime uint64) uint64 {
	ntTimestamp := unixTime*1000000 + 11644473600000000
	return ntTimestamp
}

// 时间字符串转时间戳
func TimeStrToTimestamp(timeStr string) (int64, error) {
	layout := "2006-01-02 15:04:05"
	t, err := time.ParseInLocation(layout, timeStr, time.Local)
	if err != nil {
		return 0, err
	}
	// 获取时间戳（单位是秒），并转换为毫秒
	millisecondTimestamp := t.UnixMilli()
	return millisecondTimestamp, nil
}

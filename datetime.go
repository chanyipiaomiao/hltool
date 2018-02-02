package hltool

import "time"

const (
	// Oneday 一天
	Oneday = 24 * time.Hour
)

// GetNowTime 获取当前时间
func GetNowTime() string {
	return time.Now().Format("20060102150405")
}

// GetNowTime2 获取当前时间
func GetNowTime2() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// GetNowTimeStamp 获取当前的时间戳
func GetNowTimeStamp() int64 {
	return time.Now().Unix()
}

package hltool

import (
	"time"
)

// HLogger 定义
type HLogger struct {
	LogPath      string
	FileName     string
	MaxAge       time.Duration
	RotationTime time.Duration
}

// NewHLogger 返回HLogger对象
func NewHLogger(logpath, filename string) *HLogger {
	return &HLogger{
		LogPath:      logpath,
		FileName:     filename,
		MaxAge:       15 * Oneday, // 默认保留15天日志
		RotationTime: Oneday,      // 默认24小时轮转一次日志
	}
}

// SetMaxAge 设置最大保留时间
// 单位: 天
func (hl *HLogger) SetMaxAge(day time.Duration) {
	hl.MaxAge = day * Oneday
}

// SetRotationTime 设置日志多久轮转一次
func (hl *HLogger) SetRotationTime(day time.Duration) {
	hl.RotationTime = day * Oneday
}

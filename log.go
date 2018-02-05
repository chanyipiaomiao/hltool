package hltool

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
)

// HLog 定义
type HLog struct {
	LogPath      string
	FileName     string
	LogType      string
	DateFormat   string
	LogLevel     log.Level
	MaxAge       time.Duration
	RotationTime time.Duration
}

// NewHLog 返回HLog对象
func NewHLog(logpath, filename, logType string) *HLog {
	return &HLog{
		LogPath:      logpath,
		FileName:     filename,
		LogType:      logType,
		LogLevel:     log.InfoLevel,
		DateFormat:   "%Y-%m-%d",
		MaxAge:       15 * Oneday, // 默认保留15天日志
		RotationTime: Oneday,      // 默认24小时轮转一次日志
	}
}

// SetLogType 设置日志格式 json|text
func (hl *HLog) SetLogType(logType string) {
	hl.LogType = logType
}

// SetMaxAge 设置最大保留时间
// 单位: 天
func (hl *HLog) SetMaxAge(day time.Duration) {
	hl.MaxAge = day * Oneday
}

// SetRotationTime 设置日志多久轮转一次
// 单位: 天
func (hl *HLog) SetRotationTime(day time.Duration) {
	hl.RotationTime = day * Oneday
}

// SetLevel 设置log level
// debug|info|warn|error|fatal|panic
func (hl *HLog) SetLevel(level string) {
	switch strings.ToLower(level) {
	case "panic":
		hl.LogLevel = log.PanicLevel
	case "fatal":
		hl.LogLevel = log.FatalLevel
	case "error":
		hl.LogLevel = log.ErrorLevel
	case "warn", "warning":
		hl.LogLevel = log.WarnLevel
	case "info":
		hl.LogLevel = log.InfoLevel
	default:
		hl.LogLevel = log.DebugLevel
	}
}

// SetDateFormat 设置日期格式
// format "%Y-%m-%d" | "%Y%m%d"
func (hl *HLog) SetDateFormat(format string) {
	hl.DateFormat = format
}

// setNull 相当于/dev/null
func setNull() *bufio.Writer {
	src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("err", err)
	}
	return bufio.NewWriter(src)
}

// GetLogger getlogger
func (hl *HLog) GetLogger() *log.Logger {

	logger := log.New()

	switch hl.LogType {
	case "text":
		logger.Formatter = &log.TextFormatter{}
	default:
		logger.Formatter = &log.JSONFormatter{}
	}

	logger.Level = hl.LogLevel

	filename := path.Join(hl.LogPath, hl.FileName)
	maxage := rotatelogs.WithMaxAge(hl.MaxAge)
	rotate := rotatelogs.WithRotationTime(hl.RotationTime)

	debugFileName := filename + ".debug"
	debugWriter, err := rotatelogs.New(
		fmt.Sprintf("%s.%s", debugFileName, hl.DateFormat),
		rotatelogs.WithLinkName(debugFileName),
		maxage,
		rotate,
	)

	infoFileName := filename + ".info"
	infoWriter, err := rotatelogs.New(
		fmt.Sprintf("%s.%s", infoFileName, hl.DateFormat),
		rotatelogs.WithLinkName(infoFileName),
		maxage,
		rotate,
	)

	warningFileName := filename + ".warn"
	warningWriter, err := rotatelogs.New(
		fmt.Sprintf("%s.%s", warningFileName, hl.DateFormat),
		rotatelogs.WithLinkName(warningFileName),
		maxage,
		rotate,
	)

	errorFileName := filename + ".error"
	errorWriter, err := rotatelogs.New(
		fmt.Sprintf("%s.%s", errorFileName, hl.DateFormat),
		rotatelogs.WithLinkName(errorFileName),
		maxage,
		rotate,
	)

	if err != nil {
		panic("error")
	}

	fileHook := lfshook.NewHook(lfshook.WriterMap{
		log.DebugLevel: debugWriter, // 为不同级别设置不同的输出目的
		log.InfoLevel:  infoWriter,
		log.WarnLevel:  warningWriter,
		log.ErrorLevel: errorWriter,
		log.FatalLevel: errorWriter,
		log.PanicLevel: errorWriter,
	}, &log.JSONFormatter{})

	logger.Hooks.Add(fileHook)

	if hl.LogLevel != log.DebugLevel {
		logger.Out = setNull()
	}

	return logger

}

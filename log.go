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

	// log 路径
	LogPath string

	// log文件名称
	FileName string

	// 日志类型  json|text
	LogType string

	// 文件名的日期格式 默认: %Y-%m-%d|%Y%m%d
	FileNameDateFormat string

	// 是否分离不同级别的日志 默认: true
	IsSeparateLevelLog bool

	// 日志条目中的公共字段
	CommonFields map[string]interface{}

	// 日志级别 默认: log.InfoLevel
	LogLevel log.Level

	// 日志最长保存多久 默认: 15天
	MaxAge time.Duration

	// 日志默认多长时间轮转一次 默认: 24小时
	RotationTime time.Duration
}

// NewHLog 返回HLog对象 和 error 目录创建失败
func NewHLog(logpath, filename string) (*HLog, error) {
	if !IsExist(logpath) {
		err := os.MkdirAll(logpath, os.ModePerm)
		if err != nil {
			return nil, fmt.Errorf("create <%s> error: %s", logpath, err)
		}
	}
	return &HLog{
		LogPath:            logpath,
		FileName:           filename,
		LogType:            "json",
		LogLevel:           log.InfoLevel,
		FileNameDateFormat: "%Y-%m-%d",
		IsSeparateLevelLog: true,
		MaxAge:             15 * Oneday,
		RotationTime:       Oneday,
	}, nil
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
	hl.FileNameDateFormat = format
}

// SetCommonFields 设置公共字段
func (hl *HLog) SetCommonFields(fields map[string]interface{}) {
	hl.CommonFields = fields
}

// SetSeparateLevelLog 设置是否分离不同级别的日志到不同的文件
func (hl *HLog) SetSeparateLevelLog(yes bool) {
	hl.IsSeparateLevelLog = yes
}

// setNull 相当于/dev/null
func setNull() *bufio.Writer {
	src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return nil
	}
	return bufio.NewWriter(src)
}

// GetLogger getlogger
func (hl *HLog) GetLogger() (*log.Entry, error) {

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

	if hl.IsSeparateLevelLog {
		debugFileName := filename + ".debug"
		debugWriter, err := rotatelogs.New(
			fmt.Sprintf("%s.%s", debugFileName, hl.FileNameDateFormat),
			rotatelogs.WithLinkName(debugFileName),
			maxage,
			rotate,
		)
		if err != nil {
			return nil, err
		}

		infoFileName := filename + ".info"
		infoWriter, err := rotatelogs.New(
			fmt.Sprintf("%s.%s", infoFileName, hl.FileNameDateFormat),
			rotatelogs.WithLinkName(infoFileName),
			maxage,
			rotate,
		)
		if err != nil {
			return nil, err
		}

		warningFileName := filename + ".warn"
		warningWriter, err := rotatelogs.New(
			fmt.Sprintf("%s.%s", warningFileName, hl.FileNameDateFormat),
			rotatelogs.WithLinkName(warningFileName),
			maxage,
			rotate,
		)
		if err != nil {
			return nil, err
		}

		errorFileName := filename + ".error"
		errorWriter, err := rotatelogs.New(
			fmt.Sprintf("%s.%s", errorFileName, hl.FileNameDateFormat),
			rotatelogs.WithLinkName(errorFileName),
			maxage,
			rotate,
		)

		if err != nil {
			return nil, err
		}

		// 文件 hook, 不同的级别 设置输出不同的文件
		fileHook := lfshook.NewHook(lfshook.WriterMap{
			log.DebugLevel: debugWriter, // 为不同级别设置不同的输出目的
			log.InfoLevel:  infoWriter,
			log.WarnLevel:  warningWriter,
			log.ErrorLevel: errorWriter,
			log.FatalLevel: errorWriter,
			log.PanicLevel: errorWriter,
		}, logger.Formatter)

		logger.Hooks.Add(fileHook)

	} else {
		writer, err := rotatelogs.New(
			fmt.Sprintf("%s.%s", filename, hl.FileNameDateFormat),
			maxage,
			rotate,
		)
		if err != nil {
			return nil, err
		}

		fileHook := lfshook.NewHook(lfshook.WriterMap{
			log.DebugLevel: writer,
			log.InfoLevel:  writer,
			log.WarnLevel:  writer,
			log.ErrorLevel: writer,
			log.FatalLevel: writer,
			log.PanicLevel: writer,
		}, logger.Formatter)

		logger.Hooks.Add(fileHook)
	}

	if hl.LogLevel != log.DebugLevel {
		if out := setNull(); out != nil {
			logger.Out = setNull()
		} else {
			logger.Out = os.Stdout
		}
	}

	loggerEntry := logger.WithFields(log.Fields(hl.CommonFields))

	return loggerEntry, nil

}

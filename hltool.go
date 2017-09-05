package hltool

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"regexp"
	"time"
)

// GetNowTime 获取当前时间
func GetNowTime() string {
	t := time.Now()
	return fmt.Sprintf("%d%d%d%d%d%d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}

// IsNumber 检查输入的字符串是否匹配数字
func IsNumber(input string) bool {

	if isNumber, _ := regexp.Match("^\\d+$", []byte(input)); isNumber {
		return true
	}
	return false
}

// CurrentUser 获取当前SSH连接的用户
func CurrentUser() string {
	currentUser, _ := user.Current()
	return currentUser.Username
}

// CurrentDir 获取当前路径
func CurrentDir() string {
	currentDir, _ := os.Getwd()
	return currentDir
}

// UserHome 获取用户的家目录
func UserHome() (string, error) {
	currentUser, err := user.Current()
	if err == nil {
		return currentUser.HomeDir, nil
	}

	home := os.Getenv("HOME")
	if home == "" {
		log.Fatal("User < HOME > Env not Found")
	}
	return home, nil
}

// IsFileExist 文件或目录是否存在
func IsFileExist(name string) bool {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return true
	}
	return false
}

// CheckError 错误检查
func CheckError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s error: %v", msg, err)
	}
}

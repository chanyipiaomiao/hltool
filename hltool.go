package hltool

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"regexp"
	"time"
	"math/rand"
	"encoding/hex"
	"crypto/md5"
	"golang.org/x/crypto/scrypt"
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
func GetNowTimeStamp() int64{ 
	return time.Now().Unix()
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

// IsExist 文件或目录是否存在
// return false 表示文件不存在
func IsExist(name string) bool {
	_, err := os.Stat(name)
	return err == nil || os.IsExist(err)
}

// CheckError 错误检查
func CheckError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s error: %v", msg, err)
	}
}

// CryptPassword 加密密码
func CryptPassword(password, salt string) string {
	dk, _ := scrypt.Key([]byte(password), []byte(salt), 16384, 8, 1, 32)
	return fmt.Sprintf("%x", dk)
}

// GetRandomString 生成随机字符串
func GetRandomString(length int) string{
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
	   result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
 }

// GetMD5 生成32位MD5
func GetMD5(text string) string{
	ctx := md5.New()
	ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}
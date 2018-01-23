package hltool

import (
	"fmt"
	"math/rand"
	"time"
	"encoding/hex"
	"crypto/md5"
	"golang.org/x/crypto/scrypt"
)

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

 // CryptPassword 加密密码
func CryptPassword(password, salt string) string {
	dk, _ := scrypt.Key([]byte(password), []byte(salt), 16384, 8, 1, 32)
	return fmt.Sprintf("%x", dk)
}


// GetMD5 生成32位MD5
func GetMD5(text string) string{
	ctx := md5.New()
	ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}
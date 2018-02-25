package hltool

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/scrypt"
)

// GenRandomString 生成随机字符串
// length 生成长度
// specialChar 是否生成特殊字符
func GenRandomString(length int, specialChar string) string {

	letterBytes := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	special := "!@#%$*.="

	if specialChar == "yes" {
		letterBytes = letterBytes + special
	}

	chars := []byte(letterBytes)

	if length == 0 {
		return ""
	}

	clen := len(chars)
	maxrb := 255 - (256 % clen)
	b := make([]byte, length)
	r := make([]byte, length+(length/4)) // storage for random bytes.
	i := 0
	for {
		if _, err := rand.Read(r); err != nil {
			return ""
		}
		for _, rb := range r {
			c := int(rb)
			if c > maxrb {
				continue // Skip this number to avoid modulo bias.
			}
			b[i] = chars[c%clen]
			i++
			if i == length {
				return string(b)
			}
		}
	}
}

// CryptPassword 加密密码
func CryptPassword(password, salt string) string {
	dk, _ := scrypt.Key([]byte(password), []byte(salt), 16384, 8, 1, 32)
	return fmt.Sprintf("%x", dk)
}

// GetMD5 生成32位MD5
func GetMD5(text string) string {
	ctx := md5.New()
	ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}

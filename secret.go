package hltool

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"

	"golang.org/x/crypto/scrypt"
)

// var (
// 	chars   = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
// 	special = "!@#%$*.="
// )

// GenRandomString 生成随机字符串
// chars  指定的字符串
// length 长度
func GenRandomString(length int, chars []byte) string {

	if length == 0 {
		return ""
	}

	clen := len(chars)
	if clen < 2 || clen > 256 {
		panic("Wrong charset length for NewLenChars()")
	}

	maxrb := 255 - (256 % clen)
	b := make([]byte, length)
	r := make([]byte, length+(length/4)) // storage for random bytes.
	i := 0

	for {
		if _, err := rand.Read(r); err != nil {
			panic("Error reading random bytes: " + err.Error())
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

package hltool

// 模拟 Google Authenticator 验证器
// https://github.com/robbiev/two-factor-auth

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base32"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

func toBytes(value int64) []byte {
	var result []byte
	mask := int64(0xFF)
	shifts := [8]uint16{56, 48, 40, 32, 24, 16, 8, 0}
	for _, shift := range shifts {
		result = append(result, byte((value>>shift)&mask))
	}
	return result
}

func toUint32(bytes []byte) uint32 {
	return (uint32(bytes[0]) << 24) + (uint32(bytes[1]) << 16) +
		(uint32(bytes[2]) << 8) + uint32(bytes[3])
}

// checkSecret 检查输入的secret的长度是否符合base32编码规则
func checkSecret(secret string) string {
	length := len(secret)
	if length%8 == 0 {
		return secret
	}
	n := length/8*8 + 8 - length
	return secret + strings.Repeat("=", n)
}

func convertSecret(secret string) ([]byte, error) {
	inputNoSpaces := strings.Replace(secret, " ", "", -1)
	decodeKey, err := base32.StdEncoding.DecodeString(checkSecret(strings.ToUpper(inputNoSpaces)))
	if err != nil {
		return nil, err
	}
	return decodeKey, nil
}

func oneTimePassword(key []byte, value []byte) uint32 {
	// sign the value using HMAC-SHA1
	hmacSha1 := hmac.New(sha1.New, key)
	hmacSha1.Write(value)
	hash := hmacSha1.Sum(nil)

	// We're going to use a subset of the generated hash.
	// Using the last nibble (half-byte) to choose the index to start from.
	// This number is always appropriate as it's maximum decimal 15, the hash will
	// have the maximum index 19 (20 bytes of SHA1) and we need 4 bytes.
	offset := hash[len(hash)-1] & 0x0F

	// get a 32-bit (4-byte) chunk from the hash starting at offset
	hashParts := hash[offset : offset+4]

	// ignore the most significant bit as per RFC 4226
	hashParts[0] = hashParts[0] & 0x7F

	number := toUint32(hashParts)

	// size to 6 digits
	// one million is the first number with 7 digits so the remainder
	// of the division will always return < 7 digits
	pwd := number % 1000000

	return pwd
}

func oneTimePasswordSHA256(key []byte, value []byte) uint32 {
	// sign the value using HMAC-SHA1
	hmacSha256 := hmac.New(sha256.New, key)
	hmacSha256.Write(value)
	hash := hmacSha256.Sum(nil)

	// We're going to use a subset of the generated hash.
	// Using the last nibble (half-byte) to choose the index to start from.
	// This number is always appropriate as it's maximum decimal 15, the hash will
	// have the maximum index 19 (20 bytes of SHA1) and we need 4 bytes.
	offset := hash[len(hash)-1] & 0x0F

	// get a 32-bit (4-byte) chunk from the hash starting at offset
	hashParts := hash[offset : offset+4]

	// ignore the most significant bit as per RFC 4226
	hashParts[0] = hashParts[0] & 0x7F

	number := toUint32(hashParts)

	// size to 6 digits
	// one million is the first number with 7 digits so the remainder
	// of the division will always return < 7 digits
	pwd := number % 1000000

	return pwd
}

// totpURLParse  totpURL解析
func totpURLParse(totpURL string) (url.Values, error) {
	s := strings.Split(totpURL, "?")
	m, err := url.ParseQuery(s[1])
	if err != nil {
		return nil, err
	}
	return m, nil
}

type TOTP struct {
	SecretKey string // secret
	Algorithm string // 加密算法
	Issuer    string // 发行者
	Name      string // 名称
	Digits    int    // 位数
}

// TwoStepAuthGenNumber 根据提供的 secret 来生成6位数字
// 返回 6位数字、剩余时间
func TwoStepAuthGenNumber(t *TOTP) (string, int64, error) {

	decodeSecret, err := convertSecret(t.SecretKey)
	if err != nil {
		return "", 0, err
	}
	// generate a one-time password using the time at 30-second intervals
	var pwd uint32
	epochSeconds := time.Now().Unix()
	if t.Algorithm == "SHA256" {
		pwd = oneTimePasswordSHA256(decodeSecret, toBytes(epochSeconds/30))
	} else if t.Algorithm == "SHA1" || t.Algorithm == "" {
		pwd = oneTimePassword(decodeSecret, toBytes(epochSeconds/30))
	} else {
		return "", 0, fmt.Errorf("not support algorithm: %s", t.Algorithm)
	}
	secondsRemaining := 30 - (epochSeconds % 30)

	return fmt.Sprintf("%06d", pwd), secondsRemaining, nil
}

// TwoStepAuthGenByQRCode 解析二维码图片
func TwoStepAuthParseQRCode(qrcodePath string) (*TOTP, error) {
	f, err := os.Open(qrcodePath)
	if err != nil {
		return nil, err
	}
	content, err := QRCodeParse(f)
	if err != nil {
		return nil, err
	}

	r, err := totpURLParse(content)
	if err != nil {
		return nil, err
	}

	totp := &TOTP{SecretKey: r["secret"][0]}

	if _, ok := r["algorithm"]; ok {
		totp.Algorithm = r["algorithm"][0]
	}

	if _, ok := r["issuer"]; ok {
		totp.Issuer = r["issuer"][0]
	}

	if _, ok := r["digits"]; ok {
		i, _ := strconv.Atoi(r["digits"][0])
		totp.Digits = i
	}

	return totp, nil
}

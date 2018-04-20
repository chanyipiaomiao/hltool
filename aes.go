package hltool

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

func pKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// GoAES 加密
type GoAES struct {
	Key []byte
}

// NewGoAES 返回GoAES
func NewGoAES(key []byte) *GoAES {
	return &GoAES{Key: key}
}

// Encrypt 加密数据
func (a *GoAES) Encrypt(origData []byte) ([]byte, error) {
	block, err := aes.NewCipher(a.Key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = pKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, a.Key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

// Decrypt 解密数据
func (a *GoAES) Decrypt(crypted []byte) ([]byte, error) {
	block, err := aes.NewCipher(a.Key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, a.Key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = pKCS7UnPadding(origData)
	return origData, nil
}

package hltool

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
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
	opt *AesOption

	Key []byte
}

type AesOption struct {
	Mode string // ecb, cbc
}

// NewGoAES 返回GoAES
func NewGoAES(key []byte) *GoAES {
	return &GoAES{Key: key}
}

// NewGoAES 返回GoAES
func NewAES(key []byte, opt *AesOption) *GoAES {
	a := &GoAES{Key: key}
	if opt == nil {
		a.opt = &AesOption{Mode: "cbc"}
	} else {
		a.opt = opt
	}
	return a
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

func (a *GoAES) EncryptV2(origData []byte) ([]byte, error) {
	switch a.opt.Mode {
	case "ecb":
		block, err := aes.NewCipher(a.Key)
		if err != nil {
			return nil, err
		}
		blockSize := block.BlockSize()
		origData = pKCS7Padding(origData, blockSize)
		crypted := make([]byte, len(origData))
		if len(origData)%blockSize != 0 {
			return nil, errors.New("crypto/cipher: input not full blocks")
		}
		if len(crypted) < len(origData) {
			return nil, errors.New("crypto/cipher: output smaller than input")
		}
		for len(origData) > 0 {
			block.Encrypt(crypted, origData[:blockSize])
			origData = origData[blockSize:]
			crypted = crypted[blockSize:]
		}
		return crypted, nil
	case "cbc":
		return a.Encrypt(origData)
	}
	return nil, errors.New("unsupported mode")
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

func (a *GoAES) DecryptV2(crypted []byte) ([]byte, error) {
	switch a.opt.Mode {
	case "ecb":
		origData := make([]byte, len(crypted))
		block, err := aes.NewCipher(a.Key)
		if err != nil {
			return nil, err
		}
		blockSize := block.BlockSize()
		if len(crypted)%blockSize != 0 {
			return nil, errors.New("crypto/cipher: input not full blocks")
		}
		if len(origData) < len(crypted) {
			return nil, errors.New("crypto/cipher: output smaller than input")
		}
		for len(crypted) > 0 {
			block.Decrypt(origData, crypted[:blockSize])
			origData = pKCS7UnPadding(origData)
		}
		return origData, nil
	case "cbc":
		return a.Decrypt(crypted)
	}
	return nil, errors.New("unsupported mode")
}

package hltool

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io"
	"io/ioutil"
	"os"
)

const (

	// RSAAlgorithmSign RSA签名算法
	RSAAlgorithmSign = crypto.SHA256
)

// NewRSAFile 生成密钥对文件
// pubKeyFilename: 公钥文件名 priKeyFilename: 私钥文件名 kekeyLength: 密钥长度
func NewRSAFile(pubKeyFilename, priKeyFilename string, keyLength int) error {
	if pubKeyFilename == "" {
		pubKeyFilename = "id_rsa.pub"
	}
	if priKeyFilename == "" {
		priKeyFilename = "id_rsa"
	}

	if keyLength == 0 || keyLength < 1024 {
		keyLength = 1024
	}

	// 创建公钥文件
	pubWriter, err := os.Create(pubKeyFilename)
	if err != nil {
		return err
	}
	defer pubWriter.Close()

	// 创建私钥文件
	priWriter, err := os.Create(priKeyFilename)
	if err != nil {
		return err
	}
	defer priWriter.Close()

	// 生成密钥对
	err = WriteRSAKeyPair(pubWriter, priWriter, keyLength)
	if err != nil {
		return err
	}

	return nil
}

// NewRSAString 生成密钥对字符串
// keyLength 密钥的长度
func NewRSAString(keyLength int) (string, string, error) {

	if keyLength == 0 || keyLength < 1024 {
		keyLength = 1024
	}

	bufPub := make([]byte, 1024*5)
	pubuffer := bytes.NewBuffer(bufPub)

	bufPri := make([]byte, 1024*5)
	pribuffer := bytes.NewBuffer(bufPri)

	err := WriteRSAKeyPair(pubuffer, pribuffer, keyLength)
	if err != nil {
		return "", "", nil
	}

	return pubuffer.String(), pribuffer.String(), nil
}

// WriteRSAKeyPair 生成RSA密钥对
func WriteRSAKeyPair(publicKeyWriter, privateKeyWriter io.Writer, keyLength int) error {

	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, keyLength)
	if err != nil {
		return err
	}

	derStream := MarshalPKCS8PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: derStream,
	}
	err = pem.Encode(privateKeyWriter, block)
	if err != nil {
		return err
	}

	// 生成公钥文件
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)

	if err != nil {
		return err
	}

	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}

	err = pem.Encode(publicKeyWriter, block)
	if err != nil {
		return err
	}

	return nil

}

// ReadRSAKeyPair 读取RSA密钥对
// pubKeyFilename: 公钥文件名称   priKeyFilename: 私钥文件名
func ReadRSAKeyPair(pubKeyFilename, priKeyFilename string) ([]byte, []byte, error) {
	pub, err := ioutil.ReadFile(pubKeyFilename)
	if err != nil {
		return nil, nil, err
	}

	pri, err := ioutil.ReadFile(priKeyFilename)
	if err != nil {
		return nil, nil, err
	}
	return pub, pri, nil
}

// GoRSA RSA加密解密
type GoRSA struct {
	PublicKey  *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
}

// NewGoRSA 初始化 GoRSA对象
func NewGoRSA(pubKeyFilename, priKeyFilename string) (*GoRSA, error) {

	publicKey, privateKey, err := ReadRSAKeyPair(pubKeyFilename, priKeyFilename)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	block, _ = pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error")
	}

	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	pri, ok := priv.(*rsa.PrivateKey)
	if ok {
		return &GoRSA{
			PublicKey:  pub,
			PrivateKey: pri,
		}, nil
	}
	return nil, errors.New("private key not supported")
}

// PublicEncrypt 公钥加密
func (r *GoRSA) PublicEncrypt(data string) (string, error) {
	partLen := r.PublicKey.N.BitLen()/8 - 11
	chunks := split([]byte(data), partLen)
	buffer := bytes.NewBufferString("")

	for _, chunk := range chunks {
		bytes, err := rsa.EncryptPKCS1v15(rand.Reader, r.PublicKey, chunk)
		if err != nil {
			return "", err
		}
		buffer.Write(bytes)
	}
	return base64.RawURLEncoding.EncodeToString(buffer.Bytes()), nil

}

// PrivateDecrypt 私钥解密
func (r *GoRSA) PrivateDecrypt(encrypted string) (string, error) {

	partLen := r.PublicKey.N.BitLen() / 8
	raw, err := base64.RawURLEncoding.DecodeString(encrypted)
	chunks := split([]byte(raw), partLen)
	buffer := bytes.NewBufferString("")

	for _, chunk := range chunks {
		decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, r.PrivateKey, chunk)
		if err != nil {
			return "", err
		}
		buffer.Write(decrypted)
	}
	return buffer.String(), err
}

// Sign 数据进行签名
func (r *GoRSA) Sign(data string) (string, error) {

	h := RSAAlgorithmSign.New()
	h.Write([]byte(data))
	hashed := h.Sum(nil)
	sign, err := rsa.SignPKCS1v15(rand.Reader, r.PrivateKey, RSAAlgorithmSign, hashed)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(sign), err
}

// Verify 数据验证签名
func (r *GoRSA) Verify(data string, sign string) error {

	h := RSAAlgorithmSign.New()
	h.Write([]byte(data))
	hashed := h.Sum(nil)
	decodedSign, err := base64.RawURLEncoding.DecodeString(sign)
	if err != nil {
		return err
	}
	return rsa.VerifyPKCS1v15(r.PublicKey, RSAAlgorithmSign, hashed, decodedSign)
}

// MarshalPKCS8PrivateKey 私钥解析
func MarshalPKCS8PrivateKey(key *rsa.PrivateKey) []byte {

	info := struct {
		Version             int
		PrivateKeyAlgorithm []asn1.ObjectIdentifier
		PrivateKey          []byte
	}{}

	info.Version = 0
	info.PrivateKeyAlgorithm = make([]asn1.ObjectIdentifier, 1)
	info.PrivateKeyAlgorithm[0] = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 1}
	info.PrivateKey = x509.MarshalPKCS1PrivateKey(key)
	k, _ := asn1.Marshal(info)
	return k

}

func split(buf []byte, lim int) [][]byte {
	var chunk []byte
	chunks := make([][]byte, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf[:len(buf)])
	}
	return chunks
}

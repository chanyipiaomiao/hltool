package hltool

import (
	"bytes"
	"encoding/gob"
	"os"
)

// StructToBinFile 结构体序列化成二级制文件
// structObj 结构体对象
// filepath 文件路径
func StructToBinFile(structObj interface{}, filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	enc := gob.NewEncoder(file)
	err = enc.Encode(structObj)
	if err != nil {
		return err
	}
	return nil
}

// BinFileToStruct 二级制文件反序列化为结构体,结构体必须要和转换前的结构一致
// filepath 二进制文件路径
// to 结构体对象 结构必须和序列化前的一样
func BinFileToStruct(filepath string, to interface{}) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	dec := gob.NewDecoder(file)
	err = dec.Decode(to)
	if err != nil {
		return err
	}
	return nil
}

// StructToBytes 结构体转换为[]byte
// structObj 结构体对象
func StructToBytes(structObj interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(structObj)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// BytesToStruct []byte转换为结构体,必须事先知道结构体的结构，而且必须一样
// data 转换后的字节数组
// to 结构体对象 结构必须和序列化前的一样
func BytesToStruct(data []byte, to interface{}) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(to)
}

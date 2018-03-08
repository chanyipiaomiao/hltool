package hltool

import (
	"encoding/json"
	"os"
)

// JSONFileToBytes 从json文件中转换为[]byte
func JSONFileToBytes(filepath string) ([]byte, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	fileinfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	buf := make([]byte, fileinfo.Size())
	_, err = file.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

// JSONBytesToStruct json []byte 转换为 struct
func JSONBytesToStruct(data []byte, structObj interface{}) error {
	err := json.Unmarshal(data, structObj)
	if err != nil {
		return err
	}
	return nil
}

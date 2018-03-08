package hltool

import (
	"bufio"
	"os"
)

// BytesToFile 字节数组写入到文件
func BytesToFile(data []byte, filepath string) error {
	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = writer.Write(data)
	if err != nil {
		return err
	}
	writer.Flush()
	return nil
}

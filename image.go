package hltool

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/png"
	"net/http"
	"os"
)

// ImageType 探测图片的类型
// imgbytes 图片字节数组
func ImageType(imgbytes []byte) string {
	return http.DetectContentType(imgbytes)
}

// BytesToPng 字节生成png图片
// imgbytes 图片字节数组
// filepath 文件路径名称
func BytesToPng(imgbytes []byte, filepath string) error {

	if ImageType(imgbytes) != "image/png" {
		return fmt.Errorf("it seem not a png image type")
	}
	img, _, err := image.Decode(bytes.NewReader(imgbytes))
	if err != nil {
		return err
	}

	fd, err := os.Create(filepath)
	if err != nil {
		return err
	}

	err = png.Encode(fd, img)
	if err != nil {
		return err
	}

	return err
}

// ImageToBytes 图片转换为字节数组
// filepath 图片的路径
func ImageToBytes(filepath string) ([]byte, error) {

	fd, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	fileInfo, err := fd.Stat()
	if err != nil {
		return nil, err
	}

	bytes := make([]byte, fileInfo.Size())
	buffer := bufio.NewReader(fd)
	_, err = buffer.Read(bytes)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

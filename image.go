package hltool

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
)

// ImageType 探测图片的类型
// imgbytes 图片字节数组
func ImageType(imgbytes []byte) string {
	return http.DetectContentType(imgbytes)
}

func decodeBytesCreateFile(imgbytes []byte, filepath string) (image.Image, *os.File, error) {
	img, _, err := image.Decode(bytes.NewReader(imgbytes))
	if err != nil {
		return nil, nil, err
	}

	fd, err := os.Create(filepath)
	if err != nil {
		return nil, nil, err
	}
	return img, fd, nil
}

// BytesToImage []byte生成图片
// imgbytes 图片[]byte数组
// filepath 文件路径名称
func BytesToImage(imgbytes []byte, filepath string) error {

	switch ImageType(imgbytes) {
	case "image/png":
		img, fd, err := decodeBytesCreateFile(imgbytes, filepath)
		if err != nil {
			return err
		}
		defer fd.Close()
		err = png.Encode(fd, img)
		if err != nil {
			return err
		}
	case "image/jpeg":
		img, fd, err := decodeBytesCreateFile(imgbytes, filepath)
		if err != nil {
			return err
		}
		defer fd.Close()
		err = jpeg.Encode(fd, img, nil)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown image type")
	}

	return nil
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

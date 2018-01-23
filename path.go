package hltool

import (
	"os"
)

// CurrentDir 获取当前路径
func CurrentDir() string {
	currentDir, _ := os.Getwd()
	return currentDir
}


// IsExist 文件或目录是否存在
// return false 表示文件不存在
func IsExist(name string) bool {
	_, err := os.Stat(name)
	return err == nil || os.IsExist(err)
}

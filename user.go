package hltool

import (
	"os/user"
	"os"
)

// CurrentUser 获取当前SSH连接的用户
func CurrentUser() string {
	currentUser, _ := user.Current()
	return currentUser.Username
}

// UserHome 获取用户的家目录
func UserHome() (string, error) {
	currentUser, err := user.Current()
	if err == nil {
		return currentUser.HomeDir, nil
	}
	return os.Getenv("HOME"), nil
}
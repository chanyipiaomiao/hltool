package hltool

import (
	"regexp"
)

// IsNumber 检查输入的字符串是否匹配数字
func IsNumber(input string) bool {

	if isNumber, _ := regexp.Match("^\\d+$", []byte(input)); isNumber {
		return true
	}
	return false
}
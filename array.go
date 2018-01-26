package hltool

// InStringSlice 元素是否在一个string类型的slice里面
func InStringSlice(s []string, x string) (bool, int) {
	for i, v := range s {
		if x == v {
			return true, i
		}
	}
	return false, -1
}
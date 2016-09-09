// 过滤器
package filter

var FilersMap = map[string]func([]byte) bool{
	"point": PointFilter,
	"coma":  ComaFilter,
}

func Register(name string, f func([]byte) bool) {
	FilersMap[name] = f
}

// 过滤函数
//    data: 待过滤的数据
//    filterList: 使用的过滤器名称列表，这些名称一定是出现在FilersMap中的才有效
func Filter(data []byte, filterList []string) bool {
	for _, filterName := range filterList {
		if f, ok := FilersMap[filterName]; ok {
			if f(data) {
				return true
			}
		}
	}
	return false
}

// '.'号过滤器，如果字符串中包含了'.'，就过滤掉
func PointFilter(data []byte) bool {
	for _, b := range data {
		if b == 46 {
			return true
		}
	}
	return false
}

// ','号过滤器
func ComaFilter(data []byte) bool {
	for _, b := range data {
		if b == 44 {
			return true
		}
	}
	return false
}

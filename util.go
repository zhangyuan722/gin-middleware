package m

import (
	"fmt"
	"path/filepath"
	"reflect"
)

// firstNonZeroValue 接受任意数量的参数，并返回第一个非零值参数
func firstNonZeroValue[T any](args ...T) *T {
	if len(args) == 0 {
		return nil
	}

	for _, arg := range args {
		zeroValue := reflect.Zero(reflect.TypeOf(arg)).Interface()

		if reflect.DeepEqual(arg, zeroValue) {
			continue
		}

		return &arg
	}
	return nil
}

// matchPattern 判断给定的路径是否与指定的模式匹配。
//
// 使用 filepath.Match 函数执行匹配，并返回匹配结果。如果匹配成功，
// 返回 true，否则返回 false。如果在匹配过程中发生错误，将记录日志并返回 false。
//
// 参数：
//   - pattern: 要匹配的模式，支持通配符。
//   - path: 要检查的路径。
//
// 返回值：
//   - bool: 如果路径匹配模式，返回 true；否则返回 false。
//
// 示例：
//
//	matched := matchPattern("/swagger/*", "/swagger/index.html")
//	if matched {
//	    fmt.Println("Path matches the pattern.")
//	} else {
//	    fmt.Println("Path does not match the pattern.")
//	}
func matchPattern(pattern, path string) bool {
	matched, err := filepath.Match(pattern, path)
	if err != nil {
		fmt.Println(fmt.Sprintf("match pattern error, %s", err.Error()))
		return false
	}
	return matched
}

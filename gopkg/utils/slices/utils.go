package slices

import "strings"

func InSlice[T comparable](slice []T, v T) bool {
	for _, val := range slice {
		if v == val {
			return true
		}
	}
	return false
}

// RemoveDuplicateElement 数组去重
func RemoveDuplicateElement[T comparable](list []T) []T {
	var result = make([]T, 0, len(list))
	var m = map[T]struct{}{}

	for _, v := range list {
		if _, ok := m[v]; !ok {
			m[v] = struct{}{}
			result = append(result, v)
		}
	}

	return result
}

func Split(strs []string, sep string) (results []string) {
	for _, str := range strs {
		results = append(results, strings.Split(str, sep)...)
	}
	return
}

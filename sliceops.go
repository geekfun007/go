package sliceops

// sliceops.go - Slice 操作函数库

// InsertAt 在指定位置插入元素
func InsertAt(s []int, index, value int) []int {
	s = append(s, 0)
	copy(s[index+1:], s[index:])
	s[index] = value
	return s
}

// InsertMultiple 在指定位置插入多个元素
func InsertMultiple(s []int, index int, values []int) []int {
	return append(s[:index], append(values, s[index:]...)...)
}

// RemoveAt 删除指定位置的元素（保持顺序）
func RemoveAt(s []int, index int) []int {
	return append(s[:index], s[index+1:]...)
}

// RemoveAtFast 删除指定位置的元素（不保持顺序）
func RemoveAtFast(s []int, index int) []int {
	s[index] = s[len(s)-1]
	return s[:len(s)-1]
}

// RemoveRange 删除范围内的元素
func RemoveRange(s []int, start, end int) []int {
	return append(s[:start], s[end:]...)
}

// RemoveIf 删除满足条件的元素
func RemoveIf(s []int, condition func(int) bool) []int {
	result := s[:0]
	for _, v := range s {
		if !condition(v) {
			result = append(result, v)
		}
	}
	return result
}

// RemoveDuplicates 删除重复元素（保持顺序）
func RemoveDuplicates(s []int) []int {
	seen := make(map[int]bool)
	result := make([]int, 0, len(s))
	for _, v := range s {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return result
}

// Reverse 原地反转 slice
func Reverse(s []int) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// ReverseCopy 创建反转的副本
func ReverseCopy(s []string) []string {
	result := make([]string, len(s))
	for i, v := range s {
		result[len(s)-1-i] = v
	}
	return result
}

// UniqueOrdered 去重（保持顺序）
func UniqueOrdered(s []int) []int {
	seen := make(map[int]bool)
	result := make([]int, 0, len(s))
	for _, v := range s {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return result
}

// UniqueFast 去重（不保证顺序）
func UniqueFast(s []int) []int {
	seen := make(map[int]bool)
	for _, v := range s {
		seen[v] = true
	}
	result := make([]int, 0, len(seen))
	for v := range seen {
		result = append(result, v)
	}
	return result
}

// UniqueStrings 字符串去重
func UniqueStrings(s []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0, len(s))
	for _, v := range s {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return result
}

// Filter 过滤整数 slice
func Filter(s []int, fn func(int) bool) []int {
	result := make([]int, 0, len(s))
	for _, v := range s {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}

// FilterStrings 过滤字符串 slice
func FilterStrings(s []string, fn func(string) bool) []string {
	result := make([]string, 0, len(s))
	for _, v := range s {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}

// FilterInPlace 原地过滤（复用底层数组）
func FilterInPlace(s []int, fn func(int) bool) []int {
	n := 0
	for _, v := range s {
		if fn(v) {
			s[n] = v
			n++
		}
	}
	return s[:n]
}

// Map 映射转换整数 slice
func Map(s []int, fn func(int) int) []int {
	result := make([]int, len(s))
	for i, v := range s {
		result[i] = fn(v)
	}
	return result
}

// MapString 映射转换字符串 slice
func MapString(s []string, fn func(string) string) []string {
	result := make([]string, len(s))
	for i, v := range s {
		result[i] = fn(v)
	}
	return result
}

// MapStringToInt 字符串 slice 映射到整数 slice
func MapStringToInt(s []string, fn func(string) int) []int {
	result := make([]int, len(s))
	for i, v := range s {
		result[i] = fn(v)
	}
	return result
}

// MapIntToString 整数 slice 映射到字符串 slice
func MapIntToString(s []int, fn func(int) string) []string {
	result := make([]string, len(s))
	for i, v := range s {
		result[i] = fn(v)
	}
	return result
}

// Reduce 归约操作
func Reduce(s []int, initial int, fn func(int, int) int) int {
	acc := initial
	for _, v := range s {
		acc = fn(acc, v)
	}
	return acc
}

// ReduceString 字符串归约操作
func ReduceString(s []string, initial string, fn func(string, string) string) string {
	acc := initial
	for _, v := range s {
		acc = fn(acc, v)
	}
	return acc
}

// GroupBy 按条件分组
func GroupBy(s []int, fn func(int) string) map[string][]int {
	result := make(map[string][]int)
	for _, v := range s {
		key := fn(v)
		result[key] = append(result[key], v)
	}
	return result
}

// GroupByString 字符串按条件分组
func GroupByString(s []string, fn func(string) int) map[int][]string {
	result := make(map[int][]string)
	for _, v := range s {
		key := fn(v)
		result[key] = append(result[key], v)
	}
	return result
}

// Chunk 将 slice 分块
func Chunk(s []int, size int) [][]int {
	var result [][]int
	for i := 0; i < len(s); i += size {
		end := i + size
		if end > len(s) {
			end = len(s)
		}
		result = append(result, s[i:end])
	}
	return result
}

// Contains 检查是否包含元素
func Contains(s []int, target int) bool {
	for _, v := range s {
		if v == target {
			return true
		}
	}
	return false
}

// ContainsString 检查是否包含字符串
func ContainsString(s []string, target string) bool {
	for _, v := range s {
		if v == target {
			return true
		}
	}
	return false
}

// IndexOf 查找第一个匹配的索引
func IndexOf(s []int, target int) int {
	for i, v := range s {
		if v == target {
			return i
		}
	}
	return -1
}

// LastIndexOf 查找最后一个匹配的索引
func LastIndexOf(s []int, target int) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == target {
			return i
		}
	}
	return -1
}

// FindAll 查找所有匹配的索引
func FindAll(s []int, target int) []int {
	var indices []int
	for i, v := range s {
		if v == target {
			indices = append(indices, i)
		}
	}
	return indices
}

// FindFirst 查找第一个满足条件的元素索引
func FindFirst(s []int, condition func(int) bool) int {
	for i, v := range s {
		if condition(v) {
			return i
		}
	}
	return -1
}

// FindMinMax 查找最小和最大值
func FindMinMax(s []int) (min, max int) {
	if len(s) == 0 {
		return 0, 0
	}
	min, max = s[0], s[0]
	for _, v := range s[1:] {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return
}

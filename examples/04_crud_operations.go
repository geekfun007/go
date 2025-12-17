package main

import "fmt"

// 04_crud_operations.go - Slice 增删改查操作

func main() {
	fmt.Println("=== Slice 增删改查操作 ===\n")

	// 1. 插入元素
	demonstrateInsert()

	// 2. 删除元素
	demonstrateDelete()

	// 3. 更新元素
	demonstrateUpdate()

	// 4. 查找元素
	demonstrateSearch()
}

func demonstrateInsert() {
	fmt.Println("--- 1. 插入元素 ---")

	// 在开头插入
	s := []int{2, 3, 4}
	s = append([]int{1}, s...)
	fmt.Printf("在开头插入 1: %v\n", s)

	// 在中间插入单个元素
	s = []int{1, 2, 4, 5}
	fmt.Printf("原始: %v\n", s)
	s = insertAt(s, 2, 3)
	fmt.Printf("在索引 2 插入 3: %v\n", s)

	// 在中间插入多个元素
	s = []int{1, 2, 5, 6}
	fmt.Printf("\n原始: %v\n", s)
	s = insertMultiple(s, 2, []int{3, 4})
	fmt.Printf("在索引 2 插入 [3, 4]: %v\n", s)

	// 在末尾插入
	s = append(s, 7, 8, 9)
	fmt.Printf("在末尾插入 7, 8, 9: %v\n", s)

	// 高效插入（预分配）
	fmt.Println("\n高效插入方式:")
	original := []int{1, 2, 5, 6}
	toInsert := []int{3, 4}
	index := 2
	
	result := make([]int, len(original)+len(toInsert))
	copy(result, original[:index])
	copy(result[index:], toInsert)
	copy(result[index+len(toInsert):], original[index:])
	fmt.Printf("结果: %v\n", result)

	fmt.Println()
}

func demonstrateDelete() {
	fmt.Println("--- 2. 删除元素 ---")

	// 删除第一个元素
	s := []int{1, 2, 3, 4, 5}
	fmt.Printf("原始: %v\n", s)
	s = s[1:]
	fmt.Printf("删除第一个: %v\n", s)

	// 删除最后一个元素
	s = []int{1, 2, 3, 4, 5}
	s = s[:len(s)-1]
	fmt.Printf("删除最后一个: %v\n", s)

	// 删除中间元素（保持顺序）
	s = []int{1, 2, 3, 4, 5}
	fmt.Printf("\n原始: %v\n", s)
	s = removeAt(s, 2)
	fmt.Printf("删除索引 2（保持顺序）: %v\n", s)

	// 删除中间元素（不保持顺序，更快）
	s = []int{1, 2, 3, 4, 5}
	fmt.Printf("\n原始: %v\n", s)
	s = removeAtFast(s, 2)
	fmt.Printf("删除索引 2（不保持顺序）: %v\n", s)

	// 删除范围
	s = []int{1, 2, 3, 4, 5, 6, 7, 8}
	fmt.Printf("\n原始: %v\n", s)
	s = removeRange(s, 2, 5)
	fmt.Printf("删除索引 2-4: %v\n", s)

	// 删除所有满足条件的元素
	s = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	fmt.Printf("\n原始: %v\n", s)
	s = removeIf(s, func(n int) bool { return n%2 == 0 })
	fmt.Printf("删除偶数: %v\n", s)

	// 删除重复元素
	s = []int{1, 2, 2, 3, 3, 3, 4, 5, 5}
	fmt.Printf("\n原始: %v\n", s)
	s = removeDuplicates(s)
	fmt.Printf("删除重复: %v\n", s)

	fmt.Println()
}

func demonstrateUpdate() {
	fmt.Println("--- 3. 更新元素 ---")

	// 更新单个元素
	s := []int{1, 2, 3, 4, 5}
	fmt.Printf("原始: %v\n", s)
	s[2] = 100
	fmt.Printf("更新索引 2: %v\n", s)

	// 更新范围
	s = []int{1, 2, 3, 4, 5}
	fmt.Printf("\n原始: %v\n", s)
	for i := 1; i <= 3; i++ {
		s[i] = s[i] * 10
	}
	fmt.Printf("将索引 1-3 的元素 *10: %v\n", s)

	// 使用函数更新所有元素
	s = []int{1, 2, 3, 4, 5}
	fmt.Printf("\n原始: %v\n", s)
	updateAll(s, func(n int) int { return n * n })
	fmt.Printf("所有元素平方: %v\n", s)

	// 条件更新
	s = []int{1, 2, 3, 4, 5, 6, 7, 8}
	fmt.Printf("\n原始: %v\n", s)
	updateIf(s, func(n int) bool { return n%2 == 0 }, func(n int) int { return n * 10 })
	fmt.Printf("偶数元素 *10: %v\n", s)

	// 交换元素
	s = []int{1, 2, 3, 4, 5}
	fmt.Printf("\n原始: %v\n", s)
	s[1], s[3] = s[3], s[1]
	fmt.Printf("交换索引 1 和 3: %v\n", s)

	fmt.Println()
}

func demonstrateSearch() {
	fmt.Println("--- 4. 查找元素 ---")

	s := []string{"apple", "banana", "cherry", "date", "elderberry"}

	// 查找第一个匹配
	index := indexOf(s, "cherry")
	fmt.Printf("slice: %v\n", s)
	fmt.Printf("'cherry' 的索引: %d\n", index)

	// 查找最后一个匹配
	s2 := []int{1, 2, 3, 2, 4, 2, 5}
	lastIdx := lastIndexOf(s2, 2)
	fmt.Printf("\nslice: %v\n", s2)
	fmt.Printf("2 的最后索引: %d\n", lastIdx)

	// 检查是否包含
	hasBanana := contains(s, "banana")
	fmt.Printf("\n包含 'banana': %v\n", hasBanana)
	hasGrape := contains(s, "grape")
	fmt.Printf("包含 'grape': %v\n", hasGrape)

	// 查找所有匹配的索引
	s2 = []int{1, 2, 3, 2, 4, 2, 5}
	indices := findAll(s2, 2)
	fmt.Printf("\nslice: %v\n", s2)
	fmt.Printf("2 的所有索引: %v\n", indices)

	// 条件查找
	s3 := []int{1, 3, 5, 7, 8, 9, 11}
	fmt.Printf("\nslice: %v\n", s3)
	idx := findFirst(s3, func(n int) bool { return n%2 == 0 })
	fmt.Printf("第一个偶数的索引: %d, 值: %d\n", idx, s3[idx])

	// 查找最大最小值
	s4 := []int{3, 7, 2, 9, 1, 5, 8}
	min, max := findMinMax(s4)
	fmt.Printf("\nslice: %v\n", s4)
	fmt.Printf("最小值: %d, 最大值: %d\n", min, max)

	fmt.Println()
}

// === 辅助函数 ===

// insertAt 在指定位置插入元素
func insertAt(s []int, index, value int) []int {
	s = append(s, 0)
	copy(s[index+1:], s[index:])
	s[index] = value
	return s
}

// insertMultiple 在指定位置插入多个元素
func insertMultiple(s []int, index int, values []int) []int {
	return append(s[:index], append(values, s[index:]...)...)
}

// removeAt 删除指定位置的元素（保持顺序）
func removeAt(s []int, index int) []int {
	return append(s[:index], s[index+1:]...)
}

// removeAtFast 删除指定位置的元素（不保持顺序）
func removeAtFast(s []int, index int) []int {
	s[index] = s[len(s)-1]
	return s[:len(s)-1]
}

// removeRange 删除范围内的元素
func removeRange(s []int, start, end int) []int {
	return append(s[:start], s[end:]...)
}

// removeIf 删除满足条件的元素
func removeIf(s []int, condition func(int) bool) []int {
	result := s[:0]
	for _, v := range s {
		if !condition(v) {
			result = append(result, v)
		}
	}
	return result
}

// removeDuplicates 删除重复元素（保持顺序）
func removeDuplicates(s []int) []int {
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

// updateAll 更新所有元素
func updateAll(s []int, fn func(int) int) {
	for i, v := range s {
		s[i] = fn(v)
	}
}

// updateIf 条件更新
func updateIf(s []int, condition func(int) bool, fn func(int) int) {
	for i, v := range s {
		if condition(v) {
			s[i] = fn(v)
		}
	}
}

// indexOf 查找第一个匹配的索引
func indexOf(s []string, target string) int {
	for i, v := range s {
		if v == target {
			return i
		}
	}
	return -1
}

// lastIndexOf 查找最后一个匹配的索引
func lastIndexOf(s []int, target int) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == target {
			return i
		}
	}
	return -1
}

// contains 检查是否包含元素
func contains(s []string, target string) bool {
	return indexOf(s, target) != -1
}

// findAll 查找所有匹配的索引
func findAll(s []int, target int) []int {
	var indices []int
	for i, v := range s {
		if v == target {
			indices = append(indices, i)
		}
	}
	return indices
}

// findFirst 查找第一个满足条件的元素索引
func findFirst(s []int, condition func(int) bool) int {
	for i, v := range s {
		if condition(v) {
			return i
		}
	}
	return -1
}

// findMinMax 查找最小和最大值
func findMinMax(s []int) (min, max int) {
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

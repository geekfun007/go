package main

import (
	"fmt"
	"sort"
	"strings"
)

// 05_advanced_techniques.go - Slice 高级技巧

func main() {
	fmt.Println("=== Slice 高级技巧 ===\n")

	// 1. 反转
	demonstrateReverse()

	// 2. 去重
	demonstrateUnique()

	// 3. 过滤
	demonstrateFilter()

	// 4. 映射转换
	demonstrateMap()

	// 5. 归约
	demonstrateReduce()

	// 6. 分组
	demonstrateGrouping()

	// 7. 排序
	demonstrateSorting()

	// 8. 二维切片
	demonstrate2DSlice()
}

func demonstrateReverse() {
	fmt.Println("--- 1. 反转 Slice ---")

	// 原地反转
	s := []int{1, 2, 3, 4, 5}
	fmt.Printf("原始: %v\n", s)
	reverse(s)
	fmt.Printf("反转后: %v\n", s)

	// 创建反转副本
	s2 := []string{"Go", "Python", "Java", "Rust"}
	fmt.Printf("\n原始: %v\n", s2)
	reversed := reverseCopy(s2)
	fmt.Printf("反转副本: %v\n", reversed)
	fmt.Printf("原始保持不变: %v\n", s2)

	fmt.Println()
}

func demonstrateUnique() {
	fmt.Println("--- 2. 去重 ---")

	// 保持顺序的去重
	s := []int{1, 2, 2, 3, 3, 3, 4, 5, 5, 1}
	fmt.Printf("原始: %v\n", s)
	unique := uniqueOrdered(s)
	fmt.Printf("去重（保持顺序）: %v\n", unique)

	// 不保证顺序的去重（更快）
	unique2 := uniqueFast(s)
	fmt.Printf("去重（不保证顺序）: %v\n", unique2)

	// 字符串去重
	words := []string{"apple", "banana", "apple", "cherry", "banana"}
	fmt.Printf("\n原始: %v\n", words)
	uniqueWords := uniqueStrings(words)
	fmt.Printf("去重: %v\n", uniqueWords)

	fmt.Println()
}

func demonstrateFilter() {
	fmt.Println("--- 3. 过滤 ---")

	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	fmt.Printf("原始: %v\n", numbers)

	// 过滤偶数
	evens := filter(numbers, func(n int) bool { return n%2 == 0 })
	fmt.Printf("偶数: %v\n", evens)

	// 过滤大于 5 的数
	gt5 := filter(numbers, func(n int) bool { return n > 5 })
	fmt.Printf("大于 5: %v\n", gt5)

	// 链式过滤
	result := filter(
		filter(numbers, func(n int) bool { return n%2 == 0 }),
		func(n int) bool { return n > 5 },
	)
	fmt.Printf("偶数且大于 5: %v\n", result)

	// 字符串过滤
	words := []string{"apple", "apricot", "banana", "avocado", "cherry"}
	fmt.Printf("\n原始单词: %v\n", words)
	aWords := filterStrings(words, func(s string) bool { 
		return strings.HasPrefix(s, "a") 
	})
	fmt.Printf("以 'a' 开头: %v\n", aWords)

	fmt.Println()
}

func demonstrateMap() {
	fmt.Println("--- 4. 映射转换 ---")

	// 数字转换
	numbers := []int{1, 2, 3, 4, 5}
	fmt.Printf("原始: %v\n", numbers)
	
	squared := mapInt(numbers, func(n int) int { return n * n })
	fmt.Printf("平方: %v\n", squared)

	doubled := mapInt(numbers, func(n int) int { return n * 2 })
	fmt.Printf("翻倍: %v\n", doubled)

	// 字符串转换
	words := []string{"hello", "world", "go"}
	fmt.Printf("\n原始: %v\n", words)
	
	upper := mapString(words, strings.ToUpper)
	fmt.Printf("大写: %v\n", upper)

	lengths := mapStringToInt(words, func(s string) int { return len(s) })
	fmt.Printf("长度: %v\n", lengths)

	// 类型转换
	ints := []int{1, 2, 3, 4, 5}
	strs := mapIntToString(ints, func(n int) string { 
		return fmt.Sprintf("数字%d", n) 
	})
	fmt.Printf("\n整数转字符串: %v\n", strs)

	fmt.Println()
}

func demonstrateReduce() {
	fmt.Println("--- 5. 归约 ---")

	numbers := []int{1, 2, 3, 4, 5}
	fmt.Printf("数组: %v\n", numbers)

	// 求和
	sum := reduce(numbers, 0, func(acc, n int) int { return acc + n })
	fmt.Printf("求和: %d\n", sum)

	// 求积
	product := reduce(numbers, 1, func(acc, n int) int { return acc * n })
	fmt.Printf("求积: %d\n", product)

	// 找最大值
	max := reduce(numbers, numbers[0], func(acc, n int) int {
		if n > acc {
			return n
		}
		return acc
	})
	fmt.Printf("最大值: %d\n", max)

	// 字符串连接
	words := []string{"Go", "is", "awesome"}
	fmt.Printf("\n单词: %v\n", words)
	sentence := reduceString(words, "", func(acc, s string) string {
		if acc == "" {
			return s
		}
		return acc + " " + s
	})
	fmt.Printf("连接: %s\n", sentence)

	fmt.Println()
}

func demonstrateGrouping() {
	fmt.Println("--- 6. 分组 ---")

	// 按奇偶分组
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	fmt.Printf("原始: %v\n", numbers)
	
	grouped := groupBy(numbers, func(n int) string {
		if n%2 == 0 {
			return "even"
		}
		return "odd"
	})
	fmt.Printf("按奇偶分组: %v\n", grouped)

	// 按长度分组
	words := []string{"a", "go", "hi", "cat", "dog", "hello", "world"}
	fmt.Printf("\n单词: %v\n", words)
	
	groupedWords := groupByString(words, func(s string) int {
		return len(s)
	})
	fmt.Printf("按长度分组:\n")
	for length, ws := range groupedWords {
		fmt.Printf("  长度 %d: %v\n", length, ws)
	}

	// 分块
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	fmt.Printf("\n原始: %v\n", data)
	chunks := chunk(data, 3)
	fmt.Printf("分块（每组 3 个）:\n")
	for i, c := range chunks {
		fmt.Printf("  块 %d: %v\n", i, c)
	}

	fmt.Println()
}

func demonstrateSorting() {
	fmt.Println("--- 7. 排序 ---")

	// 升序排序
	numbers := []int{5, 2, 8, 1, 9, 3, 7}
	fmt.Printf("原始: %v\n", numbers)
	sort.Ints(numbers)
	fmt.Printf("升序: %v\n", numbers)

	// 降序排序
	numbers = []int{5, 2, 8, 1, 9, 3, 7}
	sort.Sort(sort.Reverse(sort.IntSlice(numbers)))
	fmt.Printf("降序: %v\n", numbers)

	// 字符串排序
	words := []string{"banana", "apple", "cherry", "date"}
	fmt.Printf("\n原始: %v\n", words)
	sort.Strings(words)
	fmt.Printf("排序: %v\n", words)

	// 自定义排序
	people := []Person{
		{"Alice", 30},
		{"Bob", 25},
		{"Charlie", 35},
		{"David", 25},
	}
	fmt.Printf("\n原始:\n")
	printPeople(people)

	sort.Slice(people, func(i, j int) bool {
		return people[i].Age < people[j].Age
	})
	fmt.Printf("\n按年龄排序:\n")
	printPeople(people)

	sort.Slice(people, func(i, j int) bool {
		return people[i].Name < people[j].Name
	})
	fmt.Printf("\n按姓名排序:\n")
	printPeople(people)

	fmt.Println()
}

func demonstrate2DSlice() {
	fmt.Println("--- 8. 二维切片 ---")

	// 创建二维切片
	matrix := make([][]int, 3)
	for i := range matrix {
		matrix[i] = make([]int, 4)
	}
	
	// 填充数据
	count := 1
	for i := range matrix {
		for j := range matrix[i] {
			matrix[i][j] = count
			count++
		}
	}

	fmt.Println("3x4 矩阵:")
	printMatrix(matrix)

	// 不规则二维切片
	irregular := [][]int{
		{1, 2},
		{3, 4, 5},
		{6},
		{7, 8, 9, 10},
	}
	fmt.Println("\n不规则切片:")
	printMatrix(irregular)

	// 转置矩阵
	matrix2 := [][]int{
		{1, 2, 3},
		{4, 5, 6},
	}
	fmt.Println("\n原始矩阵:")
	printMatrix(matrix2)
	
	transposed := transpose(matrix2)
	fmt.Println("\n转置后:")
	printMatrix(transposed)

	fmt.Println()
}

// === 辅助函数 ===

type Person struct {
	Name string
	Age  int
}

func reverse(s []int) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func reverseCopy(s []string) []string {
	result := make([]string, len(s))
	for i, v := range s {
		result[len(s)-1-i] = v
	}
	return result
}

func uniqueOrdered(s []int) []int {
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

func uniqueFast(s []int) []int {
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

func uniqueStrings(s []string) []string {
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

func filter(s []int, fn func(int) bool) []int {
	result := make([]int, 0, len(s))
	for _, v := range s {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}

func filterStrings(s []string, fn func(string) bool) []string {
	result := make([]string, 0, len(s))
	for _, v := range s {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}

func mapInt(s []int, fn func(int) int) []int {
	result := make([]int, len(s))
	for i, v := range s {
		result[i] = fn(v)
	}
	return result
}

func mapString(s []string, fn func(string) string) []string {
	result := make([]string, len(s))
	for i, v := range s {
		result[i] = fn(v)
	}
	return result
}

func mapStringToInt(s []string, fn func(string) int) []int {
	result := make([]int, len(s))
	for i, v := range s {
		result[i] = fn(v)
	}
	return result
}

func mapIntToString(s []int, fn func(int) string) []string {
	result := make([]string, len(s))
	for i, v := range s {
		result[i] = fn(v)
	}
	return result
}

func reduce(s []int, initial int, fn func(int, int) int) int {
	acc := initial
	for _, v := range s {
		acc = fn(acc, v)
	}
	return acc
}

func reduceString(s []string, initial string, fn func(string, string) string) string {
	acc := initial
	for _, v := range s {
		acc = fn(acc, v)
	}
	return acc
}

func groupBy(s []int, fn func(int) string) map[string][]int {
	result := make(map[string][]int)
	for _, v := range s {
		key := fn(v)
		result[key] = append(result[key], v)
	}
	return result
}

func groupByString(s []string, fn func(string) int) map[int][]string {
	result := make(map[int][]string)
	for _, v := range s {
		key := fn(v)
		result[key] = append(result[key], v)
	}
	return result
}

func chunk(s []int, size int) [][]int {
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

func printPeople(people []Person) {
	for _, p := range people {
		fmt.Printf("  %s (年龄: %d)\n", p.Name, p.Age)
	}
}

func printMatrix(matrix [][]int) {
	for _, row := range matrix {
		fmt.Printf("  %v\n", row)
	}
}

func transpose(matrix [][]int) [][]int {
	if len(matrix) == 0 {
		return nil
	}
	rows, cols := len(matrix), len(matrix[0])
	result := make([][]int, cols)
	for i := range result {
		result[i] = make([]int, rows)
		for j := range result[i] {
			result[i][j] = matrix[j][i]
		}
	}
	return result
}

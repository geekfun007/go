// 示例 2：双值迭代器 (Seq2)
// 演示返回两个值的迭代器，类似遍历 map 的 key-value
package main

import "fmt"

// Enumerate 为切片添加索引
// 返回 (索引, 值) 的迭代器
func Enumerate[T any](slice []T) func(yield func(int, T) bool) {
	return func(yield func(int, T) bool) {
		for i, v := range slice {
			if !yield(i, v) {
				return
			}
		}
	}
}

// Zip 将两个切片合并为一个迭代器
// 返回 (slice1[i], slice2[i]) 直到其中一个耗尽
func Zip[T, U any](s1 []T, s2 []U) func(yield func(T, U) bool) {
	return func(yield func(T, U) bool) {
		minLen := len(s1)
		if len(s2) < minLen {
			minLen = len(s2)
		}
		for i := 0; i < minLen; i++ {
			if !yield(s1[i], s2[i]) {
				return
			}
		}
	}
}

// Entries 将 map 转换为迭代器
func Entries[K comparable, V any](m map[K]V) func(yield func(K, V) bool) {
	return func(yield func(K, V) bool) {
		for k, v := range m {
			if !yield(k, v) {
				return
			}
		}
	}
}

// Reverse 反向遍历切片，返回 (索引, 值)
func Reverse[T any](slice []T) func(yield func(int, T) bool) {
	return func(yield func(int, T) bool) {
		for i := len(slice) - 1; i >= 0; i-- {
			if !yield(i, slice[i]) {
				return
			}
		}
	}
}

func main() {
	fmt.Println("=== 双值迭代器 (Seq2) 示例 ===")
	fmt.Println()

	// 示例 1: Enumerate
	fmt.Println("1. Enumerate - 为切片添加索引:")
	fruits := []string{"apple", "banana", "cherry", "date"}
	for i, fruit := range Enumerate(fruits) {
		fmt.Printf("  [%d] %s\n", i, fruit)
	}
	fmt.Println()

	// 示例 2: Zip
	fmt.Println("2. Zip - 合并两个切片:")
	names := []string{"Alice", "Bob", "Charlie"}
	scores := []int{95, 87, 92, 88} // 注意：比 names 多一个元素
	for name, score := range Zip(names, scores) {
		fmt.Printf("  %s: %d 分\n", name, score)
	}
	fmt.Println()

	// 示例 3: Entries
	fmt.Println("3. Entries - 遍历 map:")
	config := map[string]string{
		"host":     "localhost",
		"port":     "8080",
		"protocol": "https",
	}
	for k, v := range Entries(config) {
		fmt.Printf("  %s = %s\n", k, v)
	}
	fmt.Println()

	// 示例 4: Reverse
	fmt.Println("4. Reverse - 反向遍历切片:")
	numbers := []int{1, 2, 3, 4, 5}
	for i, v := range Reverse(numbers) {
		fmt.Printf("  [%d] %d\n", i, v)
	}
	fmt.Println()

	// 示例 5: 只使用其中一个值
	fmt.Println("5. 使用 _ 忽略不需要的值:")
	for _, fruit := range Enumerate(fruits) {
		fmt.Printf("  水果: %s\n", fruit)
	}
}

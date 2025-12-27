// 示例 1：基础迭代器
// 演示如何创建和使用最简单的迭代器
package main

import "fmt"

// Range 创建一个从 0 到 n-1 的迭代器
// 这是最基础的迭代器模式
func Range(n int) func(yield func(int) bool) {
	return func(yield func(int) bool) {
		for i := 0; i < n; i++ {
			// yield 返回 false 表示消费者调用了 break
			if !yield(i) {
				return
			}
		}
	}
}

// RangeFrom 创建一个从 start 到 end-1 的迭代器
func RangeFrom(start, end int) func(yield func(int) bool) {
	return func(yield func(int) bool) {
		for i := start; i < end; i++ {
			if !yield(i) {
				return
			}
		}
	}
}

// RangeStep 创建一个带步长的迭代器
func RangeStep(start, end, step int) func(yield func(int) bool) {
	return func(yield func(int) bool) {
		for i := start; i < end; i += step {
			if !yield(i) {
				return
			}
		}
	}
}

func main() {
	fmt.Println("=== 基础迭代器示例 ===")
	fmt.Println()

	// 示例 1: 基本使用
	fmt.Println("1. Range(5) - 遍历 0 到 4:")
	for v := range Range(5) {
		fmt.Printf("  %d\n", v)
	}
	fmt.Println()

	// 示例 2: 带起始值
	fmt.Println("2. RangeFrom(3, 8) - 遍历 3 到 7:")
	for v := range RangeFrom(3, 8) {
		fmt.Printf("  %d\n", v)
	}
	fmt.Println()

	// 示例 3: 带步长
	fmt.Println("3. RangeStep(0, 10, 2) - 步长为 2:")
	for v := range RangeStep(0, 10, 2) {
		fmt.Printf("  %d\n", v)
	}
	fmt.Println()

	// 示例 4: 提前退出 (break)
	fmt.Println("4. Range(10) 但遇到 5 就 break:")
	for v := range Range(10) {
		if v == 5 {
			fmt.Println("  遇到 5，退出循环")
			break
		}
		fmt.Printf("  %d\n", v)
	}
	fmt.Println()

	// 示例 5: 在循环中使用 continue
	fmt.Println("5. Range(10) 跳过奇数:")
	for v := range Range(10) {
		if v%2 == 1 {
			continue
		}
		fmt.Printf("  %d\n", v)
	}
}

package main

import (
	"fmt"
	"time"
)

// 06_performance.go - Slice 性能优化

func main() {
	fmt.Println("=== Slice 性能优化 ===\n")

	// 1. 预分配容量
	demonstratePreallocation()

	// 2. append vs copy 性能对比
	demonstrateAppendVsCopy()

	// 3. 原地操作 vs 新建切片
	demonstrateInPlace()

	// 4. 完整切片表达式避免内存泄漏
	demonstrateMemoryLeak()

	// 5. 复用底层数组
	demonstrateReuseArray()
}

func demonstratePreallocation() {
	fmt.Println("--- 1. 预分配容量的重要性 ---")

	n := 100000

	// 不预分配
	start := time.Now()
	s1 := []int{}
	for i := 0; i < n; i++ {
		s1 = append(s1, i)
	}
	duration1 := time.Since(start)

	// 预分配
	start = time.Now()
	s2 := make([]int, 0, n)
	for i := 0; i < n; i++ {
		s2 = append(s2, i)
	}
	duration2 := time.Since(start)

	// 预分配并直接赋值
	start = time.Now()
	s3 := make([]int, n)
	for i := 0; i < n; i++ {
		s3[i] = i
	}
	duration3 := time.Since(start)

	fmt.Printf("添加 %d 个元素:\n", n)
	fmt.Printf("  不预分配:           %v\n", duration1)
	fmt.Printf("  预分配容量:         %v (提速 %.2fx)\n", duration2, float64(duration1)/float64(duration2))
	fmt.Printf("  预分配长度直接赋值: %v (提速 %.2fx)\n", duration3, float64(duration1)/float64(duration3))
	fmt.Println("\n结论: 预分配可以显著提高性能，避免多次扩容")

	fmt.Println()
}

func demonstrateAppendVsCopy() {
	fmt.Println("--- 2. append vs copy 性能对比 ---")

	n := 50000
	src := make([]int, n)
	for i := range src {
		src[i] = i
	}

	// 使用 append 合并
	start := time.Now()
	result1 := []int{}
	result1 = append(result1, src...)
	result1 = append(result1, src...)
	duration1 := time.Since(start)

	// 使用 copy 合并（预分配）
	start = time.Now()
	result2 := make([]int, len(src)*2)
	copy(result2, src)
	copy(result2[len(src):], src)
	duration2 := time.Since(start)

	fmt.Printf("合并两个长度为 %d 的 slice:\n", n)
	fmt.Printf("  使用 append:       %v\n", duration1)
	fmt.Printf("  使用 copy (预分配): %v (提速 %.2fx)\n", duration2, float64(duration1)/float64(duration2))
	fmt.Println("\n结论: 已知大小时，预分配 + copy 比 append 更快")

	fmt.Println()
}

func demonstrateInPlace() {
	fmt.Println("--- 3. 原地操作 vs 新建切片 ---")

	n := 100000

	// 新建切片过滤
	data1 := make([]int, n)
	for i := range data1 {
		data1[i] = i
	}
	start := time.Now()
	filtered1 := filterNewSlice(data1, func(x int) bool { return x%2 == 0 })
	duration1 := time.Since(start)

	// 原地过滤（复用底层数组）
	data2 := make([]int, n)
	for i := range data2 {
		data2[i] = i
	}
	start = time.Now()
	filtered2 := filterInPlace(data2, func(x int) bool { return x%2 == 0 })
	duration2 := time.Since(start)

	fmt.Printf("过滤 %d 个元素（保留偶数）:\n", n)
	fmt.Printf("  新建切片:   %v, 结果长度: %d\n", duration1, len(filtered1))
	fmt.Printf("  原地操作:   %v, 结果长度: %d (提速 %.2fx)\n", duration2, len(filtered2), float64(duration1)/float64(duration2))
	fmt.Println("\n结论: 原地操作避免内存分配，性能更好")

	fmt.Println()
}

func demonstrateMemoryLeak() {
	fmt.Println("--- 4. 完整切片表达式避免内存泄漏 ---")

	// 场景：从大切片中提取小切片
	fmt.Println("问题场景: 从 1MB 的 slice 中只需要前 10 个元素")

	largeSlice := make([]byte, 1024*1024) // 1MB
	for i := range largeSlice {
		largeSlice[i] = byte(i % 256)
	}

	// 方式 1: 简单切片（内存泄漏！）
	smallSlice1 := largeSlice[:10]
	fmt.Printf("\n方式 1 - 简单切片:\n")
	fmt.Printf("  smallSlice len=%d, cap=%d\n", len(smallSlice1), cap(smallSlice1))
	fmt.Printf("  问题: 仍然引用整个 1MB 的底层数组！\n")

	// 方式 2: 完整切片表达式
	smallSlice2 := largeSlice[:10:10]
	fmt.Printf("\n方式 2 - 完整切片表达式:\n")
	fmt.Printf("  smallSlice len=%d, cap=%d\n", len(smallSlice2), cap(smallSlice2))
	fmt.Printf("  改进: 限制了容量，但仍引用原数组\n")

	// 方式 3: 复制到新 slice（最佳）
	smallSlice3 := make([]byte, 10)
	copy(smallSlice3, largeSlice[:10])
	fmt.Printf("\n方式 3 - 复制到新 slice:\n")
	fmt.Printf("  smallSlice len=%d, cap=%d\n", len(smallSlice3), cap(smallSlice3))
	fmt.Printf("  最佳: 完全独立，原数组可以被垃圾回收\n")

	fmt.Println("\n结论:")
	fmt.Println("- 如果 largeSlice 不再使用，方式 1 和 2 会导致内存泄漏")
	fmt.Println("- 方式 3 允许 GC 回收大数组，节省内存")

	fmt.Println()
}

func demonstrateReuseArray() {
	fmt.Println("--- 5. 复用底层数组 ---")

	n := 100000

	// 每次分配新数组
	start := time.Now()
	for i := 0; i < 100; i++ {
		s := make([]int, 0, n)
		for j := 0; j < n; j++ {
			s = append(s, j)
		}
		// 使用 s...
		_ = s
	}
	duration1 := time.Since(start)

	// 复用底层数组
	start = time.Now()
	s := make([]int, 0, n)
	for i := 0; i < 100; i++ {
		s = s[:0] // 重置长度，保留容量
		for j := 0; j < n; j++ {
			s = append(s, j)
		}
		// 使用 s...
		_ = s
	}
	duration2 := time.Since(start)

	fmt.Printf("循环 100 次，每次填充 %d 个元素:\n", n)
	fmt.Printf("  每次新建:   %v\n", duration1)
	fmt.Printf("  复用数组:   %v (提速 %.2fx)\n", duration2, float64(duration1)/float64(duration2))
	fmt.Println("\n结论: 复用 slice 可以减少内存分配")

	// 实际应用：批处理
	fmt.Println("\n实际应用示例:")
	demonstrateBatchProcessing()

	fmt.Println()
}

func demonstrateBatchProcessing() {
	fmt.Println("批处理优化:")

	data := make([]int, 10000)
	for i := range data {
		data[i] = i
	}

	// 不好的做法
	start := time.Now()
	for i := 0; i < len(data); i += 100 {
		end := i + 100
		if end > len(data) {
			end = len(data)
		}
		batch := make([]int, end-i) // 每次分配新内存
		copy(batch, data[i:end])
		processBatch(batch)
	}
	duration1 := time.Since(start)

	// 好的做法
	start = time.Now()
	batch := make([]int, 100) // 预分配一次
	for i := 0; i < len(data); i += 100 {
		end := i + 100
		if end > len(data) {
			end = len(data)
		}
		n := copy(batch, data[i:end]) // 复用同一个 slice
		processBatch(batch[:n])
	}
	duration2 := time.Since(start)

	fmt.Printf("  每次分配: %v\n", duration1)
	fmt.Printf("  复用缓冲: %v (提速 %.2fx)\n", duration2, float64(duration1)/float64(duration2))
}

// === 辅助函数 ===

func filterNewSlice(s []int, fn func(int) bool) []int {
	result := make([]int, 0, len(s))
	for _, v := range s {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}

func filterInPlace(s []int, fn func(int) bool) []int {
	n := 0
	for _, v := range s {
		if fn(v) {
			s[n] = v
			n++
		}
	}
	return s[:n]
}

func processBatch(batch []int) {
	// 模拟处理
	sum := 0
	for _, v := range batch {
		sum += v
	}
	_ = sum
}

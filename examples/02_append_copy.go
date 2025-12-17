package main

import "fmt"

// 02_append_copy.go - Append 和 Copy 操作详解

func main() {
	fmt.Println("=== Append 和 Copy 操作 ===\n")

	// 1. append 基础用法
	demonstrateAppend()

	// 2. append 和底层数组
	demonstrateAppendUnderlyingArray()

	// 3. copy 函数
	demonstrateCopy()

	// 4. append vs copy 性能对比
	demonstrateAppendVsCopy()
}

func demonstrateAppend() {
	fmt.Println("--- 1. append 基础用法 ---")

	// 追加单个元素
	s := []int{1, 2, 3}
	fmt.Printf("原始: %v\n", s)

	s = append(s, 4)
	fmt.Printf("append(s, 4): %v\n", s)

	// 追加多个元素
	s = append(s, 5, 6, 7)
	fmt.Printf("append(s, 5, 6, 7): %v\n", s)

	// 追加另一个 slice
	s2 := []int{8, 9, 10}
	s = append(s, s2...)  // 注意 ... 运算符
	fmt.Printf("append(s, s2...): %v\n", s)

	// 连接多个 slice
	s3 := []int{11, 12}
	s4 := []int{13, 14}
	s = append(s, append(s3, s4...)...)
	fmt.Printf("连接多个 slice: %v\n", s)

	fmt.Println()
}

func demonstrateAppendUnderlyingArray() {
	fmt.Println("--- 2. append 和底层数组 ---")

	// 案例 1: 有足够容量，不会重新分配
	fmt.Println("案例 1: 容量充足")
	s1 := make([]int, 3, 10)
	s1[0], s1[1], s1[2] = 1, 2, 3
	fmt.Printf("s1: %v, len=%d, cap=%d, ptr=%p\n", s1, len(s1), cap(s1), s1)

	s2 := s1  // s2 和 s1 共享底层数组
	fmt.Printf("s2: %v, len=%d, cap=%d, ptr=%p\n", s2, len(s2), cap(s2), s2)

	s1 = append(s1, 4)  // 容量足够，不会重新分配
	fmt.Printf("append 后 s1: %v, len=%d, cap=%d, ptr=%p\n", s1, len(s1), cap(s1), s1)
	fmt.Printf("append 后 s2: %v, len=%d, cap=%d, ptr=%p\n", s2, len(s2), cap(s2), s2)

	// 修改 s1 会影响共享部分
	s1[0] = 100
	fmt.Printf("修改 s1[0] 后 s2: %v\n", s2)

	// 案例 2: 容量不足，会重新分配
	fmt.Println("\n案例 2: 容量不足")
	s3 := []int{1, 2, 3}
	fmt.Printf("s3: %v, len=%d, cap=%d, ptr=%p\n", s3, len(s3), cap(s3), s3)

	s4 := s3
	fmt.Printf("s4: %v, len=%d, cap=%d, ptr=%p\n", s4, len(s4), cap(s4), s4)

	// append 导致扩容，重新分配底层数组
	s3 = append(s3, 4, 5, 6, 7)
	fmt.Printf("append 后 s3: %v, len=%d, cap=%d, ptr=%p\n", s3, len(s3), cap(s3), s3)
	fmt.Printf("append 后 s4: %v, len=%d, cap=%d, ptr=%p\n", s4, len(s4), cap(s4), s4)

	// 修改 s3 不再影响 s4
	s3[0] = 100
	fmt.Printf("修改 s3[0] 后 s4: %v (不受影响)\n", s4)

	fmt.Println()
}

func demonstrateCopy() {
	fmt.Println("--- 3. copy 函数 ---")

	// 基本用法
	src := []int{1, 2, 3, 4, 5}
	dst := make([]int, len(src))
	n := copy(dst, src)
	fmt.Printf("源: %v\n", src)
	fmt.Printf("目标: %v, 复制了 %d 个元素\n", dst, n)

	// 目标 slice 较小
	dst2 := make([]int, 3)
	n = copy(dst2, src)
	fmt.Printf("\n目标较小: %v, 复制了 %d 个元素\n", dst2, n)

	// 目标 slice 较大
	dst3 := make([]int, 10)
	n = copy(dst3, src)
	fmt.Printf("目标较大: %v, 复制了 %d 个元素\n", dst3, n)

	// copy 到自身（移动元素）
	s := []int{1, 2, 3, 4, 5}
	fmt.Printf("\n原始: %v\n", s)
	copy(s[2:], s[:3])  // 将前 3 个元素复制到索引 2 开始的位置
	fmt.Printf("copy(s[2:], s[:3]): %v\n", s)

	// 删除前 n 个元素
	s = []int{1, 2, 3, 4, 5}
	s = s[copy(s, s[2:]):]  // 删除前 2 个元素
	fmt.Printf("删除前 2 个元素: %v\n", s)

	// copy 创建独立副本
	fmt.Println("\ncopy 创建独立副本:")
	original := []int{1, 2, 3}
	duplicate := make([]int, len(original))
	copy(duplicate, original)
	duplicate[0] = 100
	fmt.Printf("原始: %v, 副本: %v (独立)\n", original, duplicate)

	fmt.Println()
}

func demonstrateAppendVsCopy() {
	fmt.Println("--- 4. append vs copy 使用场景 ---")

	// 使用 append 合并 slice
	s1 := []int{1, 2, 3}
	s2 := []int{4, 5, 6}
	result := append(s1, s2...)
	fmt.Printf("使用 append 合并: %v\n", result)
	fmt.Printf("s1 可能被修改: %v\n", s1)  // 如果容量足够，s1 会被修改

	// 使用 copy 合并 slice（保护原始数据）
	s3 := []int{1, 2, 3}
	s4 := []int{4, 5, 6}
	result2 := make([]int, len(s3)+len(s4))
	copy(result2, s3)
	copy(result2[len(s3):], s4)
	fmt.Printf("使用 copy 合并: %v\n", result2)
	fmt.Printf("s3 保持不变: %v\n", s3)

	// 性能建议
	fmt.Println("\n性能建议:")
	fmt.Println("- append: 适合动态增长，自动扩容")
	fmt.Println("- copy: 适合已知大小，安全复制")
	fmt.Println("- 预分配容量可以显著提高 append 性能")

	fmt.Println()
}

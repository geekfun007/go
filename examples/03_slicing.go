package main

import "fmt"

// 03_slicing.go - 切片操作详解

func main() {
	fmt.Println("=== 切片操作详解 ===\n")

	// 1. 基本切片表达式
	demonstrateBasicSlicing()

	// 2. 完整切片表达式
	demonstrateFullSlicing()

	// 3. 切片和底层数组的关系
	demonstrateSlicingAndArray()

	// 4. 切片陷阱
	demonstrateSlicingPitfalls()
}

func demonstrateBasicSlicing() {
	fmt.Println("--- 1. 基本切片表达式 s[low:high] ---")

	s := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	fmt.Printf("原始 slice: %v, len=%d, cap=%d\n\n", s, len(s), cap(s))

	// s[low:high] 包含 low，不包含 high
	fmt.Println("各种切片操作:")
	
	s1 := s[2:5]
	fmt.Printf("s[2:5]:   %v, len=%d, cap=%d\n", s1, len(s1), cap(s1))

	s2 := s[:5]
	fmt.Printf("s[:5]:    %v, len=%d, cap=%d\n", s2, len(s2), cap(s2))

	s3 := s[5:]
	fmt.Printf("s[5:]:    %v, len=%d, cap=%d\n", s3, len(s3), cap(s3))

	s4 := s[:]
	fmt.Printf("s[:]:     %v, len=%d, cap=%d\n", s4, len(s4), cap(s4))

	// 链式切片
	s5 := s[2:8][1:4]
	fmt.Printf("s[2:8][1:4]: %v, len=%d, cap=%d\n", s5, len(s5), cap(s5))

	fmt.Println()
}

func demonstrateFullSlicing() {
	fmt.Println("--- 2. 完整切片表达式 s[low:high:max] ---")

	s := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	fmt.Printf("原始 slice: %v, len=%d, cap=%d\n\n", s, len(s), cap(s))

	// s[low:high:max]
	// len = high - low
	// cap = max - low
	
	s1 := s[2:5:7]
	fmt.Printf("s[2:5:7]: %v, len=%d, cap=%d\n", s1, len(s1), cap(s1))
	fmt.Printf("  len = high - low = 5 - 2 = %d\n", len(s1))
	fmt.Printf("  cap = max - low = 7 - 2 = %d\n", cap(s1))

	s2 := s[1:3:3]
	fmt.Printf("\ns[1:3:3]: %v, len=%d, cap=%d\n", s2, len(s2), cap(s2))

	// 限制容量的好处：防止 append 影响原 slice
	fmt.Println("\n完整切片表达式的应用：")
	original := []int{1, 2, 3, 4, 5}
	
	// 不使用完整切片表达式
	sub1 := original[:2]
	fmt.Printf("sub1 := original[:2]: %v, len=%d, cap=%d\n", sub1, len(sub1), cap(sub1))
	sub1 = append(sub1, 100)
	fmt.Printf("append 后 original: %v (被修改了！)\n", original)

	// 使用完整切片表达式
	original = []int{1, 2, 3, 4, 5}
	sub2 := original[:2:2]
	fmt.Printf("\nsub2 := original[:2:2]: %v, len=%d, cap=%d\n", sub2, len(sub2), cap(sub2))
	sub2 = append(sub2, 100)
	fmt.Printf("append 后 original: %v (未被修改)\n", original)
	fmt.Printf("append 后 sub2: %v (重新分配了底层数组)\n", sub2)

	fmt.Println()
}

func demonstrateSlicingAndArray() {
	fmt.Println("--- 3. 切片和底层数组的关系 ---")

	// 从数组创建切片
	arr := [10]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	fmt.Printf("数组: %v\n", arr)

	s1 := arr[2:5]
	fmt.Printf("s1 := arr[2:5]: %v, len=%d, cap=%d\n", s1, len(s1), cap(s1))

	// 修改切片会影响数组
	s1[0] = 100
	fmt.Printf("修改 s1[0] 后数组: %v\n", arr)
	fmt.Printf("修改 s1[0] 后 s1: %v\n", s1)

	// 多个切片可以共享同一个底层数组
	s2 := arr[3:7]
	fmt.Printf("\ns2 := arr[3:7]: %v\n", s2)
	s2[0] = 200  // arr[3] = 200
	fmt.Printf("修改 s2[0] 后数组: %v\n", arr)
	fmt.Printf("修改 s2[0] 后 s1: %v (s1 也受影响)\n", s1)

	// 切片的切片
	s3 := s1[1:]  // 基于 s1 创建新切片
	fmt.Printf("\ns3 := s1[1:]: %v\n", s3)
	s3[0] = 300
	fmt.Printf("修改 s3[0] 后数组: %v\n", arr)
	fmt.Printf("修改 s3[0] 后 s1: %v\n", s1)
	fmt.Printf("修改 s3[0] 后 s2: %v\n", s2)

	fmt.Println()
}

func demonstrateSlicingPitfalls() {
	fmt.Println("--- 4. 切片陷阱 ---")

	// 陷阱 1: 切片共享底层数组
	fmt.Println("陷阱 1: 切片共享底层数组")
	s := []int{1, 2, 3, 4, 5}
	s1 := s[0:3]
	s2 := s[2:5]
	fmt.Printf("s:  %v\n", s)
	fmt.Printf("s1: %v\n", s1)
	fmt.Printf("s2: %v\n", s2)
	
	s1[2] = 100  // 同时影响 s 和 s2
	fmt.Printf("修改 s1[2] 后:\n")
	fmt.Printf("s:  %v\n", s)
	fmt.Printf("s1: %v\n", s1)
	fmt.Printf("s2: %v (第一个元素也被修改)\n", s2)

	// 陷阱 2: 大切片的小切片导致内存泄漏
	fmt.Println("\n陷阱 2: 大切片的小切片导致内存泄漏")
	fmt.Println("问题代码:")
	fmt.Println("  largeSlice := make([]byte, 1000000)")
	fmt.Println("  smallSlice := largeSlice[:10]  // 仍然引用整个大数组！")
	fmt.Println("\n解决方案 1: 使用 copy")
	fmt.Println("  smallSlice := make([]byte, 10)")
	fmt.Println("  copy(smallSlice, largeSlice[:10])")
	fmt.Println("\n解决方案 2: 使用完整切片表达式")
	fmt.Println("  smallSlice := largeSlice[:10:10]")

	// 陷阱 3: 循环中的切片引用
	fmt.Println("\n陷阱 3: 循环中的切片引用")
	data := []int{1, 2, 3, 4, 5}
	var slices [][]int
	
	// 错误方式（仅作演示）
	fmt.Println("注意: 不要在循环中引用同一个切片的不同视图并期望保存各自的状态")
	
	// 正确方式
	for i := 0; i < len(data); i++ {
		// 创建独立的副本
		slice := make([]int, i+1)
		copy(slice, data[:i+1])
		slices = append(slices, slice)
	}
	
	fmt.Println("正确的切片集合:")
	for i, s := range slices {
		fmt.Printf("  slices[%d]: %v\n", i, s)
	}

	// 陷阱 4: 对 nil slice 进行切片
	fmt.Println("\n陷阱 4: 对 nil slice 进行切片")
	var nilSlice []int
	fmt.Printf("nilSlice: %v, len=%d, cap=%d\n", nilSlice, len(nilSlice), cap(nilSlice))
	
	emptySlice := nilSlice[0:0]
	fmt.Printf("nilSlice[0:0]: %v, len=%d, cap=%d\n", emptySlice, len(emptySlice), cap(emptySlice))
	fmt.Println("对 nil slice 切片 [0:0] 是安全的，返回 nil slice")
	
	// 但这会 panic:
	// badSlice := nilSlice[1:2]  // panic: runtime error: slice bounds out of range

	fmt.Println()
}

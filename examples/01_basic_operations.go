package main

import "fmt"

// 01_basic_operations.go - Slice 基础操作示例

func main() {
	fmt.Println("=== Go Slice 基础操作 ===\n")

	// 1. 声明和初始化
	demonstrateDeclaration()

	// 2. 访问和修改元素
	demonstrateAccess()

	// 3. 长度和容量
	demonstrateLenCap()

	// 4. nil slice vs empty slice
	demonstrateNilVsEmpty()
}

func demonstrateDeclaration() {
	fmt.Println("--- 1. 声明和初始化 ---")

	// 方式 1: 声明 nil slice
	var s1 []int
	fmt.Printf("nil slice: %v, len=%d, cap=%d, is nil? %v\n", 
		s1, len(s1), cap(s1), s1 == nil)

	// 方式 2: 使用 make 创建
	s2 := make([]int, 5)      // len=5, cap=5
	s3 := make([]int, 3, 10)  // len=3, cap=10
	fmt.Printf("make([]int, 5): %v, len=%d, cap=%d\n", s2, len(s2), cap(s2))
	fmt.Printf("make([]int, 3, 10): %v, len=%d, cap=%d\n", s3, len(s3), cap(s3))

	// 方式 3: 字面量初始化
	s4 := []int{1, 2, 3, 4, 5}
	fmt.Printf("字面量: %v, len=%d, cap=%d\n", s4, len(s4), cap(s4))

	// 方式 4: 从数组创建
	arr := [5]int{10, 20, 30, 40, 50}
	s5 := arr[1:4]
	fmt.Printf("从数组创建: %v, len=%d, cap=%d\n", s5, len(s5), cap(s5))

	fmt.Println()
}

func demonstrateAccess() {
	fmt.Println("--- 2. 访问和修改元素 ---")

	s := []string{"Go", "Python", "Java", "JavaScript", "Rust"}
	
	// 访问元素
	fmt.Printf("原始 slice: %v\n", s)
	fmt.Printf("第一个元素: %s\n", s[0])
	fmt.Printf("最后一个元素: %s\n", s[len(s)-1])
	fmt.Printf("第三个元素: %s\n", s[2])

	// 修改元素
	s[1] = "C++"
	s[3] = "TypeScript"
	fmt.Printf("修改后: %v\n", s)

	// 遍历
	fmt.Println("\n使用 range 遍历:")
	for i, lang := range s {
		fmt.Printf("  [%d] %s\n", i, lang)
	}

	// 只遍历值
	fmt.Println("\n只遍历值:")
	for _, lang := range s {
		fmt.Printf("  - %s\n", lang)
	}

	fmt.Println()
}

func demonstrateLenCap() {
	fmt.Println("--- 3. 长度和容量 ---")

	// 长度 vs 容量
	s := make([]int, 3, 10)
	fmt.Printf("初始: len=%d, cap=%d, %v\n", len(s), cap(s), s)

	// 修改现有元素
	s[0] = 100
	s[1] = 200
	s[2] = 300
	fmt.Printf("修改后: len=%d, cap=%d, %v\n", len(s), cap(s), s)

	// append 会增加长度
	s = append(s, 400)
	fmt.Printf("append 一个元素: len=%d, cap=%d, %v\n", len(s), cap(s), s)

	s = append(s, 500, 600, 700)
	fmt.Printf("append 多个元素: len=%d, cap=%d, %v\n", len(s), cap(s), s)

	// 当长度达到容量时，append 会扩容
	for i := 8; i <= 11; i++ {
		s = append(s, i*100)
		fmt.Printf("append %d: len=%d, cap=%d\n", i*100, len(s), cap(s))
	}

	fmt.Println()
}

func demonstrateNilVsEmpty() {
	fmt.Println("--- 4. nil slice vs empty slice ---")

	var nilSlice []int
	emptySlice1 := []int{}
	emptySlice2 := make([]int, 0)

	fmt.Printf("nil slice:     %v, len=%d, cap=%d, is nil? %v\n", 
		nilSlice, len(nilSlice), cap(nilSlice), nilSlice == nil)
	fmt.Printf("empty slice 1: %v, len=%d, cap=%d, is nil? %v\n", 
		emptySlice1, len(emptySlice1), cap(emptySlice1), emptySlice1 == nil)
	fmt.Printf("empty slice 2: %v, len=%d, cap=%d, is nil? %v\n", 
		emptySlice2, len(emptySlice2), cap(emptySlice2), emptySlice2 == nil)

	// 都可以 append
	nilSlice = append(nilSlice, 1)
	emptySlice1 = append(emptySlice1, 1)
	emptySlice2 = append(emptySlice2, 1)

	fmt.Println("\n追加元素后:")
	fmt.Printf("nil slice:     %v, is nil? %v\n", nilSlice, nilSlice == nil)
	fmt.Printf("empty slice 1: %v, is nil? %v\n", emptySlice1, emptySlice1 == nil)
	fmt.Printf("empty slice 2: %v, is nil? %v\n", emptySlice2, emptySlice2 == nil)

	fmt.Println()
}

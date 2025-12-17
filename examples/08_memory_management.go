package main

import (
	"fmt"
	"runtime"
)

// 08_memory_management.go - Slice 内存管理详解

func main() {
	fmt.Println("=== Slice 内存管理详解 ===\n")

	// 1. Slice 的内存结构
	demonstrateMemoryStructure()

	// 2. 扩容机制
	demonstrateGrowth()

	// 3. 内存泄漏场景
	demonstrateMemoryLeaks()

	// 4. 避免内存泄漏的方法
	demonstrateAvoidLeaks()

	// 5. 内存使用监控
	demonstrateMemoryMonitoring()
}

func demonstrateMemoryStructure() {
	fmt.Println("--- 1. Slice 的内存结构 ---")

	fmt.Println("Slice 内部结构:")
	fmt.Println("  type slice struct {")
	fmt.Println("      array unsafe.Pointer  // 指向底层数组")
	fmt.Println("      len   int             // 当前长度")
	fmt.Println("      cap   int             // 容量")
	fmt.Println("  }")

	// 示例
	s := make([]int, 3, 5)
	s[0], s[1], s[2] = 10, 20, 30

	fmt.Printf("\ns := make([]int, 3, 5)\n")
	fmt.Printf("s = %v\n", s)
	fmt.Printf("len(s) = %d\n", len(s))
	fmt.Printf("cap(s) = %d\n", cap(s))
	fmt.Printf("底层数组地址: %p\n", s)

	// 切片共享底层数组
	s2 := s[1:3]
	fmt.Printf("\ns2 := s[1:3]\n")
	fmt.Printf("s2 = %v\n", s2)
	fmt.Printf("len(s2) = %d\n", len(s2))
	fmt.Printf("cap(s2) = %d\n", cap(s2))
	fmt.Printf("底层数组地址: %p\n", s2)

	fmt.Println("\n结论: s 和 s2 共享底层数组（地址相同）")

	fmt.Println()
}

func demonstrateGrowth() {
	fmt.Println("--- 2. 扩容机制 ---")

	fmt.Println("Go 的 slice 扩容策略:")
	fmt.Println("  - 当 cap < 256: 容量翻倍")
	fmt.Println("  - 当 cap >= 256: 容量增长约 1.25 倍")
	fmt.Println("  - 实际容量会向上取整到合适的内存大小\n")

	s := make([]int, 0, 1)
	fmt.Printf("初始: len=%d, cap=%d, addr=%p\n", len(s), cap(s), s)

	for i := 0; i < 20; i++ {
		prevCap := cap(s)
		prevAddr := fmt.Sprintf("%p", s)
		s = append(s, i)
		newCap := cap(s)
		newAddr := fmt.Sprintf("%p", s)

		if newCap != prevCap {
			growthRatio := float64(newCap) / float64(prevCap)
			addrChanged := prevAddr != newAddr
			fmt.Printf("扩容: len=%2d, cap: %2d -> %2d (%.2fx), 地址变化: %v\n",
				len(s), prevCap, newCap, growthRatio, addrChanged)
		}
	}

	// 大容量扩容
	fmt.Println("\n大容量扩容:")
	largeSlice := make([]int, 0, 200)
	for i := 0; i < 600; i++ {
		prevCap := cap(largeSlice)
		largeSlice = append(largeSlice, i)
		newCap := cap(largeSlice)

		if newCap != prevCap {
			growthRatio := float64(newCap) / float64(prevCap)
			fmt.Printf("len=%3d, cap: %3d -> %3d (%.2fx)\n",
				len(largeSlice), prevCap, newCap, growthRatio)
		}
	}

	fmt.Println()
}

func demonstrateMemoryLeaks() {
	fmt.Println("--- 3. 内存泄漏场景 ---")

	// 场景 1: 大 slice 的小切片
	fmt.Println("场景 1: 从大 slice 中提取小切片")
	fmt.Println("问题代码:")
	fmt.Println("  largeSlice := make([]byte, 1000000)  // 1MB")
	fmt.Println("  smallSlice := largeSlice[:10]        // 只需要 10 字节")
	fmt.Println("  return smallSlice                     // 但整个 1MB 无法释放！")

	// 场景 2: 切片作为缓存
	fmt.Println("\n场景 2: 切片作为缓存")
	cache := make([][]byte, 0)
	for i := 0; i < 5; i++ {
		data := make([]byte, 100000) // 100KB
		cache = append(cache, data[:1000])
	}
	fmt.Printf("缓存了 %d 个小 slice，但占用了约 %.2f MB\n",
		len(cache), float64(5*100000)/1024/1024)

	// 场景 3: 循环中的切片累积
	fmt.Println("\n场景 3: 循环中的切片累积")
	var results [][]int
	bigData := make([]int, 10000)
	for i := 0; i < 100; i++ {
		// 问题：每个结果都引用同一个大数组
		result := bigData[i : i+10]
		results = append(results, result)
	}
	fmt.Printf("results 中有 %d 个 slice，但都引用同一个大数组\n", len(results))

	fmt.Println()
}

func demonstrateAvoidLeaks() {
	fmt.Println("--- 4. 避免内存泄漏的方法 ---")

	// 方法 1: 使用 copy 创建独立副本
	fmt.Println("方法 1: 使用 copy 创建独立副本")
	largeSlice := make([]byte, 1000000)
	for i := range largeSlice {
		largeSlice[i] = byte(i % 256)
	}

	// 不好的做法
	bad := largeSlice[:10]
	fmt.Printf("  不好: len=%d, cap=%d (引用整个大数组)\n", len(bad), cap(bad))

	// 好的做法
	good := make([]byte, 10)
	copy(good, largeSlice[:10])
	fmt.Printf("  好: len=%d, cap=%d (独立副本)\n", len(good), cap(good))

	// 方法 2: 使用完整切片表达式
	fmt.Println("\n方法 2: 使用完整切片表达式限制容量")
	s := make([]int, 100)
	sub1 := s[:10]
	sub2 := s[:10:10]
	fmt.Printf("  s[:10]:    len=%d, cap=%d\n", len(sub1), cap(sub1))
	fmt.Printf("  s[:10:10]: len=%d, cap=%d (容量被限制)\n", len(sub2), cap(sub2))

	// 方法 3: 及时释放不需要的引用
	fmt.Println("\n方法 3: 及时释放不需要的引用")
	fmt.Println("  处理完大数据后，将 slice 设置为 nil:")
	fmt.Println("  largeSlice = nil")
	fmt.Println("  runtime.GC()  // 手动触发 GC（通常不需要）")

	// 方法 4: 使用 append 时注意容量
	fmt.Println("\n方法 4: 使用 append 时注意容量")
	original := make([]int, 3, 10)
	original[0], original[1], original[2] = 1, 2, 3

	// 不好：可能修改原 slice
	derived1 := original[:2]
	derived1 = append(derived1, 999)
	fmt.Printf("  原 slice 被修改: %v\n", original)

	// 好：使用完整切片表达式
	original = []int{1, 2, 3}
	original = append(original, 4, 5, 6, 7) // 增加容量
	derived2 := original[:2:2]
	derived2 = append(derived2, 999)
	fmt.Printf("  原 slice 未被修改: %v\n", original)

	fmt.Println()
}

func demonstrateMemoryMonitoring() {
	fmt.Println("--- 5. 内存使用监控 ---")

	// 获取内存统计
	var m1 runtime.MemStats
	runtime.ReadMemStats(&m1)

	fmt.Printf("分配前:\n")
	printMemStats(&m1)

	// 分配大量内存
	slices := make([]*[]byte, 100)
	for i := range slices {
		s := make([]byte, 1024*1024) // 1MB
		slices[i] = &s
	}

	var m2 runtime.MemStats
	runtime.ReadMemStats(&m2)

	fmt.Printf("\n分配 100MB 后:\n")
	printMemStats(&m2)

	fmt.Printf("\n内存增长: %.2f MB\n",
		float64(m2.Alloc-m1.Alloc)/1024/1024)

	// 释放内存
	slices = nil
	runtime.GC()

	var m3 runtime.MemStats
	runtime.ReadMemStats(&m3)

	fmt.Printf("\nGC 后:\n")
	printMemStats(&m3)

	fmt.Printf("\n回收的内存: %.2f MB\n",
		float64(m2.Alloc-m3.Alloc)/1024/1024)

	// 最佳实践总结
	fmt.Println("\n=== 最佳实践总结 ===")
	fmt.Println("1. ✅ 预分配足够的容量，避免频繁扩容")
	fmt.Println("2. ✅ 从大 slice 提取小片段时，使用 copy 创建独立副本")
	fmt.Println("3. ✅ 使用完整切片表达式 [low:high:max] 限制容量")
	fmt.Println("4. ✅ 及时释放不再使用的大 slice 引用")
	fmt.Println("5. ✅ 注意 slice 是引用类型，修改会影响共享底层数组的其他 slice")
	fmt.Println("6. ❌ 避免在长期存活的对象中保留大 slice 的小切片")
	fmt.Println("7. ❌ 避免无意义的 slice 复制，直接传递 slice（引用传递）")

	fmt.Println()
}

// === 辅助函数 ===

func printMemStats(m *runtime.MemStats) {
	fmt.Printf("  已分配内存: %.2f MB\n", float64(m.Alloc)/1024/1024)
	fmt.Printf("  总分配内存: %.2f MB\n", float64(m.TotalAlloc)/1024/1024)
	fmt.Printf("  系统内存:   %.2f MB\n", float64(m.Sys)/1024/1024)
	fmt.Printf("  GC 次数:    %d\n", m.NumGC)
}

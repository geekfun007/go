// 示例 9：Pull 迭代器（推模式转拉模式）
// 演示如何将推模式迭代器转换为拉模式
//
// 推模式（Push）：迭代器主动推送值给消费者（for-range 默认方式）
// 拉模式（Pull）：消费者主动拉取值（类似传统迭代器）
package main

import (
	"fmt"
)

// Seq 单值迭代器类型
type Seq[T any] func(yield func(T) bool)

// Seq2 双值迭代器类型
type Seq2[K, V any] func(yield func(K, V) bool)

// Pull 将推模式迭代器转换为拉模式
// 返回：
//   - next: 获取下一个值的函数，返回 (值, 是否有效)
//   - stop: 停止迭代的函数（必须调用以释放资源）
//
// 注意：Go 1.23 中这个功能由 iter.Pull 提供
func Pull[T any](seq Seq[T]) (next func() (T, bool), stop func()) {
	// 使用 channel 在 goroutine 之间传递值
	ch := make(chan T)
	done := make(chan struct{})
	stopped := make(chan struct{})

	go func() {
		defer close(ch)
		seq(func(v T) bool {
			select {
			case ch <- v:
				return true
			case <-done:
				return false
			}
		})
	}()

	next = func() (T, bool) {
		select {
		case v, ok := <-ch:
			return v, ok
		case <-stopped:
			var zero T
			return zero, false
		}
	}

	stop = func() {
		select {
		case <-stopped:
			// 已经停止
		default:
			close(stopped)
			close(done)
		}
	}

	return next, stop
}

// Pull2 将双值推模式迭代器转换为拉模式
func Pull2[K, V any](seq Seq2[K, V]) (next func() (K, V, bool), stop func()) {
	type pair struct {
		k K
		v V
	}
	ch := make(chan pair)
	done := make(chan struct{})
	stopped := make(chan struct{})

	go func() {
		defer close(ch)
		seq(func(k K, v V) bool {
			select {
			case ch <- pair{k, v}:
				return true
			case <-done:
				return false
			}
		})
	}()

	next = func() (K, V, bool) {
		select {
		case p, ok := <-ch:
			return p.k, p.v, ok
		case <-stopped:
			var zeroK K
			var zeroV V
			return zeroK, zeroV, false
		}
	}

	stop = func() {
		select {
		case <-stopped:
		default:
			close(stopped)
			close(done)
		}
	}

	return next, stop
}

// ========== 示例迭代器 ==========

// Range 数字范围迭代器
func Range(start, end int) Seq[int] {
	return func(yield func(int) bool) {
		for i := start; i < end; i++ {
			if !yield(i) {
				return
			}
		}
	}
}

// Enumerate 带索引的切片迭代器
func Enumerate[T any](slice []T) Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for i, v := range slice {
			if !yield(i, v) {
				return
			}
		}
	}
}

// Fibonacci 斐波那契数列迭代器
func Fibonacci() Seq[int] {
	return func(yield func(int) bool) {
		a, b := 0, 1
		for {
			if !yield(a) {
				return
			}
			a, b = b, a+b
		}
	}
}

func main() {
	fmt.Println("=== Pull 迭代器示例 ===")
	fmt.Println()

	// 示例 1: 基本 Pull 使用
	fmt.Println("1. 基本 Pull 使用:")
	fmt.Println("   推模式 (for-range):")
	for v := range Range(1, 5) {
		fmt.Printf("     %d\n", v)
	}

	fmt.Println("   拉模式 (Pull):")
	next, stop := Pull(Range(1, 5))
	defer stop()
	for {
		v, ok := next()
		if !ok {
			break
		}
		fmt.Printf("     %d\n", v)
	}
	fmt.Println()

	// 示例 2: 交替拉取两个迭代器
	fmt.Println("2. 交替拉取两个迭代器:")
	next1, stop1 := Pull(Range(1, 5))
	defer stop1()
	next2, stop2 := Pull(Range(100, 105))
	defer stop2()

	for {
		v1, ok1 := next1()
		v2, ok2 := next2()
		if !ok1 && !ok2 {
			break
		}
		if ok1 {
			fmt.Printf("   迭代器1: %d\n", v1)
		}
		if ok2 {
			fmt.Printf("   迭代器2: %d\n", v2)
		}
	}
	fmt.Println()

	// 示例 3: 提前停止
	fmt.Println("3. 提前停止 (只取前 3 个):")
	next3, stop3 := Pull(Range(1, 100))
	for i := 0; i < 3; i++ {
		v, ok := next3()
		if !ok {
			break
		}
		fmt.Printf("   %d\n", v)
	}
	stop3() // 重要：必须调用 stop 释放资源
	fmt.Println()

	// 示例 4: Pull2 双值迭代器
	fmt.Println("4. Pull2 双值迭代器:")
	fruits := []string{"apple", "banana", "cherry"}
	next4, stop4 := Pull2(Enumerate(fruits))
	defer stop4()

	for {
		i, v, ok := next4()
		if !ok {
			break
		}
		fmt.Printf("   [%d] %s\n", i, v)
	}
	fmt.Println()

	// 示例 5: 无限序列 + Pull
	fmt.Println("5. 斐波那契数列 (取前 10 个):")
	nextFib, stopFib := Pull(Fibonacci())
	defer stopFib()

	for i := 0; i < 10; i++ {
		v, _ := nextFib()
		fmt.Printf("   F(%d) = %d\n", i, v)
	}
	fmt.Println()

	// 示例 6: 合并两个有序序列
	fmt.Println("6. 合并两个有序序列:")
	seq1 := Range(1, 10)  // 1,2,3,4,5,6,7,8,9
	seq2 := Range(5, 15)  // 5,6,7,8,9,10,11,12,13,14

	// 使用 Pull 合并（类似归并排序的合并步骤）
	n1, s1 := Pull(seq1)
	n2, s2 := Pull(seq2)
	defer s1()
	defer s2()

	v1, ok1 := n1()
	v2, ok2 := n2()
	merged := []int{}

	for ok1 || ok2 {
		if !ok2 || (ok1 && v1 <= v2) {
			merged = append(merged, v1)
			v1, ok1 = n1()
		} else {
			merged = append(merged, v2)
			v2, ok2 = n2()
		}
	}
	fmt.Printf("   合并结果: %v\n", merged)
	fmt.Println()

	// 示例 7: 推模式 vs 拉模式对比
	fmt.Println("7. 推模式 vs 拉模式对比:")
	fmt.Println("   推模式特点:")
	fmt.Println("     - 使用 for-range 语法")
	fmt.Println("     - 迭代器控制流程")
	fmt.Println("     - 语法简洁")
	fmt.Println("     - 无额外 goroutine 开销")
	fmt.Println()
	fmt.Println("   拉模式特点:")
	fmt.Println("     - 使用 next()/stop() 函数")
	fmt.Println("     - 消费者控制流程")
	fmt.Println("     - 可以交替消费多个迭代器")
	fmt.Println("     - 有 goroutine 开销")
	fmt.Println("     - 必须记得调用 stop()")
}

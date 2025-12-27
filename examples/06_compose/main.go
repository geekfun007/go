// 示例 6：组合迭代器（函数式风格）
// 演示如何组合多个迭代器构建数据处理管道
package main

import (
	"fmt"
	"strings"
)

// Seq 单值迭代器类型
type Seq[T any] func(yield func(T) bool)

// ========== 基础生成器 ==========

// Naturals 自然数无限序列
func Naturals() Seq[int] {
	return func(yield func(int) bool) {
		for i := 1; ; i++ {
			if !yield(i) {
				return
			}
		}
	}
}

// Repeat 重复同一个值
func Repeat[T any](value T) Seq[T] {
	return func(yield func(T) bool) {
		for {
			if !yield(value) {
				return
			}
		}
	}
}

// FromSlice 从切片创建迭代器
func FromSlice[T any](slice []T) Seq[T] {
	return func(yield func(T) bool) {
		for _, v := range slice {
			if !yield(v) {
				return
			}
		}
	}
}

// Range 数字范围
func Range(start, end int) Seq[int] {
	return func(yield func(int) bool) {
		for i := start; i < end; i++ {
			if !yield(i) {
				return
			}
		}
	}
}

// ========== 转换操作 ==========

// Map 映射转换
func Map[T, U any](seq Seq[T], f func(T) U) Seq[U] {
	return func(yield func(U) bool) {
		seq(func(v T) bool {
			return yield(f(v))
		})
	}
}

// Filter 过滤
func Filter[T any](seq Seq[T], predicate func(T) bool) Seq[T] {
	return func(yield func(T) bool) {
		seq(func(v T) bool {
			if predicate(v) {
				return yield(v)
			}
			return true
		})
	}
}

// FlatMap 扁平化映射
func FlatMap[T, U any](seq Seq[T], f func(T) Seq[U]) Seq[U] {
	return func(yield func(U) bool) {
		seq(func(v T) bool {
			shouldContinue := true
			f(v)(func(u U) bool {
				if !yield(u) {
					shouldContinue = false
					return false
				}
				return true
			})
			return shouldContinue
		})
	}
}

// ========== 限制操作 ==========

// Take 取前 n 个
func Take[T any](seq Seq[T], n int) Seq[T] {
	return func(yield func(T) bool) {
		count := 0
		seq(func(v T) bool {
			if count >= n {
				return false
			}
			count++
			return yield(v)
		})
	}
}

// Skip 跳过前 n 个
func Skip[T any](seq Seq[T], n int) Seq[T] {
	return func(yield func(T) bool) {
		count := 0
		seq(func(v T) bool {
			count++
			if count <= n {
				return true
			}
			return yield(v)
		})
	}
}

// TakeWhile 取满足条件的元素
func TakeWhile[T any](seq Seq[T], predicate func(T) bool) Seq[T] {
	return func(yield func(T) bool) {
		seq(func(v T) bool {
			if !predicate(v) {
				return false
			}
			return yield(v)
		})
	}
}

// DropWhile 跳过满足条件的元素
func DropWhile[T any](seq Seq[T], predicate func(T) bool) Seq[T] {
	return func(yield func(T) bool) {
		dropping := true
		seq(func(v T) bool {
			if dropping && predicate(v) {
				return true
			}
			dropping = false
			return yield(v)
		})
	}
}

// ========== 聚合操作 ==========

// Collect 收集为切片
func Collect[T any](seq Seq[T]) []T {
	var result []T
	seq(func(v T) bool {
		result = append(result, v)
		return true
	})
	return result
}

// Reduce 归约
func Reduce[T, U any](seq Seq[T], initial U, f func(U, T) U) U {
	result := initial
	seq(func(v T) bool {
		result = f(result, v)
		return true
	})
	return result
}

// Count 计数
func Count[T any](seq Seq[T]) int {
	count := 0
	seq(func(v T) bool {
		count++
		return true
	})
	return count
}

// Any 检查是否有任意元素满足条件
func Any[T any](seq Seq[T], predicate func(T) bool) bool {
	found := false
	seq(func(v T) bool {
		if predicate(v) {
			found = true
			return false
		}
		return true
	})
	return found
}

// All 检查是否所有元素都满足条件
func All[T any](seq Seq[T], predicate func(T) bool) bool {
	allMatch := true
	seq(func(v T) bool {
		if !predicate(v) {
			allMatch = false
			return false
		}
		return true
	})
	return allMatch
}

// ========== 组合操作 ==========

// Chain 连接多个迭代器
func Chain[T any](seqs ...Seq[T]) Seq[T] {
	return func(yield func(T) bool) {
		for _, seq := range seqs {
			shouldContinue := true
			seq(func(v T) bool {
				if !yield(v) {
					shouldContinue = false
					return false
				}
				return true
			})
			if !shouldContinue {
				return
			}
		}
	}
}

// Unique 去重（基于 comparable）
func Unique[T comparable](seq Seq[T]) Seq[T] {
	return func(yield func(T) bool) {
		seen := make(map[T]struct{})
		seq(func(v T) bool {
			if _, ok := seen[v]; ok {
				return true
			}
			seen[v] = struct{}{}
			return yield(v)
		})
	}
}

func main() {
	fmt.Println("=== 组合迭代器示例 ===")
	fmt.Println()

	// 示例 1: 基本管道
	fmt.Println("1. 基本管道: 自然数 -> 取10个 -> 过滤偶数 -> 平方")
	result1 := Collect(
		Map(
			Filter(
				Take(Naturals(), 10),
				func(n int) bool { return n%2 == 0 },
			),
			func(n int) int { return n * n },
		),
	)
	fmt.Printf("  结果: %v\n", result1)
	fmt.Println()

	// 示例 2: 字符串处理管道
	fmt.Println("2. 字符串处理: 过滤空字符串 -> 转大写 -> 添加前缀")
	words := []string{"hello", "", "world", "go", "", "iterator"}
	result2 := Collect(
		Map(
			Map(
				Filter(FromSlice(words), func(s string) bool { return s != "" }),
				strings.ToUpper,
			),
			func(s string) string { return "[" + s + "]" },
		),
	)
	fmt.Printf("  结果: %v\n", result2)
	fmt.Println()

	// 示例 3: 数学运算
	fmt.Println("3. 计算前 100 个自然数中 3 的倍数的和:")
	sum := Reduce(
		Filter(Take(Naturals(), 100), func(n int) bool { return n%3 == 0 }),
		0,
		func(acc, n int) int { return acc + n },
	)
	fmt.Printf("  结果: %d\n", sum)
	fmt.Println()

	// 示例 4: FlatMap
	fmt.Println("4. FlatMap: 将每个数字展开为 [n, n*10, n*100]")
	result4 := Collect(
		FlatMap(Range(1, 4), func(n int) Seq[int] {
			return FromSlice([]int{n, n * 10, n * 100})
		}),
	)
	fmt.Printf("  结果: %v\n", result4)
	fmt.Println()

	// 示例 5: Chain
	fmt.Println("5. Chain: 连接多个序列")
	result5 := Collect(
		Chain(
			FromSlice([]int{1, 2, 3}),
			FromSlice([]int{10, 20, 30}),
			FromSlice([]int{100, 200, 300}),
		),
	)
	fmt.Printf("  结果: %v\n", result5)
	fmt.Println()

	// 示例 6: Unique
	fmt.Println("6. Unique: 去重")
	result6 := Collect(
		Unique(FromSlice([]int{1, 2, 2, 3, 3, 3, 4, 4, 4, 4})),
	)
	fmt.Printf("  结果: %v\n", result6)
	fmt.Println()

	// 示例 7: Any 和 All
	fmt.Println("7. Any 和 All:")
	nums := FromSlice([]int{2, 4, 6, 8, 10})
	hasOdd := Any(nums, func(n int) bool { return n%2 == 1 })
	allEven := All(nums, func(n int) bool { return n%2 == 0 })
	fmt.Printf("  [2,4,6,8,10] 有奇数: %v\n", hasOdd)
	fmt.Printf("  [2,4,6,8,10] 全是偶数: %v\n", allEven)
	fmt.Println()

	// 示例 8: DropWhile + TakeWhile
	fmt.Println("8. DropWhile + TakeWhile:")
	result8 := Collect(
		TakeWhile(
			DropWhile(Range(1, 20), func(n int) bool { return n < 5 }),
			func(n int) bool { return n < 15 },
		),
	)
	fmt.Printf("  跳过 <5，取 <15: %v\n", result8)
	fmt.Println()

	// 示例 9: 复杂管道 - 找出所有素数的前 10 个
	fmt.Println("9. 前 10 个素数:")
	isPrime := func(n int) bool {
		if n < 2 {
			return false
		}
		for i := 2; i*i <= n; i++ {
			if n%i == 0 {
				return false
			}
		}
		return true
	}
	primes := Collect(Take(Filter(Naturals(), isPrime), 10))
	fmt.Printf("  结果: %v\n", primes)
}

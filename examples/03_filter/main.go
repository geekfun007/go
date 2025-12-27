// 示例 3：过滤迭代器
// 演示如何创建过滤、映射等高阶迭代器
package main

import "fmt"

// Seq 是单值迭代器的类型别名
type Seq[T any] func(yield func(T) bool)

// Range 创建数字范围迭代器
func Range(start, end int) Seq[int] {
	return func(yield func(int) bool) {
		for i := start; i < end; i++ {
			if !yield(i) {
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

// Filter 过滤迭代器，只保留满足条件的元素
func Filter[T any](seq Seq[T], predicate func(T) bool) Seq[T] {
	return func(yield func(T) bool) {
		seq(func(v T) bool {
			if predicate(v) {
				return yield(v)
			}
			return true // 继续迭代，但不 yield 此值
		})
	}
}

// Map 映射迭代器，将每个元素转换为新值
func Map[T, U any](seq Seq[T], transform func(T) U) Seq[U] {
	return func(yield func(U) bool) {
		seq(func(v T) bool {
			return yield(transform(v))
		})
	}
}

// Take 只取前 n 个元素
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

// Skip 跳过前 n 个元素
func Skip[T any](seq Seq[T], n int) Seq[T] {
	return func(yield func(T) bool) {
		count := 0
		seq(func(v T) bool {
			count++
			if count <= n {
				return true // 继续但不 yield
			}
			return yield(v)
		})
	}
}

// TakeWhile 取元素直到条件不满足
func TakeWhile[T any](seq Seq[T], predicate func(T) bool) Seq[T] {
	return func(yield func(T) bool) {
		seq(func(v T) bool {
			if !predicate(v) {
				return false // 停止迭代
			}
			return yield(v)
		})
	}
}

// Collect 将迭代器收集为切片
func Collect[T any](seq Seq[T]) []T {
	var result []T
	seq(func(v T) bool {
		result = append(result, v)
		return true
	})
	return result
}

func main() {
	fmt.Println("=== 过滤迭代器示例 ===")
	fmt.Println()

	// 示例 1: Filter - 过滤偶数
	fmt.Println("1. Filter - 过滤出 1-10 中的偶数:")
	evenNumbers := Filter(Range(1, 11), func(n int) bool {
		return n%2 == 0
	})
	for n := range evenNumbers {
		fmt.Printf("  %d\n", n)
	}
	fmt.Println()

	// 示例 2: Map - 平方
	fmt.Println("2. Map - 将 1-5 映射为平方:")
	squares := Map(Range(1, 6), func(n int) int {
		return n * n
	})
	for n := range squares {
		fmt.Printf("  %d\n", n)
	}
	fmt.Println()

	// 示例 3: Take
	fmt.Println("3. Take - 取前 3 个:")
	first3 := Take(Range(1, 100), 3)
	for n := range first3 {
		fmt.Printf("  %d\n", n)
	}
	fmt.Println()

	// 示例 4: Skip
	fmt.Println("4. Skip - 跳过前 5 个，取后面的:")
	skipped := Take(Skip(Range(1, 20), 5), 5)
	for n := range skipped {
		fmt.Printf("  %d\n", n)
	}
	fmt.Println()

	// 示例 5: 组合使用
	fmt.Println("5. 组合: 1-100 中能被 3 整除的数，取前 5 个，然后平方:")
	result := Map(
		Take(
			Filter(Range(1, 100), func(n int) bool { return n%3 == 0 }),
			5,
		),
		func(n int) int { return n * n },
	)
	for n := range result {
		fmt.Printf("  %d\n", n) // 9, 36, 81, 144, 225
	}
	fmt.Println()

	// 示例 6: TakeWhile
	fmt.Println("6. TakeWhile - 取数字直到大于等于 5:")
	lessThan5 := TakeWhile(Range(1, 100), func(n int) bool {
		return n < 5
	})
	for n := range lessThan5 {
		fmt.Printf("  %d\n", n)
	}
	fmt.Println()

	// 示例 7: Collect
	fmt.Println("7. Collect - 将迭代器收集为切片:")
	slice := Collect(Filter(Range(1, 10), func(n int) bool { return n%2 == 0 }))
	fmt.Printf("  结果切片: %v\n", slice)
	fmt.Println()

	// 示例 8: 字符串处理
	fmt.Println("8. 字符串过滤 - 长度大于 4 的单词:")
	words := []string{"go", "rust", "python", "java", "typescript", "c"}
	longWords := Filter(FromSlice(words), func(s string) bool {
		return len(s) > 4
	})
	for w := range longWords {
		fmt.Printf("  %s\n", w)
	}
}

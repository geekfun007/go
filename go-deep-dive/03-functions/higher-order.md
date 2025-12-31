# 高阶函数 / Higher-Order Functions

```go
package main

import (
    "fmt"
    "sort"
    "strings"
)

// Map 函数 / Map function
func Map[T, U any](slice []T, fn func(T) U) []U {
    result := make([]U, len(slice))
    for i, v := range slice {
        result[i] = fn(v)
    }
    return result
}

// Filter 函数 / Filter function
func Filter[T any](slice []T, predicate func(T) bool) []T {
    result := make([]T, 0)
    for _, v := range slice {
        if predicate(v) {
            result = append(result, v)
        }
    }
    return result
}

// Reduce 函数 / Reduce function
func Reduce[T, U any](slice []T, initial U, fn func(U, T) U) U {
    result := initial
    for _, v := range slice {
        result = fn(result, v)
    }
    return result
}

// ForEach 函数 / ForEach function
func ForEach[T any](slice []T, fn func(T)) {
    for _, v := range slice {
        fn(v)
    }
}

// Find 函数 / Find function
func Find[T any](slice []T, predicate func(T) bool) (T, bool) {
    for _, v := range slice {
        if predicate(v) {
            return v, true
        }
    }
    var zero T
    return zero, false
}

// Any 函数 / Any function
func Any[T any](slice []T, predicate func(T) bool) bool {
    for _, v := range slice {
        if predicate(v) {
            return true
        }
    }
    return false
}

// All 函数 / All function
func All[T any](slice []T, predicate func(T) bool) bool {
    for _, v := range slice {
        if !predicate(v) {
            return false
        }
    }
    return true
}

func main() {
    nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    
    // Map: 将每个元素乘以2
    doubled := Map(nums, func(n int) int { return n * 2 })
    fmt.Println("Doubled:", doubled)
    
    // Filter: 筛选偶数
    evens := Filter(nums, func(n int) bool { return n%2 == 0 })
    fmt.Println("Evens:", evens)
    
    // Reduce: 求和
    sum := Reduce(nums, 0, func(acc, n int) int { return acc + n })
    fmt.Println("Sum:", sum)
    
    // ForEach: 打印每个元素
    fmt.Print("ForEach: ")
    ForEach(nums, func(n int) { fmt.Print(n, " ") })
    fmt.Println()
    
    // Find: 找第一个大于5的数
    found, ok := Find(nums, func(n int) bool { return n > 5 })
    fmt.Println("Find >5:", found, ok)
    
    // Any: 是否有大于5的数
    fmt.Println("Any >5:", Any(nums, func(n int) bool { return n > 5 }))
    
    // All: 是否全部大于0
    fmt.Println("All >0:", All(nums, func(n int) bool { return n > 0 }))
    
    // 链式操作 / Chained operations
    result := Reduce(
        Filter(
            Map(nums, func(n int) int { return n * n }),  // 平方
            func(n int) bool { return n > 10 },           // 大于10
        ),
        0,
        func(acc, n int) int { return acc + n },  // 求和
    )
    fmt.Println("Chained result:", result)
    
    // 使用标准库的排序 / Using standard library sort
    strs := []string{"banana", "apple", "cherry", "date"}
    sort.Slice(strs, func(i, j int) bool {
        return strings.ToLower(strs[i]) < strings.ToLower(strs[j])
    })
    fmt.Println("Sorted:", strs)
}
```

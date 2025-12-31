# 迭代器与 slices 包 / Iterator & slices Package (Go 1.21+/1.23+)

## 1. slices 包详解 / slices Package Details (Go 1.21+)

```go
package main

import (
    "cmp"
    "fmt"
    "slices"
)

func main() {
    // ========================================
    // 基本操作 / Basic Operations
    // ========================================
    
    s := []int{3, 1, 4, 1, 5, 9, 2, 6}
    
    // Clone - 复制切片 / Clone slice
    cloned := slices.Clone(s)
    fmt.Println("Clone:", cloned)
    
    // Equal - 比较切片是否相等 / Compare slices
    s1 := []int{1, 2, 3}
    s2 := []int{1, 2, 3}
    s3 := []int{1, 2, 4}
    fmt.Println("Equal s1==s2:", slices.Equal(s1, s2))  // true
    fmt.Println("Equal s1==s3:", slices.Equal(s1, s3))  // false
    
    // EqualFunc - 自定义比较函数 / Custom comparison
    strings1 := []string{"hello", "world"}
    strings2 := []string{"HELLO", "WORLD"}
    equalIgnoreCase := slices.EqualFunc(strings1, strings2, func(a, b string) bool {
        return strings.EqualFold(a, b)
    })
    fmt.Println("EqualFunc (ignore case):", equalIgnoreCase)  // true
    
    // Compare - 比较切片 (返回 -1, 0, 1) / Compare slices
    fmt.Println("Compare [1,2] vs [1,3]:", slices.Compare([]int{1, 2}, []int{1, 3}))  // -1
    fmt.Println("Compare [1,3] vs [1,2]:", slices.Compare([]int{1, 3}, []int{1, 2}))  // 1
    fmt.Println("Compare [1,2] vs [1,2]:", slices.Compare([]int{1, 2}, []int{1, 2}))  // 0
    
    // ========================================
    // 查找操作 / Search Operations
    // ========================================
    
    nums := []int{10, 20, 30, 40, 50}
    
    // Index - 查找元素索引 / Find element index
    fmt.Println("Index of 30:", slices.Index(nums, 30))   // 2
    fmt.Println("Index of 99:", slices.Index(nums, 99))   // -1
    
    // IndexFunc - 使用函数查找 / Find with function
    idx := slices.IndexFunc(nums, func(n int) bool {
        return n > 25
    })
    fmt.Println("IndexFunc (>25):", idx)  // 2
    
    // Contains - 检查是否包含元素 / Check if contains
    fmt.Println("Contains 30:", slices.Contains(nums, 30))  // true
    fmt.Println("Contains 99:", slices.Contains(nums, 99))  // false
    
    // ContainsFunc - 使用函数检查 / Check with function
    hasEven := slices.ContainsFunc(nums, func(n int) bool {
        return n%2 == 0
    })
    fmt.Println("ContainsFunc (even):", hasEven)  // true
    
    // ========================================
    // 排序操作 / Sorting Operations
    // ========================================
    
    unsorted := []int{3, 1, 4, 1, 5, 9, 2, 6}
    
    // Sort - 排序 (原地修改) / Sort in place
    slices.Sort(unsorted)
    fmt.Println("Sort:", unsorted)  // [1 1 2 3 4 5 6 9]
    
    // SortFunc - 自定义排序 / Custom sort
    people := []struct {
        Name string
        Age  int
    }{
        {"Alice", 30},
        {"Bob", 25},
        {"Charlie", 35},
    }
    slices.SortFunc(people, func(a, b struct{ Name string; Age int }) int {
        return cmp.Compare(a.Age, b.Age)
    })
    fmt.Println("SortFunc by age:", people)
    
    // SortStableFunc - 稳定排序 / Stable sort
    slices.SortStableFunc(people, func(a, b struct{ Name string; Age int }) int {
        return cmp.Compare(a.Age, b.Age)
    })
    
    // IsSorted - 检查是否已排序 / Check if sorted
    fmt.Println("IsSorted:", slices.IsSorted(unsorted))  // true
    
    // IsSortedFunc - 使用函数检查排序 / Check sort with function
    isSortedByAge := slices.IsSortedFunc(people, func(a, b struct{ Name string; Age int }) int {
        return cmp.Compare(a.Age, b.Age)
    })
    fmt.Println("IsSortedFunc by age:", isSortedByAge)
    
    // ========================================
    // 二分查找 (需要已排序) / Binary Search (requires sorted)
    // ========================================
    
    sorted := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    
    // BinarySearch - 二分查找 / Binary search
    idx, found := slices.BinarySearch(sorted, 5)
    fmt.Printf("BinarySearch 5: index=%d, found=%t\n", idx, found)  // 4, true
    
    idx, found = slices.BinarySearch(sorted, 11)
    fmt.Printf("BinarySearch 11: index=%d, found=%t\n", idx, found)  // 10, false (插入位置)
    
    // BinarySearchFunc - 自定义二分查找 / Custom binary search
    idx, found = slices.BinarySearchFunc(sorted, 5, func(elem, target int) int {
        return cmp.Compare(elem, target)
    })
    fmt.Printf("BinarySearchFunc 5: index=%d, found=%t\n", idx, found)
    
    // ========================================
    // 最值操作 / Min/Max Operations
    // ========================================
    
    values := []int{3, 1, 4, 1, 5, 9, 2, 6}
    
    // Min - 最小值 / Minimum
    fmt.Println("Min:", slices.Min(values))  // 1
    
    // Max - 最大值 / Maximum
    fmt.Println("Max:", slices.Max(values))  // 9
    
    // MinFunc / MaxFunc - 自定义比较 / Custom comparison
    strSlice := []string{"apple", "pie", "go"}
    shortest := slices.MinFunc(strSlice, func(a, b string) int {
        return cmp.Compare(len(a), len(b))
    })
    fmt.Println("MinFunc (shortest):", shortest)  // "go"
    
    longest := slices.MaxFunc(strSlice, func(a, b string) int {
        return cmp.Compare(len(a), len(b))
    })
    fmt.Println("MaxFunc (longest):", longest)  // "apple"
    
    // ========================================
    // 修改操作 / Modification Operations
    // ========================================
    
    // Insert - 插入元素 / Insert elements
    original := []int{1, 2, 5, 6}
    inserted := slices.Insert(original, 2, 3, 4)
    fmt.Println("Insert:", inserted)  // [1 2 3 4 5 6]
    
    // Delete - 删除元素 / Delete elements
    toDelete := []int{1, 2, 3, 4, 5, 6}
    deleted := slices.Delete(toDelete, 2, 4)  // 删除索引 2-3
    fmt.Println("Delete [2:4]:", deleted)  // [1 2 5 6]
    
    // DeleteFunc - 条件删除 / Delete with condition
    withEvens := []int{1, 2, 3, 4, 5, 6}
    withoutEvens := slices.DeleteFunc(withEvens, func(n int) bool {
        return n%2 == 0
    })
    fmt.Println("DeleteFunc (evens):", withoutEvens)  // [1 3 5]
    
    // Replace - 替换元素 / Replace elements
    toReplace := []int{1, 2, 3, 4, 5}
    replaced := slices.Replace(toReplace, 1, 3, 10, 20, 30)
    fmt.Println("Replace [1:3]:", replaced)  // [1 10 20 30 4 5]
    
    // Reverse - 反转切片 / Reverse slice
    toReverse := []int{1, 2, 3, 4, 5}
    slices.Reverse(toReverse)
    fmt.Println("Reverse:", toReverse)  // [5 4 3 2 1]
    
    // Compact - 去除相邻重复 / Remove adjacent duplicates
    withDups := []int{1, 1, 2, 2, 2, 3, 3, 4}
    compacted := slices.Compact(withDups)
    fmt.Println("Compact:", compacted)  // [1 2 3 4]
    
    // CompactFunc - 自定义去重 / Custom compact
    strWithDups := []string{"hello", "HELLO", "world", "WORLD"}
    compactedStr := slices.CompactFunc(strWithDups, func(a, b string) bool {
        return strings.EqualFold(a, b)
    })
    fmt.Println("CompactFunc:", compactedStr)  // ["hello" "world"]
    
    // Clip - 移除未使用容量 / Remove unused capacity
    large := make([]int, 3, 100)
    large[0], large[1], large[2] = 1, 2, 3
    clipped := slices.Clip(large)
    fmt.Printf("Clip: len=%d, cap=%d\n", len(clipped), cap(clipped))  // len=3, cap=3
    
    // Grow - 增加容量 / Grow capacity
    small := []int{1, 2, 3}
    grown := slices.Grow(small, 100)
    fmt.Printf("Grow: len=%d, cap>=%d\n", len(grown), 100)
}
```

## 2. slices 迭代器函数 / slices Iterator Functions (Go 1.23+)

```go
package main

import (
    "fmt"
    "slices"
)

func main() {
    s := []int{1, 2, 3, 4, 5}
    
    // ========================================
    // 迭代器生成函数 / Iterator Generator Functions
    // ========================================
    
    // All - 返回索引和值的迭代器 / Returns iterator of (index, value)
    fmt.Println("slices.All:")
    for i, v := range slices.All(s) {
        fmt.Printf("  [%d]: %d\n", i, v)
    }
    
    // Values - 只返回值的迭代器 / Returns iterator of values only
    fmt.Println("slices.Values:")
    for v := range slices.Values(s) {
        fmt.Printf("  %d\n", v)
    }
    
    // Backward - 反向迭代器 / Reverse iterator
    fmt.Println("slices.Backward:")
    for i, v := range slices.Backward(s) {
        fmt.Printf("  [%d]: %d\n", i, v)
    }
    // 输出: [4]:5, [3]:4, [2]:3, [1]:2, [0]:1
    
    // ========================================
    // 迭代器收集函数 / Iterator Collection Functions
    // ========================================
    
    // Collect - 将迭代器收集到切片 / Collect iterator to slice
    doubled := func(yield func(int) bool) {
        for _, v := range s {
            if !yield(v * 2) {
                return
            }
        }
    }
    collected := slices.Collect(doubled)
    fmt.Println("Collect:", collected)  // [2 4 6 8 10]
    
    // AppendSeq - 将迭代器追加到切片 / Append iterator to slice
    existing := []int{100, 200}
    tripled := func(yield func(int) bool) {
        for _, v := range []int{1, 2, 3} {
            if !yield(v * 3) {
                return
            }
        }
    }
    appended := slices.AppendSeq(existing, tripled)
    fmt.Println("AppendSeq:", appended)  // [100 200 3 6 9]
    
    // Sorted - 从迭代器创建排序后的切片 / Create sorted slice from iterator
    unsortedIter := func(yield func(int) bool) {
        for _, v := range []int{3, 1, 4, 1, 5} {
            if !yield(v) {
                return
            }
        }
    }
    sortedSlice := slices.Sorted(unsortedIter)
    fmt.Println("Sorted:", sortedSlice)  // [1 1 3 4 5]
    
    // SortedFunc - 自定义排序收集 / Custom sorted collection
    sortedDesc := slices.SortedFunc(unsortedIter, func(a, b int) int {
        return b - a  // 降序
    })
    fmt.Println("SortedFunc (desc):", sortedDesc)  // [5 4 3 1 1]
    
    // SortedStableFunc - 稳定排序收集 / Stable sorted collection
    stableSorted := slices.SortedStableFunc(unsortedIter, func(a, b int) int {
        return a - b
    })
    fmt.Println("SortedStableFunc:", stableSorted)
    
    // ========================================
    // Chunk - 分块迭代器 / Chunk iterator
    // ========================================
    
    data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    
    fmt.Println("slices.Chunk (size=3):")
    for chunk := range slices.Chunk(data, 3) {
        fmt.Printf("  %v\n", chunk)
    }
    // 输出:
    //   [1 2 3]
    //   [4 5 6]
    //   [7 8 9]
    //   [10]
    
    // 配合 Collect 使用 / Use with Collect
    chunks := slices.Collect(slices.Chunk(data, 4))
    fmt.Println("Collected chunks:", chunks)
    // [[1 2 3 4] [5 6 7 8] [9 10]]
}
```

## 3. iter 包详解 / iter Package Details (Go 1.23+)

```go
package main

import (
    "fmt"
    "iter"
    "slices"
)

// ========================================
// 迭代器类型 / Iterator Types
// ========================================

// iter.Seq[V] - 单值迭代器
// type Seq[V any] func(yield func(V) bool)

// iter.Seq2[K, V] - 双值迭代器 (如 index, value)
// type Seq2[K, V any] func(yield func(K, V) bool)

// ========================================
// 创建自定义迭代器 / Create Custom Iterators
// ========================================

// 范围迭代器 / Range iterator
func Range(start, end int) iter.Seq[int] {
    return func(yield func(int) bool) {
        for i := start; i < end; i++ {
            if !yield(i) {
                return
            }
        }
    }
}

// 步进范围迭代器 / Range with step
func RangeStep(start, end, step int) iter.Seq[int] {
    return func(yield func(int) bool) {
        for i := start; i < end; i += step {
            if !yield(i) {
                return
            }
        }
    }
}

// 无限迭代器 / Infinite iterator
func Naturals() iter.Seq[int] {
    return func(yield func(int) bool) {
        for i := 0; ; i++ {
            if !yield(i) {
                return
            }
        }
    }
}

// 重复迭代器 / Repeat iterator
func Repeat[T any](value T, n int) iter.Seq[T] {
    return func(yield func(T) bool) {
        for i := 0; i < n; i++ {
            if !yield(value) {
                return
            }
        }
    }
}

// Enumerate - 添加索引 / Add index
func Enumerate[T any](seq iter.Seq[T]) iter.Seq2[int, T] {
    return func(yield func(int, T) bool) {
        i := 0
        for v := range seq {
            if !yield(i, v) {
                return
            }
            i++
        }
    }
}

// ========================================
// 迭代器转换函数 / Iterator Transformation Functions
// ========================================

// Filter - 过滤 / Filter
func Filter[T any](seq iter.Seq[T], predicate func(T) bool) iter.Seq[T] {
    return func(yield func(T) bool) {
        for v := range seq {
            if predicate(v) {
                if !yield(v) {
                    return
                }
            }
        }
    }
}

// Map - 映射 / Map
func Map[T, U any](seq iter.Seq[T], transform func(T) U) iter.Seq[U] {
    return func(yield func(U) bool) {
        for v := range seq {
            if !yield(transform(v)) {
                return
            }
        }
    }
}

// Take - 取前 n 个 / Take first n
func Take[T any](seq iter.Seq[T], n int) iter.Seq[T] {
    return func(yield func(T) bool) {
        count := 0
        for v := range seq {
            if count >= n {
                return
            }
            if !yield(v) {
                return
            }
            count++
        }
    }
}

// Skip - 跳过前 n 个 / Skip first n
func Skip[T any](seq iter.Seq[T], n int) iter.Seq[T] {
    return func(yield func(T) bool) {
        count := 0
        for v := range seq {
            if count < n {
                count++
                continue
            }
            if !yield(v) {
                return
            }
        }
    }
}

// TakeWhile - 取直到条件不满足 / Take while condition is true
func TakeWhile[T any](seq iter.Seq[T], predicate func(T) bool) iter.Seq[T] {
    return func(yield func(T) bool) {
        for v := range seq {
            if !predicate(v) {
                return
            }
            if !yield(v) {
                return
            }
        }
    }
}

// DropWhile - 跳过直到条件不满足 / Drop while condition is true
func DropWhile[T any](seq iter.Seq[T], predicate func(T) bool) iter.Seq[T] {
    return func(yield func(T) bool) {
        dropping := true
        for v := range seq {
            if dropping && predicate(v) {
                continue
            }
            dropping = false
            if !yield(v) {
                return
            }
        }
    }
}

// Zip - 合并两个迭代器 / Zip two iterators
func Zip[T, U any](seq1 iter.Seq[T], seq2 iter.Seq[U]) iter.Seq2[T, U] {
    return func(yield func(T, U) bool) {
        next1, stop1 := iter.Pull(seq1)
        next2, stop2 := iter.Pull(seq2)
        defer stop1()
        defer stop2()
        
        for {
            v1, ok1 := next1()
            v2, ok2 := next2()
            if !ok1 || !ok2 {
                return
            }
            if !yield(v1, v2) {
                return
            }
        }
    }
}

// Chain - 连接多个迭代器 / Chain multiple iterators
func Chain[T any](seqs ...iter.Seq[T]) iter.Seq[T] {
    return func(yield func(T) bool) {
        for _, seq := range seqs {
            for v := range seq {
                if !yield(v) {
                    return
                }
            }
        }
    }
}

// Flatten - 展平嵌套迭代器 / Flatten nested iterator
func Flatten[T any](seq iter.Seq[iter.Seq[T]]) iter.Seq[T] {
    return func(yield func(T) bool) {
        for innerSeq := range seq {
            for v := range innerSeq {
                if !yield(v) {
                    return
                }
            }
        }
    }
}

// ========================================
// 聚合函数 / Aggregation Functions
// ========================================

// Reduce - 归约 / Reduce
func Reduce[T, U any](seq iter.Seq[T], initial U, fn func(U, T) U) U {
    result := initial
    for v := range seq {
        result = fn(result, v)
    }
    return result
}

// Count - 计数 / Count
func Count[T any](seq iter.Seq[T]) int {
    count := 0
    for range seq {
        count++
    }
    return count
}

// Any - 任一满足 / Any matches
func Any[T any](seq iter.Seq[T], predicate func(T) bool) bool {
    for v := range seq {
        if predicate(v) {
            return true
        }
    }
    return false
}

// All - 全部满足 / All match
func All[T any](seq iter.Seq[T], predicate func(T) bool) bool {
    for v := range seq {
        if !predicate(v) {
            return false
        }
    }
    return true
}

// First - 获取第一个元素 / Get first element
func First[T any](seq iter.Seq[T]) (T, bool) {
    for v := range seq {
        return v, true
    }
    var zero T
    return zero, false
}

// Last - 获取最后一个元素 / Get last element
func Last[T any](seq iter.Seq[T]) (T, bool) {
    var last T
    found := false
    for v := range seq {
        last = v
        found = true
    }
    return last, found
}

// ========================================
// Pull 迭代器 / Pull Iterator
// ========================================

func demonstratePull() {
    fmt.Println("\n=== Pull Iterator ===")
    
    seq := Range(0, 5)
    
    // iter.Pull 将 push 迭代器转换为 pull 迭代器
    // iter.Pull converts push iterator to pull iterator
    next, stop := iter.Pull(seq)
    defer stop()  // 必须调用 stop 释放资源 / Must call stop to release resources
    
    // 手动迭代 / Manual iteration
    for {
        v, ok := next()
        if !ok {
            break
        }
        fmt.Printf("  Pull: %d\n", v)
    }
    
    // Pull2 用于双值迭代器 / Pull2 for two-value iterator
    seq2 := slices.All([]string{"a", "b", "c"})
    next2, stop2 := iter.Pull2(seq2)
    defer stop2()
    
    for {
        i, v, ok := next2()
        if !ok {
            break
        }
        fmt.Printf("  Pull2: [%d]=%s\n", i, v)
    }
}

func main() {
    // ========================================
    // 使用示例 / Usage Examples
    // ========================================
    
    fmt.Println("=== Basic Range ===")
    for n := range Range(0, 5) {
        fmt.Printf("  %d\n", n)
    }
    
    fmt.Println("\n=== Range with Step ===")
    for n := range RangeStep(0, 10, 2) {
        fmt.Printf("  %d\n", n)
    }
    
    fmt.Println("\n=== Filter and Map ===")
    // 从 1-10 中过滤偶数，然后平方
    result := Map(
        Filter(Range(1, 11), func(n int) bool { return n%2 == 0 }),
        func(n int) int { return n * n },
    )
    for n := range result {
        fmt.Printf("  %d\n", n)  // 4, 16, 36, 64, 100
    }
    
    fmt.Println("\n=== Take from Infinite ===")
    // 从无限序列中取前 5 个
    first5 := Take(Naturals(), 5)
    collected := slices.Collect(first5)
    fmt.Println("  First 5 naturals:", collected)
    
    fmt.Println("\n=== Chain Iterators ===")
    chained := Chain(Range(0, 3), Range(10, 13), Range(100, 103))
    for n := range chained {
        fmt.Printf("  %d\n", n)
    }
    
    fmt.Println("\n=== Enumerate ===")
    words := slices.Values([]string{"hello", "world", "go"})
    for i, w := range Enumerate(words) {
        fmt.Printf("  [%d]: %s\n", i, w)
    }
    
    fmt.Println("\n=== Zip ===")
    names := slices.Values([]string{"Alice", "Bob", "Charlie"})
    ages := slices.Values([]int{25, 30, 35})
    for name, age := range Zip(names, ages) {
        fmt.Printf("  %s: %d\n", name, age)
    }
    
    fmt.Println("\n=== Aggregations ===")
    nums := Range(1, 11)
    sum := Reduce(nums, 0, func(acc, n int) int { return acc + n })
    fmt.Printf("  Sum 1-10: %d\n", sum)
    
    nums2 := Range(1, 11)
    hasEven := Any(nums2, func(n int) bool { return n%2 == 0 })
    fmt.Printf("  Any even: %t\n", hasEven)
    
    nums3 := Range(2, 11)
    allPositive := All(nums3, func(n int) bool { return n > 0 })
    fmt.Printf("  All positive: %t\n", allPositive)
    
    // Pull 迭代器演示
    demonstratePull()
}
```

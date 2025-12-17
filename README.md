# Go Slice 操作方法详解 / Go Slice Operations Explained

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

一个全面的 Go Slice 操作教程，包含详细解释和实战示例。

**快速链接:**
- 📚 [快速入门](QUICKSTART.md) - 5 分钟上手
- 📋 [速查表](CHEATSHEET.md) - 快速查找常用操作
- 💻 [示例代码](examples/) - 8 个完整示例
- 🧪 [测试代码](sliceops_test.go) - 单元测试和基准测试

## 目录
- [什么是 Slice](#什么是-slice)
- [Slice 的内部结构](#slice-的内部结构)
- [声明和初始化](#声明和初始化)
- [基本操作](#基本操作)
- [高级操作](#高级操作)
- [性能优化](#性能优化)
- [实战案例](#实战案例)
- [常见陷阱](#常见陷阱)

## 什么是 Slice

Slice（切片）是 Go 语言中最重要的数据结构之一，它是对数组的抽象和封装。与数组不同，slice 是动态的、灵活的，可以自动扩容。

**特点：**
- 长度可变
- 引用类型
- 基于数组实现
- 零值为 nil

## Slice 的内部结构

Slice 在底层是一个结构体，包含三个字段：

```go
type slice struct {
    array unsafe.Pointer  // 指向底层数组的指针
    len   int             // 当前长度
    cap   int             // 容量
}
```

## 声明和初始化

### 1. 声明但不初始化（nil slice）

```go
var s []int
// s == nil
// len(s) == 0
// cap(s) == 0
```

### 2. 使用 make 函数创建

```go
// make([]T, length, capacity)
s1 := make([]int, 5)      // len=5, cap=5, [0,0,0,0,0]
s2 := make([]int, 5, 10)  // len=5, cap=10, [0,0,0,0,0]
```

### 3. 使用字面量初始化

```go
s := []int{1, 2, 3, 4, 5}  // len=5, cap=5
```

### 4. 从数组或 slice 切片

```go
arr := [5]int{1, 2, 3, 4, 5}
s1 := arr[1:4]   // [2, 3, 4], len=3, cap=4
s2 := arr[:3]    // [1, 2, 3], len=3, cap=5
s3 := arr[2:]    // [3, 4, 5], len=3, cap=3
s4 := arr[:]     // [1, 2, 3, 4, 5], len=5, cap=5
```

## 基本操作

### 1. 访问元素

```go
s := []int{1, 2, 3, 4, 5}
first := s[0]        // 1
last := s[len(s)-1]  // 5
```

### 2. 修改元素

```go
s := []int{1, 2, 3}
s[0] = 10  // [10, 2, 3]
```

### 3. 追加元素（append）

```go
s := []int{1, 2, 3}
s = append(s, 4)           // [1, 2, 3, 4]
s = append(s, 5, 6, 7)     // [1, 2, 3, 4, 5, 6, 7]

// 追加另一个 slice
s2 := []int{8, 9}
s = append(s, s2...)       // [1, 2, 3, 4, 5, 6, 7, 8, 9]
```

### 4. 复制（copy）

```go
src := []int{1, 2, 3}
dst := make([]int, len(src))
n := copy(dst, src)  // n = 3, dst = [1, 2, 3]

// 部分复制
dst2 := make([]int, 2)
copy(dst2, src)  // dst2 = [1, 2]
```

### 5. 切片操作

```go
s := []int{1, 2, 3, 4, 5}

// s[low:high] - 不包含 high
sub1 := s[1:3]   // [2, 3]
sub2 := s[:3]    // [1, 2, 3]
sub3 := s[2:]    // [3, 4, 5]
sub4 := s[:]     // [1, 2, 3, 4, 5]

// s[low:high:max] - 完整切片表达式
sub5 := s[1:3:4] // [2, 3], len=2, cap=3
```

## 高级操作

### 1. 删除元素

```go
// 删除索引 i 处的元素
func remove(s []int, i int) []int {
    return append(s[:i], s[i+1:]...)
}

// 删除索引 i 处的元素（不保持顺序，更快）
func removeUnordered(s []int, i int) []int {
    s[i] = s[len(s)-1]
    return s[:len(s)-1]
}
```

### 2. 插入元素

```go
// 在索引 i 处插入元素 v
func insert(s []int, i int, v int) []int {
    s = append(s, 0)       // 扩展 slice
    copy(s[i+1:], s[i:])   // 移动元素
    s[i] = v               // 设置新值
    return s
}

// 在索引 i 处插入多个元素
func insertMultiple(s []int, i int, values ...int) []int {
    return append(s[:i], append(values, s[i:]...)...)
}
```

### 3. 过滤元素

```go
func filter(s []int, fn func(int) bool) []int {
    result := s[:0]  // 复用底层数组
    for _, v := range s {
        if fn(v) {
            result = append(result, v)
        }
    }
    return result
}

// 使用示例
s := []int{1, 2, 3, 4, 5, 6}
evens := filter(s, func(n int) bool { return n%2 == 0 })
// evens = [2, 4, 6]
```

### 4. 反转 slice

```go
func reverse(s []int) {
    for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
        s[i], s[j] = s[j], s[i]
    }
}
```

### 5. 去重

```go
func unique(s []int) []int {
    seen := make(map[int]bool)
    result := make([]int, 0, len(s))
    for _, v := range s {
        if !seen[v] {
            seen[v] = true
            result = append(result, v)
        }
    }
    return result
}
```

## 性能优化

### 1. 预分配容量

```go
// 不好：多次扩容
var s []int
for i := 0; i < 1000; i++ {
    s = append(s, i)
}

// 好：预分配容量
s := make([]int, 0, 1000)
for i := 0; i < 1000; i++ {
    s = append(s, i)
}
```

### 2. 使用完整切片表达式避免内存泄漏

```go
// 可能导致内存泄漏
func getFirstTwo(s []int) []int {
    return s[:2]  // 仍然引用整个底层数组
}

// 好：限制容量
func getFirstTwo(s []int) []int {
    return s[:2:2]  // 容量也是 2
}

// 或者复制到新 slice
func getFirstTwo(s []int) []int {
    result := make([]int, 2)
    copy(result, s[:2])
    return result
}
```

### 3. 复用底层数组

```go
// 在原地过滤，避免分配新内存
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
```

## 实战案例

查看 `examples/` 目录下的实战代码示例：

1. **基础操作** (`examples/01_basic_operations.go`) - 创建、访问、修改 slice
2. **追加和复制** (`examples/02_append_copy.go`) - append 和 copy 的各种用法
3. **切片操作** (`examples/03_slicing.go`) - 切片表达式和子切片
4. **增删改查** (`examples/04_crud_operations.go`) - 插入、删除、更新、查找
5. **高级技巧** (`examples/05_advanced_techniques.go`) - 去重、反转、过滤等
6. **性能对比** (`examples/06_performance.go`) - 不同操作的性能测试
7. **数据处理** (`examples/07_data_processing.go`) - 实际数据处理场景
8. **内存管理** (`examples/08_memory_management.go`) - 内存优化技巧

## 常见陷阱

### 1. Slice 是引用类型

```go
s1 := []int{1, 2, 3}
s2 := s1        // s2 和 s1 共享底层数组
s2[0] = 100
// s1 = [100, 2, 3]，s1 也被修改了！
```

### 2. Append 可能导致重新分配

```go
s1 := []int{1, 2, 3}
s2 := s1
s1 = append(s1, 4)  // 可能分配新数组
s1[0] = 100
// s2 可能不受影响，取决于是否重新分配
```

### 3. 循环中的 append

```go
// 错误：可能导致无限循环或意外行为
for _, v := range s {
    s = append(s, v)  // 不要在循环中修改正在遍历的 slice
}

// 正确：先复制
original := make([]int, len(s))
copy(original, s)
for _, v := range original {
    s = append(s, v)
}
```

### 4. 切片和底层数组

```go
arr := [5]int{1, 2, 3, 4, 5}
s := arr[1:3]  // [2, 3]
s[0] = 100
// arr = [1, 100, 3, 4, 5]，原数组被修改了！
```

### 5. 空 slice vs nil slice

```go
var s1 []int        // nil slice
s2 := []int{}       // 空 slice
s3 := make([]int, 0) // 空 slice

s1 == nil  // true
s2 == nil  // false
s3 == nil  // false

// 但是它们的 len 和 cap 都是 0
// 大多数情况下可以互换使用
```

## 总结

**最佳实践：**

1. ✅ 尽可能预分配容量
2. ✅ 使用 `copy` 创建独立副本
3. ✅ 注意 slice 的引用特性
4. ✅ 使用完整切片表达式控制容量
5. ✅ 避免在循环中修改正在遍历的 slice
6. ❌ 不要假设 append 后的 slice 和原来共享底层数组
7. ❌ 不要保留对大 slice 的小切片引用（可能导致内存泄漏）

**性能提示：**

- `append` 的平摊时间复杂度是 O(1)
- `copy` 的时间复杂度是 O(n)
- 预分配可以显著提高性能
- 使用 `s[:0]` 复用底层数组进行原地操作

## 运行示例

### 方式 1: 交互式菜单

```bash
cd examples
go run main.go
# 然后选择要运行的示例 (1-8)
```

### 方式 2: 直接运行

```bash
# 运行所有示例
go run examples/01_basic_operations.go
go run examples/02_append_copy.go
go run examples/03_slicing.go
# ... 等等

# 运行测试
go test -v

# 运行性能测试
go test -bench=. -benchmem
```

### 方式 3: 使用库函数

```go
package main

import (
    "fmt"
    "github.com/go-slice-operations"
)

func main() {
    s := []int{1, 2, 3, 4, 5}
    
    // 过滤偶数
    evens := sliceops.Filter(s, func(n int) bool { 
        return n%2 == 0 
    })
    fmt.Println(evens)  // [2, 4]
    
    // 去重
    dup := []int{1, 2, 2, 3, 3, 3}
    unique := sliceops.UniqueOrdered(dup)
    fmt.Println(unique)  // [1, 2, 3]
}
```

## 项目结构

```
go-slice-operations/
├── README.md              # 主文档（本文件）
├── QUICKSTART.md          # 快速入门指南
├── CHEATSHEET.md          # 操作速查表
├── go.mod                 # Go 模块文件
├── sliceops.go            # Slice 操作函数库
├── sliceops_test.go       # 单元测试和基准测试
└── examples/              # 示例代码目录
    ├── main.go                      # 交互式运行器
    ├── 01_basic_operations.go       # 基础操作
    ├── 02_append_copy.go            # 追加和复制
    ├── 03_slicing.go                # 切片操作
    ├── 04_crud_operations.go        # 增删改查
    ├── 05_advanced_techniques.go    # 高级技巧
    ├── 06_performance.go            # 性能优化
    ├── 07_data_processing.go        # 数据处理
    └── 08_memory_management.go      # 内存管理
```

## 基准测试结果

```
BenchmarkAppendWithoutPrealloc-4         	  229255	      5716 ns/op	   25208 B/op	      12 allocs/op
BenchmarkAppendWithPrealloc-4            	 3604560	       345.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkAppendWithPreallocAndAssign-4   	 4179957	       265.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkFilterNewSlice-4                	  301584	      4263 ns/op	    8192 B/op	       1 allocs/op
BenchmarkFilterInPlace-4                 	  585618	      2021 ns/op	       0 B/op	       0 allocs/op
BenchmarkRemoveAtOrdered-4               	170704249	         7.191 ns/op	       0 B/op	       0 allocs/op
BenchmarkRemoveAtFast-4                  	1000000000	         0.2548 ns/op	       0 B/op	       0 allocs/op
```

**结论:**
- 预分配可提速 **16x**
- 原地操作可提速 **2x**
- 快速删除比保序删除快 **28x**

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT License

## 参考资源

- [Go Blog: Go Slices: usage and internals](https://go.dev/blog/slices-intro)
- [Effective Go: Slices](https://go.dev/doc/effective_go#slices)
- [Go by Example: Slices](https://gobyexample.com/slices)
- [Go Slice Tricks](https://github.com/golang/go/wiki/SliceTricks)

---

⭐ 如果这个项目对你有帮助，请给个星标支持！

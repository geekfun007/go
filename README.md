# Go Range Iterator 详解 & 实战

Go 1.23 引入了 **Range over Function**（范围遍历函数）特性，允许我们在 `for-range` 循环中使用自定义迭代器函数。这是 Go 语言迭代能力的重大增强。

## 📖 目录

- [基本概念](#基本概念)
- [迭代器函数签名](#迭代器函数签名)
- [实战示例](#实战示例)
- [标准库中的迭代器](#标准库中的迭代器)
- [最佳实践](#最佳实践)

---

## 基本概念

### 什么是 Range Iterator？

在 Go 1.23 之前，`for-range` 只能遍历以下类型：
- 数组 / 切片
- map
- string
- channel

Go 1.23 新增了对**函数类型**的 range 支持，称为 **Iterator（迭代器）**。

### 核心思想

迭代器是一个函数，它接收一个 **yield 函数** 作为参数：
- 迭代器负责产生值
- yield 函数负责消费值
- 当 yield 返回 `false` 时，迭代提前终止（对应 `break`）

```
┌─────────────────┐         ┌─────────────────┐
│    Iterator     │  yield  │   for-range     │
│   (生产者)       │ ──────> │   (消费者)       │
└─────────────────┘         └─────────────────┘
        │                           │
        │  yield(value)             │  处理 value
        │  <── true/false ──        │  继续/break
```

---

## 迭代器函数签名

Go 1.23 在 `iter` 包中定义了两种标准迭代器类型：

### 1. `iter.Seq[V]` - 单值迭代器

```go
type Seq[V any] func(yield func(V) bool)
```

用于遍历单个值的序列，类似遍历切片的值。

### 2. `iter.Seq2[K, V]` - 双值迭代器

```go
type Seq2[K, V any] func(yield func(K, V) bool)
```

用于遍历键值对，类似遍历 map。

### 3. 无值迭代器（不常用）

```go
func(yield func() bool)
```

---

## 实战示例

### 示例 1：基础迭代器

```go
// examples/01_basic/main.go
package main

import "fmt"

// 自定义迭代器：生成 0 到 n-1 的数字
func Range(n int) func(yield func(int) bool) {
    return func(yield func(int) bool) {
        for i := 0; i < n; i++ {
            if !yield(i) {
                return // 消费者调用了 break
            }
        }
    }
}

func main() {
    // 使用 for-range 遍历自定义迭代器
    for v := range Range(5) {
        fmt.Println(v)
    }
    // 输出: 0, 1, 2, 3, 4
}
```

### 示例 2：双值迭代器

```go
// examples/02_seq2/main.go
package main

import "fmt"

// 带索引的切片迭代器
func Enumerate[T any](slice []T) func(yield func(int, T) bool) {
    return func(yield func(int, T) bool) {
        for i, v := range slice {
            if !yield(i, v) {
                return
            }
        }
    }
}

func main() {
    fruits := []string{"apple", "banana", "cherry"}
    
    for i, fruit := range Enumerate(fruits) {
        fmt.Printf("%d: %s\n", i, fruit)
    }
}
```

### 示例 3：过滤迭代器

```go
// examples/03_filter/main.go
package main

import "fmt"

// 过滤迭代器
func Filter[T any](seq func(yield func(T) bool), predicate func(T) bool) func(yield func(T) bool) {
    return func(yield func(T) bool) {
        seq(func(v T) bool {
            if predicate(v) {
                return yield(v)
            }
            return true // 继续迭代
        })
    }
}

// 数字范围迭代器
func Range(start, end int) func(yield func(int) bool) {
    return func(yield func(int) bool) {
        for i := start; i < end; i++ {
            if !yield(i) {
                return
            }
        }
    }
}

func main() {
    // 过滤出偶数
    evenNumbers := Filter(Range(1, 10), func(n int) bool {
        return n%2 == 0
    })
    
    for n := range evenNumbers {
        fmt.Println(n) // 输出: 2, 4, 6, 8
    }
}
```

### 示例 4：链表遍历

```go
// examples/04_linkedlist/main.go
package main

import "fmt"

type Node[T any] struct {
    Value T
    Next  *Node[T]
}

type LinkedList[T any] struct {
    Head *Node[T]
}

// 为链表实现迭代器
func (l *LinkedList[T]) All() func(yield func(T) bool) {
    return func(yield func(T) bool) {
        for node := l.Head; node != nil; node = node.Next {
            if !yield(node.Value) {
                return
            }
        }
    }
}

func main() {
    list := &LinkedList[int]{
        Head: &Node[int]{Value: 1, Next: &Node[int]{Value: 2, Next: &Node[int]{Value: 3}}},
    }
    
    for v := range list.All() {
        fmt.Println(v) // 输出: 1, 2, 3
    }
}
```

### 示例 5：二叉树遍历

```go
// examples/05_tree/main.go
package main

import "fmt"

type TreeNode struct {
    Value int
    Left  *TreeNode
    Right *TreeNode
}

// 中序遍历迭代器
func (t *TreeNode) InOrder() func(yield func(int) bool) {
    return func(yield func(int) bool) {
        var traverse func(*TreeNode) bool
        traverse = func(node *TreeNode) bool {
            if node == nil {
                return true
            }
            // 左 -> 根 -> 右
            if !traverse(node.Left) {
                return false
            }
            if !yield(node.Value) {
                return false
            }
            return traverse(node.Right)
        }
        traverse(t)
    }
}

func main() {
    //       4
    //      / \
    //     2   6
    //    / \ / \
    //   1  3 5  7
    tree := &TreeNode{
        Value: 4,
        Left: &TreeNode{
            Value: 2,
            Left:  &TreeNode{Value: 1},
            Right: &TreeNode{Value: 3},
        },
        Right: &TreeNode{
            Value: 6,
            Left:  &TreeNode{Value: 5},
            Right: &TreeNode{Value: 7},
        },
    }
    
    fmt.Print("中序遍历: ")
    for v := range tree.InOrder() {
        fmt.Printf("%d ", v) // 输出: 1 2 3 4 5 6 7
    }
    fmt.Println()
}
```

### 示例 6：组合迭代器（Map + Filter + Take）

```go
// examples/06_compose/main.go
package main

import "fmt"

// Map 转换迭代器
func Map[T, U any](seq func(yield func(T) bool), transform func(T) U) func(yield func(U) bool) {
    return func(yield func(U) bool) {
        seq(func(v T) bool {
            return yield(transform(v))
        })
    }
}

// Take 限制数量
func Take[T any](seq func(yield func(T) bool), n int) func(yield func(T) bool) {
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

// Filter 过滤
func Filter[T any](seq func(yield func(T) bool), predicate func(T) bool) func(yield func(T) bool) {
    return func(yield func(T) bool) {
        seq(func(v T) bool {
            if predicate(v) {
                return yield(v)
            }
            return true
        })
    }
}

// 无限序列生成器
func Naturals() func(yield func(int) bool) {
    return func(yield func(int) bool) {
        for i := 1; ; i++ {
            if !yield(i) {
                return
            }
        }
    }
}

func main() {
    // 组合：取前10个自然数 -> 过滤偶数 -> 平方
    result := Map(
        Filter(
            Take(Naturals(), 10),
            func(n int) bool { return n%2 == 0 },
        ),
        func(n int) int { return n * n },
    )
    
    fmt.Print("前10个自然数中偶数的平方: ")
    for v := range result {
        fmt.Printf("%d ", v) // 输出: 4 16 36 64 100
    }
    fmt.Println()
}
```

### 示例 7：并发安全迭代器

```go
// examples/07_concurrent/main.go
package main

import (
    "fmt"
    "sync"
)

type SafeMap[K comparable, V any] struct {
    mu sync.RWMutex
    m  map[K]V
}

func NewSafeMap[K comparable, V any]() *SafeMap[K, V] {
    return &SafeMap[K, V]{m: make(map[K]V)}
}

func (s *SafeMap[K, V]) Set(key K, value V) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.m[key] = value
}

// 安全迭代（快照方式）
func (s *SafeMap[K, V]) All() func(yield func(K, V) bool) {
    return func(yield func(K, V) bool) {
        // 获取快照
        s.mu.RLock()
        snapshot := make([]struct{ k K; v V }, 0, len(s.m))
        for k, v := range s.m {
            snapshot = append(snapshot, struct{ k K; v V }{k, v})
        }
        s.mu.RUnlock()
        
        // 遍历快照
        for _, item := range snapshot {
            if !yield(item.k, item.v) {
                return
            }
        }
    }
}

func main() {
    m := NewSafeMap[string, int]()
    m.Set("a", 1)
    m.Set("b", 2)
    m.Set("c", 3)
    
    for k, v := range m.All() {
        fmt.Printf("%s: %d\n", k, v)
    }
}
```

### 示例 8：文件行迭代器

```go
// examples/08_file/main.go
package main

import (
    "bufio"
    "fmt"
    "os"
)

// 文件行迭代器
func Lines(filename string) func(yield func(int, string) bool) {
    return func(yield func(int, string) bool) {
        file, err := os.Open(filename)
        if err != nil {
            return
        }
        defer file.Close()
        
        scanner := bufio.NewScanner(file)
        lineNum := 1
        for scanner.Scan() {
            if !yield(lineNum, scanner.Text()) {
                return
            }
            lineNum++
        }
    }
}

func main() {
    // 遍历文件每一行
    for lineNum, line := range Lines("example.txt") {
        fmt.Printf("%4d: %s\n", lineNum, line)
    }
}
```

---

## 标准库中的迭代器

Go 1.23+ 标准库已经添加了许多迭代器支持：

### `slices` 包

```go
import "slices"

s := []int{1, 2, 3, 4, 5}

// slices.All - 返回索引和值的迭代器
for i, v := range slices.All(s) {
    fmt.Println(i, v)
}

// slices.Values - 返回值的迭代器
for v := range slices.Values(s) {
    fmt.Println(v)
}

// slices.Backward - 反向迭代
for i, v := range slices.Backward(s) {
    fmt.Println(i, v) // 4:5, 3:4, 2:3, 1:2, 0:1
}

// slices.Collect - 从迭代器收集为切片
result := slices.Collect(someIterator)
```

### `maps` 包

```go
import "maps"

m := map[string]int{"a": 1, "b": 2}

// maps.All - 返回键值对迭代器
for k, v := range maps.All(m) {
    fmt.Println(k, v)
}

// maps.Keys - 返回键的迭代器
for k := range maps.Keys(m) {
    fmt.Println(k)
}

// maps.Values - 返回值的迭代器
for v := range maps.Values(m) {
    fmt.Println(v)
}
```

### `iter` 包

```go
import "iter"

// iter.Pull - 将推模式迭代器转为拉模式
next, stop := iter.Pull(someSeq)
defer stop()

for {
    v, ok := next()
    if !ok {
        break
    }
    fmt.Println(v)
}

// iter.Pull2 - 双值版本
next2, stop2 := iter.Pull2(someSeq2)
defer stop2()

for {
    k, v, ok := next2()
    if !ok {
        break
    }
    fmt.Println(k, v)
}
```

---

## 最佳实践

### ✅ 推荐做法

1. **使用标准签名**
```go
// 好：使用 iter.Seq 和 iter.Seq2
func MyIterator() iter.Seq[int]
func MyIterator2() iter.Seq2[string, int]
```

2. **正确处理 yield 返回值**
```go
// 好：检查 yield 返回值
for i := 0; i < n; i++ {
    if !yield(i) {
        return // 立即退出
    }
}
```

3. **迭代器应该是无状态的**
```go
// 好：每次调用返回新的迭代器
func Range(n int) iter.Seq[int] {
    return func(yield func(int) bool) {
        for i := 0; i < n; i++ {
            if !yield(i) { return }
        }
    }
}
```

### ❌ 避免的做法

1. **忽略 yield 返回值**
```go
// 坏：忽略返回值会导致 break 失效
for i := 0; i < n; i++ {
    yield(i) // 错误！
}
```

2. **在迭代器中 panic**
```go
// 坏：应该优雅处理错误
func BadIterator() iter.Seq[int] {
    return func(yield func(int) bool) {
        panic("error") // 不推荐
    }
}
```

3. **重复使用有状态迭代器**
```go
// 坏：有状态迭代器只能用一次
counter := 0
badIter := func(yield func(int) bool) {
    for counter < 10 {
        if !yield(counter) { return }
        counter++
    }
}
```

---

## 与传统方式对比

| 特性 | Channel | Callback | Range Iterator |
|------|---------|----------|----------------|
| 语法 | `for v := range ch` | `fn(func(v))` | `for v := range iter` |
| 性能 | 有 goroutine 开销 | 最优 | 接近最优 |
| break 支持 | ✅ (需关闭) | ❌ 困难 | ✅ 原生支持 |
| 内存 | 需要 buffer | 无 | 无 |
| 适用场景 | 并发 | 简单回调 | 通用迭代 |

---

## 总结

Go 1.23 的 Range Iterator 特性：

1. **语法自然** - 与现有 for-range 无缝结合
2. **类型安全** - 泛型支持
3. **惰性求值** - 按需生成值
4. **组合性强** - 易于组合多个迭代器
5. **性能优秀** - 无 goroutine 开销

适用场景：
- 自定义数据结构遍历
- 惰性序列生成
- 数据流处理管道
- 文件/网络流处理

## 运行示例

**Go 1.23+** (推荐):
```bash
go run examples/01_basic/main.go
```

**Go 1.22** (需要启用实验性特性):
```bash
GOEXPERIMENT=rangefunc go run examples/01_basic/main.go
```

### 所有示例

```bash
# 01 - 基础迭代器
GOEXPERIMENT=rangefunc go run examples/01_basic/main.go

# 02 - 双值迭代器 (Seq2)
GOEXPERIMENT=rangefunc go run examples/02_seq2/main.go

# 03 - 过滤/映射迭代器
GOEXPERIMENT=rangefunc go run examples/03_filter/main.go

# 04 - 链表遍历
GOEXPERIMENT=rangefunc go run examples/04_linkedlist/main.go

# 05 - 二叉树遍历
GOEXPERIMENT=rangefunc go run examples/05_tree/main.go

# 06 - 组合迭代器
GOEXPERIMENT=rangefunc go run examples/06_compose/main.go

# 07 - 并发安全迭代器
GOEXPERIMENT=rangefunc go run examples/07_concurrent/main.go

# 08 - 文件迭代器
GOEXPERIMENT=rangefunc go run examples/08_file/main.go

# 09 - Pull 迭代器 (推模式转拉模式)
GOEXPERIMENT=rangefunc go run examples/09_pull/main.go
```

## 参考资料

- [Go 1.23 Release Notes](https://go.dev/doc/go1.23)
- [Range Over Function Types](https://go.dev/wiki/RangefuncExperiment)
- [iter package](https://pkg.go.dev/iter)

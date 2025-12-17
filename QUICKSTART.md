# Go Slice 快速入门

## 快速开始

### 1. 克隆或下载项目

```bash
git clone https://github.com/username/go-slice-operations
cd go-slice-operations
```

### 2. 运行示例

#### 方式 1: 交互式菜单

```bash
cd examples
go run main.go
```

然后按提示选择要运行的示例（1-8）。

#### 方式 2: 直接运行单个示例

```bash
cd examples

# 基础操作
go run 01_basic_operations.go

# 追加和复制
go run 02_append_copy.go

# 切片操作
go run 03_slicing.go

# 增删改查
go run 04_crud_operations.go

# 高级技巧
go run 05_advanced_techniques.go

# 性能优化
go run 06_performance.go

# 数据处理
go run 07_data_processing.go

# 内存管理
go run 08_memory_management.go
```

### 3. 运行测试

```bash
# 运行所有单元测试
go test -v

# 运行基准测试
go test -bench=. -benchmem

# 运行特定测试
go test -run TestInsertAt -v

# 运行特定基准测试
go test -bench=BenchmarkAppend -benchmem
```

### 4. 在你的项目中使用

#### 方式 1: 复制函数

从 `sliceops.go` 中复制你需要的函数到你的项目中。

#### 方式 2: 导入包（如果发布为模块）

```go
import "github.com/username/go-slice-operations"

func main() {
    s := []int{1, 2, 3, 4, 5}
    
    // 过滤偶数
    evens := sliceops.Filter(s, func(n int) bool { 
        return n%2 == 0 
    })
    
    // 去重
    unique := sliceops.UniqueOrdered(s)
    
    // 反转
    sliceops.Reverse(s)
}
```

## 常用操作速查

### 创建 Slice

```go
// nil slice
var s []int

// 空 slice
s := []int{}
s := make([]int, 0)

// 指定长度和容量
s := make([]int, 5)       // len=5, cap=5
s := make([]int, 5, 10)   // len=5, cap=10

// 字面量
s := []int{1, 2, 3, 4, 5}
```

### 追加元素

```go
s := []int{1, 2, 3}
s = append(s, 4)           // [1, 2, 3, 4]
s = append(s, 5, 6, 7)     // [1, 2, 3, 4, 5, 6, 7]

s2 := []int{8, 9}
s = append(s, s2...)       // [1, 2, 3, 4, 5, 6, 7, 8, 9]
```

### 复制 Slice

```go
src := []int{1, 2, 3}
dst := make([]int, len(src))
copy(dst, src)  // dst = [1, 2, 3]
```

### 切片

```go
s := []int{0, 1, 2, 3, 4, 5}
s[1:4]    // [1, 2, 3]
s[:3]     // [0, 1, 2]
s[3:]     // [3, 4, 5]
s[:]      // [0, 1, 2, 3, 4, 5]
s[1:4:5]  // [1, 2, 3], cap=4
```

### 插入元素

```go
// 在索引 i 处插入元素 v
s = append(s[:i], append([]int{v}, s[i:]...)...)

// 使用库函数
s = sliceops.InsertAt(s, i, v)
```

### 删除元素

```go
// 删除索引 i 处的元素
s = append(s[:i], s[i+1:]...)

// 使用库函数
s = sliceops.RemoveAt(s, i)
```

### 过滤

```go
// 保留偶数
evens := sliceops.Filter(s, func(n int) bool { 
    return n%2 == 0 
})
```

### 映射

```go
// 平方
squared := sliceops.Map(s, func(n int) int { 
    return n * n 
})
```

### 归约

```go
// 求和
sum := sliceops.Reduce(s, 0, func(acc, n int) int { 
    return acc + n 
})
```

### 去重

```go
unique := sliceops.UniqueOrdered(s)
```

### 反转

```go
sliceops.Reverse(s)  // 原地反转
```

### 查找

```go
// 查找索引
idx := sliceops.IndexOf(s, target)

// 检查是否包含
exists := sliceops.Contains(s, target)

// 查找最小最大值
min, max := sliceops.FindMinMax(s)
```

## 性能提示

### ✅ 推荐做法

```go
// 1. 预分配容量
s := make([]int, 0, expectedSize)
for i := 0; i < expectedSize; i++ {
    s = append(s, i)
}

// 2. 使用 copy 创建独立副本
newSlice := make([]int, len(oldSlice))
copy(newSlice, oldSlice)

// 3. 原地操作节省内存
s = s[:0]  // 重置长度，保留容量
for _, v := range data {
    if condition(v) {
        s = append(s, v)
    }
}

// 4. 使用完整切片表达式限制容量
sub := s[:n:n]  // 防止 append 修改原 slice
```

### ❌ 避免的做法

```go
// 1. 不预分配容量（会多次扩容）
var s []int
for i := 0; i < 10000; i++ {
    s = append(s, i)  // 慢！
}

// 2. 不必要的复制
for range times {
    newSlice := make([]int, len(s))
    copy(newSlice, s)
}

// 3. 大 slice 的小切片（内存泄漏）
huge := make([]byte, 1000000)
small := huge[:10]  // 仍引用整个 1MB！

// 4. 在循环中修改正在遍历的 slice
for _, v := range s {
    s = append(s, v)  // 危险！
}
```

## 学习路径

1. **初学者**: 从 `01_basic_operations.go` 开始
2. **进阶**: 学习 `02-04` 的 append、copy 和 CRUD 操作
3. **高级**: 掌握 `05-06` 的高级技巧和性能优化
4. **实战**: 参考 `07-08` 的实际场景和内存管理

## 常见问题

### Q1: Slice 和数组有什么区别？

- **数组**: 固定长度，值类型
- **Slice**: 可变长度，引用类型（指向底层数组）

### Q2: 为什么 append 后有时原 slice 会被修改？

Append 在容量足够时会复用底层数组，否则会重新分配。使用完整切片表达式或 copy 可避免这个问题。

### Q3: nil slice 和空 slice 有什么区别？

```go
var nilSlice []int        // nil
emptySlice := []int{}     // 非 nil，但 len=0

nilSlice == nil    // true
emptySlice == nil  // false

// 但功能上几乎相同，都可以 append
```

### Q4: 如何避免内存泄漏？

- 从大 slice 提取小片段时使用 `copy`
- 使用完整切片表达式限制容量
- 及时释放不需要的大 slice 引用

### Q5: Slice 是值传递还是引用传递？

Slice 本身是值传递（复制 slice 结构），但因为包含指向底层数组的指针，所以修改元素会影响原 slice。

## 下一步

- 阅读完整的 [README.md](README.md)
- 查看 [示例代码](examples/)
- 运行 [测试和基准测试](sliceops_test.go)
- 参考 Go 官方文档: https://go.dev/blog/slices-intro

## 反馈和贡献

欢迎提交 Issue 和 Pull Request！

# Go Slice 操作速查表

## 创建和初始化

| 操作 | 代码 | 说明 |
|------|------|------|
| nil slice | `var s []int` | 零值为 nil |
| 空 slice | `s := []int{}` | 非 nil，长度为 0 |
| 使用 make | `s := make([]int, 5)` | len=5, cap=5 |
| 指定容量 | `s := make([]int, 5, 10)` | len=5, cap=10 |
| 字面量 | `s := []int{1, 2, 3}` | 直接初始化值 |
| 从数组 | `s := arr[1:4]` | 从数组创建切片 |

## 访问和修改

| 操作 | 代码 | 说明 |
|------|------|------|
| 访问元素 | `v := s[i]` | 获取索引 i 的元素 |
| 修改元素 | `s[i] = v` | 设置索引 i 的值 |
| 获取长度 | `n := len(s)` | 当前元素个数 |
| 获取容量 | `n := cap(s)` | 底层数组容量 |
| 遍历 | `for i, v := range s {}` | 同时获取索引和值 |
| 只遍历值 | `for _, v := range s {}` | 只获取值 |
| 只遍历索引 | `for i := range s {}` | 只获取索引 |

## 追加 (Append)

| 操作 | 代码 | 说明 |
|------|------|------|
| 追加单个 | `s = append(s, v)` | 添加一个元素 |
| 追加多个 | `s = append(s, v1, v2, v3)` | 添加多个元素 |
| 追加 slice | `s = append(s, s2...)` | 合并两个 slice |
| 在开头插入 | `s = append([]int{v}, s...)` | 在开头添加元素 |

## 复制 (Copy)

| 操作 | 代码 | 说明 |
|------|------|------|
| 完整复制 | `dst := make([]int, len(src))`<br>`copy(dst, src)` | 创建独立副本 |
| 部分复制 | `copy(dst, src[:n])` | 只复制前 n 个 |
| 原地移动 | `copy(s[2:], s[:3])` | 移动元素位置 |

## 切片 (Slicing)

| 操作 | 代码 | 说明 |
|------|------|------|
| 基本切片 | `s[low:high]` | 包含 low，不包含 high |
| 从开始 | `s[:high]` | 等价于 s[0:high] |
| 到结尾 | `s[low:]` | 等价于 s[low:len(s)] |
| 完整复制 | `s[:]` | 等价于 s[0:len(s)] |
| 限制容量 | `s[low:high:max]` | len=high-low, cap=max-low |

## 插入

| 操作 | 代码 | 说明 |
|------|------|------|
| 在索引 i 插入 | `s = append(s[:i], append([]int{v}, s[i:]...)...)` | 保持顺序 |
| 优化版本 | `s = append(s, 0)`<br>`copy(s[i+1:], s[i:])`<br>`s[i] = v` | 更高效 |

## 删除

| 操作 | 代码 | 说明 |
|------|------|------|
| 删除第一个 | `s = s[1:]` | O(1) |
| 删除最后一个 | `s = s[:len(s)-1]` | O(1) |
| 删除索引 i | `s = append(s[:i], s[i+1:]...)` | O(n)，保持顺序 |
| 快速删除 | `s[i] = s[len(s)-1]`<br>`s = s[:len(s)-1]` | O(1)，不保持顺序 |
| 删除范围 | `s = append(s[:i], s[j:]...)` | 删除 [i, j) |

## 查找

| 操作 | 代码 | 说明 |
|------|------|------|
| 线性查找 | `for i, v := range s {`<br>`  if v == target { return i }`<br>`}` | O(n) |
| 检查包含 | `contains := false`<br>`for _, v := range s {`<br>`  if v == target { contains = true }`<br>`}` | O(n) |
| 查找最大值 | `max := s[0]`<br>`for _, v := range s[1:] {`<br>`  if v > max { max = v }`<br>`}` | O(n) |

## 高级操作

### 反转

```go
for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
    s[i], s[j] = s[j], s[i]
}
```

### 去重（保持顺序）

```go
seen := make(map[int]bool)
result := make([]int, 0, len(s))
for _, v := range s {
    if !seen[v] {
        seen[v] = true
        result = append(result, v)
    }
}
```

### 过滤

```go
result := s[:0]  // 复用底层数组
for _, v := range s {
    if condition(v) {
        result = append(result, v)
    }
}
```

### 映射 (Map)

```go
result := make([]int, len(s))
for i, v := range s {
    result[i] = transform(v)
}
```

### 归约 (Reduce)

```go
acc := initialValue
for _, v := range s {
    acc = combine(acc, v)
}
```

### 分组

```go
groups := make(map[string][]int)
for _, v := range s {
    key := getKey(v)
    groups[key] = append(groups[key], v)
}
```

## 性能对比

| 操作 | 不预分配 | 预分配容量 | 预分配长度 |
|------|----------|------------|------------|
| Append | ~5700 ns | ~345 ns ⚡ | ~265 ns ⚡⚡ |
| 内存分配 | 12 次 | 0 次 | 0 次 |

| 操作 | 新建 slice | 原地操作 |
|------|------------|----------|
| Filter | ~4260 ns | ~2020 ns ⚡ |
| 内存分配 | 1 次 | 0 次 |

| 操作 | 保持顺序 | 不保持顺序 |
|------|----------|------------|
| Remove | ~7 ns | ~0.25 ns ⚡⚡ |

## 常见陷阱

### ⚠️ 陷阱 1: Slice 共享底层数组

```go
s1 := []int{1, 2, 3}
s2 := s1        // 共享底层数组
s2[0] = 100
// s1 = [100, 2, 3]  ❌ s1 也被修改了！
```

**解决方案**: 使用 copy 创建独立副本

```go
s2 := make([]int, len(s1))
copy(s2, s1)
```

### ⚠️ 陷阱 2: Append 可能重新分配

```go
s1 := []int{1, 2, 3}
s2 := s1
s1 = append(s1, 4)  // 可能分配新数组
s1[0] = 100
// s2 可能不受影响，取决于是否重新分配
```

**解决方案**: 使用完整切片表达式

```go
s2 := s1[:len(s1):len(s1)]  // 限制容量
```

### ⚠️ 陷阱 3: 大 Slice 的小切片

```go
huge := make([]byte, 1000000)  // 1MB
small := huge[:10]              // 仍引用整个 1MB！
```

**解决方案**: 复制到新 slice

```go
small := make([]byte, 10)
copy(small, huge[:10])
```

### ⚠️ 陷阱 4: 循环中修改 Slice

```go
for _, v := range s {
    s = append(s, v)  // ❌ 危险！可能无限循环
}
```

**解决方案**: 先复制

```go
original := make([]int, len(s))
copy(original, s)
for _, v := range original {
    s = append(s, v)
}
```

## 最佳实践

### ✅ 推荐

1. **预分配容量**
   ```go
   s := make([]int, 0, expectedSize)
   ```

2. **使用 copy 创建独立副本**
   ```go
   dst := make([]int, len(src))
   copy(dst, src)
   ```

3. **完整切片表达式限制容量**
   ```go
   sub := s[:n:n]
   ```

4. **原地操作节省内存**
   ```go
   s = s[:0]  // 重置长度，保留容量
   ```

5. **及时释放大对象**
   ```go
   largeSlice = nil
   ```

### ❌ 避免

1. ❌ 不预分配容量（频繁扩容）
2. ❌ 不必要的复制
3. ❌ 保留大 slice 的小切片引用
4. ❌ 在循环中修改正在遍历的 slice
5. ❌ 假设 append 后仍共享底层数组

## 扩容策略

```
容量 < 256:   新容量 = 旧容量 × 2
容量 >= 256:  新容量 ≈ 旧容量 × 1.25
```

**示例:**
- 1 → 2 → 4 → 8 → 16 → 32 → 64 → 128 → 256
- 256 → 320 → 400 → 500 → 625 → ...

## 时间复杂度

| 操作 | 时间复杂度 | 说明 |
|------|------------|------|
| 访问 | O(1) | 直接索引 |
| Append (容量足够) | O(1) | 无需扩容 |
| Append (需扩容) | O(n) | 需要复制 |
| Insert | O(n) | 需要移动元素 |
| Delete (保持顺序) | O(n) | 需要移动元素 |
| Delete (不保持顺序) | O(1) | 交换删除 |
| Search | O(n) | 线性查找 |
| Copy | O(n) | 复制所有元素 |

## 内存布局

```
slice 结构 (24 字节):
┌─────────────────┬─────────┬─────────┐
│ array (8 字节)  │ len (8) │ cap (8) │
└────────┬────────┴─────────┴─────────┘
         │
         ▼
    ┌─────┬─────┬─────┬─────┬─────┐
    │  1  │  2  │  3  │  4  │  5  │  底层数组
    └─────┴─────┴─────┴─────┴─────┘
```

## 快速参考

```go
// 创建
s := make([]int, 0, 100)    // 最佳：预分配
s := []int{1, 2, 3}         // 字面量

// 追加
s = append(s, 4, 5, 6)      // 追加多个
s = append(s, s2...)        // 合并

// 插入 (在索引 i)
s = append(s[:i], append([]int{v}, s[i:]...)...)

// 删除 (索引 i)
s = append(s[:i], s[i+1:]...)     // 保持顺序
s[i] = s[len(s)-1]; s = s[:len(s)-1]  // 快速

// 复制
dst := make([]int, len(src))
copy(dst, src)

// 查找
idx := -1
for i, v := range s {
    if v == target {
        idx = i
        break
    }
}

// 过滤
result := s[:0]
for _, v := range s {
    if v%2 == 0 {
        result = append(result, v)
    }
}

// 去重
seen := make(map[int]bool)
unique := make([]int, 0)
for _, v := range s {
    if !seen[v] {
        seen[v] = true
        unique = append(unique, v)
    }
}

// 反转
for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
    s[i], s[j] = s[j], s[i]
}
```

## 相关资源

- [Go Blog: Slices](https://go.dev/blog/slices-intro)
- [Effective Go](https://go.dev/doc/effective_go#slices)
- [Go by Example](https://gobyexample.com/slices)

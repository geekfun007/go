# 数组与切片 (Array & Slice) / Array & Slice

## 1. 数组 (Array)

```go
package main

import "fmt"

func main() {
    // 数组是固定长度的同类型元素序列
    // Arrays are fixed-length sequences of same-type elements
    
    // 声明数组 / Declare array
    var arr1 [5]int                      // 零值初始化
    arr2 := [5]int{1, 2, 3, 4, 5}        // 字面量初始化
    arr3 := [...]int{1, 2, 3}            // 自动计算长度
    arr4 := [5]int{0: 10, 4: 50}         // 指定索引初始化
    
    fmt.Println("arr1:", arr1)  // [0 0 0 0 0]
    fmt.Println("arr2:", arr2)  // [1 2 3 4 5]
    fmt.Println("arr3:", arr3)  // [1 2 3]
    fmt.Println("arr4:", arr4)  // [10 0 0 0 50]
    
    // 数组长度 / Array length
    fmt.Println("len(arr2):", len(arr2))
    
    // 访问和修改元素 / Access and modify elements
    fmt.Println("arr2[0]:", arr2[0])
    arr2[0] = 100
    fmt.Println("Modified arr2:", arr2)
    
    // 遍历数组 / Iterate array
    for i, v := range arr2 {
        fmt.Printf("index=%d, value=%d\n", i, v)
    }
    
    // 多维数组 / Multi-dimensional array
    var matrix [3][3]int
    matrix[0][0] = 1
    matrix[1][1] = 1
    matrix[2][2] = 1
    fmt.Println("Matrix:", matrix)
    
    // 数组是值类型 (赋值会复制) / Arrays are value types (assignment copies)
    arr5 := arr2
    arr5[0] = 999
    fmt.Println("arr2[0]:", arr2[0])  // 100 (未改变)
    fmt.Println("arr5[0]:", arr5[0])  // 999
    
    // 数组比较 / Array comparison
    a := [3]int{1, 2, 3}
    b := [3]int{1, 2, 3}
    fmt.Println("a == b:", a == b)  // true
}
```

## 2. 切片 (Slice)

```go
package main

import "fmt"

func main() {
    // 切片是动态大小的，对数组的引用
    // Slices are dynamic-sized references to arrays
    
    // 创建切片 / Create slice
    var s1 []int                        // nil 切片
    s2 := []int{1, 2, 3, 4, 5}          // 字面量创建
    s3 := make([]int, 5)                // make 创建，长度5
    s4 := make([]int, 5, 10)            // 长度5，容量10
    
    fmt.Println("s1:", s1, "len:", len(s1), "cap:", cap(s1))
    fmt.Println("s2:", s2, "len:", len(s2), "cap:", cap(s2))
    fmt.Println("s3:", s3, "len:", len(s3), "cap:", cap(s3))
    fmt.Println("s4:", s4, "len:", len(s4), "cap:", cap(s4))
    
    // 从数组创建切片 / Create slice from array
    arr := [5]int{1, 2, 3, 4, 5}
    slice := arr[1:4]  // [2, 3, 4], 索引 1 到 3
    fmt.Println("slice from array:", slice)
    
    // 切片操作 / Slice operations
    fmt.Println("s2[:]:", s2[:])      // 全部
    fmt.Println("s2[:3]:", s2[:3])    // 前3个
    fmt.Println("s2[2:]:", s2[2:])    // 从索引2开始
    fmt.Println("s2[1:4]:", s2[1:4])  // 索引1到3
    
    // append 追加元素 / Append elements
    s := []int{1, 2, 3}
    s = append(s, 4)           // 追加一个元素
    s = append(s, 5, 6, 7)     // 追加多个元素
    s = append(s, []int{8, 9}...)  // 追加另一个切片
    fmt.Println("After append:", s)
    
    // copy 复制切片 / Copy slice
    src := []int{1, 2, 3, 4, 5}
    dst := make([]int, 3)
    n := copy(dst, src)  // 复制 min(len(dst), len(src)) 个元素
    fmt.Println("Copied:", n, "elements, dst:", dst)
    
    // 删除元素 / Delete element
    s = []int{1, 2, 3, 4, 5}
    i := 2  // 删除索引2的元素
    s = append(s[:i], s[i+1:]...)
    fmt.Println("After delete:", s)  // [1 2 4 5]
    
    // 插入元素 / Insert element
    s = []int{1, 2, 4, 5}
    i = 2  // 在索引2处插入
    s = append(s[:i], append([]int{3}, s[i:]...)...)
    fmt.Println("After insert:", s)  // [1 2 3 4 5]
    
    // 切片是引用类型 / Slices are reference types
    original := []int{1, 2, 3}
    copied := original
    copied[0] = 999
    fmt.Println("original:", original)  // [999 2 3] (也被修改)
    
    // 遍历切片 / Iterate slice
    for i, v := range s2 {
        fmt.Printf("index=%d, value=%d\n", i, v)
    }
    
    // 多维切片 / Multi-dimensional slice
    matrix := [][]int{
        {1, 2, 3},
        {4, 5, 6},
        {7, 8, 9},
    }
    fmt.Println("Matrix:", matrix)
    
    // 切片的底层原理 / Slice internals
    // 切片包含三个字段: 指向底层数组的指针、长度、容量
    // A slice has three fields: pointer to underlying array, length, capacity
}
```

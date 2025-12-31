# 映射 (map) / Map

```go
package main

import (
    "fmt"
    "maps"
    "sort"
)

func main() {
    // map 是键值对的无序集合
    // map is an unordered collection of key-value pairs
    
    // 创建 map / Create map
    var m1 map[string]int              // nil map (不能直接使用)
    m2 := map[string]int{}             // 空 map
    m3 := make(map[string]int)         // make 创建
    m4 := make(map[string]int, 100)    // 指定初始容量
    m5 := map[string]int{              // 字面量初始化
        "one":   1,
        "two":   2,
        "three": 3,
    }
    
    fmt.Println("m1:", m1)
    fmt.Println("m2:", m2)
    fmt.Println("m3:", m3)
    fmt.Println("m4:", m4)
    fmt.Println("m5:", m5)
    
    // 添加/修改元素 / Add/modify elements
    m3["apple"] = 1
    m3["banana"] = 2
    m3["apple"] = 10  // 修改已存在的键
    fmt.Println("m3:", m3)
    
    // 获取元素 / Get element
    value := m3["apple"]
    fmt.Println("apple:", value)
    
    // 检查键是否存在 / Check if key exists
    value, ok := m3["orange"]
    if ok {
        fmt.Println("orange:", value)
    } else {
        fmt.Println("orange not found")
    }
    
    // 简写形式 / Short form
    if v, ok := m3["banana"]; ok {
        fmt.Println("banana:", v)
    }
    
    // 删除元素 / Delete element
    delete(m3, "apple")
    fmt.Println("After delete:", m3)
    
    // 清空 map / Clear map (Go 1.21+)
    clear(m3)
    fmt.Println("After clear:", m3)
    
    // 获取长度 / Get length
    fmt.Println("len(m5):", len(m5))
    
    // 遍历 map / Iterate map
    for key, value := range m5 {
        fmt.Printf("%s: %d\n", key, value)
    }
    
    // 只遍历键 / Iterate keys only
    for key := range m5 {
        fmt.Println("Key:", key)
    }
    
    // 按键排序遍历 / Iterate in sorted key order
    keys := make([]string, 0, len(m5))
    for k := range m5 {
        keys = append(keys, k)
    }
    sort.Strings(keys)
    for _, k := range keys {
        fmt.Printf("%s: %d\n", k, m5[k])
    }
    
    // 使用 maps 包 (Go 1.21+) / Use maps package
    m6 := map[string]int{"a": 1, "b": 2}
    m7 := map[string]int{"b": 3, "c": 4}
    
    // 复制 map / Copy map
    m6Copy := maps.Clone(m6)
    fmt.Println("Clone:", m6Copy)
    
    // 合并 map / Merge maps
    maps.Copy(m6, m7)  // m7 的值复制到 m6
    fmt.Println("After Copy:", m6)
    
    // 比较 map / Compare maps
    m8 := map[string]int{"a": 1, "b": 2}
    m9 := map[string]int{"a": 1, "b": 2}
    fmt.Println("Equal:", maps.Equal(m8, m9))
    
    // 删除满足条件的元素 / Delete elements matching condition
    maps.DeleteFunc(m8, func(k string, v int) bool {
        return v > 1
    })
    fmt.Println("After DeleteFunc:", m8)
    
    // map 的值可以是任何类型 / Map values can be any type
    funcMap := map[string]func(int) int{
        "double": func(x int) int { return x * 2 },
        "square": func(x int) int { return x * x },
    }
    fmt.Println("double(5):", funcMap["double"](5))
    fmt.Println("square(5):", funcMap["square"](5))
    
    // 嵌套 map / Nested map
    nested := map[string]map[string]int{
        "group1": {"a": 1, "b": 2},
        "group2": {"c": 3, "d": 4},
    }
    fmt.Println("nested:", nested)
    fmt.Println("group1.a:", nested["group1"]["a"])
    
    // map 是引用类型 / Maps are reference types
    original := map[string]int{"x": 1}
    ref := original
    ref["x"] = 999
    fmt.Println("original:", original)  // map[x:999]
}
```

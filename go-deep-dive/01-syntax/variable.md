# 变量声明 / Variable Declaration

```go
package main

import "fmt"

func main() {
    // 方式1: var 关键字 + 类型
    // Method 1: var keyword + type
    var name string = "Go"
    var age int = 14
    
    // 方式2: var 关键字 + 类型推断
    // Method 2: var keyword + type inference
    var language = "Golang"  // 自动推断为 string
    
    // 方式3: 短变量声明 (只能在函数内使用)
    // Method 3: Short variable declaration (only inside functions)
    version := 1.21  // 自动推断为 float64
    
    // 多变量声明 / Multiple variable declaration
    var (
        x int    = 10
        y int    = 20
        z string = "result"
    )
    
    // 并行赋值 / Parallel assignment
    a, b, c := 1, 2, "three"
    
    // 零值 / Zero values
    var defaultInt int       // 0
    var defaultFloat float64 // 0.0
    var defaultBool bool     // false
    var defaultString string // "" (空字符串)
    
    fmt.Println(name, age, language, version)
    fmt.Println(x, y, z)
    fmt.Println(a, b, c)
    fmt.Println(defaultInt, defaultFloat, defaultBool, defaultString)
}
```

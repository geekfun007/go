# Go 基础语法 / Go Basic Syntax

## 1. 程序结构 / Program Structure

```go
// 每个 Go 程序都以 package 声明开始
// Every Go program starts with a package declaration
package main

// 导入其他包 / Import other packages
import (
    "fmt"
    "math"
)

// main 函数是程序入口 / main function is the entry point
func main() {
    fmt.Println("Hello, Go!")
}
```

## 2. 变量声明 / Variable Declaration

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

## 3. 常量 / Constants

```go
package main

import "fmt"

// 常量可以在包级别声明
// Constants can be declared at package level
const Pi = 3.14159
const (
    StatusOK       = 200
    StatusNotFound = 404
)

// iota: 常量计数器，从0开始
// iota: constant counter, starts from 0
const (
    Sunday    = iota // 0
    Monday           // 1
    Tuesday          // 2
    Wednesday        // 3
    Thursday         // 4
    Friday           // 5
    Saturday         // 6
)

// iota 的高级用法 / Advanced iota usage
const (
    _  = iota             // 0, 忽略
    KB = 1 << (10 * iota) // 1 << 10 = 1024
    MB                    // 1 << 20
    GB                    // 1 << 30
    TB                    // 1 << 40
)

func main() {
    fmt.Printf("Pi = %v\n", Pi)
    fmt.Printf("HTTP Status: %d, %d\n", StatusOK, StatusNotFound)
    fmt.Printf("Days: Sunday=%d, Monday=%d\n", Sunday, Monday)
    fmt.Printf("KB=%d, MB=%d, GB=%d\n", KB, MB, GB)
}
```

## 4. 基本控制流 / Basic Control Flow

### 4.1 if 语句 / if Statement

```go
package main

import "fmt"

func main() {
    x := 10
    
    // 基本 if-else
    // Basic if-else
    if x > 5 {
        fmt.Println("x is greater than 5")
    } else if x == 5 {
        fmt.Println("x is equal to 5")
    } else {
        fmt.Println("x is less than 5")
    }
    
    // if 带初始化语句 (Go 特有)
    // if with initialization statement (Go specific)
    if y := x * 2; y > 15 {
        fmt.Println("y is greater than 15")
    }
    // 注意: y 只在 if 块内可用
    // Note: y is only available within the if block
}
```

### 4.2 for 循环 / for Loop

```go
package main

import "fmt"

func main() {
    // Go 只有 for 循环，没有 while
    // Go only has for loop, no while
    
    // 基本 for 循环
    // Basic for loop
    for i := 0; i < 5; i++ {
        fmt.Printf("i = %d\n", i)
    }
    
    // 类似 while 的用法
    // while-like usage
    j := 0
    for j < 3 {
        fmt.Printf("j = %d\n", j)
        j++
    }
    
    // 无限循环
    // Infinite loop
    // for {
    //     // 使用 break 退出
    //     break
    // }
    
    // range 遍历
    // range iteration
    nums := []int{1, 2, 3, 4, 5}
    for index, value := range nums {
        fmt.Printf("index=%d, value=%d\n", index, value)
    }
    
    // 忽略索引
    // Ignore index
    for _, value := range nums {
        fmt.Printf("value=%d\n", value)
    }
    
    // 只要索引
    // Only index
    for index := range nums {
        fmt.Printf("index=%d\n", index)
    }
}
```

### 4.3 switch 语句 / switch Statement

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    // 基本 switch
    // Basic switch
    day := "Monday"
    switch day {
    case "Monday":
        fmt.Println("周一 / Monday")
    case "Tuesday":
        fmt.Println("周二 / Tuesday")
    case "Saturday", "Sunday":  // 多值匹配
        fmt.Println("周末 / Weekend")
    default:
        fmt.Println("其他 / Other")
    }
    
    // switch 不需要 break (自动 break)
    // switch doesn't need break (auto break)
    
    // 使用 fallthrough 继续执行下一个 case
    // Use fallthrough to continue to next case
    num := 1
    switch num {
    case 1:
        fmt.Println("One")
        fallthrough
    case 2:
        fmt.Println("Two")
    }
    // 输出: One Two
    
    // 无表达式 switch (类似 if-else 链)
    // Expressionless switch (like if-else chain)
    hour := time.Now().Hour()
    switch {
    case hour < 12:
        fmt.Println("Good morning!")
    case hour < 17:
        fmt.Println("Good afternoon!")
    default:
        fmt.Println("Good evening!")
    }
    
    // 类型 switch
    // Type switch
    var i interface{} = "hello"
    switch v := i.(type) {
    case int:
        fmt.Printf("Integer: %d\n", v)
    case string:
        fmt.Printf("String: %s\n", v)
    default:
        fmt.Printf("Unknown type: %T\n", v)
    }
}
```

## 5. 指针 / Pointers

```go
package main

import "fmt"

func main() {
    // 指针存储变量的内存地址
    // Pointers store memory addresses of variables
    
    x := 10
    var p *int = &x  // p 是指向 x 的指针
    
    fmt.Println("x =", x)    // 10
    fmt.Println("p =", p)    // 内存地址
    fmt.Println("*p =", *p)  // 10 (解引用)
    
    // 通过指针修改值
    // Modify value through pointer
    *p = 20
    fmt.Println("x =", x)  // 20
    
    // 指针的零值是 nil
    // Zero value of pointer is nil
    var nilPtr *int
    fmt.Println("nilPtr =", nilPtr)  // <nil>
    
    // new 函数创建指针
    // new function creates pointer
    ptr := new(int)  // 分配内存，返回指针
    *ptr = 100
    fmt.Println("*ptr =", *ptr)
}

// 函数参数传递
// Function parameter passing
func increment(x int) {
    x++ // 不会影响原值 / Won't affect original
}

func incrementPtr(x *int) {
    (*x)++ // 会影响原值 / Will affect original
}
```

## 6. defer 语句 / defer Statement

```go
package main

import "fmt"

func main() {
    // defer 延迟执行，直到函数返回
    // defer delays execution until function returns
    
    defer fmt.Println("World")  // 最后执行
    fmt.Println("Hello")        // 先执行
    // 输出: Hello World
    
    // 多个 defer 按 LIFO (后进先出) 顺序执行
    // Multiple defers execute in LIFO order
    defer fmt.Println("1")
    defer fmt.Println("2")
    defer fmt.Println("3")
    // 输出: 3 2 1
    
    // 常用于资源清理
    // Commonly used for resource cleanup
    // file, _ := os.Open("file.txt")
    // defer file.Close()
}

// defer 在 panic 时也会执行
// defer executes even during panic
func deferWithPanic() {
    defer fmt.Println("defer executed")
    panic("something went wrong")
}
```

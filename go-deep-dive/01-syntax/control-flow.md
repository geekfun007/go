# 基本控制流 / Basic Control Flow

## 1. if 语句 / if Statement

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

## 2. for 循环 / for Loop

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

## 3. switch 语句 / switch Statement

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

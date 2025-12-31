# 常量 / Constants

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

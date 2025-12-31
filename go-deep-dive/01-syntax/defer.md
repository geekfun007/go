# defer 语句 / defer Statement

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

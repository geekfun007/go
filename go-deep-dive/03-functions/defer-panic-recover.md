# defer、panic、recover 模式 / defer, panic, recover Patterns

```go
package main

import (
    "fmt"
    "io"
    "os"
    "sync"
)

// 模式1: 资源清理 / Pattern 1: Resource cleanup
func copyFile(src, dst string) (err error) {
    srcFile, err := os.Open(src)
    if err != nil {
        return err
    }
    defer srcFile.Close()
    
    dstFile, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer func() {
        closeErr := dstFile.Close()
        if err == nil {
            err = closeErr
        }
    }()
    
    _, err = io.Copy(dstFile, srcFile)
    return err
}

// 模式2: 解锁互斥锁 / Pattern 2: Unlock mutex
func safeOperation(mu *sync.Mutex, fn func()) {
    mu.Lock()
    defer mu.Unlock()
    fn()
}

// 模式3: 计时 / Pattern 3: Timing
func timeTrack(name string) func() {
    start := time.Now()
    return func() {
        fmt.Printf("%s took %v\n", name, time.Since(start))
    }
}

func slowOperation() {
    defer timeTrack("slowOperation")()
    time.Sleep(100 * time.Millisecond)
}

// 模式4: 修改命名返回值 / Pattern 4: Modify named return
func divideWithRecovery(a, b int) (result int, err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("panic recovered: %v", r)
        }
    }()
    
    result = a / b  // 如果 b=0 会 panic
    return result, nil
}

// 模式5: 日志记录 / Pattern 5: Logging
func loggedOperation(name string) func() {
    fmt.Printf("Starting %s\n", name)
    return func() {
        fmt.Printf("Finished %s\n", name)
    }
}

// 模式6: 错误处理增强 / Pattern 6: Error handling enhancement
func processWithContext(ctx string, fn func() error) (err error) {
    defer func() {
        if err != nil {
            err = fmt.Errorf("%s: %w", ctx, err)
        }
    }()
    return fn()
}

// defer 执行顺序 / defer execution order
func deferOrder() {
    defer fmt.Println("First defer")
    defer fmt.Println("Second defer")
    defer fmt.Println("Third defer")
    fmt.Println("Function body")
    // 输出顺序: Function body -> Third defer -> Second defer -> First defer
}

func main() {
    // 模式4 示例
    result, err := divideWithRecovery(10, 0)
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("Result:", result)
    }
    
    // defer 执行顺序
    deferOrder()
    
    // 使用 loggedOperation
    defer loggedOperation("main")()
    fmt.Println("Main body")
}

// 注意事项 / Notes:
// 1. defer 的参数在声明时求值，而非执行时
// 2. 多个 defer 按 LIFO (后进先出) 顺序执行
// 3. defer 可以访问和修改命名返回值
// 4. recover 只能在 defer 函数中调用才有效
```

# Goroutine 基础 / Goroutine Basics

```go
package main

import (
    "fmt"
    "runtime"
    "sync"
    "time"
)

func main() {
    // 启动 goroutine / Start goroutine
    go func() {
        fmt.Println("Hello from goroutine!")
    }()
    
    // goroutine 是轻量级线程 / goroutines are lightweight threads
    // 初始栈大小约 2KB，可动态增长
    // Initial stack size ~2KB, can grow dynamically
    
    // 查看 CPU 核心数 / Check CPU cores
    fmt.Println("NumCPU:", runtime.NumCPU())
    
    // 设置使用的 CPU 核心数 / Set number of CPUs to use
    runtime.GOMAXPROCS(runtime.NumCPU())
    
    // 并发执行示例 / Concurrent execution example
    var wg sync.WaitGroup
    
    for i := 1; i <= 5; i++ {
        wg.Add(1)
        go func(n int) {
            defer wg.Done()
            fmt.Printf("Goroutine %d starting\n", n)
            time.Sleep(time.Duration(n) * 100 * time.Millisecond)
            fmt.Printf("Goroutine %d done\n", n)
        }(i)
    }
    
    wg.Wait()
    fmt.Println("All goroutines completed")
    
    // 查看当前 goroutine 数量 / Check current goroutine count
    fmt.Println("NumGoroutine:", runtime.NumGoroutine())
}
```


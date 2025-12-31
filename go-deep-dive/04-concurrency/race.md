# 数据竞争与调试 / Data Race & Debugging

```go
package main

import (
    "fmt"
    "sync"
)

// 使用 go run -race 或 go build -race 检测数据竞争
// Use go run -race or go build -race to detect data races

func main() {
    // 数据竞争示例 (错误) / Data race example (wrong)
    // counter := 0
    // var wg sync.WaitGroup
    // for i := 0; i < 1000; i++ {
    //     wg.Add(1)
    //     go func() {
    //         counter++  // 数据竞争!
    //         wg.Done()
    //     }()
    // }
    // wg.Wait()
    
    // 正确方式1: 使用 Mutex / Correct way 1: Use Mutex
    var mu sync.Mutex
    counter1 := 0
    var wg sync.WaitGroup
    
    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            mu.Lock()
            counter1++
            mu.Unlock()
            wg.Done()
        }()
    }
    wg.Wait()
    fmt.Println("Counter (Mutex):", counter1)
    
    // 正确方式2: 使用 channel / Correct way 2: Use channel
    counter2 := 0
    ch := make(chan int, 1000)
    
    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            ch <- 1
            wg.Done()
        }()
    }
    
    go func() {
        wg.Wait()
        close(ch)
    }()
    
    for n := range ch {
        counter2 += n
    }
    fmt.Println("Counter (Channel):", counter2)
    
    // 调试技巧 / Debugging tips:
    // 1. go run -race main.go  // 检测数据竞争
    // 2. runtime.NumGoroutine()  // 查看 goroutine 数量
    // 3. pprof 性能分析
    // 4. 使用 trace 工具: go tool trace
}
```

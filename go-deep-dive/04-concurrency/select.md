# Select 语句 / Select Statement

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    ch1 := make(chan string)
    ch2 := make(chan string)
    
    // 启动两个 goroutine / Start two goroutines
    go func() {
        time.Sleep(100 * time.Millisecond)
        ch1 <- "from channel 1"
    }()
    
    go func() {
        time.Sleep(200 * time.Millisecond)
        ch2 <- "from channel 2"
    }()
    
    // select 等待多个 channel / select waits on multiple channels
    for i := 0; i < 2; i++ {
        select {
        case msg1 := <-ch1:
            fmt.Println("Received:", msg1)
        case msg2 := <-ch2:
            fmt.Println("Received:", msg2)
        }
    }
    
    // 带超时的 select / select with timeout
    ch3 := make(chan string)
    
    go func() {
        time.Sleep(2 * time.Second)
        ch3 <- "result"
    }()
    
    select {
    case result := <-ch3:
        fmt.Println("Got result:", result)
    case <-time.After(1 * time.Second):
        fmt.Println("Timeout!")
    }
    
    // 非阻塞 select (使用 default) / Non-blocking select (using default)
    ch4 := make(chan int, 1)
    
    select {
    case v := <-ch4:
        fmt.Println("Received:", v)
    default:
        fmt.Println("No value available")
    }
    
    // 用 select 实现超时控制 / Implement timeout with select
    done := make(chan bool)
    
    go func() {
        time.Sleep(500 * time.Millisecond)
        done <- true
    }()
    
    select {
    case <-done:
        fmt.Println("Task completed")
    case <-time.After(1 * time.Second):
        fmt.Println("Task timed out")
    }
    
    // select 随机选择 / select chooses randomly
    ch5 := make(chan int, 1)
    ch6 := make(chan int, 1)
    ch5 <- 1
    ch6 <- 2
    
    for i := 0; i < 10; i++ {
        select {
        case v := <-ch5:
            fmt.Print(v)
            ch5 <- 1
        case v := <-ch6:
            fmt.Print(v)
            ch6 <- 2
        }
    }
    fmt.Println()
}
```


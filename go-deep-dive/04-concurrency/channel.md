# Channel 基础 / Channel Basics

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    // 创建 channel / Create channel
    // 无缓冲 channel / Unbuffered channel
    ch1 := make(chan int)
    
    // 有缓冲 channel / Buffered channel
    ch2 := make(chan string, 3)
    
    // 无缓冲 channel - 发送和接收必须同时准备好
    // Unbuffered channel - send and receive must be ready simultaneously
    go func() {
        ch1 <- 42  // 发送 / send
    }()
    value := <-ch1  // 接收 / receive
    fmt.Println("Received:", value)
    
    // 有缓冲 channel - 缓冲区满之前发送不会阻塞
    // Buffered channel - send doesn't block until buffer is full
    ch2 <- "first"
    ch2 <- "second"
    ch2 <- "third"
    // ch2 <- "fourth"  // 这会阻塞，因为缓冲区已满
    
    fmt.Println(<-ch2)  // first
    fmt.Println(<-ch2)  // second
    fmt.Println(<-ch2)  // third
    
    // Channel 长度和容量 / Channel length and capacity
    ch3 := make(chan int, 5)
    ch3 <- 1
    ch3 <- 2
    fmt.Printf("len: %d, cap: %d\n", len(ch3), cap(ch3))  // len: 2, cap: 5
    
    // 关闭 channel / Close channel
    ch4 := make(chan int, 3)
    ch4 <- 1
    ch4 <- 2
    ch4 <- 3
    close(ch4)  // 关闭后不能再发送，但可以继续接收
    
    // 检查 channel 是否关闭 / Check if channel is closed
    for {
        v, ok := <-ch4
        if !ok {
            fmt.Println("Channel closed")
            break
        }
        fmt.Println("Value:", v)
    }
    
    // 使用 range 遍历 channel / Iterate channel with range
    ch5 := make(chan int, 3)
    ch5 <- 10
    ch5 <- 20
    ch5 <- 30
    close(ch5)
    
    for v := range ch5 {
        fmt.Println("Range value:", v)
    }
    
    // 单向 channel / Directional channels
    sendOnly := make(chan<- int)  // 只能发送
    recvOnly := make(<-chan int)  // 只能接收
    _ = sendOnly
    _ = recvOnly
    
    // 通常用于函数签名 / Usually used in function signatures
    produce := func(ch chan<- int) {
        for i := 0; i < 3; i++ {
            ch <- i
        }
        close(ch)
    }
    
    consume := func(ch <-chan int) {
        for v := range ch {
            fmt.Println("Consumed:", v)
        }
    }
    
    ch6 := make(chan int)
    go produce(ch6)
    consume(ch6)
}
```


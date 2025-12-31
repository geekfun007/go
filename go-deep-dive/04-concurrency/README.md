# Go 并发编程 / Go Concurrency

## 1. Goroutine 基础 / Goroutine Basics

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

## 2. Channel 基础 / Channel Basics

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

## 3. Select 语句 / Select Statement

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

## 4. 同步原语 / Synchronization Primitives

```go
package main

import (
    "fmt"
    "sync"
    "sync/atomic"
    "time"
)

func main() {
    // WaitGroup - 等待一组 goroutine 完成
    // WaitGroup - wait for a group of goroutines to complete
    var wg sync.WaitGroup
    
    for i := 0; i < 5; i++ {
        wg.Add(1)
        go func(n int) {
            defer wg.Done()
            fmt.Printf("Worker %d done\n", n)
        }(i)
    }
    wg.Wait()
    fmt.Println("All workers completed")
    
    // Mutex - 互斥锁 / Mutex - mutual exclusion lock
    var mu sync.Mutex
    counter := 0
    
    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            mu.Lock()
            counter++
            mu.Unlock()
        }()
    }
    wg.Wait()
    fmt.Println("Counter (Mutex):", counter)
    
    // RWMutex - 读写锁 / RWMutex - read-write lock
    var rwmu sync.RWMutex
    data := make(map[string]int)
    
    // 写操作 / Write operation
    write := func(key string, value int) {
        rwmu.Lock()
        defer rwmu.Unlock()
        data[key] = value
    }
    
    // 读操作 / Read operation
    read := func(key string) int {
        rwmu.RLock()
        defer rwmu.RUnlock()
        return data[key]
    }
    
    write("key1", 100)
    fmt.Println("Read key1:", read("key1"))
    
    // Once - 只执行一次 / Once - execute only once
    var once sync.Once
    
    for i := 0; i < 5; i++ {
        wg.Add(1)
        go func(n int) {
            defer wg.Done()
            once.Do(func() {
                fmt.Printf("Only once, called by goroutine %d\n", n)
            })
        }(i)
    }
    wg.Wait()
    
    // Cond - 条件变量 / Cond - condition variable
    var cond = sync.NewCond(&sync.Mutex{})
    ready := false
    
    go func() {
        time.Sleep(100 * time.Millisecond)
        cond.L.Lock()
        ready = true
        cond.L.Unlock()
        cond.Broadcast()  // 唤醒所有等待的 goroutine
    }()
    
    cond.L.Lock()
    for !ready {
        cond.Wait()  // 等待条件满足
    }
    fmt.Println("Condition met!")
    cond.L.Unlock()
    
    // Atomic - 原子操作 / Atomic - atomic operations
    var atomicCounter int64 = 0
    
    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            atomic.AddInt64(&atomicCounter, 1)
        }()
    }
    wg.Wait()
    fmt.Println("Counter (Atomic):", atomicCounter)
    
    // atomic.Value - 存储任意值 / atomic.Value - store any value
    var config atomic.Value
    config.Store(map[string]string{"env": "production"})
    
    go func() {
        cfg := config.Load().(map[string]string)
        fmt.Println("Config env:", cfg["env"])
    }()
    
    time.Sleep(100 * time.Millisecond)
    
    // Pool - 对象池 / Pool - object pool
    pool := &sync.Pool{
        New: func() interface{} {
            return make([]byte, 1024)
        },
    }
    
    buf := pool.Get().([]byte)
    // 使用 buffer...
    pool.Put(buf)  // 归还到池中
}
```

## 5. 并发模式 / Concurrency Patterns

### 5.1 Worker Pool 模式

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

// Job 代表一个任务 / Job represents a task
type Job struct {
    ID   int
    Data string
}

// Result 代表任务结果 / Result represents task result
type Result struct {
    JobID  int
    Output string
}

// Worker Pool 实现 / Worker Pool implementation
func workerPool(numWorkers int, jobs <-chan Job, results chan<- Result) {
    var wg sync.WaitGroup
    
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func(workerID int) {
            defer wg.Done()
            for job := range jobs {
                // 处理任务 / Process job
                time.Sleep(100 * time.Millisecond)
                results <- Result{
                    JobID:  job.ID,
                    Output: fmt.Sprintf("Worker %d processed: %s", workerID, job.Data),
                }
            }
        }(i)
    }
    
    wg.Wait()
    close(results)
}

func main() {
    jobs := make(chan Job, 100)
    results := make(chan Result, 100)
    
    // 启动 worker pool / Start worker pool
    go workerPool(3, jobs, results)
    
    // 发送任务 / Send jobs
    for i := 1; i <= 10; i++ {
        jobs <- Job{ID: i, Data: fmt.Sprintf("task-%d", i)}
    }
    close(jobs)
    
    // 收集结果 / Collect results
    for result := range results {
        fmt.Printf("Result: Job %d -> %s\n", result.JobID, result.Output)
    }
}
```

### 5.2 Fan-Out/Fan-In 模式

```go
package main

import (
    "fmt"
    "sync"
)

// Fan-Out: 一个 channel 分发到多个 goroutine
// Fan-In: 多个 channel 合并到一个 channel

func fanOut(input <-chan int, numWorkers int) []<-chan int {
    outputs := make([]<-chan int, numWorkers)
    
    for i := 0; i < numWorkers; i++ {
        output := make(chan int)
        outputs[i] = output
        
        go func(out chan<- int) {
            for n := range input {
                out <- n * n  // 处理：计算平方
            }
            close(out)
        }(output)
    }
    
    return outputs
}

func fanIn(inputs ...<-chan int) <-chan int {
    output := make(chan int)
    var wg sync.WaitGroup
    
    for _, input := range inputs {
        wg.Add(1)
        go func(in <-chan int) {
            defer wg.Done()
            for n := range in {
                output <- n
            }
        }(input)
    }
    
    go func() {
        wg.Wait()
        close(output)
    }()
    
    return output
}

func main() {
    // 创建输入 channel / Create input channel
    input := make(chan int)
    go func() {
        for i := 1; i <= 10; i++ {
            input <- i
        }
        close(input)
    }()
    
    // Fan-Out 到 3 个 worker
    workers := fanOut(input, 3)
    
    // Fan-In 合并结果
    results := fanIn(workers...)
    
    // 收集结果
    for result := range results {
        fmt.Println("Result:", result)
    }
}
```

### 5.3 Pipeline 模式

```go
package main

import "fmt"

// Pipeline 阶段 / Pipeline stages

// 生成数字 / Generate numbers
func generate(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        for _, n := range nums {
            out <- n
        }
        close(out)
    }()
    return out
}

// 平方 / Square
func square(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            out <- n * n
        }
        close(out)
    }()
    return out
}

// 过滤偶数 / Filter even
func filterEven(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            if n%2 == 0 {
                out <- n
            }
        }
        close(out)
    }()
    return out
}

// 打印 / Print
func printer(in <-chan int) {
    for n := range in {
        fmt.Println("Output:", n)
    }
}

func main() {
    // 构建 pipeline
    // generate -> square -> filterEven -> print
    
    nums := generate(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
    squared := square(nums)
    evens := filterEven(squared)
    printer(evens)
}
```

### 5.4 取消和超时模式 / Cancellation & Timeout Patterns

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func main() {
    // Context 取消 / Context cancellation
    ctx, cancel := context.WithCancel(context.Background())
    
    go func() {
        for {
            select {
            case <-ctx.Done():
                fmt.Println("Worker cancelled")
                return
            default:
                fmt.Println("Working...")
                time.Sleep(200 * time.Millisecond)
            }
        }
    }()
    
    time.Sleep(500 * time.Millisecond)
    cancel()  // 取消
    time.Sleep(100 * time.Millisecond)
    
    // Context 超时 / Context timeout
    ctx2, cancel2 := context.WithTimeout(context.Background(), 300*time.Millisecond)
    defer cancel2()
    
    select {
    case <-time.After(500 * time.Millisecond):
        fmt.Println("Operation completed")
    case <-ctx2.Done():
        fmt.Println("Operation timed out:", ctx2.Err())
    }
    
    // Context 截止时间 / Context deadline
    deadline := time.Now().Add(200 * time.Millisecond)
    ctx3, cancel3 := context.WithDeadline(context.Background(), deadline)
    defer cancel3()
    
    select {
    case <-time.After(500 * time.Millisecond):
        fmt.Println("Completed before deadline")
    case <-ctx3.Done():
        fmt.Println("Deadline exceeded:", ctx3.Err())
    }
    
    // 使用 context 传递值 / Pass values with context
    ctx4 := context.WithValue(context.Background(), "userID", 12345)
    
    processRequest := func(ctx context.Context) {
        if userID := ctx.Value("userID"); userID != nil {
            fmt.Println("Processing request for user:", userID)
        }
    }
    
    processRequest(ctx4)
}
```

### 5.5 错误组模式 / Error Group Pattern

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "golang.org/x/sync/errgroup"
)

func main() {
    // errgroup - 带错误处理的并发
    // errgroup - concurrency with error handling
    
    g, ctx := errgroup.WithContext(context.Background())
    
    // 启动多个并发任务 / Start multiple concurrent tasks
    for i := 1; i <= 3; i++ {
        i := i
        g.Go(func() error {
            select {
            case <-ctx.Done():
                return ctx.Err()
            case <-time.After(time.Duration(i) * 100 * time.Millisecond):
                if i == 2 {
                    return fmt.Errorf("task %d failed", i)
                }
                fmt.Printf("Task %d completed\n", i)
                return nil
            }
        })
    }
    
    // 等待所有任务完成或第一个错误
    // Wait for all tasks or first error
    if err := g.Wait(); err != nil {
        fmt.Println("Error:", err)
    }
    
    // 带限制的 errgroup / errgroup with limit
    g2, _ := errgroup.WithContext(context.Background())
    g2.SetLimit(2)  // 最多 2 个并发
    
    for i := 0; i < 5; i++ {
        i := i
        g2.Go(func() error {
            fmt.Printf("Limited task %d running\n", i)
            time.Sleep(100 * time.Millisecond)
            return nil
        })
    }
    g2.Wait()
}
```

### 5.6 信号量模式 / Semaphore Pattern

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "golang.org/x/sync/semaphore"
)

func main() {
    // 使用 channel 实现信号量 / Implement semaphore with channel
    sem := make(chan struct{}, 3)  // 最多 3 个并发
    
    for i := 0; i < 10; i++ {
        sem <- struct{}{}  // 获取信号量
        go func(n int) {
            defer func() { <-sem }()  // 释放信号量
            fmt.Printf("Task %d running\n", n)
            time.Sleep(100 * time.Millisecond)
        }(i)
    }
    
    // 等待所有完成
    for i := 0; i < cap(sem); i++ {
        sem <- struct{}{}
    }
    
    fmt.Println("\n--- Using golang.org/x/sync/semaphore ---")
    
    // 使用 semaphore 包 / Using semaphore package
    sem2 := semaphore.NewWeighted(3)
    ctx := context.Background()
    
    for i := 0; i < 10; i++ {
        if err := sem2.Acquire(ctx, 1); err != nil {
            break
        }
        go func(n int) {
            defer sem2.Release(1)
            fmt.Printf("Weighted task %d running\n", n)
            time.Sleep(100 * time.Millisecond)
        }(i)
    }
    
    // 等待所有完成
    sem2.Acquire(ctx, 3)
}
```

## 6. 数据竞争与调试 / Data Race & Debugging

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

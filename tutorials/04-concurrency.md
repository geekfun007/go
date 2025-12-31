# Go 并发编程详解

## 1. Goroutine 基础

### 1.1 创建 Goroutine

```go
package main

import (
    "fmt"
    "runtime"
    "sync"
    "time"
)

func sayHello(name string) {
    for i := 0; i < 3; i++ {
        fmt.Printf("Hello, %s! (%d)\n", name, i)
        time.Sleep(100 * time.Millisecond)
    }
}

func main() {
    // 设置 GOMAXPROCS
    fmt.Println("CPUs:", runtime.NumCPU())
    fmt.Println("GOMAXPROCS:", runtime.GOMAXPROCS(0))
    
    // 基本 goroutine
    go sayHello("Alice")
    go sayHello("Bob")
    
    // 匿名函数 goroutine
    go func() {
        fmt.Println("Anonymous goroutine")
    }()
    
    // 带参数的匿名函数
    message := "Hello from closure"
    go func(msg string) {
        fmt.Println(msg)
    }(message)
    
    // 等待 goroutine 完成 (简单方式)
    time.Sleep(500 * time.Millisecond)
    
    // 使用 WaitGroup
    var wg sync.WaitGroup
    
    for i := 0; i < 5; i++ {
        wg.Add(1)
        go func(n int) {
            defer wg.Done()
            fmt.Printf("Worker %d starting\n", n)
            time.Sleep(100 * time.Millisecond)
            fmt.Printf("Worker %d done\n", n)
        }(i)
    }
    
    wg.Wait()
    fmt.Println("All workers done")
    
    // 查看 goroutine 数量
    fmt.Println("Goroutines:", runtime.NumGoroutine())
}
```

### 1.2 Goroutine 泄漏

```go
package main

import (
    "context"
    "fmt"
    "runtime"
    "time"
)

func main() {
    // 演示 goroutine 泄漏
    fmt.Println("=== Goroutine 泄漏示例 ===")
    leakyFunction()
    time.Sleep(100 * time.Millisecond)
    fmt.Println("After leak, goroutines:", runtime.NumGoroutine())
    
    // 正确处理
    fmt.Println("\n=== 正确处理 ===")
    properFunction()
    time.Sleep(100 * time.Millisecond)
    fmt.Println("After proper, goroutines:", runtime.NumGoroutine())
    
    // 使用 context 取消
    fmt.Println("\n=== 使用 Context ===")
    contextExample()
}

// 泄漏示例 - 通道没有接收者
func leakyFunction() {
    ch := make(chan int)
    
    go func() {
        // 这个 goroutine 会永远阻塞
        ch <- 1
        fmt.Println("This will never print")
    }()
    
    // 没有人读取 ch，goroutine 泄漏
}

// 正确处理 - 使用带缓冲的通道或确保接收
func properFunction() {
    ch := make(chan int, 1) // 带缓冲
    
    go func() {
        ch <- 1
        fmt.Println("Sent successfully")
    }()
    
    // 或者使用 select + done
    done := make(chan struct{})
    go func() {
        select {
        case <-done:
            return
        default:
            // 工作
        }
    }()
    close(done)
}

// 使用 context 取消
func contextExample() {
    ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
    defer cancel()
    
    resultCh := make(chan string, 1)
    
    go func() {
        // 模拟长时间操作
        time.Sleep(500 * time.Millisecond)
        select {
        case resultCh <- "result":
        case <-ctx.Done():
            fmt.Println("Worker cancelled")
            return
        }
    }()
    
    select {
    case result := <-resultCh:
        fmt.Println("Got result:", result)
    case <-ctx.Done():
        fmt.Println("Timeout!")
    }
}
```

## 2. Channel 通道

### 2.1 Channel 基础

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    // 无缓冲通道
    ch1 := make(chan int)
    
    go func() {
        ch1 <- 42
        fmt.Println("Sent 42")
    }()
    
    value := <-ch1
    fmt.Println("Received:", value)
    
    // 带缓冲通道
    ch2 := make(chan int, 3)
    ch2 <- 1
    ch2 <- 2
    ch2 <- 3
    // ch2 <- 4 // 会阻塞
    
    fmt.Println(<-ch2, <-ch2, <-ch2)
    
    // 通道方向
    sendOnly := make(chan<- int, 1) // 只发送
    recvOnly := make(<-chan int, 1) // 只接收
    
    sendOnly <- 1
    _ = recvOnly
    
    // 关闭通道
    ch3 := make(chan int, 3)
    ch3 <- 1
    ch3 <- 2
    ch3 <- 3
    close(ch3)
    
    // 读取已关闭的通道
    for v := range ch3 {
        fmt.Println("From closed:", v)
    }
    
    // 检测通道是否关闭
    ch4 := make(chan int)
    close(ch4)
    v, ok := <-ch4
    fmt.Printf("value: %d, ok: %v\n", v, ok)
    
    // nil 通道
    var nilCh chan int
    // nilCh <- 1    // 永久阻塞
    // <-nilCh       // 永久阻塞
    // close(nilCh)  // panic
    _ = nilCh
}
```

### 2.2 Select 语句

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    ch1 := make(chan string)
    ch2 := make(chan string)
    
    go func() {
        time.Sleep(100 * time.Millisecond)
        ch1 <- "one"
    }()
    
    go func() {
        time.Sleep(200 * time.Millisecond)
        ch2 <- "two"
    }()
    
    // 基本 select
    for i := 0; i < 2; i++ {
        select {
        case msg1 := <-ch1:
            fmt.Println("Received", msg1)
        case msg2 := <-ch2:
            fmt.Println("Received", msg2)
        }
    }
    
    // select 带超时
    ch3 := make(chan string)
    go func() {
        time.Sleep(2 * time.Second)
        ch3 <- "result"
    }()
    
    select {
    case res := <-ch3:
        fmt.Println("Result:", res)
    case <-time.After(500 * time.Millisecond):
        fmt.Println("Timeout!")
    }
    
    // select 带 default (非阻塞)
    ch4 := make(chan int, 1)
    
    select {
    case v := <-ch4:
        fmt.Println("Received:", v)
    default:
        fmt.Println("No value available")
    }
    
    select {
    case ch4 <- 42:
        fmt.Println("Sent 42")
    default:
        fmt.Println("Channel full")
    }
    
    // 使用 select 进行超时控制
    fmt.Println("\n=== 超时控制 ===")
    timeoutExample()
    
    // 使用 select 进行取消
    fmt.Println("\n=== 取消操作 ===")
    cancelExample()
}

func timeoutExample() {
    result := make(chan string, 1)
    
    go func() {
        time.Sleep(200 * time.Millisecond)
        result <- "operation complete"
    }()
    
    select {
    case res := <-result:
        fmt.Println(res)
    case <-time.After(100 * time.Millisecond):
        fmt.Println("Operation timed out")
    }
}

func cancelExample() {
    done := make(chan struct{})
    result := make(chan int)
    
    go func() {
        for {
            select {
            case <-done:
                fmt.Println("Worker cancelled")
                return
            default:
                // 工作...
                result <- 1
                time.Sleep(50 * time.Millisecond)
            }
        }
    }()
    
    // 接收一些结果
    for i := 0; i < 3; i++ {
        fmt.Println("Result:", <-result)
    }
    
    // 取消
    close(done)
    time.Sleep(100 * time.Millisecond)
}
```

### 2.3 Channel 模式

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

func main() {
    // 生成器模式
    fmt.Println("=== Generator ===")
    for n := range generator(5) {
        fmt.Println(n)
    }
    
    // 扇出模式
    fmt.Println("\n=== Fan-Out ===")
    fanOutExample()
    
    // 扇入模式
    fmt.Println("\n=== Fan-In ===")
    fanInExample()
    
    // Pipeline 模式
    fmt.Println("\n=== Pipeline ===")
    pipelineExample()
    
    // 信号量模式
    fmt.Println("\n=== Semaphore ===")
    semaphoreExample()
    
    // Or-Done 模式
    fmt.Println("\n=== Or-Done ===")
    orDoneExample()
}

// 生成器模式
func generator(max int) <-chan int {
    ch := make(chan int)
    go func() {
        defer close(ch)
        for i := 0; i < max; i++ {
            ch <- i
        }
    }()
    return ch
}

// 扇出 - 多个 goroutine 从同一个 channel 读取
func fanOutExample() {
    jobs := make(chan int, 10)
    var wg sync.WaitGroup
    
    // 启动多个 worker
    for w := 0; w < 3; w++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            for job := range jobs {
                fmt.Printf("Worker %d processing job %d\n", id, job)
                time.Sleep(50 * time.Millisecond)
            }
        }(w)
    }
    
    // 发送任务
    for j := 0; j < 9; j++ {
        jobs <- j
    }
    close(jobs)
    
    wg.Wait()
}

// 扇入 - 多个 channel 合并到一个
func fanInExample() {
    ch1 := make(chan string)
    ch2 := make(chan string)
    ch3 := make(chan string)
    
    go func() {
        for i := 0; i < 3; i++ {
            ch1 <- fmt.Sprintf("ch1: %d", i)
            time.Sleep(50 * time.Millisecond)
        }
        close(ch1)
    }()
    
    go func() {
        for i := 0; i < 3; i++ {
            ch2 <- fmt.Sprintf("ch2: %d", i)
            time.Sleep(70 * time.Millisecond)
        }
        close(ch2)
    }()
    
    go func() {
        for i := 0; i < 3; i++ {
            ch3 <- fmt.Sprintf("ch3: %d", i)
            time.Sleep(90 * time.Millisecond)
        }
        close(ch3)
    }()
    
    merged := fanIn(ch1, ch2, ch3)
    for msg := range merged {
        fmt.Println(msg)
    }
}

func fanIn(channels ...<-chan string) <-chan string {
    var wg sync.WaitGroup
    merged := make(chan string)
    
    output := func(ch <-chan string) {
        defer wg.Done()
        for msg := range ch {
            merged <- msg
        }
    }
    
    wg.Add(len(channels))
    for _, ch := range channels {
        go output(ch)
    }
    
    go func() {
        wg.Wait()
        close(merged)
    }()
    
    return merged
}

// Pipeline 模式
func pipelineExample() {
    // stage 1: 生成数字
    nums := make(chan int)
    go func() {
        defer close(nums)
        for i := 1; i <= 5; i++ {
            nums <- i
        }
    }()
    
    // stage 2: 平方
    squares := make(chan int)
    go func() {
        defer close(squares)
        for n := range nums {
            squares <- n * n
        }
    }()
    
    // stage 3: 打印
    for sq := range squares {
        fmt.Println("Square:", sq)
    }
}

// 信号量模式 - 限制并发数
func semaphoreExample() {
    sem := make(chan struct{}, 3) // 最多3个并发
    var wg sync.WaitGroup
    
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            
            sem <- struct{}{} // 获取信号量
            defer func() { <-sem }() // 释放信号量
            
            fmt.Printf("Worker %d starting\n", id)
            time.Sleep(100 * time.Millisecond)
            fmt.Printf("Worker %d done\n", id)
        }(i)
    }
    
    wg.Wait()
}

// Or-Done 模式
func orDoneExample() {
    done := make(chan struct{})
    
    go func() {
        time.Sleep(200 * time.Millisecond)
        close(done)
    }()
    
    stream := make(chan int)
    go func() {
        defer close(stream)
        for i := 0; ; i++ {
            select {
            case <-done:
                return
            case stream <- i:
                time.Sleep(50 * time.Millisecond)
            }
        }
    }()
    
    for v := range orDone(done, stream) {
        fmt.Println("Value:", v)
    }
}

func orDone[T any](done <-chan struct{}, stream <-chan T) <-chan T {
    valStream := make(chan T)
    go func() {
        defer close(valStream)
        for {
            select {
            case <-done:
                return
            case v, ok := <-stream:
                if !ok {
                    return
                }
                select {
                case valStream <- v:
                case <-done:
                    return
                }
            }
        }
    }()
    return valStream
}
```

## 3. 同步原语

### 3.1 Mutex 互斥锁

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

// 计数器 - 使用 Mutex
type SafeCounter struct {
    mu    sync.Mutex
    value int
}

func (c *SafeCounter) Inc() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.value++
}

func (c *SafeCounter) Value() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.value
}

// 读写锁
type SafeCache struct {
    mu    sync.RWMutex
    cache map[string]string
}

func NewSafeCache() *SafeCache {
    return &SafeCache{cache: make(map[string]string)}
}

func (c *SafeCache) Get(key string) (string, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    v, ok := c.cache[key]
    return v, ok
}

func (c *SafeCache) Set(key, value string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.cache[key] = value
}

func (c *SafeCache) Delete(key string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    delete(c.cache, key)
}

func main() {
    // Mutex 示例
    counter := &SafeCounter{}
    var wg sync.WaitGroup
    
    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            counter.Inc()
        }()
    }
    
    wg.Wait()
    fmt.Println("Counter:", counter.Value())
    
    // RWMutex 示例
    cache := NewSafeCache()
    
    // 写入
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(n int) {
            defer wg.Done()
            key := fmt.Sprintf("key%d", n)
            cache.Set(key, fmt.Sprintf("value%d", n))
        }(i)
    }
    
    wg.Wait()
    
    // 并发读取
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(n int) {
            defer wg.Done()
            key := fmt.Sprintf("key%d", n)
            if v, ok := cache.Get(key); ok {
                fmt.Println(key, "=", v)
            }
        }(i)
    }
    
    wg.Wait()
    
    // 避免死锁示例
    fmt.Println("\n=== 避免死锁 ===")
    avoidDeadlock()
}

// 死锁避免
func avoidDeadlock() {
    var mu1, mu2 sync.Mutex
    
    // 总是按相同顺序获取锁
    go func() {
        mu1.Lock()
        defer mu1.Unlock()
        time.Sleep(10 * time.Millisecond)
        mu2.Lock()
        defer mu2.Unlock()
        fmt.Println("Goroutine 1 done")
    }()
    
    go func() {
        mu1.Lock() // 与上面相同顺序
        defer mu1.Unlock()
        time.Sleep(10 * time.Millisecond)
        mu2.Lock()
        defer mu2.Unlock()
        fmt.Println("Goroutine 2 done")
    }()
    
    time.Sleep(100 * time.Millisecond)
}
```

### 3.2 其他同步原语

```go
package main

import (
    "fmt"
    "sync"
    "sync/atomic"
    "time"
)

func main() {
    // WaitGroup
    fmt.Println("=== WaitGroup ===")
    waitGroupExample()
    
    // Once
    fmt.Println("\n=== Once ===")
    onceExample()
    
    // Cond
    fmt.Println("\n=== Cond ===")
    condExample()
    
    // Pool
    fmt.Println("\n=== Pool ===")
    poolExample()
    
    // Atomic
    fmt.Println("\n=== Atomic ===")
    atomicExample()
    
    // Map
    fmt.Println("\n=== sync.Map ===")
    syncMapExample()
}

func waitGroupExample() {
    var wg sync.WaitGroup
    
    for i := 0; i < 5; i++ {
        wg.Add(1)
        go func(n int) {
            defer wg.Done()
            fmt.Printf("Worker %d\n", n)
        }(i)
    }
    
    wg.Wait()
    fmt.Println("All workers done")
}

func onceExample() {
    var once sync.Once
    var config string
    
    initConfig := func() {
        fmt.Println("Initializing config...")
        config = "initialized"
    }
    
    // 多次调用，只执行一次
    var wg sync.WaitGroup
    for i := 0; i < 5; i++ {
        wg.Add(1)
        go func(n int) {
            defer wg.Done()
            once.Do(initConfig)
            fmt.Printf("Goroutine %d: config = %s\n", n, config)
        }(i)
    }
    
    wg.Wait()
}

func condExample() {
    var mu sync.Mutex
    cond := sync.NewCond(&mu)
    ready := false
    
    // 等待条件
    go func() {
        mu.Lock()
        for !ready {
            cond.Wait()
        }
        fmt.Println("Worker: condition met, proceeding")
        mu.Unlock()
    }()
    
    // 设置条件并通知
    time.Sleep(100 * time.Millisecond)
    mu.Lock()
    ready = true
    cond.Signal() // 或 cond.Broadcast() 通知所有
    mu.Unlock()
    
    time.Sleep(100 * time.Millisecond)
}

func poolExample() {
    pool := &sync.Pool{
        New: func() interface{} {
            fmt.Println("Creating new object")
            return make([]byte, 1024)
        },
    }
    
    // 获取对象
    obj1 := pool.Get().([]byte)
    fmt.Printf("Got object: len=%d\n", len(obj1))
    
    // 放回对象
    pool.Put(obj1)
    
    // 再次获取 (可能复用)
    obj2 := pool.Get().([]byte)
    fmt.Printf("Got object: len=%d\n", len(obj2))
    
    // 并发使用
    var wg sync.WaitGroup
    for i := 0; i < 5; i++ {
        wg.Add(1)
        go func(n int) {
            defer wg.Done()
            buf := pool.Get().([]byte)
            defer pool.Put(buf)
            // 使用 buf...
            fmt.Printf("Goroutine %d got buffer\n", n)
        }(i)
    }
    wg.Wait()
}

func atomicExample() {
    var counter int64 = 0
    var wg sync.WaitGroup
    
    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            atomic.AddInt64(&counter, 1)
        }()
    }
    
    wg.Wait()
    fmt.Println("Atomic counter:", atomic.LoadInt64(&counter))
    
    // CompareAndSwap
    var value int32 = 100
    swapped := atomic.CompareAndSwapInt32(&value, 100, 200)
    fmt.Printf("CAS: swapped=%v, value=%d\n", swapped, value)
    
    // atomic.Value
    var config atomic.Value
    config.Store(map[string]string{"key": "value"})
    
    cfg := config.Load().(map[string]string)
    fmt.Println("Config:", cfg)
}

func syncMapExample() {
    var m sync.Map
    
    // Store
    m.Store("key1", "value1")
    m.Store("key2", "value2")
    
    // Load
    if v, ok := m.Load("key1"); ok {
        fmt.Println("key1:", v)
    }
    
    // LoadOrStore
    v, loaded := m.LoadOrStore("key3", "value3")
    fmt.Printf("key3: %v, loaded: %v\n", v, loaded)
    
    // Range
    m.Range(func(key, value interface{}) bool {
        fmt.Printf("%v: %v\n", key, value)
        return true
    })
    
    // Delete
    m.Delete("key1")
}
```

## 4. Context 上下文

```go
package main

import (
    "context"
    "fmt"
    "net/http"
    "time"
)

func main() {
    // 基本 context
    fmt.Println("=== Basic Context ===")
    basicContext()
    
    // WithCancel
    fmt.Println("\n=== WithCancel ===")
    cancelContext()
    
    // WithTimeout
    fmt.Println("\n=== WithTimeout ===")
    timeoutContext()
    
    // WithDeadline
    fmt.Println("\n=== WithDeadline ===")
    deadlineContext()
    
    // WithValue
    fmt.Println("\n=== WithValue ===")
    valueContext()
    
    // HTTP 请求 context
    fmt.Println("\n=== HTTP Context ===")
    httpContext()
}

func basicContext() {
    ctx := context.Background()
    fmt.Printf("Background: %v\n", ctx)
    
    ctx = context.TODO()
    fmt.Printf("TODO: %v\n", ctx)
}

func cancelContext() {
    ctx, cancel := context.WithCancel(context.Background())
    
    go func() {
        for {
            select {
            case <-ctx.Done():
                fmt.Println("Worker cancelled:", ctx.Err())
                return
            default:
                fmt.Println("Working...")
                time.Sleep(50 * time.Millisecond)
            }
        }
    }()
    
    time.Sleep(150 * time.Millisecond)
    cancel()
    time.Sleep(50 * time.Millisecond)
}

func timeoutContext() {
    ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
    defer cancel()
    
    select {
    case <-time.After(200 * time.Millisecond):
        fmt.Println("Operation completed")
    case <-ctx.Done():
        fmt.Println("Timeout:", ctx.Err())
    }
}

func deadlineContext() {
    deadline := time.Now().Add(100 * time.Millisecond)
    ctx, cancel := context.WithDeadline(context.Background(), deadline)
    defer cancel()
    
    select {
    case <-time.After(200 * time.Millisecond):
        fmt.Println("Operation completed")
    case <-ctx.Done():
        fmt.Println("Deadline exceeded:", ctx.Err())
    }
}

type contextKey string

func valueContext() {
    ctx := context.WithValue(context.Background(), contextKey("userID"), 12345)
    ctx = context.WithValue(ctx, contextKey("requestID"), "abc-123")
    
    processRequest(ctx)
}

func processRequest(ctx context.Context) {
    userID := ctx.Value(contextKey("userID"))
    requestID := ctx.Value(contextKey("requestID"))
    
    fmt.Printf("Processing request %v for user %v\n", requestID, userID)
}

func httpContext() {
    handler := func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()
        
        select {
        case <-time.After(2 * time.Second):
            fmt.Fprintln(w, "Hello, World!")
        case <-ctx.Done():
            fmt.Println("Request cancelled:", ctx.Err())
            http.Error(w, "Request cancelled", http.StatusRequestTimeout)
        }
    }
    
    // 模拟请求
    req, _ := http.NewRequest("GET", "/", nil)
    ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
    defer cancel()
    req = req.WithContext(ctx)
    
    fmt.Println("Request context deadline:", ctx.Err())
    _ = handler
}
```

## 5. 并发模式

### 5.1 Worker Pool

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

type Job struct {
    ID   int
    Data string
}

type Result struct {
    JobID  int
    Output string
}

func worker(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
    defer wg.Done()
    for job := range jobs {
        fmt.Printf("Worker %d processing job %d\n", id, job.ID)
        time.Sleep(100 * time.Millisecond) // 模拟工作
        results <- Result{
            JobID:  job.ID,
            Output: fmt.Sprintf("Processed: %s", job.Data),
        }
    }
}

func main() {
    const numWorkers = 3
    const numJobs = 10
    
    jobs := make(chan Job, numJobs)
    results := make(chan Result, numJobs)
    
    var wg sync.WaitGroup
    
    // 启动 worker
    for w := 1; w <= numWorkers; w++ {
        wg.Add(1)
        go worker(w, jobs, results, &wg)
    }
    
    // 发送任务
    for j := 1; j <= numJobs; j++ {
        jobs <- Job{ID: j, Data: fmt.Sprintf("data-%d", j)}
    }
    close(jobs)
    
    // 等待完成并关闭结果通道
    go func() {
        wg.Wait()
        close(results)
    }()
    
    // 收集结果
    for result := range results {
        fmt.Printf("Result: Job %d -> %s\n", result.JobID, result.Output)
    }
}
```

### 5.2 Rate Limiter 限流器

```go
package main

import (
    "context"
    "fmt"
    "sync"
    "time"

    "golang.org/x/time/rate"
)

func main() {
    // 简单令牌桶
    fmt.Println("=== Simple Token Bucket ===")
    simpleRateLimiter()
    
    // 使用 golang.org/x/time/rate
    fmt.Println("\n=== rate.Limiter ===")
    rateLimiterExample()
    
    // 漏桶算法
    fmt.Println("\n=== Leaky Bucket ===")
    leakyBucket()
    
    // 滑动窗口
    fmt.Println("\n=== Sliding Window ===")
    slidingWindow()
}

// 简单令牌桶
func simpleRateLimiter() {
    ticker := time.NewTicker(100 * time.Millisecond)
    defer ticker.Stop()
    
    requests := make(chan int, 5)
    
    // 生成令牌
    go func() {
        for range ticker.C {
            select {
            case requests <- 1:
            default:
            }
        }
    }()
    
    // 消费令牌
    for i := 0; i < 10; i++ {
        select {
        case <-requests:
            fmt.Printf("Request %d: allowed\n", i)
        case <-time.After(50 * time.Millisecond):
            fmt.Printf("Request %d: rate limited\n", i)
        }
    }
}

// 使用 rate.Limiter
func rateLimiterExample() {
    // 每秒5个请求，突发最多10个
    limiter := rate.NewLimiter(rate.Limit(5), 10)
    
    var wg sync.WaitGroup
    
    for i := 0; i < 15; i++ {
        wg.Add(1)
        go func(n int) {
            defer wg.Done()
            
            // 方式1: Wait - 阻塞直到允许
            ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
            defer cancel()
            
            if err := limiter.Wait(ctx); err != nil {
                fmt.Printf("Request %d: rate limited (%v)\n", n, err)
                return
            }
            fmt.Printf("Request %d: allowed\n", n)
        }(i)
        
        time.Sleep(50 * time.Millisecond)
    }
    
    wg.Wait()
    
    // 方式2: Allow - 非阻塞
    fmt.Println("\nUsing Allow():")
    limiter2 := rate.NewLimiter(rate.Limit(2), 5)
    for i := 0; i < 10; i++ {
        if limiter2.Allow() {
            fmt.Printf("Request %d: allowed\n", i)
        } else {
            fmt.Printf("Request %d: denied\n", i)
        }
    }
    
    // 方式3: Reserve - 获取预约
    fmt.Println("\nUsing Reserve():")
    limiter3 := rate.NewLimiter(rate.Limit(1), 1)
    for i := 0; i < 5; i++ {
        r := limiter3.Reserve()
        if !r.OK() {
            fmt.Println("Cannot reserve")
            continue
        }
        delay := r.Delay()
        fmt.Printf("Request %d: wait %v\n", i, delay)
        time.Sleep(delay)
        fmt.Printf("Request %d: processed\n", i)
    }
}

// 漏桶算法
type LeakyBucket struct {
    capacity   int
    remaining  int
    leakRate   time.Duration
    lastLeak   time.Time
    mu         sync.Mutex
}

func NewLeakyBucket(capacity int, leakRate time.Duration) *LeakyBucket {
    return &LeakyBucket{
        capacity:  capacity,
        remaining: capacity,
        leakRate:  leakRate,
        lastLeak:  time.Now(),
    }
}

func (lb *LeakyBucket) Allow() bool {
    lb.mu.Lock()
    defer lb.mu.Unlock()
    
    // 漏水
    now := time.Now()
    elapsed := now.Sub(lb.lastLeak)
    leaked := int(elapsed / lb.leakRate)
    
    if leaked > 0 {
        lb.remaining = min(lb.capacity, lb.remaining+leaked)
        lb.lastLeak = now
    }
    
    if lb.remaining > 0 {
        lb.remaining--
        return true
    }
    return false
}

func leakyBucket() {
    bucket := NewLeakyBucket(5, 100*time.Millisecond)
    
    for i := 0; i < 10; i++ {
        if bucket.Allow() {
            fmt.Printf("Request %d: allowed\n", i)
        } else {
            fmt.Printf("Request %d: denied\n", i)
        }
        time.Sleep(30 * time.Millisecond)
    }
}

// 滑动窗口
type SlidingWindow struct {
    windowSize time.Duration
    maxCount   int
    requests   []time.Time
    mu         sync.Mutex
}

func NewSlidingWindow(windowSize time.Duration, maxCount int) *SlidingWindow {
    return &SlidingWindow{
        windowSize: windowSize,
        maxCount:   maxCount,
        requests:   make([]time.Time, 0),
    }
}

func (sw *SlidingWindow) Allow() bool {
    sw.mu.Lock()
    defer sw.mu.Unlock()
    
    now := time.Now()
    windowStart := now.Add(-sw.windowSize)
    
    // 移除窗口外的请求
    valid := make([]time.Time, 0)
    for _, t := range sw.requests {
        if t.After(windowStart) {
            valid = append(valid, t)
        }
    }
    sw.requests = valid
    
    if len(sw.requests) < sw.maxCount {
        sw.requests = append(sw.requests, now)
        return true
    }
    return false
}

func slidingWindow() {
    sw := NewSlidingWindow(500*time.Millisecond, 3)
    
    for i := 0; i < 10; i++ {
        if sw.Allow() {
            fmt.Printf("Request %d: allowed\n", i)
        } else {
            fmt.Printf("Request %d: denied\n", i)
        }
        time.Sleep(100 * time.Millisecond)
    }
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}
```

### 5.3 Circuit Breaker 断路器

```go
package main

import (
    "errors"
    "fmt"
    "sync"
    "time"
)

type State int

const (
    StateClosed State = iota
    StateOpen
    StateHalfOpen
)

func (s State) String() string {
    switch s {
    case StateClosed:
        return "Closed"
    case StateOpen:
        return "Open"
    case StateHalfOpen:
        return "HalfOpen"
    }
    return "Unknown"
}

type CircuitBreaker struct {
    name          string
    maxFailures   int
    timeout       time.Duration
    halfOpenMax   int
    
    mu            sync.Mutex
    state         State
    failures      int
    successes     int
    lastFailure   time.Time
}

func NewCircuitBreaker(name string, maxFailures int, timeout time.Duration) *CircuitBreaker {
    return &CircuitBreaker{
        name:        name,
        maxFailures: maxFailures,
        timeout:     timeout,
        halfOpenMax: 3,
        state:       StateClosed,
    }
}

func (cb *CircuitBreaker) Execute(fn func() error) error {
    cb.mu.Lock()
    
    switch cb.state {
    case StateOpen:
        if time.Since(cb.lastFailure) > cb.timeout {
            cb.state = StateHalfOpen
            cb.successes = 0
            fmt.Printf("[%s] State: Open -> HalfOpen\n", cb.name)
        } else {
            cb.mu.Unlock()
            return errors.New("circuit breaker is open")
        }
    }
    
    cb.mu.Unlock()
    
    err := fn()
    
    cb.mu.Lock()
    defer cb.mu.Unlock()
    
    if err != nil {
        cb.failures++
        cb.lastFailure = time.Now()
        
        if cb.state == StateHalfOpen {
            cb.state = StateOpen
            fmt.Printf("[%s] State: HalfOpen -> Open\n", cb.name)
        } else if cb.failures >= cb.maxFailures {
            cb.state = StateOpen
            fmt.Printf("[%s] State: Closed -> Open\n", cb.name)
        }
        return err
    }
    
    if cb.state == StateHalfOpen {
        cb.successes++
        if cb.successes >= cb.halfOpenMax {
            cb.state = StateClosed
            cb.failures = 0
            fmt.Printf("[%s] State: HalfOpen -> Closed\n", cb.name)
        }
    } else {
        cb.failures = 0
    }
    
    return nil
}

func (cb *CircuitBreaker) State() State {
    cb.mu.Lock()
    defer cb.mu.Unlock()
    return cb.state
}

func main() {
    cb := NewCircuitBreaker("test", 3, 1*time.Second)
    
    // 模拟失败
    failingCall := func() error {
        return errors.New("service unavailable")
    }
    
    successCall := func() error {
        return nil
    }
    
    // 触发断路器打开
    fmt.Println("=== Triggering failures ===")
    for i := 0; i < 5; i++ {
        err := cb.Execute(failingCall)
        fmt.Printf("Call %d: err=%v, state=%s\n", i, err, cb.State())
    }
    
    // 等待超时
    fmt.Println("\n=== Waiting for timeout ===")
    time.Sleep(1100 * time.Millisecond)
    
    // 半开状态测试
    fmt.Println("\n=== Half-open testing ===")
    for i := 0; i < 5; i++ {
        err := cb.Execute(successCall)
        fmt.Printf("Call %d: err=%v, state=%s\n", i, err, cb.State())
    }
}
```

### 5.4 Pub/Sub 发布订阅

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

type Message struct {
    Topic   string
    Payload interface{}
}

type Subscriber struct {
    id     int
    ch     chan Message
    topics map[string]bool
}

type PubSub struct {
    mu          sync.RWMutex
    subscribers map[int]*Subscriber
    nextID      int
}

func NewPubSub() *PubSub {
    return &PubSub{
        subscribers: make(map[int]*Subscriber),
    }
}

func (ps *PubSub) Subscribe(topics ...string) *Subscriber {
    ps.mu.Lock()
    defer ps.mu.Unlock()
    
    sub := &Subscriber{
        id:     ps.nextID,
        ch:     make(chan Message, 100),
        topics: make(map[string]bool),
    }
    ps.nextID++
    
    for _, topic := range topics {
        sub.topics[topic] = true
    }
    
    ps.subscribers[sub.id] = sub
    return sub
}

func (ps *PubSub) Unsubscribe(sub *Subscriber) {
    ps.mu.Lock()
    defer ps.mu.Unlock()
    
    if _, ok := ps.subscribers[sub.id]; ok {
        close(sub.ch)
        delete(ps.subscribers, sub.id)
    }
}

func (ps *PubSub) Publish(topic string, payload interface{}) {
    ps.mu.RLock()
    defer ps.mu.RUnlock()
    
    msg := Message{Topic: topic, Payload: payload}
    
    for _, sub := range ps.subscribers {
        if sub.topics[topic] || sub.topics["*"] {
            select {
            case sub.ch <- msg:
            default:
                // 通道已满，丢弃消息
            }
        }
    }
}

func (s *Subscriber) Messages() <-chan Message {
    return s.ch
}

func main() {
    ps := NewPubSub()
    
    // 订阅者1: 订阅 "news" 和 "sports"
    sub1 := ps.Subscribe("news", "sports")
    go func() {
        for msg := range sub1.Messages() {
            fmt.Printf("Sub1 received [%s]: %v\n", msg.Topic, msg.Payload)
        }
    }()
    
    // 订阅者2: 订阅所有
    sub2 := ps.Subscribe("*")
    go func() {
        for msg := range sub2.Messages() {
            fmt.Printf("Sub2 received [%s]: %v\n", msg.Topic, msg.Payload)
        }
    }()
    
    // 订阅者3: 只订阅 "weather"
    sub3 := ps.Subscribe("weather")
    go func() {
        for msg := range sub3.Messages() {
            fmt.Printf("Sub3 received [%s]: %v\n", msg.Topic, msg.Payload)
        }
    }()
    
    // 发布消息
    time.Sleep(100 * time.Millisecond)
    
    ps.Publish("news", "Breaking news!")
    ps.Publish("sports", "Team wins!")
    ps.Publish("weather", "Sunny day")
    ps.Publish("tech", "New release")
    
    time.Sleep(100 * time.Millisecond)
    
    // 取消订阅
    ps.Unsubscribe(sub1)
    
    ps.Publish("news", "More news!")
    
    time.Sleep(100 * time.Millisecond)
}
```

### 5.5 Semaphore 信号量

```go
package main

import (
    "context"
    "fmt"
    "sync"
    "time"

    "golang.org/x/sync/semaphore"
)

func main() {
    // 使用 channel 实现
    fmt.Println("=== Channel Semaphore ===")
    channelSemaphore()
    
    // 使用 golang.org/x/sync/semaphore
    fmt.Println("\n=== x/sync Semaphore ===")
    xSyncSemaphore()
}

// Channel 实现的信号量
type Semaphore struct {
    ch chan struct{}
}

func NewSemaphore(n int) *Semaphore {
    return &Semaphore{ch: make(chan struct{}, n)}
}

func (s *Semaphore) Acquire() {
    s.ch <- struct{}{}
}

func (s *Semaphore) TryAcquire() bool {
    select {
    case s.ch <- struct{}{}:
        return true
    default:
        return false
    }
}

func (s *Semaphore) Release() {
    <-s.ch
}

func channelSemaphore() {
    sem := NewSemaphore(3)
    var wg sync.WaitGroup
    
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(n int) {
            defer wg.Done()
            
            sem.Acquire()
            defer sem.Release()
            
            fmt.Printf("Worker %d starting\n", n)
            time.Sleep(100 * time.Millisecond)
            fmt.Printf("Worker %d done\n", n)
        }(i)
    }
    
    wg.Wait()
}

func xSyncSemaphore() {
    sem := semaphore.NewWeighted(3)
    ctx := context.Background()
    var wg sync.WaitGroup
    
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(n int) {
            defer wg.Done()
            
            if err := sem.Acquire(ctx, 1); err != nil {
                fmt.Printf("Worker %d: failed to acquire: %v\n", n, err)
                return
            }
            defer sem.Release(1)
            
            fmt.Printf("Worker %d starting\n", n)
            time.Sleep(100 * time.Millisecond)
            fmt.Printf("Worker %d done\n", n)
        }(i)
    }
    
    wg.Wait()
}
```

## 6. 并发安全的数据结构

```go
package main

import (
    "fmt"
    "sync"
)

// 并发安全的队列
type SafeQueue[T any] struct {
    mu    sync.Mutex
    items []T
}

func (q *SafeQueue[T]) Enqueue(item T) {
    q.mu.Lock()
    defer q.mu.Unlock()
    q.items = append(q.items, item)
}

func (q *SafeQueue[T]) Dequeue() (T, bool) {
    q.mu.Lock()
    defer q.mu.Unlock()
    
    var zero T
    if len(q.items) == 0 {
        return zero, false
    }
    
    item := q.items[0]
    q.items = q.items[1:]
    return item, true
}

func (q *SafeQueue[T]) Len() int {
    q.mu.Lock()
    defer q.mu.Unlock()
    return len(q.items)
}

// 并发安全的栈
type SafeStack[T any] struct {
    mu    sync.Mutex
    items []T
}

func (s *SafeStack[T]) Push(item T) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.items = append(s.items, item)
}

func (s *SafeStack[T]) Pop() (T, bool) {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    var zero T
    if len(s.items) == 0 {
        return zero, false
    }
    
    item := s.items[len(s.items)-1]
    s.items = s.items[:len(s.items)-1]
    return item, true
}

// 阻塞队列
type BlockingQueue[T any] struct {
    mu       sync.Mutex
    notEmpty *sync.Cond
    notFull  *sync.Cond
    items    []T
    capacity int
}

func NewBlockingQueue[T any](capacity int) *BlockingQueue[T] {
    bq := &BlockingQueue[T]{
        items:    make([]T, 0, capacity),
        capacity: capacity,
    }
    bq.notEmpty = sync.NewCond(&bq.mu)
    bq.notFull = sync.NewCond(&bq.mu)
    return bq
}

func (bq *BlockingQueue[T]) Put(item T) {
    bq.mu.Lock()
    defer bq.mu.Unlock()
    
    for len(bq.items) >= bq.capacity {
        bq.notFull.Wait()
    }
    
    bq.items = append(bq.items, item)
    bq.notEmpty.Signal()
}

func (bq *BlockingQueue[T]) Take() T {
    bq.mu.Lock()
    defer bq.mu.Unlock()
    
    for len(bq.items) == 0 {
        bq.notEmpty.Wait()
    }
    
    item := bq.items[0]
    bq.items = bq.items[1:]
    bq.notFull.Signal()
    
    return item
}

func main() {
    // 安全队列
    queue := &SafeQueue[int]{}
    var wg sync.WaitGroup
    
    // 生产者
    for i := 0; i < 5; i++ {
        wg.Add(1)
        go func(n int) {
            defer wg.Done()
            queue.Enqueue(n)
            fmt.Printf("Enqueued: %d\n", n)
        }(i)
    }
    
    wg.Wait()
    
    // 消费者
    for queue.Len() > 0 {
        if item, ok := queue.Dequeue(); ok {
            fmt.Printf("Dequeued: %d\n", item)
        }
    }
    
    // 阻塞队列
    fmt.Println("\n=== Blocking Queue ===")
    bq := NewBlockingQueue[string](3)
    
    // 生产者
    go func() {
        for i := 0; i < 10; i++ {
            msg := fmt.Sprintf("message-%d", i)
            bq.Put(msg)
            fmt.Printf("Put: %s\n", msg)
        }
    }()
    
    // 消费者
    for i := 0; i < 10; i++ {
        msg := bq.Take()
        fmt.Printf("Take: %s\n", msg)
    }
}
```

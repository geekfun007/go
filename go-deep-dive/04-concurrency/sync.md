# 同步原语 / Synchronization Primitives

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


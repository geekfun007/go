package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// Sync 包详解
// Comprehensive Sync Package Examples

// 示例1: Mutex（互斥锁）
func mutexExample() {
	fmt.Println("\n=== 示例1: Mutex 互斥锁 ===")
	
	var (
		counter int
		mu      sync.Mutex
		wg      sync.WaitGroup
	)
	
	// 启动 10 个 goroutine 增加计数器
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				mu.Lock()
				counter++
				mu.Unlock()
			}
		}(i)
	}
	
	wg.Wait()
	fmt.Printf("最终计数: %d (期望: 10000)\n", counter)
}

// 示例2: RWMutex（读写锁）
func rwMutexExample() {
	fmt.Println("\n=== 示例2: RWMutex 读写锁 ===")
	
	var (
		data = make(map[string]int)
		mu   sync.RWMutex
		wg   sync.WaitGroup
	)
	
	// 写入操作（少量）
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			mu.Lock()
			data[fmt.Sprintf("key%d", id)] = id
			fmt.Printf("写入: key%d = %d\n", id, id)
			mu.Unlock()
			time.Sleep(100 * time.Millisecond)
		}(i)
	}
	
	// 读取操作（大量）
	time.Sleep(200 * time.Millisecond)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			mu.RLock()
			val, ok := data[fmt.Sprintf("key%d", id%5)]
			mu.RUnlock()
			if ok {
				fmt.Printf("读取: key%d = %d\n", id%5, val)
			}
		}(i)
	}
	
	wg.Wait()
}

// 示例3: WaitGroup 详解
func waitGroupExample() {
	fmt.Println("\n=== 示例3: WaitGroup 详解 ===")
	
	var wg sync.WaitGroup
	
	tasks := []string{"任务A", "任务B", "任务C", "任务D", "任务E"}
	
	for _, task := range tasks {
		wg.Add(1)
		go func(t string) {
			defer wg.Done()
			fmt.Printf("开始 %s\n", t)
			time.Sleep(time.Duration(len(t)*100) * time.Millisecond)
			fmt.Printf("完成 %s\n", t)
		}(task)
	}
	
	wg.Wait()
	fmt.Println("所有任务完成")
}

// 示例4: Once（只执行一次）
func onceExample() {
	fmt.Println("\n=== 示例4: Once 只执行一次 ===")
	
	var (
		once sync.Once
		wg   sync.WaitGroup
	)
	
	initFunc := func() {
		fmt.Println("初始化：这个函数只会执行一次")
	}
	
	// 多个 goroutine 尝试执行初始化
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Printf("Goroutine %d 尝试初始化\n", id)
			once.Do(initFunc)
			fmt.Printf("Goroutine %d 完成\n", id)
		}(i)
	}
	
	wg.Wait()
}

// 示例5: Cond（条件变量）
func condExample() {
	fmt.Println("\n=== 示例5: Cond 条件变量 ===")
	
	var (
		mu    sync.Mutex
		cond  = sync.NewCond(&mu)
		ready = false
		wg    sync.WaitGroup
	)
	
	// 等待者
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			mu.Lock()
			for !ready {
				fmt.Printf("Goroutine %d 等待信号\n", id)
				cond.Wait() // 释放锁并等待
			}
			fmt.Printf("Goroutine %d 收到信号，开始工作\n", id)
			mu.Unlock()
		}(i)
	}
	
	// 发送信号
	time.Sleep(2 * time.Second)
	mu.Lock()
	ready = true
	fmt.Println("广播信号...")
	cond.Broadcast() // 唤醒所有等待的 goroutine
	mu.Unlock()
	
	wg.Wait()
}

// 示例6: Pool（对象池）
func poolExample() {
	fmt.Println("\n=== 示例6: Pool 对象池 ===")
	
	type Buffer struct {
		data []byte
	}
	
	pool := sync.Pool{
		New: func() interface{} {
			fmt.Println("创建新的 Buffer")
			return &Buffer{data: make([]byte, 1024)}
		},
	}
	
	// 获取和归还对象
	for i := 0; i < 5; i++ {
		buffer := pool.Get().(*Buffer)
		fmt.Printf("第 %d 次获取 Buffer, 地址: %p\n", i+1, buffer)
		
		// 使用 buffer...
		
		pool.Put(buffer) // 归还到池中
	}
}

// 示例7: Map（并发安全的 map）
func syncMapExample() {
	fmt.Println("\n=== 示例7: sync.Map 并发安全 Map ===")
	
	var (
		m  sync.Map
		wg sync.WaitGroup
	)
	
	// 并发写入
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			m.Store(fmt.Sprintf("key%d", id), id*10)
			fmt.Printf("存储: key%d = %d\n", id, id*10)
		}(i)
	}
	
	wg.Wait()
	time.Sleep(100 * time.Millisecond)
	
	// 并发读取
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			if val, ok := m.Load(fmt.Sprintf("key%d", id)); ok {
				fmt.Printf("读取: key%d = %v\n", id, val)
			}
		}(i)
	}
	
	wg.Wait()
	
	// 遍历
	fmt.Println("遍历所有键值对:")
	m.Range(func(key, value interface{}) bool {
		fmt.Printf("  %v: %v\n", key, value)
		return true
	})
}

// 示例8: Atomic 原子操作
func atomicExample() {
	fmt.Println("\n=== 示例8: Atomic 原子操作 ===")
	
	var (
		counter int64
		wg      sync.WaitGroup
	)
	
	// 使用原子操作增加计数器
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				atomic.AddInt64(&counter, 1)
			}
		}()
	}
	
	wg.Wait()
	fmt.Printf("最终计数: %d (期望: 10000)\n", atomic.LoadInt64(&counter))
	
	// 原子比较并交换
	var value int64 = 100
	swapped := atomic.CompareAndSwapInt64(&value, 100, 200)
	fmt.Printf("CAS 操作: %v, 新值: %d\n", swapped, value)
}

// 示例9: 死锁示例（错误示例）
func deadlockExample() {
	fmt.Println("\n=== 示例9: 死锁示例（注意避免）===")
	
	var (
		mu1 sync.Mutex
		mu2 sync.Mutex
	)
	
	// 这是一个死锁场景的演示，实际不会执行
	fmt.Println("死锁场景:")
	fmt.Println("  Goroutine 1: 锁 A -> 锁 B")
	fmt.Println("  Goroutine 2: 锁 B -> 锁 A")
	fmt.Println("避免方法: 始终以相同顺序获取锁")
	
	// 正确的做法
	done := make(chan bool)
	
	go func() {
		mu1.Lock()
		time.Sleep(10 * time.Millisecond)
		mu2.Lock()
		fmt.Println("Goroutine 1 获得两个锁")
		mu2.Unlock()
		mu1.Unlock()
		done <- true
	}()
	
	go func() {
		time.Sleep(5 * time.Millisecond)
		mu1.Lock() // 相同顺序
		mu2.Lock()
		fmt.Println("Goroutine 2 获得两个锁")
		mu2.Unlock()
		mu1.Unlock()
		done <- true
	}()
	
	<-done
	<-done
}

func main() {
	fmt.Println("Go 异步编程 - Sync 包详解")
	fmt.Println("===========================")
	
	mutexExample()
	rwMutexExample()
	waitGroupExample()
	onceExample()
	condExample()
	poolExample()
	syncMapExample()
	atomicExample()
	deadlockExample()
	
	fmt.Println("\n所有示例完成！")
}

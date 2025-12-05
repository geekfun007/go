package main

import (
	"fmt"
	"sync"
	"time"
)

// 高级并发模式
// Advanced Concurrency Patterns

// 示例1: Worker Pool（工作池）
func workerPoolPattern() {
	fmt.Println("\n=== 示例1: Worker Pool 工作池 ===")
	
	const numWorkers = 3
	const numJobs = 10
	
	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)
	
	// 启动 workers
	var wg sync.WaitGroup
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}
	
	// 发送任务
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs)
	
	// 等待所有 worker 完成
	go func() {
		wg.Wait()
		close(results)
	}()
	
	// 收集结果
	for result := range results {
		fmt.Printf("结果: %d\n", result)
	}
}

func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		fmt.Printf("Worker %d 处理任务 %d\n", id, job)
		time.Sleep(500 * time.Millisecond)
		results <- job * 2
	}
}

// 示例2: Pipeline（流水线）
func pipelinePattern() {
	fmt.Println("\n=== 示例2: Pipeline 流水线 ===")
	
	// 阶段1: 生成数字
	generator := func() <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for i := 1; i <= 5; i++ {
				out <- i
			}
		}()
		return out
	}
	
	// 阶段2: 平方
	square := func(in <-chan int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for n := range in {
				out <- n * n
			}
		}()
		return out
	}
	
	// 阶段3: 加倍
	double := func(in <-chan int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for n := range in {
				out <- n * 2
			}
		}()
		return out
	}
	
	// 构建流水线
	numbers := generator()
	squared := square(numbers)
	doubled := double(squared)
	
	// 消费结果
	for result := range doubled {
		fmt.Printf("结果: %d\n", result)
	}
}

// 示例3: Fan-out/Fan-in（扇出/扇入）
func fanOutFanInPattern() {
	fmt.Println("\n=== 示例3: Fan-out/Fan-in 扇出扇入 ===")
	
	// 生成输入
	input := make(chan int, 10)
	go func() {
		for i := 1; i <= 10; i++ {
			input <- i
		}
		close(input)
	}()
	
	// Fan-out: 启动多个 worker
	const numWorkers = 3
	workers := make([]<-chan int, numWorkers)
	for i := 0; i < numWorkers; i++ {
		workers[i] = expensiveOperation(input)
	}
	
	// Fan-in: 合并所有 worker 的输出
	results := merge(workers...)
	
	// 处理结果
	for result := range results {
		fmt.Printf("收到结果: %d\n", result)
	}
}

func expensiveOperation(input <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range input {
			time.Sleep(100 * time.Millisecond) // 模拟耗时操作
			out <- n * n
		}
	}()
	return out
}

func merge(channels ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)
	
	// 为每个输入 channel 启动一个 goroutine
	multiplex := func(c <-chan int) {
		defer wg.Done()
		for n := range c {
			out <- n
		}
	}
	
	wg.Add(len(channels))
	for _, c := range channels {
		go multiplex(c)
	}
	
	// 等待所有 goroutine 完成后关闭输出 channel
	go func() {
		wg.Wait()
		close(out)
	}()
	
	return out
}

// 示例4: Rate Limiting（速率限制）
func rateLimitingPattern() {
	fmt.Println("\n=== 示例4: Rate Limiting 速率限制 ===")
	
	requests := make(chan int, 10)
	for i := 1; i <= 10; i++ {
		requests <- i
	}
	close(requests)
	
	// 创建一个速率限制器：每 500ms 一个令牌
	limiter := time.Tick(500 * time.Millisecond)
	
	for req := range requests {
		<-limiter // 等待令牌
		fmt.Printf("处理请求 %d at %s\n", req, time.Now().Format("15:04:05.000"))
	}
}

// 示例5: Burst Rate Limiting（突发速率限制）
func burstRateLimitingPattern() {
	fmt.Println("\n=== 示例5: Burst Rate Limiting 突发速率限制 ===")
	
	requests := make(chan int, 10)
	for i := 1; i <= 10; i++ {
		requests <- i
	}
	close(requests)
	
	// 允许最多 3 个突发请求
	burstyLimiter := make(chan time.Time, 3)
	for i := 0; i < 3; i++ {
		burstyLimiter <- time.Now()
	}
	
	// 每 500ms 填充一个令牌
	go func() {
		for t := range time.Tick(500 * time.Millisecond) {
			burstyLimiter <- t
		}
	}()
	
	for req := range requests {
		<-burstyLimiter
		fmt.Printf("处理请求 %d at %s\n", req, time.Now().Format("15:04:05.000"))
	}
}

// 示例6: Semaphore（信号量）模式
func semaphorePattern() {
	fmt.Println("\n=== 示例6: Semaphore 信号量 ===")
	
	const maxConcurrent = 3
	semaphore := make(chan struct{}, maxConcurrent)
	
	var wg sync.WaitGroup
	
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			
			semaphore <- struct{}{} // 获取信号量
			fmt.Printf("任务 %d 开始执行\n", id)
			time.Sleep(1 * time.Second)
			fmt.Printf("任务 %d 完成\n", id)
			<-semaphore // 释放信号量
		}(i)
	}
	
	wg.Wait()
	fmt.Println("所有任务完成")
}

// 示例7: Error Group 模式
func errorGroupPattern() {
	fmt.Println("\n=== 示例7: Error Group 错误组 ===")
	
	type ErrorGroup struct {
		wg     sync.WaitGroup
		errMux sync.Mutex
		err    error
	}
	
	eg := &ErrorGroup{}
	
	tasks := []func() error{
		func() error {
			time.Sleep(100 * time.Millisecond)
			fmt.Println("任务 1 完成")
			return nil
		},
		func() error {
			time.Sleep(200 * time.Millisecond)
			fmt.Println("任务 2 失败")
			return fmt.Errorf("任务 2 错误")
		},
		func() error {
			time.Sleep(300 * time.Millisecond)
			fmt.Println("任务 3 完成")
			return nil
		},
	}
	
	for _, task := range tasks {
		task := task // 捕获循环变量
		eg.wg.Add(1)
		go func() {
			defer eg.wg.Done()
			if err := task(); err != nil {
				eg.errMux.Lock()
				if eg.err == nil {
					eg.err = err
				}
				eg.errMux.Unlock()
			}
		}()
	}
	
	eg.wg.Wait()
	
	if eg.err != nil {
		fmt.Printf("发生错误: %v\n", eg.err)
	} else {
		fmt.Println("所有任务成功")
	}
}

// 示例8: Future/Promise 模式
func futurePattern() {
	fmt.Println("\n=== 示例8: Future/Promise 模式 ===")
	
	type Future struct {
		result chan interface{}
		err    chan error
	}
	
	// 创建一个 Future
	asyncComputation := func(input int) *Future {
		future := &Future{
			result: make(chan interface{}, 1),
			err:    make(chan error, 1),
		}
		
		go func() {
			time.Sleep(1 * time.Second)
			if input < 0 {
				future.err <- fmt.Errorf("输入不能为负数")
			} else {
				future.result <- input * input
			}
		}()
		
		return future
	}
	
	// 使用 Future
	future := asyncComputation(5)
	fmt.Println("计算中...")
	
	select {
	case result := <-future.result:
		fmt.Printf("结果: %v\n", result)
	case err := <-future.err:
		fmt.Printf("错误: %v\n", err)
	}
}

// 示例9: Broadcast 模式
func broadcastPattern() {
	fmt.Println("\n=== 示例9: Broadcast 广播模式 ===")
	
	type Broadcaster struct {
		mu        sync.RWMutex
		listeners []chan string
	}
	
	b := &Broadcaster{}
	
	// 订阅
	subscribe := func(name string) chan string {
		ch := make(chan string, 1)
		b.mu.Lock()
		b.listeners = append(b.listeners, ch)
		b.mu.Unlock()
		
		go func() {
			for msg := range ch {
				fmt.Printf("%s 收到: %s\n", name, msg)
			}
		}()
		
		return ch
	}
	
	// 广播
	broadcast := func(msg string) {
		b.mu.RLock()
		defer b.mu.RUnlock()
		for _, ch := range b.listeners {
			ch <- msg
		}
	}
	
	// 订阅者
	subscribe("订阅者A")
	subscribe("订阅者B")
	subscribe("订阅者C")
	
	time.Sleep(100 * time.Millisecond)
	
	// 广播消息
	broadcast("你好，世界！")
	broadcast("第二条消息")
	
	time.Sleep(500 * time.Millisecond)
}

func main() {
	fmt.Println("Go 异步编程 - 高级并发模式")
	fmt.Println("============================")
	
	workerPoolPattern()
	pipelinePattern()
	fanOutFanInPattern()
	rateLimitingPattern()
	burstRateLimitingPattern()
	semaphorePattern()
	errorGroupPattern()
	futurePattern()
	broadcastPattern()
	
	fmt.Println("\n所有示例完成！")
}

package main

import (
	"fmt"
	"sync"
	"time"
)

// 基础 Goroutine 示例
// Basic Goroutine Examples

// 示例1: 简单的 goroutine
func simpleGoroutine() {
	fmt.Println("\n=== 示例1: 简单的 Goroutine ===")
	
	// 启动一个 goroutine
	go func() {
		fmt.Println("这是一个 goroutine")
	}()
	
	// 等待 goroutine 完成（不推荐使用 sleep，这里仅作演示）
	time.Sleep(100 * time.Millisecond)
}

// 示例2: 多个 goroutine
func multipleGoroutines() {
	fmt.Println("\n=== 示例2: 多个 Goroutine ===")
	
	for i := 0; i < 5; i++ {
		go func(id int) {
			fmt.Printf("Goroutine %d 正在运行\n", id)
		}(i) // 传递参数避免闭包陷阱
	}
	
	time.Sleep(100 * time.Millisecond)
}

// 示例3: 使用 WaitGroup 正确等待 goroutine 完成
func goroutinesWithWaitGroup() {
	fmt.Println("\n=== 示例3: 使用 WaitGroup ===")
	
	var wg sync.WaitGroup
	
	for i := 0; i < 5; i++ {
		wg.Add(1) // 增加等待计数
		go func(id int) {
			defer wg.Done() // goroutine 完成时减少计数
			fmt.Printf("Worker %d 开始工作\n", id)
			time.Sleep(time.Duration(id*100) * time.Millisecond)
			fmt.Printf("Worker %d 完成工作\n", id)
		}(i)
	}
	
	wg.Wait() // 等待所有 goroutine 完成
	fmt.Println("所有 workers 完成")
}

// 示例4: Goroutine 泄漏问题
func goroutineLeakExample() {
	fmt.Println("\n=== 示例4: Goroutine 泄漏（错误示例）===")
	
	// 不好的做法：goroutine 永远阻塞
	ch := make(chan int)
	go func() {
		val := <-ch // 如果没有发送者，这里会永远阻塞
		fmt.Println(val)
	}()
	
	// 正确的做法：使用超时或关闭通道
	fmt.Println("注意：上面的 goroutine 会泄漏（仅作演示）")
}

// 示例5: 正确处理 goroutine 生命周期
func properGoroutineLifecycle() {
	fmt.Println("\n=== 示例5: 正确的 Goroutine 生命周期 ===")
	
	done := make(chan bool)
	
	go func() {
		defer func() {
			done <- true
		}()
		
		fmt.Println("执行一些工作...")
		time.Sleep(100 * time.Millisecond)
		fmt.Println("工作完成")
	}()
	
	<-done // 等待 goroutine 完成
	fmt.Println("主程序继续")
}

// 示例6: 匿名函数 vs 命名函数
func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("命名函数 worker %d 执行\n", id)
	time.Sleep(50 * time.Millisecond)
}

func namedVsAnonymous() {
	fmt.Println("\n=== 示例6: 匿名函数 vs 命名函数 ===")
	
	var wg sync.WaitGroup
	
	// 使用命名函数
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go worker(i, &wg)
	}
	
	// 使用匿名函数
	for i := 3; i < 6; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Printf("匿名函数 worker %d 执行\n", id)
			time.Sleep(50 * time.Millisecond)
		}(i)
	}
	
	wg.Wait()
}

func main() {
	fmt.Println("Go 异步编程 - 基础 Goroutine 示例")
	fmt.Println("=====================================")
	
	simpleGoroutine()
	multipleGoroutines()
	goroutinesWithWaitGroup()
	goroutineLeakExample()
	properGoroutineLifecycle()
	namedVsAnonymous()
	
	fmt.Println("\n所有示例完成！")
}

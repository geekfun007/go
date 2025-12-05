package main

import (
	"fmt"
	"time"
)

// Channel 详解与实战
// Comprehensive Channel Examples

// 示例1: 无缓冲 channel (同步)
func unbufferedChannel() {
	fmt.Println("\n=== 示例1: 无缓冲 Channel ===")
	
	ch := make(chan string)
	
	// 发送者
	go func() {
		fmt.Println("发送者：准备发送数据...")
		ch <- "Hello, Channel!"
		fmt.Println("发送者：数据已发送")
	}()
	
	// 接收者
	time.Sleep(1 * time.Second) // 模拟接收者延迟
	fmt.Println("接收者：准备接收数据...")
	msg := <-ch
	fmt.Printf("接收者：收到消息: %s\n", msg)
}

// 示例2: 有缓冲 channel (异步)
func bufferedChannel() {
	fmt.Println("\n=== 示例2: 有缓冲 Channel ===")
	
	ch := make(chan int, 3) // 缓冲区大小为 3
	
	// 发送数据（不会阻塞，直到缓冲区满）
	ch <- 1
	ch <- 2
	ch <- 3
	fmt.Println("发送了 3 个值到缓冲 channel")
	
	// 接收数据
	fmt.Println("接收:", <-ch)
	fmt.Println("接收:", <-ch)
	fmt.Println("接收:", <-ch)
}

// 示例3: Channel 方向（单向 channel）
func channelDirection() {
	fmt.Println("\n=== 示例3: Channel 方向 ===")
	
	ch := make(chan int)
	
	// 只能发送的 channel
	go sender(ch)
	
	// 只能接收的 channel
	receiver(ch)
}

func sender(ch chan<- int) { // 只发送
	for i := 0; i < 5; i++ {
		ch <- i
		fmt.Printf("发送: %d\n", i)
	}
	close(ch)
}

func receiver(ch <-chan int) { // 只接收
	for val := range ch {
		fmt.Printf("接收: %d\n", val)
	}
}

// 示例4: 关闭 channel
func closingChannel() {
	fmt.Println("\n=== 示例4: 关闭 Channel ===")
	
	ch := make(chan int, 5)
	
	// 发送数据
	go func() {
		for i := 0; i < 5; i++ {
			ch <- i
		}
		close(ch) // 关闭 channel
	}()
	
	// 接收数据直到 channel 关闭
	for val := range ch {
		fmt.Printf("接收: %d\n", val)
	}
	
	// 检查 channel 是否关闭
	val, ok := <-ch
	if !ok {
		fmt.Println("Channel 已关闭")
	} else {
		fmt.Printf("收到值: %d\n", val)
	}
}

// 示例5: Channel 作为信号量
func channelAsSignal() {
	fmt.Println("\n=== 示例5: Channel 作为信号量 ===")
	
	done := make(chan bool)
	
	go func() {
		fmt.Println("执行耗时操作...")
		time.Sleep(2 * time.Second)
		fmt.Println("操作完成")
		done <- true
	}()
	
	fmt.Println("等待操作完成...")
	<-done
	fmt.Println("继续执行主程序")
}

// 示例6: 多个 goroutine 通信
func multipleGoroutineCommunication() {
	fmt.Println("\n=== 示例6: 多个 Goroutine 通信 ===")
	
	numbers := make(chan int)
	squares := make(chan int)
	
	// Goroutine 1: 生成数字
	go func() {
		for i := 1; i <= 5; i++ {
			numbers <- i
		}
		close(numbers)
	}()
	
	// Goroutine 2: 计算平方
	go func() {
		for num := range numbers {
			squares <- num * num
		}
		close(squares)
	}()
	
	// 主 goroutine: 接收结果
	for square := range squares {
		fmt.Printf("平方值: %d\n", square)
	}
}

// 示例7: Channel 超时处理
func channelTimeout() {
	fmt.Println("\n=== 示例7: Channel 超时处理 ===")
	
	ch := make(chan string)
	
	go func() {
		time.Sleep(3 * time.Second)
		ch <- "结果"
	}()
	
	select {
	case result := <-ch:
		fmt.Printf("收到结果: %s\n", result)
	case <-time.After(2 * time.Second):
		fmt.Println("超时：操作耗时过长")
	}
}

// 示例8: Nil channel 的行为
func nilChannelBehavior() {
	fmt.Println("\n=== 示例8: Nil Channel 行为 ===")
	
	var ch chan int // nil channel
	
	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("Goroutine 完成")
	}()
	
	select {
	case <-ch:
		fmt.Println("从 nil channel 接收（永远不会发生）")
	case <-time.After(500 * time.Millisecond):
		fmt.Println("Nil channel 的 select 会阻塞")
	}
}

func main() {
	fmt.Println("Go 异步编程 - Channel 详解")
	fmt.Println("===========================")
	
	unbufferedChannel()
	bufferedChannel()
	channelDirection()
	closingChannel()
	channelAsSignal()
	multipleGoroutineCommunication()
	channelTimeout()
	nilChannelBehavior()
	
	fmt.Println("\n所有示例完成！")
}

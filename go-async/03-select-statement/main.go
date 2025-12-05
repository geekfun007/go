package main

import (
	"fmt"
	"time"
)

// Select 语句详解
// Comprehensive Select Statement Examples

// 示例1: 基本 select 用法
func basicSelect() {
	fmt.Println("\n=== 示例1: 基本 Select ===")
	
	ch1 := make(chan string)
	ch2 := make(chan string)
	
	go func() {
		time.Sleep(1 * time.Second)
		ch1 <- "来自 channel 1"
	}()
	
	go func() {
		time.Sleep(2 * time.Second)
		ch2 <- "来自 channel 2"
	}()
	
	// select 会等待第一个就绪的 channel
	select {
	case msg1 := <-ch1:
		fmt.Println("收到:", msg1)
	case msg2 := <-ch2:
		fmt.Println("收到:", msg2)
	}
}

// 示例2: 多个 case 都就绪时的随机选择
func randomSelection() {
	fmt.Println("\n=== 示例2: 随机选择 ===")
	
	ch1 := make(chan string, 1)
	ch2 := make(chan string, 1)
	
	// 两个 channel 都立即就绪
	ch1 <- "channel 1"
	ch2 <- "channel 2"
	
	// 运行多次，观察随机性
	for i := 0; i < 5; i++ {
		select {
		case msg := <-ch1:
			fmt.Printf("第 %d 次: %s\n", i+1, msg)
			ch1 <- msg // 放回去
		case msg := <-ch2:
			fmt.Printf("第 %d 次: %s\n", i+1, msg)
			ch2 <- msg // 放回去
		}
	}
}

// 示例3: Default case (非阻塞操作)
func nonBlockingSelect() {
	fmt.Println("\n=== 示例3: 非阻塞 Select ===")
	
	ch := make(chan string)
	
	// 非阻塞发送
	select {
	case ch <- "消息":
		fmt.Println("发送成功")
	default:
		fmt.Println("无法发送，channel 未就绪")
	}
	
	// 非阻塞接收
	select {
	case msg := <-ch:
		fmt.Println("收到:", msg)
	default:
		fmt.Println("无数据可接收")
	}
}

// 示例4: 超时模式
func timeoutPattern() {
	fmt.Println("\n=== 示例4: 超时模式 ===")
	
	ch := make(chan string)
	
	go func() {
		time.Sleep(3 * time.Second)
		ch <- "延迟的结果"
	}()
	
	select {
	case result := <-ch:
		fmt.Println("收到结果:", result)
	case <-time.After(2 * time.Second):
		fmt.Println("操作超时")
	}
}

// 示例5: 定期执行（Ticker）
func tickerPattern() {
	fmt.Println("\n=== 示例5: 定期执行 ===")
	
	ticker := time.NewTicker(500 * time.Millisecond)
	done := make(chan bool)
	
	go func() {
		time.Sleep(3 * time.Second)
		done <- true
	}()
	
	count := 0
	for {
		select {
		case <-ticker.C:
			count++
			fmt.Printf("Tick %d: %s\n", count, time.Now().Format("15:04:05"))
		case <-done:
			ticker.Stop()
			fmt.Println("定时器停止")
			return
		}
	}
}

// 示例6: 多路复用（Fan-in）
func fanInPattern() {
	fmt.Println("\n=== 示例6: 多路复用 (Fan-in) ===")
	
	ch1 := make(chan string)
	ch2 := make(chan string)
	
	// 生产者 1
	go func() {
		for i := 0; i < 3; i++ {
			ch1 <- fmt.Sprintf("Source1: %d", i)
			time.Sleep(500 * time.Millisecond)
		}
		close(ch1)
	}()
	
	// 生产者 2
	go func() {
		for i := 0; i < 3; i++ {
			ch2 <- fmt.Sprintf("Source2: %d", i)
			time.Sleep(700 * time.Millisecond)
		}
		close(ch2)
	}()
	
	// 多路复用：合并两个 channel
	for ch1 != nil || ch2 != nil {
		select {
		case msg, ok := <-ch1:
			if ok {
				fmt.Println("收到:", msg)
			} else {
				ch1 = nil
			}
		case msg, ok := <-ch2:
			if ok {
				fmt.Println("收到:", msg)
			} else {
				ch2 = nil
			}
		}
	}
}

// 示例7: 退出信号模式
func quitSignalPattern() {
	fmt.Println("\n=== 示例7: 退出信号模式 ===")
	
	quit := make(chan bool)
	data := make(chan int)
	
	// 工作协程
	go func() {
		for {
			select {
			case val := <-data:
				fmt.Printf("处理数据: %d\n", val)
			case <-quit:
				fmt.Println("收到退出信号，清理资源...")
				quit <- true // 确认退出
				return
			}
		}
	}()
	
	// 发送一些数据
	for i := 0; i < 5; i++ {
		data <- i
		time.Sleep(200 * time.Millisecond)
	}
	
	// 发送退出信号
	quit <- true
	<-quit // 等待确认
	fmt.Println("工作协程已退出")
}

// 示例8: 优先级 select（模拟）
func prioritySelect() {
	fmt.Println("\n=== 示例8: 优先级 Select ===")
	
	highPriority := make(chan string, 1)
	lowPriority := make(chan string, 1)
	
	// 填充数据
	highPriority <- "高优先级消息"
	lowPriority <- "低优先级消息"
	
	// 优先检查高优先级
	for i := 0; i < 2; i++ {
		select {
		case msg := <-highPriority:
			fmt.Println("处理:", msg)
		default:
			select {
			case msg := <-lowPriority:
				fmt.Println("处理:", msg)
			default:
				fmt.Println("无消息")
			}
		}
	}
}

// 示例9: 循环中的 select
func selectInLoop() {
	fmt.Println("\n=== 示例9: 循环中的 Select ===")
	
	ch := make(chan int)
	done := make(chan bool)
	
	// 生产者
	go func() {
		for i := 0; i < 5; i++ {
			ch <- i
			time.Sleep(300 * time.Millisecond)
		}
		done <- true
	}()
	
	// 消费者
	count := 0
	for {
		select {
		case val := <-ch:
			fmt.Printf("收到值: %d\n", val)
			count++
		case <-done:
			fmt.Printf("完成，共处理 %d 个值\n", count)
			return
		case <-time.After(1 * time.Second):
			fmt.Println("超时，退出")
			return
		}
	}
}

func main() {
	fmt.Println("Go 异步编程 - Select 语句详解")
	fmt.Println("===============================")
	
	basicSelect()
	randomSelection()
	nonBlockingSelect()
	timeoutPattern()
	tickerPattern()
	fanInPattern()
	quitSignalPattern()
	prioritySelect()
	selectInLoop()
	
	fmt.Println("\n所有示例完成！")
}

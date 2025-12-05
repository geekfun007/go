package main

import (
	"context"
	"fmt"
	"time"
)

// Context 详解与实战
// Comprehensive Context Examples

// 示例1: 基本的 Context 使用
func basicContext() {
	fmt.Println("\n=== 示例1: 基本 Context ===")
	
	// 创建一个可取消的 context
	ctx, cancel := context.WithCancel(context.Background())
	
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Goroutine 收到取消信号")
				return
			default:
				fmt.Println("工作中...")
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()
	
	time.Sleep(2 * time.Second)
	fmt.Println("取消 context")
	cancel()
	time.Sleep(500 * time.Millisecond)
}

// 示例2: Context 超时
func contextWithTimeout() {
	fmt.Println("\n=== 示例2: Context 超时 ===")
	
	// 创建一个 3 秒后超时的 context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	go func() {
		select {
		case <-time.After(5 * time.Second):
			fmt.Println("操作完成")
		case <-ctx.Done():
			fmt.Println("操作超时:", ctx.Err())
		}
	}()
	
	time.Sleep(4 * time.Second)
}

// 示例3: Context 截止时间
func contextWithDeadline() {
	fmt.Println("\n=== 示例3: Context 截止时间 ===")
	
	deadline := time.Now().Add(2 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()
	
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("截止时间到达:", ctx.Err())
				return
			default:
				fmt.Println("执行任务...")
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()
	
	time.Sleep(3 * time.Second)
}

// 示例4: Context 传递值
func contextWithValue() {
	fmt.Println("\n=== 示例4: Context 传递值 ===")
	
	type key string
	
	ctx := context.Background()
	ctx = context.WithValue(ctx, key("userID"), 12345)
	ctx = context.WithValue(ctx, key("requestID"), "abc-123")
	
	processRequest(ctx)
}

func processRequest(ctx context.Context) {
	userID := ctx.Value("userID")
	requestID := ctx.Value("requestID")
	
	fmt.Printf("处理请求 - UserID: %v, RequestID: %v\n", userID, requestID)
}

// 示例5: Context 链式取消
func contextCancellationChain() {
	fmt.Println("\n=== 示例5: Context 链式取消 ===")
	
	// 创建父 context
	parentCtx, parentCancel := context.WithCancel(context.Background())
	defer parentCancel()
	
	// 创建子 context
	childCtx, childCancel := context.WithCancel(parentCtx)
	defer childCancel()
	
	go func() {
		<-childCtx.Done()
		fmt.Println("子 context 被取消")
	}()
	
	time.Sleep(1 * time.Second)
	fmt.Println("取消父 context")
	parentCancel() // 取消父 context 会自动取消子 context
	time.Sleep(500 * time.Millisecond)
}

// 示例6: 模拟 HTTP 请求超时
func simulateHTTPRequestWithTimeout() {
	fmt.Println("\n=== 示例6: HTTP 请求超时模拟 ===")
	
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	
	result := make(chan string)
	
	go func() {
		// 模拟耗时的 HTTP 请求
		time.Sleep(3 * time.Second)
		result <- "HTTP 响应数据"
	}()
	
	select {
	case res := <-result:
		fmt.Println("收到响应:", res)
	case <-ctx.Done():
		fmt.Println("请求超时:", ctx.Err())
	}
}

// 示例7: 数据库查询超时
func simulateDBQueryWithTimeout() {
	fmt.Println("\n=== 示例7: 数据库查询超时模拟 ===")
	
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	
	result, err := queryDatabase(ctx)
	if err != nil {
		fmt.Println("查询失败:", err)
	} else {
		fmt.Println("查询成功:", result)
	}
}

func queryDatabase(ctx context.Context) (string, error) {
	result := make(chan string)
	
	go func() {
		time.Sleep(500 * time.Millisecond)
		result <- "数据库记录"
	}()
	
	select {
	case data := <-result:
		return data, nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

// 示例8: 多个 goroutine 协调取消
func coordinatedCancellation() {
	fmt.Println("\n=== 示例8: 多个 Goroutine 协调取消 ===")
	
	ctx, cancel := context.WithCancel(context.Background())
	
	// 启动多个 worker
	for i := 1; i <= 3; i++ {
		go worker(ctx, i)
	}
	
	time.Sleep(2 * time.Second)
	fmt.Println("发送取消信号...")
	cancel()
	time.Sleep(1 * time.Second)
}

func worker(ctx context.Context, id int) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Worker %d: 收到取消信号，退出\n", id)
			return
		default:
			fmt.Printf("Worker %d: 工作中\n", id)
			time.Sleep(500 * time.Millisecond)
		}
	}
}

// 示例9: Context 最佳实践
func contextBestPractices() {
	fmt.Println("\n=== 示例9: Context 最佳实践 ===")
	
	// 1. 总是检查 context.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	
	// 2. 传递 context 作为第一个参数
	performOperation(ctx, "重要操作")
}

func performOperation(ctx context.Context, operationName string) {
	fmt.Printf("开始执行: %s\n", operationName)
	
	for i := 0; i < 5; i++ {
		// 检查是否应该取消
		select {
		case <-ctx.Done():
			fmt.Printf("%s 被取消: %v\n", operationName, ctx.Err())
			return
		default:
			fmt.Printf("%s: 步骤 %d\n", operationName, i+1)
			time.Sleep(500 * time.Millisecond)
		}
	}
	
	fmt.Printf("%s 完成\n", operationName)
}

func main() {
	fmt.Println("Go 异步编程 - Context 详解")
	fmt.Println("===========================")
	
	basicContext()
	contextWithTimeout()
	contextWithDeadline()
	contextWithValue()
	contextCancellationChain()
	simulateHTTPRequestWithTimeout()
	simulateDBQueryWithTimeout()
	coordinatedCancellation()
	contextBestPractices()
	
	fmt.Println("\n所有示例完成！")
}

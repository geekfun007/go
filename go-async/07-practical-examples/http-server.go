package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// 实战示例: 并发 HTTP 服务器
// Practical Example: Concurrent HTTP Server

// 示例1: 基本的并发 HTTP 处理器
type Server struct {
	mu       sync.RWMutex
	visitors map[string]int
}

func NewServer() *Server {
	return &Server{
		visitors: make(map[string]int),
	}
}

func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	s.visitors[r.RemoteAddr]++
	count := s.visitors[r.RemoteAddr]
	s.mu.Unlock()
	
	fmt.Fprintf(w, "欢迎！你已访问 %d 次\n", count)
}

// 示例2: 超时处理
func (s *Server) handleSlowOperation(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	
	result := make(chan string, 1)
	
	go func() {
		// 模拟耗时操作
		time.Sleep(3 * time.Second)
		result <- "操作完成"
	}()
	
	select {
	case res := <-result:
		fmt.Fprintf(w, "结果: %s\n", res)
	case <-ctx.Done():
		http.Error(w, "请求超时", http.StatusRequestTimeout)
	}
}

// 示例3: 并发请求处理
func (s *Server) handleConcurrentRequests(w http.ResponseWriter, r *http.Request) {
	type Result struct {
		Service string
		Data    string
		Err     error
	}
	
	results := make(chan Result, 3)
	
	// 同时调用多个服务
	services := []string{"服务A", "服务B", "服务C"}
	for _, service := range services {
		service := service
		go func() {
			time.Sleep(time.Duration(len(service)*100) * time.Millisecond)
			results <- Result{
				Service: service,
				Data:    fmt.Sprintf("%s 的数据", service),
			}
		}()
	}
	
	// 收集所有结果
	fmt.Fprintf(w, "并发请求结果:\n")
	for i := 0; i < len(services); i++ {
		result := <-results
		fmt.Fprintf(w, "- %s: %s\n", result.Service, result.Data)
	}
}

// 示例4: 请求速率限制
type RateLimiter struct {
	tokens chan struct{}
}

func NewRateLimiter(rate int) *RateLimiter {
	rl := &RateLimiter{
		tokens: make(chan struct{}, rate),
	}
	
	// 填充令牌桶
	for i := 0; i < rate; i++ {
		rl.tokens <- struct{}{}
	}
	
	// 定期补充令牌
	go func() {
		ticker := time.NewTicker(time.Second / time.Duration(rate))
		for range ticker.C {
			select {
			case rl.tokens <- struct{}{}:
			default:
			}
		}
	}()
	
	return rl
}

func (rl *RateLimiter) Allow() bool {
	select {
	case <-rl.tokens:
		return true
	default:
		return false
	}
}

func rateLimitMiddleware(rl *RateLimiter, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !rl.Allow() {
			http.Error(w, "请求过于频繁", http.StatusTooManyRequests)
			return
		}
		next(w, r)
	}
}

// 示例5: 工作池处理请求
type WorkerPool struct {
	workers int
	jobs    chan func()
	wg      sync.WaitGroup
}

func NewWorkerPool(workers int) *WorkerPool {
	wp := &WorkerPool{
		workers: workers,
		jobs:    make(chan func(), 100),
	}
	wp.start()
	return wp
}

func (wp *WorkerPool) start() {
	for i := 0; i < wp.workers; i++ {
		wp.wg.Add(1)
		go func(id int) {
			defer wp.wg.Done()
			for job := range wp.jobs {
				job()
			}
		}(i)
	}
}

func (wp *WorkerPool) Submit(job func()) {
	wp.jobs <- job
}

func (wp *WorkerPool) Stop() {
	close(wp.jobs)
	wp.wg.Wait()
}

func (s *Server) handleWithWorkerPool(pool *WorkerPool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		done := make(chan bool)
		
		pool.Submit(func() {
			defer func() { done <- true }()
			time.Sleep(500 * time.Millisecond)
			fmt.Fprintf(w, "请求已由工作池处理\n")
		})
		
		<-done
	}
}

// 示例6: 优雅关闭
func gracefulShutdown(srv *http.Server) {
	quit := make(chan struct{})
	
	go func() {
		// 模拟关闭信号
		time.Sleep(10 * time.Second)
		close(quit)
	}()
	
	<-quit
	log.Println("服务器正在关闭...")
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("服务器关闭错误: %v", err)
	}
	
	log.Println("服务器已优雅关闭")
}

// 示例7: 连接池示例
type ConnectionPool struct {
	mu    sync.Mutex
	conns chan *Connection
}

type Connection struct {
	ID int
}

func NewConnectionPool(size int) *ConnectionPool {
	pool := &ConnectionPool{
		conns: make(chan *Connection, size),
	}
	
	for i := 0; i < size; i++ {
		pool.conns <- &Connection{ID: i}
	}
	
	return pool
}

func (cp *ConnectionPool) Get() *Connection {
	return <-cp.conns
}

func (cp *ConnectionPool) Put(conn *Connection) {
	cp.conns <- conn
}

func (s *Server) handleWithConnectionPool(pool *ConnectionPool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn := pool.Get()
		defer pool.Put(conn)
		
		fmt.Fprintf(w, "使用连接 %d 处理请求\n", conn.ID)
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	fmt.Println("Go 异步编程 - HTTP 服务器实战")
	fmt.Println("==============================")
	
	server := NewServer()
	
	http.HandleFunc("/", server.handleHome)
	http.HandleFunc("/slow", server.handleSlowOperation)
	http.HandleFunc("/concurrent", server.handleConcurrentRequests)
	
	// 速率限制
	rateLimiter := NewRateLimiter(10)
	http.HandleFunc("/limited", rateLimitMiddleware(rateLimiter, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "这是受速率限制的端点\n")
	}))
	
	// 工作池
	pool := NewWorkerPool(5)
	defer pool.Stop()
	http.HandleFunc("/worker", server.handleWithWorkerPool(pool))
	
	// 连接池
	connPool := NewConnectionPool(10)
	http.HandleFunc("/pool", server.handleWithConnectionPool(connPool))
	
	fmt.Println("服务器启动在 :8080")
	fmt.Println("端点:")
	fmt.Println("  GET /           - 基本计数器")
	fmt.Println("  GET /slow       - 超时处理")
	fmt.Println("  GET /concurrent - 并发请求")
	fmt.Println("  GET /limited    - 速率限制")
	fmt.Println("  GET /worker     - 工作池")
	fmt.Println("  GET /pool       - 连接池")
	fmt.Println("\n注意: 这是一个演示程序，运行后可以通过 curl 测试")
	fmt.Println("示例: curl http://localhost:8080/")
	
	// 不实际启动服务器，只展示代码结构
	// log.Fatal(http.ListenAndServe(":8080", nil))
}

# HTTP/2 和高级特性 / HTTP/2 & Advanced Features

```go
package main

import (
    "context"
    "fmt"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func main() {
    mux := http.NewServeMux()
    
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Hello, HTTP/2!")
    })
    
    // Server-Sent Events (SSE)
    mux.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/event-stream")
        w.Header().Set("Cache-Control", "no-cache")
        w.Header().Set("Connection", "keep-alive")
        
        flusher, ok := w.(http.Flusher)
        if !ok {
            http.Error(w, "SSE not supported", http.StatusInternalServerError)
            return
        }
        
        for i := 0; i < 5; i++ {
            fmt.Fprintf(w, "data: Message %d\n\n", i)
            flusher.Flush()
            time.Sleep(time.Second)
        }
    })
    
    // 优雅关闭 / Graceful shutdown
    server := &http.Server{
        Addr:    ":8080",
        Handler: mux,
    }
    
    // 启动服务器 / Start server
    go func() {
        fmt.Println("Server starting on :8080")
        if err := server.ListenAndServe(); err != http.ErrServerClosed {
            fmt.Println("Server error:", err)
        }
    }()
    
    // 等待中断信号 / Wait for interrupt signal
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    fmt.Println("Shutting down server...")
    
    // 给予 30 秒的超时时间 / Allow 30 seconds timeout
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    if err := server.Shutdown(ctx); err != nil {
        fmt.Println("Server shutdown error:", err)
    }
    
    fmt.Println("Server stopped")
}
```


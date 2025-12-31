# HTTP 服务器 / HTTP Server

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "time"
)

// 基本处理函数 / Basic handler function
func helloHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}

// JSON 响应 / JSON response
func jsonHandler(w http.ResponseWriter, r *http.Request) {
    data := map[string]interface{}{
        "message": "Hello, JSON!",
        "status":  "success",
        "time":    time.Now().Format(time.RFC3339),
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(data)
}

// 处理不同方法 / Handle different methods
func methodHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        fmt.Fprintln(w, "GET request")
    case http.MethodPost:
        fmt.Fprintln(w, "POST request")
    case http.MethodPut:
        fmt.Fprintln(w, "PUT request")
    case http.MethodDelete:
        fmt.Fprintln(w, "DELETE request")
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}

// 读取请求体 / Read request body
func bodyHandler(w http.ResponseWriter, r *http.Request) {
    var data map[string]interface{}
    if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    fmt.Fprintf(w, "Received: %v", data)
}

// 查询参数 / Query parameters
func queryHandler(w http.ResponseWriter, r *http.Request) {
    name := r.URL.Query().Get("name")
    age := r.URL.Query().Get("age")
    fmt.Fprintf(w, "Name: %s, Age: %s", name, age)
}

// 中间件 / Middleware
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        log.Printf("Started %s %s", r.Method, r.URL.Path)
        next.ServeHTTP(w, r)
        log.Printf("Completed in %v", time.Since(start))
    })
}

func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        next.ServeHTTP(w, r)
    })
}

// 自定义 Handler 结构体 / Custom Handler struct
type APIHandler struct {
    version string
}

func (h *APIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "API Version: %s", h.version)
}

func main() {
    // 使用 DefaultServeMux / Using DefaultServeMux
    http.HandleFunc("/hello/", helloHandler)
    http.HandleFunc("/json", jsonHandler)
    http.HandleFunc("/method", methodHandler)
    http.HandleFunc("/body", bodyHandler)
    http.HandleFunc("/query", queryHandler)
    
    // 使用自定义 Handler / Using custom Handler
    http.Handle("/api/", &APIHandler{version: "1.0"})
    
    // 静态文件服务 / Static file server
    fs := http.FileServer(http.Dir("./static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))
    
    // 使用自定义 ServeMux / Using custom ServeMux
    mux := http.NewServeMux()
    mux.HandleFunc("/custom", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Custom mux")
    })
    
    // 应用中间件 / Apply middleware
    handler := loggingMiddleware(mux)
    
    // 配置服务器 / Configure server
    server := &http.Server{
        Addr:         ":8080",
        Handler:      handler,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout:  60 * time.Second,
    }
    
    fmt.Println("Server starting on :8080")
    // log.Fatal(server.ListenAndServe())
    
    // HTTPS 服务器 / HTTPS server
    // log.Fatal(server.ListenAndServeTLS("cert.pem", "key.pem"))
    
    _ = server
}
```


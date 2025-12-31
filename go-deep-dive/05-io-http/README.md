# Go IO 与 HTTP / Go IO & HTTP

## 1. 文件 IO / File IO

```go
package main

import (
    "bufio"
    "fmt"
    "io"
    "os"
    "path/filepath"
)

func main() {
    // 读取整个文件 / Read entire file
    data, err := os.ReadFile("example.txt")
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("Content:", string(data))
    }
    
    // 写入文件 / Write file
    err = os.WriteFile("output.txt", []byte("Hello, Go!"), 0644)
    if err != nil {
        fmt.Println("Error:", err)
    }
    
    // 打开文件 / Open file
    file, err := os.Open("example.txt")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer file.Close()
    
    // 逐行读取 / Read line by line
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        fmt.Println("Line:", scanner.Text())
    }
    
    // 创建文件 / Create file
    newFile, err := os.Create("new.txt")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer newFile.Close()
    
    // 写入字符串 / Write string
    newFile.WriteString("First line\n")
    newFile.WriteString("Second line\n")
    
    // 使用 bufio.Writer / Using bufio.Writer
    writer := bufio.NewWriter(newFile)
    writer.WriteString("Buffered line\n")
    writer.Flush()  // 确保写入磁盘
    
    // 追加模式 / Append mode
    appendFile, err := os.OpenFile("output.txt", os.O_APPEND|os.O_WRONLY, 0644)
    if err == nil {
        defer appendFile.Close()
        appendFile.WriteString("\nAppended content")
    }
    
    // 文件信息 / File info
    info, err := os.Stat("example.txt")
    if err == nil {
        fmt.Println("Name:", info.Name())
        fmt.Println("Size:", info.Size())
        fmt.Println("Mode:", info.Mode())
        fmt.Println("ModTime:", info.ModTime())
        fmt.Println("IsDir:", info.IsDir())
    }
    
    // 检查文件是否存在 / Check if file exists
    if _, err := os.Stat("example.txt"); os.IsNotExist(err) {
        fmt.Println("File does not exist")
    }
    
    // 删除文件 / Delete file
    // os.Remove("file.txt")
    
    // 重命名/移动文件 / Rename/move file
    // os.Rename("old.txt", "new.txt")
    
    // 复制文件 / Copy file
    copyFile := func(src, dst string) error {
        source, err := os.Open(src)
        if err != nil {
            return err
        }
        defer source.Close()
        
        dest, err := os.Create(dst)
        if err != nil {
            return err
        }
        defer dest.Close()
        
        _, err = io.Copy(dest, source)
        return err
    }
    _ = copyFile
}
```

## 2. 目录操作 / Directory Operations

```go
package main

import (
    "fmt"
    "os"
    "path/filepath"
)

func main() {
    // 创建目录 / Create directory
    err := os.Mkdir("mydir", 0755)
    if err != nil && !os.IsExist(err) {
        fmt.Println("Error:", err)
    }
    
    // 创建多级目录 / Create nested directories
    err = os.MkdirAll("path/to/nested/dir", 0755)
    if err != nil {
        fmt.Println("Error:", err)
    }
    
    // 读取目录 / Read directory
    entries, err := os.ReadDir(".")
    if err == nil {
        for _, entry := range entries {
            info, _ := entry.Info()
            fmt.Printf("%s (dir: %v, size: %d)\n", entry.Name(), entry.IsDir(), info.Size())
        }
    }
    
    // 遍历目录树 / Walk directory tree
    filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        fmt.Printf("Path: %s, IsDir: %v\n", path, info.IsDir())
        return nil
    })
    
    // 使用 WalkDir (更高效) / Using WalkDir (more efficient)
    filepath.WalkDir(".", func(path string, d os.DirEntry, err error) error {
        if err != nil {
            return err
        }
        fmt.Printf("Path: %s, IsDir: %v\n", path, d.IsDir())
        return nil
    })
    
    // 获取当前目录 / Get current directory
    cwd, _ := os.Getwd()
    fmt.Println("Current directory:", cwd)
    
    // 改变目录 / Change directory
    // os.Chdir("/tmp")
    
    // 删除目录 / Remove directory
    // os.Remove("empty_dir")  // 只能删除空目录
    // os.RemoveAll("dir_with_contents")  // 递归删除
    
    // 路径操作 / Path operations
    path := "/home/user/documents/file.txt"
    fmt.Println("Dir:", filepath.Dir(path))       // /home/user/documents
    fmt.Println("Base:", filepath.Base(path))     // file.txt
    fmt.Println("Ext:", filepath.Ext(path))       // .txt
    fmt.Println("Clean:", filepath.Clean("a//b/../c"))  // a/c
    
    // 拼接路径 / Join paths
    joined := filepath.Join("home", "user", "file.txt")
    fmt.Println("Joined:", joined)
    
    // 绝对路径 / Absolute path
    abs, _ := filepath.Abs(".")
    fmt.Println("Absolute:", abs)
    
    // 匹配模式 / Match pattern
    matched, _ := filepath.Match("*.txt", "file.txt")
    fmt.Println("Matched:", matched)
    
    // Glob 匹配 / Glob matching
    files, _ := filepath.Glob("*.go")
    fmt.Println("Go files:", files)
    
    // 临时文件和目录 / Temp files and directories
    tmpFile, err := os.CreateTemp("", "prefix-*.txt")
    if err == nil {
        fmt.Println("Temp file:", tmpFile.Name())
        tmpFile.Close()
        os.Remove(tmpFile.Name())
    }
    
    tmpDir, err := os.MkdirTemp("", "prefix-")
    if err == nil {
        fmt.Println("Temp dir:", tmpDir)
        os.RemoveAll(tmpDir)
    }
}
```

## 3. io.Reader 和 io.Writer / io.Reader & io.Writer

```go
package main

import (
    "bytes"
    "fmt"
    "io"
    "os"
    "strings"
)

func main() {
    // io.Reader 接口 / io.Reader interface
    // type Reader interface {
    //     Read(p []byte) (n int, err error)
    // }
    
    // 从 strings.Reader 读取 / Read from strings.Reader
    r := strings.NewReader("Hello, Reader!")
    buf := make([]byte, 4)
    
    for {
        n, err := r.Read(buf)
        if err == io.EOF {
            break
        }
        if err != nil {
            fmt.Println("Error:", err)
            break
        }
        fmt.Printf("Read %d bytes: %s\n", n, buf[:n])
    }
    
    // io.Writer 接口 / io.Writer interface
    // type Writer interface {
    //     Write(p []byte) (n int, err error)
    // }
    
    // 写入到 bytes.Buffer / Write to bytes.Buffer
    var buffer bytes.Buffer
    buffer.WriteString("Hello, ")
    buffer.WriteString("Writer!")
    fmt.Println("Buffer:", buffer.String())
    
    // io.Copy - 复制数据 / io.Copy - copy data
    src := strings.NewReader("Copy this content")
    dst := &bytes.Buffer{}
    n, err := io.Copy(dst, src)
    fmt.Printf("Copied %d bytes: %s\n", n, dst.String())
    
    // io.CopyN - 复制指定字节数 / io.CopyN - copy specific bytes
    src2 := strings.NewReader("Only copy part of this")
    dst2 := &bytes.Buffer{}
    io.CopyN(dst2, src2, 9)
    fmt.Println("CopyN:", dst2.String())
    
    // io.TeeReader - 同时读取和写入 / io.TeeReader - read and write simultaneously
    src3 := strings.NewReader("Tee data")
    var mirror bytes.Buffer
    tee := io.TeeReader(src3, &mirror)
    
    result, _ := io.ReadAll(tee)
    fmt.Println("Read:", string(result))
    fmt.Println("Mirror:", mirror.String())
    
    // io.MultiReader - 合并多个 Reader / io.MultiReader - combine readers
    r1 := strings.NewReader("First ")
    r2 := strings.NewReader("Second ")
    r3 := strings.NewReader("Third")
    multi := io.MultiReader(r1, r2, r3)
    
    combined, _ := io.ReadAll(multi)
    fmt.Println("Combined:", string(combined))
    
    // io.MultiWriter - 写入多个 Writer / io.MultiWriter - write to multiple writers
    var buf1, buf2 bytes.Buffer
    multiWriter := io.MultiWriter(&buf1, &buf2, os.Stdout)
    multiWriter.Write([]byte("Written to all\n"))
    
    // io.Pipe - 管道 / io.Pipe - pipe
    pr, pw := io.Pipe()
    
    go func() {
        pw.Write([]byte("Piped data"))
        pw.Close()
    }()
    
    pipeData, _ := io.ReadAll(pr)
    fmt.Println("Pipe:", string(pipeData))
    
    // io.LimitReader - 限制读取字节数 / io.LimitReader - limit bytes to read
    limited := io.LimitReader(strings.NewReader("Long content here"), 4)
    limitedData, _ := io.ReadAll(limited)
    fmt.Println("Limited:", string(limitedData))
    
    // io.ReadAll - 读取所有数据 / io.ReadAll - read all data
    fullData, _ := io.ReadAll(strings.NewReader("Full content"))
    fmt.Println("Full:", string(fullData))
}
```

## 4. HTTP 客户端 / HTTP Client

```go
package main

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "time"
)

func main() {
    // 简单 GET 请求 / Simple GET request
    resp, err := http.Get("https://httpbin.org/get")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer resp.Body.Close()
    
    body, _ := io.ReadAll(resp.Body)
    fmt.Println("Status:", resp.Status)
    fmt.Println("Body:", string(body)[:100], "...")
    
    // POST 请求 / POST request
    jsonData := []byte(`{"name":"Go","type":"language"}`)
    resp2, err := http.Post(
        "https://httpbin.org/post",
        "application/json",
        bytes.NewBuffer(jsonData),
    )
    if err == nil {
        defer resp2.Body.Close()
        fmt.Println("POST Status:", resp2.Status)
    }
    
    // 自定义请求 / Custom request
    client := &http.Client{
        Timeout: 10 * time.Second,
    }
    
    req, _ := http.NewRequest("GET", "https://httpbin.org/headers", nil)
    req.Header.Set("User-Agent", "Go-Client/1.0")
    req.Header.Set("Accept", "application/json")
    req.Header.Set("Authorization", "Bearer token123")
    
    resp3, err := client.Do(req)
    if err == nil {
        defer resp3.Body.Close()
        fmt.Println("Custom request status:", resp3.Status)
    }
    
    // 带 context 的请求 (支持取消和超时) / Request with context
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    req4, _ := http.NewRequestWithContext(ctx, "GET", "https://httpbin.org/delay/1", nil)
    resp4, err := client.Do(req4)
    if err != nil {
        if ctx.Err() == context.DeadlineExceeded {
            fmt.Println("Request timed out")
        } else {
            fmt.Println("Error:", err)
        }
    } else {
        defer resp4.Body.Close()
        fmt.Println("Context request status:", resp4.Status)
    }
    
    // POST 表单数据 / POST form data
    formData := url.Values{
        "username": {"john"},
        "password": {"secret"},
    }
    resp5, err := http.PostForm("https://httpbin.org/post", formData)
    if err == nil {
        defer resp5.Body.Close()
        fmt.Println("Form POST status:", resp5.Status)
    }
    
    // 配置 Transport / Configure Transport
    transport := &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
        IdleConnTimeout:     90 * time.Second,
    }
    
    customClient := &http.Client{
        Transport: transport,
        Timeout:   30 * time.Second,
    }
    _ = customClient
    
    // JSON 请求和响应 / JSON request and response
    type User struct {
        Name  string `json:"name"`
        Email string `json:"email"`
    }
    
    user := User{Name: "Alice", Email: "alice@example.com"}
    jsonBytes, _ := json.Marshal(user)
    
    resp6, err := http.Post(
        "https://httpbin.org/post",
        "application/json",
        bytes.NewBuffer(jsonBytes),
    )
    if err == nil {
        defer resp6.Body.Close()
        
        var result map[string]interface{}
        json.NewDecoder(resp6.Body).Decode(&result)
        fmt.Println("JSON response:", result["json"])
    }
}
```

## 5. HTTP 服务器 / HTTP Server

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

## 6. HTTP/2 和高级特性 / HTTP/2 & Advanced Features

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

## 7. 网络编程基础 / Network Programming Basics

```go
package main

import (
    "bufio"
    "fmt"
    "net"
    "time"
)

func main() {
    // TCP 客户端 / TCP client
    conn, err := net.Dial("tcp", "example.com:80")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer conn.Close()
    
    // 发送 HTTP 请求 / Send HTTP request
    fmt.Fprintf(conn, "GET / HTTP/1.1\r\nHost: example.com\r\n\r\n")
    
    // 读取响应 / Read response
    reader := bufio.NewReader(conn)
    line, _ := reader.ReadString('\n')
    fmt.Println("Response:", line)
    
    // 带超时的连接 / Connection with timeout
    conn2, err := net.DialTimeout("tcp", "example.com:80", 5*time.Second)
    if err == nil {
        conn2.Close()
    }
    
    // TCP 服务器 / TCP server
    listener, err := net.Listen("tcp", ":9000")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer listener.Close()
    
    fmt.Println("TCP server listening on :9000")
    
    // 接受连接 / Accept connections
    go func() {
        for {
            conn, err := listener.Accept()
            if err != nil {
                continue
            }
            
            go handleConnection(conn)
        }
    }()
    
    // UDP 通信 / UDP communication
    udpAddr, _ := net.ResolveUDPAddr("udp", ":9001")
    udpConn, err := net.ListenUDP("udp", udpAddr)
    if err == nil {
        defer udpConn.Close()
        fmt.Println("UDP server listening on :9001")
    }
    
    // DNS 查询 / DNS lookup
    ips, err := net.LookupIP("google.com")
    if err == nil {
        fmt.Println("IPs for google.com:")
        for _, ip := range ips {
            fmt.Println(" ", ip)
        }
    }
    
    // 解析地址 / Parse address
    host, port, _ := net.SplitHostPort("localhost:8080")
    fmt.Printf("Host: %s, Port: %s\n", host, port)
    
    time.Sleep(100 * time.Millisecond)
}

func handleConnection(conn net.Conn) {
    defer conn.Close()
    
    // 设置读取超时 / Set read timeout
    conn.SetReadDeadline(time.Now().Add(5 * time.Second))
    
    buffer := make([]byte, 1024)
    n, err := conn.Read(buffer)
    if err != nil {
        return
    }
    
    fmt.Println("Received:", string(buffer[:n]))
    conn.Write([]byte("Message received\n"))
}
```

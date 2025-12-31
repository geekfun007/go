# Go IO 与 HTTP 详解

## 1. IO 基础

### 1.1 io 包核心接口

```go
package main

import (
    "bytes"
    "fmt"
    "io"
    "strings"
)

func main() {
    // Reader 接口示例
    fmt.Println("=== Reader ===")
    readerExample()
    
    // Writer 接口示例
    fmt.Println("\n=== Writer ===")
    writerExample()
    
    // io.Copy
    fmt.Println("\n=== io.Copy ===")
    copyExample()
    
    // io.TeeReader
    fmt.Println("\n=== TeeReader ===")
    teeReaderExample()
    
    // io.MultiReader / io.MultiWriter
    fmt.Println("\n=== Multi ===")
    multiExample()
    
    // io.Pipe
    fmt.Println("\n=== Pipe ===")
    pipeExample()
}

func readerExample() {
    // strings.Reader
    r := strings.NewReader("Hello, World!")
    
    buf := make([]byte, 5)
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
    
    // 读取全部
    r.Reset("Hello again!")
    data, _ := io.ReadAll(r)
    fmt.Println("ReadAll:", string(data))
}

func writerExample() {
    // bytes.Buffer 实现了 Writer
    var buf bytes.Buffer
    
    buf.WriteString("Hello, ")
    buf.Write([]byte("World!"))
    buf.WriteByte('\n')
    
    fmt.Print("Buffer:", buf.String())
    
    // 使用 io.WriteString
    var buf2 bytes.Buffer
    io.WriteString(&buf2, "Using io.WriteString")
    fmt.Println(buf2.String())
}

func copyExample() {
    src := strings.NewReader("Data to copy")
    var dst bytes.Buffer
    
    n, err := io.Copy(&dst, src)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    
    fmt.Printf("Copied %d bytes: %s\n", n, dst.String())
    
    // CopyN - 复制指定字节数
    src2 := strings.NewReader("Only copy 5 bytes from this")
    var dst2 bytes.Buffer
    
    n, _ = io.CopyN(&dst2, src2, 5)
    fmt.Printf("CopyN %d bytes: %s\n", n, dst2.String())
}

func teeReaderExample() {
    src := strings.NewReader("Tee data")
    var log bytes.Buffer
    
    // TeeReader: 读取时同时写入另一个 Writer
    tee := io.TeeReader(src, &log)
    
    data, _ := io.ReadAll(tee)
    fmt.Println("Read:", string(data))
    fmt.Println("Log:", log.String())
}

func multiExample() {
    // MultiReader: 顺序读取多个 Reader
    r1 := strings.NewReader("Part 1, ")
    r2 := strings.NewReader("Part 2, ")
    r3 := strings.NewReader("Part 3")
    
    multi := io.MultiReader(r1, r2, r3)
    data, _ := io.ReadAll(multi)
    fmt.Println("MultiReader:", string(data))
    
    // MultiWriter: 同时写入多个 Writer
    var buf1, buf2 bytes.Buffer
    mw := io.MultiWriter(&buf1, &buf2)
    
    mw.Write([]byte("Write to both"))
    fmt.Println("buf1:", buf1.String())
    fmt.Println("buf2:", buf2.String())
}

func pipeExample() {
    pr, pw := io.Pipe()
    
    // 写入 goroutine
    go func() {
        defer pw.Close()
        for i := 0; i < 3; i++ {
            fmt.Fprintf(pw, "Message %d\n", i)
        }
    }()
    
    // 读取
    data, _ := io.ReadAll(pr)
    fmt.Print("Pipe data:\n", string(data))
}
```

### 1.2 文件操作

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
    // 创建和写入文件
    fmt.Println("=== Write File ===")
    writeFile()
    
    // 读取文件
    fmt.Println("\n=== Read File ===")
    readFile()
    
    // bufio 缓冲读写
    fmt.Println("\n=== Buffered IO ===")
    bufferedIO()
    
    // 文件操作
    fmt.Println("\n=== File Operations ===")
    fileOperations()
    
    // 目录操作
    fmt.Println("\n=== Directory Operations ===")
    directoryOperations()
    
    // 清理
    cleanup()
}

func writeFile() {
    // 方式1: os.WriteFile (简单方式)
    err := os.WriteFile("/tmp/test1.txt", []byte("Hello, World!"), 0644)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Written with WriteFile")
    
    // 方式2: os.Create + Write
    f, err := os.Create("/tmp/test2.txt")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer f.Close()
    
    f.WriteString("Line 1\n")
    f.Write([]byte("Line 2\n"))
    fmt.Fprintf(f, "Line %d\n", 3)
    
    fmt.Println("Written with Create")
    
    // 方式3: 追加模式
    f2, err := os.OpenFile("/tmp/test2.txt", os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer f2.Close()
    
    f2.WriteString("Appended line\n")
    fmt.Println("Appended to file")
}

func readFile() {
    // 方式1: os.ReadFile (简单方式)
    data, err := os.ReadFile("/tmp/test1.txt")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("ReadFile:", string(data))
    
    // 方式2: os.Open + Read
    f, err := os.Open("/tmp/test2.txt")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer f.Close()
    
    buf := make([]byte, 1024)
    n, err := f.Read(buf)
    if err != nil && err != io.EOF {
        fmt.Println("Error:", err)
        return
    }
    fmt.Printf("Read %d bytes:\n%s", n, buf[:n])
    
    // 方式3: 按行读取
    f.Seek(0, 0) // 重置到文件开头
    scanner := bufio.NewScanner(f)
    fmt.Println("\nLine by line:")
    for scanner.Scan() {
        fmt.Println(" ", scanner.Text())
    }
}

func bufferedIO() {
    // 缓冲写入
    f, _ := os.Create("/tmp/buffered.txt")
    defer f.Close()
    
    writer := bufio.NewWriter(f)
    writer.WriteString("Buffered write 1\n")
    writer.WriteString("Buffered write 2\n")
    writer.Flush() // 必须刷新
    
    fmt.Println("Buffered write complete")
    
    // 缓冲读取
    f2, _ := os.Open("/tmp/buffered.txt")
    defer f2.Close()
    
    reader := bufio.NewReader(f2)
    
    // ReadString
    line, _ := reader.ReadString('\n')
    fmt.Println("ReadString:", line)
    
    // ReadLine
    lineBytes, _, _ := reader.ReadLine()
    fmt.Println("ReadLine:", string(lineBytes))
    
    // Peek
    f2.Seek(0, 0)
    reader.Reset(f2)
    peeked, _ := reader.Peek(8)
    fmt.Println("Peek:", string(peeked))
}

func fileOperations() {
    // 文件信息
    info, err := os.Stat("/tmp/test1.txt")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    
    fmt.Println("Name:", info.Name())
    fmt.Println("Size:", info.Size())
    fmt.Println("Mode:", info.Mode())
    fmt.Println("ModTime:", info.ModTime())
    fmt.Println("IsDir:", info.IsDir())
    
    // 检查文件是否存在
    if _, err := os.Stat("/tmp/nonexistent.txt"); os.IsNotExist(err) {
        fmt.Println("File does not exist")
    }
    
    // 重命名
    os.Rename("/tmp/test1.txt", "/tmp/test1_renamed.txt")
    
    // 复制文件
    copyFile("/tmp/test1_renamed.txt", "/tmp/test1_copy.txt")
    
    // 删除
    os.Remove("/tmp/test1_copy.txt")
    
    // 还原名称
    os.Rename("/tmp/test1_renamed.txt", "/tmp/test1.txt")
}

func copyFile(src, dst string) error {
    source, err := os.Open(src)
    if err != nil {
        return err
    }
    defer source.Close()
    
    destination, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer destination.Close()
    
    _, err = io.Copy(destination, source)
    return err
}

func directoryOperations() {
    // 创建目录
    os.Mkdir("/tmp/testdir", 0755)
    os.MkdirAll("/tmp/testdir/sub1/sub2", 0755)
    
    // 创建临时文件和目录
    tmpFile, _ := os.CreateTemp("/tmp", "prefix-*.txt")
    fmt.Println("Temp file:", tmpFile.Name())
    tmpFile.Close()
    os.Remove(tmpFile.Name())
    
    tmpDir, _ := os.MkdirTemp("/tmp", "tmpdir-")
    fmt.Println("Temp dir:", tmpDir)
    os.RemoveAll(tmpDir)
    
    // 列出目录
    entries, _ := os.ReadDir("/tmp/testdir")
    fmt.Println("Directory contents:")
    for _, entry := range entries {
        info, _ := entry.Info()
        fmt.Printf("  %s (dir=%v, size=%d)\n", entry.Name(), entry.IsDir(), info.Size())
    }
    
    // 遍历目录树
    fmt.Println("\nWalk directory:")
    filepath.Walk("/tmp/testdir", func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        fmt.Println(" ", path)
        return nil
    })
    
    // 路径操作
    fmt.Println("\nPath operations:")
    path := "/tmp/testdir/file.txt"
    fmt.Println("Dir:", filepath.Dir(path))
    fmt.Println("Base:", filepath.Base(path))
    fmt.Println("Ext:", filepath.Ext(path))
    fmt.Println("Clean:", filepath.Clean("/tmp/../tmp/./testdir"))
    fmt.Println("Join:", filepath.Join("/tmp", "testdir", "file.txt"))
    
    // Glob 匹配
    matches, _ := filepath.Glob("/tmp/test*.txt")
    fmt.Println("Glob matches:", matches)
}

func cleanup() {
    os.Remove("/tmp/test1.txt")
    os.Remove("/tmp/test2.txt")
    os.Remove("/tmp/buffered.txt")
    os.RemoveAll("/tmp/testdir")
}
```

### 1.3 JSON 处理

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "os"
    "strings"
)

type Person struct {
    Name    string   `json:"name"`
    Age     int      `json:"age"`
    Email   string   `json:"email,omitempty"`
    Tags    []string `json:"tags,omitempty"`
    Private string   `json:"-"`
}

type Address struct {
    City    string `json:"city"`
    Country string `json:"country"`
}

type Employee struct {
    Person
    Company string  `json:"company"`
    Address Address `json:"address"`
}

func main() {
    // Marshal
    fmt.Println("=== Marshal ===")
    marshalExample()
    
    // Unmarshal
    fmt.Println("\n=== Unmarshal ===")
    unmarshalExample()
    
    // Encoder/Decoder (流式)
    fmt.Println("\n=== Encoder/Decoder ===")
    encoderDecoderExample()
    
    // 动态 JSON
    fmt.Println("\n=== Dynamic JSON ===")
    dynamicJSONExample()
    
    // JSON 流处理
    fmt.Println("\n=== JSON Stream ===")
    jsonStreamExample()
}

func marshalExample() {
    person := Person{
        Name:    "Alice",
        Age:     30,
        Email:   "alice@example.com",
        Tags:    []string{"developer", "golang"},
        Private: "secret",
    }
    
    // 基本 Marshal
    data, err := json.Marshal(person)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Marshal:", string(data))
    
    // 格式化输出
    prettyData, _ := json.MarshalIndent(person, "", "  ")
    fmt.Println("Pretty:\n", string(prettyData))
    
    // 嵌套结构
    emp := Employee{
        Person:  Person{Name: "Bob", Age: 25},
        Company: "TechCorp",
        Address: Address{City: "Beijing", Country: "China"},
    }
    
    empData, _ := json.MarshalIndent(emp, "", "  ")
    fmt.Println("Employee:\n", string(empData))
}

func unmarshalExample() {
    jsonStr := `{
        "name": "Charlie",
        "age": 35,
        "email": "charlie@example.com",
        "tags": ["manager", "scrum"]
    }`
    
    var person Person
    err := json.Unmarshal([]byte(jsonStr), &person)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    
    fmt.Printf("Unmarshal: %+v\n", person)
    
    // 部分解析
    partialJSON := `{"name": "David", "age": 40, "extra": "ignored"}`
    var partial Person
    json.Unmarshal([]byte(partialJSON), &partial)
    fmt.Printf("Partial: %+v\n", partial)
}

func encoderDecoderExample() {
    // Encoder - 写入 Writer
    var buf bytes.Buffer
    encoder := json.NewEncoder(&buf)
    encoder.SetIndent("", "  ")
    
    encoder.Encode(Person{Name: "Eve", Age: 28})
    fmt.Println("Encoded:", buf.String())
    
    // Decoder - 从 Reader 读取
    jsonInput := `{"name": "Frank", "age": 32}`
    decoder := json.NewDecoder(strings.NewReader(jsonInput))
    
    var person Person
    decoder.Decode(&person)
    fmt.Printf("Decoded: %+v\n", person)
    
    // 写入文件
    f, _ := os.Create("/tmp/person.json")
    defer f.Close()
    defer os.Remove("/tmp/person.json")
    
    json.NewEncoder(f).Encode(Person{Name: "Grace", Age: 29})
    
    // 从文件读取
    f.Seek(0, 0)
    var fromFile Person
    json.NewDecoder(f).Decode(&fromFile)
    fmt.Printf("From file: %+v\n", fromFile)
}

func dynamicJSONExample() {
    // map[string]interface{}
    jsonStr := `{
        "name": "Test",
        "count": 42,
        "active": true,
        "tags": ["a", "b"],
        "nested": {"key": "value"}
    }`
    
    var data map[string]interface{}
    json.Unmarshal([]byte(jsonStr), &data)
    
    fmt.Println("Name:", data["name"])
    fmt.Println("Count:", data["count"])
    fmt.Println("Tags:", data["tags"])
    
    // 类型断言
    if tags, ok := data["tags"].([]interface{}); ok {
        for i, tag := range tags {
            fmt.Printf("Tag %d: %s\n", i, tag)
        }
    }
    
    // 嵌套访问
    if nested, ok := data["nested"].(map[string]interface{}); ok {
        fmt.Println("Nested key:", nested["key"])
    }
    
    // json.RawMessage
    type Config struct {
        Type    string          `json:"type"`
        Options json.RawMessage `json:"options"`
    }
    
    configJSON := `{
        "type": "database",
        "options": {"host": "localhost", "port": 5432}
    }`
    
    var config Config
    json.Unmarshal([]byte(configJSON), &config)
    fmt.Println("Config type:", config.Type)
    fmt.Println("Raw options:", string(config.Options))
    
    // 根据 type 解析 options
    var dbOptions struct {
        Host string `json:"host"`
        Port int    `json:"port"`
    }
    json.Unmarshal(config.Options, &dbOptions)
    fmt.Printf("DB options: %+v\n", dbOptions)
}

func jsonStreamExample() {
    // 处理 JSON 数组流
    jsonArray := `[
        {"name": "A", "age": 20},
        {"name": "B", "age": 21},
        {"name": "C", "age": 22}
    ]`
    
    dec := json.NewDecoder(strings.NewReader(jsonArray))
    
    // 读取开始的 [
    _, _ = dec.Token()
    
    for dec.More() {
        var person Person
        dec.Decode(&person)
        fmt.Printf("Stream: %+v\n", person)
    }
    
    // 读取结束的 ]
    _, _ = dec.Token()
}
```

## 2. HTTP 客户端

### 2.1 基础请求

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "strings"
    "time"
)

func main() {
    // GET 请求
    fmt.Println("=== GET ===")
    getExample()
    
    // POST 请求
    fmt.Println("\n=== POST ===")
    postExample()
    
    // 自定义请求
    fmt.Println("\n=== Custom Request ===")
    customRequestExample()
    
    // 表单提交
    fmt.Println("\n=== Form ===")
    formExample()
    
    // 超时和取消
    fmt.Println("\n=== Timeout ===")
    timeoutExample()
}

func getExample() {
    // 简单 GET
    resp, err := http.Get("https://httpbin.org/get")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer resp.Body.Close()
    
    fmt.Println("Status:", resp.Status)
    fmt.Println("StatusCode:", resp.StatusCode)
    
    // 读取响应体
    body, _ := io.ReadAll(resp.Body)
    fmt.Println("Body length:", len(body))
    
    // 带查询参数
    baseURL := "https://httpbin.org/get"
    params := url.Values{}
    params.Add("name", "Alice")
    params.Add("age", "30")
    
    fullURL := baseURL + "?" + params.Encode()
    resp2, _ := http.Get(fullURL)
    defer resp2.Body.Close()
    
    body2, _ := io.ReadAll(resp2.Body)
    fmt.Println("With params:", string(body2)[:200])
}

func postExample() {
    // POST JSON
    data := map[string]interface{}{
        "name": "Bob",
        "age":  25,
    }
    jsonData, _ := json.Marshal(data)
    
    resp, err := http.Post(
        "https://httpbin.org/post",
        "application/json",
        bytes.NewBuffer(jsonData),
    )
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer resp.Body.Close()
    
    body, _ := io.ReadAll(resp.Body)
    fmt.Println("POST response:", string(body)[:300])
}

func customRequestExample() {
    // 创建请求
    req, err := http.NewRequest("PUT", "https://httpbin.org/put", strings.NewReader(`{"key":"value"}`))
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    
    // 设置头部
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer token123")
    req.Header.Set("X-Custom-Header", "custom-value")
    
    // 发送请求
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer resp.Body.Close()
    
    fmt.Println("PUT Status:", resp.Status)
    
    // 读取响应头
    fmt.Println("Response Headers:")
    for key, values := range resp.Header {
        fmt.Printf("  %s: %s\n", key, values[0])
    }
}

func formExample() {
    // 表单数据
    formData := url.Values{}
    formData.Set("username", "alice")
    formData.Set("password", "secret")
    
    resp, err := http.PostForm("https://httpbin.org/post", formData)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer resp.Body.Close()
    
    body, _ := io.ReadAll(resp.Body)
    fmt.Println("Form response:", string(body)[:300])
}

func timeoutExample() {
    // 设置超时
    client := &http.Client{
        Timeout: 5 * time.Second,
    }
    
    resp, err := client.Get("https://httpbin.org/delay/1")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer resp.Body.Close()
    
    fmt.Println("Timeout test passed, status:", resp.Status)
}
```

### 2.2 高级 HTTP 客户端

```go
package main

import (
    "context"
    "crypto/tls"
    "fmt"
    "io"
    "net"
    "net/http"
    "net/http/cookiejar"
    "time"
)

func main() {
    // 自定义 Transport
    fmt.Println("=== Custom Transport ===")
    customTransport()
    
    // Cookie 处理
    fmt.Println("\n=== Cookies ===")
    cookieExample()
    
    // Context 取消
    fmt.Println("\n=== Context Cancel ===")
    contextExample()
    
    // 重试机制
    fmt.Println("\n=== Retry ===")
    retryExample()
}

func customTransport() {
    transport := &http.Transport{
        // 连接池配置
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
        MaxConnsPerHost:     100,
        IdleConnTimeout:     90 * time.Second,
        
        // 超时配置
        DialContext: (&net.Dialer{
            Timeout:   30 * time.Second,
            KeepAlive: 30 * time.Second,
        }).DialContext,
        TLSHandshakeTimeout:   10 * time.Second,
        ResponseHeaderTimeout: 10 * time.Second,
        ExpectContinueTimeout: 1 * time.Second,
        
        // TLS 配置
        TLSClientConfig: &tls.Config{
            MinVersion: tls.VersionTLS12,
        },
    }
    
    client := &http.Client{
        Transport: transport,
        Timeout:   30 * time.Second,
    }
    
    resp, err := client.Get("https://httpbin.org/get")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer resp.Body.Close()
    
    fmt.Println("Custom transport status:", resp.Status)
}

func cookieExample() {
    // 创建 Cookie Jar
    jar, _ := cookiejar.New(nil)
    
    client := &http.Client{
        Jar: jar,
    }
    
    // 第一次请求，服务器设置 Cookie
    resp1, _ := client.Get("https://httpbin.org/cookies/set?name=alice")
    resp1.Body.Close()
    
    // 第二次请求，Cookie 会自动发送
    resp2, _ := client.Get("https://httpbin.org/cookies")
    body, _ := io.ReadAll(resp2.Body)
    resp2.Body.Close()
    
    fmt.Println("Cookies:", string(body))
}

func contextExample() {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()
    
    req, _ := http.NewRequestWithContext(ctx, "GET", "https://httpbin.org/delay/3", nil)
    
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("Request cancelled:", err)
        return
    }
    defer resp.Body.Close()
    
    fmt.Println("Status:", resp.Status)
}

func retryExample() {
    client := &http.Client{
        Timeout: 5 * time.Second,
    }
    
    maxRetries := 3
    var lastErr error
    
    for i := 0; i < maxRetries; i++ {
        resp, err := client.Get("https://httpbin.org/get")
        if err != nil {
            lastErr = err
            fmt.Printf("Attempt %d failed: %v\n", i+1, err)
            time.Sleep(time.Duration(i+1) * time.Second) // 指数退避
            continue
        }
        defer resp.Body.Close()
        
        if resp.StatusCode >= 500 {
            fmt.Printf("Attempt %d: server error %d\n", i+1, resp.StatusCode)
            continue
        }
        
        fmt.Println("Success after", i+1, "attempts")
        return
    }
    
    fmt.Println("All retries failed:", lastErr)
}
```

## 3. HTTP 服务器

### 3.1 基础服务器

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "time"
)

func main() {
    // 基础路由
    http.HandleFunc("/", homeHandler)
    http.HandleFunc("/hello", helloHandler)
    http.HandleFunc("/json", jsonHandler)
    http.HandleFunc("/users", usersHandler)
    
    // 静态文件
    fs := http.FileServer(http.Dir("./static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))
    
    fmt.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }
    fmt.Fprintf(w, "Welcome to Home!")
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
    name := r.URL.Query().Get("name")
    if name == "" {
        name = "World"
    }
    fmt.Fprintf(w, "Hello, %s!", name)
}

func jsonHandler(w http.ResponseWriter, r *http.Request) {
    data := map[string]interface{}{
        "message":   "Hello, JSON!",
        "timestamp": time.Now().Unix(),
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(data)
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        users := []map[string]string{
            {"id": "1", "name": "Alice"},
            {"id": "2", "name": "Bob"},
        }
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(users)
        
    case http.MethodPost:
        var user map[string]string
        json.NewDecoder(r.Body).Decode(&user)
        
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(map[string]string{
            "message": "User created",
            "name":    user["name"],
        })
        
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}
```

### 3.2 中间件和路由

```go
package main

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "strings"
    "time"
)

// 中间件类型
type Middleware func(http.Handler) http.Handler

// 日志中间件
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // 包装 ResponseWriter 以获取状态码
        wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
        
        next.ServeHTTP(wrapped, r)
        
        log.Printf("%s %s %d %v", r.Method, r.URL.Path, wrapped.statusCode, time.Since(start))
    })
}

// 认证中间件
func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" || !strings.HasPrefix(token, "Bearer ") {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        
        // 验证 token 并添加用户信息到 context
        userID := "user123" // 从 token 解析
        ctx := context.WithValue(r.Context(), "userID", userID)
        
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// CORS 中间件
func corsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        next.ServeHTTP(w, r)
    })
}

// 恢复中间件
func recoveryMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("Panic: %v", err)
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            }
        }()
        
        next.ServeHTTP(w, r)
    })
}

// 响应包装器
type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}

// 中间件链
func chain(handler http.Handler, middlewares ...Middleware) http.Handler {
    for i := len(middlewares) - 1; i >= 0; i-- {
        handler = middlewares[i](handler)
    }
    return handler
}

// 简单路由器
type Router struct {
    routes map[string]map[string]http.HandlerFunc
}

func NewRouter() *Router {
    return &Router{
        routes: make(map[string]map[string]http.HandlerFunc),
    }
}

func (r *Router) Handle(method, path string, handler http.HandlerFunc) {
    if r.routes[path] == nil {
        r.routes[path] = make(map[string]http.HandlerFunc)
    }
    r.routes[path][method] = handler
}

func (r *Router) GET(path string, handler http.HandlerFunc) {
    r.Handle(http.MethodGet, path, handler)
}

func (r *Router) POST(path string, handler http.HandlerFunc) {
    r.Handle(http.MethodPost, path, handler)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    if handlers, ok := r.routes[req.URL.Path]; ok {
        if handler, ok := handlers[req.Method]; ok {
            handler(w, req)
            return
        }
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    http.NotFound(w, req)
}

func main() {
    router := NewRouter()
    
    // 路由
    router.GET("/", func(w http.ResponseWriter, r *http.Request) {
        json.NewEncoder(w).Encode(map[string]string{"message": "Welcome"})
    })
    
    router.GET("/users", func(w http.ResponseWriter, r *http.Request) {
        json.NewEncoder(w).Encode([]string{"Alice", "Bob"})
    })
    
    router.POST("/users", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(map[string]string{"status": "created"})
    })
    
    // 应用中间件
    handler := chain(
        router,
        recoveryMiddleware,
        loggingMiddleware,
        corsMiddleware,
    )
    
    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", handler))
}
```

### 3.3 优雅关闭

```go
package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func main() {
    mux := http.NewServeMux()
    
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        time.Sleep(2 * time.Second) // 模拟慢请求
        w.Write([]byte("Hello, World!"))
    })
    
    server := &http.Server{
        Addr:         ":8080",
        Handler:      mux,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  60 * time.Second,
    }
    
    // 启动服务器
    go func() {
        log.Println("Server starting on :8080")
        if err := server.ListenAndServe(); err != http.ErrServerClosed {
            log.Fatalf("Server error: %v", err)
        }
    }()
    
    // 等待中断信号
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    log.Println("Shutting down server...")
    
    // 创建超时上下文
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    // 优雅关闭
    if err := server.Shutdown(ctx); err != nil {
        log.Fatalf("Server forced to shutdown: %v", err)
    }
    
    log.Println("Server stopped")
}
```

## 4. WebSocket

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "sync"

    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true // 允许所有来源
    },
}

// 连接管理器
type Hub struct {
    clients    map[*Client]bool
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
    mu         sync.Mutex
}

type Client struct {
    hub  *Hub
    conn *websocket.Conn
    send chan []byte
}

func newHub() *Hub {
    return &Hub{
        clients:    make(map[*Client]bool),
        broadcast:  make(chan []byte),
        register:   make(chan *Client),
        unregister: make(chan *Client),
    }
}

func (h *Hub) run() {
    for {
        select {
        case client := <-h.register:
            h.mu.Lock()
            h.clients[client] = true
            h.mu.Unlock()
            log.Println("Client connected")
            
        case client := <-h.unregister:
            h.mu.Lock()
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                close(client.send)
            }
            h.mu.Unlock()
            log.Println("Client disconnected")
            
        case message := <-h.broadcast:
            h.mu.Lock()
            for client := range h.clients {
                select {
                case client.send <- message:
                default:
                    close(client.send)
                    delete(h.clients, client)
                }
            }
            h.mu.Unlock()
        }
    }
}

func (c *Client) readPump() {
    defer func() {
        c.hub.unregister <- c
        c.conn.Close()
    }()
    
    for {
        _, message, err := c.conn.ReadMessage()
        if err != nil {
            break
        }
        log.Printf("Received: %s", message)
        c.hub.broadcast <- message
    }
}

func (c *Client) writePump() {
    defer c.conn.Close()
    
    for message := range c.send {
        if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
            return
        }
    }
}

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("Upgrade error:", err)
        return
    }
    
    client := &Client{
        hub:  hub,
        conn: conn,
        send: make(chan []byte, 256),
    }
    
    hub.register <- client
    
    go client.writePump()
    go client.readPump()
}

func main() {
    hub := newHub()
    go hub.run()
    
    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        serveWs(hub, w, r)
    })
    
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, `
<!DOCTYPE html>
<html>
<head><title>WebSocket Test</title></head>
<body>
    <input type="text" id="message">
    <button onclick="send()">Send</button>
    <div id="messages"></div>
    <script>
        var ws = new WebSocket("ws://localhost:8080/ws");
        ws.onmessage = function(e) {
            var div = document.createElement("div");
            div.textContent = e.data;
            document.getElementById("messages").appendChild(div);
        };
        function send() {
            ws.send(document.getElementById("message").value);
            document.getElementById("message").value = "";
        }
    </script>
</body>
</html>`)
    })
    
    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

## 5. HTTP/2 和 HTTPS

```go
package main

import (
    "crypto/tls"
    "fmt"
    "log"
    "net/http"
    "time"

    "golang.org/x/net/http2"
)

func main() {
    mux := http.NewServeMux()
    
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Protocol: %s\n", r.Proto)
        fmt.Fprintf(w, "Host: %s\n", r.Host)
        fmt.Fprintf(w, "Method: %s\n", r.Method)
    })
    
    // HTTP/2 Server Push
    mux.HandleFunc("/push", func(w http.ResponseWriter, r *http.Request) {
        pusher, ok := w.(http.Pusher)
        if ok {
            if err := pusher.Push("/style.css", nil); err != nil {
                log.Printf("Push failed: %v", err)
            }
        }
        
        fmt.Fprintf(w, `
<!DOCTYPE html>
<html>
<head><link rel="stylesheet" href="/style.css"></head>
<body><h1>HTTP/2 Server Push</h1></body>
</html>`)
    })
    
    mux.HandleFunc("/style.css", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/css")
        fmt.Fprintf(w, "body { background: #f0f0f0; }")
    })
    
    server := &http.Server{
        Addr:         ":8443",
        Handler:      mux,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
        TLSConfig: &tls.Config{
            MinVersion: tls.VersionTLS12,
        },
    }
    
    // 启用 HTTP/2
    http2.ConfigureServer(server, &http2.Server{})
    
    log.Println("HTTPS server starting on :8443")
    log.Fatal(server.ListenAndServeTLS("cert.pem", "key.pem"))
}

// 生成自签名证书:
// openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes
```

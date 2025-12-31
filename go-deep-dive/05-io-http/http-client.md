# HTTP 客户端 / HTTP Client

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


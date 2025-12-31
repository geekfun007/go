# 函数选项模式 / Functional Options Pattern

```go
package main

import (
    "fmt"
    "time"
)

// 服务配置 / Service configuration
type Server struct {
    host         string
    port         int
    timeout      time.Duration
    maxConns     int
    enableTLS    bool
    tlsCert      string
    tlsKey       string
}

// 选项函数类型 / Option function type
type ServerOption func(*Server)

// 选项函数 / Option functions
func WithHost(host string) ServerOption {
    return func(s *Server) {
        s.host = host
    }
}

func WithPort(port int) ServerOption {
    return func(s *Server) {
        s.port = port
    }
}

func WithTimeout(timeout time.Duration) ServerOption {
    return func(s *Server) {
        s.timeout = timeout
    }
}

func WithMaxConns(maxConns int) ServerOption {
    return func(s *Server) {
        s.maxConns = maxConns
    }
}

func WithTLS(cert, key string) ServerOption {
    return func(s *Server) {
        s.enableTLS = true
        s.tlsCert = cert
        s.tlsKey = key
    }
}

// 构造函数 / Constructor
func NewServer(options ...ServerOption) *Server {
    // 默认值 / Default values
    server := &Server{
        host:     "localhost",
        port:     8080,
        timeout:  30 * time.Second,
        maxConns: 100,
    }
    
    // 应用选项 / Apply options
    for _, opt := range options {
        opt(server)
    }
    
    return server
}

func (s *Server) Start() {
    fmt.Printf("Starting server on %s:%d\n", s.host, s.port)
    fmt.Printf("Timeout: %v, MaxConns: %d, TLS: %v\n", s.timeout, s.maxConns, s.enableTLS)
}

// 另一种变体: Builder 模式 / Another variant: Builder pattern
type ServerBuilder struct {
    server *Server
}

func NewServerBuilder() *ServerBuilder {
    return &ServerBuilder{
        server: &Server{
            host:     "localhost",
            port:     8080,
            timeout:  30 * time.Second,
            maxConns: 100,
        },
    }
}

func (b *ServerBuilder) Host(host string) *ServerBuilder {
    b.server.host = host
    return b
}

func (b *ServerBuilder) Port(port int) *ServerBuilder {
    b.server.port = port
    return b
}

func (b *ServerBuilder) Build() *Server {
    return b.server
}

func main() {
    // 使用函数选项模式 / Using functional options
    server1 := NewServer()  // 使用默认值
    server1.Start()
    
    server2 := NewServer(
        WithHost("0.0.0.0"),
        WithPort(443),
        WithTimeout(60*time.Second),
        WithTLS("cert.pem", "key.pem"),
    )
    server2.Start()
    
    // 使用 Builder 模式 / Using Builder pattern
    server3 := NewServerBuilder().
        Host("api.example.com").
        Port(9000).
        Build()
    server3.Start()
}
```

# 网络编程基础 / Network Programming Basics

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

# io.Reader 和 io.Writer / io.Reader & io.Writer

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


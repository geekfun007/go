# 文件 IO / File IO

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


# 目录操作 / Directory Operations

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


# 错误处理 / Error Handling

## 1. 错误基础 / Error Basics

```go
package main

import (
    "errors"
    "fmt"
)

func main() {
    // 创建错误 / Create errors
    err1 := errors.New("something went wrong")
    err2 := fmt.Errorf("failed to open %s", "file.txt")
    
    // 错误检查模式 / Error checking pattern
    result, err := doSomething()
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Result:", result)
}

func doSomething() (string, error) {
    return "", errors.New("operation failed")
}
```

## 2. 自定义错误类型 / Custom Error Types

```go
type ValidationError struct {
    Field   string
    Message string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("validation error: %s %s", e.Field, e.Message)
}

// 哨兵错误 / Sentinel errors
var (
    ErrNotFound     = errors.New("not found")
    ErrUnauthorized = errors.New("unauthorized")
)
```

## 3. 错误包装 (Go 1.13+) / Error Wrapping

```go
func readConfig(filename string) error {
    _, err := os.ReadFile(filename)
    if err != nil {
        return fmt.Errorf("read config: %w", err)  // 使用 %w 包装
    }
    return nil
}

func main() {
    err := readConfig("config.json")
    
    // errors.Is - 检查错误链
    if errors.Is(err, os.ErrNotExist) {
        fmt.Println("File not found")
    }
    
    // errors.As - 提取特定错误类型
    var pathErr *os.PathError
    if errors.As(err, &pathErr) {
        fmt.Println("Path:", pathErr.Path)
    }
    
    // errors.Unwrap - 获取被包装的错误
    unwrapped := errors.Unwrap(err)
}
```

## 4. 错误处理模式 / Error Handling Patterns

```go
// 模式1: 提前返回 / Early return
func processFile(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close()
    // ...
    return nil
}

// 模式2: 错误聚合 (Go 1.20+) / Error aggregation
func validate(form Form) error {
    var errs []error
    if form.Name == "" {
        errs = append(errs, errors.New("name required"))
    }
    if form.Email == "" {
        errs = append(errs, errors.New("email required"))
    }
    return errors.Join(errs...)
}

// 模式3: 重试 / Retry
func withRetry(attempts int, fn func() error) error {
    var lastErr error
    for i := 0; i < attempts; i++ {
        if err := fn(); err != nil {
            lastErr = err
            continue
        }
        return nil
    }
    return lastErr
}
```

## 5. panic 和 recover / panic & recover

```go
func main() {
    // recover 必须在 defer 中调用
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered:", r)
        }
    }()
    
    panic("something terrible")
}

// 安全调用 / Safe call
func safeCall(fn func()) (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("panic: %v", r)
        }
    }()
    fn()
    return nil
}
```

## 6. 最佳实践 / Best Practices

```go
// 1. 总是检查错误
result, err := mightFail()
if err != nil {
    return fmt.Errorf("context: %w", err)
}

// 2. 添加上下文
return fmt.Errorf("failed to process user %d: %w", userID, err)

// 3. 使用定义的错误变量
var ErrNotFound = errors.New("not found")
if errors.Is(err, ErrNotFound) { ... }

// 4. 错误只处理一次 (要么记录日志，要么返回，不要两者都做)

// 5. 使用 defer 清理资源
file, err := os.Open(name)
if err != nil {
    return err
}
defer file.Close()
```

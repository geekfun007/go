# Go 错误处理 / Go Error Handling

## 1. 错误基础 / Error Basics

```go
package main

import (
    "errors"
    "fmt"
)

func main() {
    // error 是一个接口 / error is an interface
    // type error interface {
    //     Error() string
    // }
    
    // 创建简单错误 / Create simple error
    err1 := errors.New("something went wrong")
    fmt.Println("Error 1:", err1)
    
    // 使用 fmt.Errorf 创建格式化错误 / Create formatted error
    name := "file.txt"
    err2 := fmt.Errorf("failed to open %s", name)
    fmt.Println("Error 2:", err2)
    
    // 错误检查模式 / Error checking pattern
    result, err := doSomething()
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Result:", result)
}

func doSomething() (string, error) {
    // 模拟可能失败的操作
    success := false
    if !success {
        return "", errors.New("operation failed")
    }
    return "success", nil
}
```

## 2. 自定义错误类型 / Custom Error Types

```go
package main

import (
    "fmt"
    "time"
)

// 自定义错误类型 / Custom error type
type ValidationError struct {
    Field   string
    Message string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("validation error: field '%s' %s", e.Field, e.Message)
}

// 带更多信息的错误 / Error with more info
type APIError struct {
    Code       int
    Message    string
    Timestamp  time.Time
    RequestID  string
}

func (e APIError) Error() string {
    return fmt.Sprintf("[%s] Error %d: %s (request: %s)",
        e.Timestamp.Format(time.RFC3339), e.Code, e.Message, e.RequestID)
}

// 可以添加辅助方法 / Can add helper methods
func (e APIError) IsNotFound() bool {
    return e.Code == 404
}

func (e APIError) IsServerError() bool {
    return e.Code >= 500
}

// 哨兵错误 / Sentinel errors
var (
    ErrNotFound     = errors.New("resource not found")
    ErrUnauthorized = errors.New("unauthorized access")
    ErrInvalidInput = errors.New("invalid input")
)

func main() {
    // 使用 ValidationError
    err := validateUser("", 15)
    if err != nil {
        fmt.Println(err)
    }
    
    // 使用 APIError
    apiErr := APIError{
        Code:      404,
        Message:   "User not found",
        Timestamp: time.Now(),
        RequestID: "abc123",
    }
    fmt.Println(apiErr)
    fmt.Println("Is not found:", apiErr.IsNotFound())
    
    // 使用哨兵错误
    err = findResource("xyz")
    if err == ErrNotFound {
        fmt.Println("Resource does not exist")
    }
}

func validateUser(name string, age int) error {
    if name == "" {
        return ValidationError{Field: "name", Message: "cannot be empty"}
    }
    if age < 18 {
        return ValidationError{Field: "age", Message: "must be at least 18"}
    }
    return nil
}

func findResource(id string) error {
    // 模拟资源查找
    return ErrNotFound
}
```

## 3. 错误包装 (Go 1.13+) / Error Wrapping

```go
package main

import (
    "errors"
    "fmt"
    "os"
)

func main() {
    // 使用 %w 包装错误 / Wrap error with %w
    err := readConfig("config.json")
    if err != nil {
        fmt.Println("Error:", err)
        
        // errors.Unwrap - 获取被包装的错误
        unwrapped := errors.Unwrap(err)
        fmt.Println("Unwrapped:", unwrapped)
        
        // errors.Is - 检查错误链中是否包含特定错误
        if errors.Is(err, os.ErrNotExist) {
            fmt.Println("File does not exist")
        }
        
        // errors.As - 检查并转换为特定错误类型
        var pathErr *os.PathError
        if errors.As(err, &pathErr) {
            fmt.Println("Path:", pathErr.Path)
            fmt.Println("Op:", pathErr.Op)
        }
    }
}

func readConfig(filename string) error {
    data, err := os.ReadFile(filename)
    if err != nil {
        // 包装错误，保留原始错误信息
        return fmt.Errorf("failed to read config file: %w", err)
    }
    fmt.Println("Config:", string(data))
    return nil
}

// 多层包装示例 / Multi-level wrapping example
func processData() error {
    if err := validateData(); err != nil {
        return fmt.Errorf("processData: %w", err)
    }
    return nil
}

func validateData() error {
    if err := loadData(); err != nil {
        return fmt.Errorf("validateData: %w", err)
    }
    return nil
}

func loadData() error {
    return fmt.Errorf("loadData: %w", ErrInvalidInput)
}

var ErrInvalidInput = errors.New("invalid input data")
```

## 4. errors.Is 和 errors.As / errors.Is & errors.As

```go
package main

import (
    "errors"
    "fmt"
    "io/fs"
    "os"
)

// 自定义错误类型，支持 Is 方法
type MyError struct {
    Code int
}

func (e MyError) Error() string {
    return fmt.Sprintf("error code: %d", e.Code)
}

// 实现 Is 方法以自定义比较逻辑
func (e MyError) Is(target error) bool {
    t, ok := target.(MyError)
    if !ok {
        return false
    }
    return e.Code == t.Code
}

// 自定义错误类型，支持 As 方法
type DetailedError struct {
    Inner error
    Msg   string
}

func (e DetailedError) Error() string {
    return fmt.Sprintf("%s: %v", e.Msg, e.Inner)
}

func (e DetailedError) Unwrap() error {
    return e.Inner
}

func main() {
    // errors.Is 示例 / errors.Is example
    
    // 比较哨兵错误
    err1 := fmt.Errorf("wrapped: %w", os.ErrNotExist)
    fmt.Println("Is ErrNotExist:", errors.Is(err1, os.ErrNotExist))  // true
    
    // 使用自定义 Is 方法
    err2 := MyError{Code: 404}
    err3 := fmt.Errorf("wrapped: %w", err2)
    fmt.Println("Is MyError{404}:", errors.Is(err3, MyError{Code: 404}))  // true
    fmt.Println("Is MyError{500}:", errors.Is(err3, MyError{Code: 500}))  // false
    
    // errors.As 示例 / errors.As example
    
    // 提取特定错误类型
    _, err := os.Open("nonexistent.txt")
    var pathErr *os.PathError
    if errors.As(err, &pathErr) {
        fmt.Println("Path error - Op:", pathErr.Op, "Path:", pathErr.Path)
    }
    
    // 检查是否为 fs.PathError
    var fsPathErr *fs.PathError
    if errors.As(err, &fsPathErr) {
        fmt.Println("FS Path error:", fsPathErr)
    }
    
    // 多层包装中提取错误
    detailedErr := DetailedError{
        Inner: os.ErrPermission,
        Msg:   "access denied",
    }
    wrappedErr := fmt.Errorf("operation failed: %w", detailedErr)
    
    var de DetailedError
    if errors.As(wrappedErr, &de) {
        fmt.Println("DetailedError msg:", de.Msg)
    }
    
    // 检查内部错误
    if errors.Is(wrappedErr, os.ErrPermission) {
        fmt.Println("Permission denied")
    }
}
```

## 5. 错误处理模式 / Error Handling Patterns

```go
package main

import (
    "errors"
    "fmt"
    "log"
)

// 模式1: 提前返回 / Pattern 1: Early return
func processFile(filename string) error {
    file, err := openFile(filename)
    if err != nil {
        return err
    }
    defer file.Close()
    
    data, err := readData(file)
    if err != nil {
        return err
    }
    
    return processData(data)
}

// 模式2: 错误聚合 / Pattern 2: Error aggregation
func validateForm(form Form) error {
    var errs []error
    
    if form.Name == "" {
        errs = append(errs, errors.New("name is required"))
    }
    if form.Email == "" {
        errs = append(errs, errors.New("email is required"))
    }
    if form.Age < 0 {
        errs = append(errs, errors.New("age must be positive"))
    }
    
    if len(errs) > 0 {
        return errors.Join(errs...)  // Go 1.20+
    }
    return nil
}

// 模式3: 错误回调 / Pattern 3: Error callback
type ErrorHandler func(error)

func withErrorHandler(handler ErrorHandler, fn func() error) {
    if err := fn(); err != nil {
        handler(err)
    }
}

// 模式4: 结果类型 / Pattern 4: Result type
type Result[T any] struct {
    Value T
    Err   error
}

func NewResult[T any](value T, err error) Result[T] {
    return Result[T]{Value: value, Err: err}
}

func (r Result[T]) Unwrap() (T, error) {
    return r.Value, r.Err
}

func (r Result[T]) UnwrapOr(defaultValue T) T {
    if r.Err != nil {
        return defaultValue
    }
    return r.Value
}

// 模式5: 重试 / Pattern 5: Retry
func withRetry(attempts int, fn func() error) error {
    var lastErr error
    for i := 0; i < attempts; i++ {
        if err := fn(); err != nil {
            lastErr = err
            log.Printf("Attempt %d failed: %v", i+1, err)
            continue
        }
        return nil
    }
    return fmt.Errorf("all %d attempts failed, last error: %w", attempts, lastErr)
}

type Form struct {
    Name  string
    Email string
    Age   int
}

func main() {
    // 模式2: 错误聚合
    form := Form{Name: "", Email: "", Age: -1}
    if err := validateForm(form); err != nil {
        fmt.Println("Validation errors:", err)
    }
    
    // 模式3: 错误回调
    withErrorHandler(func(err error) {
        fmt.Println("Error occurred:", err)
    }, func() error {
        return errors.New("something went wrong")
    })
    
    // 模式4: 结果类型
    result := NewResult(divide(10, 2))
    value := result.UnwrapOr(0)
    fmt.Println("Result:", value)
    
    // 模式5: 重试
    err := withRetry(3, func() error {
        return errors.New("temporary failure")
    })
    fmt.Println("Retry result:", err)
}

func divide(a, b int) (int, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

// 占位函数
func openFile(name string) (*struct{ Close func() error }, error) { return nil, nil }
func readData(file interface{}) ([]byte, error)                   { return nil, nil }
func processData(data []byte) error                               { return nil }
```

## 6. panic 和 recover / panic & recover

```go
package main

import (
    "fmt"
    "runtime/debug"
)

func main() {
    // panic 用于不可恢复的错误 / panic for unrecoverable errors
    // 应该很少使用，大多数情况返回 error
    
    // recover 必须在 defer 中调用 / recover must be called in defer
    fmt.Println("Starting...")
    
    safeCall(func() {
        fmt.Println("This will run")
    })
    
    safeCall(func() {
        panic("something terrible happened")
    })
    
    fmt.Println("Program continues...")
    
    // recover 示例 / recover example
    result := safeDivide(10, 0)
    fmt.Println("Safe divide result:", result)
}

// 安全调用函数 / Safe function call
func safeCall(fn func()) {
    defer func() {
        if r := recover(); r != nil {
            fmt.Printf("Recovered from panic: %v\n", r)
            fmt.Println("Stack trace:")
            debug.PrintStack()
        }
    }()
    fn()
}

// 将 panic 转换为 error / Convert panic to error
func safeDivide(a, b int) (result int) {
    defer func() {
        if r := recover(); r != nil {
            result = 0
        }
    }()
    return a / b  // 如果 b=0，会 panic
}

// 何时使用 panic / When to use panic:
// 1. 初始化时无法恢复的错误 (如配置错误)
// 2. 程序员错误 (如索引越界，这应该在开发时修复)
// 3. 真正不可能的情况

// 不应该使用 panic 的情况 / When NOT to use panic:
// 1. 普通的错误处理 (使用 error)
// 2. 可预期的失败 (如文件不存在)
// 3. 用户输入验证
```

## 7. 错误处理最佳实践 / Error Handling Best Practices

```go
package main

import (
    "errors"
    "fmt"
    "log"
)

// 1. 总是检查错误 / Always check errors
func example1() {
    // 不好 / Bad
    // result, _ := mightFail()
    
    // 好 / Good
    result, err := mightFail()
    if err != nil {
        log.Printf("operation failed: %v", err)
        return
    }
    _ = result
}

// 2. 添加上下文 / Add context
func example2() error {
    data, err := fetchData()
    if err != nil {
        // 不好 / Bad
        // return err
        
        // 好 / Good
        return fmt.Errorf("failed to fetch user data: %w", err)
    }
    _ = data
    return nil
}

// 3. 错误只处理一次 / Handle errors only once
func example3() error {
    err := doSomething()
    if err != nil {
        // 不好: 既记录又返回 / Bad: both logging and returning
        // log.Printf("error: %v", err)
        // return err
        
        // 好: 选择其一 / Good: choose one
        // 选项1: 返回错误让调用者处理
        return fmt.Errorf("example3: %w", err)
        // 或
        // 选项2: 处理错误并返回 nil 或继续
        // log.Printf("error handled: %v", err)
        // return nil
    }
    return nil
}

// 4. 使用定义的错误变量 / Use defined error variables
var (
    ErrNotFound     = errors.New("not found")
    ErrInvalidInput = errors.New("invalid input")
)

func example4(id string) error {
    if id == "" {
        return ErrInvalidInput
    }
    // ...
    return ErrNotFound
}

// 5. 错误类型要有意义 / Error types should be meaningful
type DatabaseError struct {
    Query   string
    Err     error
}

func (e DatabaseError) Error() string {
    return fmt.Sprintf("database error on query %q: %v", e.Query, e.Err)
}

func (e DatabaseError) Unwrap() error {
    return e.Err
}

// 6. 使用 defer 清理资源 / Use defer for cleanup
func example6() error {
    resource, err := acquireResource()
    if err != nil {
        return err
    }
    defer resource.Release()  // 确保资源被释放
    
    return useResource(resource)
}

// 7. 不要忽略 defer 中的错误 / Don't ignore errors in defer
func example7() (err error) {
    file, err := openFile()
    if err != nil {
        return err
    }
    defer func() {
        closeErr := file.Close()
        if err == nil {
            err = closeErr
        }
    }()
    
    return writeToFile(file)
}

// 占位函数
func mightFail() (string, error)              { return "", nil }
func fetchData() ([]byte, error)              { return nil, nil }
func doSomething() error                      { return nil }
func acquireResource() (*Resource, error)     { return nil, nil }
func useResource(r *Resource) error           { return nil }
func openFile() (*File, error)                { return nil, nil }
func writeToFile(f *File) error               { return nil }

type Resource struct{}
func (r *Resource) Release() {}

type File struct{}
func (f *File) Close() error { return nil }

func main() {
    // 示例错误处理流程
    err := example2()
    if err != nil {
        if errors.Is(err, ErrNotFound) {
            fmt.Println("Resource not found")
        } else {
            fmt.Println("Error:", err)
        }
    }
}
```

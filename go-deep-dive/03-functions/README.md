# Go 逻辑与函数 / Go Logic & Functions

## 1. 函数基础 / Function Basics

```go
package main

import "fmt"

// 基本函数声明 / Basic function declaration
func add(a int, b int) int {
    return a + b
}

// 参数类型相同时可以简写 / Shortened when types are same
func multiply(a, b int) int {
    return a * b
}

// 多返回值 / Multiple return values
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, fmt.Errorf("division by zero")
    }
    return a / b, nil
}

// 命名返回值 / Named return values
func swap(a, b int) (x, y int) {
    x = b
    y = a
    return  // 裸返回 / naked return
}

// 可变参数 / Variadic parameters
func sum(nums ...int) int {
    total := 0
    for _, n := range nums {
        total += n
    }
    return total
}

// 函数作为参数 / Function as parameter
func apply(a, b int, op func(int, int) int) int {
    return op(a, b)
}

// 函数作为返回值 / Function as return value
func makeMultiplier(factor int) func(int) int {
    return func(x int) int {
        return x * factor
    }
}

func main() {
    // 基本调用 / Basic call
    fmt.Println("add(3, 5):", add(3, 5))
    fmt.Println("multiply(4, 6):", multiply(4, 6))
    
    // 多返回值 / Multiple returns
    result, err := divide(10, 3)
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Printf("divide(10, 3): %.2f\n", result)
    }
    
    // 命名返回值 / Named returns
    x, y := swap(1, 2)
    fmt.Printf("swap(1, 2): %d, %d\n", x, y)
    
    // 可变参数 / Variadic
    fmt.Println("sum(1, 2, 3, 4, 5):", sum(1, 2, 3, 4, 5))
    
    // 展开切片 / Spread slice
    nums := []int{1, 2, 3, 4, 5}
    fmt.Println("sum(nums...):", sum(nums...))
    
    // 函数作为参数 / Function as parameter
    fmt.Println("apply with add:", apply(3, 5, add))
    fmt.Println("apply with multiply:", apply(3, 5, multiply))
    
    // 函数作为返回值 / Function as return
    double := makeMultiplier(2)
    triple := makeMultiplier(3)
    fmt.Println("double(5):", double(5))
    fmt.Println("triple(5):", triple(5))
}
```

## 2. 匿名函数与闭包 / Anonymous Functions & Closures

```go
package main

import "fmt"

func main() {
    // 匿名函数 / Anonymous function
    greet := func(name string) {
        fmt.Println("Hello,", name)
    }
    greet("World")
    
    // 立即执行的匿名函数 / IIFE (Immediately Invoked Function Expression)
    result := func(a, b int) int {
        return a + b
    }(3, 5)
    fmt.Println("IIFE result:", result)
    
    // 闭包 - 捕获外部变量 / Closure - capture outer variable
    counter := 0
    increment := func() int {
        counter++
        return counter
    }
    
    fmt.Println(increment())  // 1
    fmt.Println(increment())  // 2
    fmt.Println(increment())  // 3
    
    // 闭包陷阱 - 循环变量 / Closure pitfall - loop variable
    funcs := make([]func(), 3)
    
    // 错误方式 / Wrong way (pre Go 1.22)
    // for i := 0; i < 3; i++ {
    //     funcs[i] = func() { fmt.Println(i) }  // 都打印 3
    // }
    
    // 正确方式1: 参数传递 / Correct way 1: pass as parameter
    for i := 0; i < 3; i++ {
        funcs[i] = func(n int) func() {
            return func() { fmt.Println(n) }
        }(i)
    }
    
    // 正确方式2: 局部变量 / Correct way 2: local variable
    // for i := 0; i < 3; i++ {
    //     j := i
    //     funcs[i] = func() { fmt.Println(j) }
    // }
    
    // Go 1.22+ 修复了这个问题
    // Go 1.22+ fixes this issue
    
    for _, f := range funcs {
        f()  // 打印 0, 1, 2
    }
    
    // 使用闭包实现状态封装 / Use closure for state encapsulation
    newAccount := func(initial int) func(int) int {
        balance := initial
        return func(amount int) int {
            balance += amount
            return balance
        }
    }
    
    account := newAccount(100)
    fmt.Println("Deposit 50:", account(50))   // 150
    fmt.Println("Withdraw 30:", account(-30)) // 120
}
```

## 3. 方法 / Methods

```go
package main

import (
    "fmt"
    "math"
)

// 类型定义 / Type definition
type Point struct {
    X, Y float64
}

// 值接收者方法 / Value receiver method
func (p Point) Distance(q Point) float64 {
    return math.Sqrt(math.Pow(q.X-p.X, 2) + math.Pow(q.Y-p.Y, 2))
}

// 指针接收者方法 / Pointer receiver method
func (p *Point) Scale(factor float64) {
    p.X *= factor
    p.Y *= factor
}

func (p Point) String() string {
    return fmt.Sprintf("(%v, %v)", p.X, p.Y)
}

// 为基本类型定义方法 (需要先定义新类型)
// Define methods for basic types (need to define new type first)
type MyInt int

func (n MyInt) Double() MyInt {
    return n * 2
}

func (n *MyInt) Increment() {
    *n++
}

// 嵌入类型的方法提升 / Method promotion with embedded types
type ColoredPoint struct {
    Point  // 匿名嵌入 / Anonymous embedding
    Color string
}

func main() {
    // 值接收者 / Value receiver
    p1 := Point{3, 4}
    p2 := Point{0, 0}
    fmt.Println("Distance:", p1.Distance(p2))  // 5
    
    // 指针接收者 / Pointer receiver
    p1.Scale(2)  // Go 自动取地址
    fmt.Println("After scale:", p1)  // (6, 8)
    
    // 显式使用指针 / Explicit pointer
    (&p1).Scale(0.5)
    fmt.Println("After scale back:", p1)  // (3, 4)
    
    // 值接收者 vs 指针接收者 / Value vs Pointer receiver
    // 值接收者: 方法操作副本，不影响原值
    // 指针接收者: 方法可以修改原值，且对大结构体更高效
    
    // 自定义类型方法 / Custom type methods
    var n MyInt = 5
    fmt.Println("Double:", n.Double())  // 10
    n.Increment()
    fmt.Println("After increment:", n)  // 6
    
    // 方法提升 / Method promotion
    cp := ColoredPoint{
        Point: Point{1, 2},
        Color: "red",
    }
    fmt.Println("ColoredPoint:", cp)
    fmt.Println("Distance from origin:", cp.Distance(Point{0, 0}))  // 继承 Point 的方法
    cp.Scale(3)  // 继承的方法
    fmt.Println("After scale:", cp.Point)
    
    // 方法值 / Method value
    distanceFunc := p1.Distance
    fmt.Println("Method value:", distanceFunc(Point{0, 0}))
    
    // 方法表达式 / Method expression
    distanceExpr := Point.Distance
    fmt.Println("Method expression:", distanceExpr(p1, Point{0, 0}))
}
```

## 4. 高阶函数 / Higher-Order Functions

```go
package main

import (
    "fmt"
    "sort"
    "strings"
)

// Map 函数 / Map function
func Map[T, U any](slice []T, fn func(T) U) []U {
    result := make([]U, len(slice))
    for i, v := range slice {
        result[i] = fn(v)
    }
    return result
}

// Filter 函数 / Filter function
func Filter[T any](slice []T, predicate func(T) bool) []T {
    result := make([]T, 0)
    for _, v := range slice {
        if predicate(v) {
            result = append(result, v)
        }
    }
    return result
}

// Reduce 函数 / Reduce function
func Reduce[T, U any](slice []T, initial U, fn func(U, T) U) U {
    result := initial
    for _, v := range slice {
        result = fn(result, v)
    }
    return result
}

// ForEach 函数 / ForEach function
func ForEach[T any](slice []T, fn func(T)) {
    for _, v := range slice {
        fn(v)
    }
}

// Find 函数 / Find function
func Find[T any](slice []T, predicate func(T) bool) (T, bool) {
    for _, v := range slice {
        if predicate(v) {
            return v, true
        }
    }
    var zero T
    return zero, false
}

// Any 函数 / Any function
func Any[T any](slice []T, predicate func(T) bool) bool {
    for _, v := range slice {
        if predicate(v) {
            return true
        }
    }
    return false
}

// All 函数 / All function
func All[T any](slice []T, predicate func(T) bool) bool {
    for _, v := range slice {
        if !predicate(v) {
            return false
        }
    }
    return true
}

func main() {
    nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    
    // Map: 将每个元素乘以2
    doubled := Map(nums, func(n int) int { return n * 2 })
    fmt.Println("Doubled:", doubled)
    
    // Filter: 筛选偶数
    evens := Filter(nums, func(n int) bool { return n%2 == 0 })
    fmt.Println("Evens:", evens)
    
    // Reduce: 求和
    sum := Reduce(nums, 0, func(acc, n int) int { return acc + n })
    fmt.Println("Sum:", sum)
    
    // ForEach: 打印每个元素
    fmt.Print("ForEach: ")
    ForEach(nums, func(n int) { fmt.Print(n, " ") })
    fmt.Println()
    
    // Find: 找第一个大于5的数
    found, ok := Find(nums, func(n int) bool { return n > 5 })
    fmt.Println("Find >5:", found, ok)
    
    // Any: 是否有大于5的数
    fmt.Println("Any >5:", Any(nums, func(n int) bool { return n > 5 }))
    
    // All: 是否全部大于0
    fmt.Println("All >0:", All(nums, func(n int) bool { return n > 0 }))
    
    // 链式操作 / Chained operations
    result := Reduce(
        Filter(
            Map(nums, func(n int) int { return n * n }),  // 平方
            func(n int) bool { return n > 10 },           // 大于10
        ),
        0,
        func(acc, n int) int { return acc + n },  // 求和
    )
    fmt.Println("Chained result:", result)
    
    // 使用标准库的排序 / Using standard library sort
    strs := []string{"banana", "apple", "cherry", "date"}
    sort.Slice(strs, func(i, j int) bool {
        return strings.ToLower(strs[i]) < strings.ToLower(strs[j])
    })
    fmt.Println("Sorted:", strs)
}
```

## 5. 函数选项模式 / Functional Options Pattern

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

## 6. 递归与尾递归 / Recursion & Tail Recursion

```go
package main

import "fmt"

// 基本递归 - 阶乘 / Basic recursion - factorial
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n-1)
}

// 尾递归版本 / Tail recursive version
// 注意: Go 编译器不保证尾调用优化
// Note: Go compiler doesn't guarantee tail call optimization
func factorialTail(n, acc int) int {
    if n <= 1 {
        return acc
    }
    return factorialTail(n-1, n*acc)
}

// 斐波那契数列 - 普通递归 / Fibonacci - normal recursion
func fibonacci(n int) int {
    if n <= 1 {
        return n
    }
    return fibonacci(n-1) + fibonacci(n-2)
}

// 斐波那契 - 记忆化递归 / Fibonacci - memoized
func fibonacciMemo(n int, memo map[int]int) int {
    if n <= 1 {
        return n
    }
    if v, ok := memo[n]; ok {
        return v
    }
    memo[n] = fibonacciMemo(n-1, memo) + fibonacciMemo(n-2, memo)
    return memo[n]
}

// 斐波那契 - 迭代版本 (推荐) / Fibonacci - iterative (recommended)
func fibonacciIter(n int) int {
    if n <= 1 {
        return n
    }
    a, b := 0, 1
    for i := 2; i <= n; i++ {
        a, b = b, a+b
    }
    return b
}

// 树遍历递归示例 / Tree traversal recursion example
type TreeNode struct {
    Value int
    Left  *TreeNode
    Right *TreeNode
}

func inorderTraversal(node *TreeNode, result *[]int) {
    if node == nil {
        return
    }
    inorderTraversal(node.Left, result)
    *result = append(*result, node.Value)
    inorderTraversal(node.Right, result)
}

// 快速排序 / Quick sort
func quickSort(arr []int) []int {
    if len(arr) <= 1 {
        return arr
    }
    
    pivot := arr[0]
    var left, right []int
    
    for _, v := range arr[1:] {
        if v < pivot {
            left = append(left, v)
        } else {
            right = append(right, v)
        }
    }
    
    result := append(quickSort(left), pivot)
    return append(result, quickSort(right)...)
}

func main() {
    // 阶乘 / Factorial
    fmt.Println("factorial(5):", factorial(5))
    fmt.Println("factorialTail(5, 1):", factorialTail(5, 1))
    
    // 斐波那契 / Fibonacci
    fmt.Println("fibonacci(10):", fibonacci(10))
    
    memo := make(map[int]int)
    fmt.Println("fibonacciMemo(10):", fibonacciMemo(10, memo))
    fmt.Println("fibonacciIter(10):", fibonacciIter(10))
    
    // 树遍历 / Tree traversal
    tree := &TreeNode{
        Value: 4,
        Left: &TreeNode{
            Value: 2,
            Left:  &TreeNode{Value: 1},
            Right: &TreeNode{Value: 3},
        },
        Right: &TreeNode{
            Value: 6,
            Left:  &TreeNode{Value: 5},
            Right: &TreeNode{Value: 7},
        },
    }
    var result []int
    inorderTraversal(tree, &result)
    fmt.Println("Inorder:", result)
    
    // 快速排序 / Quick sort
    arr := []int{64, 34, 25, 12, 22, 11, 90}
    fmt.Println("QuickSort:", quickSort(arr))
}
```

## 7. defer、panic、recover 模式 / defer, panic, recover Patterns

```go
package main

import (
    "fmt"
    "io"
    "os"
    "sync"
)

// 模式1: 资源清理 / Pattern 1: Resource cleanup
func copyFile(src, dst string) (err error) {
    srcFile, err := os.Open(src)
    if err != nil {
        return err
    }
    defer srcFile.Close()
    
    dstFile, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer func() {
        closeErr := dstFile.Close()
        if err == nil {
            err = closeErr
        }
    }()
    
    _, err = io.Copy(dstFile, srcFile)
    return err
}

// 模式2: 解锁互斥锁 / Pattern 2: Unlock mutex
func safeOperation(mu *sync.Mutex, fn func()) {
    mu.Lock()
    defer mu.Unlock()
    fn()
}

// 模式3: 计时 / Pattern 3: Timing
func timeTrack(name string) func() {
    start := time.Now()
    return func() {
        fmt.Printf("%s took %v\n", name, time.Since(start))
    }
}

func slowOperation() {
    defer timeTrack("slowOperation")()
    time.Sleep(100 * time.Millisecond)
}

// 模式4: 修改命名返回值 / Pattern 4: Modify named return
func divideWithRecovery(a, b int) (result int, err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("panic recovered: %v", r)
        }
    }()
    
    result = a / b  // 如果 b=0 会 panic
    return result, nil
}

// 模式5: 日志记录 / Pattern 5: Logging
func loggedOperation(name string) func() {
    fmt.Printf("Starting %s\n", name)
    return func() {
        fmt.Printf("Finished %s\n", name)
    }
}

// 模式6: 错误处理增强 / Pattern 6: Error handling enhancement
func processWithContext(ctx string, fn func() error) (err error) {
    defer func() {
        if err != nil {
            err = fmt.Errorf("%s: %w", ctx, err)
        }
    }()
    return fn()
}

// defer 执行顺序 / defer execution order
func deferOrder() {
    defer fmt.Println("First defer")
    defer fmt.Println("Second defer")
    defer fmt.Println("Third defer")
    fmt.Println("Function body")
    // 输出顺序: Function body -> Third defer -> Second defer -> First defer
}

func main() {
    // 模式4 示例
    result, err := divideWithRecovery(10, 0)
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("Result:", result)
    }
    
    // defer 执行顺序
    deferOrder()
    
    // 使用 loggedOperation
    defer loggedOperation("main")()
    fmt.Println("Main body")
}

// 注意事项 / Notes:
// 1. defer 的参数在声明时求值，而非执行时
// 2. 多个 defer 按 LIFO (后进先出) 顺序执行
// 3. defer 可以访问和修改命名返回值
// 4. recover 只能在 defer 函数中调用才有效
```

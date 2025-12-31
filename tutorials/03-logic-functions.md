# Go 逻辑与函数详解

## 1. 函数基础

### 1.1 函数声明

```go
package main

import "fmt"

// 基本函数
func greet() {
    fmt.Println("Hello!")
}

// 带参数的函数
func greetName(name string) {
    fmt.Printf("Hello, %s!\n", name)
}

// 带返回值的函数
func add(a, b int) int {
    return a + b
}

// 多返回值
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, fmt.Errorf("division by zero")
    }
    return a / b, nil
}

// 命名返回值
func split(sum int) (x, y int) {
    x = sum * 4 / 9
    y = sum - x
    return // 裸返回
}

// 可变参数
func sum(nums ...int) int {
    total := 0
    for _, n := range nums {
        total += n
    }
    return total
}

// 多类型参数
func printInfo(name string, age int, scores ...float64) {
    fmt.Printf("Name: %s, Age: %d\n", name, age)
    for i, s := range scores {
        fmt.Printf("Score %d: %.2f\n", i+1, s)
    }
}

func main() {
    greet()
    greetName("Alice")
    
    result := add(3, 5)
    fmt.Println("Add:", result)
    
    quotient, err := divide(10, 3)
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Printf("Quotient: %.2f\n", quotient)
    }
    
    x, y := split(17)
    fmt.Println("Split:", x, y)
    
    fmt.Println("Sum:", sum(1, 2, 3, 4, 5))
    
    // 展开切片
    nums := []int{1, 2, 3, 4, 5}
    fmt.Println("Sum slice:", sum(nums...))
    
    printInfo("Bob", 25, 90.5, 85.0, 92.5)
}
```

### 1.2 函数作为值

```go
package main

import (
    "fmt"
    "math"
)

// 函数类型
type MathFunc func(float64) float64

// 返回函数
func getOperation(op string) MathFunc {
    switch op {
    case "square":
        return func(x float64) float64 { return x * x }
    case "sqrt":
        return math.Sqrt
    case "double":
        return func(x float64) float64 { return x * 2 }
    default:
        return func(x float64) float64 { return x }
    }
}

// 高阶函数
func apply(f MathFunc, x float64) float64 {
    return f(x)
}

// 函数组合
func compose(f, g MathFunc) MathFunc {
    return func(x float64) float64 {
        return f(g(x))
    }
}

func main() {
    // 函数赋值给变量
    add := func(a, b int) int {
        return a + b
    }
    fmt.Println("Add:", add(3, 5))
    
    // 获取不同操作
    square := getOperation("square")
    sqrt := getOperation("sqrt")
    
    fmt.Println("Square(5):", square(5))
    fmt.Println("Sqrt(16):", sqrt(16))
    
    // 函数作为参数
    fmt.Println("Apply square to 4:", apply(square, 4))
    
    // 函数组合
    sqrtOfSquare := compose(sqrt, square)
    fmt.Println("sqrt(square(4)):", sqrtOfSquare(4)) // sqrt(16) = 4
    
    // 函数切片
    operations := []MathFunc{
        func(x float64) float64 { return x + 1 },
        func(x float64) float64 { return x * 2 },
        func(x float64) float64 { return x * x },
    }
    
    x := 2.0
    for _, op := range operations {
        x = op(x)
    }
    fmt.Println("连续操作结果:", x) // ((2+1)*2)^2 = 36
    
    // 立即调用函数
    result := func(a, b int) int {
        return a + b
    }(3, 5)
    fmt.Println("立即调用:", result)
}
```

### 1.3 闭包

```go
package main

import "fmt"

// 计数器闭包
func counter() func() int {
    count := 0
    return func() int {
        count++
        return count
    }
}

// 累加器闭包
func accumulator(initial int) func(int) int {
    sum := initial
    return func(n int) int {
        sum += n
        return sum
    }
}

// 带状态的乘法器
func multiplier(factor int) func(int) int {
    return func(x int) int {
        return x * factor
    }
}

// 斐波那契闭包
func fibonacci() func() int {
    a, b := 0, 1
    return func() int {
        result := a
        a, b = b, a+b
        return result
    }
}

// 闭包陷阱演示
func closureTrap() {
    // 错误示例
    fmt.Println("闭包陷阱 - 错误:")
    funcs := make([]func(), 3)
    for i := 0; i < 3; i++ {
        funcs[i] = func() {
            fmt.Println(i) // 捕获的是变量 i，不是值
        }
    }
    for _, f := range funcs {
        f() // 都打印 3
    }
    
    // 正确示例
    fmt.Println("闭包陷阱 - 正确:")
    funcs2 := make([]func(), 3)
    for i := 0; i < 3; i++ {
        i := i // 创建新变量
        funcs2[i] = func() {
            fmt.Println(i)
        }
    }
    for _, f := range funcs2 {
        f() // 打印 0, 1, 2
    }
    
    // Go 1.22+ 自动修复
    fmt.Println("Go 1.22+ 自动修复:")
    funcs3 := make([]func(), 3)
    for i := range 3 {
        funcs3[i] = func() {
            fmt.Println(i) // Go 1.22+ 每次迭代 i 是新变量
        }
    }
    for _, f := range funcs3 {
        f()
    }
}

func main() {
    // 计数器
    c1 := counter()
    c2 := counter()
    
    fmt.Println("c1:", c1(), c1(), c1()) // 1, 2, 3
    fmt.Println("c2:", c2(), c2())       // 1, 2
    
    // 累加器
    acc := accumulator(10)
    fmt.Println("累加:", acc(5))   // 15
    fmt.Println("累加:", acc(3))   // 18
    fmt.Println("累加:", acc(2))   // 20
    
    // 乘法器
    double := multiplier(2)
    triple := multiplier(3)
    
    fmt.Println("Double 5:", double(5))
    fmt.Println("Triple 5:", triple(5))
    
    // 斐波那契
    fib := fibonacci()
    fmt.Print("Fibonacci: ")
    for i := 0; i < 10; i++ {
        fmt.Print(fib(), " ")
    }
    fmt.Println()
    
    // 闭包陷阱
    closureTrap()
}
```

## 2. 递归

```go
package main

import "fmt"

// 阶乘
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n-1)
}

// 尾递归优化版本
func factorialTail(n, acc int) int {
    if n <= 1 {
        return acc
    }
    return factorialTail(n-1, n*acc)
}

// 斐波那契
func fib(n int) int {
    if n <= 1 {
        return n
    }
    return fib(n-1) + fib(n-2)
}

// 带记忆化的斐波那契
func fibMemo() func(int) int {
    cache := make(map[int]int)
    var f func(int) int
    f = func(n int) int {
        if n <= 1 {
            return n
        }
        if v, ok := cache[n]; ok {
            return v
        }
        result := f(n-1) + f(n-2)
        cache[n] = result
        return result
    }
    return f
}

// 二分查找
func binarySearch(arr []int, target, low, high int) int {
    if low > high {
        return -1
    }
    mid := (low + high) / 2
    if arr[mid] == target {
        return mid
    } else if arr[mid] > target {
        return binarySearch(arr, target, low, mid-1)
    } else {
        return binarySearch(arr, target, mid+1, high)
    }
}

// 快速排序
func quickSort(arr []int) []int {
    if len(arr) <= 1 {
        return arr
    }
    
    pivot := arr[0]
    var left, right []int
    
    for _, v := range arr[1:] {
        if v <= pivot {
            left = append(left, v)
        } else {
            right = append(right, v)
        }
    }
    
    result := append(quickSort(left), pivot)
    return append(result, quickSort(right)...)
}

// 归并排序
func mergeSort(arr []int) []int {
    if len(arr) <= 1 {
        return arr
    }
    
    mid := len(arr) / 2
    left := mergeSort(arr[:mid])
    right := mergeSort(arr[mid:])
    
    return merge(left, right)
}

func merge(left, right []int) []int {
    result := make([]int, 0, len(left)+len(right))
    i, j := 0, 0
    
    for i < len(left) && j < len(right) {
        if left[i] <= right[j] {
            result = append(result, left[i])
            i++
        } else {
            result = append(result, right[j])
            j++
        }
    }
    
    result = append(result, left[i:]...)
    result = append(result, right[j:]...)
    
    return result
}

// 树结构遍历
type TreeNode struct {
    Value int
    Left  *TreeNode
    Right *TreeNode
}

func (t *TreeNode) PreOrder() []int {
    if t == nil {
        return nil
    }
    result := []int{t.Value}
    result = append(result, t.Left.PreOrder()...)
    result = append(result, t.Right.PreOrder()...)
    return result
}

func (t *TreeNode) InOrder() []int {
    if t == nil {
        return nil
    }
    result := t.Left.InOrder()
    result = append(result, t.Value)
    result = append(result, t.Right.InOrder()...)
    return result
}

func (t *TreeNode) PostOrder() []int {
    if t == nil {
        return nil
    }
    result := t.Left.PostOrder()
    result = append(result, t.Right.PostOrder()...)
    result = append(result, t.Value)
    return result
}

func main() {
    // 阶乘
    fmt.Println("5! =", factorial(5))
    fmt.Println("5! (tail) =", factorialTail(5, 1))
    
    // 斐波那契
    fmt.Println("fib(10) =", fib(10))
    
    // 带记忆化的斐波那契
    fibM := fibMemo()
    fmt.Println("fibMemo(40) =", fibM(40)) // 快速
    
    // 二分查找
    arr := []int{1, 3, 5, 7, 9, 11, 13, 15}
    fmt.Println("Binary search 7:", binarySearch(arr, 7, 0, len(arr)-1))
    fmt.Println("Binary search 8:", binarySearch(arr, 8, 0, len(arr)-1))
    
    // 快速排序
    unsorted := []int{3, 1, 4, 1, 5, 9, 2, 6}
    fmt.Println("QuickSort:", quickSort(unsorted))
    
    // 归并排序
    unsorted2 := []int{3, 1, 4, 1, 5, 9, 2, 6}
    fmt.Println("MergeSort:", mergeSort(unsorted2))
    
    // 树遍历
    root := &TreeNode{
        Value: 1,
        Left: &TreeNode{
            Value: 2,
            Left:  &TreeNode{Value: 4},
            Right: &TreeNode{Value: 5},
        },
        Right: &TreeNode{
            Value: 3,
            Left:  &TreeNode{Value: 6},
            Right: &TreeNode{Value: 7},
        },
    }
    
    fmt.Println("PreOrder:", root.PreOrder())   // [1 2 4 5 3 6 7]
    fmt.Println("InOrder:", root.InOrder())     // [4 2 5 1 6 3 7]
    fmt.Println("PostOrder:", root.PostOrder()) // [4 5 2 6 7 3 1]
}
```

## 3. 方法

### 3.1 值接收器 vs 指针接收器

```go
package main

import (
    "fmt"
    "math"
)

type Point struct {
    X, Y float64
}

// 值接收器 - 不修改原值
func (p Point) Distance(q Point) float64 {
    dx := p.X - q.X
    dy := p.Y - q.Y
    return math.Sqrt(dx*dx + dy*dy)
}

// 值接收器 - 返回新值
func (p Point) Add(q Point) Point {
    return Point{p.X + q.X, p.Y + q.Y}
}

// 指针接收器 - 修改原值
func (p *Point) Move(dx, dy float64) {
    p.X += dx
    p.Y += dy
}

// 指针接收器 - 避免复制大结构体
func (p *Point) ScaleBy(factor float64) {
    p.X *= factor
    p.Y *= factor
}

// 方法集规则演示
type Counter struct {
    Value int
}

func (c Counter) Get() int {
    return c.Value
}

func (c *Counter) Increment() {
    c.Value++
}

func (c *Counter) Add(n int) {
    c.Value += n
}

func main() {
    // Point 方法
    p1 := Point{0, 0}
    p2 := Point{3, 4}
    
    fmt.Println("Distance:", p1.Distance(p2))
    
    // 值方法不修改原值
    p3 := p1.Add(p2)
    fmt.Println("p1 after Add:", p1) // 不变
    fmt.Println("p3:", p3)
    
    // 指针方法修改原值
    p1.Move(1, 1)
    fmt.Println("p1 after Move:", p1) // 改变
    
    p1.ScaleBy(2)
    fmt.Println("p1 after Scale:", p1)
    
    // Counter 方法
    c := Counter{Value: 0}
    
    // 值类型可以调用值方法
    fmt.Println("Get:", c.Get())
    
    // 值类型调用指针方法 - Go 自动取地址
    c.Increment()
    fmt.Println("After Increment:", c.Value)
    
    // 指针类型
    cp := &Counter{Value: 10}
    cp.Increment()
    cp.Add(5)
    fmt.Println("Counter pointer:", cp.Value)
    
    // 方法值
    distance := p1.Distance
    fmt.Println("Method value:", distance(Point{0, 0}))
    
    // 方法表达式
    distanceExpr := Point.Distance
    fmt.Println("Method expression:", distanceExpr(p1, Point{0, 0}))
}
```

### 3.2 嵌入类型方法

```go
package main

import "fmt"

// 基础类型
type Animal struct {
    Name string
}

func (a Animal) Speak() string {
    return "..."
}

func (a Animal) Move() string {
    return "moving"
}

// 嵌入类型
type Dog struct {
    Animal
    Breed string
}

// 重写方法
func (d Dog) Speak() string {
    return "Woof!"
}

// 添加新方法
func (d Dog) Fetch() string {
    return fmt.Sprintf("%s is fetching", d.Name)
}

// 多级嵌入
type ServiceDog struct {
    Dog
    Service string
}

func (s ServiceDog) Work() string {
    return fmt.Sprintf("%s is doing %s", s.Name, s.Service)
}

// 接口实现
type Speaker interface {
    Speak() string
}

type Mover interface {
    Move() string
}

func makeSpeak(s Speaker) {
    fmt.Println(s.Speak())
}

func main() {
    // Dog
    dog := Dog{
        Animal: Animal{Name: "Buddy"},
        Breed:  "Golden Retriever",
    }
    
    // 直接访问嵌入字段
    fmt.Println("Name:", dog.Name)
    fmt.Println("Breed:", dog.Breed)
    
    // 调用重写的方法
    fmt.Println("Speak:", dog.Speak()) // Woof!
    
    // 调用继承的方法
    fmt.Println("Move:", dog.Move())
    
    // 调用父类原方法
    fmt.Println("Animal.Speak:", dog.Animal.Speak()) // ...
    
    // 新方法
    fmt.Println(dog.Fetch())
    
    // ServiceDog
    sd := ServiceDog{
        Dog: Dog{
            Animal: Animal{Name: "Max"},
            Breed:  "German Shepherd",
        },
        Service: "guide dog",
    }
    
    fmt.Println("ServiceDog Name:", sd.Name)
    fmt.Println("ServiceDog Speak:", sd.Speak())
    fmt.Println("ServiceDog Work:", sd.Work())
    
    // 接口
    makeSpeak(dog)
    makeSpeak(sd)
    
    var speaker Speaker = dog
    fmt.Println("Interface:", speaker.Speak())
}
```

## 4. 控制流进阶

### 4.1 defer 详解

```go
package main

import (
    "fmt"
    "os"
    "sync"
)

func main() {
    // defer 执行顺序 (LIFO)
    fmt.Println("=== defer 顺序 ===")
    defer fmt.Println("1")
    defer fmt.Println("2")
    defer fmt.Println("3")
    fmt.Println("main")
    
    // defer 参数立即求值
    fmt.Println("\n=== defer 参数 ===")
    x := 10
    defer fmt.Println("defer x =", x) // 10
    x = 20
    fmt.Println("x =", x)
    
    // defer 闭包捕获变量
    fmt.Println("\n=== defer 闭包 ===")
    y := 10
    defer func() {
        fmt.Println("defer closure y =", y) // 20
    }()
    y = 20
    
    // defer 修改命名返回值
    fmt.Println("\n=== defer 修改返回值 ===")
    result := deferReturn()
    fmt.Println("result:", result) // 2
    
    // defer 资源管理
    fmt.Println("\n=== defer 资源管理 ===")
    manageFile()
    
    // defer 锁管理
    fmt.Println("\n=== defer 锁管理 ===")
    manageLock()
    
    // defer 循环陷阱
    fmt.Println("\n=== defer 循环 ===")
    deferLoop()
}

func deferReturn() (n int) {
    defer func() {
        n++
    }()
    return 1 // 实际返回 2
}

func manageFile() {
    f, err := os.Create("/tmp/test.txt")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer f.Close() // 确保关闭
    
    f.WriteString("Hello, World!")
    fmt.Println("File written")
}

func manageLock() {
    var mu sync.Mutex
    
    mu.Lock()
    defer mu.Unlock()
    
    fmt.Println("Critical section")
}

func deferLoop() {
    // 错误: 文件句柄会累积
    // for i := 0; i < 5; i++ {
    //     f, _ := os.Open("file")
    //     defer f.Close() // 循环结束后才执行
    // }
    
    // 正确: 使用匿名函数
    for i := 0; i < 3; i++ {
        func(n int) {
            fmt.Println("Processing", n)
            defer fmt.Println("Done", n)
        }(i)
    }
}
```

### 4.2 panic 和 recover

```go
package main

import (
    "fmt"
    "runtime/debug"
)

func main() {
    // 基本 panic/recover
    fmt.Println("=== 基本 panic/recover ===")
    safeDivide(10, 2)
    safeDivide(10, 0)
    fmt.Println("程序继续运行")
    
    // recover 返回值
    fmt.Println("\n=== recover 返回值 ===")
    result := recoverReturn()
    fmt.Println("result:", result)
    
    // panic 链
    fmt.Println("\n=== panic 链 ===")
    panicChain()
    
    // 自定义错误类型
    fmt.Println("\n=== 自定义 panic ===")
    customPanic()
    
    // 打印堆栈
    fmt.Println("\n=== 堆栈跟踪 ===")
    stackTrace()
}

func safeDivide(a, b int) {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered:", r)
        }
    }()
    
    result := a / b
    fmt.Println("Result:", result)
}

func recoverReturn() (result int) {
    defer func() {
        if r := recover(); r != nil {
            result = -1
        }
    }()
    
    panic("error")
    return 0
}

func panicChain() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Outer recover:", r)
        }
    }()
    
    func() {
        defer func() {
            if r := recover(); r != nil {
                fmt.Println("Inner recover:", r)
                panic("re-panic: " + r.(string)) // 重新 panic
            }
        }()
        
        panic("original panic")
    }()
}

type MyError struct {
    Code    int
    Message string
}

func (e MyError) Error() string {
    return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
}

func customPanic() {
    defer func() {
        if r := recover(); r != nil {
            switch e := r.(type) {
            case MyError:
                fmt.Printf("MyError: Code=%d, Message=%s\n", e.Code, e.Message)
            case error:
                fmt.Println("Error:", e.Error())
            case string:
                fmt.Println("String:", e)
            default:
                fmt.Println("Unknown:", r)
            }
        }
    }()
    
    panic(MyError{Code: 500, Message: "Internal Server Error"})
}

func stackTrace() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Panic:", r)
            fmt.Println("Stack trace:")
            debug.PrintStack()
        }
    }()
    
    level1()
}

func level1() {
    level2()
}

func level2() {
    level3()
}

func level3() {
    panic("deep panic")
}
```

## 5. 函数式编程模式

```go
package main

import (
    "fmt"
    "sort"
    "strings"
)

// Map 函数
func Map[T, U any](slice []T, f func(T) U) []U {
    result := make([]U, len(slice))
    for i, v := range slice {
        result[i] = f(v)
    }
    return result
}

// Filter 函数
func Filter[T any](slice []T, f func(T) bool) []T {
    result := make([]T, 0)
    for _, v := range slice {
        if f(v) {
            result = append(result, v)
        }
    }
    return result
}

// Reduce 函数
func Reduce[T, U any](slice []T, init U, f func(U, T) U) U {
    result := init
    for _, v := range slice {
        result = f(result, v)
    }
    return result
}

// ForEach 函数
func ForEach[T any](slice []T, f func(T)) {
    for _, v := range slice {
        f(v)
    }
}

// Find 函数
func Find[T any](slice []T, f func(T) bool) (T, bool) {
    for _, v := range slice {
        if f(v) {
            return v, true
        }
    }
    var zero T
    return zero, false
}

// Any 函数
func Any[T any](slice []T, f func(T) bool) bool {
    for _, v := range slice {
        if f(v) {
            return true
        }
    }
    return false
}

// All 函数
func All[T any](slice []T, f func(T) bool) bool {
    for _, v := range slice {
        if !f(v) {
            return false
        }
    }
    return true
}

// GroupBy 函数
func GroupBy[T any, K comparable](slice []T, keyFunc func(T) K) map[K][]T {
    result := make(map[K][]T)
    for _, v := range slice {
        key := keyFunc(v)
        result[key] = append(result[key], v)
    }
    return result
}

// Partition 函数
func Partition[T any](slice []T, f func(T) bool) ([]T, []T) {
    trueSlice := make([]T, 0)
    falseSlice := make([]T, 0)
    for _, v := range slice {
        if f(v) {
            trueSlice = append(trueSlice, v)
        } else {
            falseSlice = append(falseSlice, v)
        }
    }
    return trueSlice, falseSlice
}

// 柯里化
func Curry2[A, B, C any](f func(A, B) C) func(A) func(B) C {
    return func(a A) func(B) C {
        return func(b B) C {
            return f(a, b)
        }
    }
}

// 管道操作
type Pipeline[T any] struct {
    data []T
}

func NewPipeline[T any](data []T) *Pipeline[T] {
    return &Pipeline[T]{data: data}
}

func (p *Pipeline[T]) Filter(f func(T) bool) *Pipeline[T] {
    p.data = Filter(p.data, f)
    return p
}

func (p *Pipeline[T]) Result() []T {
    return p.data
}

func main() {
    // Map
    nums := []int{1, 2, 3, 4, 5}
    squares := Map(nums, func(n int) int { return n * n })
    fmt.Println("Squares:", squares)
    
    // Filter
    evens := Filter(nums, func(n int) bool { return n%2 == 0 })
    fmt.Println("Evens:", evens)
    
    // Reduce
    sum := Reduce(nums, 0, func(acc, n int) int { return acc + n })
    fmt.Println("Sum:", sum)
    
    product := Reduce(nums, 1, func(acc, n int) int { return acc * n })
    fmt.Println("Product:", product)
    
    // 链式操作
    words := []string{"hello", "world", "go", "programming"}
    result := Map(
        Filter(words, func(s string) bool { return len(s) > 3 }),
        strings.ToUpper,
    )
    fmt.Println("Filtered & Upper:", result)
    
    // Find
    if found, ok := Find(nums, func(n int) bool { return n > 3 }); ok {
        fmt.Println("Found:", found)
    }
    
    // Any & All
    fmt.Println("Any > 3:", Any(nums, func(n int) bool { return n > 3 }))
    fmt.Println("All > 0:", All(nums, func(n int) bool { return n > 0 }))
    
    // GroupBy
    type Person struct {
        Name string
        Age  int
    }
    
    people := []Person{
        {"Alice", 25},
        {"Bob", 30},
        {"Charlie", 25},
        {"David", 30},
    }
    
    byAge := GroupBy(people, func(p Person) int { return p.Age })
    fmt.Println("By Age:", byAge)
    
    // Partition
    adults, minors := Partition([]int{15, 18, 20, 12, 25, 16}, func(age int) bool {
        return age >= 18
    })
    fmt.Println("Adults:", adults)
    fmt.Println("Minors:", minors)
    
    // 柯里化
    add := func(a, b int) int { return a + b }
    curriedAdd := Curry2(add)
    add5 := curriedAdd(5)
    fmt.Println("5 + 3 =", add5(3))
    
    // Pipeline
    numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    pipelineResult := NewPipeline(numbers).
        Filter(func(n int) bool { return n%2 == 0 }).
        Result()
    fmt.Println("Pipeline:", pipelineResult)
    
    // 排序示例
    sort.Slice(people, func(i, j int) bool {
        return people[i].Age < people[j].Age
    })
    fmt.Println("Sorted people:", people)
}
```

## 6. 选项模式

```go
package main

import (
    "fmt"
    "time"
)

// Server 配置
type Server struct {
    Host         string
    Port         int
    Timeout      time.Duration
    MaxConns     int
    TLS          bool
    ReadTimeout  time.Duration
    WriteTimeout time.Duration
}

// Option 函数类型
type Option func(*Server)

// 选项函数
func WithHost(host string) Option {
    return func(s *Server) {
        s.Host = host
    }
}

func WithPort(port int) Option {
    return func(s *Server) {
        s.Port = port
    }
}

func WithTimeout(timeout time.Duration) Option {
    return func(s *Server) {
        s.Timeout = timeout
    }
}

func WithMaxConns(maxConns int) Option {
    return func(s *Server) {
        s.MaxConns = maxConns
    }
}

func WithTLS(enabled bool) Option {
    return func(s *Server) {
        s.TLS = enabled
    }
}

func WithReadTimeout(timeout time.Duration) Option {
    return func(s *Server) {
        s.ReadTimeout = timeout
    }
}

func WithWriteTimeout(timeout time.Duration) Option {
    return func(s *Server) {
        s.WriteTimeout = timeout
    }
}

// 默认配置
func defaultServer() *Server {
    return &Server{
        Host:         "localhost",
        Port:         8080,
        Timeout:      30 * time.Second,
        MaxConns:     100,
        TLS:          false,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
    }
}

// 构造函数
func NewServer(opts ...Option) *Server {
    server := defaultServer()
    for _, opt := range opts {
        opt(server)
    }
    return server
}

// 另一种选项模式: Builder 模式
type ServerBuilder struct {
    server *Server
}

func NewServerBuilder() *ServerBuilder {
    return &ServerBuilder{server: defaultServer()}
}

func (b *ServerBuilder) Host(host string) *ServerBuilder {
    b.server.Host = host
    return b
}

func (b *ServerBuilder) Port(port int) *ServerBuilder {
    b.server.Port = port
    return b
}

func (b *ServerBuilder) Timeout(timeout time.Duration) *ServerBuilder {
    b.server.Timeout = timeout
    return b
}

func (b *ServerBuilder) MaxConns(maxConns int) *ServerBuilder {
    b.server.MaxConns = maxConns
    return b
}

func (b *ServerBuilder) TLS(enabled bool) *ServerBuilder {
    b.server.TLS = enabled
    return b
}

func (b *ServerBuilder) Build() *Server {
    return b.server
}

func main() {
    // 使用选项模式
    server1 := NewServer()
    fmt.Printf("Default: %+v\n", server1)
    
    server2 := NewServer(
        WithHost("0.0.0.0"),
        WithPort(443),
        WithTLS(true),
        WithTimeout(60*time.Second),
    )
    fmt.Printf("Custom: %+v\n", server2)
    
    // 使用 Builder 模式
    server3 := NewServerBuilder().
        Host("api.example.com").
        Port(8443).
        TLS(true).
        MaxConns(1000).
        Build()
    fmt.Printf("Builder: %+v\n", server3)
}
```

## 7. 依赖注入

```go
package main

import (
    "fmt"
    "time"
)

// 接口定义
type Logger interface {
    Log(message string)
}

type Database interface {
    Query(sql string) []string
    Execute(sql string) error
}

type Cache interface {
    Get(key string) (string, bool)
    Set(key string, value string, ttl time.Duration)
}

// 实现
type ConsoleLogger struct{}

func (l *ConsoleLogger) Log(message string) {
    fmt.Println("[LOG]", message)
}

type FileLogger struct {
    Filename string
}

func (l *FileLogger) Log(message string) {
    fmt.Printf("[FILE:%s] %s\n", l.Filename, message)
}

type MockDatabase struct{}

func (d *MockDatabase) Query(sql string) []string {
    return []string{"result1", "result2"}
}

func (d *MockDatabase) Execute(sql string) error {
    fmt.Println("Executing:", sql)
    return nil
}

type MemoryCache struct {
    data map[string]string
}

func NewMemoryCache() *MemoryCache {
    return &MemoryCache{data: make(map[string]string)}
}

func (c *MemoryCache) Get(key string) (string, bool) {
    v, ok := c.data[key]
    return v, ok
}

func (c *MemoryCache) Set(key string, value string, ttl time.Duration) {
    c.data[key] = value
}

// 服务 - 依赖注入
type UserService struct {
    logger Logger
    db     Database
    cache  Cache
}

func NewUserService(logger Logger, db Database, cache Cache) *UserService {
    return &UserService{
        logger: logger,
        db:     db,
        cache:  cache,
    }
}

func (s *UserService) GetUser(id string) string {
    // 先查缓存
    if user, ok := s.cache.Get(id); ok {
        s.logger.Log("Cache hit: " + id)
        return user
    }
    
    // 查数据库
    s.logger.Log("Cache miss, querying DB: " + id)
    results := s.db.Query("SELECT * FROM users WHERE id = " + id)
    
    if len(results) > 0 {
        s.cache.Set(id, results[0], 5*time.Minute)
        return results[0]
    }
    
    return ""
}

// 容器
type Container struct {
    logger Logger
    db     Database
    cache  Cache
}

func NewContainer() *Container {
    return &Container{
        logger: &ConsoleLogger{},
        db:     &MockDatabase{},
        cache:  NewMemoryCache(),
    }
}

func (c *Container) UserService() *UserService {
    return NewUserService(c.logger, c.db, c.cache)
}

func main() {
    // 手动依赖注入
    logger := &ConsoleLogger{}
    db := &MockDatabase{}
    cache := NewMemoryCache()
    
    userService := NewUserService(logger, db, cache)
    user := userService.GetUser("123")
    fmt.Println("User:", user)
    
    // 第二次查询会命中缓存
    user = userService.GetUser("123")
    fmt.Println("User:", user)
    
    // 使用容器
    fmt.Println("\n=== Using Container ===")
    container := NewContainer()
    service := container.UserService()
    service.GetUser("456")
    
    // 测试时可以使用不同的实现
    fmt.Println("\n=== Testing with FileLogger ===")
    testLogger := &FileLogger{Filename: "test.log"}
    testService := NewUserService(testLogger, db, cache)
    testService.GetUser("789")
}
```

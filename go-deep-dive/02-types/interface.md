# Go 接口 / Go Interface

## 1. 接口基础 / Interface Basics

```go
package main

import (
    "fmt"
    "math"
)

// 接口定义了一组方法签名
// Interface defines a set of method signatures
type Shape interface {
    Area() float64
    Perimeter() float64
}

// 实现接口 (隐式实现)
// Implement interface (implicit implementation)
type Rectangle struct {
    Width, Height float64
}

func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
    return 2 * (r.Width + r.Height)
}

type Circle struct {
    Radius float64
}

func (c Circle) Area() float64 {
    return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
    return 2 * math.Pi * c.Radius
}

// 使用接口作为参数 / Use interface as parameter
func PrintShapeInfo(s Shape) {
    fmt.Printf("Area: %.2f, Perimeter: %.2f\n", s.Area(), s.Perimeter())
}

func main() {
    rect := Rectangle{Width: 10, Height: 5}
    circle := Circle{Radius: 7}
    
    // 多态性 / Polymorphism
    PrintShapeInfo(rect)
    PrintShapeInfo(circle)
    
    // 接口类型的变量 / Interface type variable
    var shape Shape
    shape = rect
    fmt.Println("Rectangle area:", shape.Area())
    
    shape = circle
    fmt.Println("Circle area:", shape.Area())
    
    // 接口切片 / Interface slice
    shapes := []Shape{rect, circle}
    for _, s := range shapes {
        PrintShapeInfo(s)
    }
}
```

## 2. 空接口 (any) / Empty Interface

```go
package main

import "fmt"

func main() {
    // 空接口可以存储任何类型的值
    // Empty interface can store any value
    var anything interface{}  // 或 any (Go 1.18+)
    
    anything = 42
    fmt.Printf("int: %v, type: %T\n", anything, anything)
    
    anything = "hello"
    fmt.Printf("string: %v, type: %T\n", anything, anything)
    
    anything = []int{1, 2, 3}
    fmt.Printf("slice: %v, type: %T\n", anything, anything)
    
    // 通用函数 / Generic function (pre-generics)
    printAnything := func(v interface{}) {
        fmt.Printf("Value: %v, Type: %T\n", v, v)
    }
    
    printAnything(100)
    printAnything("world")
    printAnything(true)
    
    // 空接口切片 / Empty interface slice
    mixedSlice := []interface{}{1, "two", 3.0, true}
    for _, v := range mixedSlice {
        fmt.Println(v)
    }
    
    // map 的值为空接口 / Map with empty interface values
    data := map[string]interface{}{
        "name":   "Alice",
        "age":    30,
        "active": true,
    }
    fmt.Println("data:", data)
}
```

## 3. 类型断言 / Type Assertion

```go
package main

import "fmt"

func main() {
    var i interface{} = "hello"
    
    // 基本类型断言 / Basic type assertion
    s := i.(string)
    fmt.Println("String:", s)
    
    // 安全的类型断言 (带 ok) / Safe type assertion (with ok)
    s, ok := i.(string)
    if ok {
        fmt.Println("String:", s)
    }
    
    // 断言失败时不会 panic
    n, ok := i.(int)
    if !ok {
        fmt.Println("Not an int, got:", n)  // n 是零值
    }
    
    // 不安全的断言会 panic / Unsafe assertion panics
    // n := i.(int)  // panic: interface conversion
    
    // 类型 switch / Type switch
    checkType := func(v interface{}) {
        switch x := v.(type) {
        case int:
            fmt.Printf("int: %d\n", x)
        case string:
            fmt.Printf("string: %s\n", x)
        case bool:
            fmt.Printf("bool: %t\n", x)
        case []int:
            fmt.Printf("[]int: %v\n", x)
        case nil:
            fmt.Println("nil")
        default:
            fmt.Printf("unknown type: %T\n", x)
        }
    }
    
    checkType(42)
    checkType("hello")
    checkType(true)
    checkType([]int{1, 2, 3})
    checkType(nil)
    checkType(3.14)  // unknown type
}
```

## 4. 接口组合 / Interface Composition

```go
package main

import "fmt"

// 小接口 / Small interfaces
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

type Closer interface {
    Close() error
}

// 组合接口 / Composed interface
type ReadWriter interface {
    Reader
    Writer
}

type ReadWriteCloser interface {
    Reader
    Writer
    Closer
}

// 也可以添加额外方法 / Can also add extra methods
type File interface {
    ReadWriter
    Closer
    Stat() (FileInfo, error)
}

// 实现组合接口 / Implement composed interface
type Buffer struct {
    data []byte
}

func (b *Buffer) Read(p []byte) (n int, err error) {
    n = copy(p, b.data)
    b.data = b.data[n:]
    return n, nil
}

func (b *Buffer) Write(p []byte) (n int, err error) {
    b.data = append(b.data, p...)
    return len(p), nil
}

func main() {
    buf := &Buffer{}
    
    // 可以赋值给任何它实现的接口
    // Can assign to any interface it implements
    var r Reader = buf
    var w Writer = buf
    var rw ReadWriter = buf
    
    // 写入数据
    w.Write([]byte("Hello"))
    
    // 读取数据
    p := make([]byte, 5)
    r.Read(p)
    fmt.Println("Read:", string(p))
    
    // 使用 ReadWriter
    rw.Write([]byte("World"))
    p2 := make([]byte, 5)
    rw.Read(p2)
    fmt.Println("Read:", string(p2))
}
```

## 5. 接口值的内部结构 / Interface Value Internals

```go
package main

import (
    "fmt"
    "unsafe"
)

type Speaker interface {
    Speak() string
}

type Dog struct {
    Name string
}

func (d Dog) Speak() string {
    return "Woof! I'm " + d.Name
}

func main() {
    // 接口值包含两个指针: (type, value)
    // Interface value contains two pointers: (type, value)
    
    var s Speaker
    fmt.Printf("nil interface: (%v, %T)\n", s, s)
    
    var d *Dog = nil
    s = d  // s 现在不是 nil！
    fmt.Printf("interface with nil value: (%v, %T)\n", s, s)
    fmt.Println("s == nil:", s == nil)  // false!
    
    // 这是一个常见陷阱 / This is a common pitfall
    // 即使值是 nil，接口本身不是 nil，因为它有类型信息
    
    // 正确的 nil 检查 / Correct nil check
    if s != nil {
        // 需要检查内部值是否为 nil
        // Need to check if internal value is nil
        switch v := s.(type) {
        case *Dog:
            if v == nil {
                fmt.Println("Dog is nil")
            }
        }
    }
    
    // 设置实际值 / Set actual value
    s = Dog{Name: "Buddy"}
    fmt.Printf("interface with value: (%v, %T)\n", s, s)
    fmt.Println(s.Speak())
}
```

## 6. 常见标准库接口 / Common Standard Library Interfaces

```go
package main

import (
    "fmt"
    "io"
    "sort"
    "strings"
)

// fmt.Stringer - 自定义字符串表示 / Custom string representation
type Person struct {
    Name string
    Age  int
}

func (p Person) String() string {
    return fmt.Sprintf("%s (%d years old)", p.Name, p.Age)
}

// error 接口 / error interface
type MyError struct {
    Code    int
    Message string
}

func (e MyError) Error() string {
    return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
}

// sort.Interface / sort.Interface
type ByAge []Person

func (a ByAge) Len() int           { return len(a) }
func (a ByAge) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByAge) Less(i, j int) bool { return a[i].Age < a[j].Age }

// io.Reader 和 io.Writer / io.Reader and io.Writer
type UpperWriter struct {
    w io.Writer
}

func (uw UpperWriter) Write(p []byte) (n int, err error) {
    return uw.w.Write([]byte(strings.ToUpper(string(p))))
}

func main() {
    // Stringer
    p := Person{Name: "Alice", Age: 30}
    fmt.Println(p)  // Alice (30 years old)
    
    // error
    err := MyError{Code: 404, Message: "Not Found"}
    fmt.Println(err)  // Error 404: Not Found
    
    // sort.Interface
    people := []Person{
        {"Charlie", 35},
        {"Alice", 25},
        {"Bob", 30},
    }
    sort.Sort(ByAge(people))
    fmt.Println("Sorted:", people)
    
    // io.Writer
    var builder strings.Builder
    uw := UpperWriter{w: &builder}
    uw.Write([]byte("hello world"))
    fmt.Println("Upper:", builder.String())
}
```

## 7. 接口最佳实践 / Interface Best Practices

```go
package main

import "fmt"

// 1. 保持接口小而专注 / Keep interfaces small and focused
// 好 / Good
type Reader interface {
    Read(p []byte) (n int, err error)
}

// 不好 / Bad - 接口太大
// type FileHandler interface {
//     Read() error
//     Write() error
//     Delete() error
//     Copy() error
//     Move() error
//     ... 更多方法
// }

// 2. 接口应该由使用者定义，而不是实现者
// Interfaces should be defined by consumers, not implementors

// 服务实现 / Service implementation
type UserService struct{}

func (s *UserService) GetUser(id int) (User, error) {
    return User{ID: id, Name: "Alice"}, nil
}

func (s *UserService) SaveUser(u User) error {
    return nil
}

// 消费者定义需要的接口 / Consumer defines needed interface
type UserGetter interface {
    GetUser(id int) (User, error)
}

// 只需要 GetUser 功能的处理器
type Handler struct {
    users UserGetter  // 依赖接口，而非具体类型
}

// 3. 接受接口，返回具体类型
// Accept interfaces, return concrete types
func ProcessReader(r Reader) error {
    // 接受任何实现 Reader 的类型
    return nil
}

func NewUserService() *UserService {  // 返回具体类型
    return &UserService{}
}

// 4. 使用接口进行依赖注入和测试
// Use interfaces for dependency injection and testing
type EmailSender interface {
    Send(to, subject, body string) error
}

type NotificationService struct {
    emailer EmailSender
}

func (ns *NotificationService) Notify(user User, message string) error {
    return ns.emailer.Send(user.Email, "Notification", message)
}

// 测试时可以使用 mock / Can use mock for testing
type MockEmailSender struct {
    LastTo      string
    LastSubject string
    LastBody    string
}

func (m *MockEmailSender) Send(to, subject, body string) error {
    m.LastTo = to
    m.LastSubject = subject
    m.LastBody = body
    return nil
}

type User struct {
    ID    int
    Name  string
    Email string
}

func main() {
    // 依赖注入示例
    mock := &MockEmailSender{}
    service := &NotificationService{emailer: mock}
    
    user := User{ID: 1, Name: "Alice", Email: "alice@example.com"}
    service.Notify(user, "Hello!")
    
    fmt.Println("Sent to:", mock.LastTo)
    fmt.Println("Subject:", mock.LastSubject)
    fmt.Println("Body:", mock.LastBody)
}
```

## 8. 泛型与接口约束 (Go 1.18+) / Generics & Interface Constraints

```go
package main

import (
    "fmt"
    "golang.org/x/exp/constraints"
)

// 类型约束接口 / Type constraint interface
type Number interface {
    int | int32 | int64 | float32 | float64
}

// 使用约束的泛型函数 / Generic function with constraint
func Sum[T Number](numbers []T) T {
    var sum T
    for _, n := range numbers {
        sum += n
    }
    return sum
}

// 使用 constraints 包 / Using constraints package
func Min[T constraints.Ordered](a, b T) T {
    if a < b {
        return a
    }
    return b
}

// 接口约束可以包含方法 / Interface constraints can include methods
type Stringer interface {
    String() string
}

func Stringify[T Stringer](items []T) []string {
    result := make([]string, len(items))
    for i, item := range items {
        result[i] = item.String()
    }
    return result
}

// 结合类型和方法约束 / Combine type and method constraints
type OrderedStringer interface {
    constraints.Ordered
    String() string
}

// 类型集 / Type sets
type SignedInteger interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64
}

// ~ 表示底层类型 / ~ means underlying type
type MyInt int

func Double[T SignedInteger](x T) T {
    return x * 2
}

func main() {
    // Sum
    ints := []int{1, 2, 3, 4, 5}
    floats := []float64{1.1, 2.2, 3.3}
    
    fmt.Println("Sum of ints:", Sum(ints))
    fmt.Println("Sum of floats:", Sum(floats))
    
    // Min
    fmt.Println("Min(3, 5):", Min(3, 5))
    fmt.Println("Min(\"apple\", \"banana\"):", Min("apple", "banana"))
    
    // Double with custom type
    var myNum MyInt = 10
    fmt.Println("Double:", Double(myNum))
}
```

# 接口 / Interface

## 1. 接口基础 / Interface Basics

```go
package main

import (
    "fmt"
    "math"
)

// 接口定义方法签名集合 / Interface defines method signatures
type Shape interface {
    Area() float64
    Perimeter() float64
}

// 隐式实现接口 / Implicit implementation
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

func PrintShapeInfo(s Shape) {
    fmt.Printf("Area: %.2f, Perimeter: %.2f\n", s.Area(), s.Perimeter())
}

func main() {
    rect := Rectangle{Width: 10, Height: 5}
    circle := Circle{Radius: 7}
    
    PrintShapeInfo(rect)
    PrintShapeInfo(circle)
    
    // 接口切片 / Interface slice
    shapes := []Shape{rect, circle}
    for _, s := range shapes {
        PrintShapeInfo(s)
    }
}
```

## 2. 空接口与类型断言 / Empty Interface & Type Assertion

```go
func main() {
    // 空接口可存储任何类型 / Empty interface stores any type
    var anything interface{}  // 或 any (Go 1.18+)
    
    anything = 42
    anything = "hello"
    anything = []int{1, 2, 3}
    
    // 类型断言 / Type assertion
    var i interface{} = "hello"
    s, ok := i.(string)
    if ok {
        fmt.Println("String:", s)
    }
    
    // 类型 switch / Type switch
    switch v := i.(type) {
    case int:
        fmt.Printf("int: %d\n", v)
    case string:
        fmt.Printf("string: %s\n", v)
    default:
        fmt.Printf("unknown: %T\n", v)
    }
}
```

## 3. 接口组合 / Interface Composition

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

type ReadWriter interface {
    Reader
    Writer
}
```

## 4. 常见标准库接口 / Common Standard Library Interfaces

```go
// fmt.Stringer
type Person struct {
    Name string
    Age  int
}

func (p Person) String() string {
    return fmt.Sprintf("%s (%d)", p.Name, p.Age)
}

// error
type MyError struct {
    Code    int
    Message string
}

func (e MyError) Error() string {
    return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
}

// sort.Interface
type ByAge []Person

func (a ByAge) Len() int           { return len(a) }
func (a ByAge) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByAge) Less(i, j int) bool { return a[i].Age < a[j].Age }
```

## 5. 泛型接口约束 / Generic Interface Constraints (Go 1.18+)

```go
type Number interface {
    int | int32 | int64 | float32 | float64
}

func Sum[T Number](numbers []T) T {
    var sum T
    for _, n := range numbers {
        sum += n
    }
    return sum
}

// ~ 表示底层类型 / ~ means underlying type
type SignedInteger interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64
}
```

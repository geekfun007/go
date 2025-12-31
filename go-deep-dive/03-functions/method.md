# 方法 / Methods

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

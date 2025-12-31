# 函数基础 / Function Basics

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

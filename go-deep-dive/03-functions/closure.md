# 匿名函数与闭包 / Anonymous Functions & Closures

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

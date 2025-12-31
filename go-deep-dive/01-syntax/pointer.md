# 指针 / Pointers

```go
package main

import "fmt"

func main() {
    // 指针存储变量的内存地址
    // Pointers store memory addresses of variables
    
    x := 10
    var p *int = &x  // p 是指向 x 的指针
    
    fmt.Println("x =", x)    // 10
    fmt.Println("p =", p)    // 内存地址
    fmt.Println("*p =", *p)  // 10 (解引用)
    
    // 通过指针修改值
    // Modify value through pointer
    *p = 20
    fmt.Println("x =", x)  // 20
    
    // 指针的零值是 nil
    // Zero value of pointer is nil
    var nilPtr *int
    fmt.Println("nilPtr =", nilPtr)  // <nil>
    
    // new 函数创建指针
    // new function creates pointer
    ptr := new(int)  // 分配内存，返回指针
    *ptr = 100
    fmt.Println("*ptr =", *ptr)
}

// 函数参数传递
// Function parameter passing
func increment(x int) {
    x++ // 不会影响原值 / Won't affect original
}

func incrementPtr(x *int) {
    (*x)++ // 会影响原值 / Will affect original
}
```

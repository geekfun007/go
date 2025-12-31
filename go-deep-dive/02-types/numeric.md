# 数值类型 / Numeric Types

## 1. 整数 (int) / Integer

```go
package main

import (
    "fmt"
    "math"
)

func main() {
    // 有符号整数 / Signed integers
    var i8 int8 = 127                    // -128 to 127
    var i16 int16 = 32767                // -32768 to 32767
    var i32 int32 = 2147483647           // -2^31 to 2^31-1
    var i64 int64 = 9223372036854775807  // -2^63 to 2^63-1
    var i int = 100                      // 平台相关 (32或64位)
    
    // 无符号整数 / Unsigned integers
    var u8 uint8 = 255      // 0 to 255 (byte 的别名)
    var u16 uint16 = 65535  // 0 to 65535
    var u32 uint32 = 4294967295
    var u64 uint64 = 18446744073709551615
    var u uint = 100        // 平台相关
    
    // 特殊类型 / Special types
    var b byte = 255        // uint8 的别名
    var r rune = '中'       // int32 的别名，表示 Unicode 码点
    
    fmt.Println(i8, i16, i32, i64, i)
    fmt.Println(u8, u16, u32, u64, u)
    fmt.Println(b, r)
    
    // 整数运算 / Integer operations
    a, b := 17, 5
    fmt.Println("加法 / Addition:", a+b)       // 22
    fmt.Println("减法 / Subtraction:", a-b)    // 12
    fmt.Println("乘法 / Multiplication:", a*b) // 85
    fmt.Println("除法 / Division:", a/b)       // 3 (整数除法)
    fmt.Println("取模 / Modulo:", a%b)         // 2
    
    // 位运算 / Bit operations
    x, y := 12, 25  // 12 = 1100, 25 = 11001
    fmt.Println("AND:", x&y)   // 8  (01000)
    fmt.Println("OR:", x|y)    // 29 (11101)
    fmt.Println("XOR:", x^y)   // 21 (10101)
    fmt.Println("左移:", x<<2) // 48 (110000)
    fmt.Println("右移:", x>>2) // 3  (11)
    
    // 类型转换 / Type conversion
    var num int = 100
    var numFloat float64 = float64(num)
    var numInt64 int64 = int64(num)
    fmt.Println(numFloat, numInt64)
}
```

## 2. 浮点数 (float) / Float

```go
package main

import (
    "fmt"
    "math"
)

func main() {
    // 浮点类型 / Float types
    var f32 float32 = 3.14159265358979
    var f64 float64 = 3.14159265358979323846
    
    fmt.Printf("float32: %.10f\n", f32)  // 精度较低
    fmt.Printf("float64: %.15f\n", f64)  // 精度较高
    
    // 特殊值 / Special values
    inf := math.Inf(1)      // 正无穷
    negInf := math.Inf(-1)  // 负无穷
    nan := math.NaN()       // 非数字
    
    fmt.Println("Inf:", inf)
    fmt.Println("-Inf:", negInf)
    fmt.Println("NaN:", nan)
    fmt.Println("IsNaN:", math.IsNaN(nan))
    fmt.Println("IsInf:", math.IsInf(inf, 1))
    
    // 浮点运算 / Float operations
    a, b := 10.5, 3.2
    fmt.Printf("加法: %.2f\n", a+b)
    fmt.Printf("减法: %.2f\n", a-b)
    fmt.Printf("乘法: %.2f\n", a*b)
    fmt.Printf("除法: %.2f\n", a/b)
    
    // 数学函数 / Math functions
    fmt.Printf("Sqrt(16): %.2f\n", math.Sqrt(16))
    fmt.Printf("Pow(2, 10): %.0f\n", math.Pow(2, 10))
    fmt.Printf("Abs(-5.5): %.2f\n", math.Abs(-5.5))
    fmt.Printf("Ceil(3.2): %.0f\n", math.Ceil(3.2))
    fmt.Printf("Floor(3.8): %.0f\n", math.Floor(3.8))
    fmt.Printf("Round(3.5): %.0f\n", math.Round(3.5))
    fmt.Printf("Sin(π/2): %.2f\n", math.Sin(math.Pi/2))
    fmt.Printf("Log(e): %.2f\n", math.Log(math.E))
    fmt.Printf("Log10(100): %.2f\n", math.Log10(100))
    
    // 浮点数比较 (注意精度问题)
    // Float comparison (beware of precision issues)
    epsilon := 1e-9
    x, y := 0.1+0.2, 0.3
    if math.Abs(x-y) < epsilon {
        fmt.Println("x 约等于 y")
    }
}
```

## 3. 布尔值 (bool) / Boolean

```go
package main

import "fmt"

func main() {
    // 布尔类型 / Boolean type
    var t bool = true
    var f bool = false
    
    fmt.Println("true:", t)
    fmt.Println("false:", f)
    
    // 布尔运算 / Boolean operations
    fmt.Println("AND (&&):", true && false)  // false
    fmt.Println("OR (||):", true || false)   // true
    fmt.Println("NOT (!):", !true)           // false
    
    // 比较运算符返回布尔值
    // Comparison operators return bool
    a, b := 5, 10
    fmt.Println("a == b:", a == b)  // false
    fmt.Println("a != b:", a != b)  // true
    fmt.Println("a < b:", a < b)    // true
    fmt.Println("a > b:", a > b)    // false
    fmt.Println("a <= b:", a <= b)  // true
    fmt.Println("a >= b:", a >= b)  // false
    
    // 短路求值 / Short-circuit evaluation
    // 如果第一个条件已经确定结果，第二个条件不会执行
    x := 5
    if x > 0 || expensiveCheck() {
        fmt.Println("短路: expensiveCheck 不会被调用")
    }
    
    // 布尔值不能与整数互换
    // Boolean cannot be converted to/from int
    // var i int = true  // 编译错误
}

func expensiveCheck() bool {
    fmt.Println("expensiveCheck called")
    return true
}
```

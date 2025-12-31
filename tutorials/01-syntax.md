# Go 基础语法详解

## 1. 程序结构

### 1.1 包 (Package)

```go
// 每个 Go 文件必须以 package 声明开始
package main

// 导入单个包
import "fmt"

// 导入多个包
import (
    "fmt"
    "os"
    "strings"
)

// 别名导入
import (
    f "fmt"                    // 使用别名 f
    . "strings"                // 点导入，可以直接使用函数名
    _ "database/sql"           // 空白导入，只执行 init()
)

// main 函数是程序入口
func main() {
    f.Println("Hello")         // 使用别名
    fmt.Println(ToUpper("hi")) // 点导入可直接使用
}
```

### 1.2 init 函数

```go
package main

import "fmt"

// init 函数在 main 之前自动执行
// 一个包可以有多个 init 函数
func init() {
    fmt.Println("init 1")
}

func init() {
    fmt.Println("init 2")
}

func main() {
    fmt.Println("main")
}

// 输出:
// init 1
// init 2
// main
```

### 1.3 可见性规则

```go
package example

// 大写字母开头：公开 (可被其他包访问)
var PublicVar = "public"
func PublicFunc() {}
type PublicStruct struct {
    PublicField  string  // 公开字段
    privateField string  // 私有字段
}

// 小写字母开头：私有 (只能在当前包访问)
var privateVar = "private"
func privateFunc() {}
type privateStruct struct{}
```

## 2. 变量与常量

### 2.1 变量声明方式

```go
package main

import "fmt"

// 包级变量
var globalVar = "global"

func main() {
    // 方式1: 完整声明
    var a int = 10
    
    // 方式2: 类型推断
    var b = 20
    
    // 方式3: 短变量声明 (只能在函数内)
    c := 30
    
    // 方式4: 多变量声明
    var d, e, f int = 1, 2, 3
    g, h, i := 4, "hello", true
    
    // 方式5: 变量块
    var (
        name   string = "Alice"
        age    int    = 25
        active bool   = true
    )
    
    // 零值初始化
    var (
        zeroInt    int     // 0
        zeroFloat  float64 // 0.0
        zeroBool   bool    // false
        zeroString string  // ""
        zeroPtr    *int    // nil
        zeroSlice  []int   // nil
        zeroMap    map[string]int // nil
        zeroChan   chan int       // nil
        zeroFunc   func()         // nil
        zeroIface  interface{}    // nil
    )
    
    fmt.Println(a, b, c, d, e, f, g, h, i)
    fmt.Println(name, age, active)
    fmt.Println(zeroInt, zeroFloat, zeroBool, zeroString)
    fmt.Println(zeroPtr, zeroSlice, zeroMap, zeroChan, zeroFunc, zeroIface)
}
```

### 2.2 常量与 iota

```go
package main

import "fmt"

// 普通常量
const Pi = 3.14159
const (
    MaxInt = 1<<31 - 1
    MinInt = -1 << 31
)

// 类型常量
const TypedConst int = 100

// iota: 常量计数器
const (
    A = iota  // 0
    B         // 1
    C         // 2
)

// iota 重置
const (
    X = iota  // 0
    Y         // 1
)

// iota 高级用法
const (
    _  = iota             // 0 (跳过)
    KB = 1 << (10 * iota) // 1 << 10 = 1024
    MB                    // 1 << 20
    GB                    // 1 << 30
    TB                    // 1 << 40
)

// 位掩码
const (
    Read    = 1 << iota  // 1
    Write                // 2
    Execute              // 4
)

// 跳值
const (
    P = iota  // 0
    Q         // 1
    _         // 2 (跳过)
    R         // 3
)

// 同行 iota 相同
const (
    M, N = iota, iota + 10  // 0, 10
    O, P2 = iota, iota + 10 // 1, 11
)

func main() {
    fmt.Printf("A=%d, B=%d, C=%d\n", A, B, C)
    fmt.Printf("KB=%d, MB=%d, GB=%d\n", KB, MB, GB)
    fmt.Printf("Read=%d, Write=%d, Execute=%d\n", Read, Write, Execute)
    
    // 位运算使用
    permissions := Read | Write
    fmt.Printf("权限: %b\n", permissions)
    fmt.Printf("可读: %v\n", permissions&Read != 0)
    fmt.Printf("可执行: %v\n", permissions&Execute != 0)
}
```

## 3. 控制流

### 3.1 if 语句

```go
package main

import "fmt"

func main() {
    x := 10
    
    // 基本 if
    if x > 5 {
        fmt.Println("x > 5")
    }
    
    // if-else
    if x > 20 {
        fmt.Println("x > 20")
    } else {
        fmt.Println("x <= 20")
    }
    
    // if-else if-else
    if x > 20 {
        fmt.Println("大")
    } else if x > 10 {
        fmt.Println("中")
    } else {
        fmt.Println("小")
    }
    
    // if 带初始化语句
    if y := x * 2; y > 15 {
        fmt.Println("y > 15:", y)
    }
    // y 的作用域仅在 if 块内
    
    // 常见模式: 错误处理
    if err := doSomething(); err != nil {
        fmt.Println("错误:", err)
    }
}

func doSomething() error {
    return nil
}
```

### 3.2 for 循环

```go
package main

import "fmt"

func main() {
    // 基本 for
    for i := 0; i < 5; i++ {
        fmt.Print(i, " ")
    }
    fmt.Println()
    
    // while 风格
    j := 0
    for j < 5 {
        fmt.Print(j, " ")
        j++
    }
    fmt.Println()
    
    // 无限循环
    k := 0
    for {
        if k >= 5 {
            break
        }
        fmt.Print(k, " ")
        k++
    }
    fmt.Println()
    
    // range 遍历切片
    nums := []int{10, 20, 30}
    for i, v := range nums {
        fmt.Printf("索引%d: %d\n", i, v)
    }
    
    // range 遍历字符串
    for i, r := range "Hello" {
        fmt.Printf("%d: %c\n", i, r)
    }
    
    // range 遍历 map
    m := map[string]int{"a": 1, "b": 2}
    for k, v := range m {
        fmt.Printf("%s: %d\n", k, v)
    }
    
    // continue
    for i := 0; i < 10; i++ {
        if i%2 == 0 {
            continue  // 跳过偶数
        }
        fmt.Print(i, " ")
    }
    fmt.Println()
    
    // 标签和跳转
outer:
    for i := 0; i < 3; i++ {
        for j := 0; j < 3; j++ {
            if j == 2 {
                break outer  // 跳出外层循环
            }
            fmt.Printf("(%d,%d) ", i, j)
        }
    }
    fmt.Println()
    
    // Go 1.22+: range 整数
    for i := range 5 {
        fmt.Print(i, " ")
    }
    fmt.Println()
}
```

### 3.3 switch 语句

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

func main() {
    // 基本 switch
    day := 3
    switch day {
    case 1:
        fmt.Println("Monday")
    case 2:
        fmt.Println("Tuesday")
    case 3:
        fmt.Println("Wednesday")
    default:
        fmt.Println("Other day")
    }
    
    // 多值 case
    switch day {
    case 1, 2, 3, 4, 5:
        fmt.Println("工作日")
    case 6, 7:
        fmt.Println("周末")
    }
    
    // switch 带初始化
    switch os := runtime.GOOS; os {
    case "darwin":
        fmt.Println("macOS")
    case "linux":
        fmt.Println("Linux")
    case "windows":
        fmt.Println("Windows")
    default:
        fmt.Println(os)
    }
    
    // 无条件 switch (相当于 if-else)
    hour := time.Now().Hour()
    switch {
    case hour < 12:
        fmt.Println("上午")
    case hour < 18:
        fmt.Println("下午")
    default:
        fmt.Println("晚上")
    }
    
    // fallthrough
    n := 1
    switch n {
    case 1:
        fmt.Println("一")
        fallthrough  // 继续执行下一个 case
    case 2:
        fmt.Println("二")
    case 3:
        fmt.Println("三")
    }
    
    // 类型 switch
    var i interface{} = "hello"
    switch v := i.(type) {
    case int:
        fmt.Printf("整数: %d\n", v)
    case string:
        fmt.Printf("字符串: %s\n", v)
    case bool:
        fmt.Printf("布尔: %v\n", v)
    default:
        fmt.Printf("未知类型: %T\n", v)
    }
}
```

### 3.4 defer, panic, recover

```go
package main

import "fmt"

func main() {
    // defer: 延迟执行，LIFO 顺序
    fmt.Println("=== defer 基础 ===")
    defer fmt.Println("defer 1")
    defer fmt.Println("defer 2")
    defer fmt.Println("defer 3")
    fmt.Println("main")
    // 输出: main, defer 3, defer 2, defer 1
    
    // defer 常用于资源清理
    fmt.Println("\n=== defer 资源清理 ===")
    processFile()
    
    // defer 参数立即求值
    fmt.Println("\n=== defer 参数求值 ===")
    x := 10
    defer fmt.Println("defer x =", x)  // x=10 (立即求值)
    x = 20
    fmt.Println("current x =", x)
    
    // 使用闭包捕获变量
    y := 10
    defer func() {
        fmt.Println("defer y =", y)  // y=20 (使用最终值)
    }()
    y = 20
    
    // panic 和 recover
    fmt.Println("\n=== panic/recover ===")
    safeOperation()
    fmt.Println("程序继续运行")
    
    // defer 可以修改返回值
    fmt.Println("\n=== defer 修改返回值 ===")
    result := modifyReturn()
    fmt.Println("result =", result)  // 2
}

func processFile() {
    fmt.Println("打开文件")
    defer fmt.Println("关闭文件")  // 确保文件被关闭
    fmt.Println("处理文件")
}

func safeOperation() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered:", r)
        }
    }()
    
    fmt.Println("开始操作")
    panic("出错了!")
    fmt.Println("这行不会执行")
}

func modifyReturn() (result int) {
    defer func() {
        result++  // 修改命名返回值
    }()
    return 1  // 实际返回 2
}
```

## 4. 指针

```go
package main

import "fmt"

func main() {
    // 声明指针
    var p *int
    fmt.Println("nil指针:", p)
    
    // 获取地址
    x := 10
    p = &x
    fmt.Println("x的地址:", p)
    fmt.Println("p指向的值:", *p)
    
    // 通过指针修改值
    *p = 20
    fmt.Println("修改后x:", x)
    
    // new 分配内存
    p2 := new(int)
    *p2 = 100
    fmt.Println("new分配:", *p2)
    
    // 指针和结构体
    type Person struct {
        Name string
        Age  int
    }
    
    person := &Person{Name: "Alice", Age: 25}
    fmt.Println("person:", person)
    fmt.Println("Name:", person.Name)  // 自动解引用
    fmt.Println("Name:", (*person).Name)  // 显式解引用
    
    // 指针作为函数参数
    a := 10
    increment(&a)
    fmt.Println("increment后:", a)
    
    // 值传递 vs 指针传递
    b := 100
    modifyValue(b)
    fmt.Println("值传递后:", b)  // 不变
    
    modifyPointer(&b)
    fmt.Println("指针传递后:", b)  // 改变
    
    // Go 没有指针运算
    // p++  // 编译错误
}

func increment(n *int) {
    *n++
}

func modifyValue(n int) {
    n = 999
}

func modifyPointer(n *int) {
    *n = 999
}
```

## 5. 注释与文档

```go
// Package example 提供示例功能
// 这是包文档
package example

import "fmt"

// MaxSize 是最大尺寸常量
const MaxSize = 100

// User 表示用户信息
// 
// 示例:
//
//	user := User{Name: "Alice", Age: 25}
//	fmt.Println(user.Greet())
type User struct {
    Name string // 用户名
    Age  int    // 年龄
}

// Greet 返回问候语
//
// 参数: 无
//
// 返回值: 问候字符串
func (u User) Greet() string {
    return fmt.Sprintf("Hello, I'm %s", u.Name)
}

// Add 计算两个整数的和
//
// Deprecated: 请使用 Sum 代替
func Add(a, b int) int {
    return a + b
}

/*
多行注释
可以用于大段说明
*/
func multiLineComment() {}
```

运行 `go doc` 查看文档:
```bash
go doc example
go doc example.User
go doc example.User.Greet
```

## 6. 代码组织最佳实践

### 6.1 命名规范

```go
// 包名: 小写，简短
package httputil

// 变量名: 驼峰式
var userName string
var maxRetryCount int

// 常量: 驼峰式或全大写
const MaxConnections = 100
const DEFAULT_TIMEOUT = 30

// 函数名: 驼峰式
func getUserByID(id int) {}
func ParseJSON(data []byte) {}

// 接口名: 通常以 er 结尾
type Reader interface{}
type Writer interface{}
type Closer interface{}

// 缩写词保持一致大小写
var userID int        // 或 userId
var httpClient Client // 或 HTTPClient
var xmlParser Parser  // 或 XMLParser
```

### 6.2 文件组织

```go
// 一个文件通常包含:
// 1. 包声明
// 2. import 声明
// 3. 常量声明
// 4. 变量声明
// 5. 类型声明
// 6. 函数声明

package example

import (
    "fmt"
    "strings"
)

const maxSize = 100

var defaultConfig = Config{}

type Config struct {
    Name string
}

func NewConfig() *Config {
    return &Config{}
}

func (c *Config) Validate() error {
    return nil
}
```

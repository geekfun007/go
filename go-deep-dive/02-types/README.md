# Go 类型与方法 / Go Types & Methods

## 1. 数值类型 / Numeric Types

### 1.1 整数 (int) / Integer

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

### 1.2 浮点数 (float) / Float

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

### 1.3 布尔值 (bool) / Boolean

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
    // If first condition determines result, second won't execute
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

## 2. 字符串 (string) / String

```go
package main

import (
    "fmt"
    "strings"
    "unicode/utf8"
)

func main() {
    // 字符串是不可变的 UTF-8 编码字节序列
    // Strings are immutable UTF-8 encoded byte sequences
    
    // 字符串声明 / String declaration
    s1 := "Hello, 世界"
    s2 := `原始字符串字面量
可以包含换行
不会解释转义字符 \n`
    
    fmt.Println(s1)
    fmt.Println(s2)
    
    // 字符串长度 / String length
    fmt.Println("字节长度 / Byte length:", len(s1))              // 13
    fmt.Println("字符长度 / Rune length:", utf8.RuneCountInString(s1))  // 9
    
    // 字符串索引 (返回字节) / String indexing (returns byte)
    fmt.Printf("s1[0] = %c (byte)\n", s1[0])  // H
    
    // 遍历字符串 / Iterate string
    // 按字节遍历 / By byte
    for i := 0; i < len(s1); i++ {
        fmt.Printf("%d: %c ", i, s1[i])
    }
    fmt.Println()
    
    // 按 rune 遍历 (推荐用于中文等) / By rune (recommended for non-ASCII)
    for i, r := range s1 {
        fmt.Printf("%d: %c ", i, r)
    }
    fmt.Println()
    
    // 字符串操作 / String operations
    str := "  Hello, Go World!  "
    
    // 包含检查 / Contains check
    fmt.Println("Contains 'Go':", strings.Contains(str, "Go"))
    fmt.Println("HasPrefix '  He':", strings.HasPrefix(str, "  He"))
    fmt.Println("HasSuffix '!  ':", strings.HasSuffix(str, "!  "))
    
    // 查找 / Find
    fmt.Println("Index 'Go':", strings.Index(str, "Go"))
    fmt.Println("LastIndex 'o':", strings.LastIndex(str, "o"))
    fmt.Println("Count 'o':", strings.Count(str, "o"))
    
    // 转换 / Transform
    fmt.Println("ToUpper:", strings.ToUpper(str))
    fmt.Println("ToLower:", strings.ToLower(str))
    fmt.Println("TrimSpace:", strings.TrimSpace(str))
    fmt.Println("Trim:", strings.Trim(str, " !"))
    fmt.Println("Replace:", strings.Replace(str, "o", "0", -1))
    fmt.Println("ReplaceAll:", strings.ReplaceAll(str, "o", "0"))
    
    // 分割与连接 / Split and Join
    parts := strings.Split("a,b,c,d", ",")
    fmt.Println("Split:", parts)
    fmt.Println("Join:", strings.Join(parts, "-"))
    
    // 重复 / Repeat
    fmt.Println("Repeat:", strings.Repeat("Go", 3))
    
    // 字符串拼接 / String concatenation
    // 方式1: + 运算符 (小量拼接)
    result := "Hello" + " " + "World"
    
    // 方式2: strings.Builder (大量拼接，高效)
    var builder strings.Builder
    builder.WriteString("Hello")
    builder.WriteString(" ")
    builder.WriteString("World")
    result = builder.String()
    fmt.Println("Builder result:", result)
    
    // 方式3: fmt.Sprintf (格式化)
    name := "Go"
    version := 1.21
    result = fmt.Sprintf("%s version %.2f", name, version)
    fmt.Println("Sprintf result:", result)
    
    // 字符串与字节切片转换 / String to/from byte slice
    bytes := []byte("Hello")
    str2 := string(bytes)
    fmt.Println("Bytes:", bytes)
    fmt.Println("String:", str2)
    
    // 字符串与 rune 切片转换 / String to/from rune slice
    runes := []rune("Hello, 世界")
    fmt.Println("Runes:", runes)
    fmt.Println("String from runes:", string(runes))
}
```

## 3. 日期时间 (datetime) / Date & Time

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    // 获取当前时间 / Get current time
    now := time.Now()
    fmt.Println("当前时间 / Now:", now)
    
    // 获取时间组件 / Get time components
    fmt.Println("年 / Year:", now.Year())
    fmt.Println("月 / Month:", now.Month())
    fmt.Println("日 / Day:", now.Day())
    fmt.Println("时 / Hour:", now.Hour())
    fmt.Println("分 / Minute:", now.Minute())
    fmt.Println("秒 / Second:", now.Second())
    fmt.Println("纳秒 / Nanosecond:", now.Nanosecond())
    fmt.Println("星期 / Weekday:", now.Weekday())
    fmt.Println("年中第几天 / YearDay:", now.YearDay())
    
    // 创建特定时间 / Create specific time
    t := time.Date(2024, time.December, 25, 10, 30, 0, 0, time.UTC)
    fmt.Println("特定时间 / Specific time:", t)
    
    // 时间格式化 / Time formatting
    // Go 使用特殊的参考时间: Mon Jan 2 15:04:05 MST 2006
    // Go uses special reference time: Mon Jan 2 15:04:05 MST 2006
    fmt.Println("格式1:", now.Format("2006-01-02"))
    fmt.Println("格式2:", now.Format("2006-01-02 15:04:05"))
    fmt.Println("格式3:", now.Format("2006/01/02 03:04:05 PM"))
    fmt.Println("格式4:", now.Format(time.RFC3339))
    fmt.Println("格式5:", now.Format(time.RFC1123))
    
    // 解析时间字符串 / Parse time string
    parsed, err := time.Parse("2006-01-02", "2024-12-25")
    if err == nil {
        fmt.Println("解析结果 / Parsed:", parsed)
    }
    
    // 带时区解析 / Parse with location
    loc, _ := time.LoadLocation("Asia/Shanghai")
    parsedWithLoc, _ := time.ParseInLocation("2006-01-02 15:04:05", "2024-12-25 10:30:00", loc)
    fmt.Println("带时区解析 / Parsed with location:", parsedWithLoc)
    
    // 时间运算 / Time arithmetic
    // Duration 表示时间间隔
    oneHour := time.Hour
    oneDay := 24 * time.Hour
    
    future := now.Add(oneDay)
    past := now.Add(-oneHour)
    fmt.Println("一天后 / One day later:", future)
    fmt.Println("一小时前 / One hour ago:", past)
    
    // 时间差 / Time difference
    diff := future.Sub(past)
    fmt.Println("时间差 / Difference:", diff)
    fmt.Println("小时数 / Hours:", diff.Hours())
    fmt.Println("分钟数 / Minutes:", diff.Minutes())
    fmt.Println("秒数 / Seconds:", diff.Seconds())
    
    // 时间比较 / Time comparison
    t1 := time.Now()
    t2 := t1.Add(time.Hour)
    fmt.Println("t1.Before(t2):", t1.Before(t2))  // true
    fmt.Println("t1.After(t2):", t1.After(t2))    // false
    fmt.Println("t1.Equal(t1):", t1.Equal(t1))    // true
    
    // Unix 时间戳 / Unix timestamp
    fmt.Println("Unix 秒:", now.Unix())
    fmt.Println("Unix 毫秒:", now.UnixMilli())
    fmt.Println("Unix 纳秒:", now.UnixNano())
    
    // 从时间戳创建时间 / Create time from timestamp
    fromUnix := time.Unix(1703500800, 0)
    fmt.Println("从时间戳创建 / From Unix:", fromUnix)
    
    // 时区操作 / Timezone operations
    utc := now.UTC()
    shanghai, _ := time.LoadLocation("Asia/Shanghai")
    newYork, _ := time.LoadLocation("America/New_York")
    
    fmt.Println("UTC:", utc)
    fmt.Println("上海 / Shanghai:", now.In(shanghai))
    fmt.Println("纽约 / New York:", now.In(newYork))
    
    // 定时器 / Timer
    // timer := time.NewTimer(2 * time.Second)
    // <-timer.C  // 阻塞等待 2 秒
    
    // 延时 / Sleep
    // time.Sleep(time.Second)  // 休眠 1 秒
    
    // Ticker (定期触发)
    // ticker := time.NewTicker(time.Second)
    // defer ticker.Stop()
    // for t := range ticker.C {
    //     fmt.Println("Tick at", t)
    // }
}
```

## 4. 正则表达式 (regexp) / Regular Expression

```go
package main

import (
    "fmt"
    "regexp"
)

func main() {
    // 编译正则表达式 / Compile regex
    // 如果模式有误，Compile 返回 error
    // MustCompile 如果有误会 panic
    re := regexp.MustCompile(`\d+`)
    
    // 匹配检查 / Match check
    fmt.Println("MatchString:", re.MatchString("abc123def"))  // true
    fmt.Println("MatchString:", re.MatchString("abcdef"))     // false
    
    // 查找第一个匹配 / Find first match
    text := "Phone: 123-456-7890, Code: 42"
    fmt.Println("FindString:", re.FindString(text))  // "123"
    
    // 查找所有匹配 / Find all matches
    fmt.Println("FindAllString:", re.FindAllString(text, -1))  // [123 456 7890 42]
    
    // 查找匹配的位置 / Find match positions
    fmt.Println("FindStringIndex:", re.FindStringIndex(text))  // [7 10]
    
    // 替换 / Replace
    result := re.ReplaceAllString(text, "XXX")
    fmt.Println("ReplaceAllString:", result)  // "Phone: XXX-XXX-XXX, Code: XXX"
    
    // 使用函数替换 / Replace with function
    result = re.ReplaceAllStringFunc(text, func(s string) string {
        return "[" + s + "]"
    })
    fmt.Println("ReplaceAllStringFunc:", result)
    
    // 分组捕获 / Capture groups
    reEmail := regexp.MustCompile(`(\w+)@(\w+)\.(\w+)`)
    email := "user@example.com"
    
    matches := reEmail.FindStringSubmatch(email)
    fmt.Println("Submatch:", matches)  // [user@example.com user example com]
    
    // 命名捕获组 / Named capture groups
    reNamed := regexp.MustCompile(`(?P<user>\w+)@(?P<domain>\w+)\.(?P<tld>\w+)`)
    matches = reNamed.FindStringSubmatch(email)
    names := reNamed.SubexpNames()
    
    for i, name := range names {
        if i > 0 && name != "" {
            fmt.Printf("%s: %s\n", name, matches[i])
        }
    }
    
    // 分割字符串 / Split string
    reSplit := regexp.MustCompile(`[,;\s]+`)
    parts := reSplit.Split("a,b;c d  e", -1)
    fmt.Println("Split:", parts)  // [a b c d e]
    
    // 常用正则模式 / Common regex patterns
    patterns := map[string]string{
        "email":    `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`,
        "phone":    `^1[3-9]\d{9}$`,  // 中国手机号
        "ipv4":     `^(\d{1,3}\.){3}\d{1,3}$`,
        "url":      `^https?://[^\s]+$`,
        "date":     `^\d{4}-\d{2}-\d{2}$`,
        "chinese":  `[\x{4e00}-\x{9fa5}]+`,  // 中文字符
    }
    
    for name, pattern := range patterns {
        re := regexp.MustCompile(pattern)
        fmt.Printf("%s pattern compiled: %v\n", name, re != nil)
    }
    
    // 验证邮箱 / Validate email
    emailRe := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    testEmails := []string{"test@example.com", "invalid-email", "user@domain.co.uk"}
    for _, e := range testEmails {
        fmt.Printf("%s: %v\n", e, emailRe.MatchString(e))
    }
}
```

## 5. 数组与切片 (Array & Slice) / Array & Slice

### 5.1 数组 (Array)

```go
package main

import "fmt"

func main() {
    // 数组是固定长度的同类型元素序列
    // Arrays are fixed-length sequences of same-type elements
    
    // 声明数组 / Declare array
    var arr1 [5]int                      // 零值初始化
    arr2 := [5]int{1, 2, 3, 4, 5}        // 字面量初始化
    arr3 := [...]int{1, 2, 3}            // 自动计算长度
    arr4 := [5]int{0: 10, 4: 50}         // 指定索引初始化
    
    fmt.Println("arr1:", arr1)  // [0 0 0 0 0]
    fmt.Println("arr2:", arr2)  // [1 2 3 4 5]
    fmt.Println("arr3:", arr3)  // [1 2 3]
    fmt.Println("arr4:", arr4)  // [10 0 0 0 50]
    
    // 数组长度 / Array length
    fmt.Println("len(arr2):", len(arr2))
    
    // 访问和修改元素 / Access and modify elements
    fmt.Println("arr2[0]:", arr2[0])
    arr2[0] = 100
    fmt.Println("Modified arr2:", arr2)
    
    // 遍历数组 / Iterate array
    for i, v := range arr2 {
        fmt.Printf("index=%d, value=%d\n", i, v)
    }
    
    // 多维数组 / Multi-dimensional array
    var matrix [3][3]int
    matrix[0][0] = 1
    matrix[1][1] = 1
    matrix[2][2] = 1
    fmt.Println("Matrix:", matrix)
    
    // 数组是值类型 (赋值会复制) / Arrays are value types (assignment copies)
    arr5 := arr2
    arr5[0] = 999
    fmt.Println("arr2[0]:", arr2[0])  // 100 (未改变)
    fmt.Println("arr5[0]:", arr5[0])  // 999
    
    // 数组比较 / Array comparison
    a := [3]int{1, 2, 3}
    b := [3]int{1, 2, 3}
    fmt.Println("a == b:", a == b)  // true
}
```

### 5.2 切片 (Slice)

```go
package main

import "fmt"

func main() {
    // 切片是动态大小的，对数组的引用
    // Slices are dynamic-sized references to arrays
    
    // 创建切片 / Create slice
    var s1 []int                        // nil 切片
    s2 := []int{1, 2, 3, 4, 5}          // 字面量创建
    s3 := make([]int, 5)                // make 创建，长度5
    s4 := make([]int, 5, 10)            // 长度5，容量10
    
    fmt.Println("s1:", s1, "len:", len(s1), "cap:", cap(s1))
    fmt.Println("s2:", s2, "len:", len(s2), "cap:", cap(s2))
    fmt.Println("s3:", s3, "len:", len(s3), "cap:", cap(s3))
    fmt.Println("s4:", s4, "len:", len(s4), "cap:", cap(s4))
    
    // 从数组创建切片 / Create slice from array
    arr := [5]int{1, 2, 3, 4, 5}
    slice := arr[1:4]  // [2, 3, 4], 索引 1 到 3
    fmt.Println("slice from array:", slice)
    
    // 切片操作 / Slice operations
    fmt.Println("s2[:]:", s2[:])      // 全部
    fmt.Println("s2[:3]:", s2[:3])    // 前3个
    fmt.Println("s2[2:]:", s2[2:])    // 从索引2开始
    fmt.Println("s2[1:4]:", s2[1:4])  // 索引1到3
    
    // append 追加元素 / Append elements
    s := []int{1, 2, 3}
    s = append(s, 4)           // 追加一个元素
    s = append(s, 5, 6, 7)     // 追加多个元素
    s = append(s, []int{8, 9}...)  // 追加另一个切片
    fmt.Println("After append:", s)
    
    // copy 复制切片 / Copy slice
    src := []int{1, 2, 3, 4, 5}
    dst := make([]int, 3)
    n := copy(dst, src)  // 复制 min(len(dst), len(src)) 个元素
    fmt.Println("Copied:", n, "elements, dst:", dst)
    
    // 删除元素 / Delete element
    s = []int{1, 2, 3, 4, 5}
    i := 2  // 删除索引2的元素
    s = append(s[:i], s[i+1:]...)
    fmt.Println("After delete:", s)  // [1 2 4 5]
    
    // 插入元素 / Insert element
    s = []int{1, 2, 4, 5}
    i = 2  // 在索引2处插入
    s = append(s[:i], append([]int{3}, s[i:]...)...)
    fmt.Println("After insert:", s)  // [1 2 3 4 5]
    
    // 切片是引用类型 / Slices are reference types
    original := []int{1, 2, 3}
    copied := original
    copied[0] = 999
    fmt.Println("original:", original)  // [999 2 3] (也被修改)
    
    // 遍历切片 / Iterate slice
    for i, v := range s2 {
        fmt.Printf("index=%d, value=%d\n", i, v)
    }
    
    // 多维切片 / Multi-dimensional slice
    matrix := [][]int{
        {1, 2, 3},
        {4, 5, 6},
        {7, 8, 9},
    }
    fmt.Println("Matrix:", matrix)
    
    // 切片的底层原理 / Slice internals
    // 切片包含三个字段: 指向底层数组的指针、长度、容量
    // A slice has three fields: pointer to underlying array, length, capacity
}
```

## 6. 迭代器 (iter) / Iterator (Go 1.23+)

```go
package main

import (
    "fmt"
    "iter"
    "slices"
)

// Go 1.23 引入了 iter 包和范围函数
// Go 1.23 introduced iter package and range functions

// 自定义迭代器类型 / Custom iterator types
// iter.Seq[V] 是单值迭代器
// iter.Seq2[K, V] 是双值迭代器

// 创建简单的数字范围迭代器 / Create simple number range iterator
func Range(start, end int) iter.Seq[int] {
    return func(yield func(int) bool) {
        for i := start; i < end; i++ {
            if !yield(i) {
                return
            }
        }
    }
}

// 创建键值对迭代器 / Create key-value iterator
func Enumerate[T any](s []T) iter.Seq2[int, T] {
    return func(yield func(int, T) bool) {
        for i, v := range s {
            if !yield(i, v) {
                return
            }
        }
    }
}

// 过滤迭代器 / Filter iterator
func Filter[T any](seq iter.Seq[T], predicate func(T) bool) iter.Seq[T] {
    return func(yield func(T) bool) {
        for v := range seq {
            if predicate(v) {
                if !yield(v) {
                    return
                }
            }
        }
    }
}

// 映射迭代器 / Map iterator
func Map[T, U any](seq iter.Seq[T], transform func(T) U) iter.Seq[U] {
    return func(yield func(U) bool) {
        for v := range seq {
            if !yield(transform(v)) {
                return
            }
        }
    }
}

func main() {
    // 使用 slices 包的迭代器函数 / Use slices package iterator functions
    s := []int{1, 2, 3, 4, 5}
    
    // slices.All - 返回索引和值的迭代器
    // slices.All - returns iterator of index and value
    fmt.Println("slices.All:")
    for i, v := range slices.All(s) {
        fmt.Printf("  [%d]: %d\n", i, v)
    }
    
    // slices.Values - 只返回值的迭代器
    // slices.Values - returns iterator of values only
    fmt.Println("slices.Values:")
    for v := range slices.Values(s) {
        fmt.Printf("  %d\n", v)
    }
    
    // slices.Backward - 反向迭代
    // slices.Backward - iterate in reverse
    fmt.Println("slices.Backward:")
    for i, v := range slices.Backward(s) {
        fmt.Printf("  [%d]: %d\n", i, v)
    }
    
    // 使用自定义迭代器 / Use custom iterator
    fmt.Println("Custom Range(0, 5):")
    for n := range Range(0, 5) {
        fmt.Printf("  %d\n", n)
    }
    
    // 链式迭代器操作 / Chained iterator operations
    fmt.Println("Filter and Map:")
    numbers := Range(1, 10)
    evens := Filter(numbers, func(n int) bool { return n%2 == 0 })
    doubled := Map(evens, func(n int) int { return n * 2 })
    
    for n := range doubled {
        fmt.Printf("  %d\n", n)  // 4, 8, 12, 16
    }
    
    // 将迭代器收集到切片 / Collect iterator to slice
    collected := slices.Collect(Range(0, 5))
    fmt.Println("Collected:", collected)
    
    // 使用 Enumerate / Use Enumerate
    words := []string{"hello", "world", "go"}
    fmt.Println("Enumerate:")
    for i, w := range Enumerate(words) {
        fmt.Printf("  [%d]: %s\n", i, w)
    }
}

// Pull 迭代器 - 手动控制迭代 / Pull iterator - manual iteration control
func demonstratePull() {
    seq := Range(0, 5)
    next, stop := iter.Pull(seq)
    defer stop()
    
    for {
        v, ok := next()
        if !ok {
            break
        }
        fmt.Println(v)
    }
}
```

## 7. 映射 (map) / Map

```go
package main

import (
    "fmt"
    "maps"
    "sort"
)

func main() {
    // map 是键值对的无序集合
    // map is an unordered collection of key-value pairs
    
    // 创建 map / Create map
    var m1 map[string]int              // nil map (不能直接使用)
    m2 := map[string]int{}             // 空 map
    m3 := make(map[string]int)         // make 创建
    m4 := make(map[string]int, 100)    // 指定初始容量
    m5 := map[string]int{              // 字面量初始化
        "one":   1,
        "two":   2,
        "three": 3,
    }
    
    fmt.Println("m1:", m1)
    fmt.Println("m2:", m2)
    fmt.Println("m3:", m3)
    fmt.Println("m4:", m4)
    fmt.Println("m5:", m5)
    
    // 添加/修改元素 / Add/modify elements
    m3["apple"] = 1
    m3["banana"] = 2
    m3["apple"] = 10  // 修改已存在的键
    fmt.Println("m3:", m3)
    
    // 获取元素 / Get element
    value := m3["apple"]
    fmt.Println("apple:", value)
    
    // 检查键是否存在 / Check if key exists
    value, ok := m3["orange"]
    if ok {
        fmt.Println("orange:", value)
    } else {
        fmt.Println("orange not found")
    }
    
    // 简写形式 / Short form
    if v, ok := m3["banana"]; ok {
        fmt.Println("banana:", v)
    }
    
    // 删除元素 / Delete element
    delete(m3, "apple")
    fmt.Println("After delete:", m3)
    
    // 清空 map / Clear map (Go 1.21+)
    clear(m3)
    fmt.Println("After clear:", m3)
    
    // 获取长度 / Get length
    fmt.Println("len(m5):", len(m5))
    
    // 遍历 map / Iterate map
    for key, value := range m5 {
        fmt.Printf("%s: %d\n", key, value)
    }
    
    // 只遍历键 / Iterate keys only
    for key := range m5 {
        fmt.Println("Key:", key)
    }
    
    // 按键排序遍历 / Iterate in sorted key order
    keys := make([]string, 0, len(m5))
    for k := range m5 {
        keys = append(keys, k)
    }
    sort.Strings(keys)
    for _, k := range keys {
        fmt.Printf("%s: %d\n", k, m5[k])
    }
    
    // 使用 maps 包 (Go 1.21+) / Use maps package
    m6 := map[string]int{"a": 1, "b": 2}
    m7 := map[string]int{"b": 3, "c": 4}
    
    // 复制 map / Copy map
    m6Copy := maps.Clone(m6)
    fmt.Println("Clone:", m6Copy)
    
    // 合并 map / Merge maps
    maps.Copy(m6, m7)  // m7 的值复制到 m6
    fmt.Println("After Copy:", m6)
    
    // 比较 map / Compare maps
    m8 := map[string]int{"a": 1, "b": 2}
    m9 := map[string]int{"a": 1, "b": 2}
    fmt.Println("Equal:", maps.Equal(m8, m9))
    
    // 删除满足条件的元素 / Delete elements matching condition
    maps.DeleteFunc(m8, func(k string, v int) bool {
        return v > 1
    })
    fmt.Println("After DeleteFunc:", m8)
    
    // map 的值可以是任何类型 / Map values can be any type
    funcMap := map[string]func(int) int{
        "double": func(x int) int { return x * 2 },
        "square": func(x int) int { return x * x },
    }
    fmt.Println("double(5):", funcMap["double"](5))
    fmt.Println("square(5):", funcMap["square"](5))
    
    // 嵌套 map / Nested map
    nested := map[string]map[string]int{
        "group1": {"a": 1, "b": 2},
        "group2": {"c": 3, "d": 4},
    }
    fmt.Println("nested:", nested)
    fmt.Println("group1.a:", nested["group1"]["a"])
    
    // map 是引用类型 / Maps are reference types
    original := map[string]int{"x": 1}
    ref := original
    ref["x"] = 999
    fmt.Println("original:", original)  // map[x:999]
}
```

## 8. 结构体 (struct) / Struct

```go
package main

import "fmt"

// 定义结构体 / Define struct
type Person struct {
    Name    string
    Age     int
    Email   string
    private string  // 小写字段，包外不可访问
}

// 带标签的结构体 / Struct with tags
type User struct {
    ID       int    `json:"id" db:"user_id"`
    Username string `json:"username" db:"user_name"`
    Password string `json:"-"`  // JSON 序列化时忽略
}

// 嵌套结构体 / Nested struct
type Address struct {
    City    string
    Country string
}

type Employee struct {
    Person  // 匿名嵌入 / Anonymous embedding
    Address Address
    Title   string
}

// 方法 / Methods
func (p Person) Greet() string {
    return fmt.Sprintf("Hello, I'm %s", p.Name)
}

// 指针接收者方法 (可以修改结构体) / Pointer receiver method
func (p *Person) SetAge(age int) {
    p.Age = age
}

// 值接收者 vs 指针接收者 / Value vs Pointer receiver
// 值接收者: 方法操作的是副本
// 指针接收者: 方法可以修改原值，且对大结构体更高效

func main() {
    // 创建结构体实例 / Create struct instance
    // 方式1: 字段顺序初始化 (不推荐)
    p1 := Person{"Alice", 30, "alice@example.com", "secret"}
    
    // 方式2: 命名字段初始化 (推荐)
    p2 := Person{
        Name:  "Bob",
        Age:   25,
        Email: "bob@example.com",
    }
    
    // 方式3: 使用 new (返回指针)
    p3 := new(Person)
    p3.Name = "Charlie"
    p3.Age = 35
    
    // 方式4: 取址操作
    p4 := &Person{Name: "Diana", Age: 28}
    
    fmt.Println("p1:", p1)
    fmt.Println("p2:", p2)
    fmt.Println("p3:", *p3)
    fmt.Println("p4:", *p4)
    
    // 访问字段 / Access fields
    fmt.Println("Name:", p1.Name)
    fmt.Println("Age:", p1.Age)
    
    // 修改字段 / Modify fields
    p2.Age = 26
    fmt.Println("Updated age:", p2.Age)
    
    // 调用方法 / Call methods
    fmt.Println(p1.Greet())
    p1.SetAge(31)
    fmt.Println("After SetAge:", p1.Age)
    
    // 嵌套结构体 / Nested struct
    emp := Employee{
        Person: Person{Name: "Eve", Age: 32},
        Address: Address{
            City:    "Beijing",
            Country: "China",
        },
        Title: "Engineer",
    }
    
    // 访问嵌入字段 / Access embedded fields
    fmt.Println("Employee name:", emp.Name)  // 可以直接访问
    fmt.Println("Employee city:", emp.Address.City)
    fmt.Println("Greet:", emp.Greet())  // 继承方法
    
    // 匿名结构体 / Anonymous struct
    point := struct {
        X, Y int
    }{10, 20}
    fmt.Println("Point:", point)
    
    // 结构体比较 / Struct comparison
    // 如果所有字段都可比较，结构体就可比较
    s1 := struct{ X int }{1}
    s2 := struct{ X int }{1}
    fmt.Println("s1 == s2:", s1 == s2)
    
    // 结构体是值类型 / Structs are value types
    original := Person{Name: "Original", Age: 20}
    copy := original
    copy.Name = "Copy"
    fmt.Println("original.Name:", original.Name)  // "Original"
    fmt.Println("copy.Name:", copy.Name)          // "Copy"
}
```

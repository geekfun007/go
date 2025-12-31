# Go 语言深度详解与实战教程

本教程全面覆盖 Go 语言从基础到高级的所有核心知识点，包含大量实战代码示例。

## 📚 目录结构

| 章节 | 内容 | 文件 |
|------|------|------|
| 01 | [基础语法](#01-基础语法) | [01-syntax.md](./tutorials/01-syntax.md) |
| 02 | [类型与方法](#02-类型与方法) | [02-types.md](./tutorials/02-types.md) |
| 03 | [逻辑与函数](#03-逻辑与函数) | [03-logic-functions.md](./tutorials/03-logic-functions.md) |
| 04 | [并发编程](#04-并发编程) | [04-concurrency.md](./tutorials/04-concurrency.md) |
| 05 | [IO与HTTP](#05-io与http) | [05-io-http.md](./tutorials/05-io-http.md) |
| 06 | [数据库与框架](#06-数据库与框架) | [06-frameworks.md](./tutorials/06-frameworks.md) |

---

## 01. 基础语法

### 1.1 Hello World

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, Go!")
}
```

### 1.2 变量声明

```go
package main

import "fmt"

func main() {
    // 方式1: var 关键字
    var name string = "Alice"
    var age int = 25
    
    // 方式2: 类型推断
    var city = "Beijing"
    
    // 方式3: 短变量声明 (只能在函数内使用)
    country := "China"
    
    // 方式4: 多变量声明
    var (
        x int    = 10
        y string = "hello"
        z bool   = true
    )
    
    // 方式5: 同时声明多个同类型变量
    var a, b, c int = 1, 2, 3
    
    fmt.Println(name, age, city, country)
    fmt.Println(x, y, z)
    fmt.Println(a, b, c)
}
```

### 1.3 常量

```go
package main

import "fmt"

// 普通常量
const Pi = 3.14159
const MaxSize = 100

// 常量组
const (
    StatusOK       = 200
    StatusNotFound = 404
    StatusError    = 500
)

// iota 枚举
const (
    Sunday    = iota // 0
    Monday           // 1
    Tuesday          // 2
    Wednesday        // 3
    Thursday         // 4
    Friday           // 5
    Saturday         // 6
)

// iota 高级用法
const (
    _  = iota             // 忽略第一个值
    KB = 1 << (10 * iota) // 1 << 10 = 1024
    MB                    // 1 << 20
    GB                    // 1 << 30
    TB                    // 1 << 40
)

func main() {
    fmt.Printf("Pi = %v\n", Pi)
    fmt.Printf("Sunday = %d, Monday = %d\n", Sunday, Monday)
    fmt.Printf("KB = %d, MB = %d, GB = %d\n", KB, MB, GB)
}
```

### 1.4 运算符

```go
package main

import "fmt"

func main() {
    // 算术运算符
    a, b := 10, 3
    fmt.Println("加法:", a+b)  // 13
    fmt.Println("减法:", a-b)  // 7
    fmt.Println("乘法:", a*b)  // 30
    fmt.Println("除法:", a/b)  // 3
    fmt.Println("取模:", a%b)  // 1
    
    // 比较运算符
    fmt.Println("等于:", a == b)   // false
    fmt.Println("不等于:", a != b) // true
    fmt.Println("大于:", a > b)    // true
    fmt.Println("小于:", a < b)    // false
    
    // 逻辑运算符
    x, y := true, false
    fmt.Println("与:", x && y) // false
    fmt.Println("或:", x || y) // true
    fmt.Println("非:", !x)     // false
    
    // 位运算符
    c, d := 5, 3 // 二进制: 101, 011
    fmt.Println("位与:", c&d)  // 1 (001)
    fmt.Println("位或:", c|d)  // 7 (111)
    fmt.Println("异或:", c^d)  // 6 (110)
    fmt.Println("左移:", c<<1) // 10 (1010)
    fmt.Println("右移:", c>>1) // 2 (10)
}
```

---

## 02. 类型与方法

### 2.1 基础类型 (Base Types)

#### 整数类型 (int)

```go
package main

import (
    "fmt"
    "math"
)

func main() {
    // 有符号整数
    var i8 int8 = 127                  // -128 到 127
    var i16 int16 = 32767              // -32768 到 32767
    var i32 int32 = 2147483647         // -2^31 到 2^31-1
    var i64 int64 = 9223372036854775807 // -2^63 到 2^63-1
    var i int = 100                     // 平台相关 (32位或64位)
    
    // 无符号整数
    var u8 uint8 = 255                  // 0 到 255
    var u16 uint16 = 65535              // 0 到 65535
    var u32 uint32 = 4294967295         // 0 到 2^32-1
    var u64 uint64 = 18446744073709551615
    
    fmt.Printf("int8: %d, int16: %d, int32: %d, int64: %d\n", i8, i16, i32, i64)
    fmt.Printf("uint8: %d, uint16: %d, uint32: %d, uint64: %d\n", u8, u16, u32, u64)
    fmt.Printf("int: %d\n", i)
    
    // 类型转换
    var x int32 = 100
    var y int64 = int64(x)
    fmt.Printf("转换后: %d\n", y)
    
    // 整数范围
    fmt.Printf("int8 范围: %d 到 %d\n", math.MinInt8, math.MaxInt8)
    fmt.Printf("int64 范围: %d 到 %d\n", math.MinInt64, math.MaxInt64)
}
```

#### 浮点数类型 (float)

```go
package main

import (
    "fmt"
    "math"
)

func main() {
    var f32 float32 = 3.14159
    var f64 float64 = 3.141592653589793
    
    fmt.Printf("float32: %.5f\n", f32)
    fmt.Printf("float64: %.15f\n", f64)
    
    // 科学计数法
    var big = 1.5e10
    var small = 1.5e-10
    fmt.Printf("大数: %e\n", big)
    fmt.Printf("小数: %e\n", small)
    
    // 特殊值
    var inf = math.Inf(1)  // 正无穷
    var ninf = math.Inf(-1) // 负无穷
    var nan = math.NaN()    // 非数字
    
    fmt.Printf("正无穷: %v\n", inf)
    fmt.Printf("负无穷: %v\n", ninf)
    fmt.Printf("NaN: %v\n", nan)
    
    // 浮点数比较 (避免直接使用 ==)
    a, b := 0.1+0.2, 0.3
    epsilon := 1e-9
    if math.Abs(a-b) < epsilon {
        fmt.Println("a 和 b 相等")
    }
    
    // 常用数学函数
    fmt.Printf("向下取整: %v\n", math.Floor(3.7))  // 3
    fmt.Printf("向上取整: %v\n", math.Ceil(3.2))   // 4
    fmt.Printf("四舍五入: %v\n", math.Round(3.5))  // 4
    fmt.Printf("绝对值: %v\n", math.Abs(-5.5))     // 5.5
    fmt.Printf("平方根: %v\n", math.Sqrt(16))      // 4
    fmt.Printf("幂运算: %v\n", math.Pow(2, 10))    // 1024
}
```

#### 布尔类型 (bool)

```go
package main

import "fmt"

func main() {
    var b1 bool = true
    var b2 bool = false
    b3 := true
    
    fmt.Println("b1:", b1)
    fmt.Println("b2:", b2)
    fmt.Println("b3:", b3)
    
    // 布尔运算
    fmt.Println("AND:", b1 && b2)  // false
    fmt.Println("OR:", b1 || b2)   // true
    fmt.Println("NOT:", !b1)       // false
    
    // 比较表达式返回布尔值
    x, y := 10, 20
    fmt.Println("x > y:", x > y)   // false
    fmt.Println("x < y:", x < y)   // true
    fmt.Println("x == y:", x == y) // false
    
    // 注意: Go 中布尔值不能转换为整数
    // var i int = int(true) // 编译错误
}
```

#### 字符串类型 (string)

```go
package main

import (
    "fmt"
    "strings"
    "unicode/utf8"
)

func main() {
    // 字符串声明
    s1 := "Hello, World!"
    s2 := `这是
多行
字符串`
    
    fmt.Println(s1)
    fmt.Println(s2)
    
    // 字符串长度
    s := "Hello, 世界"
    fmt.Println("字节长度:", len(s))                    // 13
    fmt.Println("字符长度:", utf8.RuneCountInString(s)) // 9
    
    // 字符串索引和切片
    fmt.Println("第一个字节:", s[0])      // 72 (H的ASCII码)
    fmt.Println("切片:", s[0:5])          // Hello
    
    // 遍历字符串
    fmt.Println("\n按字节遍历:")
    for i := 0; i < len(s); i++ {
        fmt.Printf("%d: %c\n", i, s[i])
    }
    
    fmt.Println("\n按字符遍历:")
    for i, r := range s {
        fmt.Printf("%d: %c\n", i, r)
    }
    
    // 字符串操作 (strings 包)
    str := "Hello, World!"
    
    fmt.Println("包含:", strings.Contains(str, "World"))     // true
    fmt.Println("前缀:", strings.HasPrefix(str, "Hello"))    // true
    fmt.Println("后缀:", strings.HasSuffix(str, "!"))        // true
    fmt.Println("索引:", strings.Index(str, "World"))        // 7
    fmt.Println("替换:", strings.Replace(str, "World", "Go", 1)) // Hello, Go!
    fmt.Println("大写:", strings.ToUpper(str))               // HELLO, WORLD!
    fmt.Println("小写:", strings.ToLower(str))               // hello, world!
    fmt.Println("去空格:", strings.TrimSpace("  hello  "))   // hello
    fmt.Println("分割:", strings.Split("a,b,c", ","))        // [a b c]
    fmt.Println("连接:", strings.Join([]string{"a", "b", "c"}, "-")) // a-b-c
    
    // 字符串拼接
    // 方式1: + 运算符 (少量拼接)
    result1 := "Hello" + " " + "World"
    
    // 方式2: fmt.Sprintf
    result2 := fmt.Sprintf("%s %s", "Hello", "World")
    
    // 方式3: strings.Builder (大量拼接推荐)
    var builder strings.Builder
    builder.WriteString("Hello")
    builder.WriteString(" ")
    builder.WriteString("World")
    result3 := builder.String()
    
    fmt.Println(result1, result2, result3)
}
```

#### 时间日期类型 (time)

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    // 获取当前时间
    now := time.Now()
    fmt.Println("当前时间:", now)
    
    // 获取时间组件
    fmt.Println("年:", now.Year())
    fmt.Println("月:", now.Month())
    fmt.Println("日:", now.Day())
    fmt.Println("时:", now.Hour())
    fmt.Println("分:", now.Minute())
    fmt.Println("秒:", now.Second())
    fmt.Println("纳秒:", now.Nanosecond())
    fmt.Println("星期:", now.Weekday())
    fmt.Println("今年第几天:", now.YearDay())
    
    // 创建特定时间
    t1 := time.Date(2024, 12, 25, 10, 30, 0, 0, time.Local)
    fmt.Println("指定时间:", t1)
    
    // 时间格式化 (Go 使用特殊的参考时间: 2006-01-02 15:04:05)
    fmt.Println("格式1:", now.Format("2006-01-02"))
    fmt.Println("格式2:", now.Format("2006-01-02 15:04:05"))
    fmt.Println("格式3:", now.Format("2006/01/02 03:04:05 PM"))
    fmt.Println("格式4:", now.Format(time.RFC3339))
    
    // 时间解析
    timeStr := "2024-12-25 10:30:00"
    parsed, err := time.Parse("2006-01-02 15:04:05", timeStr)
    if err != nil {
        fmt.Println("解析错误:", err)
    } else {
        fmt.Println("解析结果:", parsed)
    }
    
    // 时间戳
    fmt.Println("Unix秒:", now.Unix())
    fmt.Println("Unix毫秒:", now.UnixMilli())
    fmt.Println("Unix纳秒:", now.UnixNano())
    
    // 从时间戳创建时间
    t2 := time.Unix(1735084200, 0)
    fmt.Println("从时间戳:", t2)
    
    // 时间运算
    // 加减时间
    future := now.Add(24 * time.Hour)
    past := now.Add(-24 * time.Hour)
    fmt.Println("明天:", future)
    fmt.Println("昨天:", past)
    
    // 添加年月日
    nextYear := now.AddDate(1, 0, 0)
    nextMonth := now.AddDate(0, 1, 0)
    fmt.Println("明年:", nextYear)
    fmt.Println("下月:", nextMonth)
    
    // 时间差
    duration := future.Sub(now)
    fmt.Println("时间差:", duration)
    fmt.Println("小时数:", duration.Hours())
    
    // 时间比较
    fmt.Println("future在now之后:", future.After(now))   // true
    fmt.Println("past在now之前:", past.Before(now))      // true
    fmt.Println("now等于now:", now.Equal(now))           // true
    
    // Duration 类型
    d1 := 5 * time.Second
    d2 := 100 * time.Millisecond
    d3 := time.Duration(1500) * time.Millisecond
    
    fmt.Println("d1:", d1)
    fmt.Println("d2:", d2)
    fmt.Println("d3:", d3)
    
    // 定时器
    fmt.Println("\n等待1秒...")
    time.Sleep(1 * time.Second)
    fmt.Println("完成!")
    
    // 时区
    loc, _ := time.LoadLocation("America/New_York")
    nyTime := now.In(loc)
    fmt.Println("纽约时间:", nyTime)
    
    utcTime := now.UTC()
    fmt.Println("UTC时间:", utcTime)
}
```

#### 正则表达式 (regexp)

```go
package main

import (
    "fmt"
    "regexp"
)

func main() {
    // 编译正则表达式
    // MustCompile 如果正则无效会 panic
    re := regexp.MustCompile(`\d+`)
    
    // Compile 返回错误而不是 panic
    re2, err := regexp.Compile(`\d+`)
    if err != nil {
        fmt.Println("正则编译错误:", err)
        return
    }
    _ = re2
    
    text := "我的电话是 13812345678，座机是 010-12345678"
    
    // 查找第一个匹配
    match := re.FindString(text)
    fmt.Println("第一个匹配:", match) // 13812345678
    
    // 查找所有匹配
    matches := re.FindAllString(text, -1)
    fmt.Println("所有匹配:", matches) // [13812345678 010 12345678]
    
    // 限制匹配数量
    matches2 := re.FindAllString(text, 2)
    fmt.Println("限制2个:", matches2) // [13812345678 010]
    
    // 测试是否匹配
    fmt.Println("是否包含数字:", re.MatchString(text)) // true
    
    // 带分组的正则
    rePhone := regexp.MustCompile(`(\d{3})-(\d{8})`)
    result := rePhone.FindStringSubmatch(text)
    if len(result) > 0 {
        fmt.Println("完整匹配:", result[0]) // 010-12345678
        fmt.Println("区号:", result[1])     // 010
        fmt.Println("号码:", result[2])     // 12345678
    }
    
    // 查找所有分组匹配
    allResults := rePhone.FindAllStringSubmatch(text, -1)
    for i, r := range allResults {
        fmt.Printf("匹配 %d: %v\n", i, r)
    }
    
    // 替换
    reDigit := regexp.MustCompile(`\d`)
    replaced := reDigit.ReplaceAllString(text, "*")
    fmt.Println("替换后:", replaced)
    
    // 使用函数替换
    replaceFunc := reDigit.ReplaceAllStringFunc(text, func(s string) string {
        return "[" + s + "]"
    })
    fmt.Println("函数替换:", replaceFunc)
    
    // 分割
    reSplit := regexp.MustCompile(`[,\s]+`)
    parts := reSplit.Split("a, b,  c,d", -1)
    fmt.Println("分割结果:", parts) // [a b c d]
    
    // 常用正则表达式示例
    patterns := map[string]string{
        "email":    `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`,
        "phone":    `^1[3-9]\d{9}$`,
        "ip":       `^(\d{1,3}\.){3}\d{1,3}$`,
        "url":      `^https?://[^\s]+$`,
        "chinese":  `[\x{4e00}-\x{9fa5}]+`,
        "password": `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[a-zA-Z\d]{8,}$`,
    }
    
    // 验证邮箱
    emailRe := regexp.MustCompile(patterns["email"])
    fmt.Println("验证邮箱 test@example.com:", emailRe.MatchString("test@example.com"))
    
    // 验证手机号
    phoneRe := regexp.MustCompile(patterns["phone"])
    fmt.Println("验证手机 13812345678:", phoneRe.MatchString("13812345678"))
    
    // 提取中文
    chineseRe := regexp.MustCompile(patterns["chinese"])
    chineseMatches := chineseRe.FindAllString("Hello世界Go语言", -1)
    fmt.Println("中文:", chineseMatches) // [世界 语言]
}
```

### 2.2 引用类型 (Reference Types)

#### 数组与切片 (Array & Slice)

```go
package main

import (
    "fmt"
    "sort"
)

func main() {
    // ========== 数组 (固定长度) ==========
    // 声明数组
    var arr1 [5]int                      // 零值初始化
    arr2 := [5]int{1, 2, 3, 4, 5}        // 字面量初始化
    arr3 := [...]int{1, 2, 3}            // 自动计算长度
    arr4 := [5]int{0: 1, 4: 5}           // 指定索引初始化
    
    fmt.Println("arr1:", arr1)
    fmt.Println("arr2:", arr2)
    fmt.Println("arr3:", arr3)
    fmt.Println("arr4:", arr4)
    
    // 数组长度
    fmt.Println("arr2长度:", len(arr2))
    
    // 数组是值类型
    arr5 := arr2
    arr5[0] = 100
    fmt.Println("arr2[0]:", arr2[0]) // 1 (原数组不变)
    fmt.Println("arr5[0]:", arr5[0]) // 100
    
    // ========== 切片 (动态长度) ==========
    // 创建切片
    var s1 []int                         // nil 切片
    s2 := []int{1, 2, 3, 4, 5}           // 字面量
    s3 := make([]int, 5)                 // make(类型, 长度)
    s4 := make([]int, 3, 10)             // make(类型, 长度, 容量)
    
    fmt.Println("\ns1:", s1, "是nil:", s1 == nil)
    fmt.Println("s2:", s2)
    fmt.Println("s3:", s3)
    fmt.Println("s4:", s4, "容量:", cap(s4))
    
    // 从数组创建切片
    arr := [5]int{1, 2, 3, 4, 5}
    slice1 := arr[1:4]   // [2, 3, 4]
    slice2 := arr[:3]    // [1, 2, 3]
    slice3 := arr[2:]    // [3, 4, 5]
    slice4 := arr[:]     // [1, 2, 3, 4, 5]
    
    fmt.Println("slice1:", slice1)
    fmt.Println("slice2:", slice2)
    fmt.Println("slice3:", slice3)
    fmt.Println("slice4:", slice4)
    
    // 切片是引用类型
    s5 := []int{1, 2, 3}
    s6 := s5
    s6[0] = 100
    fmt.Println("s5:", s5) // [100 2 3] (原切片也变了)
    
    // 追加元素
    s7 := []int{1, 2, 3}
    s7 = append(s7, 4)           // 追加一个
    s7 = append(s7, 5, 6, 7)     // 追加多个
    s8 := []int{8, 9}
    s7 = append(s7, s8...)       // 追加另一个切片
    fmt.Println("s7:", s7)
    
    // 复制切片
    src := []int{1, 2, 3}
    dst := make([]int, len(src))
    copy(dst, src)
    fmt.Println("dst:", dst)
    
    // 删除元素
    s9 := []int{1, 2, 3, 4, 5}
    // 删除索引2的元素
    s9 = append(s9[:2], s9[3:]...)
    fmt.Println("删除后:", s9) // [1 2 4 5]
    
    // 切片排序
    nums := []int{3, 1, 4, 1, 5, 9, 2, 6}
    sort.Ints(nums)
    fmt.Println("排序后:", nums)
    
    // 降序排序
    sort.Sort(sort.Reverse(sort.IntSlice(nums)))
    fmt.Println("降序:", nums)
    
    // 字符串切片排序
    strs := []string{"banana", "apple", "cherry"}
    sort.Strings(strs)
    fmt.Println("字符串排序:", strs)
    
    // 自定义排序
    type Person struct {
        Name string
        Age  int
    }
    people := []Person{
        {"Alice", 30},
        {"Bob", 25},
        {"Charlie", 35},
    }
    sort.Slice(people, func(i, j int) bool {
        return people[i].Age < people[j].Age
    })
    fmt.Println("按年龄排序:", people)
    
    // 切片容量扩展规则
    s10 := make([]int, 0)
    for i := 0; i < 10; i++ {
        s10 = append(s10, i)
        fmt.Printf("len=%d, cap=%d\n", len(s10), cap(s10))
    }
}
```

#### Map

```go
package main

import (
    "fmt"
    "sort"
)

func main() {
    // 创建 Map
    var m1 map[string]int              // nil map (不能直接写入)
    m2 := make(map[string]int)         // 空 map
    m3 := map[string]int{              // 字面量初始化
        "apple":  1,
        "banana": 2,
        "cherry": 3,
    }
    
    fmt.Println("m1:", m1, "是nil:", m1 == nil)
    fmt.Println("m2:", m2)
    fmt.Println("m3:", m3)
    
    // 添加/修改元素
    m2["one"] = 1
    m2["two"] = 2
    m2["one"] = 100  // 修改
    fmt.Println("m2:", m2)
    
    // 获取元素
    value := m3["apple"]
    fmt.Println("apple:", value)
    
    // 检查键是否存在
    if val, ok := m3["orange"]; ok {
        fmt.Println("orange:", val)
    } else {
        fmt.Println("orange 不存在")
    }
    
    // 删除元素
    delete(m3, "banana")
    fmt.Println("删除后:", m3)
    
    // 遍历 Map
    m4 := map[string]int{"a": 1, "b": 2, "c": 3}
    for key, value := range m4 {
        fmt.Printf("%s: %d\n", key, value)
    }
    
    // 只遍历键
    for key := range m4 {
        fmt.Println("键:", key)
    }
    
    // Map 长度
    fmt.Println("m4长度:", len(m4))
    
    // 按键排序遍历
    m5 := map[string]int{"c": 3, "a": 1, "b": 2}
    keys := make([]string, 0, len(m5))
    for k := range m5 {
        keys = append(keys, k)
    }
    sort.Strings(keys)
    for _, k := range keys {
        fmt.Printf("%s: %d\n", k, m5[k])
    }
    
    // 嵌套 Map
    nested := map[string]map[string]int{
        "group1": {"a": 1, "b": 2},
        "group2": {"c": 3, "d": 4},
    }
    fmt.Println("嵌套:", nested["group1"]["a"])
    
    // Map 作为集合使用
    set := make(map[string]struct{})
    set["apple"] = struct{}{}
    set["banana"] = struct{}{}
    
    if _, exists := set["apple"]; exists {
        fmt.Println("apple 在集合中")
    }
    
    // 统计词频
    text := "hello world hello go go go"
    words := []string{"hello", "world", "hello", "go", "go", "go"}
    freq := make(map[string]int)
    for _, word := range words {
        freq[word]++
    }
    fmt.Println("词频:", freq)
    _ = text
}
```

#### 结构体 (Struct)

```go
package main

import (
    "encoding/json"
    "fmt"
)

// 定义结构体
type Person struct {
    Name    string
    Age     int
    Email   string
    private string // 小写字母开头，私有字段
}

// 嵌套结构体
type Address struct {
    City    string
    Country string
}

type Employee struct {
    Person           // 匿名嵌入
    Company string
    Address Address  // 命名嵌入
}

// 带标签的结构体 (用于序列化)
type User struct {
    ID        int    `json:"id"`
    Username  string `json:"username"`
    Password  string `json:"-"`                    // 忽略此字段
    Email     string `json:"email,omitempty"`      // 空值时省略
    CreatedAt string `json:"created_at,omitempty"`
}

// 结构体方法
func (p Person) Greet() string {
    return fmt.Sprintf("Hello, I'm %s, %d years old", p.Name, p.Age)
}

// 指针接收器 (可以修改结构体)
func (p *Person) SetAge(age int) {
    p.Age = age
}

// 值接收器 vs 指针接收器
func (p Person) GetAge() int {
    return p.Age
}

func main() {
    // 创建结构体实例
    // 方式1: 字面量
    p1 := Person{Name: "Alice", Age: 25, Email: "alice@example.com"}
    
    // 方式2: 按顺序初始化 (不推荐)
    p2 := Person{"Bob", 30, "bob@example.com", ""}
    
    // 方式3: new (返回指针)
    p3 := new(Person)
    p3.Name = "Charlie"
    p3.Age = 35
    
    // 方式4: 取地址
    p4 := &Person{Name: "David", Age: 40}
    
    fmt.Println("p1:", p1)
    fmt.Println("p2:", p2)
    fmt.Println("p3:", p3)
    fmt.Println("p4:", p4)
    
    // 访问字段
    fmt.Println("p1.Name:", p1.Name)
    fmt.Println("p1.Age:", p1.Age)
    
    // 调用方法
    fmt.Println(p1.Greet())
    
    // 指针方法
    p1.SetAge(26)
    fmt.Println("修改后年龄:", p1.Age)
    
    // 嵌套结构体
    emp := Employee{
        Person:  Person{Name: "Eve", Age: 28},
        Company: "TechCorp",
        Address: Address{City: "Beijing", Country: "China"},
    }
    fmt.Println("员工:", emp)
    fmt.Println("姓名:", emp.Name)         // 直接访问嵌入字段
    fmt.Println("城市:", emp.Address.City)
    
    // 结构体比较
    p5 := Person{Name: "Alice", Age: 25, Email: "alice@example.com"}
    fmt.Println("p1 == p5:", p1 == p5) // true (所有字段相等)
    
    // 匿名结构体
    anonymous := struct {
        X int
        Y int
    }{X: 10, Y: 20}
    fmt.Println("匿名结构体:", anonymous)
    
    // JSON 序列化
    user := User{
        ID:       1,
        Username: "john",
        Password: "secret123",
        Email:    "john@example.com",
    }
    
    jsonData, _ := json.Marshal(user)
    fmt.Println("JSON:", string(jsonData))
    
    // 格式化 JSON
    jsonPretty, _ := json.MarshalIndent(user, "", "  ")
    fmt.Println("格式化JSON:\n", string(jsonPretty))
    
    // JSON 反序列化
    jsonStr := `{"id":2,"username":"jane","email":"jane@example.com"}`
    var user2 User
    json.Unmarshal([]byte(jsonStr), &user2)
    fmt.Println("反序列化:", user2)
}
```

#### 接口 (Interface)

```go
package main

import (
    "fmt"
    "math"
)

// 定义接口
type Shape interface {
    Area() float64
    Perimeter() float64
}

// 实现接口的结构体
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

// 使用接口的函数
func PrintShapeInfo(s Shape) {
    fmt.Printf("面积: %.2f, 周长: %.2f\n", s.Area(), s.Perimeter())
}

// 空接口 (可以接受任何类型)
func PrintAnything(v interface{}) {
    fmt.Printf("类型: %T, 值: %v\n", v, v)
}

// any 是 interface{} 的别名 (Go 1.18+)
func PrintAny(v any) {
    fmt.Printf("类型: %T, 值: %v\n", v, v)
}

// 接口组合
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

// Stringer 接口 (类似 toString)
type Person struct {
    Name string
    Age  int
}

func (p Person) String() string {
    return fmt.Sprintf("Person{Name: %s, Age: %d}", p.Name, p.Age)
}

func main() {
    // 接口变量
    var s Shape
    
    s = Rectangle{Width: 10, Height: 5}
    PrintShapeInfo(s)
    
    s = Circle{Radius: 5}
    PrintShapeInfo(s)
    
    // 接口切片
    shapes := []Shape{
        Rectangle{Width: 3, Height: 4},
        Circle{Radius: 2},
        Rectangle{Width: 5, Height: 5},
    }
    
    for _, shape := range shapes {
        PrintShapeInfo(shape)
    }
    
    // 类型断言
    var i interface{} = "hello"
    
    // 方式1: 直接断言 (失败会panic)
    str := i.(string)
    fmt.Println("字符串:", str)
    
    // 方式2: 安全断言
    if str, ok := i.(string); ok {
        fmt.Println("是字符串:", str)
    }
    
    // 类型switch
    checkType := func(v interface{}) {
        switch t := v.(type) {
        case int:
            fmt.Println("整数:", t)
        case string:
            fmt.Println("字符串:", t)
        case bool:
            fmt.Println("布尔:", t)
        case []int:
            fmt.Println("整数切片:", t)
        default:
            fmt.Printf("未知类型: %T\n", t)
        }
    }
    
    checkType(42)
    checkType("hello")
    checkType(true)
    checkType([]int{1, 2, 3})
    checkType(3.14)
    
    // 空接口
    PrintAnything(42)
    PrintAnything("hello")
    PrintAnything([]int{1, 2, 3})
    
    // Stringer 接口
    p := Person{Name: "Alice", Age: 25}
    fmt.Println(p) // 自动调用 String() 方法
    
    // 检查是否实现接口
    var _ Shape = Rectangle{} // 编译时检查
    var _ Shape = Circle{}    // 编译时检查
}
```

### 2.3 迭代器与Range (Go 1.23+)

```go
package main

import (
    "fmt"
    "iter"
    "maps"
    "slices"
)

// 自定义迭代器 (Go 1.23+)
func Range(start, end int) iter.Seq[int] {
    return func(yield func(int) bool) {
        for i := start; i < end; i++ {
            if !yield(i) {
                return
            }
        }
    }
}

// 带索引的迭代器
func Enumerate[T any](s []T) iter.Seq2[int, T] {
    return func(yield func(int, T) bool) {
        for i, v := range s {
            if !yield(i, v) {
                return
            }
        }
    }
}

// 过滤迭代器
func Filter[T any](s iter.Seq[T], predicate func(T) bool) iter.Seq[T] {
    return func(yield func(T) bool) {
        for v := range s {
            if predicate(v) {
                if !yield(v) {
                    return
                }
            }
        }
    }
}

// 映射迭代器
func Map[T, U any](s iter.Seq[T], transform func(T) U) iter.Seq[U] {
    return func(yield func(U) bool) {
        for v := range s {
            if !yield(transform(v)) {
                return
            }
        }
    }
}

func main() {
    // ========== 基础 range 用法 ==========
    
    // 1. 遍历切片
    nums := []int{10, 20, 30, 40, 50}
    fmt.Println("遍历切片:")
    for i, v := range nums {
        fmt.Printf("  索引 %d: %d\n", i, v)
    }
    
    // 只要值
    for _, v := range nums {
        fmt.Print(v, " ")
    }
    fmt.Println()
    
    // 只要索引
    for i := range nums {
        fmt.Print(i, " ")
    }
    fmt.Println()
    
    // 2. 遍历字符串
    str := "Hello, 世界"
    fmt.Println("\n遍历字符串:")
    for i, r := range str {
        fmt.Printf("  字节位置 %d: %c (Unicode: %U)\n", i, r, r)
    }
    
    // 3. 遍历 Map
    m := map[string]int{"a": 1, "b": 2, "c": 3}
    fmt.Println("\n遍历Map:")
    for k, v := range m {
        fmt.Printf("  %s: %d\n", k, v)
    }
    
    // 4. 遍历通道
    ch := make(chan int, 3)
    ch <- 1
    ch <- 2
    ch <- 3
    close(ch)
    
    fmt.Println("\n遍历通道:")
    for v := range ch {
        fmt.Print(v, " ")
    }
    fmt.Println()
    
    // ========== Go 1.22+ range over integers ==========
    fmt.Println("\nrange整数:")
    for i := range 5 {
        fmt.Print(i, " ")
    }
    fmt.Println()
    
    // ========== Go 1.23+ 迭代器 ==========
    fmt.Println("\n自定义迭代器:")
    for v := range Range(1, 6) {
        fmt.Print(v, " ")
    }
    fmt.Println()
    
    // 带索引的迭代
    fruits := []string{"apple", "banana", "cherry"}
    fmt.Println("\n带索引迭代:")
    for i, v := range Enumerate(fruits) {
        fmt.Printf("  %d: %s\n", i, v)
    }
    
    // 过滤迭代
    fmt.Println("\n过滤偶数:")
    evens := Filter(Range(1, 11), func(n int) bool {
        return n%2 == 0
    })
    for v := range evens {
        fmt.Print(v, " ")
    }
    fmt.Println()
    
    // 映射迭代
    fmt.Println("\n平方映射:")
    squares := Map(Range(1, 6), func(n int) int {
        return n * n
    })
    for v := range squares {
        fmt.Print(v, " ")
    }
    fmt.Println()
    
    // ========== slices 包 (Go 1.21+) ==========
    fmt.Println("\n=== slices 包 ===")
    
    s := []int{3, 1, 4, 1, 5, 9, 2, 6}
    
    // 排序
    slices.Sort(s)
    fmt.Println("排序后:", s)
    
    // 二分查找
    idx, found := slices.BinarySearch(s, 4)
    fmt.Printf("查找4: 索引=%d, 找到=%v\n", idx, found)
    
    // 反转
    slices.Reverse(s)
    fmt.Println("反转后:", s)
    
    // 最大值/最小值
    s2 := []int{3, 1, 4, 1, 5}
    fmt.Println("最大值:", slices.Max(s2))
    fmt.Println("最小值:", slices.Min(s2))
    
    // 包含
    fmt.Println("包含4:", slices.Contains(s2, 4))
    
    // 索引
    fmt.Println("4的索引:", slices.Index(s2, 4))
    
    // 克隆
    s3 := slices.Clone(s2)
    fmt.Println("克隆:", s3)
    
    // 紧凑 (去除连续重复)
    s4 := []int{1, 1, 2, 2, 2, 3, 3}
    s4 = slices.Compact(s4)
    fmt.Println("去重后:", s4)
    
    // 比较
    fmt.Println("相等:", slices.Equal([]int{1, 2}, []int{1, 2}))
    
    // ========== maps 包 (Go 1.21+) ==========
    fmt.Println("\n=== maps 包 ===")
    
    m1 := map[string]int{"a": 1, "b": 2, "c": 3}
    
    // 克隆
    m2 := maps.Clone(m1)
    fmt.Println("克隆:", m2)
    
    // 相等
    fmt.Println("相等:", maps.Equal(m1, m2))
    
    // 获取所有键
    fmt.Println("所有键:", slices.Collect(maps.Keys(m1)))
    
    // 获取所有值
    fmt.Println("所有值:", slices.Collect(maps.Values(m1)))
    
    // 删除符合条件的键值对
    maps.DeleteFunc(m2, func(k string, v int) bool {
        return v > 1
    })
    fmt.Println("删除后:", m2)
}
```

### 2.4 错误处理 (Error Handling)

```go
package main

import (
    "errors"
    "fmt"
    "io"
    "os"
)

// 自定义错误类型
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error: %s - %s", e.Field, e.Message)
}

// 哨兵错误
var (
    ErrNotFound     = errors.New("not found")
    ErrUnauthorized = errors.New("unauthorized")
    ErrInvalidInput = errors.New("invalid input")
)

// 返回错误的函数
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

// 使用自定义错误
func validateAge(age int) error {
    if age < 0 {
        return &ValidationError{
            Field:   "age",
            Message: "age cannot be negative",
        }
    }
    if age > 150 {
        return &ValidationError{
            Field:   "age",
            Message: "age is unrealistic",
        }
    }
    return nil
}

// 错误包装 (Go 1.13+)
func readConfig(filename string) ([]byte, error) {
    data, err := os.ReadFile(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to read config %s: %w", filename, err)
    }
    return data, nil
}

// 使用哨兵错误
func findUser(id int) (string, error) {
    if id <= 0 {
        return "", fmt.Errorf("invalid user id %d: %w", id, ErrInvalidInput)
    }
    if id > 1000 {
        return "", fmt.Errorf("user %d: %w", id, ErrNotFound)
    }
    return "User", nil
}

// panic 和 recover
func mayPanic() {
    panic("something went wrong")
}

func safeCall() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered from panic:", r)
        }
    }()
    mayPanic()
    fmt.Println("This won't be printed")
}

// 多错误处理 (Go 1.20+)
func multipleErrors() error {
    var errs []error
    
    if err := validateAge(-1); err != nil {
        errs = append(errs, err)
    }
    
    if _, err := divide(1, 0); err != nil {
        errs = append(errs, err)
    }
    
    return errors.Join(errs...)
}

func main() {
    // 基本错误处理
    result, err := divide(10, 2)
    if err != nil {
        fmt.Println("错误:", err)
    } else {
        fmt.Println("结果:", result)
    }
    
    // 错误会发生
    _, err = divide(10, 0)
    if err != nil {
        fmt.Println("除法错误:", err)
    }
    
    // 自定义错误
    if err := validateAge(-5); err != nil {
        fmt.Println("验证错误:", err)
        
        // 类型断言获取详细信息
        if ve, ok := err.(*ValidationError); ok {
            fmt.Printf("  字段: %s, 消息: %s\n", ve.Field, ve.Message)
        }
    }
    
    // 错误包装和解包
    _, err = readConfig("nonexistent.json")
    if err != nil {
        fmt.Println("配置错误:", err)
        
        // 检查底层错误
        if errors.Is(err, os.ErrNotExist) {
            fmt.Println("  -> 文件不存在")
        }
        
        // 解包错误
        unwrapped := errors.Unwrap(err)
        fmt.Println("  解包后:", unwrapped)
    }
    
    // 哨兵错误检查
    _, err = findUser(2000)
    if err != nil {
        switch {
        case errors.Is(err, ErrNotFound):
            fmt.Println("用户未找到")
        case errors.Is(err, ErrInvalidInput):
            fmt.Println("输入无效")
        default:
            fmt.Println("其他错误:", err)
        }
    }
    
    // errors.As 获取特定类型的错误
    if err := validateAge(200); err != nil {
        var ve *ValidationError
        if errors.As(err, &ve) {
            fmt.Printf("验证失败: 字段=%s\n", ve.Field)
        }
    }
    
    // panic 和 recover
    fmt.Println("\n测试 panic/recover:")
    safeCall()
    fmt.Println("程序继续运行")
    
    // 多错误
    fmt.Println("\n多错误处理:")
    if err := multipleErrors(); err != nil {
        fmt.Println("多个错误:", err)
    }
    
    // defer 中的错误处理
    fmt.Println("\ndefer 错误处理:")
    if err := processFile("test.txt"); err != nil {
        fmt.Println("处理文件错误:", err)
    }
    
    // 错误处理最佳实践
    fmt.Println("\n错误处理最佳实践:")
    fmt.Println("1. 总是检查错误")
    fmt.Println("2. 尽早返回错误")
    fmt.Println("3. 为错误添加上下文")
    fmt.Println("4. 使用自定义错误类型提供更多信息")
    fmt.Println("5. 使用哨兵错误进行已知错误比较")
    fmt.Println("6. 谨慎使用 panic，仅用于不可恢复的错误")
}

// defer 中处理错误
func processFile(filename string) (err error) {
    f, err := os.Open(filename)
    if err != nil {
        return fmt.Errorf("opening file: %w", err)
    }
    defer func() {
        if cerr := f.Close(); cerr != nil && err == nil {
            err = fmt.Errorf("closing file: %w", cerr)
        }
    }()
    
    // 处理文件...
    buf := make([]byte, 1024)
    _, err = f.Read(buf)
    if err != nil && err != io.EOF {
        return fmt.Errorf("reading file: %w", err)
    }
    
    return nil
}
```

---

## 03. 逻辑与函数

详见 [03-logic-functions.md](./tutorials/03-logic-functions.md)

---

## 04. 并发编程

详见 [04-concurrency.md](./tutorials/04-concurrency.md)

---

## 05. IO与HTTP

详见 [05-io-http.md](./tutorials/05-io-http.md)

---

## 06. 数据库与框架

详见 [06-frameworks.md](./tutorials/06-frameworks.md)

---

## 快速参考

### 常用命令

```bash
# 初始化模块
go mod init module-name

# 运行程序
go run main.go

# 编译程序
go build -o app main.go

# 格式化代码
go fmt ./...

# 代码检查
go vet ./...

# 运行测试
go test ./...

# 下载依赖
go mod tidy

# 查看文档
go doc fmt.Println
```

### 项目结构

```
myproject/
├── cmd/
│   └── app/
│       └── main.go
├── internal/
│   ├── handler/
│   ├── service/
│   └── repository/
├── pkg/
│   └── utils/
├── api/
├── configs/
├── docs/
├── go.mod
├── go.sum
└── README.md
```

---

## 许可证

MIT License

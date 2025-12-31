# Go 类型系统详解

## 1. 基础类型

### 1.1 数值类型

```go
package main

import (
    "fmt"
    "math"
    "strconv"
)

func main() {
    // ========== 整数类型 ==========
    var i8 int8 = 127                    // -128 ~ 127
    var i16 int16 = 32767                // -32768 ~ 32767
    var i32 int32 = 2147483647           // -2^31 ~ 2^31-1
    var i64 int64 = 9223372036854775807  // -2^63 ~ 2^63-1
    var i int = 100                      // 32或64位
    
    var u8 uint8 = 255                   // 0 ~ 255
    var u16 uint16 = 65535               // 0 ~ 65535
    var u32 uint32 = 4294967295          // 0 ~ 2^32-1
    var u64 uint64 = 18446744073709551615
    var u uint = 100                     // 32或64位
    
    // byte 是 uint8 的别名
    var b byte = 'A'
    
    // rune 是 int32 的别名，表示 Unicode 码点
    var r rune = '中'
    
    fmt.Printf("int8: %d, int16: %d\n", i8, i16)
    fmt.Printf("byte: %c, rune: %c\n", b, r)
    
    // 整数字面量
    decimal := 42      // 十进制
    octal := 0o52      // 八进制 (Go 1.13+)
    hex := 0x2a        // 十六进制
    binary := 0b101010 // 二进制 (Go 1.13+)
    
    // 数字分隔符 (Go 1.13+)
    million := 1_000_000
    hexBytes := 0xFF_FF_FF_FF
    
    fmt.Printf("各进制: %d %d %d %d\n", decimal, octal, hex, binary)
    fmt.Printf("分隔符: %d %d\n", million, hexBytes)
    
    // ========== 浮点类型 ==========
    var f32 float32 = 3.14159
    var f64 float64 = 3.141592653589793
    
    // 科学计数法
    scientific := 6.022e23
    small := 1.6e-19
    
    fmt.Printf("float32: %f, float64: %.15f\n", f32, f64)
    fmt.Printf("科学: %e, %e\n", scientific, small)
    
    // 特殊浮点值
    inf := math.Inf(1)
    ninf := math.Inf(-1)
    nan := math.NaN()
    
    fmt.Printf("Inf: %v, -Inf: %v, NaN: %v\n", inf, ninf, nan)
    fmt.Printf("IsNaN: %v, IsInf: %v\n", math.IsNaN(nan), math.IsInf(inf, 1))
    
    // ========== 复数类型 ==========
    var c64 complex64 = 1 + 2i
    var c128 complex128 = 3 + 4i
    
    fmt.Printf("complex64: %v, complex128: %v\n", c64, c128)
    fmt.Printf("实部: %f, 虚部: %f\n", real(c128), imag(c128))
    
    // ========== 类型转换 ==========
    var x int = 100
    var y float64 = float64(x)
    var z int64 = int64(x)
    
    fmt.Printf("转换: int=%d, float64=%f, int64=%d\n", x, y, z)
    
    // 字符串转换
    s := strconv.Itoa(42)           // int -> string
    n, _ := strconv.Atoi("42")      // string -> int
    f, _ := strconv.ParseFloat("3.14", 64)
    
    fmt.Printf("字符串转换: %s, %d, %f\n", s, n, f)
    
    _ = i32
    _ = i64
    _ = i
    _ = u8
    _ = u16
    _ = u32
    _ = u64
    _ = u
}
```

### 1.2 布尔类型

```go
package main

import "fmt"

func main() {
    var b1 bool = true
    var b2 bool = false
    b3 := true
    
    // 零值是 false
    var b4 bool
    fmt.Println("零值:", b4)
    
    // 逻辑运算
    fmt.Println("AND:", b1 && b2)   // false
    fmt.Println("OR:", b1 || b2)    // true
    fmt.Println("NOT:", !b1)        // false
    
    // 比较运算返回布尔值
    x, y := 10, 20
    fmt.Println("x > y:", x > y)
    fmt.Println("x < y:", x < y)
    fmt.Println("x == y:", x == y)
    fmt.Println("x != y:", x != y)
    
    // 短路求值
    result := false && expensiveOperation()  // 不会调用
    fmt.Println("短路:", result)
    
    result = true || expensiveOperation()    // 不会调用
    fmt.Println("短路:", result)
    
    _ = b3
}

func expensiveOperation() bool {
    fmt.Println("执行了昂贵操作")
    return true
}
```

### 1.3 字符串类型

```go
package main

import (
    "fmt"
    "strings"
    "unicode/utf8"
)

func main() {
    // 字符串是不可变的字节序列
    s := "Hello, 世界"
    
    // 长度
    fmt.Println("字节长度:", len(s))                    // 13
    fmt.Println("字符长度:", utf8.RuneCountInString(s)) // 9
    
    // 索引访问 (返回字节)
    fmt.Printf("s[0] = %d (%c)\n", s[0], s[0])
    
    // 切片
    fmt.Println("切片 [0:5]:", s[0:5])  // Hello
    fmt.Println("切片 [:5]:", s[:5])
    fmt.Println("切片 [7:]:", s[7:])
    
    // 遍历 - 按字节
    fmt.Println("\n按字节遍历:")
    for i := 0; i < len(s); i++ {
        fmt.Printf("%d: %x ", i, s[i])
    }
    fmt.Println()
    
    // 遍历 - 按字符 (rune)
    fmt.Println("\n按字符遍历:")
    for i, r := range s {
        fmt.Printf("%d: %c ", i, r)
    }
    fmt.Println()
    
    // 字符串拼接
    // 方式1: + 运算符
    s1 := "Hello" + " " + "World"
    
    // 方式2: fmt.Sprintf
    s2 := fmt.Sprintf("%s %s", "Hello", "World")
    
    // 方式3: strings.Builder (高效)
    var builder strings.Builder
    builder.WriteString("Hello")
    builder.WriteString(" ")
    builder.WriteString("World")
    s3 := builder.String()
    
    // 方式4: strings.Join
    s4 := strings.Join([]string{"Hello", "World"}, " ")
    
    fmt.Println(s1, s2, s3, s4)
    
    // strings 包常用函数
    str := "  Hello, World!  "
    
    fmt.Println("Contains:", strings.Contains(str, "World"))
    fmt.Println("HasPrefix:", strings.HasPrefix(str, "  H"))
    fmt.Println("HasSuffix:", strings.HasSuffix(str, "!  "))
    fmt.Println("Index:", strings.Index(str, "World"))
    fmt.Println("Replace:", strings.Replace(str, "World", "Go", 1))
    fmt.Println("ToUpper:", strings.ToUpper(str))
    fmt.Println("ToLower:", strings.ToLower(str))
    fmt.Println("TrimSpace:", strings.TrimSpace(str))
    fmt.Println("Split:", strings.Split("a,b,c", ","))
    fmt.Println("Count:", strings.Count("hello", "l"))
    fmt.Println("Repeat:", strings.Repeat("ab", 3))
    
    // 字符串和字节切片转换
    bytes := []byte(str)
    str2 := string(bytes)
    fmt.Printf("bytes: %v\nstring: %s\n", bytes, str2)
    
    // 字符串和 rune 切片转换
    runes := []rune(s)
    fmt.Printf("runes: %v\n", runes)
    fmt.Printf("rune[7]: %c\n", runes[7])  // 世
    
    // 多行字符串
    multiline := `这是
多行
字符串`
    fmt.Println(multiline)
    
    // 原始字符串 (不转义)
    raw := `C:\Users\name\file.txt`
    fmt.Println("原始:", raw)
}
```

## 2. 复合类型

### 2.1 数组

```go
package main

import "fmt"

func main() {
    // 数组声明
    var arr1 [5]int                      // 零值
    arr2 := [5]int{1, 2, 3, 4, 5}        // 初始化
    arr3 := [...]int{1, 2, 3}            // 自动长度
    arr4 := [5]int{0: 1, 4: 5}           // 指定索引
    
    fmt.Println("arr1:", arr1)
    fmt.Println("arr2:", arr2)
    fmt.Println("arr3:", arr3)
    fmt.Println("arr4:", arr4)
    
    // 长度和容量
    fmt.Println("长度:", len(arr2))
    
    // 访问和修改
    arr2[0] = 100
    fmt.Println("修改后:", arr2)
    
    // 数组是值类型
    arr5 := arr2
    arr5[0] = 999
    fmt.Println("arr2[0]:", arr2[0])  // 100 (不变)
    fmt.Println("arr5[0]:", arr5[0])  // 999
    
    // 遍历
    for i, v := range arr2 {
        fmt.Printf("arr2[%d] = %d\n", i, v)
    }
    
    // 多维数组
    var matrix [3][3]int
    matrix[0] = [3]int{1, 2, 3}
    matrix[1] = [3]int{4, 5, 6}
    matrix[2] = [3]int{7, 8, 9}
    
    fmt.Println("矩阵:", matrix)
    
    // 多维数组初始化
    matrix2 := [2][3]int{
        {1, 2, 3},
        {4, 5, 6},
    }
    fmt.Println("matrix2:", matrix2)
    
    // 数组比较
    a1 := [3]int{1, 2, 3}
    a2 := [3]int{1, 2, 3}
    a3 := [3]int{1, 2, 4}
    
    fmt.Println("a1 == a2:", a1 == a2)  // true
    fmt.Println("a1 == a3:", a1 == a3)  // false
    
    // 数组作为函数参数 (值传递)
    modifyArray(arr2)
    fmt.Println("函数后arr2:", arr2)  // 不变
    
    // 传递指针
    modifyArrayPtr(&arr2)
    fmt.Println("指针后arr2:", arr2)  // 改变
}

func modifyArray(arr [5]int) {
    arr[0] = 0
}

func modifyArrayPtr(arr *[5]int) {
    arr[0] = 0
}
```

### 2.2 切片

```go
package main

import (
    "fmt"
    "sort"
)

func main() {
    // 切片声明
    var s1 []int              // nil 切片
    s2 := []int{}             // 空切片
    s3 := []int{1, 2, 3}      // 初始化
    s4 := make([]int, 5)      // make(类型, 长度)
    s5 := make([]int, 3, 10)  // make(类型, 长度, 容量)
    
    fmt.Printf("s1: %v, nil: %v, len: %d\n", s1, s1 == nil, len(s1))
    fmt.Printf("s2: %v, nil: %v, len: %d\n", s2, s2 == nil, len(s2))
    fmt.Println("s3:", s3)
    fmt.Println("s4:", s4)
    fmt.Printf("s5: %v, len: %d, cap: %d\n", s5, len(s5), cap(s5))
    
    // 从数组创建切片
    arr := [5]int{1, 2, 3, 4, 5}
    slice := arr[1:4]  // [2, 3, 4]
    fmt.Println("切片:", slice)
    
    // 切片语法
    fmt.Println("arr[:]:", arr[:])    // 全部
    fmt.Println("arr[:3]:", arr[:3])  // [1,2,3]
    fmt.Println("arr[2:]:", arr[2:])  // [3,4,5]
    fmt.Println("arr[1:4]:", arr[1:4])// [2,3,4]
    
    // 切片是引用类型
    s6 := []int{1, 2, 3}
    s7 := s6
    s7[0] = 100
    fmt.Println("s6:", s6)  // [100, 2, 3]
    fmt.Println("s7:", s7)  // [100, 2, 3]
    
    // append
    s8 := []int{1, 2, 3}
    s8 = append(s8, 4)           // 追加一个
    s8 = append(s8, 5, 6, 7)     // 追加多个
    s9 := []int{8, 9}
    s8 = append(s8, s9...)       // 追加切片
    fmt.Println("append后:", s8)
    
    // copy
    src := []int{1, 2, 3, 4, 5}
    dst := make([]int, 3)
    n := copy(dst, src)
    fmt.Printf("复制了 %d 个元素: %v\n", n, dst)
    
    // 深拷贝
    original := []int{1, 2, 3}
    copied := make([]int, len(original))
    copy(copied, original)
    copied[0] = 100
    fmt.Println("original:", original)  // [1, 2, 3]
    fmt.Println("copied:", copied)      // [100, 2, 3]
    
    // 删除元素
    s10 := []int{1, 2, 3, 4, 5}
    i := 2  // 删除索引2
    s10 = append(s10[:i], s10[i+1:]...)
    fmt.Println("删除后:", s10)  // [1, 2, 4, 5]
    
    // 删除但不保持顺序 (更高效)
    s11 := []int{1, 2, 3, 4, 5}
    i = 2
    s11[i] = s11[len(s11)-1]
    s11 = s11[:len(s11)-1]
    fmt.Println("无序删除:", s11)  // [1, 2, 5, 4]
    
    // 插入元素
    s12 := []int{1, 2, 4, 5}
    i = 2
    s12 = append(s12[:i], append([]int{3}, s12[i:]...)...)
    fmt.Println("插入后:", s12)  // [1, 2, 3, 4, 5]
    
    // 排序
    nums := []int{3, 1, 4, 1, 5, 9, 2, 6}
    sort.Ints(nums)
    fmt.Println("排序后:", nums)
    
    // 降序
    sort.Sort(sort.Reverse(sort.IntSlice(nums)))
    fmt.Println("降序:", nums)
    
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
    fmt.Println("按年龄:", people)
    
    // 搜索
    sorted := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
    idx := sort.SearchInts(sorted, 5)
    fmt.Println("5的索引:", idx)
    
    // 切片容量增长
    s13 := make([]int, 0)
    for i := 0; i < 10; i++ {
        s13 = append(s13, i)
        fmt.Printf("len=%d, cap=%d\n", len(s13), cap(s13))
    }
}
```

### 2.3 Map

```go
package main

import (
    "fmt"
    "sort"
    "sync"
)

func main() {
    // Map 声明
    var m1 map[string]int             // nil map
    m2 := make(map[string]int)        // 空 map
    m3 := make(map[string]int, 10)    // 预分配容量
    m4 := map[string]int{             // 初始化
        "apple":  1,
        "banana": 2,
    }
    
    fmt.Println("m1 nil:", m1 == nil)
    fmt.Println("m2:", m2)
    fmt.Println("m3:", m3)
    fmt.Println("m4:", m4)
    
    // 添加和修改
    m2["one"] = 1
    m2["two"] = 2
    m2["one"] = 100  // 修改
    fmt.Println("m2:", m2)
    
    // 获取
    val := m4["apple"]
    fmt.Println("apple:", val)
    
    // 检查键是否存在
    val, ok := m4["orange"]
    if ok {
        fmt.Println("orange:", val)
    } else {
        fmt.Println("orange 不存在")
    }
    
    // 删除
    delete(m4, "banana")
    fmt.Println("删除后:", m4)
    
    // 长度
    fmt.Println("长度:", len(m4))
    
    // 遍历
    m5 := map[string]int{"a": 1, "b": 2, "c": 3}
    for k, v := range m5 {
        fmt.Printf("%s: %d\n", k, v)
    }
    
    // 有序遍历
    keys := make([]string, 0, len(m5))
    for k := range m5 {
        keys = append(keys, k)
    }
    sort.Strings(keys)
    for _, k := range keys {
        fmt.Printf("%s: %d\n", k, m5[k])
    }
    
    // 嵌套 map
    nested := map[string]map[string]int{
        "group1": {"a": 1, "b": 2},
        "group2": {"c": 3, "d": 4},
    }
    fmt.Println("nested:", nested["group1"]["a"])
    
    // Map 作为集合
    set := make(map[string]struct{})
    set["apple"] = struct{}{}
    set["banana"] = struct{}{}
    
    if _, exists := set["apple"]; exists {
        fmt.Println("apple 存在")
    }
    
    // 并发安全的 map
    var sm sync.Map
    sm.Store("key1", "value1")
    sm.Store("key2", "value2")
    
    if v, ok := sm.Load("key1"); ok {
        fmt.Println("sync.Map:", v)
    }
    
    sm.Range(func(k, v interface{}) bool {
        fmt.Printf("sync.Map: %v = %v\n", k, v)
        return true
    })
    
    // 统计词频
    words := []string{"go", "java", "go", "python", "go"}
    freq := make(map[string]int)
    for _, w := range words {
        freq[w]++
    }
    fmt.Println("词频:", freq)
}
```

### 2.4 结构体

```go
package main

import (
    "encoding/json"
    "fmt"
)

// 结构体定义
type Person struct {
    Name    string
    Age     int
    Email   string
    private string  // 私有字段
}

// 嵌套结构体
type Address struct {
    City    string
    Country string
}

type Employee struct {
    Person           // 匿名嵌入
    Company string
    Address Address  // 具名嵌入
}

// 带标签的结构体
type User struct {
    ID        int    `json:"id" db:"user_id"`
    Username  string `json:"username" db:"user_name"`
    Password  string `json:"-"`                    // JSON 忽略
    Email     string `json:"email,omitempty"`      // 空值省略
    CreatedAt string `json:"created_at,omitempty"`
}

// 结构体方法
func (p Person) Greet() string {
    return fmt.Sprintf("Hi, I'm %s", p.Name)
}

// 指针接收器
func (p *Person) SetAge(age int) {
    p.Age = age
}

func main() {
    // 创建实例
    p1 := Person{Name: "Alice", Age: 25, Email: "alice@example.com"}
    p2 := Person{"Bob", 30, "bob@example.com", ""}  // 按顺序
    p3 := new(Person)
    p3.Name = "Charlie"
    p4 := &Person{Name: "David"}
    
    fmt.Println("p1:", p1)
    fmt.Println("p2:", p2)
    fmt.Println("p3:", p3)
    fmt.Println("p4:", p4)
    
    // 访问字段
    fmt.Println("Name:", p1.Name)
    
    // 调用方法
    fmt.Println(p1.Greet())
    
    // 指针方法
    p1.SetAge(26)
    fmt.Println("新年龄:", p1.Age)
    
    // 嵌套结构体
    emp := Employee{
        Person:  Person{Name: "Eve", Age: 28},
        Company: "TechCorp",
        Address: Address{City: "Beijing", Country: "China"},
    }
    fmt.Println("员工:", emp)
    fmt.Println("姓名:", emp.Name)           // 直接访问嵌入字段
    fmt.Println("城市:", emp.Address.City)
    
    // 匿名结构体
    anon := struct {
        X, Y int
    }{10, 20}
    fmt.Println("匿名:", anon)
    
    // 结构体比较
    p5 := Person{Name: "Alice", Age: 25, Email: "alice@example.com"}
    fmt.Println("p1 == p5:", p1 == p5)
    
    // JSON 序列化
    user := User{
        ID:       1,
        Username: "john",
        Password: "secret",
        Email:    "john@example.com",
    }
    
    data, _ := json.Marshal(user)
    fmt.Println("JSON:", string(data))
    
    pretty, _ := json.MarshalIndent(user, "", "  ")
    fmt.Println("Pretty:\n", string(pretty))
    
    // JSON 反序列化
    jsonStr := `{"id":2,"username":"jane","email":"jane@example.com"}`
    var user2 User
    json.Unmarshal([]byte(jsonStr), &user2)
    fmt.Println("反序列化:", user2)
    
    // 结构体切片
    people := []Person{
        {Name: "A", Age: 30},
        {Name: "B", Age: 25},
        {Name: "C", Age: 35},
    }
    fmt.Println("People:", people)
}
```

## 3. 接口

```go
package main

import (
    "fmt"
    "math"
)

// 接口定义
type Shape interface {
    Area() float64
    Perimeter() float64
}

// 实现接口
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

// Stringer 接口
type Person struct {
    Name string
    Age  int
}

func (p Person) String() string {
    return fmt.Sprintf("Person(%s, %d)", p.Name, p.Age)
}

func main() {
    // 接口变量
    var s Shape
    s = Rectangle{10, 5}
    fmt.Printf("Rectangle: Area=%.2f, Perimeter=%.2f\n", s.Area(), s.Perimeter())
    
    s = Circle{5}
    fmt.Printf("Circle: Area=%.2f, Perimeter=%.2f\n", s.Area(), s.Perimeter())
    
    // 接口切片
    shapes := []Shape{
        Rectangle{3, 4},
        Circle{2},
    }
    for _, shape := range shapes {
        fmt.Printf("Area: %.2f\n", shape.Area())
    }
    
    // 类型断言
    var i interface{} = "hello"
    
    // 安全断言
    if str, ok := i.(string); ok {
        fmt.Println("字符串:", str)
    }
    
    // 类型 switch
    checkType := func(v interface{}) {
        switch t := v.(type) {
        case int:
            fmt.Println("int:", t)
        case string:
            fmt.Println("string:", t)
        case bool:
            fmt.Println("bool:", t)
        default:
            fmt.Printf("unknown: %T\n", t)
        }
    }
    
    checkType(42)
    checkType("hello")
    checkType(true)
    checkType(3.14)
    
    // 空接口
    var any interface{}
    any = 42
    any = "hello"
    any = []int{1, 2, 3}
    fmt.Println("any:", any)
    
    // Stringer
    p := Person{"Alice", 25}
    fmt.Println(p)
    
    // 接口零值是 nil
    var s2 Shape
    fmt.Println("nil接口:", s2 == nil)
    
    // 编译时检查接口实现
    var _ Shape = Rectangle{}
    var _ Shape = Circle{}
}
```

## 4. 类型方法

```go
package main

import (
    "fmt"
    "math"
)

// 自定义类型
type Celsius float64
type Fahrenheit float64

// 类型方法
func (c Celsius) ToFahrenheit() Fahrenheit {
    return Fahrenheit(c*9/5 + 32)
}

func (f Fahrenheit) ToCelsius() Celsius {
    return Celsius((f - 32) * 5 / 9)
}

func (c Celsius) String() string {
    return fmt.Sprintf("%.2f°C", c)
}

// 自定义切片类型
type IntSlice []int

func (s IntSlice) Sum() int {
    sum := 0
    for _, v := range s {
        sum += v
    }
    return sum
}

func (s IntSlice) Max() int {
    if len(s) == 0 {
        return 0
    }
    max := s[0]
    for _, v := range s[1:] {
        if v > max {
            max = v
        }
    }
    return max
}

// 函数类型
type MathFunc func(float64) float64

func (f MathFunc) Apply(x float64) float64 {
    return f(x)
}

func (f MathFunc) Compose(g MathFunc) MathFunc {
    return func(x float64) float64 {
        return f(g(x))
    }
}

func main() {
    // 温度转换
    c := Celsius(100)
    f := c.ToFahrenheit()
    fmt.Printf("%v = %.2f°F\n", c, f)
    fmt.Printf("%.2f°F = %v\n", f, f.ToCelsius())
    
    // 自定义切片
    nums := IntSlice{1, 2, 3, 4, 5}
    fmt.Println("Sum:", nums.Sum())
    fmt.Println("Max:", nums.Max())
    
    // 函数类型方法
    square := MathFunc(func(x float64) float64 { return x * x })
    addOne := MathFunc(func(x float64) float64 { return x + 1 })
    
    fmt.Println("square(5):", square.Apply(5))
    
    // 函数组合
    composed := square.Compose(addOne)  // square(addOne(x))
    fmt.Println("square(addOne(3)):", composed(3))  // (3+1)^2 = 16
    
    // 值接收器 vs 指针接收器
    type Counter struct {
        Value int
    }
    
    // 值接收器 - 不修改原值
    increment := func(c Counter) Counter {
        c.Value++
        return c
    }
    
    c1 := Counter{0}
    c2 := increment(c1)
    fmt.Printf("c1: %d, c2: %d\n", c1.Value, c2.Value)
    
    // 指针接收器方法
    type Counter2 struct {
        Value int
    }
    
    incrementPtr := func(c *Counter2) {
        c.Value++
    }
    
    c3 := &Counter2{0}
    incrementPtr(c3)
    fmt.Printf("c3: %d\n", c3.Value)
    
    _ = math.Pi
}
```

## 5. 泛型 (Go 1.18+)

```go
package main

import (
    "fmt"
    "golang.org/x/exp/constraints"
)

// 泛型函数
func Min[T constraints.Ordered](a, b T) T {
    if a < b {
        return a
    }
    return b
}

func Max[T constraints.Ordered](a, b T) T {
    if a > b {
        return a
    }
    return b
}

// 泛型切片函数
func Map[T, U any](s []T, f func(T) U) []U {
    result := make([]U, len(s))
    for i, v := range s {
        result[i] = f(v)
    }
    return result
}

func Filter[T any](s []T, f func(T) bool) []T {
    result := make([]T, 0)
    for _, v := range s {
        if f(v) {
            result = append(result, v)
        }
    }
    return result
}

func Reduce[T, U any](s []T, init U, f func(U, T) U) U {
    result := init
    for _, v := range s {
        result = f(result, v)
    }
    return result
}

// 泛型结构体
type Stack[T any] struct {
    items []T
}

func (s *Stack[T]) Push(item T) {
    s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() (T, bool) {
    if len(s.items) == 0 {
        var zero T
        return zero, false
    }
    item := s.items[len(s.items)-1]
    s.items = s.items[:len(s.items)-1]
    return item, true
}

func (s *Stack[T]) Peek() (T, bool) {
    if len(s.items) == 0 {
        var zero T
        return zero, false
    }
    return s.items[len(s.items)-1], true
}

func (s *Stack[T]) Size() int {
    return len(s.items)
}

// 类型约束
type Number interface {
    ~int | ~int32 | ~int64 | ~float32 | ~float64
}

func Sum[T Number](nums []T) T {
    var sum T
    for _, n := range nums {
        sum += n
    }
    return sum
}

// 自定义约束
type Comparable[T any] interface {
    Compare(T) int
}

func main() {
    // 泛型函数
    fmt.Println("Min(3, 5):", Min(3, 5))
    fmt.Println("Min(3.14, 2.71):", Min(3.14, 2.71))
    fmt.Println("Max(\"abc\", \"xyz\"):", Max("abc", "xyz"))
    
    // Map
    nums := []int{1, 2, 3, 4, 5}
    squares := Map(nums, func(n int) int { return n * n })
    fmt.Println("Squares:", squares)
    
    // Filter
    evens := Filter(nums, func(n int) bool { return n%2 == 0 })
    fmt.Println("Evens:", evens)
    
    // Reduce
    sum := Reduce(nums, 0, func(acc, n int) int { return acc + n })
    fmt.Println("Sum:", sum)
    
    // 泛型栈
    stack := &Stack[int]{}
    stack.Push(1)
    stack.Push(2)
    stack.Push(3)
    
    fmt.Println("Size:", stack.Size())
    if v, ok := stack.Pop(); ok {
        fmt.Println("Pop:", v)
    }
    if v, ok := stack.Peek(); ok {
        fmt.Println("Peek:", v)
    }
    
    // 字符串栈
    strStack := &Stack[string]{}
    strStack.Push("hello")
    strStack.Push("world")
    
    // 类型约束
    intNums := []int{1, 2, 3, 4, 5}
    fmt.Println("Int Sum:", Sum(intNums))
    
    floatNums := []float64{1.1, 2.2, 3.3}
    fmt.Println("Float Sum:", Sum(floatNums))
}
```

## 6. 指针与内存

```go
package main

import (
    "fmt"
    "unsafe"
)

func main() {
    // 基本指针
    x := 10
    p := &x
    fmt.Println("x:", x)
    fmt.Println("p:", p)
    fmt.Println("*p:", *p)
    
    *p = 20
    fmt.Println("修改后x:", x)
    
    // new 分配
    p2 := new(int)
    *p2 = 100
    fmt.Println("new:", *p2)
    
    // 结构体指针
    type Point struct{ X, Y int }
    pt := &Point{10, 20}
    fmt.Println("Point:", pt.X, pt.Y)  // 自动解引用
    
    // 指针接收器
    pt.Move(5, 5)
    fmt.Println("移动后:", pt)
    
    // unsafe 指针
    fmt.Println("\n=== unsafe ===")
    
    var i int64 = 100
    iptr := unsafe.Pointer(&i)
    fptr := (*float64)(iptr)
    fmt.Printf("int64 as float64: %f\n", *fptr)
    
    // 获取结构体字段偏移
    type Example struct {
        A int32
        B int64
        C int32
    }
    
    fmt.Println("A offset:", unsafe.Offsetof(Example{}.A))
    fmt.Println("B offset:", unsafe.Offsetof(Example{}.B))
    fmt.Println("C offset:", unsafe.Offsetof(Example{}.C))
    fmt.Println("Size:", unsafe.Sizeof(Example{}))
    fmt.Println("Align:", unsafe.Alignof(Example{}))
}

func (p *Point) Move(dx, dy int) {
    p.X += dx
    p.Y += dy
}

type Point struct{ X, Y int }
```

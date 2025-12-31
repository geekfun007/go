// Go 语言实战示例 / Go Practical Examples
// 本文件演示了 Go 语言的各种特性

package main

import (
    "encoding/json"
    "fmt"
    "math"
    "regexp"
    "sort"
    "strings"
    "sync"
    "time"
)

func main() {
    fmt.Println("=== Go 语言实战示例 ===\n")
    
    // 1. 基本类型示例 / Basic types example
    basicTypesDemo()
    
    // 2. 字符串操作 / String operations
    stringDemo()
    
    // 3. 时间日期 / Date time
    dateTimeDemo()
    
    // 4. 正则表达式 / Regular expressions
    regexpDemo()
    
    // 5. 切片操作 / Slice operations
    sliceDemo()
    
    // 6. Map 操作 / Map operations
    mapDemo()
    
    // 7. 结构体和方法 / Struct and methods
    structDemo()
    
    // 8. JSON 处理 / JSON handling
    jsonDemo()
    
    // 9. 接口示例 / Interface example
    interfaceDemo()
    
    // 10. 错误处理 / Error handling
    errorDemo()
    
    // 11. 并发示例 / Concurrency example
    concurrencyDemo()
    
    // 12. 高阶函数 / Higher-order functions
    functionalDemo()
}

// ========================================
// 1. 基本类型示例 / Basic Types Demo
// ========================================

func basicTypesDemo() {
    fmt.Println("--- 1. 基本类型 ---")
    
    // 整数 / Integer
    var i int = 42
    var i64 int64 = 9223372036854775807
    
    // 浮点数 / Float
    var f32 float32 = 3.14
    var f64 float64 = 3.141592653589793
    
    // 布尔 / Boolean
    var b bool = true
    
    // 字符串 / String
    var s string = "Hello, Go!"
    
    // 运算 / Operations
    fmt.Printf("整数: %d, int64 最大值: %d\n", i, i64)
    fmt.Printf("浮点数: %.2f, %.10f\n", f32, f64)
    fmt.Printf("布尔: %t\n", b)
    fmt.Printf("字符串: %s\n", s)
    fmt.Printf("数学运算: sqrt(16)=%.0f, pow(2,10)=%.0f\n", math.Sqrt(16), math.Pow(2, 10))
    fmt.Println()
}

// ========================================
// 2. 字符串操作 / String Operations Demo
// ========================================

func stringDemo() {
    fmt.Println("--- 2. 字符串操作 ---")
    
    str := "  Hello, 世界! Welcome to Go!  "
    
    // 基本操作 / Basic operations
    fmt.Printf("原始: %q\n", str)
    fmt.Printf("长度(字节): %d\n", len(str))
    fmt.Printf("TrimSpace: %q\n", strings.TrimSpace(str))
    fmt.Printf("ToUpper: %s\n", strings.ToUpper(str))
    fmt.Printf("ToLower: %s\n", strings.ToLower(str))
    fmt.Printf("Contains 'Go': %t\n", strings.Contains(str, "Go"))
    fmt.Printf("Index 'Go': %d\n", strings.Index(str, "Go"))
    fmt.Printf("Replace: %s\n", strings.Replace(str, "Go", "Golang", 1))
    
    // 分割和连接 / Split and Join
    parts := strings.Split("a,b,c,d", ",")
    fmt.Printf("Split: %v\n", parts)
    fmt.Printf("Join: %s\n", strings.Join(parts, "-"))
    
    // 字符串构建 / String building
    var builder strings.Builder
    builder.WriteString("Hello")
    builder.WriteString(" ")
    builder.WriteString("Builder")
    fmt.Printf("Builder: %s\n", builder.String())
    
    // 遍历 rune / Iterate runes
    fmt.Print("Runes: ")
    for _, r := range "Go世界" {
        fmt.Printf("%c ", r)
    }
    fmt.Println("\n")
}

// ========================================
// 3. 时间日期 / Date Time Demo
// ========================================

func dateTimeDemo() {
    fmt.Println("--- 3. 时间日期 ---")
    
    // 当前时间 / Current time
    now := time.Now()
    fmt.Printf("当前时间: %s\n", now.Format("2006-01-02 15:04:05"))
    fmt.Printf("年: %d, 月: %s, 日: %d\n", now.Year(), now.Month(), now.Day())
    fmt.Printf("时: %d, 分: %d, 秒: %d\n", now.Hour(), now.Minute(), now.Second())
    fmt.Printf("星期: %s\n", now.Weekday())
    
    // 创建特定时间 / Create specific time
    t := time.Date(2024, time.December, 25, 10, 30, 0, 0, time.UTC)
    fmt.Printf("圣诞节: %s\n", t.Format("2006-01-02"))
    
    // 时间运算 / Time arithmetic
    oneDay := 24 * time.Hour
    tomorrow := now.Add(oneDay)
    yesterday := now.Add(-oneDay)
    fmt.Printf("明天: %s\n", tomorrow.Format("2006-01-02"))
    fmt.Printf("昨天: %s\n", yesterday.Format("2006-01-02"))
    
    // 时间差 / Time difference
    diff := tomorrow.Sub(yesterday)
    fmt.Printf("时间差: %.0f 小时\n", diff.Hours())
    
    // Unix 时间戳 / Unix timestamp
    fmt.Printf("Unix 时间戳: %d\n", now.Unix())
    fmt.Println()
}

// ========================================
// 4. 正则表达式 / Regular Expression Demo
// ========================================

func regexpDemo() {
    fmt.Println("--- 4. 正则表达式 ---")
    
    // 编译正则 / Compile regex
    emailRe := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    phoneRe := regexp.MustCompile(`^1[3-9]\d{9}$`)
    
    // 匹配测试 / Match test
    emails := []string{"test@example.com", "invalid-email", "user@domain.co.uk"}
    for _, e := range emails {
        fmt.Printf("邮箱 %s: %t\n", e, emailRe.MatchString(e))
    }
    
    phones := []string{"13812345678", "12345678901", "138123456789"}
    for _, p := range phones {
        fmt.Printf("手机 %s: %t\n", p, phoneRe.MatchString(p))
    }
    
    // 查找和替换 / Find and replace
    text := "电话: 13812345678, 邮箱: test@example.com"
    numberRe := regexp.MustCompile(`\d+`)
    fmt.Printf("找到数字: %v\n", numberRe.FindAllString(text, -1))
    fmt.Printf("替换数字: %s\n", numberRe.ReplaceAllString(text, "***"))
    
    // 分组捕获 / Group capture
    dateRe := regexp.MustCompile(`(\d{4})-(\d{2})-(\d{2})`)
    matches := dateRe.FindStringSubmatch("日期: 2024-12-25")
    if len(matches) > 0 {
        fmt.Printf("完整匹配: %s, 年: %s, 月: %s, 日: %s\n",
            matches[0], matches[1], matches[2], matches[3])
    }
    fmt.Println()
}

// ========================================
// 5. 切片操作 / Slice Operations Demo
// ========================================

func sliceDemo() {
    fmt.Println("--- 5. 切片操作 ---")
    
    // 创建切片 / Create slice
    nums := []int{5, 2, 8, 1, 9, 3, 7, 4, 6}
    fmt.Printf("原始切片: %v\n", nums)
    
    // 追加 / Append
    nums = append(nums, 10)
    fmt.Printf("追加后: %v\n", nums)
    
    // 切片操作 / Slicing
    fmt.Printf("前3个: %v\n", nums[:3])
    fmt.Printf("后3个: %v\n", nums[len(nums)-3:])
    fmt.Printf("中间: %v\n", nums[2:5])
    
    // 复制 / Copy
    numsCopy := make([]int, len(nums))
    copy(numsCopy, nums)
    
    // 排序 / Sort
    sort.Ints(numsCopy)
    fmt.Printf("排序后: %v\n", numsCopy)
    
    // 反转 / Reverse
    for i, j := 0, len(numsCopy)-1; i < j; i, j = i+1, j-1 {
        numsCopy[i], numsCopy[j] = numsCopy[j], numsCopy[i]
    }
    fmt.Printf("反转后: %v\n", numsCopy)
    
    // 删除元素 / Delete element
    idx := 2
    nums = append(nums[:idx], nums[idx+1:]...)
    fmt.Printf("删除索引2后: %v\n", nums)
    
    // 二维切片 / 2D slice
    matrix := [][]int{
        {1, 2, 3},
        {4, 5, 6},
        {7, 8, 9},
    }
    fmt.Printf("矩阵: %v\n", matrix)
    fmt.Println()
}

// ========================================
// 6. Map 操作 / Map Operations Demo
// ========================================

func mapDemo() {
    fmt.Println("--- 6. Map 操作 ---")
    
    // 创建 map / Create map
    scores := map[string]int{
        "Alice":   95,
        "Bob":     87,
        "Charlie": 92,
    }
    fmt.Printf("成绩: %v\n", scores)
    
    // 添加和修改 / Add and modify
    scores["Diana"] = 88
    scores["Alice"] = 98
    fmt.Printf("更新后: %v\n", scores)
    
    // 访问 / Access
    if score, ok := scores["Bob"]; ok {
        fmt.Printf("Bob 的成绩: %d\n", score)
    }
    
    // 删除 / Delete
    delete(scores, "Charlie")
    fmt.Printf("删除后: %v\n", scores)
    
    // 遍历 / Iterate
    fmt.Print("遍历: ")
    for name, score := range scores {
        fmt.Printf("%s=%d ", name, score)
    }
    fmt.Println()
    
    // 按键排序遍历 / Iterate in sorted key order
    keys := make([]string, 0, len(scores))
    for k := range scores {
        keys = append(keys, k)
    }
    sort.Strings(keys)
    fmt.Print("排序遍历: ")
    for _, k := range keys {
        fmt.Printf("%s=%d ", k, scores[k])
    }
    fmt.Println("\n")
}

// ========================================
// 7. 结构体和方法 / Struct and Methods Demo
// ========================================

type Person struct {
    Name string
    Age  int
    City string
}

func (p Person) Introduce() string {
    return fmt.Sprintf("我是 %s, %d 岁, 来自 %s", p.Name, p.Age, p.City)
}

func (p *Person) Birthday() {
    p.Age++
}

type Employee struct {
    Person
    Department string
    Salary     float64
}

func (e Employee) Info() string {
    return fmt.Sprintf("%s, 部门: %s, 薪资: %.2f", e.Introduce(), e.Department, e.Salary)
}

func structDemo() {
    fmt.Println("--- 7. 结构体和方法 ---")
    
    // 创建结构体 / Create struct
    p := Person{Name: "Alice", Age: 25, City: "Beijing"}
    fmt.Println(p.Introduce())
    
    // 修改方法 / Modifier method
    p.Birthday()
    fmt.Printf("生日后: %s, %d 岁\n", p.Name, p.Age)
    
    // 嵌入结构体 / Embedded struct
    emp := Employee{
        Person:     Person{Name: "Bob", Age: 30, City: "Shanghai"},
        Department: "Engineering",
        Salary:     15000.0,
    }
    fmt.Println(emp.Info())
    fmt.Println()
}

// ========================================
// 8. JSON 处理 / JSON Handling Demo
// ========================================

type User struct {
    ID       int      `json:"id"`
    Username string   `json:"username"`
    Email    string   `json:"email,omitempty"`
    Tags     []string `json:"tags"`
    Active   bool     `json:"active"`
}

func jsonDemo() {
    fmt.Println("--- 8. JSON 处理 ---")
    
    // 结构体到 JSON / Struct to JSON
    user := User{
        ID:       1,
        Username: "alice",
        Email:    "alice@example.com",
        Tags:     []string{"golang", "developer"},
        Active:   true,
    }
    
    jsonData, _ := json.Marshal(user)
    fmt.Printf("序列化: %s\n", jsonData)
    
    prettyJSON, _ := json.MarshalIndent(user, "", "  ")
    fmt.Printf("格式化:\n%s\n", prettyJSON)
    
    // JSON 到结构体 / JSON to struct
    jsonStr := `{"id":2,"username":"bob","tags":["backend"],"active":false}`
    var user2 User
    json.Unmarshal([]byte(jsonStr), &user2)
    fmt.Printf("反序列化: %+v\n", user2)
    
    // 动态 JSON / Dynamic JSON
    var data map[string]interface{}
    json.Unmarshal([]byte(jsonStr), &data)
    fmt.Printf("动态解析 username: %v\n", data["username"])
    fmt.Println()
}

// ========================================
// 9. 接口示例 / Interface Demo
// ========================================

type Shape interface {
    Area() float64
    Perimeter() float64
}

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

func printShapeInfo(s Shape) {
    fmt.Printf("面积: %.2f, 周长: %.2f\n", s.Area(), s.Perimeter())
}

func interfaceDemo() {
    fmt.Println("--- 9. 接口示例 ---")
    
    rect := Rectangle{Width: 10, Height: 5}
    circle := Circle{Radius: 7}
    
    fmt.Print("矩形 - ")
    printShapeInfo(rect)
    
    fmt.Print("圆形 - ")
    printShapeInfo(circle)
    
    // 接口切片 / Interface slice
    shapes := []Shape{rect, circle}
    var totalArea float64
    for _, s := range shapes {
        totalArea += s.Area()
    }
    fmt.Printf("总面积: %.2f\n", totalArea)
    
    // 类型断言 / Type assertion
    var s Shape = rect
    if r, ok := s.(Rectangle); ok {
        fmt.Printf("类型断言成功: 宽=%v, 高=%v\n", r.Width, r.Height)
    }
    fmt.Println()
}

// ========================================
// 10. 错误处理 / Error Handling Demo
// ========================================

type ValidationError struct {
    Field   string
    Message string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("验证错误 [%s]: %s", e.Field, e.Message)
}

func validateAge(age int) error {
    if age < 0 {
        return ValidationError{Field: "age", Message: "年龄不能为负数"}
    }
    if age > 150 {
        return ValidationError{Field: "age", Message: "年龄不能超过150"}
    }
    return nil
}

func divide(a, b int) (int, error) {
    if b == 0 {
        return 0, fmt.Errorf("除数不能为零")
    }
    return a / b, nil
}

func errorDemo() {
    fmt.Println("--- 10. 错误处理 ---")
    
    // 简单错误 / Simple error
    result, err := divide(10, 0)
    if err != nil {
        fmt.Printf("除法错误: %v\n", err)
    } else {
        fmt.Printf("结果: %d\n", result)
    }
    
    result, err = divide(10, 3)
    if err == nil {
        fmt.Printf("10/3 = %d\n", result)
    }
    
    // 自定义错误 / Custom error
    if err := validateAge(-5); err != nil {
        fmt.Println(err)
    }
    
    if err := validateAge(200); err != nil {
        fmt.Println(err)
    }
    
    if err := validateAge(25); err == nil {
        fmt.Println("年龄 25 验证通过")
    }
    fmt.Println()
}

// ========================================
// 11. 并发示例 / Concurrency Demo
// ========================================

func concurrencyDemo() {
    fmt.Println("--- 11. 并发示例 ---")
    
    // WaitGroup 示例 / WaitGroup example
    var wg sync.WaitGroup
    results := make([]int, 5)
    
    for i := 0; i < 5; i++ {
        wg.Add(1)
        go func(idx int) {
            defer wg.Done()
            results[idx] = idx * idx
        }(i)
    }
    wg.Wait()
    fmt.Printf("平方结果: %v\n", results)
    
    // Channel 示例 / Channel example
    ch := make(chan int, 3)
    go func() {
        for i := 1; i <= 3; i++ {
            ch <- i * 10
        }
        close(ch)
    }()
    
    fmt.Print("Channel 接收: ")
    for v := range ch {
        fmt.Printf("%d ", v)
    }
    fmt.Println()
    
    // Select 示例 / Select example
    ch1 := make(chan string)
    ch2 := make(chan string)
    
    go func() {
        time.Sleep(50 * time.Millisecond)
        ch1 <- "from ch1"
    }()
    go func() {
        time.Sleep(100 * time.Millisecond)
        ch2 <- "from ch2"
    }()
    
    for i := 0; i < 2; i++ {
        select {
        case msg := <-ch1:
            fmt.Printf("Select 收到: %s\n", msg)
        case msg := <-ch2:
            fmt.Printf("Select 收到: %s\n", msg)
        }
    }
    
    // Mutex 示例 / Mutex example
    var mu sync.Mutex
    counter := 0
    
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            mu.Lock()
            counter++
            mu.Unlock()
        }()
    }
    wg.Wait()
    fmt.Printf("计数器 (Mutex): %d\n", counter)
    fmt.Println()
}

// ========================================
// 12. 高阶函数 / Functional Demo
// ========================================

func Map[T, U any](slice []T, fn func(T) U) []U {
    result := make([]U, len(slice))
    for i, v := range slice {
        result[i] = fn(v)
    }
    return result
}

func Filter[T any](slice []T, predicate func(T) bool) []T {
    result := make([]T, 0)
    for _, v := range slice {
        if predicate(v) {
            result = append(result, v)
        }
    }
    return result
}

func Reduce[T, U any](slice []T, initial U, fn func(U, T) U) U {
    result := initial
    for _, v := range slice {
        result = fn(result, v)
    }
    return result
}

func functionalDemo() {
    fmt.Println("--- 12. 高阶函数 ---")
    
    nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    
    // Map: 平方 / Square
    squared := Map(nums, func(n int) int { return n * n })
    fmt.Printf("平方: %v\n", squared)
    
    // Filter: 偶数 / Even numbers
    evens := Filter(nums, func(n int) bool { return n%2 == 0 })
    fmt.Printf("偶数: %v\n", evens)
    
    // Reduce: 求和 / Sum
    sum := Reduce(nums, 0, func(acc, n int) int { return acc + n })
    fmt.Printf("求和: %d\n", sum)
    
    // 链式操作 / Chained operations
    // 1. 过滤偶数 2. 平方 3. 求和
    result := Reduce(
        Map(
            Filter(nums, func(n int) bool { return n%2 == 0 }),
            func(n int) int { return n * n },
        ),
        0,
        func(acc, n int) int { return acc + n },
    )
    fmt.Printf("偶数平方和: %d\n", result)
}

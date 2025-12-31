# Go JSON 处理 / Go JSON Handling

## 1. 基本序列化与反序列化 / Basic Serialization & Deserialization

```go
package main

import (
    "encoding/json"
    "fmt"
)

type Person struct {
    Name    string   `json:"name"`
    Age     int      `json:"age"`
    Email   string   `json:"email,omitempty"`  // 空值时省略
    Friends []string `json:"friends"`
    private string   // 小写字段不会被序列化
}

func main() {
    // 序列化 (struct -> JSON) / Marshal
    p := Person{
        Name:    "Alice",
        Age:     30,
        Email:   "alice@example.com",
        Friends: []string{"Bob", "Charlie"},
    }
    
    // json.Marshal 返回紧凑的 JSON
    data, err := json.Marshal(p)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Marshal:", string(data))
    
    // json.MarshalIndent 返回格式化的 JSON
    prettyData, _ := json.MarshalIndent(p, "", "  ")
    fmt.Println("MarshalIndent:")
    fmt.Println(string(prettyData))
    
    // 反序列化 (JSON -> struct) / Unmarshal
    jsonStr := `{"name":"Bob","age":25,"friends":["Alice","Diana"]}`
    var p2 Person
    err = json.Unmarshal([]byte(jsonStr), &p2)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Printf("Unmarshal: %+v\n", p2)
}
```

## 2. JSON 标签详解 / JSON Tags

```go
package main

import (
    "encoding/json"
    "fmt"
)

type Example struct {
    // 基本重命名 / Basic rename
    Field1 string `json:"field_1"`
    
    // 空值时省略 / Omit when empty
    Field2 string `json:"field_2,omitempty"`
    
    // 忽略字段 / Ignore field
    Field3 string `json:"-"`
    
    // 保留字段名但可省略 / Keep name but omit if empty
    Field4 string `json:",omitempty"`
    
    // 字符串类型的数字 / Number as string
    Field5 int `json:"field_5,string"`
    
    // 内联结构体 / Inline struct
    Nested struct {
        Inner string `json:"inner"`
    } `json:"nested"`
}

type EmbeddedExample struct {
    // 匿名嵌入会被展平 / Anonymous embedding is flattened
    Person
    Title string `json:"title"`
}

type Person struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

func main() {
    // omitempty 示例 / omitempty example
    type User struct {
        Name  string `json:"name"`
        Email string `json:"email,omitempty"`
        Phone string `json:"phone,omitempty"`
    }
    
    u := User{Name: "Alice"}  // Email 和 Phone 为空
    data, _ := json.Marshal(u)
    fmt.Println("With omitempty:", string(data))
    // 输出: {"name":"Alice"}
    
    // 嵌入示例 / Embedding example
    e := EmbeddedExample{
        Person: Person{Name: "Bob", Age: 30},
        Title:  "Engineer",
    }
    data, _ = json.Marshal(e)
    fmt.Println("Embedded:", string(data))
    // 输出: {"name":"Bob","age":30,"title":"Engineer"}
}
```

## 3. 处理动态 JSON / Handle Dynamic JSON

```go
package main

import (
    "encoding/json"
    "fmt"
)

func main() {
    // 使用 map[string]interface{} / Using map[string]interface{}
    jsonStr := `{
        "name": "Alice",
        "age": 30,
        "active": true,
        "scores": [95, 87, 92],
        "address": {"city": "Beijing"}
    }`
    
    var data map[string]interface{}
    json.Unmarshal([]byte(jsonStr), &data)
    
    // 访问值需要类型断言 / Accessing values requires type assertion
    name := data["name"].(string)
    age := data["age"].(float64)  // JSON 数字默认为 float64
    active := data["active"].(bool)
    scores := data["scores"].([]interface{})
    address := data["address"].(map[string]interface{})
    
    fmt.Println("Name:", name)
    fmt.Println("Age:", int(age))
    fmt.Println("Active:", active)
    fmt.Println("First score:", scores[0])
    fmt.Println("City:", address["city"])
    
    // 使用 json.RawMessage 延迟解析 / Defer parsing with json.RawMessage
    type Response struct {
        Type    string          `json:"type"`
        Payload json.RawMessage `json:"payload"`
    }
    
    jsonResp := `{"type":"user","payload":{"name":"Bob","age":25}}`
    var resp Response
    json.Unmarshal([]byte(jsonResp), &resp)
    
    fmt.Println("Type:", resp.Type)
    fmt.Println("Payload (raw):", string(resp.Payload))
    
    // 根据类型解析 payload
    if resp.Type == "user" {
        var user struct {
            Name string `json:"name"`
            Age  int    `json:"age"`
        }
        json.Unmarshal(resp.Payload, &user)
        fmt.Printf("User: %+v\n", user)
    }
    
    // 使用 json.Number 保持数字精度 / Use json.Number for precision
    type Transaction struct {
        ID     string      `json:"id"`
        Amount json.Number `json:"amount"`
    }
    
    jsonTx := `{"id":"123","amount":"99999999999999999.99"}`
    var tx Transaction
    decoder := json.NewDecoder(strings.NewReader(jsonTx))
    decoder.UseNumber()
    decoder.Decode(&tx)
    
    // 可以转换为不同类型
    amountFloat, _ := tx.Amount.Float64()
    amountInt, _ := tx.Amount.Int64()
    amountStr := tx.Amount.String()
    
    fmt.Println("Amount as float:", amountFloat)
    fmt.Println("Amount as int:", amountInt)
    fmt.Println("Amount as string:", amountStr)
}
```

## 4. 自定义 JSON 序列化 / Custom JSON Serialization

```go
package main

import (
    "encoding/json"
    "fmt"
    "time"
)

// 自定义时间格式 / Custom time format
type CustomTime struct {
    time.Time
}

const timeLayout = "2006-01-02 15:04:05"

func (ct CustomTime) MarshalJSON() ([]byte, error) {
    return json.Marshal(ct.Format(timeLayout))
}

func (ct *CustomTime) UnmarshalJSON(data []byte) error {
    var s string
    if err := json.Unmarshal(data, &s); err != nil {
        return err
    }
    t, err := time.Parse(timeLayout, s)
    if err != nil {
        return err
    }
    ct.Time = t
    return nil
}

// 自定义枚举类型 / Custom enum type
type Status int

const (
    StatusPending Status = iota
    StatusActive
    StatusInactive
)

func (s Status) MarshalJSON() ([]byte, error) {
    names := map[Status]string{
        StatusPending:  "pending",
        StatusActive:   "active",
        StatusInactive: "inactive",
    }
    return json.Marshal(names[s])
}

func (s *Status) UnmarshalJSON(data []byte) error {
    var name string
    if err := json.Unmarshal(data, &name); err != nil {
        return err
    }
    values := map[string]Status{
        "pending":  StatusPending,
        "active":   StatusActive,
        "inactive": StatusInactive,
    }
    *s = values[name]
    return nil
}

type User struct {
    Name      string     `json:"name"`
    Status    Status     `json:"status"`
    CreatedAt CustomTime `json:"created_at"`
}

func main() {
    user := User{
        Name:      "Alice",
        Status:    StatusActive,
        CreatedAt: CustomTime{time.Now()},
    }
    
    data, _ := json.MarshalIndent(user, "", "  ")
    fmt.Println("Custom Marshal:")
    fmt.Println(string(data))
    
    // 反序列化
    jsonStr := `{
        "name": "Bob",
        "status": "pending",
        "created_at": "2024-01-15 10:30:00"
    }`
    
    var user2 User
    json.Unmarshal([]byte(jsonStr), &user2)
    fmt.Printf("Custom Unmarshal: %+v\n", user2)
    fmt.Println("Time:", user2.CreatedAt.Format(time.RFC3339))
}
```

## 5. 流式 JSON 处理 / Streaming JSON

```go
package main

import (
    "encoding/json"
    "fmt"
    "os"
    "strings"
)

type Item struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

func main() {
    // Encoder - 写入流 / Write to stream
    encoder := json.NewEncoder(os.Stdout)
    encoder.SetIndent("", "  ")
    
    items := []Item{
        {ID: 1, Name: "Apple"},
        {ID: 2, Name: "Banana"},
    }
    
    for _, item := range items {
        encoder.Encode(item)
    }
    
    // Decoder - 从流读取 / Read from stream
    jsonStream := `{"id":1,"name":"One"}
{"id":2,"name":"Two"}
{"id":3,"name":"Three"}`
    
    decoder := json.NewDecoder(strings.NewReader(jsonStream))
    
    fmt.Println("\nDecoding stream:")
    for decoder.More() {
        var item Item
        if err := decoder.Decode(&item); err != nil {
            break
        }
        fmt.Printf("  %+v\n", item)
    }
    
    // 处理大型 JSON 数组 / Handle large JSON array
    largeJSON := `[
        {"id": 1, "name": "Item1"},
        {"id": 2, "name": "Item2"},
        {"id": 3, "name": "Item3"}
    ]`
    
    decoder = json.NewDecoder(strings.NewReader(largeJSON))
    
    // 读取开始的 '[' / Read opening '['
    _, err := decoder.Token()
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    
    fmt.Println("\nDecoding array:")
    for decoder.More() {
        var item Item
        decoder.Decode(&item)
        fmt.Printf("  %+v\n", item)
    }
    
    // 读取结束的 ']' / Read closing ']'
    decoder.Token()
}
```

## 6. JSON 验证与错误处理 / JSON Validation & Error Handling

```go
package main

import (
    "encoding/json"
    "errors"
    "fmt"
)

func main() {
    // 验证 JSON 是否有效 / Validate JSON
    validJSON := `{"name":"Alice","age":30}`
    invalidJSON := `{"name":"Alice","age":}`
    
    fmt.Println("Valid JSON:", json.Valid([]byte(validJSON)))
    fmt.Println("Invalid JSON:", json.Valid([]byte(invalidJSON)))
    
    // 处理解析错误 / Handle parse errors
    type User struct {
        Name string `json:"name"`
        Age  int    `json:"age"`
    }
    
    testCases := []string{
        `{"name":"Alice","age":30}`,           // 正确
        `{"name":"Bob","age":"twenty"}`,       // 类型错误
        `{"name":"Charlie","age":999999999999999999999}`, // 溢出
        `{"name":"Diana"}`,                    // 缺少字段
    }
    
    for i, jsonStr := range testCases {
        var user User
        err := json.Unmarshal([]byte(jsonStr), &user)
        if err != nil {
            // 检查错误类型 / Check error type
            var syntaxErr *json.SyntaxError
            var typeErr *json.UnmarshalTypeError
            var numErr *json.InvalidUnmarshalError
            
            switch {
            case errors.As(err, &syntaxErr):
                fmt.Printf("Case %d: Syntax error at offset %d\n", i, syntaxErr.Offset)
            case errors.As(err, &typeErr):
                fmt.Printf("Case %d: Type error - expected %v, got %v for field %s\n",
                    i, typeErr.Type, typeErr.Value, typeErr.Field)
            case errors.As(err, &numErr):
                fmt.Printf("Case %d: Invalid unmarshal error\n", i)
            default:
                fmt.Printf("Case %d: Other error: %v\n", i, err)
            }
        } else {
            fmt.Printf("Case %d: Success - %+v\n", i, user)
        }
    }
    
    // 严格模式 - 不允许未知字段 / Strict mode - disallow unknown fields
    type StrictUser struct {
        Name string `json:"name"`
    }
    
    jsonWithExtra := `{"name":"Eve","unknown":"field"}`
    
    decoder := json.NewDecoder(strings.NewReader(jsonWithExtra))
    decoder.DisallowUnknownFields()
    
    var strictUser StrictUser
    if err := decoder.Decode(&strictUser); err != nil {
        fmt.Println("Strict mode error:", err)
    }
}
```

## 7. 性能优化 / Performance Optimization

```go
package main

import (
    "encoding/json"
    "fmt"
    "sync"
)

// 使用 sync.Pool 复用 buffer / Reuse buffer with sync.Pool
var bufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

func marshalWithPool(v interface{}) ([]byte, error) {
    buf := bufferPool.Get().(*bytes.Buffer)
    defer func() {
        buf.Reset()
        bufferPool.Put(buf)
    }()
    
    encoder := json.NewEncoder(buf)
    if err := encoder.Encode(v); err != nil {
        return nil, err
    }
    
    // 复制结果，因为 buffer 会被复用
    result := make([]byte, buf.Len())
    copy(result, buf.Bytes())
    return result, nil
}

// 预定义结构体以避免反射 / Pre-define struct to avoid reflection
// 对于高性能场景，考虑使用:
// - github.com/json-iterator/go
// - github.com/goccy/go-json
// - github.com/bytedance/sonic

func main() {
    type User struct {
        ID   int    `json:"id"`
        Name string `json:"name"`
    }
    
    user := User{ID: 1, Name: "Alice"}
    
    // 标准库 / Standard library
    data1, _ := json.Marshal(user)
    fmt.Println("Standard:", string(data1))
    
    // 使用 Pool / Using Pool
    data2, _ := marshalWithPool(user)
    fmt.Println("With Pool:", string(data2))
    
    // 性能提示 / Performance tips:
    // 1. 预分配切片/map
    // 2. 使用 json.RawMessage 延迟解析
    // 3. 考虑使用第三方高性能 JSON 库
    // 4. 避免频繁的小对象分配
    // 5. 对于已知结构，手写序列化代码
}
```

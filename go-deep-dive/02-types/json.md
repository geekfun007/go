# JSON 处理 / JSON Handling

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
    
    data, _ := json.Marshal(p)
    fmt.Println("Marshal:", string(data))
    
    // 格式化输出 / Pretty print
    prettyData, _ := json.MarshalIndent(p, "", "  ")
    fmt.Println("MarshalIndent:\n", string(prettyData))
    
    // 反序列化 (JSON -> struct) / Unmarshal
    jsonStr := `{"name":"Bob","age":25,"friends":["Alice"]}`
    var p2 Person
    json.Unmarshal([]byte(jsonStr), &p2)
    fmt.Printf("Unmarshal: %+v\n", p2)
}
```

## 2. JSON 标签 / JSON Tags

```go
type Example struct {
    Field1 string `json:"field_1"`           // 重命名
    Field2 string `json:"field_2,omitempty"` // 空值时省略
    Field3 string `json:"-"`                 // 忽略字段
    Field4 int    `json:"field_4,string"`    // 数字作为字符串
}
```

## 3. 动态 JSON / Dynamic JSON

```go
func dynamicJSON() {
    jsonStr := `{"name":"Alice","age":30,"active":true}`
    
    var data map[string]interface{}
    json.Unmarshal([]byte(jsonStr), &data)
    
    name := data["name"].(string)
    age := data["age"].(float64)  // JSON 数字默认为 float64
    fmt.Println(name, int(age))
    
    // 使用 json.RawMessage 延迟解析
    type Response struct {
        Type    string          `json:"type"`
        Payload json.RawMessage `json:"payload"`
    }
}
```

## 4. 自定义序列化 / Custom Serialization

```go
type CustomTime struct {
    time.Time
}

func (ct CustomTime) MarshalJSON() ([]byte, error) {
    return json.Marshal(ct.Format("2006-01-02 15:04:05"))
}

func (ct *CustomTime) UnmarshalJSON(data []byte) error {
    var s string
    json.Unmarshal(data, &s)
    t, _ := time.Parse("2006-01-02 15:04:05", s)
    ct.Time = t
    return nil
}
```

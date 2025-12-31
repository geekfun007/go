# 结构体 (struct) / Struct

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

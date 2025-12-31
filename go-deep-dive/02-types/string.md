# 字符串 (string) / String

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

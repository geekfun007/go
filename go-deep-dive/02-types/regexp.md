# 正则表达式 (regexp) / Regular Expression

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

# 日期时间 (datetime) / Date & Time

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

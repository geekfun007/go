# 示例索引 / Examples Index

快速查找和运行示例 / Quick reference for examples

---

## Go 异步编程示例 / Go Async Examples

### 01. 基础 Goroutines / Basic Goroutines

**位置**: `go-async/01-basic-goroutines/main.go`

**内容**:
- 简单的 goroutine
- 多个 goroutine
- 使用 WaitGroup
- Goroutine 泄漏问题
- 正确的生命周期管理
- 匿名函数 vs 命名函数

**运行**:
```bash
cd go-async/01-basic-goroutines
go run main.go
```

**关键概念**:
- `go func()`
- `sync.WaitGroup`
- 避免 goroutine 泄漏

---

### 02. Channel 详解 / Channels

**位置**: `go-async/02-channels/main.go`

**内容**:
- 无缓冲 channel
- 有缓冲 channel
- Channel 方向
- 关闭 channel
- Channel 作为信号量
- 多个 goroutine 通信
- 超时处理
- Nil channel 行为

**运行**:
```bash
cd go-async/02-channels
go run main.go
```

**关键概念**:
- `make(chan Type)`
- `make(chan Type, size)`
- `close(ch)`
- `for val := range ch`

---

### 03. Select 语句 / Select Statement

**位置**: `go-async/03-select-statement/main.go`

**内容**:
- 基本 select 用法
- 随机选择
- Default case（非阻塞）
- 超时模式
- Ticker 定时器
- 多路复用（Fan-in）
- 退出信号
- 优先级 select
- 循环中的 select

**运行**:
```bash
cd go-async/03-select-statement
go run main.go
```

**关键概念**:
- `select { case ... }`
- `time.After()`
- `time.NewTicker()`

---

### 04. Context 上下文 / Context

**位置**: `go-async/04-context/main.go`

**内容**:
- 基本 context 使用
- 超时控制
- 截止时间
- 传递值
- 链式取消
- HTTP 请求超时模拟
- 数据库查询超时
- 多 goroutine 协调
- Context 最佳实践

**运行**:
```bash
cd go-async/04-context
go run main.go
```

**关键概念**:
- `context.WithCancel()`
- `context.WithTimeout()`
- `context.WithDeadline()`
- `context.WithValue()`

---

### 05. Sync 包 / Sync Package

**位置**: `go-async/05-sync-package/main.go`

**内容**:
- Mutex 互斥锁
- RWMutex 读写锁
- WaitGroup 详解
- Once 只执行一次
- Cond 条件变量
- Pool 对象池
- sync.Map 并发安全 Map
- Atomic 原子操作
- 死锁示例

**运行**:
```bash
cd go-async/05-sync-package
go run main.go
```

**关键概念**:
- `sync.Mutex`
- `sync.RWMutex`
- `sync.Once`
- `sync.Pool`
- `sync.Map`
- `atomic.*`

---

### 06. 高级并发模式 / Advanced Patterns

**位置**: `go-async/06-advanced-patterns/main.go`

**内容**:
- Worker Pool 工作池
- Pipeline 流水线
- Fan-out/Fan-in 扇出扇入
- Rate Limiting 速率限制
- Burst Rate Limiting 突发速率限制
- Semaphore 信号量
- Error Group 错误组
- Future/Promise 模式
- Broadcast 广播模式

**运行**:
```bash
cd go-async/06-advanced-patterns
go run main.go
```

**关键模式**:
- 工作池
- 流水线
- 扇出/扇入
- 速率限制

---

### 07. 实战示例 / Practical Examples

#### HTTP 服务器

**位置**: `go-async/07-practical-examples/http-server.go`

**内容**:
- 并发 HTTP 处理
- 超时处理
- 并发请求
- 速率限制中间件
- Worker Pool
- 连接池
- 优雅关闭

**运行**:
```bash
cd go-async/07-practical-examples
go run http-server.go
```

#### 数据处理

**位置**: `go-async/07-practical-examples/data-processing.go`

**内容**:
- 批量数据处理
- Map-Reduce 模式
- 流式处理管道
- 并发缓存
- 并发下载器
- 生产者-消费者

**运行**:
```bash
cd go-async/07-practical-examples
go run data-processing.go
```

---

## Python Asyncio 示例 / Python Asyncio Examples

### 01. Asyncio 基础 / Asyncio Basics

**位置**: `python-asyncio/01-asyncio-basics/basic_coroutines.py`

**内容**:
- 基本协程
- 并发执行
- 顺序 vs 并发
- 创建任务
- 超时处理
- 异常处理
- 等待第一个完成
- 协程链式调用
- 按完成顺序处理

**运行**:
```bash
cd python-asyncio/01-asyncio-basics
python basic_coroutines.py
```

**关键概念**:
- `async def`
- `await`
- `asyncio.create_task()`
- `asyncio.gather()`
- `asyncio.wait_for()`

---

### 02. 事件循环 / Event Loop

**位置**: `python-asyncio/02-event-loop/event_loop_basics.py`

**内容**:
- 获取事件循环
- 调度回调
- 运行协程
- Future 对象
- 自定义策略
- 异常处理
- 运行直到完成
- 调试模式
- 运行同步代码
- 生命周期回调

**运行**:
```bash
cd python-asyncio/02-event-loop
python event_loop_basics.py
```

**关键概念**:
- `asyncio.get_event_loop()`
- `loop.create_task()`
- `loop.run_until_complete()`
- `loop.run_in_executor()`

---

### 03. 流和协议 / Streams and Protocols

**位置**: `python-asyncio/03-streams-and-protocols/async_streams.py`

**内容**:
- TCP 客户端
- TCP 服务器
- 聊天服务器
- 异步文件操作
- 流式数据处理
- 协议类
- UDP 协议
- 子进程通信
- 流量控制
- 异步上下文管理器

**运行**:
```bash
cd python-asyncio/03-streams-and-protocols
python async_streams.py
```

**关键概念**:
- `asyncio.open_connection()`
- `asyncio.start_server()`
- `asyncio.Protocol`
- `async with`

---

### 04. 实战示例 - HTTP 客户端 / Practical HTTP Client

**位置**: `python-asyncio/04-practical-examples/async_http_client.py`

**内容**:
- 并发 HTTP 请求
- 批量请求
- 超时处理
- 重试机制
- 速率限制
- 连接池
- 错误处理
- 流式处理
- 取消和清理

**运行**:
```bash
cd python-asyncio/04-practical-examples
python async_http_client.py
```

**实用模式**:
- 并发请求
- 速率限制器
- 连接池
- 重试机制

---

### 05. 自定义实现 / Custom Implementation

**位置**: `python-asyncio/05-custom-implementation/simple_event_loop.py`

**内容**:
- Future 类实现
- Task 类实现
- 事件循环实现
- 协程调度
- 回调机制
- 实现原理讲解

**运行**:
```bash
cd python-asyncio/05-custom-implementation
python simple_event_loop.py
```

**学习重点**:
- 理解事件循环原理
- 协程调度机制
- Future 和 Task 的关系

---

### 06. 高级示例 - 网络爬虫 / Advanced Web Crawler

**位置**: `python-asyncio/06-advanced-examples/async_web_crawler.py`

**内容**:
- 简单爬虫
- 优先级爬虫
- 带缓存的爬虫
- 分布式爬虫（生产者-消费者）

**运行**:
```bash
cd python-asyncio/06-advanced-examples
python async_web_crawler.py
```

**高级特性**:
- 优先级队列
- 缓存机制
- 生产者-消费者模式

---

## 运行所有示例 / Run All Examples

### Go 示例

```bash
# 使用脚本运行所有 Go 示例
./run-go-examples.sh

# 或手动运行
cd go-async/01-basic-goroutines && go run main.go
cd ../02-channels && go run main.go
# ... 依此类推
```

### Python 示例

```bash
# 使用脚本运行所有 Python 示例
./run-python-examples.sh

# 或手动运行
cd python-asyncio/01-asyncio-basics && python basic_coroutines.py
cd ../02-event-loop && python event_loop_basics.py
# ... 依此类推
```

---

## 按主题查找 / Find by Topic

### 并发控制 / Concurrency Control

- Go: Worker Pool (`06-advanced-patterns`)
- Go: Semaphore (`06-advanced-patterns`)
- Python: Connection Pool (`04-practical-examples`)
- Python: Rate Limiter (`04-practical-examples`)

### 错误处理 / Error Handling

- Go: Error Group (`06-advanced-patterns`)
- Python: Exception Handling (`01-asyncio-basics`, `04-practical-examples`)

### 超时和取消 / Timeout and Cancellation

- Go: Context (`04-context`)
- Python: wait_for (`01-asyncio-basics`)

### 通信模式 / Communication Patterns

- Go: Channels (`02-channels`)
- Go: Fan-out/Fan-in (`06-advanced-patterns`)
- Python: Queue (`06-advanced-examples`)

### 数据处理 / Data Processing

- Go: Pipeline (`06-advanced-patterns`)
- Go: Map-Reduce (`07-practical-examples`)
- Python: Stream Processing (`03-streams-and-protocols`)

---

## 难度等级 / Difficulty Levels

### 初级 / Beginner
- ⭐ Go: 01-basic-goroutines
- ⭐ Go: 02-channels
- ⭐ Python: 01-asyncio-basics
- ⭐ Python: 02-event-loop

### 中级 / Intermediate
- ⭐⭐ Go: 03-select-statement
- ⭐⭐ Go: 04-context
- ⭐⭐ Go: 05-sync-package
- ⭐⭐ Python: 03-streams-and-protocols
- ⭐⭐ Python: 04-practical-examples

### 高级 / Advanced
- ⭐⭐⭐ Go: 06-advanced-patterns
- ⭐⭐⭐ Go: 07-practical-examples
- ⭐⭐⭐ Python: 05-custom-implementation
- ⭐⭐⭐ Python: 06-advanced-examples

---

## 学习路径 / Learning Path

### Go 异步编程 / Go Async

1. **基础** → 01, 02
2. **进阶** → 03, 04, 05
3. **高级** → 06, 07

### Python Asyncio

1. **基础** → 01, 02
2. **进阶** → 03, 04
3. **高级** → 05, 06

---

## 常见问题 / FAQ

### Go

**Q: Goroutine 什么时候会退出？**

A: Goroutine 在以下情况退出：
- 函数正常返回
- 发生 panic
- 主程序退出

**Q: Channel 什么时候应该关闭？**

A: 由发送者关闭 channel，当不再有数据要发送时。

**Q: Context 必须用吗？**

A: 不是必须，但强烈推荐用于：
- 超时控制
- 取消操作
- 传递请求范围的值

### Python

**Q: async/await 和多线程有什么区别？**

A: 
- async/await: 单线程并发（协作式）
- 多线程: 真正并行（抢占式）

**Q: 什么时候使用 asyncio？**

A: I/O 密集型任务，如：
- 网络请求
- 文件操作
- 数据库查询

**Q: 如何调试异步代码？**

A:
- 启用调试模式：`loop.set_debug(True)`
- 使用 `PYTHONASYNCIODEBUG=1`
- 添加日志

---

**需要帮助？查看 [README.md](README.md) 或提交 Issue！**

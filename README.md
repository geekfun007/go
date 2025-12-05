# Go 异步详解与实战 + Python Asyncio 实现

> 全面的 Go 并发编程和 Python 异步编程教程与实战

[English](#english) | [中文](#中文)

---

## 中文

### 📚 目录

- [项目简介](#项目简介)
- [项目结构](#项目结构)
- [Go 异步编程](#go-异步编程)
  - [基础概念](#基础概念)
  - [核心特性](#核心特性)
  - [实战示例](#实战示例)
- [Python Asyncio](#python-asyncio)
  - [基础概念](#python-基础概念)
  - [事件循环](#事件循环)
  - [实战应用](#实战应用)
- [对比分析](#对比分析)
- [快速开始](#快速开始)
- [最佳实践](#最佳实践)

---

### 项目简介

本项目提供了全面的 **Go 并发编程** 和 **Python 异步编程** 教程，包括：

- ✨ **Go 异步编程完整指南**：从 Goroutines 到高级并发模式
- 🐍 **Python Asyncio 深度解析**：包含自定义事件循环实现
- 🚀 **丰富的实战示例**：HTTP 服务器、数据处理、流式通信等
- 📊 **详细对比分析**：两种语言异步模型的异同
- 💡 **最佳实践指南**：生产环境可用的设计模式

---

### 项目结构

```
.
├── go-async/                          # Go 异步编程
│   ├── 01-basic-goroutines/          # Goroutine 基础
│   ├── 02-channels/                   # Channel 详解
│   ├── 03-select-statement/          # Select 语句
│   ├── 04-context/                    # Context 上下文
│   ├── 05-sync-package/              # Sync 包
│   ├── 06-advanced-patterns/         # 高级模式
│   └── 07-practical-examples/        # 实战示例
│
├── python-asyncio/                    # Python Asyncio
│   ├── 01-asyncio-basics/            # Asyncio 基础
│   ├── 02-event-loop/                # 事件循环
│   ├── 03-streams-and-protocols/     # 流和协议
│   ├── 04-practical-examples/        # 实战示例
│   └── 05-custom-implementation/     # 自定义实现
│
└── README.md                          # 本文件
```

---

## Go 异步编程

### 基础概念

Go 使用 **CSP (Communicating Sequential Processes)** 模型实现并发：

- **Goroutine**：轻量级线程，由 Go 运行时管理
- **Channel**：用于 Goroutine 间通信的管道
- **Select**：多路复用，监听多个 Channel
- **Context**：管理 Goroutine 生命周期和取消

#### Goroutine 示例

```go
// 启动一个 Goroutine
go func() {
    fmt.Println("Hello from goroutine")
}()

// 使用 WaitGroup 等待完成
var wg sync.WaitGroup
wg.Add(1)
go func() {
    defer wg.Done()
    // 执行任务
}()
wg.Wait()
```

#### Channel 示例

```go
// 无缓冲 channel（同步）
ch := make(chan int)

// 有缓冲 channel（异步）
ch := make(chan int, 10)

// 发送和接收
ch <- 42        // 发送
value := <-ch   // 接收

// 关闭 channel
close(ch)
```

### 核心特性

#### 1. **Goroutines**（轻量级并发）

- **特点**：
  - 初始栈大小仅 2KB
  - 由 Go 调度器管理，非操作系统线程
  - 可以创建数十万个 Goroutine

- **使用场景**：
  - I/O 密集型任务
  - 并发处理请求
  - 后台任务

#### 2. **Channels**（通信机制）

- **类型**：
  - 无缓冲：同步通信
  - 有缓冲：异步通信
  - 单向：只读或只写

- **模式**：
  - 生产者-消费者
  - Fan-out/Fan-in
  - Pipeline（流水线）

#### 3. **Select**（多路复用）

```go
select {
case msg := <-ch1:
    // 处理 ch1
case msg := <-ch2:
    // 处理 ch2
case <-time.After(1 * time.Second):
    // 超时
default:
    // 非阻塞
}
```

#### 4. **Context**（上下文管理）

```go
// 创建带超时的 context
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// 检查取消
select {
case <-ctx.Done():
    return ctx.Err()
default:
    // 继续工作
}
```

#### 5. **Sync 包**（同步原语）

- `sync.Mutex`：互斥锁
- `sync.RWMutex`：读写锁
- `sync.WaitGroup`：等待组
- `sync.Once`：只执行一次
- `sync.Pool`：对象池
- `sync.Map`：并发安全的 Map

### 实战示例

#### Worker Pool（工作池）

```go
func workerPool(jobs <-chan int, results chan<- int, workers int) {
    var wg sync.WaitGroup
    
    for w := 0; w < workers; w++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for job := range jobs {
                results <- processJob(job)
            }
        }()
    }
    
    wg.Wait()
    close(results)
}
```

#### Rate Limiting（速率限制）

```go
limiter := time.Tick(500 * time.Millisecond)

for req := range requests {
    <-limiter  // 等待令牌
    go handleRequest(req)
}
```

#### Pipeline（流水线）

```go
func pipeline() {
    // 阶段1：生成
    gen := func() <-chan int {
        out := make(chan int)
        go func() {
            defer close(out)
            for i := 1; i <= 10; i++ {
                out <- i
            }
        }()
        return out
    }
    
    // 阶段2：处理
    square := func(in <-chan int) <-chan int {
        out := make(chan int)
        go func() {
            defer close(out)
            for n := range in {
                out <- n * n
            }
        }()
        return out
    }
    
    // 构建流水线
    for n := range square(gen()) {
        fmt.Println(n)
    }
}
```

---

## Python Asyncio

### Python 基础概念

Python 使用 **事件循环 + 协程** 模型实现异步：

- **Coroutine**：协程函数（async def）
- **Event Loop**：事件循环，调度协程执行
- **Task**：封装协程的任务对象
- **Future**：代表未来的结果

#### 基本示例

```python
import asyncio

# 定义协程
async def hello():
    print("Hello")
    await asyncio.sleep(1)
    print("World")
    return "Done"

# 运行协程
result = asyncio.run(hello())
```

#### 并发执行

```python
async def main():
    # 并发执行多个协程
    results = await asyncio.gather(
        task1(),
        task2(),
        task3()
    )
    return results
```

### 事件循环

事件循环是 asyncio 的核心：

```python
# 获取事件循环
loop = asyncio.get_event_loop()

# 创建任务
task = loop.create_task(my_coroutine())

# 调度回调
loop.call_soon(callback)
loop.call_later(5, callback)

# 运行直到完成
result = loop.run_until_complete(coroutine)
```

### 核心 API

#### 1. **协程和任务**

```python
# 创建任务
task = asyncio.create_task(coroutine())

# 等待完成
result = await task

# 并发执行
results = await asyncio.gather(*tasks)

# 等待第一个完成
done, pending = await asyncio.wait(
    tasks,
    return_when=asyncio.FIRST_COMPLETED
)
```

#### 2. **超时和取消**

```python
# 超时
try:
    result = await asyncio.wait_for(coroutine(), timeout=5.0)
except asyncio.TimeoutError:
    print("Timeout!")

# 取消任务
task.cancel()
try:
    await task
except asyncio.CancelledError:
    print("Task cancelled")
```

#### 3. **流和协议**

```python
# TCP 客户端
reader, writer = await asyncio.open_connection('host', port)
writer.write(b'data')
await writer.drain()
data = await reader.read(1024)
writer.close()
await writer.wait_closed()

# TCP 服务器
async def handle_client(reader, writer):
    data = await reader.read(1024)
    writer.write(data)
    await writer.drain()
    writer.close()

server = await asyncio.start_server(handle_client, 'localhost', 8888)
```

#### 4. **同步代码集成**

```python
import concurrent.futures

# 在线程池中运行阻塞代码
loop = asyncio.get_event_loop()
result = await loop.run_in_executor(None, blocking_function)

# 使用自定义执行器
with concurrent.futures.ThreadPoolExecutor() as executor:
    result = await loop.run_in_executor(executor, blocking_function)
```

### 实战应用

#### 异步 HTTP 客户端

```python
import aiohttp

async def fetch_urls(urls):
    async with aiohttp.ClientSession() as session:
        tasks = [fetch_one(session, url) for url in urls]
        return await asyncio.gather(*tasks)

async def fetch_one(session, url):
    async with session.get(url) as response:
        return await response.text()
```

#### 速率限制

```python
class RateLimiter:
    def __init__(self, rate, per=1.0):
        self.rate = rate
        self.per = per
        self.allowance = rate
        self.last_check = time.time()
    
    async def acquire(self):
        current = time.time()
        time_passed = current - self.last_check
        self.last_check = current
        self.allowance += time_passed * (self.rate / self.per)
        
        if self.allowance > self.rate:
            self.allowance = self.rate
        
        if self.allowance < 1.0:
            await asyncio.sleep((1.0 - self.allowance) * (self.per / self.rate))
            self.allowance = 0.0
        else:
            self.allowance -= 1.0
```

#### 连接池

```python
class ConnectionPool:
    def __init__(self, size):
        self.pool = asyncio.Queue(maxsize=size)
        for i in range(size):
            self.pool.put_nowait(Connection(i))
    
    async def acquire(self):
        return await self.pool.get()
    
    async def release(self, conn):
        await self.pool.put(conn)
```

---

## 对比分析

### Go vs Python Asyncio

| 特性 | Go | Python Asyncio |
|-----|-----|---------------|
| **并发模型** | CSP (goroutines + channels) | 事件循环 + 协程 |
| **语法** | `go func()` | `async def` / `await` |
| **调度** | M:N 调度器（抢占式） | 事件循环（协作式） |
| **性能** | 更高（编译型，真正并行） | 较低（解释型，单线程） |
| **内存** | Goroutine 2KB 起始栈 | 协程更轻量 |
| **I/O** | 内置网络库，自动异步 | 需要异步库（aiohttp等） |
| **学习曲线** | 相对简单 | 需要理解事件循环 |
| **生态** | 强大的标准库 | 需要第三方库 |
| **适用场景** | 高性能服务、系统编程 | I/O 密集型、脚本任务 |

### 相似之处

1. **都支持并发**：多任务同时执行
2. **都有通信机制**：Go 的 Channel，Python 的 Queue
3. **都支持超时和取消**：Context vs asyncio.wait_for
4. **都有同步原语**：Mutex, WaitGroup vs Lock, Event

### 主要区别

#### 1. **并发 vs 并行**

- **Go**：真正的并行执行（多核）
  ```go
  // 自动利用多核
  runtime.GOMAXPROCS(runtime.NumCPU())
  ```

- **Python**：并发但不并行（GIL限制）
  ```python
  # 单线程事件循环
  # 需要多进程才能利用多核
  ```

#### 2. **调度模型**

- **Go**：抢占式调度
  - Goroutine 可以被强制切换
  - 运行时自动管理

- **Python**：协作式调度
  - 必须显式 `await` 让出控制权
  - 长时间运行会阻塞事件循环

#### 3. **内存模型**

- **Go**：共享内存 + Channel
  ```go
  // 鼓励使用 channel 通信
  // "不要通过共享内存来通信，而是通过通信来共享内存"
  ```

- **Python**：事件循环内共享状态
  ```python
  # 单线程内可以安全共享
  # 但需要小心异步操作的顺序
  ```

---

## 快速开始

### Go 示例

```bash
# 进入 Go 示例目录
cd go-async/01-basic-goroutines

# 运行示例
go run main.go

# 运行所有 Go 示例
for dir in go-async/*/; do
    echo "Running $dir"
    cd "$dir"
    go run *.go
    cd ../..
done
```

### Python 示例

```bash
# 进入 Python 示例目录
cd python-asyncio/01-asyncio-basics

# 运行示例
python basic_coroutines.py

# 运行所有 Python 示例
for dir in python-asyncio/*/; do
    echo "Running $dir"
    cd "$dir"
    python *.py
    cd ../..
done
```

### 依赖安装

#### Go
```bash
# Go 标准库已包含所有需要的包
# 确保 Go 版本 >= 1.16
go version
```

#### Python
```bash
# Python >= 3.7
python --version

# 安装可选依赖（用于实际项目）
pip install aiohttp aiofiles
```

---

## 最佳实践

### Go 最佳实践

#### 1. **避免 Goroutine 泄漏**

```go
// ❌ 错误：goroutine 永远阻塞
ch := make(chan int)
go func() {
    val := <-ch  // 永远等待
}()

// ✅ 正确：使用 context 或 done channel
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

go func() {
    select {
    case val := <-ch:
        // 处理
    case <-ctx.Done():
        return
    }
}()
```

#### 2. **正确使用 WaitGroup**

```go
var wg sync.WaitGroup

for i := 0; i < 10; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        // 注意：传递 id 参数，避免闭包陷阱
        process(id)
    }(i)
}

wg.Wait()
```

#### 3. **Channel 方向**

```go
// 使用单向 channel 增强类型安全
func producer(out chan<- int) {
    for i := 0; i < 10; i++ {
        out <- i
    }
    close(out)
}

func consumer(in <-chan int) {
    for val := range in {
        process(val)
    }
}
```

#### 4. **Context 传递**

```go
// ✅ Context 作为第一个参数
func DoSomething(ctx context.Context, arg string) error {
    // 检查取消
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
    }
    
    // 继续工作
    return nil
}
```

### Python 最佳实践

#### 1. **避免阻塞事件循环**

```python
# ❌ 错误：阻塞事件循环
async def bad():
    time.sleep(1)  # 阻塞！

# ✅ 正确：使用异步版本
async def good():
    await asyncio.sleep(1)

# ✅ 运行阻塞代码在执行器中
async def run_blocking():
    loop = asyncio.get_event_loop()
    result = await loop.run_in_executor(None, blocking_function)
```

#### 2. **正确处理异常**

```python
# 使用 return_exceptions 收集所有结果
results = await asyncio.gather(
    task1(),
    task2(),
    return_exceptions=True
)

for result in results:
    if isinstance(result, Exception):
        handle_error(result)
    else:
        process_result(result)
```

#### 3. **资源管理**

```python
# ✅ 使用异步上下文管理器
async with aiohttp.ClientSession() as session:
    async with session.get(url) as response:
        data = await response.text()

# ✅ 确保清理
try:
    await operation()
finally:
    await cleanup()
```

#### 4. **任务管理**

```python
# ✅ 保持对任务的引用
tasks = set()

def create_task(coro):
    task = asyncio.create_task(coro)
    tasks.add(task)
    task.add_done_callback(tasks.discard)
    return task

# ✅ 取消所有任务
async def shutdown():
    tasks = [t for t in asyncio.all_tasks() if t is not asyncio.current_task()]
    for task in tasks:
        task.cancel()
    await asyncio.gather(*tasks, return_exceptions=True)
```

---

## 性能考虑

### Go

- **优势**：
  - 真正的并行执行
  - 低延迟，高吞吐
  - 适合 CPU 密集型任务

- **适用场景**：
  - 高性能 API 服务
  - 微服务架构
  - 系统编程
  - 网络服务器

### Python

- **优势**：
  - 适合 I/O 密集型任务
  - 丰富的异步库生态
  - 易于集成现有 Python 代码

- **适用场景**：
  - Web 爬虫
  - 异步 API 客户端
  - 实时数据处理
  - 聊天服务器

---

## 进阶主题

### Go 进阶

1. **Channel 内部实现**：理解 hchan 结构
2. **调度器原理**：G-M-P 模型
3. **内存模型**：Happens-Before 保证
4. **性能优化**：逃逸分析、内联

### Python 进阶

1. **事件循环实现**：理解 selector 和回调
2. **协程内部机制**：生成器协议
3. **第三方库**：aiohttp, aioredis, asyncpg
4. **性能调优**：uvloop, 连接池

---

## 学习路径

### 初学者

1. 完成基础示例（01-basic-*）
2. 理解核心概念
3. 运行实战示例
4. 修改代码实验

### 中级

1. 学习高级模式
2. 阅读自定义实现
3. 构建小项目
4. 性能测试和优化

### 高级

1. 深入源码实现
2. 贡献开源项目
3. 设计复杂系统
4. 编写高性能应用

---

## 资源链接

### Go

- [Go 官方文档](https://golang.org/doc/)
- [Effective Go](https://golang.org/doc/effective_go)
- [Go 并发模式](https://go.dev/blog/pipelines)
- [Go 内存模型](https://golang.org/ref/mem)

### Python

- [Python Asyncio 文档](https://docs.python.org/3/library/asyncio.html)
- [PEP 492 - Coroutines](https://www.python.org/dev/peps/pep-0492/)
- [aiohttp 文档](https://docs.aiohttp.org/)
- [Real Python Asyncio](https://realpython.com/async-io-python/)

---

## 贡献

欢迎贡献！请：

1. Fork 本仓库
2. 创建特性分支
3. 提交更改
4. 发起 Pull Request

---

## 许可证

MIT License

---

## 联系方式

如有问题或建议，欢迎提 Issue！

---

**⭐ 如果这个项目对你有帮助，请给个 Star！**

---

## English

# Go Async Deep Dive + Python Asyncio Implementation

> Comprehensive guide to Go concurrency and Python asynchronous programming

### 📚 Table of Contents

- [Introduction](#introduction)
- [Project Structure](#project-structure-en)
- [Go Async Programming](#go-async-programming)
- [Python Asyncio](#python-asyncio-en)
- [Comparison](#comparison)
- [Getting Started](#getting-started-en)
- [Best Practices](#best-practices-en)

---

### Introduction

This project provides comprehensive tutorials on **Go concurrency** and **Python async programming**:

- ✨ Complete guide to Go async programming
- 🐍 Deep dive into Python Asyncio with custom implementation
- 🚀 Rich practical examples
- 📊 Detailed comparison analysis
- 💡 Production-ready design patterns

### Project Structure (EN)

See Chinese section above for directory structure.

### Features

#### Go Async
- Goroutines and channels
- Select statements
- Context package
- Sync primitives
- Advanced patterns (worker pools, pipelines, fan-out/fan-in)
- Practical examples (HTTP servers, data processing)

#### Python Asyncio
- Coroutines and event loop
- Tasks and futures
- Streams and protocols
- Custom event loop implementation
- Practical examples (async HTTP, concurrent tasks)

### Getting Started (EN)

#### Run Go Examples

```bash
cd go-async/01-basic-goroutines
go run main.go
```

#### Run Python Examples

```bash
cd python-asyncio/01-asyncio-basics
python basic_coroutines.py
```

### Best Practices (EN)

See detailed best practices in the Chinese section above.

### License

MIT License

---

**⭐ Star this repo if you find it helpful!**

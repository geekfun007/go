# 快速开始 / Quick Start Guide

> 5 分钟开始使用本项目 / Get started in 5 minutes

---

## 🚀 立即开始 / Start Now

### 第一步：验证环境 / Step 1: Check Environment

```bash
# 检查 Go 版本 (需要 >= 1.16)
go version

# 检查 Python 版本 (需要 >= 3.7)
python3 --version
```

### 第二步：运行第一个示例 / Step 2: Run First Example

#### Go 示例

```bash
# 进入目录
cd go-async/01-basic-goroutines

# 运行
go run main.go
```

**你会看到**:
```
Go 异步编程 - 基础 Goroutine 示例
=====================================

=== 示例1: 简单的 Goroutine ===
这是一个 goroutine

=== 示例2: 多个 Goroutine ===
Goroutine 0 正在运行
Goroutine 1 正在运行
...
```

#### Python 示例

```bash
# 进入目录
cd python-asyncio/01-asyncio-basics

# 运行
python3 basic_coroutines.py
```

**你会看到**:
```
==================================================
Python Asyncio 基础 - 协程示例
==================================================

=== 示例1: 基本协程 ===
协程开始执行
协程执行完成
返回值: 结果
...
```

---

## 📁 项目结构速览 / Project Structure

```
workspace/
│
├── go-async/              ← Go 异步编程示例
│   ├── 01-basic-goroutines/
│   ├── 02-channels/
│   ├── 03-select-statement/
│   ├── 04-context/
│   ├── 05-sync-package/
│   ├── 06-advanced-patterns/
│   └── 07-practical-examples/
│
├── python-asyncio/        ← Python Asyncio 示例
│   ├── 01-asyncio-basics/
│   ├── 02-event-loop/
│   ├── 03-streams-and-protocols/
│   ├── 04-practical-examples/
│   ├── 05-custom-implementation/
│   └── 06-advanced-examples/
│
└── 文档/
    ├── README.md          ← 主文档（从这里开始）
    ├── EXAMPLES.md        ← 示例索引
    ├── CONTRIBUTING.md    ← 贡献指南
    └── PROJECT_SUMMARY.md ← 项目总结
```

---

## 🎯 5 分钟学习路径 / 5-Minute Learning Path

### 方案 A: Go 并发入门

```bash
# 1. 基础 Goroutines (1 分钟)
cd go-async/01-basic-goroutines && go run main.go

# 2. Channels 通信 (2 分钟)
cd ../02-channels && go run main.go

# 3. Select 多路复用 (2 分钟)
cd ../03-select-statement && go run main.go
```

### 方案 B: Python Asyncio 入门

```bash
# 1. 协程基础 (1 分钟)
cd python-asyncio/01-asyncio-basics && python3 basic_coroutines.py

# 2. 事件循环 (2 分钟)
cd ../02-event-loop && python3 event_loop_basics.py

# 3. 实战应用 (2 分钟)
cd ../04-practical-examples && python3 async_http_client.py
```

### 方案 C: 对比学习

```bash
# 1. 运行 Go Worker Pool
cd go-async/06-advanced-patterns && go run main.go

# 2. 运行 Python 爬虫
cd ../../python-asyncio/06-advanced-examples && python3 async_web_crawler.py

# 3. 对比两者实现方式
```

---

## 📚 推荐阅读顺序 / Recommended Reading Order

### 1️⃣ 新手 (第 1 天)

**上午**:
1. 阅读 [README.md](README.md) 的"基础概念"部分
2. 运行 Go 的 01-03 示例
3. 运行 Python 的 01-02 示例

**下午**:
1. 阅读 [EXAMPLES.md](EXAMPLES.md)
2. 尝试修改示例代码
3. 运行所有基础示例

### 2️⃣ 进阶 (第 2-3 天)

**第 2 天**:
- Go: Context 和 Sync 包 (04, 05)
- Python: 流和协议 (03)
- 阅读最佳实践

**第 3 天**:
- Go: 高级模式 (06)
- Python: HTTP 客户端 (04)
- 实现小项目

### 3️⃣ 高级 (第 4-7 天)

**第 4-5 天**:
- 深入实战示例 (Go 07, Python 05-06)
- 理解内部实现原理
- 性能优化技巧

**第 6-7 天**:
- 构建自己的项目
- 参考高级模式
- 贡献代码

---

## 🎮 互动练习 / Interactive Exercises

### 练习 1: 修改并发数

**Go - Worker Pool**:
```go
// go-async/06-advanced-patterns/main.go
// 找到这一行:
const numWorkers = 3

// 修改为:
const numWorkers = 10

// 观察性能变化
```

**Python - 爬虫**:
```python
# python-asyncio/06-advanced-examples/async_web_crawler.py
# 找到这一行:
crawler = AsyncWebCrawler(max_concurrent=5, max_depth=2)

# 修改为:
crawler = AsyncWebCrawler(max_concurrent=20, max_depth=3)

# 观察行为变化
```

### 练习 2: 添加超时

**Go**:
```go
// 在任何示例中添加超时
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

select {
case result := <-ch:
    // 处理结果
case <-ctx.Done():
    fmt.Println("超时!")
}
```

**Python**:
```python
# 在任何协程中添加超时
try:
    result = await asyncio.wait_for(some_coroutine(), timeout=5.0)
except asyncio.TimeoutError:
    print("超时!")
```

### 练习 3: 添加错误处理

**Go**:
```go
// 使用 recover 捕获 panic
defer func() {
    if r := recover(); r != nil {
        fmt.Println("恢复自:", r)
    }
}()
```

**Python**:
```python
# 使用 try-except 处理异常
try:
    result = await some_operation()
except Exception as e:
    print(f"错误: {e}")
```

---

## 🔧 常用命令 / Common Commands

### 运行所有示例

```bash
# Go 所有示例
./run-go-examples.sh

# Python 所有示例
./run-python-examples.sh
```

### 运行单个模块

```bash
# Go
cd go-async/XX-module-name
go run main.go

# Python
cd python-asyncio/XX-module-name
python3 *.py
```

### 查看帮助

```bash
# Go
go help

# Python
python3 -h
```

---

## 🐛 遇到问题？/ Troubleshooting

### 问题 1: Go 版本太低

**症状**:
```
go: go.mod file indicates go 1.21, but maximum supported version is 1.16
```

**解决**:
```bash
# 升级 Go
# macOS
brew upgrade go

# Linux
sudo snap refresh go --classic

# 或访问 https://golang.org/dl/
```

### 问题 2: Python asyncio 不可用

**症状**:
```
ModuleNotFoundError: No module named 'asyncio'
```

**解决**:
```bash
# 确保 Python >= 3.7
python3 --version

# 如果版本过低，升级 Python
# macOS
brew upgrade python

# Linux
sudo apt-get update
sudo apt-get install python3.11
```

### 问题 3: 导入错误

**症状**:
```
ImportError: cannot import name ...
```

**解决**:
```bash
# 安装依赖
cd python-asyncio
pip3 install -r requirements.txt
```

---

## 💡 快速技巧 / Quick Tips

### Go Tips

1. **总是使用 defer**
   ```go
   defer wg.Done()
   defer mu.Unlock()
   defer cancel()
   ```

2. **传递变量给 goroutine**
   ```go
   // ❌ 错误
   for i := 0; i < 10; i++ {
       go func() { fmt.Println(i) }()
   }
   
   // ✅ 正确
   for i := 0; i < 10; i++ {
       go func(id int) { fmt.Println(id) }(i)
   }
   ```

3. **检查 channel 是否关闭**
   ```go
   val, ok := <-ch
   if !ok {
       // channel 已关闭
   }
   ```

### Python Tips

1. **总是使用 await**
   ```python
   # ❌ 错误
   async def bad():
       asyncio.sleep(1)  # 缺少 await
   
   # ✅ 正确
   async def good():
       await asyncio.sleep(1)
   ```

2. **使用 async with**
   ```python
   async with aiohttp.ClientSession() as session:
       async with session.get(url) as response:
           data = await response.text()
   ```

3. **处理所有异常**
   ```python
   results = await asyncio.gather(
       *tasks,
       return_exceptions=True
   )
   ```

---

## 📖 下一步学习 / Next Steps

### 1. 深入阅读

- 📘 [README.md](README.md) - 完整教程
- 📗 [EXAMPLES.md](EXAMPLES.md) - 示例索引
- 📙 [PROJECT_SUMMARY.md](PROJECT_SUMMARY.md) - 项目总结

### 2. 实践项目

建议实现以下项目：

**初级**:
- 简单的 HTTP 服务器
- 并发文件处理器
- 任务队列

**中级**:
- 网络爬虫
- API 代理服务
- 实时聊天服务器

**高级**:
- 分布式任务调度器
- 微服务网关
- 流数据处理系统

### 3. 深入学习资源

**Go**:
- [Effective Go](https://golang.org/doc/effective_go)
- [Go Blog - Concurrency](https://blog.golang.org/pipelines)
- [Go by Example](https://gobyexample.com/)

**Python**:
- [Python Asyncio Docs](https://docs.python.org/3/library/asyncio.html)
- [Real Python - Async IO](https://realpython.com/async-io-python/)
- [PEP 492](https://www.python.org/dev/peps/pep-0492/)

---

## 🎓 学习检查清单 / Learning Checklist

### 基础 (必须掌握)

- [ ] 理解 Goroutine / 协程的概念
- [ ] 会使用 Channel / asyncio.Queue
- [ ] 理解并发 vs 并行
- [ ] 会处理超时
- [ ] 会处理错误

### 进阶 (应该掌握)

- [ ] Context / Event Loop 的使用
- [ ] 同步原语 (Mutex, Lock)
- [ ] Select / asyncio.wait
- [ ] 并发模式 (Worker Pool, Pipeline)
- [ ] 资源管理和清理

### 高级 (深入理解)

- [ ] 调度器原理
- [ ] 内存模型
- [ ] 性能优化
- [ ] 分布式模式
- [ ] 生产环境最佳实践

---

## 🌟 社区与贡献 / Community

### 参与方式

1. ⭐ Star 本项目
2. 🐛 报告 Bug
3. 💡 提出建议
4. 📝 改进文档
5. 💻 贡献代码

详见 [CONTRIBUTING.md](CONTRIBUTING.md)

---

## 📞 获取帮助 / Get Help

### 有问题？

1. 🔍 搜索现有 Issues
2. 📖 查阅文档
3. 💬 创建新 Issue
4. 📧 联系维护者

---

**🎉 准备好了吗？开始你的异步编程之旅！**

**Ready? Start your async programming journey!**

---

**快速链接 / Quick Links**:
- [README](README.md) - 主文档
- [EXAMPLES](EXAMPLES.md) - 示例索引
- [CONTRIBUTING](CONTRIBUTING.md) - 贡献指南
- [Go 第一个示例](go-async/01-basic-goroutines/main.go)
- [Python 第一个示例](python-asyncio/01-asyncio-basics/basic_coroutines.py)

---

最后更新: 2025-12-05

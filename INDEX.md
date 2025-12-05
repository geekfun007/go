# 项目导航 / Project Navigation

> 快速找到你需要的内容 / Quickly find what you need

---

## 🚀 快速入口 / Quick Access

### 新手入门
1. 📖 [快速开始](QUICK_START.md) - 5分钟上手
2. 💡 [核心思路](ASYNC_CORE_CONCEPTS.md) - 理解异步本质 ⚡新增
3. 📚 [主文档](README.md) - 完整教程
4. 📋 [示例索引](EXAMPLES.md) - 查找示例

### 经验开发者
1. 🔄 [模式对比](PATTERNS_COMPARISON.md) - 深度对比分析 ⚡新增
2. 🔍 [示例索引](EXAMPLES.md) - 按主题查找
3. 💡 [最佳实践](README.md#最佳实践) - 生产环境指南
4. 🎯 [高级示例](EXAMPLES.md#高级--advanced) - 复杂应用

### 贡献者
1. 📝 [贡献指南](CONTRIBUTING.md) - 如何贡献
2. 📊 [项目总结](PROJECT_SUMMARY.md) - 项目概况
3. ✅ [完成报告](COMPLETION_REPORT.md) - 实现详情

---

## 📂 目录结构 / Directory Structure

```
/workspace/
│
├── 📚 文档 / Documentation
│   ├── README.md              ⭐ 主文档（从这里开始）
│   ├── QUICK_START.md         🚀 快速开始指南
│   ├── ASYNC_CORE_CONCEPTS.md 💡 异步核心思路（新增！）
│   ├── PATTERNS_COMPARISON.md 🔄 模式深度对比（新增！）
│   ├── EXAMPLES.md            📋 示例索引
│   ├── CONTRIBUTING.md        🤝 贡献指南
│   ├── PROJECT_SUMMARY.md     📊 项目总结
│   ├── COMPLETION_REPORT.md   ✅ 完成报告
│   └── INDEX.md               📍 本文件
│
├── 💻 Go 异步编程 / Go Async
│   ├── 01-basic-goroutines/   ⭐ 入门：基础 Goroutines
│   ├── 02-channels/           ⭐ 入门：Channel 详解
│   ├── 03-select-statement/   ⭐⭐ 进阶：Select 语句
│   ├── 04-context/            ⭐⭐ 进阶：Context 上下文
│   ├── 05-sync-package/       ⭐⭐ 进阶：Sync 包
│   ├── 06-advanced-patterns/  ⭐⭐⭐ 高级：并发模式
│   └── 07-practical-examples/ ⭐⭐⭐ 高级：实战应用
│
├── 🐍 Python Asyncio
│   ├── 01-asyncio-basics/     ⭐ 入门：Asyncio 基础
│   ├── 02-event-loop/         ⭐ 入门：事件循环
│   ├── 03-streams-and-protocols/ ⭐⭐ 进阶：流和协议
│   ├── 04-practical-examples/ ⭐⭐ 进阶：HTTP 客户端
│   ├── 05-custom-implementation/ ⭐⭐⭐ 高级：自定义实现
│   └── 06-advanced-examples/  ⭐⭐⭐ 高级：网络爬虫
│
└── 🛠️ 工具 / Tools
    ├── run-go-examples.sh     运行所有 Go 示例
    ├── run-python-examples.sh 运行所有 Python 示例
    ├── go.mod                 Go 模块配置
    └── requirements.txt       Python 依赖
```

---

## 🎯 按需求查找 / Find by Need

### 我想学习...

#### Go 并发基础
→ [01-basic-goroutines](go-async/01-basic-goroutines/)  
→ [02-channels](go-async/02-channels/)

#### Go 高级特性
→ [04-context](go-async/04-context/)  
→ [05-sync-package](go-async/05-sync-package/)

#### Go 设计模式
→ [06-advanced-patterns](go-async/06-advanced-patterns/)

#### Go 实战应用
→ [07-practical-examples](go-async/07-practical-examples/)

#### Python 异步基础
→ [01-asyncio-basics](python-asyncio/01-asyncio-basics/)  
→ [02-event-loop](python-asyncio/02-event-loop/)

#### Python 网络编程
→ [03-streams-and-protocols](python-asyncio/03-streams-and-protocols/)

#### Python 实战应用
→ [04-practical-examples](python-asyncio/04-practical-examples/)  
→ [06-advanced-examples](python-asyncio/06-advanced-examples/)

#### 理解底层原理
→ [05-custom-implementation](python-asyncio/05-custom-implementation/)

---

## 🔍 按主题查找 / Find by Topic

### 并发控制
- Go: [Worker Pool](go-async/06-advanced-patterns/main.go)
- Python: [Connection Pool](python-asyncio/04-practical-examples/async_http_client.py)

### 超时处理
- Go: [Context](go-async/04-context/main.go)
- Python: [wait_for](python-asyncio/01-asyncio-basics/basic_coroutines.py)

### 错误处理
- Go: [Error Group](go-async/06-advanced-patterns/main.go)
- Python: [Exception Handling](python-asyncio/04-practical-examples/async_http_client.py)

### 速率限制
- Go: [Rate Limiting](go-async/06-advanced-patterns/main.go)
- Python: [Rate Limiter](python-asyncio/04-practical-examples/async_http_client.py)

### 数据处理
- Go: [Data Processing](go-async/07-practical-examples/data-processing.go)
- Python: [Web Crawler](python-asyncio/06-advanced-examples/async_web_crawler.py)

---

## 📈 学习路径 / Learning Paths

### 路径 A: 快速入门 (1-2 天)
```
Day 1:
  → QUICK_START.md
  → go-async/01-basic-goroutines
  → python-asyncio/01-asyncio-basics

Day 2:
  → go-async/02-channels
  → python-asyncio/02-event-loop
  → 实践练习
```

### 路径 B: 系统学习 (1-2 周)
```
Week 1:
  → README.md (基础概念)
  → Go: 01, 02, 03
  → Python: 01, 02

Week 2:
  → Go: 04, 05
  → Python: 03, 04
  → 小项目实践
```

### 路径 C: 深入精通 (1 个月)
```
Week 1-2: 基础 (01-03)
Week 3: 进阶 (04-05)
Week 4: 高级 (06-07)
持续: 实战项目
```

---

## 🎓 按难度查找 / Find by Difficulty

### ⭐ 初级 / Beginner
- [Go 01-02](go-async/)
- [Python 01-02](python-asyncio/)
- [QUICK_START.md](QUICK_START.md)

### ⭐⭐ 中级 / Intermediate
- [Go 03-05](go-async/)
- [Python 03-04](python-asyncio/)
- [EXAMPLES.md](EXAMPLES.md)

### ⭐⭐⭐ 高级 / Advanced
- [Go 06-07](go-async/)
- [Python 05-06](python-asyncio/)
- [PROJECT_SUMMARY.md](PROJECT_SUMMARY.md)

---

## 📚 文档说明 / Documentation Guide

### README.md (18KB)
- ✨ 项目介绍
- 📖 完整教程
- 📊 对比分析
- 💡 最佳实践

**适合**: 系统学习

### QUICK_START.md (9.3KB)
- 🚀 5分钟上手
- 🎯 快速路径
- 💻 互动练习
- 🐛 问题排查

**适合**: 快速入门

### EXAMPLES.md (9.3KB)
- 📋 完整索引
- 🔍 主题查找
- 📈 难度分级
- ❓ 常见问题

**适合**: 查找示例

### CONTRIBUTING.md (5.4KB)
- 🤝 贡献流程
- 📝 代码规范
- ✅ 提交指南
- 🌐 行为准则

**适合**: 贡献者

### PROJECT_SUMMARY.md (11KB)
- 📊 项目统计
- 🎯 内容总结
- 🗺️ 知识体系
- 📈 学习路径

**适合**: 全面了解

---

## 🛠️ 使用工具 / Using Tools

### 运行所有示例
```bash
# Go
./run-go-examples.sh

# Python
./run-python-examples.sh
```

### 运行单个示例
```bash
# Go
cd go-async/01-basic-goroutines
go run main.go

# Python
cd python-asyncio/01-asyncio-basics
python3 basic_coroutines.py
```

### 查看统计
```bash
cat PROJECT_STATS.txt
```

---

## 💡 快速问答 / Quick Q&A

**Q: 从哪里开始？**  
A: [QUICK_START.md](QUICK_START.md)

**Q: 如何查找示例？**  
A: [EXAMPLES.md](EXAMPLES.md)

**Q: 如何贡献？**  
A: [CONTRIBUTING.md](CONTRIBUTING.md)

**Q: 项目有多大？**  
A: [PROJECT_SUMMARY.md](PROJECT_SUMMARY.md)

**Q: Go vs Python？**  
A: [README.md#对比分析](README.md#对比分析)

---

## 🔗 相关链接 / Related Links

### 官方文档
- [Go Documentation](https://golang.org/doc/)
- [Python Asyncio](https://docs.python.org/3/library/asyncio.html)

### 推荐资源
- [Effective Go](https://golang.org/doc/effective_go)
- [PEP 492](https://www.python.org/dev/peps/pep-0492/)

---

## 📞 需要帮助？/ Need Help?

1. 📖 查看 [QUICK_START.md](QUICK_START.md)
2. 🔍 搜索 [EXAMPLES.md](EXAMPLES.md)
3. 📚 阅读 [README.md](README.md)
4. 💬 提交 Issue
5. 🤝 查看 [CONTRIBUTING.md](CONTRIBUTING.md)

---

**提示**: 建议先阅读 [QUICK_START.md](QUICK_START.md) 开始 5 分钟快速体验！

**Tip**: Start with [QUICK_START.md](QUICK_START.md) for a 5-minute quick start!

---

更新时间: 2025-12-05

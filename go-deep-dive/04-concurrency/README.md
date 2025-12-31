# Go 并发编程 / Go Concurrency

## 目录 / Contents

1. [Goroutine 基础 / Goroutine Basics](./goroutine.md)
2. [Channel 基础 / Channel Basics](./channel.md)
3. [Select 语句 / Select Statement](./select.md)
4. [同步原语 / Synchronization Primitives](./sync.md) - Mutex, RWMutex, WaitGroup, Once, Cond
5. [并发模式 / Concurrency Patterns](./patterns.md) - Worker Pool, Fan-Out/Fan-In, Pipeline, Rate Limiting
6. [数据竞争与调试 / Data Race & Debugging](./race.md)

## 概述 / Overview

本章节介绍 Go 语言强大的并发特性，包括：

- **Goroutine**: 轻量级线程、启动方式
- **Channel**: 有缓冲/无缓冲通道、通信机制
- **Select**: 多路复用、超时控制
- **同步原语**: Mutex, RWMutex, WaitGroup, Once, Cond, Atomic, sync.Pool
- **并发模式**: Worker Pool, Fan-Out/Fan-In, Pipeline, Cancellation, Error Group, Semaphore, Rate Limiting
- **调试**: 数据竞争检测、性能分析

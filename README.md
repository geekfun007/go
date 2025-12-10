# Go 微服务开发完全指南

> GORM + Hertz + Kitex 详解与实战

本项目提供了 Go 语言微服务开发的完整教程，涵盖 GORM（ORM）、Hertz（HTTP框架）、Kitex（RPC框架）的详细使用说明和实战案例。

---

## 📚 目录

- [项目介绍](#项目介绍)
- [快速开始](#快速开始)
- [教程文档](#教程文档)
- [实战案例](#实战案例)
- [技术栈](#技术栈)
- [学习路径](#学习路径)
- [贡献指南](#贡献指南)

---

## 项目介绍

本项目旨在帮助开发者快速掌握 Go 微服务开发的核心技术栈：

### 🎯 目标

- 深入理解 GORM ORM 框架的使用
- 掌握 Hertz HTTP 框架的开发技巧
- 学习 Kitex RPC 框架的最佳实践
- 构建完整的微服务架构

### ✨ 特点

- **详细的文档**：每个框架都有完整的使用说明
- **命令行详解**：详细介绍 hz 和 kitex CLI 工具
- **Context 深度解析**：深入讲解 Context 和参数传递
- **实战案例**：提供完整的微服务项目示例
- **最佳实践**：总结开发中的注意事项和优化技巧

---

## 快速开始

### 环境要求

- Go 1.19+
- MySQL 8.0+
- Etcd 3.5+
- Docker (可选)

### 安装工具

```bash
# 安装 Hertz 命令行工具
go install github.com/cloudwego/hertz/cmd/hz@latest

# 安装 Kitex 命令行工具
go install github.com/cloudwego/kitex/tool/cmd/kitex@latest

# 安装 Thriftgo
go install github.com/cloudwego/thriftgo@latest

# 验证安装
hz --version
kitex --version
```

### 快速体验

```bash
# 克隆项目
git clone https://github.com/yourusername/go-microservices-guide.git
cd go-microservices-guide

# 启动基础服务（MySQL + Etcd）
docker-compose up -d

# 运行用户服务
cd user-service
go mod tidy
go run main.go

# 运行 API 网关（新终端）
cd api-gateway
go mod tidy
go run main.go

# 测试 API
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"test@example.com","password":"123456"}'
```

---

## 教程文档

### 1. GORM 详解

📖 [GORM详解.md](./GORM详解.md)

**主要内容：**
- ✅ 安装与配置
- ✅ 数据库连接（MySQL、PostgreSQL、SQLite）
- ✅ 模型定义与字段标签
- ✅ CRUD 操作详解
- ✅ 高级查询（预加载、作用域、子查询）
- ✅ 关联关系（一对一、一对多、多对多）
- ✅ 事务处理
- ✅ 钩子函数
- ✅ 性能优化
- ✅ 完整实战案例

**适合读者：**
- Go 初学者
- 需要使用 ORM 的开发者
- 想深入了解 GORM 的进阶用户

---

### 2. Hertz 详解

📖 [Hertz详解.md](./Hertz详解.md)

**主要内容：**
- ✅ 框架简介与特性
- ✅ **hz CLI 命令行详解**
  - hz new - 创建项目
  - hz update - 更新项目
  - hz model - 生成模型
  - 完整参数说明
  - 实战示例
- ✅ **Context 详解**
  - 请求参数获取（路径、Query、Form、Body、Header、Cookie）
  - 文件上传
  - 响应处理（JSON、XML、HTML、文件）
  - Context 存储
- ✅ 路由与中间件
- ✅ 参数绑定与验证
- ✅ 渲染
- ✅ 客户端使用
- ✅ 中间件开发
- ✅ 性能优化
- ✅ 完整的 RESTful API 实战

**适合读者：**
- 需要构建 HTTP API 的开发者
- 从 Gin 迁移到 Hertz 的用户
- 追求高性能的微服务开发者

---

### 3. Kitex 详解

📖 [Kitex详解.md](./Kitex详解.md)

**主要内容：**
- ✅ 框架简介与特性
- ✅ **kitex CLI 命令行详解**
  - 基本用法
  - 完整参数列表
  - Thrift vs Protobuf
  - 多服务项目管理
  - 实战示例脚本
- ✅ **Context 详解**
  - context.Context 使用
  - metainfo 元信息传递
  - 链路传递
  - RPC 信息获取
- ✅ 服务定义与代码生成
  - Thrift IDL 编写规范
  - Protobuf IDL 编写规范
- ✅ 服务端开发
- ✅ 客户端开发
- ✅ 中间件开发
- ✅ 服务治理
  - 服务注册发现（Etcd、Consul）
  - 负载均衡
  - 熔断
  - 限流
- ✅ 性能优化
- ✅ 完整的微服务实战

**适合读者：**
- 需要构建 RPC 服务的开发者
- 从 gRPC 迁移到 Kitex 的用户
- 微服务架构实践者

---

## 实战案例

### 完整微服务架构

📖 [实战案例-完整微服务架构.md](./实战案例-完整微服务架构.md)

**项目架构：**
```
Client (HTTP)
    ↓
API Gateway (Hertz)
    ↓ (RPC)
User Service (Kitex + GORM)
    ↓ (SQL)
MySQL Database
```

**技术栈：**
- **API 网关**：Hertz
- **RPC 服务**：Kitex
- **数据库 ORM**：GORM
- **服务发现**：Etcd
- **数据库**：MySQL
- **认证**：JWT

**功能特性：**
- 用户注册/登录
- JWT 认证鉴权
- 用户信息 CRUD
- 密码加密存储
- 服务注册发现
- RESTful API
- RPC 通信
- 中间件支持
- 错误处理
- 日志记录

**项目结构：**
```
microservices/
├── api-gateway/        # HTTP API 网关
├── user-service/       # 用户 RPC 服务
├── common/             # 公共库
└── docker-compose.yml  # Docker 编排
```

**快速运行：**
```bash
# 1. 启动基础服务
docker-compose up -d

# 2. 运行用户服务
cd user-service && go run main.go

# 3. 运行 API 网关
cd api-gateway && go run main.go

# 4. 测试
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"test@example.com","password":"123456"}'
```

---

## 技术栈

| 技术 | 版本 | 用途 | 官方文档 |
|------|------|------|----------|
| **GORM** | v1.25+ | ORM 框架 | [gorm.io](https://gorm.io) |
| **Hertz** | v0.8+ | HTTP 框架 | [cloudwego.io/hertz](https://www.cloudwego.io/zh/docs/hertz/) |
| **Kitex** | v0.9+ | RPC 框架 | [cloudwego.io/kitex](https://www.cloudwego.io/zh/docs/kitex/) |
| **Etcd** | v3.5+ | 服务发现 | [etcd.io](https://etcd.io) |
| **MySQL** | v8.0+ | 关系型数据库 | [mysql.com](https://www.mysql.com) |
| **JWT** | v5 | 身份认证 | [jwt.io](https://jwt.io) |

---

## 学习路径

### 初学者路径

```
1. Go 语言基础
   └─> 熟悉 Go 语法、并发、接口

2. GORM 学习
   └─> 阅读 GORM详解.md
   └─> 练习 CRUD 操作
   └─> 理解关联关系

3. Hertz 学习
   └─> 阅读 Hertz详解.md
   └─> 学习 hz CLI 工具
   └─> 掌握 Context 使用
   └─> 构建简单的 HTTP API

4. Kitex 学习
   └─> 阅读 Kitex详解.md
   └─> 学习 kitex CLI 工具
   └─> 理解 IDL 定义
   └─> 构建 RPC 服务

5. 综合实战
   └─> 跟随完整微服务架构案例
   └─> 理解服务间通信
   └─> 实现自己的微服务项目
```

### 进阶开发者路径

```
1. 深入框架原理
   └─> 源码阅读
   └─> 性能优化

2. 服务治理
   └─> 服务发现
   └─> 负载均衡
   └─> 熔断限流

3. 可观测性
   └─> 日志收集
   └─> 监控告警
   └─> 链路追踪

4. DevOps
   └─> Docker 容器化
   └─> Kubernetes 部署
   └─> CI/CD 流程
```

---

## 核心知识点对照表

### CLI 工具对比

| 功能 | hz (Hertz) | kitex (Kitex) |
|------|-----------|---------------|
| **创建项目** | `hz new` | `kitex` |
| **更新项目** | `hz update` | 重新运行 `kitex` |
| **IDL 类型** | Thrift, Protobuf | Thrift, Protobuf |
| **代码生成** | Handler, Router, Model | Client, Server, Model |
| **服务类型** | HTTP | RPC |

### Context 对比

| 功能 | Hertz Context | Kitex Context |
|------|---------------|---------------|
| **类型** | `*app.RequestContext` | `context.Context` |
| **请求参数** | `c.Param()`, `c.Query()`, `c.PostForm()` | 通过 RPC 方法参数 |
| **元信息传递** | `c.Set()`, `c.Get()` | `metainfo.WithValue()` |
| **超时控制** | 内置支持 | `context.WithTimeout()` |
| **响应处理** | `c.JSON()`, `c.String()` | 通过返回值 |

---

## 常见问题 FAQ

<details>
<summary><b>Q: GORM、Hertz、Kitex 之间是什么关系？</b></summary>

**A:** 它们是互补的：
- **GORM**：负责数据库操作（持久层）
- **Hertz**：提供 HTTP API 接口（对外服务）
- **Kitex**：提供 RPC 服务（内部服务间通信）

典型架构：`客户端 → Hertz（API网关）→ Kitex（微服务）→ GORM（数据库）`
</details>

<details>
<summary><b>Q: 何时使用 Hertz？何时使用 Kitex？</b></summary>

**A:**
- **Hertz**：对外提供 HTTP RESTful API，如移动端、Web 端调用
- **Kitex**：服务间内部通信，性能要求高的场景

可以同时使用：Hertz 作为 API 网关，Kitex 作为后端微服务。
</details>

<details>
<summary><b>Q: hz 和 kitex 命令有什么区别？</b></summary>

**A:**
- **hz**：生成 Hertz HTTP 服务代码，包括路由、Handler
- **kitex**：生成 Kitex RPC 服务代码，包括 Client、Server

两者都支持 Thrift 和 Protobuf IDL，但生成的代码类型不同。
</details>

<details>
<summary><b>Q: 如何在 Hertz 和 Kitex 之间传递用户信息？</b></summary>

**A:**
- **Hertz**：使用 `c.Set("user_id", id)` 存储用户信息
- **Kitex**：使用 `metainfo.WithValue(ctx, "user_id", id)` 传递元信息
- **跨服务**：在 Hertz 中通过 Kitex 客户端调用时传递 Context
</details>

<details>
<summary><b>Q: 生成的代码可以修改吗？</b></summary>

**A:**
- **不应修改**：`kitex_gen/`、`biz/model/` 等生成的基础代码
- **可以修改**：`handler.go`、`main.go` 等业务逻辑代码
- **更新 IDL 后**：重新运行命令会覆盖生成的代码，但保留 handler 中的业务逻辑
</details>

---

## 项目示例代码

### GORM 快速示例

```go
package main

import (
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

type User struct {
    ID       uint   `gorm:"primarykey"`
    Username string `gorm:"uniqueIndex"`
    Email    string
}

func main() {
    db, _ := gorm.Open(mysql.Open("user:pass@tcp(127.0.0.1:3306)/db"))
    db.AutoMigrate(&User{})
    
    // 创建
    db.Create(&User{Username: "test", Email: "test@example.com"})
    
    // 查询
    var user User
    db.First(&user, "username = ?", "test")
    
    // 更新
    db.Model(&user).Update("Email", "new@example.com")
    
    // 删除
    db.Delete(&user)
}
```

### Hertz 快速示例

```go
package main

import (
    "context"
    "github.com/cloudwego/hertz/pkg/app"
    "github.com/cloudwego/hertz/pkg/app/server"
    "github.com/cloudwego/hertz/pkg/common/utils"
)

func main() {
    h := server.Default()
    
    h.GET("/ping", func(ctx context.Context, c *app.RequestContext) {
        c.JSON(200, utils.H{"message": "pong"})
    })
    
    h.POST("/users", func(ctx context.Context, c *app.RequestContext) {
        var req struct {
            Username string `json:"username" binding:"required"`
        }
        c.BindAndValidate(&req)
        c.JSON(201, utils.H{"username": req.Username})
    })
    
    h.Spin()
}
```

### Kitex 快速示例

```go
// user.thrift
namespace go user
struct User { 1: i64 id, 2: string name }
service UserService { User GetUser(1: i64 id) }

// 生成代码
// kitex -module yourmodule -service userservice user.thrift

// handler.go
func (s *UserServiceImpl) GetUser(ctx context.Context, id int64) (*user.User, error) {
    return &user.User{Id: id, Name: "test"}, nil
}

// main.go
func main() {
    svr := userservice.NewServer(new(UserServiceImpl))
    svr.Run()
}

// client.go
func main() {
    cli, _ := userservice.NewClient("userservice", client.WithHostPorts("127.0.0.1:8888"))
    resp, _ := cli.GetUser(context.Background(), 1)
    fmt.Printf("User: %+v\n", resp)
}
```

---

## 资源链接

### 官方文档
- [GORM 官方文档](https://gorm.io/zh_CN/docs/)
- [Hertz 官方文档](https://www.cloudwego.io/zh/docs/hertz/)
- [Kitex 官方文档](https://www.cloudwego.io/zh/docs/kitex/)
- [CloudWeGo 官网](https://www.cloudwego.io/zh/)

### 社区资源
- [GORM GitHub](https://github.com/go-gorm/gorm)
- [Hertz GitHub](https://github.com/cloudwego/hertz)
- [Kitex GitHub](https://github.com/cloudwego/kitex)
- [CloudWeGo GitHub](https://github.com/cloudwego)

### 相关教程
- [Go 语言官方教程](https://go.dev/doc/tutorial/)
- [Thrift IDL 语法](https://thrift.apache.org/docs/idl)
- [Protobuf 语言指南](https://protobuf.dev/programming-guides/proto3/)

---

## 贡献指南

欢迎提交 Issue 和 Pull Request！

### 如何贡献

1. Fork 本仓库
2. 创建你的特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交你的修改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启一个 Pull Request

### 贡献内容

- 修正文档错误
- 补充示例代码
- 添加新的实战案例
- 优化项目结构
- 翻译文档

---

## 许可证

本项目采用 MIT 许可证。详见 [LICENSE](LICENSE) 文件。

---

## 联系方式

- **Issues**：[提交问题](https://github.com/yourusername/go-microservices-guide/issues)
- **Discussions**：[参与讨论](https://github.com/yourusername/go-microservices-guide/discussions)

---

## 更新日志

### v1.0.0 (2025-12-10)
- ✅ 完成 GORM 详解文档
- ✅ 完成 Hertz 详解文档（包含 hz CLI 和 Context 详解）
- ✅ 完成 Kitex 详解文档（包含 kitex CLI 和 Context 详解）
- ✅ 添加完整微服务架构实战案例
- ✅ 提供 Docker Compose 环境配置

---

## ⭐ Star History

如果这个项目对你有帮助，请给个 Star ⭐️

[![Star History Chart](https://api.star-history.com/svg?repos=yourusername/go-microservices-guide&type=Date)](https://star-history.com/#yourusername/go-microservices-guide&Date)

---

<div align="center">

**[⬆ 回到顶部](#go-微服务开发完全指南)**

Made with ❤️ by Go Developers

</div>

# 完整微服务项目

基于 GORM + Hertz + Kitex 的完整微服务架构实现。

> **最新版本：v2.0.0** - 已集成 Hertz IDL 注解，支持自动生成路由和参数校验

## 📖 版本说明

### v2.0.0 新特性 🎉

1. **✅ IDL 增强** - 添加 Hertz API 注解和参数校验
2. **✅ 自动生成** - 使用 `hz` CLI 自动生成 Handler 和 Router
3. **✅ 参数验证** - 内置参数校验，无需手动验证
4. **✅ 完整文档** - 新增 IDL 注解详解文档

详见 [CHANGELOG.md](./CHANGELOG.md)

---

## 项目架构

```
┌─────────────┐
│   Client    │
└──────┬──────┘
       │ HTTP
       ▼
┌─────────────────────────────────────────┐
│     API Gateway (Hertz)                 │
│  - HTTP RESTful API                     │
│  - hz 自动生成路由                      │
│  - IDL 参数校验                         │
│  - JWT 认证                             │
│  - 请求转发                             │
└──────┬──────────────────────────────────┘
       │ RPC (Thrift)
       ▼
┌─────────────────────────────────────────┐
│     User Service (Kitex)                │
│  - 用户管理                             │
│  - 业务逻辑                             │
│  - 数据持久化（GORM + MySQL）           │
└──────┬──────────────────────────────────┘
       │ SQL
       ▼
┌─────────────────────────────────────────┐
│          MySQL (GORM)                   │
└─────────────────────────────────────────┘

    Service Discovery (Etcd)
```

---

## 目录结构

```
microservices-project/
├── idl/                          # IDL 接口定义（带 Hertz 注解）
│   ├── common.thrift             # 公共定义
│   └── user.thrift               # 用户服务接口（含路由和校验注解）
│
├── api-gateway/                  # API 网关 (Hertz)
│   ├── biz/
│   │   ├── handler/              # HTTP 处理器（hz 自动生成）
│   │   │   └── user/
│   │   │       └── user_service.go
│   │   ├── model/                # 数据模型（hz 自动生成）
│   │   │   └── user/
│   │   ├── middleware/           # 中间件
│   │   │   ├── auth.go
│   │   │   ├── cors.go
│   │   │   ├── logger.go
│   │   │   └── recovery.go
│   │   ├── router/               # 路由配置（hz 自动生成）
│   │   │   └── user/
│   │   └── client/               # RPC 客户端
│   │       └── user_client.go
│   ├── conf/
│   ├── pkg/
│   ├── router_gen.go             # hz 自动生成的路由注册
│   └── main.go
│
├── user-service/                 # 用户服务 (Kitex + GORM)
│   ├── biz/
│   │   ├── handler/
│   │   ├── service/
│   │   ├── repository/
│   │   └── model/
│   ├── conf/
│   ├── pkg/
│   └── main.go
│
├── scripts/                      # 脚本
│   ├── generate.sh               # 代码生成脚本（Kitex + Hertz）
│   ├── start.sh
│   └── stop.sh
│
├── docker-compose.yml
├── Makefile                      # 构建脚本
├── IDL_ANNOTATIONS.md            # IDL 注解详解 ⭐新增
├── CHANGELOG.md                  # 更新日志 ⭐新增
└── README.md
```

---

## 🚀 快速开始

### 环境要求

- Go 1.19+
- MySQL 8.0+
- Etcd 3.5+
- Docker & Docker Compose

### 1. 安装工具

```bash
# 方式 1：使用 Make
make init

# 方式 2：手动安装
go install github.com/cloudwego/hertz/cmd/hz@latest
go install github.com/cloudwego/kitex/tool/cmd/kitex@latest
go install github.com/cloudwego/thriftgo@latest
```

### 2. 启动基础服务

```bash
# 启动 MySQL、Etcd、Redis
make docker-up

# 或手动启动
docker-compose up -d
```

### 3. 生成代码

```bash
# 方式 1：使用 Make（推荐）
make gen

# 方式 2：使用脚本
chmod +x scripts/generate.sh
./scripts/generate.sh

# 方式 3：分别生成
make gen-user      # 生成用户服务（Kitex）
make gen-gateway   # 生成网关（Hertz）
```

**生成的内容：**
- ✅ Kitex RPC 服务端代码
- ✅ Hertz HTTP Handler（自动生成）
- ✅ Hertz Router（自动生成）
- ✅ Model 结构体（自动生成）
- ✅ 参数绑定和验证（自动生成）

### 4. 启动服务

```bash
# 方式 1：使用 Make
make run

# 方式 2：手动启动
# 终端 1：用户服务
cd user-service
go mod tidy
go run main.go

# 终端 2：API 网关
cd api-gateway
go mod tidy
go run main.go
```

### 5. 测试 API

```bash
# 注册用户
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123",
    "age": 25
  }'

# 登录
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'

# 获取用户信息
curl http://localhost:8080/api/v1/users/1

# 健康检查
curl http://localhost:8080/health
```

---

## 📚 IDL 注解说明

### API 注解

项目使用 Hertz IDL 注解定义 HTTP 路由和参数校验。

#### 路由注解

```thrift
service UserService {
    // POST 请求
    RegisterResponse Register(1: RegisterRequest req) (
        api.post="/api/v1/auth/register"
    )
    
    // GET 请求（带路径参数）
    GetUserResponse GetUser(1: GetUserRequest req) (
        api.get="/api/v1/users/:id"
    )
    
    // PUT 请求
    UpdateUserResponse UpdateUser(1: UpdateUserRequest req) (
        api.put="/api/v1/users/:id"
    )
    
    // DELETE 请求
    DeleteUserResponse DeleteUser(1: DeleteUserRequest req) (
        api.delete="/api/v1/users/:id"
    )
}
```

#### 参数注解

```thrift
struct RegisterRequest {
    // 请求体参数 + 长度校验
    1: required string Username (
        api.body="username", 
        api.vd="len($) >= 3 && len($) <= 20"
    )
    
    // 请求体参数 + 正则校验
    2: required string Email (
        api.body="email", 
        api.vd="regexp($, '^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$')"
    )
    
    // 请求体参数 + 范围校验
    3: required string Password (
        api.body="password", 
        api.vd="len($) >= 6 && len($) <= 20"
    )
    
    // 可选参数 + 范围校验
    4: optional i32 Age (
        api.body="age", 
        api.vd="$ >= 0 && $ <= 150"
    )
}

struct GetUserRequest {
    // 路径参数
    1: required i64 UserID (
        api.path="id", 
        api.vd="$ > 0"
    )
}

struct ListUsersRequest {
    // 查询参数
    1: optional i32 Page (
        api.query="page", 
        api.vd="$ > 0"
    )
    2: optional i32 PageSize (
        api.query="page_size", 
        api.vd="$ > 0 && $ <= 100"
    )
}
```

### 参数校验规则

| 类型 | 语法 | 示例 |
|------|------|------|
| 字符串长度 | `len($) >= min && len($) <= max` | `len($) >= 6 && len($) <= 20` |
| 数值范围 | `$ >= min && $ <= max` | `$ >= 0 && $ <= 150` |
| 大于 | `$ > value` | `$ > 0` |
| 正则表达式 | `regexp($, pattern)` | `regexp($, '^1[3-9]\\d{9}$')` |
| 组合条件 | `condition1 && condition2` | `len($) > 0 && len($) <= 20` |

**详细说明请查看：** [IDL_ANNOTATIONS.md](./IDL_ANNOTATIONS.md)

---

## 🎯 核心特性

### ✅ 已实现功能

1. **用户注册** 
   - 邮箱/用户名唯一性验证
   - 密码加密（bcrypt）
   - 自动参数校验（IDL 注解）

2. **用户登录**
   - JWT Token 生成
   - 密码验证

3. **用户信息获取**
   - 支持公开访问（无需 token）
   - 支持认证访问（需要 token）

4. **用户信息更新**
   - 权限验证（只能更新自己）
   - 字段部分更新
   - 自动参数校验

5. **用户删除**
   - 软删除
   - 权限验证

6. **用户列表**
   - 分页查询
   - 关键字搜索
   - 自动参数校验

7. **服务治理**
   - Etcd 服务注册发现
   - 健康检查接口

### 🎨 技术亮点

1. **✅ IDL 驱动开发**
   - 使用 Thrift IDL 定义接口
   - 添加 Hertz API 注解
   - 自动生成代码

2. **✅ 自动化代码生成**
   - `hz` 自动生成 Handler
   - `hz` 自动生成 Router
   - `kitex` 生成 RPC 代码

3. **✅ 内置参数校验**
   - IDL 中定义校验规则
   - 自动进行参数验证
   - 减少手动校验代码

4. **✅ 清晰的分层架构**
   - Handler → Service → Repository → Model
   - 职责明确，易于维护

5. **✅ 完善的中间件**
   - 日志记录
   - CORS 跨域
   - JWT 认证
   - 异常恢复

---

## 📖 文档索引

| 文档 | 说明 |
|------|------|
| [README.md](./README.md) | 项目说明、快速开始 |
| [IDL_ANNOTATIONS.md](./IDL_ANNOTATIONS.md) | ⭐ Hertz IDL 注解详解 |
| [PROJECT_STRUCTURE.md](./PROJECT_STRUCTURE.md) | 项目结构详解 |
| [TESTING.md](./TESTING.md) | API 测试指南 |
| [CHANGELOG.md](./CHANGELOG.md) | ⭐ 更新日志 |

---

## 🔧 Make 命令

```bash
# 初始化
make init          # 安装工具（kitex, hz, thriftgo）

# 代码生成
make gen           # 生成所有代码（推荐）
make gen-user      # 只生成用户服务（Kitex）
make gen-gateway   # 只生成网关（Hertz）
make regen         # 清理后重新生成

# 编译
make build         # 编译所有服务
make build-user    # 编译用户服务
make build-gateway # 编译网关

# 运行
make run           # 运行所有服务
make run-user      # 运行用户服务
make run-gateway   # 运行网关

# Docker
make docker-up     # 启动基础服务
make docker-down   # 停止基础服务

# 其他
make clean         # 清理生成的代码
make test          # 运行测试
make stop          # 停止所有服务
```

---

## 🚀 扩展方向

项目提供了 **10 大扩展方向**，包括 40+ 个具体扩展项：

### 1. 功能扩展
- [ ] 用户头像上传
- [ ] 邮箱验证
- [ ] 手机验证码
- [ ] 忘记密码/重置密码
- [ ] 用户权限管理（RBAC）
- [ ] 用户等级系统
- [ ] 用户标签系统

### 2. 新增服务
- [ ] 文章服务 (Post Service)
- [ ] 评论服务 (Comment Service)
- [ ] 通知服务 (Notification Service)

### 3. 中间件服务
- [ ] Redis 缓存
- [ ] 消息队列（RabbitMQ/Kafka）
- [ ] 文件存储（MinIO/OSS）

### 4. 服务治理
- [ ] 配置中心（Consul/Nacos）
- [ ] 链路追踪（Jaeger/Zipkin）
- [ ] 监控告警（Prometheus + Grafana）
- [ ] 熔断降级（Sentinel）
- [ ] 限流

### 5. 可观测性
- [ ] 日志收集（ELK Stack）
- [ ] 性能监控
- [ ] 业务监控

### 6. 安全增强
- [ ] HTTPS 支持
- [ ] OAuth2.0 第三方登录
- [ ] 双因素认证（2FA）
- [ ] API 签名验证
- [ ] 防 XSS/CSRF

### 7. 性能优化
- [ ] 数据库优化（索引、读写分离、分库分表）
- [ ] 缓存策略（多级缓存）
- [ ] 异步处理
- [ ] CDN 加速

### 8. DevOps
- [ ] Docker 容器化
- [ ] Kubernetes 部署
- [ ] CI/CD 流程
- [ ] 自动化测试

### 9. 前端应用
- [ ] Web 管理后台
- [ ] 移动端 APP
- [ ] 小程序

### 10. 数据分析
- [ ] 用户行为分析
- [ ] 数据报表
- [ ] 实时数据看板

---

## 🎓 学习建议

### 学习路径

```
1. 理解 IDL 注解
   └─> 阅读 IDL_ANNOTATIONS.md
   └─> 查看 idl/user.thrift

2. 代码生成实践
   └─> 运行 make gen
   └─> 查看生成的代码

3. 实现业务逻辑
   └─> 修改 api-gateway/biz/handler
   └─> 实现 RPC 调用

4. 测试 API
   └─> 使用 TESTING.md 中的测试用例
   └─> 理解请求流程

5. 扩展功能
   └─> 添加新的 IDL 接口
   └─> 重新生成代码
   └─> 实现新功能
```

### 关键文件

| 文件 | 说明 | 关注点 |
|------|------|--------|
| `idl/user.thrift` | IDL 定义 | 路由注解、参数校验 |
| `scripts/generate.sh` | 代码生成脚本 | 生成流程 |
| `api-gateway/biz/handler` | HTTP Handler | hz 生成的代码 |
| `api-gateway/router_gen.go` | 路由注册 | hz 自动生成 |
| `user-service/biz/handler` | RPC Handler | 业务逻辑实现 |

---

## ❓ 常见问题

### Q: 如何添加新的 API 接口？

**A:** 
1. 在 `idl/user.thrift` 中添加新的接口定义
2. 添加 Hertz 注解（`api.get`, `api.post` 等）
3. 添加参数校验（`api.vd`）
4. 运行 `make gen` 重新生成代码
5. 在生成的 Handler 中实现业务逻辑

### Q: hz 生成的代码可以修改吗？

**A:**
- ✅ **可以修改**：`biz/handler` 中的业务逻辑
- ❌ **不要修改**：`router_gen.go`、`biz/model` 等自动生成的基础代码
- ⚠️ **注意**：更新 IDL 后重新生成，已修改的 handler 不会被覆盖

### Q: 如何添加参数校验？

**A:** 在 IDL 中使用 `api.vd` 注解：
```thrift
1: required string Username (
    api.body="username",
    api.vd="len($) >= 3 && len($) <= 20"
)
```

### Q: 如何调试生成的代码？

**A:**
1. 查看 `api-gateway/biz/handler` 中生成的 Handler
2. 查看 `api-gateway/router_gen.go` 中的路由注册
3. 运行服务，查看日志输出
4. 使用 Postman 或 curl 测试接口

---

## 技术栈

| 技术 | 版本 | 用途 |
|------|------|------|
| **GORM** | v1.25+ | ORM 框架 |
| **Hertz** | v0.8+ | HTTP 框架 |
| **Kitex** | v0.9+ | RPC 框架 |
| **Etcd** | v3.5+ | 服务发现 |
| **MySQL** | v8.0+ | 数据库 |
| **Redis** | v7+ | 缓存 |
| **JWT** | v5 | 认证 |

---

## 接口列表

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/api/v1/auth/register` | 用户注册 | ❌ |
| POST | `/api/v1/auth/login` | 用户登录 | ❌ |
| GET | `/api/v1/users` | 用户列表 | ❌ |
| GET | `/api/v1/users/:id` | 获取用户 | ❌ |
| GET | `/api/v1/users/me` | 当前用户 | ✅ |
| PUT | `/api/v1/users/:id` | 更新用户 | ✅ |
| DELETE | `/api/v1/users/:id` | 删除用户 | ✅ |
| GET | `/health` | 健康检查 | ❌ |

---

## 贡献指南

欢迎提交 Issue 和 Pull Request！

1. Fork 本仓库
2. 创建特性分支
3. 提交代码
4. 推送到分支
5. 开启 Pull Request

---

## 许可证

MIT License

---

## 联系方式

- **Issues**: [提交问题](https://github.com/yourusername/microservices-project/issues)
- **Email**: your-email@example.com

---

<div align="center">

**⭐ 如果这个项目对你有帮助，请给个 Star ⭐**

Made with ❤️ by Go Developers

</div>

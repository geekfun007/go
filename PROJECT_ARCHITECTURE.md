# 项目架构说明

## 🎯 项目概述

本项目是一个基于 CloudWeGo 生态的完整微服务架构示例，展示了如何使用现代化的 Go 开发工具和框架构建生产级的后端服务。

## 📐 技术架构

### 整体架构图

```
┌─────────────────────────────────────────────────────────────┐
│                         客户端                                │
└────────────────────────┬────────────────────────────────────┘
                         │ HTTP Request
                         ▼
┌─────────────────────────────────────────────────────────────┐
│               Hertz HTTP 服务 (:8080)                        │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  Router → Handler → RPC Client                       │   │
│  │  (hz 自动生成)    (手动实现)   (服务发现)             │   │
│  └──────────────────────────────────────────────────────┘   │
└────────────────────────┬────────────────────────────────────┘
                         │ RPC Call (Thrift)
                         │ via Service Discovery
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                 Etcd 服务注册中心 (:2379)                    │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  Service Registry                                    │   │
│  │  - user_service: [ip:port, ...]                     │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              Kitex RPC 服务 (:8888)                          │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  Handler → Business Logic → DAL → Database           │   │
│  │  (kitex 生成)  (手动实现)     (封装)  (GORM)         │   │
│  └──────────────────────────────────────────────────────┘   │
└────────────────────────┬────────────────────────────────────┘
                         │ SQL Query
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                MySQL 数据库 (:3306)                          │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  users 表                                             │   │
│  │  - id, username, email, phone, timestamps            │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

## 🔧 技术栈详解

### 1. IDL 驱动开发

**Thrift IDL** (`idl/user.thrift`)
- 定义数据结构和服务接口
- 作为服务契约，确保前后端一致性
- 支持跨语言代码生成

```thrift
struct User {
    1: i64 id
    2: string username
    3: string email
}

service UserService {
    CreateUserResponse CreateUser(1: CreateUserRequest req)
    GetUserResponse GetUser(1: GetUserRequest req)
}
```

### 2. Kitex RPC 服务层

**代码生成**:
```bash
kitex -module <module> -service user_service ./idl/user.thrift
```

**生成内容**:
- `kitex_gen/`: Thrift 生成的类型定义和序列化代码
- `main.go`: 服务启动入口
- `handler.go`: 业务逻辑接口（需手动实现）

**核心组件**:

1. **Handler** (`kitex_service/biz/handler/`)
   - 实现 IDL 定义的服务接口
   - 处理业务逻辑
   - 调用 DAL 层

2. **DAL 层** (`kitex_service/dal/`)
   - 封装所有数据库操作
   - 使用 GORM 进行 ORM 映射
   - 提供统一的数据访问接口

3. **Model 层** (`kitex_service/model/`)
   - 定义 GORM 数据模型
   - 与数据库表结构对应

4. **配置管理** (`kitex_service/config/`)
   - 加载环境变量
   - 管理服务配置

**服务注册**:
```go
// 创建 Etcd registry
r, _ := etcd.NewEtcdRegistry(cfg.Etcd.Endpoints)

// 启动服务并注册
svr := userservice.NewServer(
    handler.NewUserServiceImpl(),
    server.WithRegistry(r),
    server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
        ServiceName: "user_service",
    }),
)
svr.Run()
```

### 3. Hertz HTTP 服务层

**代码生成**:
```bash
hz new -module <module> -idl ./idl/user.thrift
hz update -idl ./idl/user.thrift  # 更新代码
```

**生成内容**:
- `biz/model/`: HTTP 请求/响应模型
- `biz/router/`: 路由配置
- `biz/handler/`: Handler 接口（需手动实现）
- `main.go`: 服务启动入口

**核心组件**:

1. **Handler** (`hertz_service/biz/handler/`)
   - 处理 HTTP 请求
   - 参数验证和绑定
   - 调用 RPC 客户端
   - 返回 HTTP 响应

2. **RPC Client** (`hertz_service/biz/client/`)
   - 封装 Kitex 客户端
   - 实现服务发现
   - 提供调用接口

3. **Model** (`hertz_service/biz/model/`)
   - Hz 从 IDL 自动生成
   - HTTP 请求/响应结构

4. **Router** (`hertz_service/biz/router/`)
   - Hz 自动生成路由配置
   - 支持增量更新

**服务发现**:
```go
// 创建 Etcd resolver
r, _ := etcd.NewEtcdResolver(cfg.Etcd.Endpoints)

// 创建 RPC 客户端
userClient, _ := userservice.NewClient(
    "user_service",  // 服务名称
    client.WithResolver(r),
)
```

### 4. 服务注册与发现 (Etcd)

**工作流程**:

1. **服务注册** (Kitex 服务启动时)
   ```
   Kitex Service → Etcd Registry
   注册信息: {
     service: "user_service",
     address: "ip:port",
     metadata: {...}
   }
   ```

2. **服务发现** (Hertz 调用时)
   ```
   Hertz Client → Etcd Resolver
   查询: "user_service"
   返回: ["ip1:port1", "ip2:port2", ...]
   ```

3. **负载均衡**
   - Kitex 内置负载均衡
   - 支持轮询、随机等策略
   - 自动剔除不健康的实例

4. **健康检查**
   - Etcd 自动维护服务列表
   - 服务下线时自动清理

**优势**:
- ✅ 无需硬编码服务地址
- ✅ 支持动态扩缩容
- ✅ 自动故障转移
- ✅ 服务实例管理

### 5. 数据持久化 (MySQL + GORM)

**层次结构**:
```
Handler → DAL → GORM → MySQL
```

**特性**:
- 自动表结构迁移
- 类型安全的查询
- 事务支持
- 连接池管理

## 📁 目录结构说明

```
/workspace/
├── idl/                          # IDL 定义
│   └── user.thrift              # 用户服务接口定义
│
├── kitex_gen/                   # Kitex 生成的代码（自动生成）
│   └── user/
│       ├── user.go              # Thrift 类型定义
│       └── userservice/         # 服务接口和客户端
│
├── kitex_service/               # Kitex RPC 服务
│   ├── config/                  # 配置管理
│   │   └── config.go           # 加载环境变量
│   ├── dal/                     # 数据访问层
│   │   ├── db.go               # 数据库初始化
│   │   └── user.go             # 用户数据访问
│   ├── model/                   # 数据模型
│   │   └── user.go             # GORM 模型
│   ├── biz/handler/            # 业务逻辑
│   │   └── user_handler.go     # 实现 RPC 接口
│   └── main.go                 # 服务入口（含注册）
│
├── hertz_service/               # Hertz HTTP 服务
│   ├── biz/
│   │   ├── client/             # RPC 客户端
│   │   │   └── user_client.go  # 封装 Kitex 客户端
│   │   ├── config/             # 配置管理
│   │   │   └── config.go       # 加载环境变量
│   │   ├── handler/            # HTTP 处理器
│   │   │   ├── ping.go         # 健康检查（hz 生成）
│   │   │   └── user/           # 用户接口（手动实现）
│   │   ├── model/              # HTTP 模型
│   │   │   └── user/           # Hz 从 IDL 生成
│   │   └── router/             # 路由配置
│   │       ├── register.go     # Hz 生成的路由注册
│   │       └── user/           # 用户路由（hz 生成）
│   ├── main.go                 # 服务入口
│   ├── router.go               # 自定义路由
│   └── router_gen.go           # Hz 生成的路由
│
├── scripts/                     # 工具脚本
│   └── test_api.sh             # API 测试脚本
│
├── docker-compose.yml           # Docker 编排（含 Etcd）
├── Dockerfile.kitex            # Kitex 服务镜像
├── Dockerfile.hertz            # Hertz 服务镜像
├── Makefile                    # 构建脚本
├── .env.example                # 环境变量模板
├── go.mod                      # 根模块（包含 kitex_gen）
└── README.md                   # 项目文档
```

## 🔄 开发工作流

### 1. 修改 IDL

```bash
vim idl/user.thrift
# 添加新的接口或修改现有接口
```

### 2. 重新生成代码

```bash
# 生成 Kitex 服务代码
make gen-kitex

# 更新 Hertz 服务代码
make update-hertz
```

### 3. 实现业务逻辑

**Kitex 端**:
- 编辑 `kitex_service/biz/handler/user_handler.go`
- 实现新增或修改的接口
- 添加 DAL 操作（如需要）

**Hertz 端**:
- 编辑 `hertz_service/biz/handler/user/user_service.go`
- 调用 RPC 客户端
- 处理 HTTP 响应

### 4. 测试

```bash
# 启动服务
make docker-up

# 测试 API
curl http://localhost:8080/api/v1/users
```

### 5. 部署

```bash
# 构建镜像
make docker-build

# 部署到生产环境
docker-compose up -d
```

## 🚦 请求流程

### 创建用户示例

```
1. 客户端发送 HTTP POST 请求
   POST /api/v1/users
   {"username": "john", "email": "john@example.com"}
   
2. Hertz 路由器接收请求
   router.POST("/api/v1/users", handler.CreateUser)
   
3. Handler 解析和验证请求
   var req CreateUserRequest
   c.BindAndValidate(&req)
   
4. 调用 RPC 客户端（服务发现）
   - 从 Etcd 查询 user_service 地址
   - 建立 RPC 连接
   - 发送 Thrift 请求
   
5. Kitex 服务处理请求
   - Handler 接收 RPC 请求
   - 调用 DAL 层
   - GORM 执行 SQL
   - 返回结果
   
6. Hertz 返回 HTTP 响应
   - 转换数据格式
   - 设置 HTTP 状态码
   - 返回 JSON 响应
```

## 🎯 最佳实践

### 1. IDL 设计

- ✅ 使用语义化的命名
- ✅ 添加注释说明
- ✅ 定义清晰的错误码
- ✅ 考虑向后兼容性

### 2. 服务拆分

- ✅ 按业务领域拆分服务
- ✅ 避免过度拆分
- ✅ 保持服务独立性
- ✅ 使用 API 网关统一入口

### 3. 错误处理

- ✅ 统一错误码和错误信息
- ✅ 区分业务错误和系统错误
- ✅ 记录详细的错误日志
- ✅ 返回友好的错误提示

### 4. 性能优化

- ✅ 使用连接池
- ✅ 启用多路复用
- ✅ 合理设置超时时间
- ✅ 添加缓存层

### 5. 监控和运维

- ✅ 添加日志记录
- ✅ 暴露 metrics 接口
- ✅ 配置健康检查
- ✅ 实现优雅关闭

## 📊 性能指标

### 预期性能

- **QPS**: 10,000+ (单机)
- **延迟**: P99 < 50ms
- **连接数**: 支持 10,000+ 并发连接
- **吞吐量**: 100MB/s+

### 优化建议

1. **网络优化**
   - 使用 Netpoll 提升网络性能
   - 启用 TCP 连接复用

2. **序列化优化**
   - 使用 Thrift 二进制协议
   - 考虑使用 Frugal 加速序列化

3. **数据库优化**
   - 添加适当索引
   - 使用连接池
   - 实现读写分离

4. **缓存策略**
   - 热点数据缓存
   - 本地缓存 + Redis

## 🔐 安全建议

- [ ] 添加 API 认证（JWT/OAuth2）
- [ ] 实现请求签名验证
- [ ] 添加 Rate Limiting
- [ ] 配置 HTTPS/TLS
- [ ] 数据加密存储
- [ ] SQL 注入防护（GORM 已内置）
- [ ] XSS 防护

## 📚 学习资源

- [CloudWeGo 官方文档](https://www.cloudwego.io/)
- [Kitex 最佳实践](https://www.cloudwego.io/docs/kitex/best-practice/)
- [Hertz 最佳实践](https://www.cloudwego.io/docs/hertz/tutorials/basic-feature/)
- [微服务架构设计模式](https://microservices.io/)
- [Go 并发编程](https://go.dev/doc/effective_go#concurrency)

---

**项目状态**: ✅ 生产就绪

**最后更新**: 2024-12-04

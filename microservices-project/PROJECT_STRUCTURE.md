# 项目结构详解

本文档详细说明项目的目录结构、文件职责和设计思路。

## 整体架构

```
microservices-project/
├── idl/                    # 接口定义层
├── api-gateway/            # HTTP API 网关
├── user-service/           # 用户 RPC 服务
├── common/                 # 公共库（待扩展）
├── scripts/                # 脚本工具
├── docker-compose.yml      # Docker 编排
├── Makefile               # 构建工具
└── README.md              # 项目说明
```

---

## IDL 层（接口定义）

```
idl/
├── common.thrift          # 公共数据结构
└── user.thrift           # 用户服务接口定义
```

### 设计思路

**职责**：
- 定义服务接口契约
- 定义数据结构
- 作为客户端和服务端的桥梁

**特点**：
- 使用 Thrift IDL 语言
- 语言无关，可生成多种语言代码
- 版本控制，便于接口演进

**文件说明**：

#### common.thrift
- 定义通用数据结构
- 统一响应格式
- 分页相关结构

#### user.thrift
- 用户相关接口定义
- 请求/响应结构体
- RPC 方法声明

---

## User Service（用户服务）

```
user-service/
├── biz/                      # 业务逻辑层
│   ├── handler/              # RPC 处理器（接口实现）
│   │   └── user_handler.go   # 实现 UserService 接口
│   ├── service/              # 业务服务层
│   │   └── user_service.go   # 核心业务逻辑
│   ├── repository/           # 数据访问层
│   │   └── user_repository.go # GORM 数据库操作
│   └── model/                # 数据模型
│       └── user.go           # GORM 模型定义
│
├── conf/                     # 配置管理
│   ├── config.go            # 配置结构和加载
│   └── config.yaml          # 配置文件
│
├── pkg/                      # 工具包
│   ├── db/                  # 数据库工具
│   │   └── mysql.go         # MySQL 连接
│   └── utils/               # 工具函数
│       └── jwt.go           # JWT 生成/解析
│
├── kitex_gen/               # Kitex 生成的代码
│   └── ...                  # （由 kitex 命令生成）
│
├── idl/                     # IDL 文件（链接或复制）
│   └── user.thrift
│
├── main.go                  # 服务入口
├── build.sh                 # 构建脚本
├── go.mod                   # Go 依赖管理
└── go.sum
```

### 分层架构

```
┌─────────────────────────────────────┐
│  Handler Layer (RPC 接口层)         │
│  - 参数验证                         │
│  - 错误处理                         │
│  - 类型转换                         │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│  Service Layer (业务逻辑层)         │
│  - 核心业务逻辑                     │
│  - 事务处理                         │
│  - 业务验证                         │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│  Repository Layer (数据访问层)      │
│  - GORM 操作                        │
│  - 数据库查询                       │
│  - 缓存操作（待扩展）                │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│  Model Layer (数据模型层)           │
│  - GORM 模型定义                    │
│  - 数据库表结构                     │
└─────────────────────────────────────┘
```

### 核心文件详解

#### main.go
**职责**：服务启动入口
```go
- 加载配置
- 初始化数据库
- 自动迁移
- 创建依赖（DI）
- 注册服务到 Etcd
- 启动 Kitex 服务器
```

#### biz/handler/user_handler.go
**职责**：RPC 接口实现
```go
- 实现 UserService 接口
- 参数验证
- 调用 Service 层
- 错误处理和响应封装
- Thrift 类型转换
```

#### biz/service/user_service.go
**职责**：业务逻辑处理
```go
- 核心业务逻辑
- 密码加密
- Token 生成
- 数据验证
- 调用 Repository 层
```

#### biz/repository/user_repository.go
**职责**：数据访问
```go
- GORM 数据库操作
- CRUD 方法
- 查询封装
- 错误转换
```

#### biz/model/user.go
**职责**：数据模型
```go
- GORM 模型定义
- 表结构映射
- 字段标签
- 钩子函数
```

---

## API Gateway（API 网关）

```
api-gateway/
├── biz/                      # 业务逻辑层
│   ├── handler/              # HTTP 处理器
│   │   ├── user.go          # 用户相关接口
│   │   └── ping.go          # 健康检查
│   ├── middleware/           # 中间件
│   │   ├── logger.go        # 日志记录
│   │   ├── cors.go          # CORS 跨域
│   │   ├── auth.go          # JWT 认证
│   │   └── recovery.go      # 异常恢复
│   ├── router/              # 路由配置
│   │   └── router.go        # 路由注册
│   └── client/              # RPC 客户端
│       └── user_client.go   # 用户服务客户端
│
├── conf/                     # 配置管理
│   ├── config.go
│   └── config.yaml
│
├── pkg/                      # 工具包
│   ├── response/            # 统一响应
│   │   └── response.go
│   └── jwt/                 # JWT 工具
│       └── jwt.go
│
├── kitex_gen/               # Kitex 生成的代码
│   └── ...                  # （复制自 user-service）
│
├── main.go                  # 服务入口
├── go.mod
└── go.sum
```

### 请求处理流程

```
┌──────────┐
│  Client  │
└─────┬────┘
      │ HTTP Request
      ▼
┌─────────────────────────────────────┐
│  Middleware Layer                   │
│  1. Recovery (异常恢复)             │
│  2. Logger (日志记录)               │
│  3. CORS (跨域处理)                 │
│  4. Auth (认证 - 可选)              │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│  Router Layer (路由层)              │
│  - 路由匹配                         │
│  - 分组路由                         │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│  Handler Layer (处理器层)           │
│  - 参数绑定验证                     │
│  - 调用 RPC 服务                    │
│  - 响应封装                         │
└──────────────┬──────────────────────┘
               │ RPC Call
               ▼
┌─────────────────────────────────────┐
│  RPC Client (Kitex 客户端)          │
│  - 服务发现                         │
│  - 负载均衡                         │
│  - 调用用户服务                     │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│  User Service (用户服务)            │
└─────────────────────────────────────┘
```

### 核心文件详解

#### main.go
**职责**：HTTP 服务启动
```go
- 加载配置
- 初始化 RPC 客户端
- 创建 Hertz 服务器
- 注册中间件
- 注册路由
- 启动服务
```

#### biz/handler/user.go
**职责**：HTTP 接口处理
```go
- HTTP 请求处理
- 参数绑定和验证
- 调用 RPC 服务
- 响应格式化
- 权限验证
```

#### biz/middleware/auth.go
**职责**：JWT 认证
```go
- 提取 Authorization header
- 解析 JWT token
- 验证 token 有效性
- 存储用户信息到上下文
```

#### biz/middleware/logger.go
**职责**：请求日志
```go
- 记录请求信息
- 记录响应状态
- 计算请求耗时
- 记录客户端 IP
```

#### biz/router/router.go
**职责**：路由配置
```go
- 定义路由规则
- 路由分组
- 中间件绑定
- RESTful 风格
```

#### biz/client/user_client.go
**职责**：RPC 客户端管理
```go
- 初始化 Kitex 客户端
- 服务发现配置
- 单例模式
```

---

## 数据流转

### 注册流程

```
1. Client
   POST /api/v1/auth/register
   Body: {username, email, password}
   ↓

2. API Gateway (Hertz)
   - CORS Middleware
   - Logger Middleware
   - Handler: handler.Register()
   - Validate Request
   ↓

3. RPC Call
   - userClient.Register(ctx, req)
   - Service Discovery (Etcd)
   - Load Balancing
   ↓

4. User Service (Kitex)
   - Handler: RegisterHandler
   - Service: userService.Register()
     * Check if email exists
     * Hash password (bcrypt)
     * Create user record
     * Generate JWT token
   - Repository: userRepo.Create()
   ↓

5. Database (MySQL via GORM)
   INSERT INTO users (...)
   ↓

6. Response
   User Service → API Gateway → Client
   {code: 0, message: "success", data: {user, token}}
```

### 认证流程

```
1. Client
   GET /api/v1/users/me
   Header: Authorization: Bearer {token}
   ↓

2. API Gateway
   - Auth Middleware
     * Extract token from header
     * Validate token (JWT)
     * Get user_id from claims
     * Store user_id in context
   - Handler: handler.GetCurrentUser()
     * Get user_id from context
     * Call RPC
   ↓

3. User Service
   - GetUser(user_id)
   - Query database
   ↓

4. Response
   {user info}
```

---

## 配置管理

### 配置文件结构

**user-service/conf/config.yaml**
```yaml
server:          # 服务配置
  address: "0.0.0.0"
  port: 8888
  service_name: "user-service"

database:        # 数据库配置
  host: "127.0.0.1"
  port: 3306
  user: "root"
  password: "root123"
  database: "userdb"
  max_idle_conns: 10
  max_open_conns: 100

etcd:            # 服务发现配置
  endpoints:
    - "127.0.0.1:2379"

jwt:             # JWT 配置
  secret: "your-secret-key"
  expire_hours: 24

log:             # 日志配置
  level: "info"
  filename: "logs/user-service.log"
```

### 环境变量支持（可扩展）

```go
// 可以添加环境变量覆盖
func LoadConfig() *Config {
    cfg := loadFromYAML()
    
    // 环境变量覆盖
    if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
        cfg.Database.Host = dbHost
    }
    
    return cfg
}
```

---

## 依赖注入

### User Service 依赖关系

```go
main()
  ↓
  db := initDB()              // 初始化数据库
  ↓
  repo := NewRepository(db)   // 创建 Repository
  ↓
  svc := NewService(repo)     // 创建 Service
  ↓
  handler := NewHandler(svc)  // 创建 Handler
  ↓
  server.Run(handler)         // 启动服务器
```

### 优点

- **松耦合**：各层独立，易于测试
- **可维护**：修改一层不影响其他层
- **可测试**：可以 mock 各层依赖
- **可扩展**：易于添加新功能

---

## 错误处理

### 错误传递链

```
Repository Layer
  ↓ 数据库错误
Service Layer
  ↓ 业务错误
Handler Layer
  ↓ RPC 错误码
API Gateway
  ↓ HTTP 状态码
Client
```

### 错误码设计

```go
// 成功
0: success

// 客户端错误 (400-499)
400: bad request
401: unauthorized
403: forbidden
404: not found

// 业务错误 (40000-49999)
40001: email already exists
40002: username already exists
40003: username taken
40004: email taken
40101: invalid password
40401: user not found (login)
40402: user not found (get)
40403: user not found (update)
40404: user not found (delete)

// 服务器错误 (500-599)
500: internal server error
```

---

## 扩展方向

### 1. 添加缓存层

```
user-service/
├── pkg/
│   ├── cache/
│   │   ├── redis.go
│   │   └── user_cache.go
```

```go
// Service Layer
func (s *UserService) GetUser(id uint) (*User, error) {
    // 1. 查缓存
    if user, ok := cache.Get(id); ok {
        return user, nil
    }
    
    // 2. 查数据库
    user, err := s.repo.GetByID(id)
    
    // 3. 写缓存
    cache.Set(id, user)
    
    return user, err
}
```

### 2. 添加消息队列

```
user-service/
├── pkg/
│   ├── mq/
│   │   ├── rabbitmq.go
│   │   └── producer.go
```

```go
// 异步任务
func (s *UserService) Register() error {
    // 创建用户
    user := createUser()
    
    // 发送欢迎邮件（异步）
    mq.Publish("email.welcome", user)
    
    return nil
}
```

### 3. 添加更多服务

```
microservices-project/
├── post-service/      # 文章服务
├── comment-service/   # 评论服务
├── notification-service/ # 通知服务
```

### 4. 添加 API 版本管理

```go
// v1
h.Group("/api/v1", v1Handler)

// v2
h.Group("/api/v2", v2Handler)
```

---

## 最佳实践

### 1. 代码组织

- **按功能分层**：handler、service、repository
- **依赖注入**：构造函数注入
- **接口抽象**：使用接口而非具体实现

### 2. 错误处理

- **自定义错误**：定义业务错误类型
- **错误包装**：使用 `fmt.Errorf` 或 `errors.Wrap`
- **统一响应**：使用统一的响应格式

### 3. 日志记录

- **结构化日志**：使用 JSON 格式
- **日志级别**：DEBUG、INFO、WARN、ERROR
- **关键信息**：请求 ID、用户 ID、耗时

### 4. 测试

- **单元测试**：测试各层逻辑
- **集成测试**：测试服务间调用
- **Mock**：使用 mock 隔离依赖

### 5. 性能优化

- **数据库索引**：合理创建索引
- **连接池**：配置合适的连接池参数
- **缓存策略**：热点数据缓存
- **批量操作**：减少数据库调用次数

---

## 总结

本项目采用清晰的分层架构，各层职责明确：

- **IDL 层**：定义接口契约
- **Handler 层**：处理 RPC/HTTP 请求
- **Service 层**：实现核心业务逻辑
- **Repository 层**：封装数据访问
- **Model 层**：定义数据结构

这种架构具有以下优点：
- 高内聚低耦合
- 易于测试和维护
- 便于扩展新功能
- 符合微服务最佳实践

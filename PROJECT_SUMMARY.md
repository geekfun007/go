# Go 服务端开发实战 - 项目总结

## 🎯 项目概述

本项目是一个完整的 Go 微服务架构示例，演示了如何使用 CloudWeGo 生态构建现代化的后端服务。项目实现了一个用户管理系统，包含完整的 CRUD 操作。

## 📊 技术架构

```
┌─────────────┐
│   客户端     │
└──────┬──────┘
       │ HTTP
       ▼
┌─────────────────────────┐
│  Hertz HTTP 服务        │
│  - RESTful API         │
│  - 请求验证             │
│  - 响应转换             │
└──────┬──────────────────┘
       │ RPC (Thrift)
       ▼
┌─────────────────────────┐
│  Kitex RPC 服务         │
│  - 业务逻辑处理          │
│  - 参数校验              │
│  - 错误处理              │
└──────┬──────────────────┘
       │ DAL
       ▼
┌─────────────────────────┐
│  MySQL 数据库            │
│  - GORM ORM             │
│  - 自动迁移              │
│  - 事务支持              │
└─────────────────────────┘
```

## 🏗️ 项目结构

```
hertz-kitex-demo/
│
├── idl/                          # Thrift IDL 定义
│   └── user.thrift              # 用户服务接口定义
│
├── kitex_service/               # Kitex RPC 服务
│   ├── config/                  # 配置管理
│   │   └── config.go           # 服务配置
│   ├── dal/                     # 数据访问层
│   │   ├── db.go               # 数据库初始化
│   │   └── user.go             # 用户数据访问
│   ├── model/                   # 数据模型
│   │   └── user.go             # 用户模型（GORM）
│   ├── kitex_gen/              # Kitex 生成的代码
│   ├── handler.go              # RPC 业务逻辑
│   └── main.go                 # 服务入口
│
├── hertz_service/               # Hertz HTTP 服务
│   ├── config/                  # 配置管理
│   │   └── config.go           # 服务配置
│   ├── client/                  # RPC 客户端
│   │   └── user_client.go      # 用户服务客户端
│   ├── handler/                 # HTTP 处理器
│   │   └── user_handler.go     # 用户接口处理
│   ├── model/                   # HTTP 模型
│   │   └── user.go             # 请求/响应模型
│   └── main.go                 # 服务入口
│
├── scripts/                     # 工具脚本
│   └── test_api.sh             # API 测试脚本
│
├── docker-compose.yml           # Docker 编排
├── Dockerfile.kitex            # Kitex 服务镜像
├── Dockerfile.hertz            # Hertz 服务镜像
├── Makefile                    # 构建脚本
├── .env.example                # 环境变量模板
├── .gitignore                  # Git 忽略文件
├── go.mod                      # Go 模块定义
└── README.md                   # 项目文档
```

## 🔑 核心功能

### 1. Thrift IDL 定义

- **数据结构**: User, CreateUserRequest, GetUserRequest 等
- **服务接口**: CreateUser, GetUser, UpdateUser, DeleteUser, ListUsers
- **类型安全**: 编译时类型检查

### 2. Kitex RPC 服务

**端口**: 8888

**功能模块**:
- ✅ 用户创建（参数校验、唯一性检查）
- ✅ 用户查询（按 ID 查询）
- ✅ 用户更新（部分字段更新）
- ✅ 用户删除（软删除支持）
- ✅ 用户列表（分页查询）

**核心特性**:
- 完善的错误处理机制
- 统一的响应格式（code, message, data）
- DAL 层封装数据库操作
- GORM 自动迁移表结构

### 3. Hertz HTTP 服务

**端口**: 8080

**API 接口**:
```
GET    /ping                    # 健康检查
POST   /api/v1/users           # 创建用户
GET    /api/v1/users/:id       # 获取用户
PUT    /api/v1/users/:id       # 更新用户
DELETE /api/v1/users/:id       # 删除用户
GET    /api/v1/users           # 用户列表（支持分页）
```

**核心特性**:
- RESTful API 设计
- 请求参数验证
- RPC 客户端封装
- HTTP 状态码映射
- JSON 格式响应

### 4. 数据访问层 (DAL)

**功能**:
- CreateUser: 创建用户记录
- GetUserByID: 根据 ID 查询
- GetUserByUsername: 根据用户名查询
- GetUserByEmail: 根据邮箱查询
- UpdateUser: 更新用户信息
- DeleteUser: 删除用户
- ListUsers: 分页查询用户列表

**特性**:
- 统一错误处理
- Context 传递
- 事务支持（可扩展）

## 📋 API 示例

### 创建用户

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "email": "john@example.com",
    "phone": "13800138000"
  }'
```

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "username": "john_doe",
    "email": "john@example.com",
    "phone": "13800138000",
    "created_at": 1701676800,
    "updated_at": 1701676800
  }
}
```

### 获取用户列表

```bash
curl "http://localhost:8080/api/v1/users?page=1&page_size=10"
```

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "users": [
      {
        "id": 1,
        "username": "john_doe",
        "email": "john@example.com",
        "phone": "13800138000",
        "created_at": 1701676800,
        "updated_at": 1701676800
      }
    ],
    "total": 1
  }
}
```

## 🚀 快速启动

### 方式一：本地运行

```bash
# 1. 启动 MySQL
docker run -d --name mysql \
  -e MYSQL_ROOT_PASSWORD=root \
  -e MYSQL_DATABASE=user_db \
  -p 3306:3306 mysql:8.0

# 2. 配置环境变量
cp .env.example .env

# 3. 启动 Kitex RPC 服务（终端 1）
make run-kitex

# 4. 启动 Hertz HTTP 服务（终端 2）
make run-hertz

# 5. 测试 API
./scripts/test_api.sh
```

### 方式二：Docker Compose

```bash
# 启动所有服务
make docker-up

# 查看日志
make docker-logs

# 停止服务
make docker-down
```

## 🔧 开发命令

```bash
make help          # 显示帮助信息
make build         # 编译服务
make gen-code      # 从 IDL 生成代码
make run-kitex     # 运行 Kitex 服务
make run-hertz     # 运行 Hertz 服务
make test          # 运行测试
make clean         # 清理构建文件
make docker-build  # 构建 Docker 镜像
make docker-up     # 启动 Docker 服务
make docker-down   # 停止 Docker 服务
make tidy          # 整理依赖
```

## 📈 技术亮点

### 1. 微服务架构
- 职责分离：HTTP 网关 + RPC 服务
- 服务间通信：高性能 Thrift 协议
- 可独立部署和扩展

### 2. 代码生成
- IDL 定义接口，自动生成代码
- 类型安全，减少手动编码错误
- 跨语言支持

### 3. 分层设计
- **HTTP 层**: 处理 HTTP 请求，参数验证
- **RPC 层**: 业务逻辑处理，错误处理
- **DAL 层**: 数据访问封装，统一接口

### 4. 数据库设计
- GORM ORM 框架
- 自动迁移表结构
- 索引优化（username, email 唯一索引）
- 时间戳自动管理

### 5. 错误处理
- 统一错误码和错误信息
- 多层错误传递
- 详细的错误日志

### 6. 配置管理
- 环境变量配置
- 默认值支持
- 多环境部署友好

## 🎓 学习要点

### Thrift IDL
```thrift
struct User {
    1: i64 id
    2: string username
    3: string email
}

service UserService {
    CreateUserResponse CreateUser(1: CreateUserRequest req)
}
```

### Kitex Handler
```go
func (s *UserServiceImpl) CreateUser(ctx context.Context, req *user.CreateUserRequest) (*user.CreateUserResponse, error) {
    // 1. 参数校验
    // 2. 业务逻辑
    // 3. DAL 调用
    // 4. 返回响应
}
```

### Hertz Handler
```go
func (h *UserHandler) CreateUser(ctx context.Context, c *app.RequestContext) {
    var req model.CreateUserRequest
    c.BindAndValidate(&req)
    
    rpcResp, err := client.GetUserClient().CreateUser(ctx, rpcReq)
    
    c.JSON(consts.StatusOK, resp)
}
```

### GORM Model
```go
type User struct {
    ID        int64     `gorm:"primaryKey;autoIncrement"`
    Username  string    `gorm:"type:varchar(64);uniqueIndex;not null"`
    Email     string    `gorm:"type:varchar(128);uniqueIndex;not null"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
}
```

## 🔍 调试建议

### 查看日志
```bash
# Kitex 服务日志
# 显示数据库连接、RPC 调用等信息

# Hertz 服务日志
# 显示 HTTP 请求、响应状态等信息
```

### 常见问题

1. **数据库连接失败**: 检查 MySQL 是否启动，配置是否正确
2. **RPC 调用失败**: 确保 Kitex 服务已启动，端口未被占用
3. **端口冲突**: 修改 .env 文件中的端口配置

## 📚 扩展方向

### 功能扩展
- [ ] 添加用户认证（JWT）
- [ ] 添加权限控制（RBAC）
- [ ] 添加缓存层（Redis）
- [ ] 添加消息队列（Kafka）
- [ ] 添加链路追踪（OpenTelemetry）
- [ ] 添加监控告警（Prometheus + Grafana）

### 性能优化
- [ ] 连接池优化
- [ ] 缓存策略
- [ ] 数据库索引优化
- [ ] 批量操作支持
- [ ] 异步处理

### 可靠性提升
- [ ] 服务限流
- [ ] 熔断降级
- [ ] 超时重试
- [ ] 优雅关闭
- [ ] 健康检查

## 📖 参考资料

- [CloudWeGo 官方文档](https://www.cloudwego.io/)
- [Hertz 文档](https://www.cloudwego.io/docs/hertz/)
- [Kitex 文档](https://www.cloudwego.io/docs/kitex/)
- [Thrift 官方文档](https://thrift.apache.org/)
- [GORM 文档](https://gorm.io/)

## ✅ 项目状态

- ✅ Thrift IDL 定义
- ✅ Kitex RPC 服务实现
- ✅ Hertz HTTP 服务实现
- ✅ MySQL 数据库集成
- ✅ GORM ORM 配置
- ✅ DAL 数据访问层
- ✅ Docker 容器化
- ✅ 完整文档
- ✅ API 测试脚本
- ✅ 编译通过

## 🎉 总结

本项目展示了一个完整的 Go 微服务架构实践，涵盖了：
- **IDL 定义**: Thrift 接口定义语言
- **RPC 服务**: Kitex 高性能 RPC 框架
- **HTTP 服务**: Hertz 高性能 HTTP 框架
- **数据库**: MySQL + GORM ORM
- **架构设计**: 分层架构、职责分离
- **工程实践**: Docker 容器化、配置管理、错误处理

这是一个生产级的代码示例，可以作为实际项目的起点。

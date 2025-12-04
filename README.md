# Go 服务端开发实战：Hertz + Kitex + MySQL

这是一个完整的 Go 微服务架构示例项目，展示了如何使用 CloudWeGo 生态的框架构建现代化的后端服务。

## 📋 项目简介

本项目实现了一个完整的用户管理系统，采用微服务架构，包含以下组件：

- **Thrift IDL**: 使用 Thrift 定义服务接口
- **Hertz**: HTTP 服务框架，处理 RESTful API 请求
- **Kitex**: RPC 服务框架，提供高性能的服务间通信
- **MySQL**: 关系型数据库
- **GORM**: ORM 框架，简化数据库操作
- **DAL**: 数据访问层，封装数据库操作

## 🏗️ 架构设计

```
客户端 
   ↓
Hertz HTTP 服务 (端口 8080)
   ↓ (RPC 调用)
Kitex RPC 服务 (端口 8888)
   ↓ (DAL + ORM)
MySQL 数据库 (端口 3306)
```

### 目录结构

```
.
├── idl/                        # Thrift IDL 文件
│   └── user.thrift            # 用户服务接口定义
├── kitex_service/             # Kitex RPC 服务
│   ├── config/                # 配置管理
│   ├── dal/                   # 数据访问层
│   ├── model/                 # 数据模型
│   ├── kitex_gen/            # Kitex 生成的代码
│   ├── handler.go            # 业务逻辑处理
│   └── main.go               # 服务入口
├── hertz_service/             # Hertz HTTP 服务
│   ├── config/                # 配置管理
│   ├── client/                # RPC 客户端
│   ├── handler/               # HTTP 处理器
│   ├── model/                 # HTTP 请求/响应模型
│   └── main.go               # 服务入口
├── docker-compose.yml         # Docker 编排文件
├── Dockerfile.kitex          # Kitex 服务镜像
├── Dockerfile.hertz          # Hertz 服务镜像
├── Makefile                  # 构建脚本
├── .env.example              # 环境变量示例
└── README.md                 # 项目文档
```

## 🚀 快速开始

### 前置要求

- Go 1.21+
- MySQL 8.0+
- Docker & Docker Compose (可选)

### 方式一：本地运行

#### 1. 克隆项目

```bash
git clone <repository-url>
cd hertz-kitex-demo
```

#### 2. 安装依赖

```bash
# 安装 Go 依赖
go mod download

# 安装 Kitex 工具
go install github.com/cloudwego/kitex/tool/cmd/kitex@latest
go install github.com/cloudwego/thriftgo@latest
```

#### 3. 配置环境变量

```bash
cp .env.example .env
# 编辑 .env 文件，修改数据库连接信息
```

#### 4. 启动 MySQL

```bash
# 使用 Docker 启动 MySQL
docker run -d \
  --name mysql \
  -e MYSQL_ROOT_PASSWORD=root \
  -e MYSQL_DATABASE=user_db \
  -p 3306:3306 \
  mysql:8.0
```

#### 5. 启动 Kitex RPC 服务

```bash
# 在终端 1 中运行
source .env
make run-kitex
# 或者
go run ./kitex_service/main.go ./kitex_service/handler.go
```

#### 6. 启动 Hertz HTTP 服务

```bash
# 在终端 2 中运行
source .env
make run-hertz
# 或者
go run ./hertz_service/main.go
```

### 方式二：Docker Compose 运行

```bash
# 构建并启动所有服务
make docker-up
# 或者
docker-compose up -d

# 查看日志
make docker-logs
# 或者
docker-compose logs -f

# 停止服务
make docker-down
# 或者
docker-compose down
```

## 📡 API 接口

### 健康检查

```bash
curl http://localhost:8080/ping
```

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

### 获取用户

```bash
curl http://localhost:8080/api/v1/users/1
```

### 更新用户

```bash
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_updated",
    "email": "john_new@example.com"
  }'
```

### 删除用户

```bash
curl -X DELETE http://localhost:8080/api/v1/users/1
```

### 获取用户列表

```bash
curl "http://localhost:8080/api/v1/users?page=1&page_size=10"
```

## 🔧 开发指南

### 修改 IDL 后重新生成代码

```bash
make gen-code
```

### 编译项目

```bash
make build
```

### 运行测试

```bash
make test
```

### 清理构建文件

```bash
make clean
```

## 📚 技术栈详解

### 1. Thrift IDL

使用 Apache Thrift 定义服务接口，实现跨语言的 RPC 通信。IDL 文件定义了：
- 数据结构（struct）
- 服务接口（service）
- 请求/响应消息

### 2. Kitex RPC 服务

**特点：**
- 高性能的 RPC 框架
- 支持多种序列化协议（Thrift、Protobuf）
- 内置服务治理功能
- 完善的中间件机制

**核心组件：**
- `handler.go`: 实现业务逻辑
- `dal/`: 数据访问层，封装数据库操作
- `model/`: GORM 数据模型

### 3. Hertz HTTP 服务

**特点：**
- 高性能的 HTTP 框架
- 兼容 net/http 标准库
- 丰富的中间件支持
- 优雅的路由设计

**核心组件：**
- `handler/`: HTTP 请求处理器
- `client/`: RPC 客户端封装
- `model/`: HTTP 请求/响应模型

### 4. MySQL + GORM

**特点：**
- 自动迁移数据库表结构
- 类型安全的查询构建器
- 支持事务、钩子等高级特性

**DAL 层职责：**
- 封装所有数据库操作
- 提供统一的错误处理
- 支持分页查询

## 🔍 关键代码说明

### IDL 定义

```thrift
service UserService {
    CreateUserResponse CreateUser(1: CreateUserRequest req)
    GetUserResponse GetUser(1: GetUserRequest req)
    UpdateUserResponse UpdateUser(1: UpdateUserRequest req)
    DeleteUserResponse DeleteUser(1: DeleteUserRequest req)
    ListUsersResponse ListUsers(1: ListUsersRequest req)
}
```

### Kitex Handler

```go
func (s *UserServiceImpl) CreateUser(ctx context.Context, req *user.CreateUserRequest) (*user.CreateUserResponse, error) {
    // 1. 参数校验
    // 2. 业务逻辑处理
    // 3. 调用 DAL 层
    // 4. 返回响应
}
```

### Hertz Handler

```go
func (h *UserHandler) CreateUser(ctx context.Context, c *app.RequestContext) {
    // 1. 解析 HTTP 请求
    // 2. 调用 Kitex RPC 服务
    // 3. 转换响应格式
    // 4. 返回 HTTP 响应
}
```

### DAL 数据访问层

```go
func (d *UserDAL) CreateUser(ctx context.Context, user *model.User) error {
    return DB.WithContext(ctx).Create(user).Error
}
```

## 🌟 最佳实践

1. **分层架构**: HTTP 层、RPC 层、DAL 层职责明确
2. **错误处理**: 统一的错误码和错误信息
3. **配置管理**: 使用环境变量，支持多环境部署
4. **日志记录**: 关键操作记录日志，便于排查问题
5. **数据校验**: 请求参数严格校验
6. **资源管理**: 正确初始化和释放资源

## 📊 性能优化建议

1. **数据库优化**
   - 为常用查询字段添加索引
   - 使用连接池复用连接
   - 合理使用分页避免大量数据查询

2. **RPC 优化**
   - 使用连接池
   - 启用多路复用
   - 合理设置超时时间

3. **缓存策略**
   - 添加 Redis 缓存热点数据
   - 使用本地缓存减少网络开销

## 🐛 故障排查

### 数据库连接失败

```bash
# 检查 MySQL 是否启动
docker ps | grep mysql

# 检查连接信息是否正确
echo $MYSQL_HOST $MYSQL_PORT $MYSQL_USER
```

### RPC 调用失败

```bash
# 检查 Kitex 服务是否启动
netstat -tlnp | grep 8888

# 检查 RPC 地址配置
echo $RPC_HOST $RPC_PORT
```

### 端口被占用

```bash
# 查找占用端口的进程
lsof -i :8080
lsof -i :8888

# 停止进程或更换端口
```

## 📖 相关资料

- [CloudWeGo 官方文档](https://www.cloudwego.io/)
- [Hertz 文档](https://www.cloudwego.io/docs/hertz/)
- [Kitex 文档](https://www.cloudwego.io/docs/kitex/)
- [Thrift 官方文档](https://thrift.apache.org/)
- [GORM 文档](https://gorm.io/)

## 📝 License

MIT License

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

---

如有问题，请提交 Issue 或联系维护者。

# Go 服务端开发实战：Hertz + Kitex + MySQL + Etcd

这是一个完整的 Go 微服务架构示例项目，展示了如何使用 CloudWeGo 生态的框架和工具构建现代化的后端服务，包含服务注册与发现功能。

## 📋 项目简介

本项目实现了一个完整的用户管理系统，采用微服务架构，包含以下组件：

- **Thrift IDL**: 使用 Thrift 定义服务接口
- **Kitex CLI**: 自动生成 RPC 服务代码
- **Hertz Hz CLI**: 自动生成 HTTP 服务代码
- **Kitex**: RPC 服务框架，提供高性能的服务间通信
- **Hertz**: HTTP 服务框架，处理 RESTful API 请求
- **Etcd**: 服务注册与发现中心
- **MySQL**: 关系型数据库
- **GORM**: ORM 框架，简化数据库操作
- **DAL**: 数据访问层，封装数据库操作

## 🏗️ 架构设计

```
客户端 
   ↓ HTTP
Hertz HTTP 服务 (端口 8080)
   ↓ RPC (通过 Etcd 服务发现)
Kitex RPC 服务 (端口 8888)
   ↓ DAL + ORM
MySQL 数据库 (端口 3306)

服务注册与发现: Etcd (端口 2379)
```

### 核心特性

✅ **代码自动生成**
- 使用 `kitex` CLI 从 IDL 生成 RPC 服务代码
- 使用 `hz` CLI 从 IDL 生成 HTTP 服务代码
- 支持 `hz update` 和 `kitex` 增量更新

✅ **服务注册与发现**
- Kitex RPC 服务启动时自动注册到 Etcd
- Hertz HTTP 服务通过 Etcd 自动发现 Kitex 服务
- 支持动态服务发现，无需硬编码服务地址

✅ **微服务架构**
- HTTP 网关与 RPC 服务分离
- 高性能 Thrift 协议通信
- 可独立部署和扩展

### 目录结构

```
.
├── idl/                        # Thrift IDL 文件
│   └── user.thrift            # 用户服务接口定义
├── kitex_gen/                 # Kitex 生成的代码
│   └── user/                  # 用户服务相关代码
├── kitex_service/             # Kitex RPC 服务
│   ├── config/                # 配置管理
│   ├── dal/                   # 数据访问层
│   ├── model/                 # 数据模型
│   ├── biz/handler/          # 业务逻辑处理
│   └── main.go               # 服务入口（含服务注册）
├── hertz_service/             # Hertz HTTP 服务
│   ├── biz/
│   │   ├── client/           # RPC 客户端（含服务发现）
│   │   ├── config/           # 配置管理
│   │   ├── handler/          # HTTP 处理器
│   │   ├── model/            # HTTP 模型（hz 生成）
│   │   └── router/           # 路由配置
│   └── main.go               # 服务入口
├── docker-compose.yml         # Docker 编排（含 Etcd）
├── Dockerfile.kitex          # Kitex 服务镜像
├── Dockerfile.hertz          # Hertz 服务镜像
├── Makefile                  # 构建脚本
├── .env.example              # 环境变量示例
└── README.md                 # 项目文档
```

## 🚀 快速开始

### 前置要求

- Go 1.23+
- MySQL 8.0+
- Etcd 3.5+（或使用 Docker Compose）
- Docker & Docker Compose (推荐)

### 方式一：Docker Compose 运行（推荐）

```bash
# 启动所有服务（MySQL + Etcd + Kitex + Hertz）
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

### 方式二：本地运行

#### 1. 安装工具

```bash
# 安装 Kitex CLI 工具
go install github.com/cloudwego/kitex/tool/cmd/kitex@latest

# 安装 Hertz Hz CLI 工具
go install github.com/cloudwego/hertz/cmd/hz@latest

# 安装 Thriftgo
go install github.com/cloudwego/thriftgo@latest
```

#### 2. 启动依赖服务

```bash
# 启动 MySQL
docker run -d --name mysql \
  -e MYSQL_ROOT_PASSWORD=root \
  -e MYSQL_DATABASE=user_db \
  -p 3306:3306 \
  mysql:8.0

# 启动 Etcd
docker run -d --name etcd \
  -p 2379:2379 \
  -p 2380:2380 \
  -e ETCD_NAME=etcd0 \
  -e ETCD_ADVERTISE_CLIENT_URLS=http://localhost:2379 \
  -e ETCD_LISTEN_CLIENT_URLS=http://0.0.0.0:2379 \
  -e ETCD_INITIAL_ADVERTISE_PEER_URLS=http://localhost:2380 \
  -e ETCD_LISTEN_PEER_URLS=http://0.0.0.0:2380 \
  -e ETCD_INITIAL_CLUSTER=etcd0=http://localhost:2380 \
  quay.io/coreos/etcd:v3.5.9
```

#### 3. 配置环境变量

```bash
cp .env.example .env
# 编辑 .env 文件（如果需要）
```

#### 4. 启动服务

```bash
# 终端 1：启动 Kitex RPC 服务
cd kitex_service
source ../.env
go run .

# 终端 2：启动 Hertz HTTP 服务
cd hertz_service
source ../.env
go run .
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

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "user": {
    "id": 1,
    "username": "john_doe",
    "email": "john@example.com",
    "phone": "13800138000",
    "created_at": 1701676800,
    "updated_at": 1701676800
  }
}
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

当你修改 `idl/user.thrift` 后：

```bash
# 重新生成 Kitex RPC 服务代码
cd /workspace
kitex -module github.com/example/hertz-kitex-demo -service user_service ./idl/user.thrift

# 重新生成 Hertz HTTP 服务代码（在 hertz_service 目录）
cd hertz_service
hz update -idl ../idl/user.thrift
```

### 编译项目

```bash
# 编译 Kitex 服务
cd kitex_service
go build -o ../bin/kitex_service .

# 编译 Hertz 服务
cd hertz_service
go build -o ../bin/hertz_service .
```

### 开发命令

```bash
make help          # 显示帮助信息
make build         # 编译服务
make gen-code      # 从 IDL 生成代码
make docker-build  # 构建 Docker 镜像
make docker-up     # 启动 Docker 服务
make docker-down   # 停止 Docker 服务
make tidy          # 整理依赖
```

## 📚 技术栈详解

### 1. Thrift IDL

使用 Apache Thrift 定义服务接口，实现跨语言的 RPC 通信。

```thrift
service UserService {
    CreateUserResponse CreateUser(1: CreateUserRequest req)
    GetUserResponse GetUser(1: GetUserRequest req)
    // ...
}
```

### 2. Kitex RPC 服务（自动生成）

**代码生成**:
```bash
kitex -module <module-name> -service <service-name> <idl-file>
```

**特点**:
- 自动生成服务端代码框架
- 内置服务注册功能（与 Etcd 集成）
- 高性能的 RPC 通信
- 支持多种序列化协议

**服务注册示例**:
```go
// kitex_service/main.go
r, err := etcd.NewEtcdRegistry(cfg.Etcd.Endpoints)
svr := userservice.NewServer(
    handler.NewUserServiceImpl(),
    server.WithRegistry(r),
    server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
        ServiceName: cfg.RPC.Name,
    }),
)
```

### 3. Hertz HTTP 服务（自动生成）

**代码生成**:
```bash
hz new -module <module-name> -idl <idl-file>
hz update -idl <idl-file>  # 更新代码
```

**特点**:
- 自动生成 HTTP handler 和 model
- 自动生成路由配置
- 支持参数绑定和验证
- 高性能的 HTTP 框架

**服务发现示例**:
```go
// hertz_service/biz/client/user_client.go
r, err := etcd.NewEtcdResolver(cfg.Etcd.Endpoints)
userClient, err = userservice.NewClient(
    cfg.RPC.Name,
    client.WithResolver(r),
)
```

### 4. Etcd 服务注册与发现

**工作流程**:
1. Kitex 服务启动时，向 Etcd 注册服务信息
2. Hertz 服务启动时，从 Etcd 查询 Kitex 服务地址
3. 支持多实例负载均衡
4. 服务下线时自动取消注册

**优势**:
- 无需硬编码服务地址
- 支持服务动态扩缩容
- 自动故障转移
- 分布式一致性保证

### 5. MySQL + GORM + DAL

**特点**:
- 自动迁移数据库表结构
- 类型安全的查询构建器
- 支持事务、钩子等高级特性
- DAL 层封装所有数据库操作

## 🎓 学习要点

### 代码自动生成

**Kitex CLI**:
```bash
# 生成 RPC 服务
kitex -module github.com/example/demo -service user_service ./idl/user.thrift

# 使用已有的 kitex_gen
kitex -module github.com/example/demo -use github.com/example/demo/kitex_gen -service user_service ./idl/user.thrift
```

**Hz CLI**:
```bash
# 新建 HTTP 服务
hz new -module github.com/example/demo -idl ./idl/user.thrift

# 更新 HTTP 服务
hz update -idl ./idl/user.thrift
```

### 服务注册

```go
// 创建 Etcd registry
r, err := etcd.NewEtcdRegistry([]string{"localhost:2379"})

// 创建服务并注册
svr := userservice.NewServer(
    handler,
    server.WithRegistry(r),
    server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
        ServiceName: "user_service",
    }),
)
```

### 服务发现

```go
// 创建 Etcd resolver
r, err := etcd.NewEtcdResolver([]string{"localhost:2379"})

// 创建客户端并启用服务发现
client, err := userservice.NewClient(
    "user_service",  // 服务名称
    client.WithResolver(r),
)
```

## 🔍 故障排查

### Etcd 连接失败

```bash
# 检查 Etcd 是否启动
docker ps | grep etcd

# 测试 Etcd 连接
etcdctl --endpoints=localhost:2379 endpoint health
```

### 服务注册失败

```bash
# 检查 Etcd 中注册的服务
etcdctl --endpoints=localhost:2379 get --prefix /kitex
```

### 数据库连接失败

```bash
# 检查 MySQL 是否启动
docker ps | grep mysql

# 检查连接信息
echo $MYSQL_HOST $MYSQL_PORT $MYSQL_USER
```

### 服务发现失败

```bash
# 查看 Hertz 日志，确认是否连接到 Etcd
# 查看 Kitex 日志，确认服务是否成功注册
```

## 📈 扩展方向

### 功能扩展
- [ ] 添加用户认证（JWT）
- [ ] 添加权限控制（RBAC）
- [ ] 添加缓存层（Redis）
- [ ] 添加消息队列（Kafka）
- [ ] 添加链路追踪（OpenTelemetry）
- [ ] 添加监控告警（Prometheus + Grafana）

### 服务治理
- [ ] 服务限流
- [ ] 熔断降级
- [ ] 超时重试
- [ ] 负载均衡策略
- [ ] 优雅关闭

## 📚 本项目文档

- [IDL_ANNOTATIONS.md](./IDL_ANNOTATIONS.md) - **IDL 注解完整说明**（新增）
  - 字段注解（api.body, api.query, api.path 等）
  - 方法注解（api.get, api.post, api.put 等）
  - 自动路由生成说明
  - 完整示例和最佳实践
- [IDL_VALIDATION.md](./IDL_VALIDATION.md) - IDL 字段定义与校验规则
- [PROJECT_ARCHITECTURE.md](./PROJECT_ARCHITECTURE.md) - 架构设计详解
- [QUICKSTART.md](./QUICKSTART.md) - 快速开始指南
- [FEATURES.md](./FEATURES.md) - 项目特性说明
- [CHANGELOG.md](./CHANGELOG.md) - 更新日志

## 📖 参考资料

### CloudWeGo 官方文档
- [CloudWeGo 官方文档](https://www.cloudwego.io/)
- [Kitex 文档](https://www.cloudwego.io/docs/kitex/)
- [Kitex CLI 工具](https://www.cloudwego.io/docs/kitex/tutorials/code-gen/code_generation/)
- [Hertz 文档](https://www.cloudwego.io/docs/hertz/)
- [Hz CLI 工具](https://www.cloudwego.io/docs/hertz/tutorials/toolkit/toolkit/)
- [**Hertz 注解说明**](https://www.cloudwego.io/zh/docs/hertz/tutorials/toolkit/annotation/) - 本项目参考
- [服务注册与发现](https://www.cloudwego.io/docs/kitex/tutorials/service-governance/service_discovery/)

### 其他技术文档
- [Thrift 官方文档](https://thrift.apache.org/)
- [Etcd 官方文档](https://etcd.io/)
- [GORM 文档](https://gorm.io/)

## ✅ 项目状态

- ✅ Thrift IDL 定义（首字母大写）
- ✅ **IDL 注解支持**（api.get/post/put/delete + api.body/query/path）
- ✅ **自动路由生成**（基于 IDL 注解）
- ✅ Kitex CLI 自动生成 RPC 服务
- ✅ Hz CLI 自动生成 HTTP 服务
- ✅ 参数绑定与校验（go.tag + vd + api 注解）
- ✅ Kitex RPC 服务实现
- ✅ Hertz HTTP 服务实现
- ✅ Etcd 服务注册与发现
- ✅ MySQL 数据库集成
- ✅ GORM ORM 配置
- ✅ DAL 数据访问层
- ✅ Docker 容器化
- ✅ 完整文档
- ✅ 编译通过

## 🎉 总结

本项目展示了一个完整的 Go 微服务架构实践，涵盖了：
- **IDL 驱动开发**: 使用 Thrift 定义接口，自动生成代码
- **代码生成工具**: kitex 和 hz CLI 提高开发效率
- **RPC 服务**: Kitex 高性能 RPC 框架
- **HTTP 服务**: Hertz 高性能 HTTP 框架
- **服务治理**: Etcd 服务注册与发现
- **数据库**: MySQL + GORM ORM
- **架构设计**: 分层架构、职责分离
- **工程实践**: Docker 容器化、配置管理、错误处理

这是一个生产级的代码示例，可以作为实际项目的起点。

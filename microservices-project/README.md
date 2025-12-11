# 完整微服务项目

基于 GORM + Hertz + Kitex 的完整微服务架构实现。

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
│  - JWT 认证                             │
│  - 请求转发                             │
│  - 限流熔断                             │
└──────┬──────────────────────────────────┘
       │ RPC (Thrift)
       ▼
┌─────────────────────────────────────────┐
│     User Service (Kitex)                │
│  - 用户管理                             │
│  - 业务逻辑                             │
│  - 数据持久化                           │
└──────┬──────────────────────────────────┘
       │ SQL
       ▼
┌─────────────────────────────────────────┐
│          MySQL (GORM)                   │
│  - 数据存储                             │
│  - 事务管理                             │
└─────────────────────────────────────────┘

    Service Discovery (Etcd)
```

## 目录结构

```
microservices-project/
├── idl/                          # IDL 接口定义
│   ├── user.thrift               # 用户服务接口定义
│   └── common.thrift             # 公共定义
│
├── api-gateway/                  # API 网关 (Hertz)
│   ├── biz/
│   │   ├── handler/              # HTTP 处理器
│   │   │   ├── user.go           # 用户相关接口
│   │   │   └── ping.go           # 健康检查
│   │   ├── middleware/           # 中间件
│   │   │   ├── auth.go           # JWT 认证
│   │   │   ├── cors.go           # CORS 跨域
│   │   │   ├── logger.go         # 日志记录
│   │   │   └── recovery.go       # 异常恢复
│   │   ├── router/               # 路由配置
│   │   │   └── router.go
│   │   └── client/               # RPC 客户端
│   │       └── user_client.go    # 用户服务客户端
│   ├── conf/                     # 配置文件
│   │   ├── config.go
│   │   └── config.yaml
│   ├── pkg/                      # 工具包
│   │   ├── response/             # 统一响应
│   │   └── jwt/                  # JWT 工具
│   ├── kitex_gen/                # Kitex 生成代码
│   ├── main.go
│   ├── go.mod
│   └── go.sum
│
├── user-service/                 # 用户服务 (Kitex + GORM)
│   ├── biz/
│   │   ├── handler/              # RPC 处理器
│   │   │   └── user_handler.go
│   │   ├── service/              # 业务逻辑层
│   │   │   └── user_service.go
│   │   ├── repository/           # 数据访问层
│   │   │   └── user_repository.go
│   │   └── model/                # 数据模型
│   │       └── user.go
│   ├── conf/                     # 配置文件
│   │   ├── config.go
│   │   └── config.yaml
│   ├── pkg/                      # 工具包
│   │   ├── db/                   # 数据库连接
│   │   └── registry/             # 服务注册
│   ├── kitex_gen/                # Kitex 生成代码
│   ├── idl/                      # IDL 文件
│   ├── main.go
│   ├── build.sh
│   ├── go.mod
│   └── go.sum
│
├── common/                       # 公共库
│   ├── utils/                    # 工具函数
│   │   ├── hash.go               # 密码加密
│   │   └── jwt.go                # JWT 工具
│   ├── constants/                # 常量定义
│   │   └── error_code.go
│   └── go.mod
│
├── scripts/                      # 脚本
│   ├── generate.sh               # 代码生成脚本
│   ├── start.sh                  # 启动脚本
│   └── stop.sh                   # 停止脚本
│
├── docker-compose.yml            # Docker 编排
├── Makefile                      # 构建脚本
└── README.md                     # 项目说明
```

## 快速开始

### 1. 环境要求

- Go 1.19+
- MySQL 8.0+
- Etcd 3.5+
- Docker & Docker Compose (可选)

### 2. 安装工具

```bash
# 安装 Hertz CLI
go install github.com/cloudwego/hertz/cmd/hz@latest

# 安装 Kitex CLI
go install github.com/cloudwego/kitex/tool/cmd/kitex@latest

# 安装 Thriftgo
go install github.com/cloudwego/thriftgo@latest
```

### 3. 启动基础服务

```bash
# 启动 MySQL 和 Etcd
docker-compose up -d

# 等待服务就绪
sleep 5
```

### 4. 生成代码

```bash
# 方式 1：使用脚本
./scripts/generate.sh

# 方式 2：手动生成
cd user-service
kitex -module github.com/example/microservices-project/user-service \
  -service userservice \
  ../idl/user.thrift

cd ../api-gateway
# 复制 kitex_gen 或使用相同的生成命令
```

### 5. 启动服务

```bash
# 方式 1：使用 Makefile
make run

# 方式 2：手动启动
# 终端 1：启动用户服务
cd user-service
go mod tidy
go run main.go

# 终端 2：启动 API 网关
cd api-gateway
go mod tidy
go run main.go
```

### 6. 测试 API

```bash
# 注册用户
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "zhangsan",
    "email": "zhangsan@example.com",
    "password": "password123",
    "phone": "13800138000",
    "age": 25
  }'

# 登录
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "zhangsan@example.com",
    "password": "password123"
  }'

# 获取用户信息（需要 token）
TOKEN="your_jwt_token"
curl -X GET http://localhost:8080/api/v1/users/1 \
  -H "Authorization: Bearer $TOKEN"

# 更新用户信息
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "zhangsan_updated",
    "age": 26
  }'

# 删除用户
curl -X DELETE http://localhost:8080/api/v1/users/1 \
  -H "Authorization: Bearer $TOKEN"
```

## 技术栈

- **API 网关**: Hertz (HTTP 框架)
- **RPC 服务**: Kitex (RPC 框架)
- **ORM**: GORM (数据库 ORM)
- **服务发现**: Etcd
- **数据库**: MySQL
- **认证**: JWT
- **序列化**: Thrift

## 功能特性

### 已实现

- ✅ 用户注册/登录
- ✅ JWT 认证鉴权
- ✅ 用户信息 CRUD
- ✅ 密码加密存储
- ✅ 服务注册发现
- ✅ 统一错误处理
- ✅ 日志记录
- ✅ CORS 支持
- ✅ 健康检查

### 扩展方向

#### 1. 功能扩展
- [ ] 用户头像上传
- [ ] 邮箱验证
- [ ] 手机验证码
- [ ] 忘记密码/重置密码
- [ ] 用户权限管理
- [ ] 用户等级系统
- [ ] 用户标签系统

#### 2. 新增服务
- [ ] 文章服务 (Post Service)
  - 文章发布/编辑/删除
  - 文章分类/标签
  - 文章搜索
- [ ] 评论服务 (Comment Service)
  - 评论发布/回复
  - 评论点赞
- [ ] 通知服务 (Notification Service)
  - 站内消息
  - 邮件通知
  - 推送通知

#### 3. 中间件服务
- [ ] 缓存服务 (Redis)
  - 用户信息缓存
  - 热点数据缓存
  - 分布式锁
- [ ] 消息队列 (RabbitMQ/Kafka)
  - 异步任务处理
  - 事件驱动
- [ ] 文件存储 (MinIO/OSS)
  - 图片上传
  - 文件管理

#### 4. 服务治理
- [ ] 配置中心 (Consul/Nacos)
- [ ] 链路追踪 (Jaeger/Zipkin)
- [ ] 监控告警 (Prometheus + Grafana)
- [ ] 熔断降级 (Sentinel)
- [ ] 限流 (Token Bucket/Leaky Bucket)
- [ ] 灰度发布
- [ ] API 网关增强
  - 请求限流
  - 黑白名单
  - IP 频率限制

#### 5. 可观测性
- [ ] 日志收集 (ELK Stack)
  - Elasticsearch
  - Logstash
  - Kibana
- [ ] 性能监控
  - QPS 统计
  - 响应时间监控
  - 错误率监控
- [ ] 业务监控
  - 用户活跃度
  - 接口调用统计

#### 6. 安全增强
- [ ] HTTPS 支持
- [ ] OAuth2.0 第三方登录
- [ ] 双因素认证 (2FA)
- [ ] API 签名验证
- [ ] 请求加密
- [ ] 防 XSS/CSRF
- [ ] SQL 注入防护
- [ ] 敏感信息脱敏

#### 7. 性能优化
- [ ] 数据库优化
  - 索引优化
  - 读写分离
  - 分库分表
- [ ] 缓存策略
  - 多级缓存
  - 缓存预热
  - 缓存穿透/雪崩防护
- [ ] 异步处理
  - 消息队列
  - 协程池
- [ ] CDN 加速
- [ ] 数据库连接池优化

#### 8. DevOps
- [ ] Docker 容器化
- [ ] Kubernetes 部署
- [ ] CI/CD 流程
  - Jenkins
  - GitLab CI
  - GitHub Actions
- [ ] 自动化测试
  - 单元测试
  - 集成测试
  - 压力测试
- [ ] 灰度发布
- [ ] 蓝绿部署
- [ ] 滚动更新

#### 9. 前端应用
- [ ] Web 管理后台 (Vue/React)
- [ ] 移动端 APP (Flutter/React Native)
- [ ] 小程序

#### 10. 数据分析
- [ ] 用户行为分析
- [ ] 数据报表
- [ ] 实时数据看板
- [ ] 埋点系统

## 接口文档

### 认证相关

#### 注册
- **URL**: `/api/v1/auth/register`
- **Method**: `POST`
- **Body**:
```json
{
  "username": "zhangsan",
  "email": "zhangsan@example.com",
  "password": "password123",
  "phone": "13800138000",
  "age": 25
}
```
- **Response**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "user": {...},
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

#### 登录
- **URL**: `/api/v1/auth/login`
- **Method**: `POST`
- **Body**:
```json
{
  "email": "zhangsan@example.com",
  "password": "password123"
}
```

### 用户管理

#### 获取用户信息
- **URL**: `/api/v1/users/:id`
- **Method**: `GET`
- **Headers**: `Authorization: Bearer {token}`

#### 更新用户信息
- **URL**: `/api/v1/users/:id`
- **Method**: `PUT`
- **Headers**: `Authorization: Bearer {token}`

#### 删除用户
- **URL**: `/api/v1/users/:id`
- **Method**: `DELETE`
- **Headers**: `Authorization: Bearer {token}`

## 开发指南

### 添加新接口

1. 在 `idl/user.thrift` 中定义接口
2. 运行 `kitex` 命令重新生成代码
3. 在 `user-service/biz/handler` 中实现业务逻辑
4. 在 `api-gateway/biz/handler` 中添加 HTTP 接口
5. 在 `api-gateway/biz/router` 中注册路由

### 添加新服务

1. 在 `idl/` 中创建新的 IDL 文件
2. 创建新的服务目录（如 `post-service/`）
3. 使用 `kitex` 生成代码
4. 实现服务逻辑
5. 在 API 网关中集成新服务客户端

### 数据库迁移

```bash
# 进入 user-service
cd user-service

# 运行服务会自动执行 AutoMigrate
go run main.go
```

## 常见问题

### Q: 如何修改数据库配置？

A: 编辑 `user-service/conf/config.yaml` 文件。

### Q: 如何修改服务端口？

A: 
- API 网关端口：编辑 `api-gateway/conf/config.yaml`
- 用户服务端口：编辑 `user-service/conf/config.yaml`

### Q: 如何添加新的中间件？

A: 在 `api-gateway/biz/middleware/` 中创建新文件，然后在 `main.go` 中注册。

### Q: 服务注册失败怎么办？

A: 检查 Etcd 是否正常运行，确保配置中的 Etcd 地址正确。

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT License

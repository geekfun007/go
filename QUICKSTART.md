# 快速开始指南

## 📦 前置准备

确保已安装：
- Docker & Docker Compose
- Go 1.23+（本地开发时）

## 🚀 一分钟快速启动

```bash
# 1. 克隆项目（如果还没有）
git clone <your-repo-url>
cd hertz-kitex-demo

# 2. 启动所有服务
docker-compose up -d

# 3. 等待服务启动（约 30 秒）
docker-compose logs -f

# 4. 测试 API
curl http://localhost:8080/ping
```

## 🧪 API 测试

### 创建用户

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","email":"alice@example.com","phone":"13800138000"}'
```

### 查询用户

```bash
curl http://localhost:8080/api/v1/users/1
```

### 更新用户

```bash
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{"username":"alice_updated"}'
```

### 获取用户列表

```bash
curl "http://localhost:8080/api/v1/users?page=1&page_size=10"
```

### 删除用户

```bash
curl -X DELETE http://localhost:8080/api/v1/users/1
```

## 🛠️ 本地开发

### 1. 安装工具

```bash
make install-tools
```

### 2. 启动依赖服务

```bash
# 只启动 MySQL 和 Etcd
docker-compose up -d mysql etcd
```

### 3. 运行服务

```bash
# 终端 1: 启动 Kitex RPC 服务
make run-kitex

# 终端 2: 启动 Hertz HTTP 服务
make run-hertz
```

## 📝 修改 IDL

```bash
# 1. 修改 IDL
vim idl/user.thrift

# 2. 重新生成代码
make gen-kitex      # 生成 RPC 服务代码
make update-hertz   # 更新 HTTP 服务代码

# 3. 实现业务逻辑
# - kitex_service/biz/handler/user_handler.go
# - hertz_service/biz/handler/user/user_service.go

# 4. 重启服务测试
```

## 🔍 查看服务状态

```bash
# 查看所有服务
docker-compose ps

# 查看日志
docker-compose logs -f kitex_service
docker-compose logs -f hertz_service

# 检查 Etcd 中注册的服务
docker exec -it etcd etcdctl get --prefix /kitex
```

## 🐛 故障排查

### 服务无法启动

```bash
# 查看详细日志
docker-compose logs

# 检查端口占用
netstat -tlnp | grep -E "8080|8888|3306|2379"
```

### 数据库连接失败

```bash
# 检查 MySQL 状态
docker-compose ps mysql

# 重启 MySQL
docker-compose restart mysql
```

### Etcd 连接失败

```bash
# 检查 Etcd 状态
docker exec -it etcd etcdctl endpoint health

# 重启 Etcd
docker-compose restart etcd
```

## 📚 下一步

- 阅读 [README.md](./README.md) 了解完整功能
- 阅读 [PROJECT_ARCHITECTURE.md](./PROJECT_ARCHITECTURE.md) 了解架构设计
- 查看 [API 文档](./docs/api.md)（如果有）

## 🆘 获取帮助

```bash
# 查看所有可用命令
make help
```

遇到问题？请提交 Issue 或查看文档。

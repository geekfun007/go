# 更新日志

## [2.0.0] - 2024-12-04

### 🎉 重大更新

#### 新增
- ✅ **Hz CLI 集成**: 使用 `hz new` 和 `hz update` 自动生成 HTTP 服务代码
- ✅ **Kitex CLI 优化**: 使用 `kitex` 命令自动生成 RPC 服务代码
- ✅ **服务注册与发现**: 集成 Etcd 实现服务治理
  - Kitex 服务自动注册到 Etcd
  - Hertz 服务通过 Etcd 自动发现 Kitex 服务
  - 支持多实例负载均衡
- ✅ **Docker Compose**: 一键启动完整的微服务环境（MySQL + Etcd + Kitex + Hertz）

#### 改进
- 📝 完善的文档体系
  - README.md: 完整的项目说明
  - PROJECT_ARCHITECTURE.md: 架构设计文档
  - QUICKSTART.md: 快速开始指南
  - CHANGELOG.md: 更新日志
- 🔧 优化的 Makefile
  - 添加 `install-tools` 命令
  - 添加 `gen-kitex` 和 `update-hertz` 命令
  - 简化开发流程
- 🐳 改进的 Docker 配置
  - 添加 Etcd 服务
  - 优化构建流程
  - 支持环境变量配置

#### 技术栈
- Go 1.23.0
- Kitex v0.15.2
- Hertz v0.10.3
- Etcd v3.5.9
- MySQL 8.0
- GORM v1.25.5

### 架构变化

**之前**:
```
HTTP → RPC (直连)
```

**现在**:
```
HTTP → Etcd (服务发现) → RPC
```

### 代码生成

**Kitex RPC 服务**:
```bash
kitex -module <module> -service user_service ./idl/user.thrift
```

**Hertz HTTP 服务**:
```bash
hz new -module <module> -idl ./idl/user.thrift
hz update -idl ./idl/user.thrift
```

### 破坏性变更

- ⚠️ 服务启动需要 Etcd（可通过配置禁用）
- ⚠️ Hertz 服务目录结构调整
- ⚠️ 需要 Go 1.23+

### 迁移指南

从 1.x 迁移到 2.0:

1. 更新 Go 版本到 1.23+
2. 安装 Hz 和 Kitex CLI 工具
3. 启动 Etcd 服务
4. 更新环境变量配置
5. 重新生成代码（可选）

## [1.0.0] - 2024-12-04

### 初始版本

- ✅ Thrift IDL 定义
- ✅ Kitex RPC 服务
- ✅ Hertz HTTP 服务  
- ✅ MySQL + GORM + DAL
- ✅ Docker 支持
- ✅ 基础文档

---

**图例**:
- 🎉 重大更新
- ✅ 新增功能
- 📝 文档改进
- 🔧 功能改进
- 🐛 Bug 修复
- ⚠️ 破坏性变更
- 🐳 Docker 相关

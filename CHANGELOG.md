# 更新日志 (Changelog)

所有重要的项目更改都会记录在此文件中。

## [v2.1.0] - 2024-12-05

### ✨ 新增 (Added)
- **IDL 注解支持**: 完整实现 Hertz 注解规范
  - 字段注解：`api.body`, `api.query`, `api.path`, `api.header`, `api.cookie`, `api.form`
  - 方法注解：`api.get`, `api.post`, `api.put`, `api.delete`, `api.patch`, `api.serializer`
- **自动路由生成**: Hz 根据 IDL 注解自动生成 RESTful 路由
  - 支持路径参数（`:id`）
  - 支持查询参数
  - 支持请求体参数
  - 自动生成路由分组和中间件挂载点
- **新增文档**:
  - `IDL_ANNOTATIONS.md`: 完整的 IDL 注解说明文档
    - 字段注解详解
    - 方法注解详解
    - 完整示例代码
    - 最佳实践指南
    - 测试方法
    - 参数绑定优先级说明

### 🔄 变更 (Changed)
- **IDL 定义优化**:
  - 为所有服务方法添加 API 注解
  - 为请求字段添加相应的参数来源注解
  - 统一使用 RESTful 风格路由
- **路由注册方式**:
  - 从手动注册改为 Hz 自动生成
  - 路由定义在 `biz/router/user/user.go`
  - 支持中间件自定义（`middleware.go`）

### 🐛 修复 (Fixed)
- 修复 `middleware.go` 缺少 `app` 包导入的问题
- 确保所有服务编译通过

### 📖 文档 (Documentation)
- 更新 `README.md`：
  - 新增"本项目文档"章节
  - 添加 Hertz 注解官方文档链接
  - 更新项目状态，标注注解支持
- `IDL_ANNOTATIONS.md`：
  - 7 种字段注解详解
  - 5 种方法注解说明
  - 参数绑定优先级
  - 完整的 CRUD API 示例
  - 测试命令示例

### 🎯 路由生成 (Routes Generated)
```
GET    /api/v1/users          → ListUsers
POST   /api/v1/users          → CreateUser
GET    /api/v1/users/:id      → GetUser
PUT    /api/v1/users/:id      → UpdateUser
DELETE /api/v1/users/:id      → DeleteUser
```

### 📝 示例 (Examples)

#### IDL 方法注解
```thrift
service UserService {
    CreateUserResponse CreateUser(1: CreateUserRequest req) (
        api.post = "/api/v1/users",
        api.serializer = "json"
    )
    
    GetUserResponse GetUser(1: GetUserRequest req) (
        api.get = "/api/v1/users/:id",
        api.serializer = "json"
    )
}
```

#### IDL 字段注解
```thrift
struct CreateUserRequest {
    1: required string Username (
        go.tag = "json:\"username\" vd:\"len($) > 0\"",
        api.body = "username"
    )
}

struct GetUserRequest {
    1: required i64 UserID (
        go.tag = "path:\"id\" vd:\"$ > 0\"",
        api.path = "id"
    )
}

struct ListUsersRequest {
    1: required i32 Page (
        go.tag = "query:\"page\" vd:\"$ > 0\"",
        api.query = "page"
    )
}
```

---

## [v2.0.0] - 2024-12-04

### ✨ 新增 (Added)
- **IDL 字段首字母大写**: 符合 Go 语言命名规范
- **参数绑定与校验**: 使用 `go.tag` + `vd` 标签
- **完整文档**:
  - `IDL_VALIDATION.md`: 字段定义与校验规则
  - `FEATURES.md`: 项目特性说明

### 🔄 变更 (Changed)
- IDL 结构体字段从小写改为大写（`id` → `ID`, `username` → `Username`）
- 添加完整的 Hertz 参数校验规则
- 更新所有 Handler 代码以适配新字段名

### 🐛 修复 (Fixed)
- 修复字段名不一致导致的编译错误
- 修复 Hertz 生成代码的导入问题

---

## [v1.0.0] - 2024-12-03

### ✨ 新增 (Added)
- **CLI 工具集成**:
  - Kitex CLI: 自动生成 RPC 服务代码
  - Hz CLI: 自动生成 HTTP 服务代码
- **服务注册与发现**:
  - Etcd 作为注册中心
  - Kitex 服务自动注册
  - Hertz 服务自动发现
- **微服务架构**:
  - Hertz HTTP 服务（端口 8080）
  - Kitex RPC 服务（端口 8888）
  - MySQL 数据库（端口 3306）
  - Etcd 服务（端口 2379）
- **数据层**:
  - GORM ORM 集成
  - DAL 数据访问层
  - 用户 CRUD 操作
- **容器化**:
  - Docker Compose 编排
  - 多阶段 Docker 构建
  - 服务健康检查
- **文档**:
  - `README.md`: 项目概述
  - `PROJECT_ARCHITECTURE.md`: 架构设计
  - `QUICKSTART.md`: 快速开始
  - `CHANGELOG.md`: 更新日志

### 🎯 功能特性
- ✅ Thrift IDL 定义
- ✅ 用户管理 API（创建、查询、更新、删除、列表）
- ✅ RPC 服务调用
- ✅ 数据库持久化
- ✅ 错误处理
- ✅ 日志记录
- ✅ 配置管理

---

## 版本号说明

本项目遵循[语义化版本](https://semver.org/lang/zh-CN/)规范：

- **主版本号 (MAJOR)**: 不兼容的 API 修改
- **次版本号 (MINOR)**: 向下兼容的功能性新增
- **修订号 (PATCH)**: 向下兼容的问题修正

## 变更类型说明

- **✨ 新增 (Added)**: 新功能
- **🔄 变更 (Changed)**: 现有功能的变更
- **🗑️ 废弃 (Deprecated)**: 即将移除的功能
- **🚮 移除 (Removed)**: 已移除的功能
- **🐛 修复 (Fixed)**: Bug 修复
- **🔒 安全 (Security)**: 安全相关的修复
- **📖 文档 (Documentation)**: 文档更新

---

**最后更新**: 2024-12-05

# 更新日志

## v2.0.0 (2025-12-10) - Hertz IDL 注解版本

### 🎉 重大更新

#### IDL 增强
- ✅ **添加 Hertz API 注解**
  - 支持路由注解：`api.get`, `api.post`, `api.put`, `api.delete`
  - 支持参数来源注解：`api.path`, `api.query`, `api.body`, `api.header`
  - 支持参数校验注解：`api.vd`

- ✅ **参数校验规则**
  - 字符串长度验证：`len($) >= min && len($) <= max`
  - 数值范围验证：`$ >= min && $ <= max`
  - 正则表达式验证：`regexp($, pattern)`
  - 组合条件验证：使用 `&&` 和 `||`

#### 代码生成
- ✅ **Hertz CLI 集成**
  - 使用 `hz` 命令自动生成 Handler
  - 自动生成 Router 配置
  - 自动生成 Model 结构体
  - 自动进行参数绑定和验证

- ✅ **生成脚本优化**
  - 新增 `scripts/generate.sh` 一键生成所有代码
  - 支持 Kitex 和 Hertz 代码同时生成
  - 自动处理依赖关系

#### 文档
- ✅ **新增 IDL_ANNOTATIONS.md**
  - 详细的 Hertz 注解说明
  - 参数校验规则完整示例
  - 路由定义最佳实践
  - 常用验证规则参考表

### 📝 IDL 变更

#### 字段命名规范
```diff
- 1: required i64 id
+ 1: required i64 ID
```

#### 参数注解
```diff
struct RegisterRequest {
-   1: required string username
+   1: required string Username (api.body="username", api.vd="len($) > 0 && len($) <= 20")
}
```

#### 路由注解
```diff
service UserService {
-   RegisterResponse Register(1: RegisterRequest req)
+   RegisterResponse Register(1: RegisterRequest req) (api.post="/api/v1/auth/register")
}
```

### 🔧 工具链更新

#### Makefile 命令
```bash
make init          # 安装 kitex, hz, thriftgo
make gen           # 生成所有代码
make gen-user      # 只生成用户服务（Kitex）
make gen-gateway   # 只生成网关（Hertz）
make regen         # 清理后重新生成
```

#### 生成脚本
```bash
./scripts/generate.sh  # 一键生成所有代码
```

### 📊 项目结构变化

#### API Gateway 新增文件
```
api-gateway/
├── biz/
│   ├── handler/
│   │   └── user/           # hz 自动生成
│   │       └── user_service.go
│   ├── model/
│   │   └── user/           # hz 自动生成
│   │       └── user.go
│   └── router/
│       └── user/           # hz 自动生成
│           └── user.go
└── router_gen.go           # hz 自动生成的路由注册
```

### ⚠️ 破坏性变更

1. **字段命名变更**
   - 所有结构体字段改为 PascalCase（首字母大写）
   - JSON 序列化时自动转换为 snake_case

2. **Handler 实现方式变更**
   - 之前：手动实现 Handler
   - 现在：基于 hz 生成的 Handler 框架实现

3. **Router 注册方式变更**
   - 之前：手动注册路由
   - 现在：hz 自动生成路由注册代码

### 🚀 迁移指南

#### 从 v1.0 迁移到 v2.0

1. **更新 IDL 文件**
   ```bash
   # 使用新的 IDL 文件
   cp idl/user.thrift.new idl/user.thrift
   ```

2. **重新生成代码**
   ```bash
   make clean
   make gen
   ```

3. **更新 Handler 实现**
   - 在 hz 生成的 Handler 中实现业务逻辑
   - 调用 RPC 客户端
   - 返回响应

4. **更新 main.go**
   - 使用 hz 生成的路由注册函数

### 📚 新增文档

- **IDL_ANNOTATIONS.md** - Hertz IDL 注解完整指南
- **CHANGELOG.md** - 项目更新日志

### 🐛 Bug 修复

- 修复字段命名不一致问题
- 修复参数验证缺失问题
- 优化代码生成流程

### 🎯 下一步计划

- [ ] 添加更多中间件示例
- [ ] 集成 Redis 缓存
- [ ] 添加消息队列支持
- [ ] 完善单元测试
- [ ] 添加集成测试
- [ ] 性能优化

---

## v1.0.0 (2025-12-10) - 初始版本

### 功能特性

- ✅ 完整的微服务架构
- ✅ GORM + MySQL 数据持久化
- ✅ Kitex RPC 服务
- ✅ Hertz HTTP 网关
- ✅ Etcd 服务注册发现
- ✅ JWT 认证
- ✅ Docker Compose 环境
- ✅ 完整文档

### 核心模块

#### User Service（用户服务）
- 用户注册
- 用户登录
- 用户信息 CRUD
- 分页查询

#### API Gateway（网关服务）
- RESTful API
- JWT 认证
- CORS 支持
- 统一响应格式

#### 基础设施
- MySQL 数据库
- Etcd 服务发现
- Redis 缓存（预留）

### 文档
- README.md - 项目说明
- PROJECT_STRUCTURE.md - 项目结构详解
- TESTING.md - API 测试指南
- GORM详解.md - GORM 教程
- Hertz详解.md - Hertz 教程
- Kitex详解.md - Kitex 教程
- 实战案例-完整微服务架构.md - 实战案例

---

## 版本规范

### 版本号格式
`major.minor.patch`

- **major**: 重大更新，不兼容旧版本
- **minor**: 新功能，向后兼容
- **patch**: Bug 修复，向后兼容

### 更新类型标识

- 🎉 重大更新
- ✨ 新功能
- 🐛 Bug 修复
- 📝 文档更新
- 🔧 工具更新
- ⚠️ 破坏性变更
- 🚀 性能优化
- 🎨 代码优化

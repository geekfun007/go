# Hertz IDL 注解详解

本文档详细说明 Hertz 支持的 IDL 注解和参数校验。

## 目录
- [API 注解](#api-注解)
- [参数校验](#参数校验)
- [路由注解](#路由注解)
- [完整示例](#完整示例)

---

## API 注解

Hertz 支持在 Thrift IDL 中使用注解来定义 HTTP 路由和参数映射。

### 1. 路由方法注解

用于 service 方法上，定义 HTTP 路由。

| 注解 | 说明 | 示例 |
|------|------|------|
| `api.get` | GET 请求 | `api.get="/api/users/:id"` |
| `api.post` | POST 请求 | `api.post="/api/users"` |
| `api.put` | PUT 请求 | `api.put="/api/users/:id"` |
| `api.delete` | DELETE 请求 | `api.delete="/api/users/:id"` |
| `api.patch` | PATCH 请求 | `api.patch="/api/users/:id"` |
| `api.head` | HEAD 请求 | `api.head="/api/users/:id"` |
| `api.options` | OPTIONS 请求 | `api.options="/api/users"` |

**示例：**
```thrift
service UserService {
    // GET 请求
    GetUserResponse GetUser(1: GetUserRequest req) (api.get="/api/v1/users/:id")
    
    // POST 请求
    CreateUserResponse CreateUser(1: CreateUserRequest req) (api.post="/api/v1/users")
    
    // PUT 请求
    UpdateUserResponse UpdateUser(1: UpdateUserRequest req) (api.put="/api/v1/users/:id")
    
    // DELETE 请求
    DeleteUserResponse DeleteUser(1: DeleteUserRequest req) (api.delete="/api/v1/users/:id")
}
```

### 2. 参数来源注解

用于 struct 字段上，定义参数来源。

| 注解 | 说明 | 位置 | 示例 |
|------|------|------|------|
| `api.path` | 路径参数 | URL 路径 | `api.path="id"` |
| `api.query` | 查询参数 | URL query string | `api.query="page"` |
| `api.body` | 请求体参数 | Request Body | `api.body="username"` |
| `api.header` | 请求头参数 | HTTP Header | `api.header="Authorization"` |
| `api.cookie` | Cookie 参数 | HTTP Cookie | `api.cookie="session_id"` |
| `api.form` | 表单参数 | Form Data | `api.form="file"` |
| `api.raw_body` | 原始请求体 | Raw Body | `api.raw_body="true"` |

**示例：**
```thrift
// 路径参数
struct GetUserRequest {
    1: required i64 UserID (api.path="id")
}

// 查询参数
struct ListUsersRequest {
    1: optional i32 Page (api.query="page")
    2: optional i32 PageSize (api.query="page_size")
    3: optional string Keyword (api.query="keyword")
}

// 请求体参数
struct CreateUserRequest {
    1: required string Username (api.body="username")
    2: required string Email (api.body="email")
    3: required string Password (api.body="password")
}

// 请求头参数
struct AuthRequest {
    1: required string Token (api.header="Authorization")
}

// Cookie 参数
struct SessionRequest {
    1: required string SessionID (api.cookie="session_id")
}

// 表单参数（文件上传）
struct UploadRequest {
    1: required string File (api.form="file")
    2: optional string Description (api.form="description")
}
```

### 3. 多参数来源混合

可以在同一个 struct 中混合使用不同的参数来源。

```thrift
struct UpdateUserRequest {
    // 路径参数：用户 ID
    1: required i64 UserID (api.path="id")
    
    // 请求体参数：更新的字段
    2: optional string Username (api.body="username")
    3: optional string Email (api.body="email")
    4: optional i32 Age (api.body="age")
    
    // 查询参数：操作选项
    5: optional bool SendNotification (api.query="send_notification")
}

// 对应路由
service UserService {
    UpdateUserResponse UpdateUser(1: UpdateUserRequest req) (api.put="/api/v1/users/:id")
}

// 使用示例：
// PUT /api/v1/users/123?send_notification=true
// Body: {"username": "newname", "email": "new@example.com"}
```

---

## 参数校验

Hertz 支持使用 `api.vd` 注解进行参数校验。

### 1. 基本校验

| 校验类型 | 语法 | 说明 | 示例 |
|---------|------|------|------|
| 必填 | `len($) > 0` | 字符串非空 | `api.vd="len($) > 0"` |
| 长度范围 | `len($) >= min && len($) <= max` | 字符串长度 | `api.vd="len($) >= 6 && len($) <= 20"` |
| 数值范围 | `$ >= min && $ <= max` | 数字范围 | `api.vd="$ >= 0 && $ <= 150"` |
| 大于 | `$ > value` | 数字大于 | `api.vd="$ > 0"` |
| 小于 | `$ < value` | 数字小于 | `api.vd="$ < 100"` |
| 等于 | `$ == value` | 等于某值 | `api.vd='$ == "active"'` |
| 不等于 | `$ != value` | 不等于某值 | `api.vd='$ != "deleted"'` |

**示例：**
```thrift
struct RegisterRequest {
    // 用户名：3-20 字符
    1: required string Username (
        api.body="username", 
        api.vd="len($) >= 3 && len($) <= 20"
    )
    
    // 邮箱：非空
    2: required string Email (
        api.body="email", 
        api.vd="len($) > 0"
    )
    
    // 密码：6-20 字符
    3: required string Password (
        api.body="password", 
        api.vd="len($) >= 6 && len($) <= 20"
    )
    
    // 年龄：0-150
    4: optional i32 Age (
        api.body="age", 
        api.vd="$ >= 0 && $ <= 150"
    )
    
    // 用户 ID：必须大于 0
    5: required i64 UserID (
        api.path="id", 
        api.vd="$ > 0"
    )
}
```

### 2. 正则表达式校验

```thrift
struct PhoneRequest {
    // 手机号：正则校验
    1: required string Phone (
        api.body="phone",
        api.vd="regexp($, '^1[3-9]\\d{9}$')"
    )
    
    // 邮箱：正则校验
    2: required string Email (
        api.body="email",
        api.vd="regexp($, '^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$')"
    )
}
```

### 3. 枚举值校验

```thrift
struct StatusRequest {
    // 状态：只能是特定值
    1: required string Status (
        api.body="status",
        api.vd='$ == "active" || $ == "inactive" || $ == "pending"'
    )
    
    // 或使用枚举
    2: required i32 Role (
        api.body="role",
        api.vd="$ >= 1 && $ <= 3"  // 1=admin, 2=user, 3=guest
    )
}
```

### 4. 列表校验

```thrift
struct BatchRequest {
    // 列表长度校验
    1: required list<i64> UserIDs (
        api.body="user_ids",
        api.vd="len($) > 0 && len($) <= 100"
    )
    
    // 列表中每个元素校验（需要自定义验证器）
    2: required list<string> Tags (
        api.body="tags",
        api.vd="len($) > 0"
    )
}
```

### 5. 组合校验

使用 `&&` 和 `||` 组合多个条件。

```thrift
struct ComplexRequest {
    // 用户名：3-20字符，且不能包含特殊字符
    1: required string Username (
        api.body="username",
        api.vd="len($) >= 3 && len($) <= 20 && regexp($, '^[a-zA-Z0-9_]+$')"
    )
    
    // 年龄：可选，如果提供则必须在 18-65 之间
    2: optional i32 Age (
        api.body="age",
        api.vd="$ >= 18 && $ <= 65"
    )
    
    // 分页大小：必须是 10, 20, 50, 100 之一
    3: required i32 PageSize (
        api.query="page_size",
        api.vd="$ == 10 || $ == 20 || $ == 50 || $ == 100"
    )
}
```

---

## 路由注解

### 1. 路径参数

使用 `:name` 定义路径参数。

```thrift
service UserService {
    // 单个路径参数
    GetUserResponse GetUser(1: GetUserRequest req) (api.get="/api/v1/users/:id")
    
    // 多个路径参数
    GetPostResponse GetPost(1: GetPostRequest req) (api.get="/api/v1/users/:user_id/posts/:post_id")
}

struct GetUserRequest {
    1: required i64 ID (api.path="id")
}

struct GetPostRequest {
    1: required i64 UserID (api.path="user_id")
    2: required i64 PostID (api.path="post_id")
}
```

### 2. 查询参数

查询参数自动从 URL query string 获取。

```thrift
service UserService {
    ListUsersResponse ListUsers(1: ListUsersRequest req) (api.get="/api/v1/users")
}

struct ListUsersRequest {
    1: optional i32 Page (api.query="page", api.vd="$ > 0")
    2: optional i32 PageSize (api.query="page_size", api.vd="$ > 0 && $ <= 100")
    3: optional string Keyword (api.query="keyword")
    4: optional string SortBy (api.query="sort_by")
    5: optional string Order (api.query="order")
}

// 请求示例：
// GET /api/v1/users?page=1&page_size=10&keyword=test&sort_by=created_at&order=desc
```

### 3. 请求体参数

POST/PUT 请求的请求体参数。

```thrift
service UserService {
    CreateUserResponse CreateUser(1: CreateUserRequest req) (api.post="/api/v1/users")
}

struct CreateUserRequest {
    1: required string Username (api.body="username", api.vd="len($) > 0")
    2: required string Email (api.body="email", api.vd="len($) > 0")
    3: required string Password (api.body="password", api.vd="len($) >= 6")
    4: optional string Phone (api.body="phone")
    5: optional i32 Age (api.body="age")
}

// 请求示例：
// POST /api/v1/users
// Content-Type: application/json
// Body: {
//   "username": "testuser",
//   "email": "test@example.com",
//   "password": "password123",
//   "age": 25
// }
```

### 4. 文件上传

```thrift
service FileService {
    UploadResponse Upload(1: UploadRequest req) (api.post="/api/v1/upload")
}

struct UploadRequest {
    1: required binary File (api.form="file")
    2: optional string Description (api.form="description")
    3: optional string Category (api.form="category")
}

// 请求示例：
// POST /api/v1/upload
// Content-Type: multipart/form-data
// 
// --boundary
// Content-Disposition: form-data; name="file"; filename="test.jpg"
// Content-Type: image/jpeg
// 
// [binary file data]
// --boundary
// Content-Disposition: form-data; name="description"
// 
// Test file upload
// --boundary--
```

---

## 完整示例

### 用户服务 IDL

```thrift
include "common.thrift"

namespace go user

// 用户信息
struct User {
    1: required i64 ID
    2: required string Username
    3: required string Email
    4: optional string Phone
    5: optional i32 Age
    6: optional string Avatar
    7: optional string CreatedAt
    8: optional string UpdatedAt
}

// 注册请求
struct RegisterRequest {
    1: required string Username (
        api.body="username", 
        api.vd="len($) >= 3 && len($) <= 20 && regexp($, '^[a-zA-Z0-9_]+$')"
    )
    2: required string Email (
        api.body="email", 
        api.vd="len($) > 0 && regexp($, '^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$')"
    )
    3: required string Password (
        api.body="password", 
        api.vd="len($) >= 6 && len($) <= 20"
    )
    4: optional string Phone (
        api.body="phone",
        api.vd="len($) == 0 || regexp($, '^1[3-9]\\d{9}$')"
    )
    5: optional i32 Age (
        api.body="age", 
        api.vd="$ >= 0 && $ <= 150"
    )
}

// 注册响应
struct RegisterResponse {
    1: required common.BaseResponse Base
    2: optional User UserInfo
    3: optional string Token
}

// 登录请求
struct LoginRequest {
    1: required string Email (
        api.body="email", 
        api.vd="len($) > 0"
    )
    2: required string Password (
        api.body="password", 
        api.vd="len($) > 0"
    )
}

// 登录响应
struct LoginResponse {
    1: required common.BaseResponse Base
    2: optional User UserInfo
    3: optional string Token
}

// 获取用户请求
struct GetUserRequest {
    1: required i64 UserID (
        api.path="id", 
        api.vd="$ > 0"
    )
}

// 获取用户响应
struct GetUserResponse {
    1: required common.BaseResponse Base
    2: optional User UserInfo
}

// 更新用户请求
struct UpdateUserRequest {
    1: required i64 UserID (
        api.path="id", 
        api.vd="$ > 0"
    )
    2: optional string Username (
        api.body="username", 
        api.vd="len($) >= 3 && len($) <= 20"
    )
    3: optional string Email (
        api.body="email"
    )
    4: optional string Phone (
        api.body="phone"
    )
    5: optional i32 Age (
        api.body="age", 
        api.vd="$ >= 0 && $ <= 150"
    )
    6: optional string Avatar (
        api.body="avatar"
    )
}

// 更新用户响应
struct UpdateUserResponse {
    1: required common.BaseResponse Base
    2: optional User UserInfo
}

// 删除用户请求
struct DeleteUserRequest {
    1: required i64 UserID (
        api.path="id", 
        api.vd="$ > 0"
    )
}

// 删除用户响应
struct DeleteUserResponse {
    1: required common.BaseResponse Base
}

// 列出用户请求
struct ListUsersRequest {
    1: optional i32 Page (
        api.query="page", 
        api.vd="$ > 0"
    )
    2: optional i32 PageSize (
        api.query="page_size", 
        api.vd="$ > 0 && $ <= 100"
    )
    3: optional string Keyword (
        api.query="keyword"
    )
}

// 列出用户响应
struct ListUsersResponse {
    1: required common.BaseResponse Base
    2: optional list<User> Users
    3: optional common.PageInfo PageInfo
}

// 获取当前用户请求
struct GetCurrentUserRequest {
    // 从 JWT token 中获取用户信息
}

// 获取当前用户响应
struct GetCurrentUserResponse {
    1: required common.BaseResponse Base
    2: optional User UserInfo
}

// 服务定义
service UserService {
    // 认证相关
    RegisterResponse Register(1: RegisterRequest req) (
        api.post="/api/v1/auth/register"
    )
    
    LoginResponse Login(1: LoginRequest req) (
        api.post="/api/v1/auth/login"
    )
    
    // 用户管理（公开）
    GetUserResponse GetUser(1: GetUserRequest req) (
        api.get="/api/v1/users/:id"
    )
    
    ListUsersResponse ListUsers(1: ListUsersRequest req) (
        api.get="/api/v1/users"
    )
    
    // 用户管理（需要认证）
    GetCurrentUserResponse GetCurrentUser(1: GetCurrentUserRequest req) (
        api.get="/api/v1/users/me"
    )
    
    UpdateUserResponse UpdateUser(1: UpdateUserRequest req) (
        api.put="/api/v1/users/:id"
    )
    
    DeleteUserResponse DeleteUser(1: DeleteUserRequest req) (
        api.delete="/api/v1/users/:id"
    )
}
```

---

## 使用 hz 生成代码

### 生成命令

```bash
# 生成 Hertz 服务端代码
hz new -module github.com/example/api-gateway \
   -idl user.thrift \
   -handler_dir biz/handler \
   -model_dir biz/model \
   -router_dir biz/router

# 更新已有项目
hz update -idl user.thrift
```

### 生成的代码结构

```
api-gateway/
├── biz/
│   ├── handler/           # 自动生成的 Handler
│   │   └── user/
│   │       └── user_service.go
│   ├── model/             # 自动生成的 Model
│   │   └── user/
│   │       └── user.go
│   └── router/            # 自动生成的 Router
│       ├── register.go
│       └── user/
│           └── user.go
├── router.go              # 路由注册
├── router_gen.go          # 自动生成的路由
└── main.go                # 主入口
```

### 生成的 Handler 示例

hz 会根据 IDL 自动生成 Handler 框架：

```go
// biz/handler/user/user_service.go
package user

import (
    "context"
    
    "github.com/cloudwego/hertz/pkg/app"
    "github.com/example/api-gateway/biz/model/user"
)

// Register 注册用户
// @router /api/v1/auth/register [POST]
func Register(ctx context.Context, c *app.RequestContext) {
    var req user.RegisterRequest
    err := c.BindAndValidate(&req)
    if err != nil {
        c.String(400, err.Error())
        return
    }
    
    // TODO: 实现业务逻辑
    resp := user.RegisterResponse{}
    
    c.JSON(200, resp)
}

// Login 用户登录
// @router /api/v1/auth/login [POST]
func Login(ctx context.Context, c *app.RequestContext) {
    var req user.LoginRequest
    err := c.BindAndValidate(&req)
    if err != nil {
        c.String(400, err.Error())
        return
    }
    
    // TODO: 实现业务逻辑
    resp := user.LoginResponse{}
    
    c.JSON(200, resp)
}
```

### 生成的 Router 示例

```go
// biz/router/user/user.go
package user

import (
    "github.com/cloudwego/hertz/pkg/app/server"
    handler "github.com/example/api-gateway/biz/handler/user"
)

func Register(h *server.Hertz) {
    userGroup := h.Group("/api/v1")
    
    // 认证路由
    authGroup := userGroup.Group("/auth")
    {
        authGroup.POST("/register", handler.Register)
        authGroup.POST("/login", handler.Login)
    }
    
    // 用户路由
    usersGroup := userGroup.Group("/users")
    {
        usersGroup.GET("", handler.ListUsers)
        usersGroup.GET("/:id", handler.GetUser)
        usersGroup.GET("/me", handler.GetCurrentUser)
        usersGroup.PUT("/:id", handler.UpdateUser)
        usersGroup.DELETE("/:id", handler.DeleteUser)
    }
}
```

---

## 总结

### Hertz IDL 注解优势

1. **自动生成**：路由、Handler、Model 自动生成
2. **参数校验**：内置参数验证，减少手动校验代码
3. **类型安全**：编译时检查参数类型
4. **文档友好**：注解即文档，易于维护
5. **统一规范**：团队协作时保持一致性

### 最佳实践

1. **合理使用校验**：在 IDL 中定义基础校验，复杂业务校验在代码中实现
2. **清晰的路由**：使用 RESTful 风格，路径语义化
3. **字段命名**：使用 PascalCase（IDL）→ snake_case（JSON）
4. **版本管理**：路由中包含版本号 `/api/v1/`
5. **文档注释**：在 IDL 中添加注释说明

### 常用验证规则

| 场景 | 验证规则 | 示例 |
|------|---------|------|
| 必填字符串 | `len($) > 0` | 用户名、邮箱 |
| 字符串长度 | `len($) >= min && len($) <= max` | 密码长度 |
| 数字范围 | `$ >= min && $ <= max` | 年龄、分页大小 |
| 正整数 | `$ > 0` | ID、数量 |
| 手机号 | `regexp($, '^1[3-9]\\d{9}$')` | 中国手机号 |
| 邮箱 | `regexp($, '^[\\w.+-]+@[\\w.-]+\\.[a-zA-Z]{2,}$')` | 邮箱格式 |
| 用户名 | `regexp($, '^[a-zA-Z0-9_]+$')` | 字母数字下划线 |

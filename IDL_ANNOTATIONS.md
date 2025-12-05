# IDL 注解说明文档

本文档详细说明了项目中使用的 Thrift IDL 注解，包括字段注解和方法注解。

## 📋 概述

本项目使用 **Hertz IDL 注解规范**，通过注解实现：
- 自动路由生成
- 参数绑定
- 参数校验
- API 文档生成

参考文档：[Hertz 注解说明](https://www.cloudwego.io/zh/docs/hertz/tutorials/toolkit/annotation/#%E6%94%AF%E6%8C%81%E7%9A%84-api-%E6%B3%A8%E8%A7%A3)

## 🏷️ 字段注解（Field Annotations）

### 1. go.tag - Go 结构体标签

用于为生成的 Go 结构体字段添加标签。

**语法**:
```thrift
field_name (go.tag = "tag_content")
```

**支持的标签**:

| 标签 | 说明 | 示例 |
|------|------|------|
| `json` | JSON 序列化字段名 | `json:"username"` |
| `form` | Form 表单字段名 | `form:"username"` |
| `query` | URL 查询参数名 | `query:"page"` |
| `path` | URL 路径参数名 | `path:"id"` |
| `header` | HTTP 头字段名 | `header:"X-Token"` |
| `cookie` | Cookie 字段名 | `cookie:"session_id"` |
| `vd` | 参数校验规则 | `vd:"len($) > 0"` |

**示例**:
```thrift
struct CreateUserRequest {
    1: required string Username (
        go.tag = "json:\"username\" form:\"username\" vd:\"len($) > 0 && len($) <= 64\""
    )
}
```

### 2. api.body - 请求体参数

标记字段从 HTTP 请求体中获取。

**语法**:
```thrift
field_name (api.body = "field_name")
```

**示例**:
```thrift
struct CreateUserRequest {
    1: required string Username (
        go.tag = "json:\"username\"",
        api.body = "username"
    )
    2: required string Email (
        go.tag = "json:\"email\"",
        api.body = "email"
    )
}
```

**生成的代码**:
```go
type CreateUserRequest struct {
    Username string `json:"username"` // 从请求体获取
    Email    string `json:"email"`    // 从请求体获取
}
```

### 3. api.query - 查询参数

标记字段从 URL 查询参数中获取。

**语法**:
```thrift
field_name (api.query = "param_name")
```

**示例**:
```thrift
struct ListUsersRequest {
    1: required i32 Page (
        go.tag = "json:\"page\" query:\"page\"",
        api.query = "page"
    )
    2: required i32 PageSize (
        go.tag = "json:\"page_size\" query:\"page_size\"",
        api.query = "page_size"
    )
}
```

**对应的 URL**: `/api/v1/users?page=1&page_size=10`

### 4. api.path - 路径参数

标记字段从 URL 路径参数中获取。

**语法**:
```thrift
field_name (api.path = "param_name")
```

**示例**:
```thrift
struct GetUserRequest {
    1: required i64 UserID (
        go.tag = "json:\"user_id\" path:\"id\"",
        api.path = "id"
    )
}
```

**对应的 URL**: `/api/v1/users/:id` → `/api/v1/users/123`

### 5. api.header - 请求头参数

标记字段从 HTTP 请求头中获取。

**语法**:
```thrift
field_name (api.header = "Header-Name")
```

**示例**:
```thrift
struct AuthRequest {
    1: required string Token (
        go.tag = "header:\"Authorization\"",
        api.header = "Authorization"
    )
}
```

### 6. api.cookie - Cookie 参数

标记字段从 Cookie 中获取。

**语法**:
```thrift
field_name (api.cookie = "cookie_name")
```

**示例**:
```thrift
struct SessionRequest {
    1: required string SessionID (
        go.tag = "cookie:\"session_id\"",
        api.cookie = "session_id"
    )
}
```

### 7. api.form - 表单参数

标记字段从表单数据中获取。

**语法**:
```thrift
field_name (api.form = "field_name")
```

**示例**:
```thrift
struct LoginRequest {
    1: required string Username (
        go.tag = "form:\"username\"",
        api.form = "username"
    )
    2: required string Password (
        go.tag = "form:\"password\"",
        api.form = "password"
    )
}
```

## 🔧 方法注解（Method Annotations）

### 1. HTTP 方法注解

指定 HTTP 方法和路径。

#### api.get - GET 请求

**语法**:
```thrift
Method(Request) (api.get = "/path")
```

**示例**:
```thrift
// 获取用户列表
ListUsersResponse ListUsers(1: ListUsersRequest req) (
    api.get = "/api/v1/users",
    api.serializer = "json"
)

// 获取单个用户
GetUserResponse GetUser(1: GetUserRequest req) (
    api.get = "/api/v1/users/:id",
    api.serializer = "json"
)
```

#### api.post - POST 请求

**语法**:
```thrift
Method(Request) (api.post = "/path")
```

**示例**:
```thrift
// 创建用户
CreateUserResponse CreateUser(1: CreateUserRequest req) (
    api.post = "/api/v1/users",
    api.serializer = "json"
)
```

#### api.put - PUT 请求

**语法**:
```thrift
Method(Request) (api.put = "/path")
```

**示例**:
```thrift
// 更新用户
UpdateUserResponse UpdateUser(1: UpdateUserRequest req) (
    api.put = "/api/v1/users/:id",
    api.serializer = "json"
)
```

#### api.delete - DELETE 请求

**语法**:
```thrift
Method(Request) (api.delete = "/path")
```

**示例**:
```thrift
// 删除用户
DeleteUserResponse DeleteUser(1: DeleteUserRequest req) (
    api.delete = "/api/v1/users/:id",
    api.serializer = "json"
)
```

#### api.patch - PATCH 请求

**语法**:
```thrift
Method(Request) (api.patch = "/path")
```

**示例**:
```thrift
// 部分更新用户
PatchUserResponse PatchUser(1: PatchUserRequest req) (
    api.patch = "/api/v1/users/:id",
    api.serializer = "json"
)
```

### 2. api.serializer - 序列化方式

指定请求和响应的序列化格式。

**支持的值**:
- `json` - JSON 格式（默认）
- `form` - 表单格式
- `thrift` - Thrift 二进制格式
- `protobuf` - Protobuf 格式

**示例**:
```thrift
CreateUserResponse CreateUser(1: CreateUserRequest req) (
    api.post = "/api/v1/users",
    api.serializer = "json"  // 使用 JSON 格式
)
```

### 3. api.baseurl - 基础 URL

为服务指定基础 URL（可选）。

**示例**:
```thrift
service UserService {
    // 方法定义...
} (
    api.baseurl = "http://localhost:8080"
)
```

## 📝 完整示例

### 用户服务 IDL

```thrift
namespace go user

// 用户信息
struct User {
    1: i64 ID (go.tag = "json:\"id\"")
    2: string Username (go.tag = "json:\"username\"")
    3: string Email (go.tag = "json:\"email\"")
    4: string Phone (go.tag = "json:\"phone\"")
    5: i64 CreatedAt (go.tag = "json:\"created_at\"")
    6: i64 UpdatedAt (go.tag = "json:\"updated_at\"")
}

// 创建用户请求
struct CreateUserRequest {
    1: required string Username (
        go.tag = "json:\"username\" form:\"username\" vd:\"len($) > 0 && len($) <= 64\"",
        api.body = "username"
    )
    2: required string Email (
        go.tag = "json:\"email\" form:\"email\" vd:\"len($) > 0 && len($) <= 128\"",
        api.body = "email"
    )
    3: optional string Phone (
        go.tag = "json:\"phone,omitempty\" form:\"phone\" vd:\"len($) == 0 || len($) == 11\"",
        api.body = "phone"
    )
}

// 创建用户响应
struct CreateUserResponse {
    1: i32 Code (go.tag = "json:\"code\"")
    2: string Message (go.tag = "json:\"message\"")
    3: optional User User (go.tag = "json:\"user,omitempty\"")
}

// 获取用户请求
struct GetUserRequest {
    1: required i64 UserID (
        go.tag = "json:\"user_id\" path:\"id\" vd:\"$ > 0\"",
        api.path = "id"
    )
}

// 获取用户响应
struct GetUserResponse {
    1: i32 Code (go.tag = "json:\"code\"")
    2: string Message (go.tag = "json:\"message\"")
    3: optional User User (go.tag = "json:\"user,omitempty\"")
}

// 用户列表请求
struct ListUsersRequest {
    1: required i32 Page (
        go.tag = "json:\"page\" form:\"page\" query:\"page\" vd:\"$ > 0\"",
        api.query = "page"
    )
    2: required i32 PageSize (
        go.tag = "json:\"page_size\" form:\"page_size\" query:\"page_size\" vd:\"$ > 0 && $ <= 100\"",
        api.query = "page_size"
    )
}

// 用户列表响应
struct ListUsersResponse {
    1: i32 Code (go.tag = "json:\"code\"")
    2: string Message (go.tag = "json:\"message\"")
    3: list<User> Users (go.tag = "json:\"users,omitempty\"")
    4: i64 Total (go.tag = "json:\"total\"")
}

// 用户服务
service UserService {
    // 创建用户
    CreateUserResponse CreateUser(1: CreateUserRequest req) (
        api.post = "/api/v1/users",
        api.serializer = "json"
    )
    
    // 获取用户详情
    GetUserResponse GetUser(1: GetUserRequest req) (
        api.get = "/api/v1/users/:id",
        api.serializer = "json"
    )
    
    // 获取用户列表
    ListUsersResponse ListUsers(1: ListUsersRequest req) (
        api.get = "/api/v1/users",
        api.serializer = "json"
    )
}
```

## 🎯 生成的路由

Hz 会根据注解自动生成路由：

```go
// biz/router/user/user.go
func Register(r *server.Hertz) {
    root := r.Group("/", rootMw()...)
    {
        _api := root.Group("/api", _apiMw()...)
        {
            _v1 := _api.Group("/v1", _v1Mw()...)
            _v1.GET("/users", append(_listusersMw(), user.ListUsers)...)
            _v1.POST("/users", append(_createuserMw(), user.CreateUser)...)
            _users := _v1.Group("/users", _usersMw()...)
            _users.DELETE("/:id", append(_deleteuserMw(), user.DeleteUser)...)
            _users.GET("/:id", append(_getuserMw(), user.GetUser)...)
            _users.PUT("/:id", append(_updateuserMw(), user.UpdateUser)...)
        }
    }
}
```

## 🔄 参数绑定优先级

当多个注解同时存在时，参数绑定的优先级：

1. `api.path` - 路径参数（最高优先级）
2. `api.query` - 查询参数
3. `api.header` - 请求头
4. `api.cookie` - Cookie
5. `api.body` - 请求体
6. `api.form` - 表单数据

**示例**:
```thrift
struct MixedRequest {
    1: required i64 ID (api.path = "id")              // 从路径获取
    2: required string Token (api.header = "Token")   // 从请求头获取
    3: optional i32 Page (api.query = "page")         // 从查询参数获取
    4: optional string Data (api.body = "data")       // 从请求体获取
}
```

对应的请求：
```bash
curl -X POST "http://localhost:8080/api/v1/resource/123?page=1" \
  -H "Token: abc123" \
  -H "Content-Type: application/json" \
  -d '{"data":"example"}'
```

## 🧪 测试注解

### 测试路径参数

```bash
# GET /api/v1/users/:id
curl http://localhost:8080/api/v1/users/1
```

### 测试查询参数

```bash
# GET /api/v1/users?page=1&page_size=10
curl "http://localhost:8080/api/v1/users?page=1&page_size=10"
```

### 测试请求体参数

```bash
# POST /api/v1/users
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "alice",
    "email": "alice@example.com",
    "phone": "13800138000"
  }'
```

### 测试 PUT 请求

```bash
# PUT /api/v1/users/:id
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "username": "alice_updated"
  }'
```

### 测试 DELETE 请求

```bash
# DELETE /api/v1/users/:id
curl -X DELETE http://localhost:8080/api/v1/users/1
```

## 📚 最佳实践

### 1. 字段注解

✅ **推荐**: 同时使用 go.tag 和 api 注解
```thrift
1: required string Username (
    go.tag = "json:\"username\" vd:\"len($) > 0\"",
    api.body = "username"
)
```

❌ **不推荐**: 只使用一种注解
```thrift
1: required string Username (go.tag = "json:\"username\"")
```

### 2. 方法注解

✅ **推荐**: 明确指定 HTTP 方法和序列化方式
```thrift
CreateUserResponse CreateUser(1: CreateUserRequest req) (
    api.post = "/api/v1/users",
    api.serializer = "json"
)
```

❌ **不推荐**: 缺少注解
```thrift
CreateUserResponse CreateUser(1: CreateUserRequest req)
```

### 3. RESTful 风格

✅ **推荐**: 遵循 RESTful 规范
```thrift
api.get = "/api/v1/users"         // 列表
api.post = "/api/v1/users"        // 创建
api.get = "/api/v1/users/:id"     // 详情
api.put = "/api/v1/users/:id"     // 更新
api.delete = "/api/v1/users/:id"  // 删除
```

### 4. 参数校验

✅ **推荐**: 添加完整的校验规则
```thrift
1: required string Email (
    go.tag = "json:\"email\" vd:\"len($) > 0 && len($) <= 128\"",
    api.body = "email"
)
```

## 🔄 更新流程

当修改 IDL 注解后：

```bash
# 1. 重新生成 Kitex 代码
cd /workspace
kitex -module github.com/example/hertz-kitex-demo -service user_service ./idl/user.thrift

# 2. 更新 Hertz 代码
cd hertz_service
hz update -idl ../idl/user.thrift

# 3. 检查生成的路由
cat biz/router/user/user.go

# 4. 编译测试
go build .
```

## 📖 参考资料

- [Hertz 注解说明](https://www.cloudwego.io/zh/docs/hertz/tutorials/toolkit/annotation/)
- [Hz CLI 工具](https://www.cloudwego.io/zh/docs/hertz/tutorials/toolkit/usage/)
- [Kitex Thrift IDL](https://www.cloudwego.io/zh/docs/kitex/tutorials/code-gen/code_generation/)
- [Thrift IDL 规范](https://thrift.apache.org/docs/idl)

---

**最后更新**: 2024-12-05

**版本**: v2.0

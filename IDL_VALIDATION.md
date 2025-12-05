# IDL 字段定义与校验说明

## 📋 概述

本项目使用 **首字母大写** 的字段定义，符合 Go 语言命名规范，并集成了 **Hertz 的参数绑定与校验** 功能。

## 🔤 字段命名规范

### IDL 定义（首字母大写）

```thrift
struct User {
    1: i64 ID (go.tag = "json:\"id\"")
    2: string Username (go.tag = "json:\"username\"")
    3: string Email (go.tag = "json:\"email\"")
    4: string Phone (go.tag = "json:\"phone\"")
    5: i64 CreatedAt (go.tag = "json:\"created_at\"")
    6: i64 UpdatedAt (go.tag = "json:\"updated_at\"")
}
```

### 生成的 Go 代码

```go
type User struct {
    ID        int64  `json:"id"`
    Username  string `json:"username"`
    Email     string `json:"email"`
    Phone     string `json:"phone"`
    CreatedAt int64  `json:"created_at"`
    UpdatedAt int64  `json:"updated_at"`
}
```

## ✅ 参数校验规则

### 1. CreateUserRequest - 创建用户

```thrift
struct CreateUserRequest {
    1: required string Username (go.tag = "json:\"username\" form:\"username\" query:\"username\" vd:\"len($) > 0 && len($) <= 64\"")
    2: required string Email (go.tag = "json:\"email\" form:\"email\" query:\"email\" vd:\"len($) > 0 && len($) <= 128\"")
    3: optional string Phone (go.tag = "json:\"phone,omitempty\" form:\"phone\" query:\"phone\" vd:\"len($) == 0 || len($) == 11\"")
}
```

**校验规则**:
- ✅ `Username`: 必填，长度 1-64 个字符
- ✅ `Email`: 必填，长度 1-128 个字符
- ✅ `Phone`: 可选，如果提供必须是 11 位

**示例请求**:
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "email": "john@example.com",
    "phone": "13800138000"
  }'
```

### 2. GetUserRequest - 获取用户

```thrift
struct GetUserRequest {
    1: required i64 UserID (go.tag = "json:\"user_id\" path:\"id\" vd:\"$ > 0\"")
}
```

**校验规则**:
- ✅ `UserID`: 必填，必须大于 0

**示例请求**:
```bash
curl http://localhost:8080/api/v1/users/1
```

### 3. UpdateUserRequest - 更新用户

```thrift
struct UpdateUserRequest {
    1: required i64 UserID (go.tag = "json:\"user_id\" path:\"id\" vd:\"$ > 0\"")
    2: optional string Username (go.tag = "json:\"username,omitempty\" form:\"username\" vd:\"len($) == 0 || len($) <= 64\"")
    3: optional string Email (go.tag = "json:\"email,omitempty\" form:\"email\" vd:\"len($) == 0 || len($) <= 128\"")
    4: optional string Phone (go.tag = "json:\"phone,omitempty\" form:\"phone\" vd:\"len($) == 0 || len($) == 11\"")
}
```

**校验规则**:
- ✅ `UserID`: 必填，必须大于 0
- ✅ `Username`: 可选，如果提供则长度不超过 64
- ✅ `Email`: 可选，如果提供则长度不超过 128
- ✅ `Phone`: 可选，如果提供则必须是 11 位

**示例请求**:
```bash
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_updated"
  }'
```

### 4. DeleteUserRequest - 删除用户

```thrift
struct DeleteUserRequest {
    1: required i64 UserID (go.tag = "json:\"user_id\" path:\"id\" vd:\"$ > 0\"")
}
```

**校验规则**:
- ✅ `UserID`: 必填，必须大于 0

**示例请求**:
```bash
curl -X DELETE http://localhost:8080/api/v1/users/1
```

### 5. ListUsersRequest - 获取用户列表

```thrift
struct ListUsersRequest {
    1: required i32 Page (go.tag = "json:\"page\" form:\"page\" query:\"page\" vd:\"$ > 0\"")
    2: required i32 PageSize (go.tag = "json:\"page_size\" form:\"page_size\" query:\"page_size\" vd:\"$ > 0 && $ <= 100\"")
}
```

**校验规则**:
- ✅ `Page`: 必填，必须大于 0
- ✅ `PageSize`: 必填，范围 1-100

**示例请求**:
```bash
curl "http://localhost:8080/api/v1/users?page=1&page_size=10"
```

## 🏷️ Go Tag 说明

### 支持的标签类型

```thrift
go.tag = "json:\"field_name\" form:\"field_name\" query:\"field_name\" path:\"field_name\" vd:\"validation_rule\""
```

| 标签 | 说明 | 用途 |
|------|------|------|
| `json` | JSON 序列化字段名 | HTTP 请求/响应体 |
| `form` | Form 表单字段名 | Form 数据绑定 |
| `query` | URL 查询参数名 | Query 参数绑定 |
| `path` | URL 路径参数名 | Path 参数绑定 |
| `vd` | 校验规则 | 参数校验 |

### 校验规则语法（vd tag）

| 规则 | 说明 | 示例 |
|------|------|------|
| `len($)` | 字符串/数组长度 | `len($) > 0` |
| `$ > 0` | 数值比较 | `$ > 0 && $ <= 100` |
| `len($) == 0 \|\| expr` | 或条件 | `len($) == 0 \|\| len($) == 11` |
| `&&` | 与条件 | `len($) > 0 && len($) <= 64` |

**注意**: 复杂的正则表达式校验建议在代码层面实现，避免 IDL 解析问题。

## 🔄 字段类型映射

### Thrift → Go

| Thrift 类型 | Go 类型 | 说明 |
|-------------|---------|------|
| `i64` | `int64` | 64位整数 |
| `i32` | `int32` | 32位整数 |
| `string` | `string` | 字符串 |
| `bool` | `bool` | 布尔值 |
| `list<T>` | `[]T` | 切片 |
| `map<K,V>` | `map[K]V` | 映射 |
| `optional T` | `*T` | 指针（可选） |
| `required T` | `T` | 非指针（必填） |

## 📝 最佳实践

### 1. 字段命名

✅ **推荐**: 使用首字母大写（Go 风格）
```thrift
struct User {
    1: i64 UserID
    2: string Username
}
```

❌ **不推荐**: 使用小写或下划线
```thrift
struct User {
    1: i64 user_id
    2: string user_name
}
```

### 2. 必填与可选

- 使用 `required` 标记必填字段
- 使用 `optional` 标记可选字段
- 可选字段在 Go 中会生成指针类型

```thrift
struct CreateUserRequest {
    1: required string Username  // string
    2: optional string Phone     // *string
}
```

### 3. 校验规则

- 保持校验规则简单明了
- 复杂校验在代码层实现
- 使用 `||` 处理可选字段

```thrift
// ✅ 简单校验
vd:"len($) > 0 && len($) <= 64"

// ✅ 可选字段校验
vd:"len($) == 0 || len($) == 11"

// ❌ 避免复杂正则（IDL 解析问题）
vd:"regexp('^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$')"
```

### 4. JSON 标签

- 使用下划线分隔的小写字段名（符合 JSON 规范）
- 可选字段添加 `omitempty`

```thrift
1: i64 UserID (go.tag = "json:\"user_id\"")
2: optional string Phone (go.tag = "json:\"phone,omitempty\"")
```

## 🧪 测试校验

### 测试正常请求

```bash
# 正常请求 - 应该成功
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "alice",
    "email": "alice@example.com",
    "phone": "13800138000"
  }'
```

### 测试校验失败

```bash
# 用户名为空 - 应该失败
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "",
    "email": "alice@example.com"
  }'

# 用户名过长 - 应该失败
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "a_very_long_username_that_exceeds_the_maximum_allowed_length_of_64_characters",
    "email": "alice@example.com"
  }'

# 手机号格式错误 - 应该失败
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "alice",
    "email": "alice@example.com",
    "phone": "123"
  }'
```

## 📚 参考资料

- [Hertz Binding and Validate](https://www.cloudwego.io/docs/hertz/tutorials/basic-feature/binding-and-validate/)
- [Thrift IDL](https://thrift.apache.org/docs/idl)
- [Go Tag Expr](https://github.com/bytedance/go-tagexpr)
- [Kitex Thrift](https://www.cloudwego.io/docs/kitex/tutorials/code-gen/code_generation/)

## 🔄 更新流程

当修改 IDL 后：

```bash
# 1. 重新生成 Kitex 代码
make gen-kitex

# 2. 重新生成 Hertz 代码
make update-hertz

# 3. 更新业务逻辑代码
# - 检查字段名是否变化
# - 检查 optional 字段（指针类型）
# - 更新校验逻辑

# 4. 编译测试
make build
```

---

**注意**: 
- Hz 生成的代码会自动包含校验标签
- Hertz 框架会自动执行参数校验
- 校验失败会返回 400 错误

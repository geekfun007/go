include "common.thrift"

namespace go user

// ==================== 数据结构 ====================

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

// ==================== 注册相关 ====================

// 注册请求
struct RegisterRequest {
    1: required string Username (api.body="username", api.vd="len($) > 0 && len($) <= 20")
    2: required string Email (api.body="email", api.vd="len($) > 0")
    3: required string Password (api.body="password", api.vd="len($) >= 6 && len($) <= 20")
    4: optional string Phone (api.body="phone")
    5: optional i32 Age (api.body="age", api.vd="$ >= 0 && $ <= 150")
}

// 注册响应
struct RegisterResponse {
    1: required common.BaseResponse Base
    2: optional User UserInfo
    3: optional string Token
}

// ==================== 登录相关 ====================

// 登录请求
struct LoginRequest {
    1: required string Email (api.body="email", api.vd="len($) > 0")
    2: required string Password (api.body="password", api.vd="len($) > 0")
}

// 登录响应
struct LoginResponse {
    1: required common.BaseResponse Base
    2: optional User UserInfo
    3: optional string Token
}

// ==================== 用户管理 ====================

// 获取用户请求
struct GetUserRequest {
    1: required i64 UserID (api.path="id", api.vd="$ > 0")
}

// 获取用户响应
struct GetUserResponse {
    1: required common.BaseResponse Base
    2: optional User UserInfo
}

// 更新用户请求
struct UpdateUserRequest {
    1: required i64 UserID (api.path="id", api.vd="$ > 0")
    2: optional string Username (api.body="username", api.vd="len($) > 0 && len($) <= 20")
    3: optional string Email (api.body="email", api.vd="len($) > 0")
    4: optional string Phone (api.body="phone")
    5: optional i32 Age (api.body="age", api.vd="$ >= 0 && $ <= 150")
    6: optional string Avatar (api.body="avatar")
}

// 更新用户响应
struct UpdateUserResponse {
    1: required common.BaseResponse Base
    2: optional User UserInfo
}

// 删除用户请求
struct DeleteUserRequest {
    1: required i64 UserID (api.path="id", api.vd="$ > 0")
}

// 删除用户响应
struct DeleteUserResponse {
    1: required common.BaseResponse Base
}

// 列出用户请求
struct ListUsersRequest {
    1: optional i32 Page (api.query="page", api.vd="$ > 0")
    2: optional i32 PageSize (api.query="page_size", api.vd="$ > 0 && $ <= 100")
    3: optional string Keyword (api.query="keyword")
}

// 列出用户响应
struct ListUsersResponse {
    1: required common.BaseResponse Base
    2: optional list<User> Users
    3: optional common.PageInfo PageInfo
}

// 获取当前用户请求（空请求，从 token 获取）
struct GetCurrentUserRequest {
    // 从 JWT token 中获取用户 ID，无需参数
}

// 获取当前用户响应
struct GetCurrentUserResponse {
    1: required common.BaseResponse Base
    2: optional User UserInfo
}

// ==================== 服务定义 ====================

service UserService {
    // 认证相关（不需要 token）
    RegisterResponse Register(1: RegisterRequest req) (api.post="/api/v1/auth/register")
    LoginResponse Login(1: LoginRequest req) (api.post="/api/v1/auth/login")
    
    // 用户管理（公开接口）
    GetUserResponse GetUser(1: GetUserRequest req) (api.get="/api/v1/users/:id")
    ListUsersResponse ListUsers(1: ListUsersRequest req) (api.get="/api/v1/users")
    
    // 用户管理（需要认证）
    GetCurrentUserResponse GetCurrentUser(1: GetCurrentUserRequest req) (api.get="/api/v1/users/me")
    UpdateUserResponse UpdateUser(1: UpdateUserRequest req) (api.put="/api/v1/users/:id")
    DeleteUserResponse DeleteUser(1: DeleteUserRequest req) (api.delete="/api/v1/users/:id")
}

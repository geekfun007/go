include "common.thrift"

namespace go user

// ==================== 数据结构 ====================

// 用户信息
struct User {
    1: required i64 id
    2: required string username
    3: required string email
    4: optional string phone
    5: optional i32 age
    6: optional string avatar
    7: optional string created_at
    8: optional string updated_at
}

// ==================== 注册相关 ====================

// 注册请求
struct RegisterRequest {
    1: required string username
    2: required string email
    3: required string password
    4: optional string phone
    5: optional i32 age
}

// 注册响应
struct RegisterResponse {
    1: required common.BaseResponse base
    2: optional User user
    3: optional string token
}

// ==================== 登录相关 ====================

// 登录请求
struct LoginRequest {
    1: required string email
    2: required string password
}

// 登录响应
struct LoginResponse {
    1: required common.BaseResponse base
    2: optional User user
    3: optional string token
}

// ==================== 用户管理 ====================

// 获取用户请求
struct GetUserRequest {
    1: required i64 user_id
}

// 获取用户响应
struct GetUserResponse {
    1: required common.BaseResponse base
    2: optional User user
}

// 更新用户请求
struct UpdateUserRequest {
    1: required i64 user_id
    2: optional string username
    3: optional string email
    4: optional string phone
    5: optional i32 age
    6: optional string avatar
}

// 更新用户响应
struct UpdateUserResponse {
    1: required common.BaseResponse base
    2: optional User user
}

// 删除用户请求
struct DeleteUserRequest {
    1: required i64 user_id
}

// 删除用户响应
struct DeleteUserResponse {
    1: required common.BaseResponse base
}

// 列出用户请求
struct ListUsersRequest {
    1: optional common.PageRequest page
    2: optional string keyword
}

// 列出用户响应
struct ListUsersResponse {
    1: required common.BaseResponse base
    2: optional list<User> users
    3: optional common.PageInfo page_info
}

// ==================== 服务定义 ====================

service UserService {
    // 认证相关
    RegisterResponse Register(1: RegisterRequest req)
    LoginResponse Login(1: LoginRequest req)
    
    // 用户管理
    GetUserResponse GetUser(1: GetUserRequest req)
    UpdateUserResponse UpdateUser(1: UpdateUserRequest req)
    DeleteUserResponse DeleteUser(1: DeleteUserRequest req)
    ListUsersResponse ListUsers(1: ListUsersRequest req)
}

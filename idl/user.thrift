namespace go user

// 用户信息
struct User {
    1: i64 id
    2: string username
    3: string email
    4: string phone
    5: i64 created_at
    6: i64 updated_at
}

// 创建用户请求
struct CreateUserRequest {
    1: string username
    2: string email
    3: string phone
}

// 创建用户响应
struct CreateUserResponse {
    1: i32 code
    2: string message
    3: User user
}

// 获取用户请求
struct GetUserRequest {
    1: i64 user_id
}

// 获取用户响应
struct GetUserResponse {
    1: i32 code
    2: string message
    3: User user
}

// 更新用户请求
struct UpdateUserRequest {
    1: i64 user_id
    2: optional string username
    3: optional string email
    4: optional string phone
}

// 更新用户响应
struct UpdateUserResponse {
    1: i32 code
    2: string message
}

// 删除用户请求
struct DeleteUserRequest {
    1: i64 user_id
}

// 删除用户响应
struct DeleteUserResponse {
    1: i32 code
    2: string message
}

// 用户列表请求
struct ListUsersRequest {
    1: i32 page
    2: i32 page_size
}

// 用户列表响应
struct ListUsersResponse {
    1: i32 code
    2: string message
    3: list<User> users
    4: i64 total
}

// 用户服务
service UserService {
    CreateUserResponse CreateUser(1: CreateUserRequest req)
    GetUserResponse GetUser(1: GetUserRequest req)
    UpdateUserResponse UpdateUser(1: UpdateUserRequest req)
    DeleteUserResponse DeleteUser(1: DeleteUserRequest req)
    ListUsersResponse ListUsers(1: ListUsersRequest req)
}

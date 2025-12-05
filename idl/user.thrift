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
    1: required string Username (go.tag = "json:\"username\" form:\"username\" vd:\"len($) > 0 && len($) <= 64\"", api.body = "username")
    2: required string Email (go.tag = "json:\"email\" form:\"email\" vd:\"len($) > 0 && len($) <= 128\"", api.body = "email")
    3: optional string Phone (go.tag = "json:\"phone,omitempty\" form:\"phone\" vd:\"len($) == 0 || len($) == 11\"", api.body = "phone")
}

// 创建用户响应
struct CreateUserResponse {
    1: i32 Code (go.tag = "json:\"code\"")
    2: string Message (go.tag = "json:\"message\"")
    3: optional User User (go.tag = "json:\"user,omitempty\"")
}

// 获取用户请求
struct GetUserRequest {
    1: required i64 UserID (go.tag = "json:\"user_id\" path:\"id\" vd:\"$ > 0\"", api.path = "id")
}

// 获取用户响应
struct GetUserResponse {
    1: i32 Code (go.tag = "json:\"code\"")
    2: string Message (go.tag = "json:\"message\"")
    3: optional User User (go.tag = "json:\"user,omitempty\"")
}

// 更新用户请求
struct UpdateUserRequest {
    1: required i64 UserID (go.tag = "json:\"user_id\" path:\"id\" vd:\"$ > 0\"", api.path = "id")
    2: optional string Username (go.tag = "json:\"username,omitempty\" form:\"username\" vd:\"len($) == 0 || len($) <= 64\"", api.body = "username")
    3: optional string Email (go.tag = "json:\"email,omitempty\" form:\"email\" vd:\"len($) == 0 || len($) <= 128\"", api.body = "email")
    4: optional string Phone (go.tag = "json:\"phone,omitempty\" form:\"phone\" vd:\"len($) == 0 || len($) == 11\"", api.body = "phone")
}

// 更新用户响应
struct UpdateUserResponse {
    1: i32 Code (go.tag = "json:\"code\"")
    2: string Message (go.tag = "json:\"message\"")
}

// 删除用户请求
struct DeleteUserRequest {
    1: required i64 UserID (go.tag = "json:\"user_id\" path:\"id\" vd:\"$ > 0\"", api.path = "id")
}

// 删除用户响应
struct DeleteUserResponse {
    1: i32 Code (go.tag = "json:\"code\"")
    2: string Message (go.tag = "json:\"message\"")
}

// 用户列表请求
struct ListUsersRequest {
    1: required i32 Page (go.tag = "json:\"page\" form:\"page\" query:\"page\" vd:\"$ > 0\"", api.query = "page")
    2: required i32 PageSize (go.tag = "json:\"page_size\" form:\"page_size\" query:\"page_size\" vd:\"$ > 0 && $ <= 100\"", api.query = "page_size")
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
    CreateUserResponse CreateUser(1: CreateUserRequest req) (api.post = "/api/v1/users", api.serializer = "json")
    
    // 获取用户详情
    GetUserResponse GetUser(1: GetUserRequest req) (api.get = "/api/v1/users/:id", api.serializer = "json")
    
    // 更新用户信息
    UpdateUserResponse UpdateUser(1: UpdateUserRequest req) (api.put = "/api/v1/users/:id", api.serializer = "json")
    
    // 删除用户
    DeleteUserResponse DeleteUser(1: DeleteUserRequest req) (api.delete = "/api/v1/users/:id", api.serializer = "json")
    
    // 获取用户列表
    ListUsersResponse ListUsers(1: ListUsersRequest req) (api.get = "/api/v1/users", api.serializer = "json")
}

package user

// User HTTP 响应模型
type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

// CreateUserRequest HTTP 请求
type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Phone    string `json:"phone"`
}

// UpdateUserRequest HTTP 请求
type UpdateUserRequest struct {
	Username *string `json:"username,omitempty"`
	Email    *string `json:"email,omitempty"`
	Phone    *string `json:"phone,omitempty"`
}

// BaseResponse 基础响应
type BaseResponse struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

// CreateUserResponse 创建用户响应
type CreateUserResponse struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
	User    *User  `json:"user,omitempty"`
}

// GetUserResponse 获取用户响应
type GetUserResponse struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
	User    *User  `json:"user,omitempty"`
}

// ListUsersResponse 用户列表响应
type ListUsersResponse struct {
	Code    int32   `json:"code"`
	Message string  `json:"message"`
	Users   []*User `json:"users,omitempty"`
	Total   int64   `json:"total"`
}

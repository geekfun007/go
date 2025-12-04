package model

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

// UserResponse HTTP 响应
type UserResponse struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

// BaseResponse 基础响应
type BaseResponse struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

// CreateUserResponse 创建用户响应
type CreateUserResponse struct {
	BaseResponse
	Data *UserResponse `json:"data,omitempty"`
}

// GetUserResponse 获取用户响应
type GetUserResponse struct {
	BaseResponse
	Data *UserResponse `json:"data,omitempty"`
}

// ListUsersResponse 用户列表响应
type ListUsersResponse struct {
	BaseResponse
	Data struct {
		Users []*UserResponse `json:"users"`
		Total int64           `json:"total"`
	} `json:"data,omitempty"`
}

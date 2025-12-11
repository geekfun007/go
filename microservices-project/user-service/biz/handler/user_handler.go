package handler

import (
	"context"
	"errors"
	"time"

	"github.com/example/microservices-project/user-service/biz/model"
	"github.com/example/microservices-project/user-service/biz/repository"
	"github.com/example/microservices-project/user-service/biz/service"
	"github.com/example/microservices-project/user-service/kitex_gen/common"
	"github.com/example/microservices-project/user-service/kitex_gen/user"
)

type UserServiceImpl struct {
	userService *service.UserService
}

func NewUserService(userService *service.UserService) *UserServiceImpl {
	return &UserServiceImpl{userService: userService}
}

// Register 用户注册
func (s *UserServiceImpl) Register(ctx context.Context, req *user.RegisterRequest) (*user.RegisterResponse, error) {
	// 参数验证
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return &user.RegisterResponse{
			Base: &common.BaseResponse{
				Code:    400,
				Message: "username, email and password are required",
			},
		}, nil
	}

	// 调用 service 层
	age := 0
	if req.Age != nil {
		age = int(*req.Age)
	}

	phone := ""
	if req.Phone != nil {
		phone = *req.Phone
	}

	newUser, token, err := s.userService.Register(ctx, req.Username, req.Email, req.Password, phone, age)
	if err != nil {
		if errors.Is(err, repository.ErrUserExists) {
			return &user.RegisterResponse{
				Base: &common.BaseResponse{
					Code:    40001,
					Message: "email already exists",
				},
			}, nil
		}
		if err.Error() == "username already exists" {
			return &user.RegisterResponse{
				Base: &common.BaseResponse{
					Code:    40002,
					Message: "username already exists",
				},
			}, nil
		}
		return &user.RegisterResponse{
			Base: &common.BaseResponse{
				Code:    500,
				Message: "internal error: " + err.Error(),
			},
		}, nil
	}

	return &user.RegisterResponse{
		Base: &common.BaseResponse{
			Code:    0,
			Message: "success",
		},
		User:  convertToThriftUser(newUser),
		Token: &token,
	}, nil
}

// Login 用户登录
func (s *UserServiceImpl) Login(ctx context.Context, req *user.LoginRequest) (*user.LoginResponse, error) {
	if req.Email == "" || req.Password == "" {
		return &user.LoginResponse{
			Base: &common.BaseResponse{
				Code:    400,
				Message: "email and password are required",
			},
		}, nil
	}

	loginUser, token, err := s.userService.Login(ctx, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return &user.LoginResponse{
				Base: &common.BaseResponse{
					Code:    40401,
					Message: "user not found",
				},
			}, nil
		}
		if err.Error() == "invalid password" {
			return &user.LoginResponse{
				Base: &common.BaseResponse{
					Code:    40101,
					Message: "invalid password",
				},
			}, nil
		}
		return &user.LoginResponse{
			Base: &common.BaseResponse{
				Code:    500,
				Message: "internal error: " + err.Error(),
			},
		}, nil
	}

	return &user.LoginResponse{
		Base: &common.BaseResponse{
			Code:    0,
			Message: "success",
		},
		User:  convertToThriftUser(loginUser),
		Token: &token,
	}, nil
}

// GetUser 获取用户信息
func (s *UserServiceImpl) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.GetUserResponse, error) {
	u, err := s.userService.GetUser(ctx, uint(req.UserId))
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return &user.GetUserResponse{
				Base: &common.BaseResponse{
					Code:    40402,
					Message: "user not found",
				},
			}, nil
		}
		return &user.GetUserResponse{
			Base: &common.BaseResponse{
				Code:    500,
				Message: "internal error: " + err.Error(),
			},
		}, nil
	}

	return &user.GetUserResponse{
		Base: &common.BaseResponse{
			Code:    0,
			Message: "success",
		},
		User: convertToThriftUser(u),
	}, nil
}

// UpdateUser 更新用户信息
func (s *UserServiceImpl) UpdateUser(ctx context.Context, req *user.UpdateUserRequest) (*user.UpdateUserResponse, error) {
	var age *int
	if req.Age != nil {
		a := int(*req.Age)
		age = &a
	}

	u, err := s.userService.UpdateUser(ctx, uint(req.UserId), req.Username, req.Email, req.Phone, req.Avatar, age)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return &user.UpdateUserResponse{
				Base: &common.BaseResponse{
					Code:    40403,
					Message: "user not found",
				},
			}, nil
		}
		if err.Error() == "username already exists" {
			return &user.UpdateUserResponse{
				Base: &common.BaseResponse{
					Code:    40003,
					Message: "username already exists",
				},
			}, nil
		}
		if err.Error() == "email already exists" {
			return &user.UpdateUserResponse{
				Base: &common.BaseResponse{
					Code:    40004,
					Message: "email already exists",
				},
			}, nil
		}
		return &user.UpdateUserResponse{
			Base: &common.BaseResponse{
				Code:    500,
				Message: "internal error: " + err.Error(),
			},
		}, nil
	}

	return &user.UpdateUserResponse{
		Base: &common.BaseResponse{
			Code:    0,
			Message: "success",
		},
		User: convertToThriftUser(u),
	}, nil
}

// DeleteUser 删除用户
func (s *UserServiceImpl) DeleteUser(ctx context.Context, req *user.DeleteUserRequest) (*user.DeleteUserResponse, error) {
	err := s.userService.DeleteUser(ctx, uint(req.UserId))
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return &user.DeleteUserResponse{
				Base: &common.BaseResponse{
					Code:    40404,
					Message: "user not found",
				},
			}, nil
		}
		return &user.DeleteUserResponse{
			Base: &common.BaseResponse{
				Code:    500,
				Message: "internal error: " + err.Error(),
			},
		}, nil
	}

	return &user.DeleteUserResponse{
		Base: &common.BaseResponse{
			Code:    0,
			Message: "success",
		},
	}, nil
}

// ListUsers 列出用户
func (s *UserServiceImpl) ListUsers(ctx context.Context, req *user.ListUsersRequest) (*user.ListUsersResponse, error) {
	page := 1
	pageSize := 10
	keyword := ""

	if req.Page != nil {
		page = int(req.Page.Page)
		pageSize = int(req.Page.PageSize)
	}

	if req.Keyword != nil {
		keyword = *req.Keyword
	}

	users, total, err := s.userService.ListUsers(ctx, page, pageSize, keyword)
	if err != nil {
		return &user.ListUsersResponse{
			Base: &common.BaseResponse{
				Code:    500,
				Message: "internal error: " + err.Error(),
			},
		}, nil
	}

	thriftUsers := make([]*user.User, 0, len(users))
	for _, u := range users {
		thriftUsers = append(thriftUsers, convertToThriftUser(u))
	}

	return &user.ListUsersResponse{
		Base: &common.BaseResponse{
			Code:    0,
			Message: "success",
		},
		Users: thriftUsers,
		PageInfo: &common.PageInfo{
			Page:     int32(page),
			PageSize: int32(pageSize),
			Total:    total,
		},
	}, nil
}

// convertToThriftUser 将 model.User 转换为 thrift User
func convertToThriftUser(u *model.User) *user.User {
	createdAt := u.CreatedAt.Format(time.RFC3339)
	updatedAt := u.UpdatedAt.Format(time.RFC3339)
	age := int32(u.Age)

	thriftUser := &user.User{
		Id:        int64(u.ID),
		Username:  u.Username,
		Email:     u.Email,
		Age:       &age,
		CreatedAt: &createdAt,
		UpdatedAt: &updatedAt,
	}

	if u.Phone != "" {
		thriftUser.Phone = &u.Phone
	}

	if u.Avatar != "" {
		thriftUser.Avatar = &u.Avatar
	}

	return thriftUser
}

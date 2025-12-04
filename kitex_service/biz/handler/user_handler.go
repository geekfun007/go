package handler

import (
	"context"
	"errors"
	"log"

	"github.com/example/hertz-kitex-demo/kitex_gen/user"
	"github.com/example/hertz-kitex-demo/kitex_service/dal"
	"github.com/example/hertz-kitex-demo/kitex_service/model"
	"gorm.io/gorm"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct {
	userDAL *dal.UserDAL
}

// NewUserServiceImpl 创建服务实例
func NewUserServiceImpl() *UserServiceImpl {
	return &UserServiceImpl{
		userDAL: dal.NewUserDAL(),
	}
}

// CreateUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) CreateUser(ctx context.Context, req *user.CreateUserRequest) (resp *user.CreateUserResponse, err error) {
	resp = &user.CreateUserResponse{}

	// 参数校验
	if req.Username == "" || req.Email == "" {
		resp.Code = 400
		resp.Message = "username and email are required"
		return resp, nil
	}

	// 检查用户名是否已存在
	existUser, _ := s.userDAL.GetUserByUsername(ctx, req.Username)
	if existUser != nil {
		resp.Code = 400
		resp.Message = "username already exists"
		return resp, nil
	}

	// 检查邮箱是否已存在
	existUser, _ = s.userDAL.GetUserByEmail(ctx, req.Email)
	if existUser != nil {
		resp.Code = 400
		resp.Message = "email already exists"
		return resp, nil
	}

	// 创建用户
	dbUser := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Phone:    req.Phone,
	}

	if err = s.userDAL.CreateUser(ctx, dbUser); err != nil {
		log.Printf("CreateUser error: %v", err)
		resp.Code = 500
		resp.Message = "failed to create user"
		return resp, nil
	}

	// 转换为返回对象
	resp.Code = 0
	resp.Message = "success"
	resp.User = &user.User{
		Id:        dbUser.ID,
		Username:  dbUser.Username,
		Email:     dbUser.Email,
		Phone:     dbUser.Phone,
		CreatedAt: dbUser.CreatedAt.Unix(),
		UpdatedAt: dbUser.UpdatedAt.Unix(),
	}

	return resp, nil
}

// GetUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetUser(ctx context.Context, req *user.GetUserRequest) (resp *user.GetUserResponse, err error) {
	resp = &user.GetUserResponse{}

	if req.UserId <= 0 {
		resp.Code = 400
		resp.Message = "invalid user_id"
		return resp, nil
	}

	dbUser, err := s.userDAL.GetUserByID(ctx, req.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			resp.Code = 404
			resp.Message = "user not found"
			return resp, nil
		}
		log.Printf("GetUser error: %v", err)
		resp.Code = 500
		resp.Message = "failed to get user"
		return resp, nil
	}

	resp.Code = 0
	resp.Message = "success"
	resp.User = &user.User{
		Id:        dbUser.ID,
		Username:  dbUser.Username,
		Email:     dbUser.Email,
		Phone:     dbUser.Phone,
		CreatedAt: dbUser.CreatedAt.Unix(),
		UpdatedAt: dbUser.UpdatedAt.Unix(),
	}

	return resp, nil
}

// UpdateUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) UpdateUser(ctx context.Context, req *user.UpdateUserRequest) (resp *user.UpdateUserResponse, err error) {
	resp = &user.UpdateUserResponse{}

	if req.UserId <= 0 {
		resp.Code = 400
		resp.Message = "invalid user_id"
		return resp, nil
	}

	// 检查用户是否存在
	_, err = s.userDAL.GetUserByID(ctx, req.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			resp.Code = 404
			resp.Message = "user not found"
			return resp, nil
		}
		log.Printf("GetUser error: %v", err)
		resp.Code = 500
		resp.Message = "failed to get user"
		return resp, nil
	}

	// 构建更新数据
	updates := make(map[string]interface{})
	if req.Username != nil {
		updates["username"] = *req.Username
	}
	if req.Email != nil {
		updates["email"] = *req.Email
	}
	if req.Phone != nil {
		updates["phone"] = *req.Phone
	}

	if len(updates) == 0 {
		resp.Code = 400
		resp.Message = "no fields to update"
		return resp, nil
	}

	// 更新用户
	if err = s.userDAL.UpdateUser(ctx, req.UserId, updates); err != nil {
		log.Printf("UpdateUser error: %v", err)
		resp.Code = 500
		resp.Message = "failed to update user"
		return resp, nil
	}

	resp.Code = 0
	resp.Message = "success"
	return resp, nil
}

// DeleteUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) DeleteUser(ctx context.Context, req *user.DeleteUserRequest) (resp *user.DeleteUserResponse, err error) {
	resp = &user.DeleteUserResponse{}

	if req.UserId <= 0 {
		resp.Code = 400
		resp.Message = "invalid user_id"
		return resp, nil
	}

	// 检查用户是否存在
	_, err = s.userDAL.GetUserByID(ctx, req.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			resp.Code = 404
			resp.Message = "user not found"
			return resp, nil
		}
		log.Printf("GetUser error: %v", err)
		resp.Code = 500
		resp.Message = "failed to get user"
		return resp, nil
	}

	// 删除用户
	if err = s.userDAL.DeleteUser(ctx, req.UserId); err != nil {
		log.Printf("DeleteUser error: %v", err)
		resp.Code = 500
		resp.Message = "failed to delete user"
		return resp, nil
	}

	resp.Code = 0
	resp.Message = "success"
	return resp, nil
}

// ListUsers implements the UserServiceImpl interface.
func (s *UserServiceImpl) ListUsers(ctx context.Context, req *user.ListUsersRequest) (resp *user.ListUsersResponse, err error) {
	resp = &user.ListUsersResponse{}

	// 参数校验
	page := req.Page
	pageSize := req.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	// 查询用户列表
	dbUsers, total, err := s.userDAL.ListUsers(ctx, int(page), int(pageSize))
	if err != nil {
		log.Printf("ListUsers error: %v", err)
		resp.Code = 500
		resp.Message = "failed to list users"
		return resp, nil
	}

	// 转换为返回对象
	users := make([]*user.User, 0, len(dbUsers))
	for _, dbUser := range dbUsers {
		users = append(users, &user.User{
			Id:        dbUser.ID,
			Username:  dbUser.Username,
			Email:     dbUser.Email,
			Phone:     dbUser.Phone,
			CreatedAt: dbUser.CreatedAt.Unix(),
			UpdatedAt: dbUser.UpdatedAt.Unix(),
		})
	}

	resp.Code = 0
	resp.Message = "success"
	resp.Users = users
	resp.Total = total

	return resp, nil
}

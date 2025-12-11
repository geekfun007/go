package handler

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/example/microservices-project/api-gateway/biz/client"
	"github.com/example/microservices-project/api-gateway/biz/middleware"
	"github.com/example/microservices-project/api-gateway/kitex_gen/common"
	"github.com/example/microservices-project/api-gateway/kitex_gen/user"
	"github.com/example/microservices-project/api-gateway/pkg/response"
)

// ==================== 请求/响应结构 ====================

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=20"`
	Phone    string `json:"phone"`
	Age      int32  `json:"age"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UpdateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Age      int32  `json:"age"`
	Avatar   string `json:"avatar"`
}

// ==================== 处理器函数 ====================

// Register 用户注册
func Register(ctx context.Context, c *app.RequestContext) {
	var req RegisterRequest
	if err := c.BindAndValidate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 调用用户服务
	userClient := client.GetUserClient()
	resp, err := userClient.Register(ctx, &user.RegisterRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Phone:    &req.Phone,
		Age:      &req.Age,
	})

	if err != nil {
		response.InternalServerError(c, "failed to call user service: "+err.Error())
		return
	}

	// 处理响应
	if resp.Base.Code != 0 {
		response.Error(c, int(resp.Base.Code), resp.Base.Message)
		return
	}

	response.Success(c, map[string]interface{}{
		"user":  resp.User,
		"token": resp.Token,
	})
}

// Login 用户登录
func Login(ctx context.Context, c *app.RequestContext) {
	var req LoginRequest
	if err := c.BindAndValidate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userClient := client.GetUserClient()
	resp, err := userClient.Login(ctx, &user.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		response.InternalServerError(c, "failed to call user service: "+err.Error())
		return
	}

	if resp.Base.Code != 0 {
		response.Error(c, int(resp.Base.Code), resp.Base.Message)
		return
	}

	response.Success(c, map[string]interface{}{
		"user":  resp.User,
		"token": resp.Token,
	})
}

// GetUser 获取用户信息
func GetUser(ctx context.Context, c *app.RequestContext) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	userClient := client.GetUserClient()
	resp, err := userClient.GetUser(ctx, &user.GetUserRequest{
		UserId: id,
	})

	if err != nil {
		response.InternalServerError(c, "failed to call user service: "+err.Error())
		return
	}

	if resp.Base.Code != 0 {
		response.Error(c, int(resp.Base.Code), resp.Base.Message)
		return
	}

	response.Success(c, resp.User)
}

// GetCurrentUser 获取当前登录用户信息
func GetCurrentUser(ctx context.Context, c *app.RequestContext) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "unauthorized")
		return
	}

	userClient := client.GetUserClient()
	resp, err := userClient.GetUser(ctx, &user.GetUserRequest{
		UserId: int64(userID),
	})

	if err != nil {
		response.InternalServerError(c, "failed to call user service: "+err.Error())
		return
	}

	if resp.Base.Code != 0 {
		response.Error(c, int(resp.Base.Code), resp.Base.Message)
		return
	}

	response.Success(c, resp.User)
}

// UpdateUser 更新用户信息
func UpdateUser(ctx context.Context, c *app.RequestContext) {
	// 获取要更新的用户 ID
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	// 验证权限：只能更新自己的信息
	currentUserID, ok := middleware.GetUserID(c)
	if !ok || uint(id) != currentUserID {
		response.Forbidden(c, "you can only update your own information")
		return
	}

	var req UpdateUserRequest
	if err := c.BindAndValidate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userClient := client.GetUserClient()
	resp, err := userClient.UpdateUser(ctx, &user.UpdateUserRequest{
		UserId:   id,
		Username: &req.Username,
		Email:    &req.Email,
		Phone:    &req.Phone,
		Age:      &req.Age,
		Avatar:   &req.Avatar,
	})

	if err != nil {
		response.InternalServerError(c, "failed to call user service: "+err.Error())
		return
	}

	if resp.Base.Code != 0 {
		response.Error(c, int(resp.Base.Code), resp.Base.Message)
		return
	}

	response.Success(c, resp.User)
}

// DeleteUser 删除用户
func DeleteUser(ctx context.Context, c *app.RequestContext) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	// 验证权限：只能删除自己的账号
	currentUserID, ok := middleware.GetUserID(c)
	if !ok || uint(id) != currentUserID {
		response.Forbidden(c, "you can only delete your own account")
		return
	}

	userClient := client.GetUserClient()
	resp, err := userClient.DeleteUser(ctx, &user.DeleteUserRequest{
		UserId: id,
	})

	if err != nil {
		response.InternalServerError(c, "failed to call user service: "+err.Error())
		return
	}

	if resp.Base.Code != 0 {
		response.Error(c, int(resp.Base.Code), resp.Base.Message)
		return
	}

	response.Success(c, nil)
}

// ListUsers 列出用户
func ListUsers(ctx context.Context, c *app.RequestContext) {
	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	keyword := c.Query("keyword")

	userClient := client.GetUserClient()
	resp, err := userClient.ListUsers(ctx, &user.ListUsersRequest{
		Page: &common.PageRequest{
			Page:     int32(page),
			PageSize: int32(pageSize),
		},
		Keyword: &keyword,
	})

	if err != nil {
		response.InternalServerError(c, "failed to call user service: "+err.Error())
		return
	}

	if resp.Base.Code != 0 {
		response.Error(c, int(resp.Base.Code), resp.Base.Message)
		return
	}

	response.Success(c, map[string]interface{}{
		"users": resp.Users,
		"page": map[string]interface{}{
			"page":      resp.PageInfo.Page,
			"page_size": resp.PageInfo.PageSize,
			"total":     resp.PageInfo.Total,
		},
	})
}

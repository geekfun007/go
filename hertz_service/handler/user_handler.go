package handler

import (
	"context"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/example/hertz-kitex-demo/hertz_service/client"
	"github.com/example/hertz-kitex-demo/hertz_service/model"
	"github.com/example/hertz-kitex-demo/kitex_service/kitex_gen/user"
)

type UserHandler struct{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// CreateUser 创建用户
func (h *UserHandler) CreateUser(ctx context.Context, c *app.RequestContext) {
	var req model.CreateUserRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"code":    400,
			"message": "invalid request: " + err.Error(),
		})
		return
	}

	// 调用 RPC 服务
	rpcReq := &user.CreateUserRequest{
		Username: req.Username,
		Email:    req.Email,
		Phone:    req.Phone,
	}

	rpcResp, err := client.GetUserClient().CreateUser(ctx, rpcReq)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{
			"code":    500,
			"message": "rpc call failed: " + err.Error(),
		})
		return
	}

	// 构造响应
	resp := model.CreateUserResponse{
		BaseResponse: model.BaseResponse{
			Code:    rpcResp.Code,
			Message: rpcResp.Message,
		},
	}

	if rpcResp.User != nil {
		resp.Data = &model.UserResponse{
			ID:        rpcResp.User.Id,
			Username:  rpcResp.User.Username,
			Email:     rpcResp.User.Email,
			Phone:     rpcResp.User.Phone,
			CreatedAt: rpcResp.User.CreatedAt,
			UpdatedAt: rpcResp.User.UpdatedAt,
		}
	}

	if resp.Code == 0 {
		c.JSON(consts.StatusOK, resp)
	} else {
		c.JSON(consts.StatusBadRequest, resp)
	}
}

// GetUser 获取用户
func (h *UserHandler) GetUser(ctx context.Context, c *app.RequestContext) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"code":    400,
			"message": "invalid user_id",
		})
		return
	}

	// 调用 RPC 服务
	rpcReq := &user.GetUserRequest{
		UserId: userID,
	}

	rpcResp, err := client.GetUserClient().GetUser(ctx, rpcReq)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{
			"code":    500,
			"message": "rpc call failed: " + err.Error(),
		})
		return
	}

	// 构造响应
	resp := model.GetUserResponse{
		BaseResponse: model.BaseResponse{
			Code:    rpcResp.Code,
			Message: rpcResp.Message,
		},
	}

	if rpcResp.User != nil {
		resp.Data = &model.UserResponse{
			ID:        rpcResp.User.Id,
			Username:  rpcResp.User.Username,
			Email:     rpcResp.User.Email,
			Phone:     rpcResp.User.Phone,
			CreatedAt: rpcResp.User.CreatedAt,
			UpdatedAt: rpcResp.User.UpdatedAt,
		}
	}

	if resp.Code == 0 {
		c.JSON(consts.StatusOK, resp)
	} else if resp.Code == 404 {
		c.JSON(consts.StatusNotFound, resp)
	} else {
		c.JSON(consts.StatusInternalServerError, resp)
	}
}

// UpdateUser 更新用户
func (h *UserHandler) UpdateUser(ctx context.Context, c *app.RequestContext) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"code":    400,
			"message": "invalid user_id",
		})
		return
	}

	var req model.UpdateUserRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"code":    400,
			"message": "invalid request: " + err.Error(),
		})
		return
	}

	// 调用 RPC 服务
	rpcReq := &user.UpdateUserRequest{
		UserId:   userID,
		Username: req.Username,
		Email:    req.Email,
		Phone:    req.Phone,
	}

	rpcResp, err := client.GetUserClient().UpdateUser(ctx, rpcReq)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{
			"code":    500,
			"message": "rpc call failed: " + err.Error(),
		})
		return
	}

	// 构造响应
	resp := model.BaseResponse{
		Code:    rpcResp.Code,
		Message: rpcResp.Message,
	}

	if resp.Code == 0 {
		c.JSON(consts.StatusOK, resp)
	} else if resp.Code == 404 {
		c.JSON(consts.StatusNotFound, resp)
	} else {
		c.JSON(consts.StatusBadRequest, resp)
	}
}

// DeleteUser 删除用户
func (h *UserHandler) DeleteUser(ctx context.Context, c *app.RequestContext) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"code":    400,
			"message": "invalid user_id",
		})
		return
	}

	// 调用 RPC 服务
	rpcReq := &user.DeleteUserRequest{
		UserId: userID,
	}

	rpcResp, err := client.GetUserClient().DeleteUser(ctx, rpcReq)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{
			"code":    500,
			"message": "rpc call failed: " + err.Error(),
		})
		return
	}

	// 构造响应
	resp := model.BaseResponse{
		Code:    rpcResp.Code,
		Message: rpcResp.Message,
	}

	if resp.Code == 0 {
		c.JSON(consts.StatusOK, resp)
	} else if resp.Code == 404 {
		c.JSON(consts.StatusNotFound, resp)
	} else {
		c.JSON(consts.StatusInternalServerError, resp)
	}
}

// ListUsers 获取用户列表
func (h *UserHandler) ListUsers(ctx context.Context, c *app.RequestContext) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// 调用 RPC 服务
	rpcReq := &user.ListUsersRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
	}

	rpcResp, err := client.GetUserClient().ListUsers(ctx, rpcReq)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{
			"code":    500,
			"message": "rpc call failed: " + err.Error(),
		})
		return
	}

	// 构造响应
	resp := model.ListUsersResponse{
		BaseResponse: model.BaseResponse{
			Code:    rpcResp.Code,
			Message: rpcResp.Message,
		},
	}

	users := make([]*model.UserResponse, 0, len(rpcResp.Users))
	for _, u := range rpcResp.Users {
		users = append(users, &model.UserResponse{
			ID:        u.Id,
			Username:  u.Username,
			Email:     u.Email,
			Phone:     u.Phone,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		})
	}

	resp.Data.Users = users
	resp.Data.Total = rpcResp.Total

	c.JSON(consts.StatusOK, resp)
}

package dal

import (
	"context"
	"fmt"

	"github.com/example/hertz-kitex-demo/kitex_service/model"
)

// UserDAL 数据访问层
type UserDAL struct{}

func NewUserDAL() *UserDAL {
	return &UserDAL{}
}

// CreateUser 创建用户
func (d *UserDAL) CreateUser(ctx context.Context, user *model.User) error {
	if err := DB.WithContext(ctx).Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// GetUserByID 根据ID获取用户
func (d *UserDAL) GetUserByID(ctx context.Context, userID int64) (*model.User, error) {
	var user model.User
	if err := DB.WithContext(ctx).Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

// UpdateUser 更新用户
func (d *UserDAL) UpdateUser(ctx context.Context, userID int64, updates map[string]interface{}) error {
	if err := DB.WithContext(ctx).Model(&model.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// DeleteUser 删除用户
func (d *UserDAL) DeleteUser(ctx context.Context, userID int64) error {
	if err := DB.WithContext(ctx).Where("id = ?", userID).Delete(&model.User{}).Error; err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// ListUsers 获取用户列表
func (d *UserDAL) ListUsers(ctx context.Context, page, pageSize int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	// 计算总数
	if err := DB.WithContext(ctx).Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := DB.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	return users, total, nil
}

// GetUserByUsername 根据用户名获取用户
func (d *UserDAL) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	if err := DB.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}
	return &user, nil
}

// GetUserByEmail 根据邮箱获取用户
func (d *UserDAL) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := DB.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return &user, nil
}

package service

import (
	"context"
	"errors"
	"time"

	"github.com/example/microservices-project/user-service/biz/model"
	"github.com/example/microservices-project/user-service/biz/repository"
	"github.com/example/microservices-project/user-service/conf"
	"github.com/example/microservices-project/user-service/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// Register 用户注册
func (s *UserService) Register(ctx context.Context, username, email, password, phone string, age int) (*model.User, string, error) {
	// 检查邮箱是否已存在
	exists, err := s.repo.ExistsByEmail(ctx, email)
	if err != nil {
		return nil, "", err
	}
	if exists {
		return nil, "", repository.ErrUserExists
	}

	// 检查用户名是否已存在
	exists, err = s.repo.ExistsByUsername(ctx, username)
	if err != nil {
		return nil, "", err
	}
	if exists {
		return nil, "", errors.New("username already exists")
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	// 创建用户
	user := &model.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
		Phone:    phone,
		Age:      age,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, "", err
	}

	// 生成 JWT token
	token, err := utils.GenerateToken(user.ID, conf.GetConfig().JWT.Secret, conf.GetConfig().JWT.ExpireHours)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// Login 用户登录
func (s *UserService) Login(ctx context.Context, email, password string) (*model.User, string, error) {
	// 查找用户
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, "", err
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", errors.New("invalid password")
	}

	// 生成 JWT token
	token, err := utils.GenerateToken(user.ID, conf.GetConfig().JWT.Secret, conf.GetConfig().JWT.ExpireHours)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// GetUser 获取用户信息
func (s *UserService) GetUser(ctx context.Context, userID uint) (*model.User, error) {
	return s.repo.GetByID(ctx, userID)
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(ctx context.Context, userID uint, username, email, phone, avatar *string, age *int) (*model.User, error) {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 更新字段
	if username != nil && *username != "" {
		// 检查用户名是否被占用
		if *username != user.Username {
			exists, err := s.repo.ExistsByUsername(ctx, *username)
			if err != nil {
				return nil, err
			}
			if exists {
				return nil, errors.New("username already exists")
			}
		}
		user.Username = *username
	}

	if email != nil && *email != "" {
		// 检查邮箱是否被占用
		if *email != user.Email {
			exists, err := s.repo.ExistsByEmail(ctx, *email)
			if err != nil {
				return nil, err
			}
			if exists {
				return nil, errors.New("email already exists")
			}
		}
		user.Email = *email
	}

	if phone != nil {
		user.Phone = *phone
	}

	if avatar != nil {
		user.Avatar = *avatar
	}

	if age != nil {
		user.Age = *age
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(ctx context.Context, userID uint) error {
	return s.repo.Delete(ctx, userID)
}

// ListUsers 列出用户
func (s *UserService) ListUsers(ctx context.Context, page, pageSize int, keyword string) ([]*model.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	return s.repo.List(ctx, page, pageSize, keyword)
}

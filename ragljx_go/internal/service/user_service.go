package service

import (
	"context"
	"ragljx/internal/model"
	"ragljx/internal/pkg/errors"
	"ragljx/internal/repository"

	"gorm.io/gorm"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		userRepo: repository.NewUserRepository(db),
	}
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
	RealName string `json:"real_name"`
	Phone    string `json:"phone"`
	IsAdmin  bool   `json:"is_admin"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Email    string `json:"email" binding:"omitempty,email"`
	RealName string `json:"real_name"`
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar"`
	Status   string `json:"status"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// Create 创建用户
func (s *UserService) Create(ctx context.Context, req *CreateUserRequest) (*model.User, error) {
	// 检查用户名是否存在
	existingUser, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err == nil && existingUser != nil {
		return nil, errors.ErrUserExists
	}

	// 检查邮箱是否存在
	existingUser, err = s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New(400, "email already exists")
	}

	// 创建用户
	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		RealName: req.RealName,
		Phone:    req.Phone,
		IsAdmin:  req.IsAdmin,
		Status:   "active",
	}

	// 设置密码
	if err := user.SetPassword(req.Password); err != nil {
		return nil, errors.Wrap(500, "failed to hash password", err)
	}

	// 保存用户
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, errors.Wrap(500, "failed to create user", err)
	}

	return user, nil
}

// GetByID 根据 ID 获取用户
func (s *UserService) GetByID(ctx context.Context, id int) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrUserNotFound
		}
		return nil, errors.Wrap(500, "failed to get user", err)
	}
	return user, nil
}

// List 获取用户列表
func (s *UserService) List(ctx context.Context, page, pageSize int, keyword string) ([]*model.User, int64, error) {
	offset := (page - 1) * pageSize
	users, total, err := s.userRepo.List(ctx, offset, pageSize, keyword)
	if err != nil {
		return nil, 0, errors.Wrap(500, "failed to list users", err)
	}
	return users, total, nil
}

// Update 更新用户
func (s *UserService) Update(ctx context.Context, id int, req *UpdateUserRequest) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrUserNotFound
		}
		return nil, errors.Wrap(500, "failed to get user", err)
	}

	// 更新字段
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.RealName != "" {
		user.RealName = req.RealName
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.Status != "" {
		user.Status = req.Status
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, errors.Wrap(500, "failed to update user", err)
	}

	return user, nil
}

// Delete 删除用户
func (s *UserService) Delete(ctx context.Context, id int) error {
	if err := s.userRepo.Delete(ctx, id); err != nil {
		return errors.Wrap(500, "failed to delete user", err)
	}
	return nil
}

// ChangePassword 修改密码
func (s *UserService) ChangePassword(ctx context.Context, userID int, req *ChangePasswordRequest) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.ErrUserNotFound
		}
		return errors.Wrap(500, "failed to get user", err)
	}

	// 验证旧密码
	if !user.CheckPassword(req.OldPassword) {
		return errors.New(400, "old password is incorrect")
	}

	// 设置新密码
	if err := user.SetPassword(req.NewPassword); err != nil {
		return errors.Wrap(500, "failed to hash password", err)
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return errors.Wrap(500, "failed to update password", err)
	}

	return nil
}

// AssignRoles 分配角色
func (s *UserService) AssignRoles(ctx context.Context, userID int, roleIDs []int) error {
	// 检查用户是否存在
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.ErrUserNotFound
		}
		return errors.Wrap(500, "failed to get user", err)
	}

	if err := s.userRepo.AssignRoles(ctx, userID, roleIDs); err != nil {
		return errors.Wrap(500, "failed to assign roles", err)
	}

	return nil
}


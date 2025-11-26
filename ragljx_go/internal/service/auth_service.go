package service

import (
	"context"
	"ragljx/internal/model"
	"ragljx/internal/pkg/errors"
	"ragljx/internal/pkg/jwt"
	"ragljx/internal/repository"
	"time"

	"gorm.io/gorm"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{
		userRepo: repository.NewUserRepository(db),
	}
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	User         *model.User `json:"user"`
}

// Login 用户登录
func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// 查找用户
	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrUserNotFound
		}
		return nil, errors.Wrap(500, "failed to get user", err)
	}

	// 验证密码
	if !user.CheckPassword(req.Password) {
		return nil, errors.ErrInvalidPassword
	}

	// 检查用户状态
	if user.Status != "active" {
		return nil, errors.New(403, "user is not active")
	}

	// 生成 token
	accessToken, err := jwt.GenerateToken(user.ID, user.Username, user.IsAdmin)
	if err != nil {
		return nil, errors.Wrap(500, "failed to generate access token", err)
	}

	refreshToken, err := jwt.GenerateRefreshToken(user.ID, user.Username, user.IsAdmin)
	if err != nil {
		return nil, errors.Wrap(500, "failed to generate refresh token", err)
	}

	// 更新最后登录时间
	now := time.Now()
	user.LastLoginAt = &now
	if err := s.userRepo.Update(ctx, user); err != nil {
		// 不影响登录流程，只记录错误
		// 实质上应该打印一条日志
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

// RefreshTokenRequest 刷新 Token 请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshToken 刷新 Token
func (s *AuthService) RefreshToken(ctx context.Context, req *RefreshTokenRequest) (*LoginResponse, error) {
	// 解析 refresh token
	claims, err := jwt.ParseToken(req.RefreshToken)
	if err != nil {
		return nil, errors.ErrInvalidToken
	}

	// 查找用户
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrUserNotFound
		}
		return nil, errors.Wrap(500, "failed to get user", err)
	}

	// 检查用户状态
	if user.Status != "active" {
		return nil, errors.New(403, "user is not active")
	}

	// 生成新的 token
	accessToken, err := jwt.GenerateToken(user.ID, user.Username, user.IsAdmin)
	if err != nil {
		return nil, errors.Wrap(500, "failed to generate access token", err)
	}

	refreshToken, err := jwt.GenerateRefreshToken(user.ID, user.Username, user.IsAdmin)
	if err != nil {
		return nil, errors.Wrap(500, "failed to generate refresh token", err)
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

// GetCurrentUser 获取当前用户信息
func (s *AuthService) GetCurrentUser(ctx context.Context, userID int) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrUserNotFound
		}
		return nil, errors.Wrap(500, "failed to get user", err)
	}
	return user, nil
}


package model

import (
	"ragljx/internal/pkg/utils"
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID          int            `gorm:"primaryKey;autoIncrement" json:"id"`
	Username    string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"username"`
	Password    string         `gorm:"type:varchar(255);not null" json:"-"`
	Email       string         `gorm:"type:varchar(128);uniqueIndex" json:"email"`
	RealName    string         `gorm:"type:varchar(64)" json:"real_name"`
	Phone       string         `gorm:"type:varchar(32)" json:"phone"`
	Avatar      string         `gorm:"type:varchar(255)" json:"avatar"`
	Status      string         `gorm:"type:varchar(32);default:'active'" json:"status"`
	IsAdmin     bool           `gorm:"default:false" json:"is_admin"`
	LastLoginAt *time.Time     `json:"last_login_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Roles []*Role `gorm:"many2many:user_roles;" json:"roles,omitempty"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// SetPassword 设置密码（加密）
func (u *User) SetPassword(password string) error {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}
	u.Password = hashedPassword
	return nil
}

// CheckPassword 检查密码是否正确
func (u *User) CheckPassword(password string) bool {
	return utils.CheckPasswordHash(password, u.Password)
}

// Role 角色模型
type Role struct {
	ID          int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"type:varchar(64);uniqueIndex;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Permissions []*Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
}

// TableName 指定表名
func (Role) TableName() string {
	return "roles"
}

// Permission 权限模型
type Permission struct {
	ID          int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"type:varchar(64);uniqueIndex;not null" json:"name"`
	Resource    string    `gorm:"type:varchar(64);not null" json:"resource"`
	Action      string    `gorm:"type:varchar(32);not null" json:"action"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName 指定表名
func (Permission) TableName() string {
	return "permissions"
}

// UserRole 用户角色关联
type UserRole struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int       `gorm:"not null;index" json:"user_id"`
	RoleID    int       `gorm:"not null;index" json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 指定表名
func (UserRole) TableName() string {
	return "user_roles"
}

// RolePermission 角色权限关联
type RolePermission struct {
	ID           int       `gorm:"primaryKey;autoIncrement" json:"id"`
	RoleID       int       `gorm:"not null;index" json:"role_id"`
	PermissionID int       `gorm:"not null;index" json:"permission_id"`
	CreatedAt    time.Time `json:"created_at"`
}

// TableName 指定表名
func (RolePermission) TableName() string {
	return "role_permissions"
}


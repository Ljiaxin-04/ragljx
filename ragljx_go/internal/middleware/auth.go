package middleware

import (
	"ragljx/internal/pkg/jwt"
	"ragljx/internal/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuth JWT 认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 优先从 Header 获取 Bearer token
		authHeader := c.GetHeader("Authorization")
		var tokenString string

		if authHeader != "" {
			// 解析 Bearer token
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				response.Unauthorized(c, "invalid authorization header")
				c.Abort()
				return
			}
			tokenString = parts[1]
		} else {
			// 兼容 EventSource 等无法自定义 header 的场景，从 query 参数 token 获取
			tokenString = c.Query("token")
			if tokenString == "" {
				response.Unauthorized(c, "missing authorization header")
				c.Abort()
				return
			}
		}

		// 验证 token
		claims, err := jwt.ParseToken(tokenString)
		if err != nil {
			response.Unauthorized(c, "invalid token")
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("is_admin", claims.IsAdmin)

		c.Next()
	}
}

// OptionalAuth 可选的认证中间件（允许未登录访问）
func OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				claims, err := jwt.ParseToken(parts[1])
				if err == nil {
					c.Set("user_id", claims.UserID)
					c.Set("username", claims.Username)
					c.Set("is_admin", claims.IsAdmin)
				}
			}
		}
		c.Next()
	}
}

// AdminAuth 管理员认证中间件
func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAdmin, exists := c.Get("is_admin")
		if !exists || !isAdmin.(bool) {
			response.Forbidden(c, "admin access required")
			c.Abort()
			return
		}
		c.Next()
	}
}

// GetUserID 从上下文获取用户 ID
func GetUserID(c *gin.Context) (int, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	return userID.(int), true
}

// GetUsername 从上下文获取用户名
func GetUsername(c *gin.Context) (string, bool) {
	username, exists := c.Get("username")
	if !exists {
		return "", false
	}
	return username.(string), true
}

// IsAdmin 判断是否为管理员
func IsAdmin(c *gin.Context) bool {
	isAdmin, exists := c.Get("is_admin")
	if !exists {
		return false
	}
	return isAdmin.(bool)
}

// ValidateToken 验证 token 并返回用户 ID
func ValidateToken(tokenString string) (int, error) {
	claims, err := jwt.ParseToken(tokenString)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}


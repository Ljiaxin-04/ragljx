package auth

import (
	"ragljx/internal/middleware"
	"ragljx/internal/pkg/response"
	"ragljx/internal/service"
	"ragljx/ioc"
	"ragljx/ioc/config/datasource"
	httpConfig "ragljx/ioc/config/http"

	"github.com/gin-gonic/gin"
)

func init() {
	ioc.Api().Registry(&AuthAPI{})
}

type AuthAPI struct {
	ioc.ObjectImpl
	authService *service.AuthService
}

func (a *AuthAPI) Name() string {
	return "auth_api"
}

func (a *AuthAPI) Init() error {
	db := datasource.Get()
	a.authService = service.NewAuthService(db)

	// 注册路由
	engine := httpConfig.RootRouter()
	a.Registry(engine)

	return nil
}

func (a *AuthAPI) Registry(r gin.IRouter) {
	api := r.Group("/api/v1/auth")
	{
		api.POST("/login", a.Login)
		api.POST("/refresh", a.RefreshToken)
		api.GET("/me", middleware.JWTAuth(), a.GetCurrentUser)
	}
}

// Login 用户登录
func (a *AuthAPI) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	resp, err := a.authService.Login(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, 401, err.Error())
		return
	}

	response.Success(c, resp)
}

// RefreshToken 刷新 Token
func (a *AuthAPI) RefreshToken(c *gin.Context) {
	var req service.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	resp, err := a.authService.RefreshToken(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, 401, err.Error())
		return
	}

	response.Success(c, resp)
}

// GetCurrentUser 获取当前用户信息
func (a *AuthAPI) GetCurrentUser(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "unauthorized")
		return
	}

	user, err := a.authService.GetCurrentUser(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, user)
}


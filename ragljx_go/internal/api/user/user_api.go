package user

import (
	"ragljx/internal/middleware"
	"ragljx/internal/pkg/response"
	"ragljx/internal/pkg/utils"
	"ragljx/internal/service"
	"ragljx/ioc"
	"ragljx/ioc/config/datasource"
	httpConfig "ragljx/ioc/config/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func init() {
	ioc.Api().Registry(&UserAPI{})
}

type UserAPI struct {
	ioc.ObjectImpl
	userService *service.UserService
}

func (u *UserAPI) Name() string {
	return "user_api"
}

func (u *UserAPI) Init() error {
	db := datasource.Get()
	u.userService = service.NewUserService(db)

	// 注册路由
	engine := httpConfig.RootRouter()
	u.Registry(engine)

	return nil
}

func (u *UserAPI) Registry(r gin.IRouter) {
	api := r.Group("/api/v1/users", middleware.JWTAuth())
	{
		api.POST("", middleware.AdminAuth(), u.Create)
		api.GET("", u.List)
		api.GET("/:id", u.GetByID)
		api.PUT("/:id", u.Update)
		api.DELETE("/:id", middleware.AdminAuth(), u.Delete)
		api.POST("/:id/password", u.ChangePassword)
		api.POST("/:id/roles", middleware.AdminAuth(), u.AssignRoles)
	}
}

// Create 创建用户
func (u *UserAPI) Create(c *gin.Context) {
	var req service.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	user, err := u.userService.Create(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, user)
}

// GetByID 根据 ID 获取用户
func (u *UserAPI) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	user, err := u.userService.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, user)
}

// List 获取用户列表
func (u *UserAPI) List(c *gin.Context) {
	var pagination utils.Pagination
	if err := c.ShouldBindQuery(&pagination); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	keyword := c.Query("keyword")

	users, total, err := u.userService.List(c.Request.Context(), pagination.Page, pagination.GetLimit(), keyword)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Page(c, total, pagination.Page, pagination.GetLimit(), users)
}

// Update 更新用户
func (u *UserAPI) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	var req service.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	user, err := u.userService.Update(c.Request.Context(), id, &req)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, user)
}

// Delete 删除用户
func (u *UserAPI) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	if err := u.userService.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, nil)
}

// ChangePassword 修改密码
func (u *UserAPI) ChangePassword(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	// 检查权限：只能修改自己的密码或管理员可以修改任何人的密码
	userID, _ := middleware.GetUserID(c)
	if userID != id && !middleware.IsAdmin(c) {
		response.Forbidden(c, "forbidden")
		return
	}

	var req service.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := u.userService.ChangePassword(c.Request.Context(), id, &req); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, nil)
}

// AssignRoles 分配角色
func (u *UserAPI) AssignRoles(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	var req struct {
		RoleIDs []int `json:"role_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := u.userService.AssignRoles(c.Request.Context(), id, req.RoleIDs); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, nil)
}


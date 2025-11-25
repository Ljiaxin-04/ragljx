package knowledge_base

import (
	"ragljx/internal/middleware"
	"ragljx/internal/pkg/response"
	"ragljx/internal/pkg/utils"
	"ragljx/internal/service"
	"ragljx/ioc"
	"ragljx/ioc/config/datasource"
	httpConfig "ragljx/ioc/config/http"

	"github.com/gin-gonic/gin"
)

func init() {
	ioc.Api().Registry(&KnowledgeBaseAPI{})
}

type KnowledgeBaseAPI struct {
	ioc.ObjectImpl
	kbService *service.KnowledgeBaseService
}

func (k *KnowledgeBaseAPI) Name() string {
	return "knowledge_base_api"
}

func (k *KnowledgeBaseAPI) Init() error {
	db := datasource.Get()
	k.kbService = service.NewKnowledgeBaseService(db)

	// 注册路由
	engine := httpConfig.RootRouter()
	k.Registry(engine)

	return nil
}

func (k *KnowledgeBaseAPI) Registry(r gin.IRouter) {
	api := r.Group("/api/v1/knowledge-bases", middleware.JWTAuth())
	{
		api.POST("", k.Create)
		api.GET("", k.List)
		api.GET("/:id", k.GetByID)
		api.PUT("/:id", k.Update)
		api.DELETE("/:id", k.Delete)
	}
}

// Create 创建知识库
func (k *KnowledgeBaseAPI) Create(c *gin.Context) {
	var req service.CreateKBRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID, _ := middleware.GetUserID(c)
	kb, err := k.kbService.Create(c.Request.Context(), &req, userID)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, kb)
}

// GetByID 根据 ID 获取知识库
func (k *KnowledgeBaseAPI) GetByID(c *gin.Context) {
	id := c.Param("id")

	kb, err := k.kbService.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, kb)
}

// List 获取知识库列表
func (k *KnowledgeBaseAPI) List(c *gin.Context) {
	var pagination utils.Pagination
	if err := c.ShouldBindQuery(&pagination); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	status := c.Query("status")

	kbs, total, err := k.kbService.List(c.Request.Context(), pagination.Page, pagination.GetLimit(), status)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Page(c, total, pagination.Page, pagination.GetLimit(), kbs)
}

// Update 更新知识库
func (k *KnowledgeBaseAPI) Update(c *gin.Context) {
	id := c.Param("id")

	var req service.UpdateKBRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	kb, err := k.kbService.Update(c.Request.Context(), id, &req)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, kb)
}

// Delete 删除知识库
func (k *KnowledgeBaseAPI) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := k.kbService.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, nil)
}


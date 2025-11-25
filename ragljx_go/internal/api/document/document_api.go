package document

import (
	"ragljx/internal/middleware"
	"ragljx/internal/pkg/response"
	"ragljx/internal/pkg/utils"
	"ragljx/internal/service"
	"ragljx/ioc"
	"ragljx/ioc/config/datasource"
	httpConfig "ragljx/ioc/config/http"
	kafkaConfig "ragljx/ioc/config/kafka"
	minioConfig "ragljx/ioc/config/minio"

	"github.com/gin-gonic/gin"
)

func init() {
	ioc.Api().Registry(&DocumentAPI{})
}

type DocumentAPI struct {
	ioc.ObjectImpl
	docService *service.DocumentService
}

func (d *DocumentAPI) Name() string {
	return "document_api"
}

func (d *DocumentAPI) Init() error {
	db := datasource.Get()
	minioObj := minioConfig.Get()
	kafkaObj := kafkaConfig.Get()

	var minioClient = minioObj.Client()
	var bucketName = minioObj.Bucket()
	var kafkaWriter = kafkaObj.Producer("document-tasks")

	d.docService = service.NewDocumentService(db, minioClient, bucketName, kafkaWriter)

	// 注册路由
	engine := httpConfig.RootRouter()
	d.Registry(engine)

	return nil
}

func (d *DocumentAPI) Registry(r gin.IRouter) {
	api := r.Group("/api/v1/documents", middleware.JWTAuth())
	{
		api.POST("/upload", d.Upload)
		api.GET("", d.List)
		api.GET("/:id", d.GetByID)
		api.GET("/:id/download", d.Download)
		api.POST("/:id/vectorize", d.Vectorize)
		api.DELETE("/:id", d.Delete)
	}
}

// Upload 上传文档
func (d *DocumentAPI) Upload(c *gin.Context) {
	var req service.UploadRequest
	if err := c.ShouldBind(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID, _ := middleware.GetUserID(c)
	doc, err := d.docService.Upload(c.Request.Context(), &req, userID)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, doc)
}

// GetByID 根据 ID 获取文档
func (d *DocumentAPI) GetByID(c *gin.Context) {
	id := c.Param("id")

	doc, err := d.docService.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, doc)
}

// List 获取文档列表
func (d *DocumentAPI) List(c *gin.Context) {
	var pagination utils.Pagination
	if err := c.ShouldBindQuery(&pagination); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	kbID := c.Query("knowledge_base_id")
	if kbID == "" {
		response.BadRequest(c, "knowledge_base_id is required")
		return
	}

	status := c.Query("status")

	docs, total, err := d.docService.ListByKnowledgeBase(c.Request.Context(), kbID, pagination.Page, pagination.GetLimit(), status)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Page(c, total, pagination.Page, pagination.GetLimit(), docs)
}

// Download 下载文档
func (d *DocumentAPI) Download(c *gin.Context) {
	id := c.Param("id")

	reader, filename, err := d.docService.Download(c.Request.Context(), id)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}
	defer reader.Close()

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.DataFromReader(200, -1, "application/octet-stream", reader, nil)
}

// Vectorize 向量化文档
func (d *DocumentAPI) Vectorize(c *gin.Context) {
	id := c.Param("id")

	if err := d.docService.Vectorize(c.Request.Context(), id); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, nil)
}

// Delete 删除文档
func (d *DocumentAPI) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := d.docService.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, nil)
}


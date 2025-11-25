package knowledge_base

import (
	pb "ragljx/proto/rag"
	"ragljx/internal/middleware"
	"ragljx/internal/pkg/response"
	"ragljx/internal/pkg/utils"
	"ragljx/internal/service"
	"ragljx/ioc"
	"ragljx/ioc/config/datasource"
	grpcConfig "ragljx/ioc/config/grpc"
	httpConfig "ragljx/ioc/config/http"
	kafkaConfig "ragljx/ioc/config/kafka"
	minioConfig "ragljx/ioc/config/minio"

	"github.com/gin-gonic/gin"
)

func init() {
	ioc.Api().Registry(&KnowledgeBaseAPI{})
}

type KnowledgeBaseAPI struct {
	ioc.ObjectImpl
	kbService  *service.KnowledgeBaseService
	docService *service.DocumentService
}

func (k *KnowledgeBaseAPI) Name() string {
	return "knowledge_base_api"
}

func (k *KnowledgeBaseAPI) Init() error {
	db := datasource.Get()
	k.kbService = service.NewKnowledgeBaseService(db)

	// 初始化文档服务
	minioObj := minioConfig.Get()
	kafkaObj := kafkaConfig.Get()
	grpcConn := grpcConfig.Get()

	var minioClient = minioObj.Client()
	var bucketName = minioObj.Bucket()
	var kafkaWriter = kafkaObj.Producer("document-tasks")
	var grpcClient = pb.NewRAGServiceClient(grpcConn)

	k.docService = service.NewDocumentService(db, minioClient, bucketName, kafkaWriter, grpcClient)

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

		// 文档相关的嵌套路由
		api.GET("/:id/documents", k.GetDocuments)
		api.POST("/:id/documents/upload", k.UploadDocument)
		api.GET("/:id/documents/:docId", k.GetDocument)
		api.DELETE("/:id/documents/:docId", k.DeleteDocument)
		api.POST("/:id/documents/:docId/vectorize", k.VectorizeDocument)
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

// GetDocuments 获取知识库的文档列表
func (k *KnowledgeBaseAPI) GetDocuments(c *gin.Context) {
	kbID := c.Param("id")

	var pagination utils.Pagination
	if err := c.ShouldBindQuery(&pagination); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	status := c.Query("status")

	docs, total, err := k.docService.ListByKnowledgeBase(c.Request.Context(), kbID, pagination.Page, pagination.GetLimit(), status)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Page(c, total, pagination.Page, pagination.GetLimit(), docs)
}

// UploadDocument 上传文档到知识库
func (k *KnowledgeBaseAPI) UploadDocument(c *gin.Context) {
	kbID := c.Param("id")

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "file is required")
		return
	}

	// 构建上传请求
	req := service.UploadRequest{
		KnowledgeBaseID: kbID,
		File:            file,
	}

	userID, _ := middleware.GetUserID(c)
	doc, err := k.docService.Upload(c.Request.Context(), &req, userID)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, doc)
}

// GetDocument 获取文档详情
func (k *KnowledgeBaseAPI) GetDocument(c *gin.Context) {
	docID := c.Param("docId")

	doc, err := k.docService.GetByID(c.Request.Context(), docID)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, doc)
}

// DeleteDocument 删除文档
func (k *KnowledgeBaseAPI) DeleteDocument(c *gin.Context) {
	docID := c.Param("docId")

	if err := k.docService.Delete(c.Request.Context(), docID); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, nil)
}

// VectorizeDocument 向量化文档
func (k *KnowledgeBaseAPI) VectorizeDocument(c *gin.Context) {
	docID := c.Param("docId")

	if err := k.docService.Vectorize(c.Request.Context(), docID); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, nil)
}


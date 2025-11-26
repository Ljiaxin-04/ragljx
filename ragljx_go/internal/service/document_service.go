package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	pb "ragljx/proto/rag"
	"ragljx/internal/model"
	"ragljx/internal/pkg/errors"
	"ragljx/internal/pkg/utils"
	"ragljx/internal/repository"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/segmentio/kafka-go"
	"gorm.io/gorm"
)

type DocumentService struct {
	docRepo     *repository.DocumentRepository
	kbRepo      *repository.KnowledgeBaseRepository
	minioClient *minio.Client
	bucketName  string
	kafkaWriter *kafka.Writer
	grpcClient  pb.RAGServiceClient
}

func NewDocumentService(db *gorm.DB, minioClient *minio.Client, bucketName string, kafkaWriter *kafka.Writer, grpcClient pb.RAGServiceClient) *DocumentService {
	return &DocumentService{
		docRepo:     repository.NewDocumentRepository(db),
		kbRepo:      repository.NewKnowledgeBaseRepository(db),
		minioClient: minioClient,
		bucketName:  bucketName,
		kafkaWriter: kafkaWriter,
		grpcClient:  grpcClient,
	}
}

// UploadRequest 上传文档请求
type UploadRequest struct {
	KnowledgeBaseID string                `form:"knowledge_base_id" binding:"required"`
	File            *multipart.FileHeader `form:"file" binding:"required"`
}

// 支持的文件扩展名（小写）
var allowedExtensions = map[string]struct{}{
	".txt": {}, ".md": {}, ".pdf": {}, ".docx": {}, ".xlsx": {}, ".pptx": {}, ".html": {}, ".htm": {},
	".csv": {}, ".json": {}, ".xml": {}, ".rtf": {},
}

// 最大文件尺寸（50MB）
const maxUploadSize = int64(50 * 1024 * 1024)

// Upload 上传文档
func (s *DocumentService) Upload(ctx context.Context, req *UploadRequest, userID int) (*model.KnowledgeDocument, error) {
	// 检查知识库是否存在
	kb, err := s.kbRepo.GetByID(ctx, req.KnowledgeBaseID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrKBNotFound
		}
		return nil, errors.Wrap(500, "failed to get knowledge base", err)
	}

	// 基础校验：大小与扩展名
	if req.File.Size > maxUploadSize {
		return nil, errors.New(400, "file too large: maximum 50MB")
	}
	ext := strings.ToLower(filepath.Ext(req.File.Filename))
	if _, ok := allowedExtensions[ext]; !ok {
		return nil, errors.New(400, fmt.Sprintf("unsupported file type: %s", ext))
	}

	// 打开文件
	file, err := req.File.Open()
	if err != nil {
		return nil, errors.Wrap(500, "failed to open file", err)
	}
	defer file.Close()

	// 计算文件哈希
	checksum, err := utils.SHA256Reader(file)
	if err != nil {
		return nil, errors.Wrap(500, "failed to calculate checksum", err)
	}

	// 检查文件是否已存在
	existingDoc, err := s.docRepo.GetByChecksum(ctx, req.KnowledgeBaseID, checksum)
	if err == nil && existingDoc != nil {
		return existingDoc, nil
	}

	// 重置文件指针
	file.Seek(0, io.SeekStart)

	// 生成对象键
	objectKey := fmt.Sprintf("documents/%s/%s_%s", req.KnowledgeBaseID, time.Now().Format("20060102150405"), req.File.Filename)

	// 上传到 MinIO
	_, err = s.minioClient.PutObject(ctx, s.bucketName, objectKey, file, req.File.Size, minio.PutObjectOptions{
		ContentType: req.File.Header.Get("Content-Type"),
	})
	if err != nil {
		return nil, errors.Wrap(500, "failed to upload file to minio", err)
	}

	// 创建文档记录
	doc := &model.KnowledgeDocument{
		KnowledgeBaseID: req.KnowledgeBaseID,
		Title:           req.File.Filename,
		ObjectKey:       objectKey,
		Size:            req.File.Size,
		Mime:            req.File.Header.Get("Content-Type"),
		Checksum:        checksum,
		ParsingStatus:   "processing",  // 使用前端期望的状态值
		CreatedByID:     &userID,
	}

	if err := s.docRepo.Create(ctx, doc); err != nil {
		return nil, errors.Wrap(500, "failed to create document", err)
	}

	// 更新知识库统计
	s.kbRepo.IncrementDocumentCount(ctx, req.KnowledgeBaseID, 1)
	s.kbRepo.UpdateTotalSize(ctx, req.KnowledgeBaseID, req.File.Size)

	// 检查 gRPC 客户端是否可用
	if s.grpcClient == nil {
		s.docRepo.UpdateStatus(ctx, doc.ID, "failed", "gRPC client is not initialized")
		return doc, nil
	}

	// 异步处理文档解析和向量化
	go s.processDocument(context.Background(), doc, kb)

	return doc, nil
}

// processDocument 处理文档解析和向量化（异步）
func (s *DocumentService) processDocument(ctx context.Context, doc *model.KnowledgeDocument, kb *model.KnowledgeBase) {
	// 创建带超时的 context（10 分钟超时）
	ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	// 1. 从 MinIO 下载文件内容
	obj, err := s.minioClient.GetObject(ctx, s.bucketName, doc.ObjectKey, minio.GetObjectOptions{})
	if err != nil {
		s.docRepo.UpdateStatus(context.Background(), doc.ID, "failed", fmt.Sprintf("failed to get file from minio: %v", err))
		return
	}
	defer obj.Close()

	fileContent, err := io.ReadAll(obj)
	if err != nil {
		s.docRepo.UpdateStatus(context.Background(), doc.ID, "failed", fmt.Sprintf("failed to read file: %v", err))
		return
	}

	// 2. 调用 Python gRPC 服务解析文档
	parseResp, err := s.grpcClient.ParseDocument(ctx, &pb.ParseDocumentRequest{
		DocumentId:  doc.ID,
		ObjectKey:   doc.ObjectKey,
		FileContent: fileContent,
		MimeType:    doc.Mime,
	})
	if err != nil {
		s.docRepo.UpdateStatus(context.Background(), doc.ID, "failed", fmt.Sprintf("failed to parse document: %v", err))
		return
	}

	if !parseResp.Success {
		s.docRepo.UpdateStatus(context.Background(), doc.ID, "failed", parseResp.ErrorMessage)
		return
	}

	// 检查解析内容是否为空
	if parseResp.Content == "" {
		s.docRepo.UpdateStatus(context.Background(), doc.ID, "failed", "parsed content is empty")
		return
	}

	// 3. 调用 Python gRPC 服务向量化文档
	// 使用知识库的 english_name 作为 collection_name
	collectionName := kb.EnglishName
	if collectionName == "" {
		s.docRepo.UpdateStatus(context.Background(), doc.ID, "failed", "knowledge base english_name is empty")
		return
	}

	vectorizeResp, err := s.grpcClient.VectorizeDocument(ctx, &pb.VectorizeDocumentRequest{
		DocumentId:      doc.ID,
		Content:         parseResp.Content,
		KnowledgeBaseId: doc.KnowledgeBaseID,
		CollectionName:  collectionName,
		Title:           doc.Title,
		ObjectKey:       doc.ObjectKey,
	})
	if err != nil {
		s.docRepo.UpdateStatus(context.Background(), doc.ID, "failed", fmt.Sprintf("failed to vectorize document: %v", err))
		return
	}

	if !vectorizeResp.Success {
		s.docRepo.UpdateStatus(context.Background(), doc.ID, "failed", vectorizeResp.ErrorMessage)
		return
	}

	// 4. 更新文档状态为成功，并更新分块数
	updates := map[string]interface{}{
		"parsing_status": "completed",  // 使用前端期望的状态值
		"chunk_count":    vectorizeResp.ChunkCount,
	}
	s.docRepo.UpdateWithMap(context.Background(), doc.ID, updates)
}

// GetByID 根据 ID 获取文档
func (s *DocumentService) GetByID(ctx context.Context, id string) (*model.KnowledgeDocument, error) {
	doc, err := s.docRepo.GetByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrDocumentNotFound
		}
		return nil, errors.Wrap(500, "failed to get document", err)
	}
	return doc, nil
}

// ListByKnowledgeBase 根据知识库获取文档列表
func (s *DocumentService) ListByKnowledgeBase(ctx context.Context, kbID string, page, pageSize int, status string) ([]*model.KnowledgeDocument, int64, error) {
	offset := (page - 1) * pageSize
	docs, total, err := s.docRepo.ListByKnowledgeBase(ctx, kbID, offset, pageSize, status)
	if err != nil {
		return nil, 0, errors.Wrap(500, "failed to list documents", err)
	}
	return docs, total, nil
}

// Delete 删除文档
func (s *DocumentService) Delete(ctx context.Context, id string) error {
	doc, err := s.docRepo.GetByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.ErrDocumentNotFound
		}
		return errors.Wrap(500, "failed to get document", err)
	}

	// 获取知识库信息（需要 english_name 作为 collection name）
	kb, err := s.kbRepo.GetByID(ctx, doc.KnowledgeBaseID)
	if err != nil {
		return errors.Wrap(500, "failed to get knowledge base", err)
	}

	// 从 Qdrant 删除向量（如果 gRPC 客户端可用且文档已向量化）
	if s.grpcClient != nil && doc.ParsingStatus == "completed" && kb.EnglishName != "" {
		deleteResp, err := s.grpcClient.DeleteDocumentVectors(ctx, &pb.DeleteDocumentVectorsRequest{
			CollectionName: kb.EnglishName,
			DocumentIds:    []string{doc.ID},
		})
		if err != nil {
			// 记录错误但不中断删除流程
			// 可以考虑记录日志
		} else if !deleteResp.Success {
			// 记录错误但不中断删除流程
		}
	}

	// 从 MinIO 删除文件
	if err := s.minioClient.RemoveObject(ctx, s.bucketName, doc.ObjectKey, minio.RemoveObjectOptions{}); err != nil {
		// 记录错误但不中断删除流程
	}

	// 删除文档记录
	if err := s.docRepo.Delete(ctx, id); err != nil {
		return errors.Wrap(500, "failed to delete document", err)
	}

	// 更新知识库统计
	s.kbRepo.IncrementDocumentCount(ctx, doc.KnowledgeBaseID, -1)
	s.kbRepo.UpdateTotalSize(ctx, doc.KnowledgeBaseID, -doc.Size)

	return nil
}

// Download 下载文档
func (s *DocumentService) Download(ctx context.Context, id string) (io.ReadCloser, string, error) {
	doc, err := s.docRepo.GetByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, "", errors.ErrDocumentNotFound
		}
		return nil, "", errors.Wrap(500, "failed to get document", err)
	}

	// 从 MinIO 下载文件
	object, err := s.minioClient.GetObject(ctx, s.bucketName, doc.ObjectKey, minio.GetObjectOptions{})
	if err != nil {
		return nil, "", errors.Wrap(500, "failed to download file from minio", err)
	}

	return object, filepath.Base(doc.Title), nil
}

// Vectorize 向量化文档
func (s *DocumentService) Vectorize(ctx context.Context, id string) error {
	doc, err := s.docRepo.GetByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.ErrDocumentNotFound
		}
		return errors.Wrap(500, "failed to get document", err)
	}

	// 检查文档状态
	if doc.ParsingStatus != "completed" {
		return errors.New(400, "document is not ready for vectorization")
	}

	// 发送向量化任务到 Kafka
	taskMsg := map[string]interface{}{
		"document_id":       doc.ID,
		"knowledge_base_id": doc.KnowledgeBaseID,
		"object_key":        doc.ObjectKey,
		"task_type":         "vectorize",
	}
	msgBytes, _ := json.Marshal(taskMsg)
	if err := s.kafkaWriter.WriteMessages(ctx, kafka.Message{
		Key:   []byte(doc.ID),
		Value: msgBytes,
	}); err != nil {
		return errors.Wrap(500, "failed to send vectorize task", err)
	}

	// 更新文档状态
	if err := s.docRepo.UpdateStatus(ctx, id, "vectorizing", ""); err != nil {
		return errors.Wrap(500, "failed to update document status", err)
	}

	return nil
}

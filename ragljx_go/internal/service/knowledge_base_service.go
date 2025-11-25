package service

import (
	"context"
	pb "ragljx/proto/rag"
	"ragljx/internal/model"
	"ragljx/internal/pkg/errors"
	"ragljx/internal/repository"

	"gorm.io/gorm"
)

type KnowledgeBaseService struct {
	kbRepo     *repository.KnowledgeBaseRepository
	docRepo    *repository.DocumentRepository
	grpcClient pb.RAGServiceClient
}

func NewKnowledgeBaseService(db *gorm.DB, grpcClient pb.RAGServiceClient) *KnowledgeBaseService {
	return &KnowledgeBaseService{
		kbRepo:     repository.NewKnowledgeBaseRepository(db),
		docRepo:    repository.NewDocumentRepository(db),
		grpcClient: grpcClient,
	}
}

// CreateKBRequest 创建知识库请求
type CreateKBRequest struct {
	Name           string `json:"name" binding:"required"`
	EnglishName    string `json:"english_name" binding:"required"`
	Description    string `json:"description"`
	EmbeddingModel string `json:"embedding_model"`
}

// UpdateKBRequest 更新知识库请求
type UpdateKBRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

// Create 创建知识库
func (s *KnowledgeBaseService) Create(ctx context.Context, req *CreateKBRequest, userID int) (*model.KnowledgeBase, error) {
	// 检查英文名是否存在
	existingKB, err := s.kbRepo.GetByEnglishName(ctx, req.EnglishName)
	if err == nil && existingKB != nil {
		return nil, errors.New(400, "english name already exists")
	}

	// 创建知识库
	kb := &model.KnowledgeBase{
		Name:           req.Name,
		EnglishName:    req.EnglishName,
		Description:    req.Description,
		EmbeddingModel: req.EmbeddingModel,
		Status:         "active",
		CreatedByID:    &userID,
	}

	if kb.EmbeddingModel == "" {
		kb.EmbeddingModel = "text-embedding-3-small"
	}

	if err := s.kbRepo.Create(ctx, kb); err != nil {
		return nil, errors.Wrap(500, "failed to create knowledge base", err)
	}

	return kb, nil
}

// GetByID 根据 ID 获取知识库
func (s *KnowledgeBaseService) GetByID(ctx context.Context, id string) (*model.KnowledgeBase, error) {
	kb, err := s.kbRepo.GetByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrKBNotFound
		}
		return nil, errors.Wrap(500, "failed to get knowledge base", err)
	}
	return kb, nil
}

// List 获取知识库列表
func (s *KnowledgeBaseService) List(ctx context.Context, page, pageSize int, status string) ([]*model.KnowledgeBase, int64, error) {
	offset := (page - 1) * pageSize
	kbs, total, err := s.kbRepo.List(ctx, offset, pageSize, status)
	if err != nil {
		return nil, 0, errors.Wrap(500, "failed to list knowledge bases", err)
	}
	return kbs, total, nil
}

// Update 更新知识库
func (s *KnowledgeBaseService) Update(ctx context.Context, id string, req *UpdateKBRequest) (*model.KnowledgeBase, error) {
	kb, err := s.kbRepo.GetByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrKBNotFound
		}
		return nil, errors.Wrap(500, "failed to get knowledge base", err)
	}

	// 更新字段
	if req.Name != "" {
		kb.Name = req.Name
	}
	if req.Description != "" {
		kb.Description = req.Description
	}
	if req.Status != "" {
		kb.Status = req.Status
	}

	if err := s.kbRepo.Update(ctx, kb); err != nil {
		return nil, errors.Wrap(500, "failed to update knowledge base", err)
	}

	return kb, nil
}

// Delete 删除知识库
func (s *KnowledgeBaseService) Delete(ctx context.Context, id string) error {
	// 检查知识库是否存在
	kb, err := s.kbRepo.GetByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.ErrKBNotFound
		}
		return errors.Wrap(500, "failed to get knowledge base", err)
	}

	// 检查是否还有文档
	_, total, err := s.docRepo.ListByKnowledgeBase(ctx, id, 0, 1, "")
	if err != nil {
		return errors.Wrap(500, "failed to check documents", err)
	}
	if total > 0 {
		return errors.New(400, "knowledge base has documents, please delete them first")
	}

	// 删除 Qdrant collection（如果 gRPC 客户端可用且 english_name 不为空）
	if s.grpcClient != nil && kb.EnglishName != "" {
		// 注意：这里需要添加一个新的 gRPC 方法来删除 collection
		// 暂时跳过，因为 proto 中可能没有定义这个方法
		// 可以在 Python 端添加 DeleteCollection 方法
	}

	// 删除知识库记录
	if err := s.kbRepo.Delete(ctx, id); err != nil {
		return errors.Wrap(500, "failed to delete knowledge base", err)
	}

	return nil
}


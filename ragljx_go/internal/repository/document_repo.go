package repository

import (
	"context"
	"ragljx/internal/model"

	"gorm.io/gorm"
)

type DocumentRepository struct {
	db *gorm.DB
}

func NewDocumentRepository(db *gorm.DB) *DocumentRepository {
	return &DocumentRepository{db: db}
}

// Create 创建文档
func (r *DocumentRepository) Create(ctx context.Context, doc *model.KnowledgeDocument) error {
	return r.db.WithContext(ctx).Create(doc).Error
}

// GetByID 根据 ID 获取文档
func (r *DocumentRepository) GetByID(ctx context.Context, id string) (*model.KnowledgeDocument, error) {
	var doc model.KnowledgeDocument
	err := r.db.WithContext(ctx).Preload("KnowledgeBase").Preload("CreatedBy").First(&doc, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

// ListByKnowledgeBase 根据知识库获取文档列表
func (r *DocumentRepository) ListByKnowledgeBase(ctx context.Context, kbID string, offset, limit int, status string) ([]*model.KnowledgeDocument, int64, error) {
	var docs []*model.KnowledgeDocument
	var total int64

	query := r.db.WithContext(ctx).Model(&model.KnowledgeDocument{}).Where("knowledge_base_id = ?", kbID)
	if status != "" {
		query = query.Where("parsing_status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("CreatedBy").Offset(offset).Limit(limit).Order("created_at DESC").Find(&docs).Error
	return docs, total, err
}

// Update 更新文档
func (r *DocumentRepository) Update(ctx context.Context, doc *model.KnowledgeDocument) error {
	return r.db.WithContext(ctx).Save(doc).Error
}

// UpdateStatus 更新文档状态
func (r *DocumentRepository) UpdateStatus(ctx context.Context, id, status string, errorMsg string) error {
	updates := map[string]interface{}{
		"parsing_status": status,
	}
	if errorMsg != "" {
		updates["error_message"] = errorMsg
	}
	return r.db.WithContext(ctx).Model(&model.KnowledgeDocument{}).Where("id = ?", id).Updates(updates).Error
}

// Delete 删除文档
func (r *DocumentRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&model.KnowledgeDocument{}, "id = ?", id).Error
}

// GetByChecksum 根据校验和获取文档
func (r *DocumentRepository) GetByChecksum(ctx context.Context, kbID, checksum string) (*model.KnowledgeDocument, error) {
	var doc model.KnowledgeDocument
	err := r.db.WithContext(ctx).Where("knowledge_base_id = ? AND checksum = ?", kbID, checksum).First(&doc).Error
	if err != nil {
		return nil, err
	}
	return &doc, nil
}


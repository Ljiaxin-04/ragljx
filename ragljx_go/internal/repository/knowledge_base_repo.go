package repository

import (
	"context"
	"ragljx/internal/model"

	"gorm.io/gorm"
)

type KnowledgeBaseRepository struct {
	db *gorm.DB
}

func NewKnowledgeBaseRepository(db *gorm.DB) *KnowledgeBaseRepository {
	return &KnowledgeBaseRepository{db: db}
}

// Create 创建知识库
func (r *KnowledgeBaseRepository) Create(ctx context.Context, kb *model.KnowledgeBase) error {
	return r.db.WithContext(ctx).Create(kb).Error
}

// GetByID 根据 ID 获取知识库
func (r *KnowledgeBaseRepository) GetByID(ctx context.Context, id string) (*model.KnowledgeBase, error) {
	var kb model.KnowledgeBase
	err := r.db.WithContext(ctx).Preload("CreatedBy").First(&kb, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &kb, nil
}

// GetByEnglishName 根据英文名获取知识库
func (r *KnowledgeBaseRepository) GetByEnglishName(ctx context.Context, englishName string) (*model.KnowledgeBase, error) {
	var kb model.KnowledgeBase
	err := r.db.WithContext(ctx).Where("english_name = ?", englishName).First(&kb).Error
	if err != nil {
		return nil, err
	}
	return &kb, nil
}

// List 获取知识库列表
func (r *KnowledgeBaseRepository) List(ctx context.Context, offset, limit int, status string) ([]*model.KnowledgeBase, int64, error) {
	var kbs []*model.KnowledgeBase
	var total int64

	query := r.db.WithContext(ctx).Model(&model.KnowledgeBase{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("CreatedBy").Offset(offset).Limit(limit).Order("created_at DESC").Find(&kbs).Error
	return kbs, total, err
}

// Update 更新知识库
func (r *KnowledgeBaseRepository) Update(ctx context.Context, kb *model.KnowledgeBase) error {
	return r.db.WithContext(ctx).Save(kb).Error
}

// Delete 删除知识库
func (r *KnowledgeBaseRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&model.KnowledgeBase{}, "id = ?", id).Error
}

// IncrementDocumentCount 增加文档计数
func (r *KnowledgeBaseRepository) IncrementDocumentCount(ctx context.Context, id string, delta int) error {
	return r.db.WithContext(ctx).Model(&model.KnowledgeBase{}).
		Where("id = ?", id).
		UpdateColumn("document_count", gorm.Expr("document_count + ?", delta)).Error
}

// UpdateTotalSize 更新总大小
func (r *KnowledgeBaseRepository) UpdateTotalSize(ctx context.Context, id string, delta int64) error {
	return r.db.WithContext(ctx).Model(&model.KnowledgeBase{}).
		Where("id = ?", id).
		UpdateColumn("total_size", gorm.Expr("total_size + ?", delta)).Error
}


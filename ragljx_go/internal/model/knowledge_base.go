package model

import (
	"time"

	"gorm.io/gorm"
)

// KnowledgeBase 知识库模型
type KnowledgeBase struct {
	ID             string         `gorm:"type:varchar(64);primaryKey" json:"id"`
	Name           string         `gorm:"type:varchar(128);not null" json:"name"`
	EnglishName    string         `gorm:"type:varchar(128);uniqueIndex;not null" json:"english_name"`
	Description    string         `gorm:"type:text" json:"description"`
	EmbeddingModel string         `gorm:"type:varchar(64);default:'text-embedding-3-small'" json:"embedding_model"`
	CollectionName string         `gorm:"type:varchar(255)" json:"collection_name"`
	DocumentCount  int            `gorm:"default:0" json:"document_count"`
	TotalSize      int64          `gorm:"default:0" json:"total_size"`
	Status         string         `gorm:"type:varchar(32);default:'active'" json:"status"`
	IsBuiltin      bool           `gorm:"default:false" json:"is_builtin"`
	CreatedByID    *int           `json:"created_by_id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	CreatedBy *User `gorm:"foreignKey:CreatedByID" json:"created_by,omitempty"`
}

// TableName 指定表名
func (KnowledgeBase) TableName() string {
	return "knowledge_bases"
}

// BeforeCreate GORM hook，创建前生成 UUID
func (kb *KnowledgeBase) BeforeCreate(tx *gorm.DB) error {
	if kb.ID == "" {
		kb.ID = generateUUID()
	}
	return nil
}

// KnowledgeDocument 知识库文档模型
type KnowledgeDocument struct {
	ID              string         `gorm:"type:varchar(64);primaryKey" json:"id"`
	KnowledgeBaseID string         `gorm:"type:varchar(64);not null;index" json:"knowledge_base_id"`
	Title           string         `gorm:"type:varchar(255);not null" json:"name"`  // 前端使用 name
	Content         string         `gorm:"type:text" json:"content"`
	ObjectKey       string         `gorm:"type:varchar(512)" json:"object_key"`
	Size            int64          `gorm:"default:0" json:"file_size"`  // 前端使用 file_size
	Mime            string         `gorm:"type:varchar(128)" json:"file_type"`  // 前端使用 file_type
	Checksum        string         `gorm:"type:varchar(128);index" json:"checksum"`
	ParsingStatus   string         `gorm:"type:varchar(32);default:'pending';index" json:"status"`  // 前端使用 status
	ChunkCount      int            `gorm:"default:0" json:"chunk_count"`  // 添加分块数字段
	ErrorMessage    string         `gorm:"type:text" json:"error_message,omitempty"`
	IsBuiltin       bool           `gorm:"default:false" json:"is_builtin"`
	CreatedByID     *int           `json:"created_by_id"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	KnowledgeBase *KnowledgeBase `gorm:"foreignKey:KnowledgeBaseID" json:"knowledge_base,omitempty"`
	CreatedBy     *User          `gorm:"foreignKey:CreatedByID" json:"created_by,omitempty"`
}

// TableName 指定表名
func (KnowledgeDocument) TableName() string {
	return "knowledge_documents"
}

// BeforeCreate GORM hook，创建前生成 UUID
func (kd *KnowledgeDocument) BeforeCreate(tx *gorm.DB) error {
	if kd.ID == "" {
		kd.ID = generateUUID()
	}
	return nil
}


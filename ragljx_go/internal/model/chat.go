package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ChatSession 对话会话模型
type ChatSession struct {
	ID                  string         `gorm:"type:varchar(64);primaryKey" json:"id"`
	UserID              int            `gorm:"not null;index" json:"user_id"`
	Title               string         `gorm:"type:varchar(255)" json:"title"`
	KnowledgeBaseIDs    StringArray    `gorm:"type:text[]" json:"knowledge_base_ids"`
	UseRAG              bool           `gorm:"default:true" json:"use_rag"`
	TopK                int            `gorm:"default:3" json:"top_k"`
	SimilarityThreshold float64        `gorm:"default:0.7" json:"similarity_threshold"`
	SimilarityWeight    float64        `gorm:"default:1.5" json:"similarity_weight"`
	Status              string         `gorm:"type:varchar(32);default:'active';index" json:"status"`
	MessageCount        int            `gorm:"default:0" json:"message_count"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"-"`

	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 指定表名
func (ChatSession) TableName() string {
	return "chat_sessions"
}

// BeforeCreate GORM hook，创建前生成 UUID
func (cs *ChatSession) BeforeCreate(tx *gorm.DB) error {
	if cs.ID == "" {
		cs.ID = generateUUID()
	}
	return nil
}

// ChatMessage 对话消息模型
type ChatMessage struct {
	ID         string     `gorm:"type:varchar(64);primaryKey" json:"id"`
	SessionID  string     `gorm:"type:varchar(64);not null;index" json:"session_id"`
	Role       string     `gorm:"type:varchar(32);not null;index" json:"role"`
	Content    string     `gorm:"type:text;not null" json:"content"`
	RAGSources JSONBArray `gorm:"type:jsonb" json:"rag_sources,omitempty"`
	TokensUsed int        `gorm:"default:0" json:"tokens_used"`
	CreatedAt  time.Time  `json:"created_at"`

	Session *ChatSession `gorm:"foreignKey:SessionID" json:"session,omitempty"`
}

// TableName 指定表名
func (ChatMessage) TableName() string {
	return "chat_messages"
}

// BeforeCreate GORM hook，创建前生成 UUID
func (cm *ChatMessage) BeforeCreate(tx *gorm.DB) error {
	if cm.ID == "" {
		cm.ID = generateUUID()
	}
	return nil
}

// StringArray 字符串数组类型，用于 PostgreSQL text[]
type StringArray []string

// Scan 实现 sql.Scanner 接口
func (s *StringArray) Scan(value interface{}) error {
	if value == nil {
		*s = []string{}
		return nil
	}
	// PostgreSQL 返回的是字符串格式，需要解析
	return json.Unmarshal(value.([]byte), s)
}

// Value 实现 driver.Valuer 接口
func (s StringArray) Value() (driver.Value, error) {
	if len(s) == 0 {
		return "{}", nil
	}
	return json.Marshal(s)
}

// JSONBArray JSONB 数组类型
type JSONBArray []map[string]interface{}

// Scan 实现 sql.Scanner 接口
func (j *JSONBArray) Scan(value interface{}) error {
	if value == nil {
		*j = []map[string]interface{}{}
		return nil
	}
	return json.Unmarshal(value.([]byte), j)
}

// Value 实现 driver.Valuer 接口
func (j JSONBArray) Value() (driver.Value, error) {
	if len(j) == 0 {
		return "[]", nil
	}
	return json.Marshal(j)
}

// generateUUID 生成 UUID
func generateUUID() string {
	return uuid.New().String()
}

// RAGSource RAG 来源信息
type RAGSource struct {
	DocumentID string  `json:"document_id"`
	Title      string  `json:"title"`
	Content    string  `json:"content"`
	Score      float64 `json:"score"`
}

// RAGSources RAG 来源列表
type RAGSources []RAGSource

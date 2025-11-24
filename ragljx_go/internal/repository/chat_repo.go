package repository

import (
	"context"
	"ragljx/internal/model"

	"gorm.io/gorm"
)

type ChatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) *ChatRepository {
	return &ChatRepository{db: db}
}

// CreateSession 创建会话
func (r *ChatRepository) CreateSession(ctx context.Context, session *model.ChatSession) error {
	return r.db.WithContext(ctx).Create(session).Error
}

// GetSessionByID 根据 ID 获取会话
func (r *ChatRepository) GetSessionByID(ctx context.Context, id string) (*model.ChatSession, error) {
	var session model.ChatSession
	err := r.db.WithContext(ctx).Preload("User").First(&session, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// ListSessions 获取会话列表
func (r *ChatRepository) ListSessions(ctx context.Context, userID int, offset, limit int) ([]*model.ChatSession, int64, error) {
	var sessions []*model.ChatSession
	var total int64

	query := r.db.WithContext(ctx).Model(&model.ChatSession{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Offset(offset).Limit(limit).Order("updated_at DESC").Find(&sessions).Error
	return sessions, total, err
}

// UpdateSession 更新会话
func (r *ChatRepository) UpdateSession(ctx context.Context, session *model.ChatSession) error {
	return r.db.WithContext(ctx).Save(session).Error
}

// DeleteSession 删除会话
func (r *ChatRepository) DeleteSession(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&model.ChatSession{}, "id = ?", id).Error
}

// IncrementMessageCount 增加消息计数
func (r *ChatRepository) IncrementMessageCount(ctx context.Context, sessionID string) error {
	return r.db.WithContext(ctx).Model(&model.ChatSession{}).
		Where("id = ?", sessionID).
		UpdateColumn("message_count", gorm.Expr("message_count + 1")).Error
}

// CreateMessage 创建消息
func (r *ChatRepository) CreateMessage(ctx context.Context, message *model.ChatMessage) error {
	return r.db.WithContext(ctx).Create(message).Error
}

// GetMessagesBySession 获取会话消息列表
func (r *ChatRepository) GetMessagesBySession(ctx context.Context, sessionID string, offset, limit int) ([]*model.ChatMessage, int64, error) {
	var messages []*model.ChatMessage
	var total int64

	query := r.db.WithContext(ctx).Model(&model.ChatMessage{}).Where("session_id = ?", sessionID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Offset(offset).Limit(limit).Order("created_at ASC").Find(&messages).Error
	return messages, total, err
}

// DeleteMessagesBySession 删除会话的所有消息
func (r *ChatRepository) DeleteMessagesBySession(ctx context.Context, sessionID string) error {
	return r.db.WithContext(ctx).Where("session_id = ?", sessionID).Delete(&model.ChatMessage{}).Error
}


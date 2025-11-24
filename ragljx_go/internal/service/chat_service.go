package service

import (
	"context"
	"ragljx/internal/model"
	"ragljx/internal/pkg/errors"
	"ragljx/internal/repository"

	"gorm.io/gorm"
)

type ChatService struct {
	chatRepo *repository.ChatRepository
	kbRepo   *repository.KnowledgeBaseRepository
}

func NewChatService(db *gorm.DB) *ChatService {
	return &ChatService{
		chatRepo: repository.NewChatRepository(db),
		kbRepo:   repository.NewKnowledgeBaseRepository(db),
	}
}

// CreateSessionRequest 创建会话请求
type CreateSessionRequest struct {
	Title               string   `json:"title"`
	KnowledgeBaseIDs    []string `json:"knowledge_base_ids"`
	UseRAG              bool     `json:"use_rag"`
	TopK                int      `json:"top_k"`
	SimilarityThreshold float64  `json:"similarity_threshold"`
	SimilarityWeight    float64  `json:"similarity_weight"`
}

// UpdateSessionRequest 更新会话请求
type UpdateSessionRequest struct {
	Title               string   `json:"title"`
	KnowledgeBaseIDs    []string `json:"knowledge_base_ids"`
	UseRAG              bool     `json:"use_rag"`
	TopK                int      `json:"top_k"`
	SimilarityThreshold float64  `json:"similarity_threshold"`
	SimilarityWeight    float64  `json:"similarity_weight"`
}

// ChatRequest 聊天请求
type ChatRequest struct {
	SessionID string `json:"session_id" binding:"required"`
	Message   string `json:"message" binding:"required"`
	Stream    bool   `json:"stream"`
}

// ChatResponse 聊天响应
type ChatResponse struct {
	MessageID  string              `json:"message_id"`
	Content    string              `json:"content"`
	RAGSources []model.RAGSource   `json:"rag_sources,omitempty"`
	TokensUsed int                 `json:"tokens_used"`
}

// CreateSession 创建会话
func (s *ChatService) CreateSession(ctx context.Context, req *CreateSessionRequest, userID int) (*model.ChatSession, error) {
	// 验证知识库是否存在
	for _, kbID := range req.KnowledgeBaseIDs {
		_, err := s.kbRepo.GetByID(ctx, kbID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, errors.New(400, "knowledge base not found: "+kbID)
			}
			return nil, errors.Wrap(500, "failed to get knowledge base", err)
		}
	}

	// 创建会话
	session := &model.ChatSession{
		UserID:              userID,
		Title:               req.Title,
		KnowledgeBaseIDs:    model.StringArray(req.KnowledgeBaseIDs),
		UseRAG:              req.UseRAG,
		TopK:                req.TopK,
		SimilarityThreshold: req.SimilarityThreshold,
		SimilarityWeight:    req.SimilarityWeight,
		Status:              "active",
	}

	// 设置默认值
	if session.Title == "" {
		session.Title = "新对话"
	}
	if session.TopK == 0 {
		session.TopK = 3
	}
	if session.SimilarityThreshold == 0 {
		session.SimilarityThreshold = 0.7
	}
	if session.SimilarityWeight == 0 {
		session.SimilarityWeight = 1.5
	}

	if err := s.chatRepo.CreateSession(ctx, session); err != nil {
		return nil, errors.Wrap(500, "failed to create session", err)
	}

	return session, nil
}

// GetSessionByID 根据 ID 获取会话
func (s *ChatService) GetSessionByID(ctx context.Context, id string, userID int) (*model.ChatSession, error) {
	session, err := s.chatRepo.GetSessionByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrSessionNotFound
		}
		return nil, errors.Wrap(500, "failed to get session", err)
	}

	// 检查权限
	if session.UserID != userID {
		return nil, errors.ErrForbidden
	}

	return session, nil
}

// ListSessions 获取会话列表
func (s *ChatService) ListSessions(ctx context.Context, userID int, page, pageSize int) ([]*model.ChatSession, int64, error) {
	offset := (page - 1) * pageSize
	sessions, total, err := s.chatRepo.ListSessions(ctx, userID, offset, pageSize)
	if err != nil {
		return nil, 0, errors.Wrap(500, "failed to list sessions", err)
	}
	return sessions, total, nil
}

// UpdateSession 更新会话
func (s *ChatService) UpdateSession(ctx context.Context, id string, req *UpdateSessionRequest, userID int) (*model.ChatSession, error) {
	session, err := s.chatRepo.GetSessionByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrSessionNotFound
		}
		return nil, errors.Wrap(500, "failed to get session", err)
	}

	// 检查权限
	if session.UserID != userID {
		return nil, errors.ErrForbidden
	}

	// 更新字段
	if req.Title != "" {
		session.Title = req.Title
	}
	if len(req.KnowledgeBaseIDs) > 0 {
		session.KnowledgeBaseIDs = model.StringArray(req.KnowledgeBaseIDs)
	}
	session.UseRAG = req.UseRAG
	if req.TopK > 0 {
		session.TopK = req.TopK
	}
	if req.SimilarityThreshold > 0 {
		session.SimilarityThreshold = req.SimilarityThreshold
	}
	if req.SimilarityWeight > 0 {
		session.SimilarityWeight = req.SimilarityWeight
	}

	if err := s.chatRepo.UpdateSession(ctx, session); err != nil {
		return nil, errors.Wrap(500, "failed to update session", err)
	}

	return session, nil
}

// DeleteSession 删除会话
func (s *ChatService) DeleteSession(ctx context.Context, id string, userID int) error {
	session, err := s.chatRepo.GetSessionByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.ErrSessionNotFound
		}
		return errors.Wrap(500, "failed to get session", err)
	}

	// 检查权限
	if session.UserID != userID {
		return errors.ErrForbidden
	}

	// 删除会话消息
	if err := s.chatRepo.DeleteMessagesBySession(ctx, id); err != nil {
		return errors.Wrap(500, "failed to delete messages", err)
	}

	// 删除会话
	if err := s.chatRepo.DeleteSession(ctx, id); err != nil {
		return errors.Wrap(500, "failed to delete session", err)
	}

	return nil
}

// GetMessages 获取会话消息列表
func (s *ChatService) GetMessages(ctx context.Context, sessionID string, userID int, page, pageSize int) ([]*model.ChatMessage, int64, error) {
	// 检查会话权限
	session, err := s.chatRepo.GetSessionByID(ctx, sessionID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, 0, errors.ErrSessionNotFound
		}
		return nil, 0, errors.Wrap(500, "failed to get session", err)
	}

	if session.UserID != userID {
		return nil, 0, errors.ErrForbidden
	}

	offset := (page - 1) * pageSize
	messages, total, err := s.chatRepo.GetMessagesBySession(ctx, sessionID, offset, pageSize)
	if err != nil {
		return nil, 0, errors.Wrap(500, "failed to get messages", err)
	}

	return messages, total, nil
}

// SaveMessage 保存消息
func (s *ChatService) SaveMessage(ctx context.Context, sessionID, role, content string, ragSources []model.RAGSource, tokensUsed int) (*model.ChatMessage, error) {
	// 将 RAGSource 转换为 JSONBArray
	var ragSourcesJSON model.JSONBArray
	if len(ragSources) > 0 {
		ragSourcesJSON = make(model.JSONBArray, len(ragSources))
		for i, src := range ragSources {
			ragSourcesJSON[i] = map[string]interface{}{
				"document_id": src.DocumentID,
				"title":       src.Title,
				"content":     src.Content,
				"score":       src.Score,
			}
		}
	}

	message := &model.ChatMessage{
		SessionID:  sessionID,
		Role:       role,
		Content:    content,
		RAGSources: ragSourcesJSON,
		TokensUsed: tokensUsed,
	}

	if err := s.chatRepo.CreateMessage(ctx, message); err != nil {
		return nil, errors.Wrap(500, "failed to save message", err)
	}

	// 增加会话消息计数
	s.chatRepo.IncrementMessageCount(ctx, sessionID)

	return message, nil
}


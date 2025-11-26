package chat

import (
	"context"
	"encoding/json"
	"io"
	"ragljx/internal/middleware"
	"ragljx/internal/model"
	"ragljx/internal/pkg/response"
	"ragljx/internal/pkg/utils"
	"ragljx/internal/service"
	"ragljx/ioc"
	"ragljx/ioc/config/datasource"
	grpcConfig "ragljx/ioc/config/grpc"
	httpConfig "ragljx/ioc/config/http"
	pb "ragljx/proto/rag"

	"github.com/gin-gonic/gin"
)

func init() {
	ioc.Api().Registry(&ChatAPI{})
}

type ChatAPI struct {
	ioc.ObjectImpl
	chatService *service.ChatService
	grpcClient  pb.RAGServiceClient
}

func (ch *ChatAPI) Name() string {
	return "chat_api"
}

func (ch *ChatAPI) Init() error {
	db := datasource.Get()
	ch.chatService = service.NewChatService(db)

	// 获取 gRPC 连接
	grpcConn := grpcConfig.Get()
	if grpcConn != nil {
		ch.grpcClient = pb.NewRAGServiceClient(grpcConn)
	}

	// 注册路由
	engine := httpConfig.RootRouter()
	ch.Registry(engine)

	return nil
}

func (ch *ChatAPI) Registry(r gin.IRouter) {
	api := r.Group("/api/v1/chat", middleware.JWTAuth())
	{
		api.POST("/sessions", ch.CreateSession)
		api.GET("/sessions", ch.ListSessions)
		api.GET("/sessions/:id", ch.GetSession)
		api.PUT("/sessions/:id", ch.UpdateSession)
		api.DELETE("/sessions/:id", ch.DeleteSession)
		api.GET("/sessions/:id/messages", ch.GetMessages)
		api.POST("/sessions/:id/messages", ch.Chat)
		api.GET("/sessions/:id/messages/stream", ch.ChatStream)
	}
}

// CreateSession 创建会话
func (ch *ChatAPI) CreateSession(c *gin.Context) {
	var req service.CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID, _ := middleware.GetUserID(c)
	session, err := ch.chatService.CreateSession(c.Request.Context(), &req, userID)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, session)
}

// GetSession 获取会话
func (ch *ChatAPI) GetSession(c *gin.Context) {
	id := c.Param("id")
	userID, _ := middleware.GetUserID(c)

	session, err := ch.chatService.GetSessionByID(c.Request.Context(), id, userID)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, session)
}

// ListSessions 获取会话列表
func (ch *ChatAPI) ListSessions(c *gin.Context) {
	var pagination utils.Pagination
	if err := c.ShouldBindQuery(&pagination); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID, _ := middleware.GetUserID(c)
	sessions, total, err := ch.chatService.ListSessions(c.Request.Context(), userID, pagination.Page, pagination.GetLimit())
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Page(c, total, pagination.Page, pagination.GetLimit(), sessions)
}

// UpdateSession 更新会话
func (ch *ChatAPI) UpdateSession(c *gin.Context) {
	id := c.Param("id")

	var req service.UpdateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID, _ := middleware.GetUserID(c)
	session, err := ch.chatService.UpdateSession(c.Request.Context(), id, &req, userID)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, session)
}

// DeleteSession 删除会话
func (ch *ChatAPI) DeleteSession(c *gin.Context) {
	id := c.Param("id")
	userID, _ := middleware.GetUserID(c)

	if err := ch.chatService.DeleteSession(c.Request.Context(), id, userID); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetMessages 获取会话消息
func (ch *ChatAPI) GetMessages(c *gin.Context) {
	id := c.Param("id")

	var pagination utils.Pagination
	if err := c.ShouldBindQuery(&pagination); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID, _ := middleware.GetUserID(c)
	messages, total, err := ch.chatService.GetMessages(c.Request.Context(), id, userID, pagination.Page, pagination.GetLimit())
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Page(c, total, pagination.Page, pagination.GetLimit(), messages)
}

// Chat 同步聊天
func (ch *ChatAPI) Chat(c *gin.Context) {
	sessionID := c.Param("id")

	var req service.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 从 URL 参数设置 session_id
	req.SessionID = sessionID

	userID, _ := middleware.GetUserID(c)

	// 获取会话信息
	session, err := ch.chatService.GetSessionByID(c.Request.Context(), req.SessionID, userID)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	// 获取历史消息，用于构造上下文
	historyMessages, _, err := ch.chatService.GetMessages(c.Request.Context(), req.SessionID, userID, 1, 50)
	if err != nil {
		response.Error(c, 500, "failed to get chat history: "+err.Error())
		return
	}
	var history []*pb.ChatMessage
	for _, m := range historyMessages {
		history = append(history, &pb.ChatMessage{
			Role:    m.Role,
			Content: m.Content,
		})
	}

	// 保存用户消息
	ch.chatService.SaveMessage(c.Request.Context(), req.SessionID, "user", req.Message, nil, 0)

	// 检查 gRPC 客户端是否可用
	if ch.grpcClient == nil {
		response.Error(c, 500, "gRPC client is not initialized")
		return
	}

	// 获取知识库的 english_name 列表（用作 collection name）
	collectionNames, err := ch.chatService.GetKnowledgeBaseCollectionNames(c.Request.Context(), []string(session.KnowledgeBaseIDs))
	if err != nil {
		response.Error(c, 500, "failed to get knowledge base collection names: "+err.Error())
		return
	}

	// 调用 Python gRPC 服务
	grpcReq := &pb.ChatRequest{
		Query:               req.Message,
		UseRag:              session.UseRAG,
		KnowledgeBaseIds:    collectionNames, // 使用 english_name 而不是 UUID
		TopK:                int32(session.TopK),
		SimilarityThreshold: float32(session.SimilarityThreshold),
		SimilarityWeight:    float32(session.SimilarityWeight),
		History:             history,
	}

	grpcResp, err := ch.grpcClient.Chat(c.Request.Context(), grpcReq)
	if err != nil {
		response.Error(c, 500, "failed to call chat service: "+err.Error())
		return
	}

	// 转换 RAG 来源
	var ragSources []model.RAGSource
	for _, src := range grpcResp.Sources {
		ragSources = append(ragSources, model.RAGSource{
			DocumentID: src.DocumentId,
			Title:      src.Title,
			Score:      float64(src.Score),
			Content:    src.Snippet,
		})
	}

	// 保存助手消息
	ch.chatService.SaveMessage(c.Request.Context(), req.SessionID, "assistant", grpcResp.Content, ragSources, int(grpcResp.TokensUsed))

	// 返回响应
	resp := service.ChatResponse{
		Content:    grpcResp.Content,
		RAGSources: ragSources,
		TokensUsed: int(grpcResp.TokensUsed),
	}

	response.Success(c, resp)
}

// ChatStream 流式聊天
func (ch *ChatAPI) ChatStream(c *gin.Context) {
	sessionID := c.Param("id")

	// 从查询参数获取数据（因为是 GET 请求）
	question := c.Query("question")
	if question == "" {
		response.BadRequest(c, "question is required")
		return
	}

	var req service.ChatRequest
	req.SessionID = sessionID
	req.Message = question

	// EventSource 不支持自定义 headers，所以从 query 参数获取 token
	token := c.Query("token")
	if token == "" {
		// 如果 query 中没有 token，尝试从 header 获取
		token = c.GetHeader("Authorization")
		if token != "" && len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}
	}

	if token == "" {
		response.Unauthorized(c, "token is required")
		return
	}

	// 验证 token 并获取 userID
	userID, err := middleware.ValidateToken(token)
	if err != nil {
		response.Unauthorized(c, "invalid token")
		return
	}

	// 获取会话信息
	session, err := ch.chatService.GetSessionByID(c.Request.Context(), req.SessionID, userID)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	// 获取历史消息，用于构造上下文
	historyMessages, _, err := ch.chatService.GetMessages(c.Request.Context(), req.SessionID, userID, 1, 50)
	if err != nil {
		response.Error(c, 500, "failed to get chat history: "+err.Error())
		return
	}
	var history []*pb.ChatMessage
	for _, m := range historyMessages {
		history = append(history, &pb.ChatMessage{
			Role:    m.Role,
			Content: m.Content,
		})
	}

	// 检查 gRPC 客户端是否可用
	if ch.grpcClient == nil {
		response.Error(c, 500, "gRPC client is not initialized")
		return
	}

	// 保存用户消息
	ch.chatService.SaveMessage(c.Request.Context(), req.SessionID, "user", req.Message, nil, 0)

	// 获取知识库的 english_name 列表（用作 collection name）
	collectionNames, err := ch.chatService.GetKnowledgeBaseCollectionNames(c.Request.Context(), []string(session.KnowledgeBaseIDs))
	if err != nil {
		response.Error(c, 500, "failed to get knowledge base collection names: "+err.Error())
		return
	}

	// 调用 Python gRPC 服务（流式）
	grpcReq := &pb.ChatRequest{
		Query:               req.Message,
		UseRag:              session.UseRAG,
		KnowledgeBaseIds:    collectionNames, // 使用 english_name 而不是 UUID
		TopK:                int32(session.TopK),
		SimilarityThreshold: float32(session.SimilarityThreshold),
		SimilarityWeight:    float32(session.SimilarityWeight),
		History:             history,
	}

	stream, err := ch.grpcClient.ChatStream(c.Request.Context(), grpcReq)
	if err != nil {
		response.Error(c, 500, "failed to call chat stream service: "+err.Error())
		return
	}

	// 设置 SSE 响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	// 流式返回
	fullContent := ""
	var ragSources []model.RAGSource
	var tokensUsed int32
	sourcesSent := false

	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			errorData := map[string]interface{}{
				"type":  "error",
				"error": err.Error(),
			}
			errorJSON, _ := json.Marshal(errorData)
			c.SSEvent("message", string(errorJSON))
			c.Writer.Flush()
			break
		}

		fullContent += chunk.Delta
		tokensUsed = chunk.TokensUsed

		// 处理 start 类型（包含 sources）
		if chunk.Type == "start" && len(chunk.Sources) > 0 && !sourcesSent {
			for _, src := range chunk.Sources {
				ragSources = append(ragSources, model.RAGSource{
					DocumentID: src.DocumentId,
					Title:      src.Title,
					Score:      float64(src.Score),
					Content:    src.Snippet,
				})
			}

			sourcesData := map[string]interface{}{
				"type":    "sources",
				"sources": ragSources,
			}
			sourcesJSON, _ := json.Marshal(sourcesData)
			c.SSEvent("message", string(sourcesJSON))
			c.Writer.Flush()
			sourcesSent = true
		}

		// 发送内容增量
		if chunk.Type == "content" && chunk.Delta != "" {
			contentData := map[string]interface{}{
				"type":    "content",
				"content": chunk.Delta,
			}
			contentJSON, _ := json.Marshal(contentData)
			c.SSEvent("message", string(contentJSON))
			c.Writer.Flush()
		}
	}

	// 保存助手消息
	ch.chatService.SaveMessage(context.Background(), req.SessionID, "assistant", fullContent, ragSources, int(tokensUsed))

	// 发送完成信号
	doneData := map[string]interface{}{
		"type": "done",
	}
	doneJSON, _ := json.Marshal(doneData)
	c.SSEvent("message", string(doneJSON))
	c.Writer.Flush()
}

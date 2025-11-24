"""RAG gRPC 服务实现"""
from loguru import logger
import grpc

from app.proto import rag_service_pb2
from app.proto import rag_service_pb2_grpc
from app.services import FileParserService, VectorService, ChatService


class RAGServicer(rag_service_pb2_grpc.RAGServiceServicer):
    """RAG 服务实现"""
    
    def __init__(self):
        self.file_parser = FileParserService()
        self.vector_service = VectorService()
        self.chat_service = ChatService(self.vector_service)
        logger.info("RAGServicer initialized")
    
    def ParseDocument(self, request, context):
        """解析文档"""
        try:
            logger.info(f"Parsing document: {request.document_id}")
            
            # 解析文件
            success, content, error_msg = self.file_parser.parse(
                file_content=request.file_content,
                mime_type=request.mime_type,
                filename=request.object_key
            )
            
            return rag_service_pb2.ParseDocumentResponse(
                success=success,
                content=content,
                error_message=error_msg
            )
        
        except Exception as e:
            logger.error(f"Error parsing document: {e}")
            return rag_service_pb2.ParseDocumentResponse(
                success=False,
                content="",
                error_message=str(e)
            )
    
    def VectorizeDocument(self, request, context):
        """向量化文档"""
        try:
            logger.info(f"Vectorizing document: {request.document_id} to collection: {request.collection_name}")
            
            # 向量化并插入
            success = self.vector_service.upsert_document(
                collection_name=request.collection_name,
                document_id=request.document_id,
                content=request.content,
                title=request.title,
                metadata={
                    "knowledge_base_id": request.knowledge_base_id,
                    "object_key": request.object_key
                }
            )
            
            return rag_service_pb2.VectorizeDocumentResponse(
                success=success,
                error_message="" if success else "Vectorization failed"
            )
        
        except Exception as e:
            logger.error(f"Error vectorizing document: {e}")
            return rag_service_pb2.VectorizeDocumentResponse(
                success=False,
                error_message=str(e)
            )
    
    def DeleteDocumentVectors(self, request, context):
        """删除文档向量"""
        try:
            logger.info(f"Deleting vectors for documents: {request.document_ids} from collection: {request.collection_name}")
            
            deleted_count = 0
            for document_id in request.document_ids:
                if self.vector_service.delete_document(request.collection_name, document_id):
                    deleted_count += 1
            
            return rag_service_pb2.DeleteDocumentVectorsResponse(
                success=True,
                deleted_count=deleted_count,
                error_message=""
            )
        
        except Exception as e:
            logger.error(f"Error deleting document vectors: {e}")
            return rag_service_pb2.DeleteDocumentVectorsResponse(
                success=False,
                deleted_count=0,
                error_message=str(e)
            )
    
    def Chat(self, request, context):
        """RAG 对话（非流式）"""
        try:
            logger.info(f"Chat request: {request.query[:50]}... use_rag={request.use_rag}")
            
            # 构建历史对话
            history = []
            for msg in request.history:
                history.append({
                    "role": msg.role,
                    "content": msg.content
                })
            
            # 构建知识库集合名称列表
            kb_collection_names = list(request.knowledge_base_ids) if request.use_rag else []
            
            # 调用对话服务
            content, rag_sources, tokens_used = self.chat_service.chat(
                query=request.query,
                use_rag=request.use_rag,
                kb_collection_names=kb_collection_names,
                top_k=request.top_k if request.top_k > 0 else 5,
                similarity_threshold=request.similarity_threshold if request.similarity_threshold > 0 else 0.7,
                history=history
            )
            
            # 构建响应
            sources = []
            for src in rag_sources:
                sources.append(rag_service_pb2.RAGSource(
                    document_id=src['document_id'],
                    title=src['title'],
                    score=src['score'],
                    snippet=src['snippet']
                ))
            
            return rag_service_pb2.ChatResponse(
                content=content,
                sources=sources,
                tokens_used=tokens_used
            )
        
        except Exception as e:
            logger.error(f"Error in chat: {e}")
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            return rag_service_pb2.ChatResponse(
                content=f"Error: {str(e)}",
                sources=[],
                tokens_used=0
            )
    
    def ChatStream(self, request, context):
        """RAG 对话（流式）"""
        try:
            logger.info(f"Chat stream request: {request.query[:50]}... use_rag={request.use_rag}")
            
            # 构建历史对话
            history = []
            for msg in request.history:
                history.append({
                    "role": msg.role,
                    "content": msg.content
                })
            
            # 构建知识库集合名称列表
            kb_collection_names = list(request.knowledge_base_ids) if request.use_rag else []
            
            # 调用流式对话服务
            for chunk in self.chat_service.chat_stream(
                query=request.query,
                use_rag=request.use_rag,
                kb_collection_names=kb_collection_names,
                top_k=request.top_k if request.top_k > 0 else 5,
                similarity_threshold=request.similarity_threshold if request.similarity_threshold > 0 else 0.7,
                history=history
            ):
                # 构建流式响应
                sources = []
                if chunk.get('sources'):
                    for src in chunk['sources']:
                        sources.append(rag_service_pb2.RAGSource(
                            document_id=src['document_id'],
                            title=src['title'],
                            score=src['score'],
                            snippet=src['snippet']
                        ))
                
                yield rag_service_pb2.ChatStreamResponse(
                    type=chunk.get('type', 'content'),
                    delta=chunk.get('delta', ''),
                    sources=sources,
                    tokens_used=0  # 流式模式下无法获取准确的 token 数
                )
        
        except Exception as e:
            logger.error(f"Error in chat stream: {e}")
            yield rag_service_pb2.ChatStreamResponse(
                type="error",
                delta=str(e),
                sources=[],
                tokens_used=0
            )
    
    def RetrieveDocuments(self, request, context):
        """检索相关文档"""
        try:
            logger.info(f"Retrieve documents: {request.query[:50]}...")
            
            # 检索文档
            retrieved_docs = self.chat_service.retrieve_documents(
                query=request.query,
                kb_collection_names=list(request.knowledge_base_ids),
                top_k=request.top_k if request.top_k > 0 else 5,
                similarity_threshold=request.similarity_threshold if request.similarity_threshold > 0 else 0.7
            )
            
            # 构建响应
            documents = []
            for doc in retrieved_docs:
                documents.append(rag_service_pb2.RAGSource(
                    document_id=doc['document_id'],
                    title=doc.get('title', ''),
                    score=doc.get('score', 0.0),
                    snippet=doc.get('content', '')[:500]  # 返回前500字符
                ))
            
            return rag_service_pb2.RetrieveDocumentsResponse(
                documents=documents
            )
        
        except Exception as e:
            logger.error(f"Error retrieving documents: {e}")
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            return rag_service_pb2.RetrieveDocumentsResponse(
                documents=[]
            )


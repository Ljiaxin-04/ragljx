"""向量服务"""
from typing import List, Dict, Any, Optional
from loguru import logger
from qdrant_client import QdrantClient
from qdrant_client.models import Distance, VectorParams, PointStruct, Filter, FieldCondition, MatchValue
from llama_index.core import Document, VectorStoreIndex, Settings
from llama_index.core.node_parser import SentenceSplitter
from llama_index.embeddings.openai import OpenAIEmbedding
from llama_index.vector_stores.qdrant import QdrantVectorStore
from tenacity import retry, stop_after_attempt, wait_exponential

from app.config import config


class VectorService:
    """向量服务 - 负责文档向量化和检索"""
    
    DEFAULT_VECTOR_NAME = "text-dense"
    
    def __init__(self):
        # 初始化 Qdrant 客户端
        self.qdrant_client = QdrantClient(
            host=config.qdrant_host,
            port=config.qdrant_port
        )

        # 初始化嵌入模型
        self.embedding = OpenAIEmbedding(
            api_key=config.embedding_api_key,
            api_base=config.embedding_api_base,
            model=config.embedding_model,
            dimensions=config.embedding_dimensions
        )

        # 配置 LlamaIndex Settings
        Settings.embed_model = self.embedding
        Settings.chunk_size = config.rag_chunk_size
        Settings.chunk_overlap = config.rag_chunk_overlap

        # 文本分割器
        self.text_splitter = SentenceSplitter(
            chunk_size=config.rag_chunk_size,
            chunk_overlap=config.rag_chunk_overlap
        )

        logger.info(f"VectorService initialized with Qdrant at {config.qdrant_host}:{config.qdrant_port}")
        logger.info(f"Using embedding model: {config.embedding_model}")
    
    def build_collection_name(self, kb_english_name: str) -> str:
        """
        构建集合名称
        
        Args:
            kb_english_name: 知识库英文名
        
        Returns:
            集合名称，格式: {prefix}_{kb_english_name}_{embed_model}
        """
        # 简化 embedding 模型名称
        model_name = config.openai_embedding_model.replace('text-embedding-', '').replace('-', '_')
        return f"{config.qdrant_collection_prefix}_{kb_english_name}_{model_name}"
    
    @retry(stop=stop_after_attempt(3), wait=wait_exponential(multiplier=1, min=2, max=10))
    def create_collection(self, collection_name: str) -> bool:
        """
        创建 Qdrant 集合
        
        Args:
            collection_name: 集合名称
        
        Returns:
            是否创建成功
        """
        try:
            # 检查集合是否已存在
            collections = self.qdrant_client.get_collections().collections
            if any(c.name == collection_name for c in collections):
                logger.info(f"Collection {collection_name} already exists")
                return True
            
            # 创建集合
            self.qdrant_client.create_collection(
                collection_name=collection_name,
                vectors_config=VectorParams(
                    size=config.openai_embedding_dimensions,
                    distance=Distance.COSINE
                )
            )
            
            logger.info(f"Created collection: {collection_name}")
            return True
        
        except Exception as e:
            logger.error(f"Error creating collection {collection_name}: {e}")
            raise
    
    def delete_collection(self, collection_name: str) -> bool:
        """
        删除集合
        
        Args:
            collection_name: 集合名称
        
        Returns:
            是否删除成功
        """
        try:
            self.qdrant_client.delete_collection(collection_name)
            logger.info(f"Deleted collection: {collection_name}")
            return True
        except Exception as e:
            logger.error(f"Error deleting collection {collection_name}: {e}")
            return False
    
    def upsert_document(
        self,
        collection_name: str,
        document_id: str,
        content: str,
        title: str = "",
        metadata: Optional[Dict[str, Any]] = None
    ) -> bool:
        """
        向量化并插入文档
        
        Args:
            collection_name: 集合名称
            document_id: 文档 ID
            content: 文档内容
            title: 文档标题
            metadata: 额外的元数据
        
        Returns:
            是否成功
        """
        try:
            # 确保集合存在
            self.create_collection(collection_name)
            
            # 创建 LlamaIndex Document
            doc_metadata = metadata or {}
            doc_metadata.update({
                "document_id": document_id,
                "title": title
            })
            
            document = Document(
                text=content,
                metadata=doc_metadata
            )
            
            # 分割文档为节点
            nodes = self.text_splitter.get_nodes_from_documents([document])
            
            # 为每个节点生成向量并插入
            points = []
            for i, node in enumerate(nodes):
                # 生成向量
                embedding_vector = self.embedding.get_text_embedding(node.get_content())
                
                # 构建点
                point_id = f"{document_id}_{i}"
                point = PointStruct(
                    id=point_id,
                    vector=embedding_vector,
                    payload={
                        "document_id": document_id,
                        "title": title,
                        "content": node.get_content(),
                        "chunk_index": i,
                        **doc_metadata
                    }
                )
                points.append(point)
            
            # 批量插入
            self.qdrant_client.upsert(
                collection_name=collection_name,
                points=points
            )
            
            logger.info(f"Upserted document {document_id} with {len(points)} chunks to {collection_name}")
            return True
        
        except Exception as e:
            logger.error(f"Error upserting document {document_id}: {e}")
            raise
    
    def delete_document(self, collection_name: str, document_id: str) -> bool:
        """
        删除文档的所有向量
        
        Args:
            collection_name: 集合名称
            document_id: 文档 ID
        
        Returns:
            是否成功
        """
        try:
            # 使用过滤器删除所有匹配的点
            self.qdrant_client.delete(
                collection_name=collection_name,
                points_selector=Filter(
                    must=[
                        FieldCondition(
                            key="document_id",
                            match=MatchValue(value=document_id)
                        )
                    ]
                )
            )
            
            logger.info(f"Deleted document {document_id} from {collection_name}")
            return True
        
        except Exception as e:
            logger.error(f"Error deleting document {document_id}: {e}")
            return False
    
    def query(
        self,
        collection_name: str,
        query_text: str,
        top_k: int = 5,
        score_threshold: float = 0.7,
        filter_dict: Optional[Dict[str, Any]] = None
    ) -> List[Dict[str, Any]]:
        """
        查询相似文档

        Args:
            collection_name: 集合名称
            query_text: 查询文本
            top_k: 返回结果数量
            score_threshold: 相似度阈值
            filter_dict: 过滤条件

        Returns:
            相似文档列表
        """
        try:
            # 检查集合是否存在
            collections = self.qdrant_client.get_collections().collections
            if not any(c.name == collection_name for c in collections):
                logger.warning(f"Collection {collection_name} does not exist, skipping query")
                return []

            # 生成查询向量
            query_vector = self.embedding.get_text_embedding(query_text)

            # 构建过滤器
            query_filter = None
            if filter_dict:
                conditions = []
                for key, value in filter_dict.items():
                    conditions.append(
                        FieldCondition(key=key, match=MatchValue(value=value))
                    )
                query_filter = Filter(must=conditions)

            # 执行查询
            search_result = self.qdrant_client.search(
                collection_name=collection_name,
                query_vector=query_vector,
                limit=top_k,
                query_filter=query_filter,
                score_threshold=score_threshold
            )

            # 转换结果
            results = []
            for hit in search_result:
                results.append({
                    "document_id": hit.payload.get("document_id"),
                    "title": hit.payload.get("title", ""),
                    "content": hit.payload.get("content", ""),
                    "score": hit.score,
                    "chunk_index": hit.payload.get("chunk_index", 0)
                })

            logger.info(f"Query collection {collection_name} returned {len(results)} results")
            return results

        except Exception as e:
            logger.error(f"Error querying collection {collection_name}: {e}", exc_info=True)
            return []


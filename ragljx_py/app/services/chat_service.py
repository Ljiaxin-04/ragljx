"""对话服务"""
from typing import List, Dict, Any, Iterator, Optional
from loguru import logger
from openai import OpenAI

from app.config import config
from app.services.vector_service import VectorService


class ChatService:
    """对话服务 - 负责 RAG 对话"""
    
    def __init__(self, vector_service: VectorService):
        self.vector_service = vector_service

        # 初始化对话模型客户端
        self.client = OpenAI(
            api_key=config.chat_api_key,
            base_url=config.chat_api_base
        )

        logger.info(f"ChatService initialized with model: {config.chat_model}")
    
    def retrieve_documents(
        self,
        query: str,
        kb_collection_names: List[str],
        top_k: int = 5,
        similarity_threshold: float = 0.7
    ) -> List[Dict[str, Any]]:
        """从多个知识库检索相关文档。

        Args:
            query: 查询文本
            kb_collection_names: 知识库集合名称列表
            top_k: 每个知识库返回的结果数
            similarity_threshold: 相似度阈值（最小相似度）。
                如果在该阈值下没有检索到结果，会自动放宽阈值再尝试一次。

        Returns:
            检索到的文档列表
        """
        all_results = []

        for collection_name in kb_collection_names:
            try:
                # 第一次使用会话配置的相似度阈值检索
                results = self.vector_service.query(
                    collection_name=collection_name,
                    query_text=query,
                    top_k=top_k,
                    score_threshold=similarity_threshold,
                )

                # 如果没有结果且设置了阈值，则再尝试一次放宽阈值（不设置阈值，只取前 top_k）
                if not results and similarity_threshold > 0:
                    logger.info(
                        f"No results from {collection_name} with threshold "
                        f"{similarity_threshold}, retrying without threshold"
                    )
                    results = self.vector_service.query(
                        collection_name=collection_name,
                        query_text=query,
                        top_k=top_k,
                        score_threshold=0.0,
                    )

                all_results.extend(results)
            except Exception as e:
                logger.error(f"Error retrieving from {collection_name}: {e}")
                continue

        # 按相似度排序并去重
        all_results.sort(key=lambda x: x['score'], reverse=True)

        # 去重（基于 document_id 和 chunk_index）
        seen = set()
        unique_results = []
        for result in all_results:
            key = (result['document_id'], result.get('chunk_index', 0))
            if key not in seen:
                seen.add(key)
                unique_results.append(result)

        # 限制总数
        return unique_results[:top_k * 2]  # 返回最多 top_k * 2 个结果

    def build_rag_context(self, retrieved_docs: List[Dict[str, Any]]) -> str:
        """
        构建 RAG 上下文
        
        Args:
            retrieved_docs: 检索到的文档列表
        
        Returns:
            格式化的上下文字符串
        """
        if not retrieved_docs:
            return ""
        
        context_parts = ["以下是相关的参考资料：\n"]
        
        for i, doc in enumerate(retrieved_docs, 1):
            title = doc.get('title', '未命名文档')
            content = doc.get('content', '')
            score = doc.get('score', 0.0)
            
            context_parts.append(f"[{i}] {title} (相似度: {score:.2f})")
            context_parts.append(content)
            context_parts.append("")  # 空行分隔
        
        return "\n".join(context_parts)
    
    def chat(
        self,
        query: str,
        use_rag: bool = False,
        kb_collection_names: Optional[List[str]] = None,
        top_k: int = 5,
        similarity_threshold: float = 0.7,
        history: Optional[List[Dict[str, str]]] = None,
        temperature: Optional[float] = None,
        max_tokens: Optional[int] = None
    ) -> tuple[str, List[Dict[str, Any]], int]:
        """
        对话（非流式）
        
        Args:
            query: 用户查询
            use_rag: 是否使用 RAG
            kb_collection_names: 知识库集合名称列表
            top_k: 检索文档数量
            similarity_threshold: 相似度阈值
            history: 历史对话
            temperature: 温度参数
            max_tokens: 最大 token 数
        
        Returns:
            (回复内容, RAG 来源列表, 使用的 token 数)
        """
        try:
            # 构建消息列表
            messages = []
            
            # 系统提示
            system_prompt = "你是一个有帮助的AI助手。"
            
            # 如果使用 RAG，检索相关文档
            rag_sources = []
            if use_rag and kb_collection_names:
                retrieved_docs = self.retrieve_documents(
                    query=query,
                    kb_collection_names=kb_collection_names,
                    top_k=top_k,
                    similarity_threshold=similarity_threshold
                )
                
                if retrieved_docs:
                    rag_context = self.build_rag_context(retrieved_docs)
                    system_prompt += f"\n\n{rag_context}\n\n请基于以上参考资料回答用户的问题。如果参考资料中没有相关信息，请明确告知用户。"
                    
                    # 保存 RAG 来源
                    rag_sources = [
                        {
                            "document_id": doc['document_id'],
                            "title": doc.get('title', ''),
                            "score": doc.get('score', 0.0),
                            "snippet": doc.get('content', '')[:200]  # 只保存前200字符作为摘要
                        }
                        for doc in retrieved_docs[:5]  # 只返回前5个来源
                    ]
            
            messages.append({"role": "system", "content": system_prompt})
            
            # 添加历史对话
            if history:
                for msg in history:
                    messages.append({
                        "role": msg.get('role', 'user'),
                        "content": msg.get('content', '')
                    })
            
            # 添加当前查询
            messages.append({"role": "user", "content": query})
            
            # 调用对话模型 API
            response = self.client.chat.completions.create(
                model=config.chat_model,
                messages=messages,
                temperature=temperature or config.chat_temperature,
                max_tokens=max_tokens or config.chat_max_tokens
            )
            
            # 提取回复
            content = response.choices[0].message.content
            tokens_used = response.usage.total_tokens if response.usage else 0
            
            # 清理输出
            content = self._sanitize_output(content)
            
            return content, rag_sources, tokens_used
        
        except Exception as e:
            logger.error(f"Error in chat: {e}")
            raise
    
    def chat_stream(
        self,
        query: str,
        use_rag: bool = False,
        kb_collection_names: Optional[List[str]] = None,
        top_k: int = 5,
        similarity_threshold: float = 0.7,
        history: Optional[List[Dict[str, str]]] = None,
        temperature: Optional[float] = None,
        max_tokens: Optional[int] = None
    ) -> Iterator[Dict[str, Any]]:
        """
        对话（流式）
        
        Args:
            query: 用户查询
            use_rag: 是否使用 RAG
            kb_collection_names: 知识库集合名称列表
            top_k: 检索文档数量
            similarity_threshold: 相似度阈值
            history: 历史对话
            temperature: 温度参数
            max_tokens: 最大 token 数
        
        Yields:
            流式响应块
        """
        try:
            # 构建消息列表（与非流式相同）
            messages = []
            system_prompt = "你是一个有帮助的AI助手。"
            
            rag_sources = []
            if use_rag and kb_collection_names:
                retrieved_docs = self.retrieve_documents(
                    query=query,
                    kb_collection_names=kb_collection_names,
                    top_k=top_k,
                    similarity_threshold=similarity_threshold
                )
                
                if retrieved_docs:
                    rag_context = self.build_rag_context(retrieved_docs)
                    system_prompt += f"\n\n{rag_context}\n\n请基于以上参考资料回答用户的问题。如果参考资料中没有相关信息，请明确告知用户。"
                    
                    rag_sources = [
                        {
                            "document_id": doc['document_id'],
                            "title": doc.get('title', ''),
                            "score": doc.get('score', 0.0),
                            "snippet": doc.get('content', '')[:200]
                        }
                        for doc in retrieved_docs[:5]
                    ]
            
            messages.append({"role": "system", "content": system_prompt})
            
            if history:
                for msg in history:
                    messages.append({
                        "role": msg.get('role', 'user'),
                        "content": msg.get('content', '')
                    })
            
            messages.append({"role": "user", "content": query})
            
            # 流式调用对话模型 API
            stream = self.client.chat.completions.create(
                model=config.chat_model,
                messages=messages,
                temperature=temperature or config.chat_temperature,
                max_tokens=max_tokens or config.chat_max_tokens,
                stream=True
            )
            
            # 发送开始信号
            yield {"type": "start", "delta": "", "sources": rag_sources}
            
            # 流式返回内容
            for chunk in stream:
                if chunk.choices[0].delta.content:
                    yield {
                        "type": "content",
                        "delta": chunk.choices[0].delta.content,
                        "sources": []
                    }
            
            # 发送结束信号
            yield {"type": "done", "delta": "", "sources": []}
        
        except Exception as e:
            logger.error(f"Error in chat_stream: {e}")
            yield {"type": "error", "delta": str(e), "sources": []}
    
    def _sanitize_output(self, text: str) -> str:
        """清理 LLM 输出"""
        # 移除可能的思考标签
        text = text.replace("<think>", "").replace("</think>", "")
        return text.strip()


"""配置管理模块"""
import os
import yaml
from pathlib import Path
from typing import Dict, Any


class Config:
    """配置类"""
    
    def __init__(self, config_path: str = "config.yaml"):
        self.config_path = config_path
        self._config: Dict[str, Any] = {}
        self.load()
    
    def load(self):
        """加载配置文件"""
        if os.path.exists(self.config_path):
            with open(self.config_path, 'r', encoding='utf-8') as f:
                self._config = yaml.safe_load(f) or {}
        
        # 从环境变量覆盖配置
        self._load_from_env()
    
    def _load_from_env(self):
        """从环境变量加载配置"""
        # gRPC
        if os.getenv('RAGLJX_GRPC_HOST'):
            self._config.setdefault('grpc', {})['host'] = os.getenv('RAGLJX_GRPC_HOST')
        if os.getenv('RAGLJX_GRPC_PORT'):
            self._config.setdefault('grpc', {})['port'] = int(os.getenv('RAGLJX_GRPC_PORT'))

        # Database
        if os.getenv('RAGLJX_DB_HOST'):
            self._config.setdefault('database', {})['host'] = os.getenv('RAGLJX_DB_HOST')
        if os.getenv('RAGLJX_DB_PORT'):
            self._config.setdefault('database', {})['port'] = int(os.getenv('RAGLJX_DB_PORT'))
        if os.getenv('RAGLJX_DB_DATABASE'):
            self._config.setdefault('database', {})['database'] = os.getenv('RAGLJX_DB_DATABASE')
        if os.getenv('RAGLJX_DB_USERNAME'):
            self._config.setdefault('database', {})['username'] = os.getenv('RAGLJX_DB_USERNAME')
        if os.getenv('RAGLJX_DB_PASSWORD'):
            self._config.setdefault('database', {})['password'] = os.getenv('RAGLJX_DB_PASSWORD')

        # Qdrant
        if os.getenv('RAGLJX_QDRANT_HOST'):
            self._config.setdefault('qdrant', {})['host'] = os.getenv('RAGLJX_QDRANT_HOST')
        if os.getenv('RAGLJX_QDRANT_PORT'):
            self._config.setdefault('qdrant', {})['port'] = int(os.getenv('RAGLJX_QDRANT_PORT'))

        # 嵌入模型配置
        if os.getenv('EMBEDDING_API_KEY'):
            self._config.setdefault('embedding', {})['api_key'] = os.getenv('EMBEDDING_API_KEY')
        if os.getenv('EMBEDDING_API_BASE'):
            self._config.setdefault('embedding', {})['api_base'] = os.getenv('EMBEDDING_API_BASE')
        if os.getenv('EMBEDDING_MODEL'):
            self._config.setdefault('embedding', {})['model'] = os.getenv('EMBEDDING_MODEL')

        # 对话模型配置
        if os.getenv('CHAT_API_KEY'):
            self._config.setdefault('chat', {})['api_key'] = os.getenv('CHAT_API_KEY')
        if os.getenv('CHAT_API_BASE'):
            self._config.setdefault('chat', {})['api_base'] = os.getenv('CHAT_API_BASE')
        if os.getenv('CHAT_MODEL'):
            self._config.setdefault('chat', {})['model'] = os.getenv('CHAT_MODEL')
        if os.getenv('CHAT_TEMPERATURE'):
            self._config.setdefault('chat', {})['temperature'] = float(os.getenv('CHAT_TEMPERATURE'))
        if os.getenv('CHAT_MAX_TOKENS'):
            self._config.setdefault('chat', {})['max_tokens'] = int(os.getenv('CHAT_MAX_TOKENS'))

        # MinIO
        if os.getenv('RAGLJX_MINIO_ENDPOINT'):
            self._config.setdefault('minio', {})['endpoint'] = os.getenv('RAGLJX_MINIO_ENDPOINT')
        if os.getenv('RAGLJX_MINIO_ACCESS_KEY'):
            self._config.setdefault('minio', {})['access_key'] = os.getenv('RAGLJX_MINIO_ACCESS_KEY')
        if os.getenv('RAGLJX_MINIO_SECRET_KEY'):
            self._config.setdefault('minio', {})['secret_key'] = os.getenv('RAGLJX_MINIO_SECRET_KEY')
    
    def get(self, key: str, default: Any = None) -> Any:
        """获取配置项"""
        keys = key.split('.')
        value = self._config
        for k in keys:
            if isinstance(value, dict):
                value = value.get(k)
            else:
                return default
            if value is None:
                return default
        return value
    
    @property
    def grpc_host(self) -> str:
        return self.get('grpc.host', '0.0.0.0')
    
    @property
    def grpc_port(self) -> int:
        return self.get('grpc.port', 50051)
    
    @property
    def db_url(self) -> str:
        """获取数据库连接 URL"""
        db_config = self._config.get('database', {})
        return f"postgresql://{db_config.get('username', 'ragljx')}:{db_config.get('password', 'ragljx_password')}@{db_config.get('host', 'localhost')}:{db_config.get('port', 5432)}/{db_config.get('database', 'ragljx')}"
    
    @property
    def qdrant_host(self) -> str:
        return self.get('qdrant.host', 'localhost')
    
    @property
    def qdrant_port(self) -> int:
        return self.get('qdrant.port', 6333)
    
    @property
    def qdrant_collection_prefix(self) -> str:
        return self.get('qdrant.collection_prefix', 'rag_collection')
    
    # 嵌入模型配置
    @property
    def embedding_api_key(self) -> str:
        return self.get('embedding.api_key', '')

    @property
    def embedding_api_base(self) -> str:
        return self.get('embedding.api_base', 'https://api.openai.com/v1')

    @property
    def embedding_model(self) -> str:
        return self.get('embedding.model', 'text-embedding-3-small')

    @property
    def embedding_dimensions(self) -> int:
        return self.get('embedding.dimensions', 1536)

    # 对话模型配置
    @property
    def chat_api_key(self) -> str:
        return self.get('chat.api_key', '')

    @property
    def chat_api_base(self) -> str:
        return self.get('chat.api_base', 'https://api.openai.com/v1')

    @property
    def chat_model(self) -> str:
        return self.get('chat.model', 'gpt-4')

    @property
    def chat_temperature(self) -> float:
        return self.get('chat.temperature', 0.7)

    @property
    def chat_max_tokens(self) -> int:
        return self.get('chat.max_tokens', 2000)

    # 兼容旧的属性名（保持向后兼容）
    @property
    def openai_api_key(self) -> str:
        """兼容旧代码"""
        return self.embedding_api_key or self.chat_api_key

    @property
    def openai_api_base(self) -> str:
        """兼容旧代码"""
        return self.chat_api_base

    @property
    def openai_embedding_model(self) -> str:
        """兼容旧代码"""
        return self.embedding_model

    @property
    def openai_embedding_dimensions(self) -> int:
        """兼容旧代码"""
        return self.embedding_dimensions

    @property
    def openai_chat_model(self) -> str:
        """兼容旧代码"""
        return self.chat_model

    @property
    def openai_temperature(self) -> float:
        """兼容旧代码"""
        return self.chat_temperature

    @property
    def openai_max_tokens(self) -> int:
        """兼容旧代码"""
        return self.chat_max_tokens
    
    @property
    def minio_endpoint(self) -> str:
        return self.get('minio.endpoint', 'localhost:9000')
    
    @property
    def minio_access_key(self) -> str:
        return self.get('minio.access_key', 'minioadmin')
    
    @property
    def minio_secret_key(self) -> str:
        return self.get('minio.secret_key', 'minioadmin')
    
    @property
    def minio_secure(self) -> bool:
        return self.get('minio.secure', False)
    
    @property
    def minio_bucket(self) -> str:
        return self.get('minio.bucket', 'ragljx')
    
    @property
    def rag_chunk_size(self) -> int:
        return self.get('rag.chunk_size', 512)
    
    @property
    def rag_chunk_overlap(self) -> int:
        return self.get('rag.chunk_overlap', 50)
    
    @property
    def rag_top_k(self) -> int:
        return self.get('rag.top_k', 5)
    
    @property
    def rag_similarity_threshold(self) -> float:
        return self.get('rag.similarity_threshold', 0.7)
    
    @property
    def rag_similarity_weight(self) -> float:
        return self.get('rag.similarity_weight', 0.5)


# 全局配置实例
config = Config()


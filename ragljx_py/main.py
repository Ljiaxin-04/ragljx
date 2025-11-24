"""RAG Python 服务主程序"""
import sys
import signal
from concurrent import futures
import grpc
from loguru import logger

from app.config import config
from app.proto import rag_service_pb2_grpc
from app.grpc_server import RAGServicer


def configure_logging():
    """配置日志"""
    logger.remove()  # 移除默认处理器
    logger.add(
        sys.stdout,
        format=config.get('logging.format', '<green>{time:YYYY-MM-DD HH:mm:ss}</green> | <level>{level: <8}</level> | <cyan>{name}</cyan>:<cyan>{function}</cyan>:<cyan>{line}</cyan> - <level>{message}</level>'),
        level=config.get('logging.level', 'INFO')
    )
    logger.add(
        "logs/ragljx_py_{time:YYYY-MM-DD}.log",
        rotation="00:00",
        retention="30 days",
        level="INFO"
    )


def serve():
    """启动 gRPC 服务器"""
    # 配置日志
    configure_logging()
    
    # 创建 gRPC 服务器
    server = grpc.server(
        futures.ThreadPoolExecutor(max_workers=10),
        options=[
            ('grpc.max_send_message_length', 100 * 1024 * 1024),  # 100MB
            ('grpc.max_receive_message_length', 100 * 1024 * 1024),  # 100MB
        ]
    )
    
    # 注册服务
    rag_service_pb2_grpc.add_RAGServiceServicer_to_server(
        RAGServicer(),
        server
    )
    
    # 绑定端口
    address = f"{config.grpc_host}:{config.grpc_port}"
    server.add_insecure_port(address)
    
    # 启动服务器
    server.start()
    logger.info(f"RAG gRPC server started on {address}")
    
    # 优雅关闭
    def handle_sigterm(*args):
        logger.info("Received shutdown signal, stopping server...")
        server.stop(grace=10)
        logger.info("Server stopped")
        sys.exit(0)
    
    signal.signal(signal.SIGTERM, handle_sigterm)
    signal.signal(signal.SIGINT, handle_sigterm)
    
    # 等待终止
    server.wait_for_termination()


if __name__ == '__main__':
    serve()


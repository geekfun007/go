"""
Python Asyncio - 流和协议
Streams and Protocols
"""

import asyncio
from typing import Optional


# 示例1: TCP 客户端 (Stream)
async def example1_tcp_client():
    """示例1: TCP 客户端"""
    print("\n=== 示例1: TCP 客户端 ===")
    
    async def fetch_http():
        """发送 HTTP 请求"""
        reader, writer = await asyncio.open_connection('example.com', 80)
        
        request = (
            "GET / HTTP/1.1\r\n"
            "Host: example.com\r\n"
            "Connection: close\r\n"
            "\r\n"
        )
        
        print("发送请求...")
        writer.write(request.encode())
        await writer.drain()
        
        print("接收响应...")
        response = await reader.read(500)
        print(f"响应 (前500字节):\n{response.decode()[:500]}")
        
        writer.close()
        await writer.wait_closed()
    
    print("注意：此示例需要网络连接")
    # await fetch_http()  # 取消注释以实际运行


# 示例2: TCP 服务器
async def example2_tcp_server():
    """示例2: TCP 服务器"""
    print("\n=== 示例2: TCP 服务器 ===")
    
    async def handle_client(reader, writer):
        """处理客户端连接"""
        addr = writer.get_extra_info('peername')
        print(f"客户端连接: {addr}")
        
        try:
            while True:
                data = await reader.readline()
                if not data:
                    break
                
                message = data.decode().strip()
                print(f"收到: {message}")
                
                # 回显消息
                response = f"Echo: {message}\n"
                writer.write(response.encode())
                await writer.drain()
        
        except Exception as e:
            print(f"错误: {e}")
        finally:
            print(f"关闭连接: {addr}")
            writer.close()
            await writer.wait_closed()
    
    async def run_server():
        """运行服务器"""
        server = await asyncio.start_server(
            handle_client, '127.0.0.1', 8888
        )
        
        addr = server.sockets[0].getsockname()
        print(f"服务器启动于 {addr}")
        
        async with server:
            # 实际应用中，这里会一直运行
            print("服务器准备接受连接...")
            print("(示例不会实际启动)")
            # await server.serve_forever()
    
    # 演示代码结构
    print("服务器代码结构已定义")


# 示例3: 简单的聊天服务器
class ChatServer:
    """简单的聊天服务器"""
    
    def __init__(self):
        self.clients = set()
    
    async def handle_client(self, reader, writer):
        """处理客户端"""
        addr = writer.get_extra_info('peername')
        print(f"新客户端: {addr}")
        
        self.clients.add(writer)
        
        try:
            # 发送欢迎消息
            welcome = f"欢迎！当前在线: {len(self.clients)}\n"
            writer.write(welcome.encode())
            await writer.drain()
            
            while True:
                data = await reader.readline()
                if not data:
                    break
                
                message = data.decode().strip()
                print(f"{addr}: {message}")
                
                # 广播消息给所有客户端
                broadcast = f"{addr}: {message}\n"
                await self.broadcast(broadcast.encode(), writer)
        
        finally:
            print(f"客户端断开: {addr}")
            self.clients.remove(writer)
            writer.close()
            await writer.wait_closed()
    
    async def broadcast(self, message: bytes, sender):
        """广播消息"""
        for client in self.clients:
            if client != sender:
                try:
                    client.write(message)
                    await client.drain()
                except Exception as e:
                    print(f"广播错误: {e}")


async def example3_chat_server():
    """示例3: 聊天服务器"""
    print("\n=== 示例3: 聊天服务器 ===")
    print("聊天服务器代码结构已定义")
    print("功能:")
    print("  - 多客户端连接")
    print("  - 消息广播")
    print("  - 连接管理")


# 示例4: 异步文件操作
async def example4_async_file_operations():
    """示例4: 异步文件操作"""
    print("\n=== 示例4: 异步文件操作 ===")
    
    import aiofiles
    
    # 注意：需要安装 aiofiles 库
    print("使用 aiofiles 进行异步文件操作:")
    print("""
    async with aiofiles.open('file.txt', 'w') as f:
        await f.write('Hello, Async World!')
    
    async with aiofiles.open('file.txt', 'r') as f:
        content = await f.read()
        print(content)
    """)


# 示例5: 流式数据处理
async def example5_stream_processing():
    """示例5: 流式数据处理"""
    print("\n=== 示例5: 流式数据处理 ===")
    
    async def data_generator():
        """生成数据流"""
        for i in range(10):
            await asyncio.sleep(0.1)
            yield f"数据 {i}"
    
    async def process_stream():
        """处理数据流"""
        async for data in data_generator():
            print(f"处理: {data}")
    
    await process_stream()


# 示例6: 协议类
class EchoProtocol(asyncio.Protocol):
    """回显协议"""
    
    def connection_made(self, transport):
        """连接建立"""
        peername = transport.get_extra_info('peername')
        print(f'连接来自 {peername}')
        self.transport = transport
    
    def data_received(self, data):
        """接收到数据"""
        message = data.decode()
        print(f'收到数据: {message}')
        
        # 回显
        self.transport.write(data)
    
    def connection_lost(self, exc):
        """连接断开"""
        print('连接关闭')


async def example6_protocol_based():
    """示例6: 基于协议的服务器"""
    print("\n=== 示例6: 基于协议的服务器 ===")
    print("协议类 EchoProtocol 已定义")
    print("使用方式:")
    print("""
    loop = asyncio.get_event_loop()
    server = await loop.create_server(
        EchoProtocol,
        '127.0.0.1', 8888
    )
    """)


# 示例7: UDP 协议
class UDPEchoProtocol(asyncio.DatagramProtocol):
    """UDP 回显协议"""
    
    def connection_made(self, transport):
        self.transport = transport
    
    def datagram_received(self, data, addr):
        message = data.decode()
        print(f'从 {addr} 收到: {message}')
        
        # 回显
        self.transport.sendto(data, addr)


async def example7_udp_protocol():
    """示例7: UDP 协议"""
    print("\n=== 示例7: UDP 协议 ===")
    print("UDP 协议类已定义")
    print("使用方式:")
    print("""
    loop = asyncio.get_event_loop()
    transport, protocol = await loop.create_datagram_endpoint(
        UDPEchoProtocol,
        local_addr=('127.0.0.1', 9999)
    )
    """)


# 示例8: 子进程通信
async def example8_subprocess():
    """示例8: 子进程通信"""
    print("\n=== 示例8: 子进程通信 ===")
    
    # 运行命令
    proc = await asyncio.create_subprocess_shell(
        'echo "Hello from subprocess"',
        stdout=asyncio.subprocess.PIPE,
        stderr=asyncio.subprocess.PIPE
    )
    
    stdout, stderr = await proc.communicate()
    
    print(f'输出: {stdout.decode()}')
    if stderr:
        print(f'错误: {stderr.decode()}')
    print(f'返回码: {proc.returncode}')


# 示例9: 流的缓冲控制
async def example9_flow_control():
    """示例9: 流量控制"""
    print("\n=== 示例9: 流量控制 ===")
    
    class ThrottledProtocol(asyncio.Protocol):
        """限流协议"""
        
        def __init__(self):
            self.pause_writing_called = False
        
        def pause_writing(self):
            """暂停写入（缓冲区满）"""
            print("缓冲区满，暂停写入")
            self.pause_writing_called = True
        
        def resume_writing(self):
            """恢复写入（缓冲区可用）"""
            print("缓冲区可用，恢复写入")
            self.pause_writing_called = False
    
    print("流量控制协议已定义")
    print("当发送缓冲区达到高水位时，会调用 pause_writing()")
    print("当缓冲区降到低水位时，会调用 resume_writing()")


# 示例10: 异步上下文管理器
class AsyncResource:
    """异步资源"""
    
    async def __aenter__(self):
        print("获取资源")
        await asyncio.sleep(0.1)
        return self
    
    async def __aexit__(self, exc_type, exc_val, exc_tb):
        print("释放资源")
        await asyncio.sleep(0.1)
    
    async def use(self):
        print("使用资源")
        await asyncio.sleep(0.1)


async def example10_async_context_manager():
    """示例10: 异步上下文管理器"""
    print("\n=== 示例10: 异步上下文管理器 ===")
    
    async with AsyncResource() as resource:
        await resource.use()
    
    print("资源已自动释放")


async def main():
    """主函数"""
    print("=" * 50)
    print("Python Asyncio - 流和协议")
    print("=" * 50)
    
    await example1_tcp_client()
    await example2_tcp_server()
    await example3_chat_server()
    await example4_async_file_operations()
    await example5_stream_processing()
    await example6_protocol_based()
    await example7_udp_protocol()
    await example8_subprocess()
    await example9_flow_control()
    await example10_async_context_manager()
    
    print("\n" + "=" * 50)
    print("所有示例完成！")
    print("=" * 50)


if __name__ == "__main__":
    asyncio.run(main())

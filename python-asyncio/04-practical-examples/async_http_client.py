"""
Python Asyncio 实战 - 异步 HTTP 客户端
Practical Example: Async HTTP Client
"""

import asyncio
import time
from typing import List, Dict
import json


# 模拟 HTTP 客户端 (实际应使用 aiohttp)
class MockHTTPClient:
    """模拟 HTTP 客户端"""
    
    @staticmethod
    async def get(url: str) -> Dict:
        """GET 请求"""
        # 模拟网络延迟
        await asyncio.sleep(0.5)
        return {
            "url": url,
            "status": 200,
            "data": f"Response from {url}"
        }
    
    @staticmethod
    async def post(url: str, data: Dict) -> Dict:
        """POST 请求"""
        await asyncio.sleep(0.7)
        return {
            "url": url,
            "status": 201,
            "data": f"Created with {data}"
        }


# 示例1: 并发 HTTP 请求
async def example1_concurrent_requests():
    """示例1: 并发 HTTP 请求"""
    print("\n=== 示例1: 并发 HTTP 请求 ===")
    
    urls = [
        "http://api.example.com/users",
        "http://api.example.com/posts",
        "http://api.example.com/comments",
        "http://api.example.com/todos",
    ]
    
    client = MockHTTPClient()
    
    start = time.time()
    
    # 并发请求
    tasks = [client.get(url) for url in urls]
    responses = await asyncio.gather(*tasks)
    
    elapsed = time.time() - start
    
    print(f"完成 {len(responses)} 个请求")
    for response in responses:
        print(f"  {response['url']}: {response['status']}")
    print(f"总耗时: {elapsed:.2f}s")


# 示例2: 批量请求
async def example2_batch_requests():
    """示例2: 批量请求（分批处理）"""
    print("\n=== 示例2: 批量请求 ===")
    
    async def fetch_batch(batch: List[str], client: MockHTTPClient):
        """获取一批数据"""
        tasks = [client.get(url) for url in batch]
        return await asyncio.gather(*tasks)
    
    # 100 个 URL
    urls = [f"http://api.example.com/item/{i}" for i in range(20)]
    
    # 分成 5 个一批
    batch_size = 5
    batches = [urls[i:i+batch_size] for i in range(0, len(urls), batch_size)]
    
    client = MockHTTPClient()
    
    all_responses = []
    for i, batch in enumerate(batches):
        print(f"处理第 {i+1} 批 (共 {len(batches)} 批)...")
        responses = await fetch_batch(batch, client)
        all_responses.extend(responses)
    
    print(f"完成所有请求，共 {len(all_responses)} 个")


# 示例3: 带超时的请求
async def example3_requests_with_timeout():
    """示例3: 带超时的请求"""
    print("\n=== 示例3: 带超时的请求 ===")
    
    client = MockHTTPClient()
    
    async def fetch_with_timeout(url: str, timeout: float):
        """带超时的请求"""
        try:
            response = await asyncio.wait_for(
                client.get(url),
                timeout=timeout
            )
            return response
        except asyncio.TimeoutError:
            return {"url": url, "status": 408, "error": "Request Timeout"}
    
    urls = [
        "http://api.example.com/fast",
        "http://api.example.com/slow",
    ]
    
    # 0.3 秒超时（模拟的请求需要 0.5 秒）
    tasks = [fetch_with_timeout(url, 0.3) for url in urls]
    responses = await asyncio.gather(*tasks)
    
    for response in responses:
        if "error" in response:
            print(f"  {response['url']}: {response['error']}")
        else:
            print(f"  {response['url']}: {response['status']}")


# 示例4: 重试机制
async def example4_retry_mechanism():
    """示例4: 重试机制"""
    print("\n=== 示例4: 重试机制 ===")
    
    class RetryableClient:
        """支持重试的客户端"""
        
        def __init__(self, max_retries: int = 3):
            self.max_retries = max_retries
            self.attempt = 0
        
        async def get_with_retry(self, url: str) -> Dict:
            """带重试的 GET 请求"""
            for attempt in range(self.max_retries):
                try:
                    print(f"  尝试 {attempt + 1}/{self.max_retries}: {url}")
                    
                    # 模拟失败
                    if attempt < 2:
                        await asyncio.sleep(0.2)
                        raise Exception("模拟网络错误")
                    
                    # 成功
                    await asyncio.sleep(0.2)
                    return {"url": url, "status": 200, "attempt": attempt + 1}
                
                except Exception as e:
                    if attempt == self.max_retries - 1:
                        return {"url": url, "status": 500, "error": str(e)}
                    
                    # 指数退避
                    backoff = 2 ** attempt * 0.1
                    await asyncio.sleep(backoff)
    
    client = RetryableClient(max_retries=3)
    response = await client.get_with_retry("http://api.example.com/unstable")
    print(f"最终结果: {response}")


# 示例5: 速率限制
class RateLimiter:
    """速率限制器"""
    
    def __init__(self, rate: int, per: float = 1.0):
        """
        rate: 每 per 秒允许的请求数
        per: 时间窗口（秒）
        """
        self.rate = rate
        self.per = per
        self.allowance = rate
        self.last_check = time.time()
        self.lock = asyncio.Lock()
    
    async def acquire(self):
        """获取许可"""
        async with self.lock:
            current = time.time()
            time_passed = current - self.last_check
            self.last_check = current
            
            # 补充令牌
            self.allowance += time_passed * (self.rate / self.per)
            if self.allowance > self.rate:
                self.allowance = self.rate
            
            # 检查是否有令牌
            if self.allowance < 1.0:
                sleep_time = (1.0 - self.allowance) * (self.per / self.rate)
                await asyncio.sleep(sleep_time)
                self.allowance = 0.0
            else:
                self.allowance -= 1.0


async def example5_rate_limiting():
    """示例5: 速率限制"""
    print("\n=== 示例5: 速率限制 ===")
    
    # 每秒最多 2 个请求
    limiter = RateLimiter(rate=2, per=1.0)
    client = MockHTTPClient()
    
    async def fetch_limited(url: str):
        """限速请求"""
        await limiter.acquire()
        return await client.get(url)
    
    urls = [f"http://api.example.com/item/{i}" for i in range(5)]
    
    start = time.time()
    tasks = [fetch_limited(url) for url in urls]
    responses = await asyncio.gather(*tasks)
    elapsed = time.time() - start
    
    print(f"完成 {len(responses)} 个请求，耗时 {elapsed:.2f}s")
    print(f"平均速率: {len(responses)/elapsed:.2f} req/s")


# 示例6: 连接池
class ConnectionPool:
    """连接池"""
    
    def __init__(self, size: int):
        self.size = size
        self.connections = asyncio.Queue(maxsize=size)
        self._initialize()
    
    def _initialize(self):
        """初始化连接池"""
        for i in range(self.size):
            self.connections.put_nowait(f"conn_{i}")
    
    async def acquire(self):
        """获取连接"""
        return await self.connections.get()
    
    async def release(self, conn):
        """释放连接"""
        await self.connections.put(conn)
    
    async def execute(self, url: str):
        """使用连接池执行请求"""
        conn = await self.acquire()
        try:
            print(f"  使用 {conn} 请求 {url}")
            await asyncio.sleep(0.3)
            return {"url": url, "conn": conn, "status": 200}
        finally:
            await self.release(conn)


async def example6_connection_pool():
    """示例6: 连接池"""
    print("\n=== 示例6: 连接池 ===")
    
    pool = ConnectionPool(size=3)
    urls = [f"http://api.example.com/item/{i}" for i in range(10)]
    
    tasks = [pool.execute(url) for url in urls]
    responses = await asyncio.gather(*tasks)
    
    print(f"完成 {len(responses)} 个请求")


# 示例7: 错误处理和日志
async def example7_error_handling():
    """示例7: 错误处理"""
    print("\n=== 示例7: 错误处理 ===")
    
    async def safe_fetch(url: str, client: MockHTTPClient):
        """安全的请求"""
        try:
            response = await client.get(url)
            return {"success": True, **response}
        except Exception as e:
            print(f"  错误: {url} - {e}")
            return {"success": False, "url": url, "error": str(e)}
    
    client = MockHTTPClient()
    urls = [
        "http://api.example.com/valid",
        "http://api.example.com/error",
    ]
    
    # 使用 return_exceptions=True 收集所有结果
    tasks = [safe_fetch(url, client) for url in urls]
    responses = await asyncio.gather(*tasks, return_exceptions=True)
    
    successful = sum(1 for r in responses if isinstance(r, dict) and r.get("success"))
    print(f"成功: {successful}/{len(responses)}")


# 示例8: 实时流式处理
async def example8_streaming():
    """示例8: 流式处理"""
    print("\n=== 示例8: 流式处理 ===")
    
    async def stream_data(url: str):
        """流式获取数据"""
        print(f"开始流式处理: {url}")
        
        for i in range(5):
            await asyncio.sleep(0.2)
            chunk = f"chunk_{i}"
            print(f"  收到数据块: {chunk}")
            yield chunk
    
    async def process_stream(url: str):
        """处理流"""
        async for chunk in stream_data(url):
            # 处理每个数据块
            print(f"  处理: {chunk}")
    
    await process_stream("http://api.example.com/stream")


# 示例9: 取消和清理
async def example9_cancellation():
    """示例9: 取消和清理"""
    print("\n=== 示例9: 取消和清理 ===")
    
    async def long_request(url: str):
        """长时间请求"""
        try:
            print(f"开始长请求: {url}")
            await asyncio.sleep(5)
            return {"url": url, "status": 200}
        except asyncio.CancelledError:
            print(f"请求被取消: {url}")
            # 清理资源
            raise
    
    task = asyncio.create_task(long_request("http://api.example.com/slow"))
    
    # 等待一段时间后取消
    await asyncio.sleep(1)
    task.cancel()
    
    try:
        await task
    except asyncio.CancelledError:
        print("任务已取消")


async def main():
    """主函数"""
    print("=" * 50)
    print("Python Asyncio - 异步 HTTP 客户端实战")
    print("=" * 50)
    
    await example1_concurrent_requests()
    await example2_batch_requests()
    await example3_requests_with_timeout()
    await example4_retry_mechanism()
    await example5_rate_limiting()
    await example6_connection_pool()
    await example7_error_handling()
    await example8_streaming()
    await example9_cancellation()
    
    print("\n" + "=" * 50)
    print("所有示例完成！")
    print("=" * 50)
    print("\n提示: 实际项目中推荐使用 aiohttp 库")


if __name__ == "__main__":
    asyncio.run(main())

"""
Python Asyncio 高级实战 - 异步网络爬虫
Advanced Example: Async Web Crawler
"""

import asyncio
import time
from typing import Set, List, Dict
from urllib.parse import urljoin, urlparse
import re


class MockHTTPClient:
    """模拟 HTTP 客户端"""
    
    @staticmethod
    async def fetch(url: str) -> Dict:
        """获取页面内容"""
        await asyncio.sleep(0.3)  # 模拟网络延迟
        
        # 模拟 HTML 内容
        html = f"""
        <html>
            <body>
                <h1>Page: {url}</h1>
                <a href="{url}/page1">Link 1</a>
                <a href="{url}/page2">Link 2</a>
                <a href="{url}/page3">Link 3</a>
            </body>
        </html>
        """
        
        return {
            "url": url,
            "status": 200,
            "html": html,
            "title": f"Page {url}"
        }


class AsyncWebCrawler:
    """异步网络爬虫"""
    
    def __init__(self, max_concurrent: int = 10, max_depth: int = 3):
        self.max_concurrent = max_concurrent
        self.max_depth = max_depth
        self.visited: Set[str] = set()
        self.semaphore = asyncio.Semaphore(max_concurrent)
        self.client = MockHTTPClient()
    
    async def crawl(self, start_url: str) -> List[Dict]:
        """开始爬取"""
        results = []
        queue = asyncio.Queue()
        await queue.put((start_url, 0))  # (url, depth)
        
        tasks = []
        for _ in range(self.max_concurrent):
            task = asyncio.create_task(self._worker(queue, results))
            tasks.append(task)
        
        # 等待队列为空
        await queue.join()
        
        # 取消所有 worker
        for task in tasks:
            task.cancel()
        
        await asyncio.gather(*tasks, return_exceptions=True)
        
        return results
    
    async def _worker(self, queue: asyncio.Queue, results: List[Dict]):
        """工作协程"""
        while True:
            try:
                url, depth = await queue.get()
                
                if url in self.visited or depth > self.max_depth:
                    queue.task_done()
                    continue
                
                self.visited.add(url)
                
                async with self.semaphore:
                    page = await self._fetch_page(url)
                    if page:
                        results.append(page)
                        
                        # 提取链接
                        if depth < self.max_depth:
                            links = self._extract_links(page["html"], url)
                            for link in links:
                                if link not in self.visited:
                                    await queue.put((link, depth + 1))
                
                queue.task_done()
            
            except asyncio.CancelledError:
                break
            except Exception as e:
                print(f"Error in worker: {e}")
                queue.task_done()
    
    async def _fetch_page(self, url: str) -> Dict:
        """获取页面"""
        try:
            print(f"Fetching: {url}")
            return await self.client.fetch(url)
        except Exception as e:
            print(f"Error fetching {url}: {e}")
            return None
    
    def _extract_links(self, html: str, base_url: str) -> List[str]:
        """提取链接"""
        links = []
        pattern = r'href=["\']([^"\']+)["\']'
        matches = re.findall(pattern, html)
        
        for match in matches:
            absolute_url = urljoin(base_url, match)
            links.append(absolute_url)
        
        return links[:3]  # 限制链接数量


async def example1_simple_crawler():
    """示例1: 简单爬虫"""
    print("\n=== 示例1: 简单网络爬虫 ===")
    
    crawler = AsyncWebCrawler(max_concurrent=5, max_depth=2)
    
    start = time.time()
    results = await crawler.crawl("http://example.com")
    elapsed = time.time() - start
    
    print(f"\n爬取完成:")
    print(f"  页面数: {len(results)}")
    print(f"  耗时: {elapsed:.2f}s")


# 示例2: 带优先级的爬虫
class PriorityQueue:
    """优先级队列"""
    
    def __init__(self):
        self.queue = asyncio.PriorityQueue()
        self.count = 0
    
    async def put(self, priority: int, item):
        """添加项目"""
        await self.queue.put((priority, self.count, item))
        self.count += 1
    
    async def get(self):
        """获取项目"""
        priority, _, item = await self.queue.get()
        return item
    
    def task_done(self):
        """标记任务完成"""
        self.queue.task_done()
    
    async def join(self):
        """等待完成"""
        await self.queue.join()


class PriorityCrawler(AsyncWebCrawler):
    """带优先级的爬虫"""
    
    async def crawl(self, start_url: str) -> List[Dict]:
        """开始爬取"""
        results = []
        queue = PriorityQueue()
        await queue.put(0, (start_url, 0))  # 优先级 0（最高）
        
        tasks = []
        for _ in range(self.max_concurrent):
            task = asyncio.create_task(self._priority_worker(queue, results))
            tasks.append(task)
        
        await queue.join()
        
        for task in tasks:
            task.cancel()
        
        await asyncio.gather(*tasks, return_exceptions=True)
        
        return results
    
    async def _priority_worker(self, queue: PriorityQueue, results: List[Dict]):
        """优先级工作协程"""
        while True:
            try:
                url, depth = await queue.get()
                
                if url in self.visited or depth > self.max_depth:
                    queue.task_done()
                    continue
                
                self.visited.add(url)
                
                async with self.semaphore:
                    page = await self._fetch_page(url)
                    if page:
                        results.append(page)
                        
                        if depth < self.max_depth:
                            links = self._extract_links(page["html"], url)
                            for i, link in enumerate(links):
                                if link not in self.visited:
                                    # 优先级随深度增加
                                    await queue.put(depth + 1, (link, depth + 1))
                
                queue.task_done()
            
            except asyncio.CancelledError:
                break
            except Exception as e:
                print(f"Error: {e}")
                queue.task_done()


async def example2_priority_crawler():
    """示例2: 优先级爬虫"""
    print("\n=== 示例2: 优先级爬虫 ===")
    
    crawler = PriorityCrawler(max_concurrent=5, max_depth=2)
    results = await crawler.crawl("http://example.com")
    
    print(f"爬取 {len(results)} 个页面")


# 示例3: 带缓存的爬虫
class CachedCrawler(AsyncWebCrawler):
    """带缓存的爬虫"""
    
    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.cache: Dict[str, Dict] = {}
        self.cache_lock = asyncio.Lock()
    
    async def _fetch_page(self, url: str) -> Dict:
        """带缓存的页面获取"""
        # 检查缓存
        async with self.cache_lock:
            if url in self.cache:
                print(f"Cache hit: {url}")
                return self.cache[url]
        
        # 获取页面
        page = await super()._fetch_page(url)
        
        # 存入缓存
        if page:
            async with self.cache_lock:
                self.cache[url] = page
        
        return page


async def example3_cached_crawler():
    """示例3: 带缓存的爬虫"""
    print("\n=== 示例3: 带缓存的爬虫 ===")
    
    crawler = CachedCrawler(max_concurrent=5, max_depth=2)
    
    # 第一次爬取
    print("第一次爬取:")
    start = time.time()
    results1 = await crawler.crawl("http://example.com")
    elapsed1 = time.time() - start
    print(f"耗时: {elapsed1:.2f}s")
    
    # 第二次爬取（使用缓存）
    print("\n第二次爬取（使用缓存）:")
    crawler.visited.clear()  # 清除访问记录
    start = time.time()
    results2 = await crawler.crawl("http://example.com")
    elapsed2 = time.time() - start
    print(f"耗时: {elapsed2:.2f}s")
    print(f"速度提升: {elapsed1/elapsed2:.2f}x")


# 示例4: 分布式爬虫（生产者-消费者模式）
class DistributedCrawler:
    """分布式爬虫"""
    
    def __init__(self, num_producers: int = 2, num_consumers: int = 3):
        self.num_producers = num_producers
        self.num_consumers = num_consumers
        self.url_queue = asyncio.Queue()
        self.result_queue = asyncio.Queue()
        self.visited: Set[str] = set()
        self.client = MockHTTPClient()
    
    async def crawl(self, start_urls: List[str]) -> List[Dict]:
        """开始爬取"""
        # 添加初始 URL
        for url in start_urls:
            await self.url_queue.put(url)
        
        # 启动生产者
        producers = [
            asyncio.create_task(self._producer(i))
            for i in range(self.num_producers)
        ]
        
        # 启动消费者
        consumers = [
            asyncio.create_task(self._consumer(i))
            for i in range(self.num_consumers)
        ]
        
        # 等待一段时间
        await asyncio.sleep(3)
        
        # 取消所有任务
        for task in producers + consumers:
            task.cancel()
        
        await asyncio.gather(*producers, *consumers, return_exceptions=True)
        
        # 收集结果
        results = []
        while not self.result_queue.empty():
            results.append(await self.result_queue.get())
        
        return results
    
    async def _producer(self, id: int):
        """生产者：提取链接"""
        print(f"生产者 {id} 启动")
        
        while True:
            try:
                url = await asyncio.wait_for(self.url_queue.get(), timeout=1.0)
                
                if url not in self.visited:
                    self.visited.add(url)
                    page = await self.client.fetch(url)
                    
                    # 提取新链接
                    links = self._extract_links(page["html"], url)
                    for link in links:
                        if link not in self.visited:
                            await self.url_queue.put(link)
                    
                    await self.result_queue.put(page)
                    print(f"生产者 {id} 处理: {url}")
            
            except asyncio.TimeoutError:
                continue
            except asyncio.CancelledError:
                break
            except Exception as e:
                print(f"生产者 {id} 错误: {e}")
    
    async def _consumer(self, id: int):
        """消费者：处理页面"""
        print(f"消费者 {id} 启动")
        
        while True:
            try:
                page = await asyncio.wait_for(
                    self.result_queue.get(),
                    timeout=1.0
                )
                
                # 处理页面（这里只是打印）
                print(f"消费者 {id} 处理: {page['url']}")
                await asyncio.sleep(0.1)
                
                # 放回队列供其他操作使用
                await self.result_queue.put(page)
            
            except asyncio.TimeoutError:
                continue
            except asyncio.CancelledError:
                break
            except Exception as e:
                print(f"消费者 {id} 错误: {e}")
    
    def _extract_links(self, html: str, base_url: str) -> List[str]:
        """提取链接"""
        links = []
        pattern = r'href=["\']([^"\']+)["\']'
        matches = re.findall(pattern, html)
        
        for match in matches:
            absolute_url = urljoin(base_url, match)
            links.append(absolute_url)
        
        return links[:2]


async def example4_distributed_crawler():
    """示例4: 分布式爬虫"""
    print("\n=== 示例4: 分布式爬虫 ===")
    
    crawler = DistributedCrawler(num_producers=2, num_consumers=3)
    
    start_urls = [
        "http://example.com/page1",
        "http://example.com/page2",
    ]
    
    results = await crawler.crawl(start_urls)
    print(f"\n爬取完成，共 {len(results)} 个页面")


async def main():
    """主函数"""
    print("=" * 50)
    print("Python Asyncio - 异步网络爬虫")
    print("=" * 50)
    
    await example1_simple_crawler()
    await example2_priority_crawler()
    await example3_cached_crawler()
    await example4_distributed_crawler()
    
    print("\n" + "=" * 50)
    print("所有示例完成！")
    print("=" * 50)


if __name__ == "__main__":
    asyncio.run(main())

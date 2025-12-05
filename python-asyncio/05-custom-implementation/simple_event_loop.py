"""
自定义 Asyncio 实现 - 简单的事件循环
Custom Asyncio Implementation - Simple Event Loop

这是一个简化的 asyncio 事件循环实现，用于教学目的
展示 asyncio 的核心原理
"""

import time
from collections import deque
from typing import Any, Callable, Coroutine, Optional
import types


class Future:
    """简单的 Future 实现"""
    
    def __init__(self):
        self._result = None
        self._exception = None
        self._done = False
        self._callbacks = []
    
    def set_result(self, result: Any):
        """设置结果"""
        if self._done:
            raise RuntimeError("Future already done")
        self._result = result
        self._done = True
        self._schedule_callbacks()
    
    def set_exception(self, exception: Exception):
        """设置异常"""
        if self._done:
            raise RuntimeError("Future already done")
        self._exception = exception
        self._done = True
        self._schedule_callbacks()
    
    def done(self) -> bool:
        """是否完成"""
        return self._done
    
    def result(self) -> Any:
        """获取结果"""
        if not self._done:
            raise RuntimeError("Future not done")
        if self._exception:
            raise self._exception
        return self._result
    
    def add_done_callback(self, callback: Callable):
        """添加完成回调"""
        if self._done:
            callback(self)
        else:
            self._callbacks.append(callback)
    
    def _schedule_callbacks(self):
        """调度回调"""
        for callback in self._callbacks:
            callback(self)
        self._callbacks.clear()


class Task(Future):
    """简单的 Task 实现"""
    
    def __init__(self, coro: Coroutine, loop):
        super().__init__()
        self._coro = coro
        self._loop = loop
        
        # 立即调度第一步
        self._loop.call_soon(self._step)
    
    def _step(self, exc: Optional[Exception] = None):
        """执行协程的一步"""
        try:
            if exc:
                # 向协程抛出异常
                result = self._coro.throw(exc)
            else:
                # 继续执行协程
                result = self._coro.send(None)
            
            # 如果返回的是 Future，等待它完成
            if isinstance(result, Future):
                result.add_done_callback(self._wakeup)
            else:
                # 不应该到这里
                raise RuntimeError(f"Unexpected yield: {result}")
        
        except StopIteration as e:
            # 协程完成
            self.set_result(e.value if e.value else None)
        
        except Exception as e:
            # 协程异常
            self.set_exception(e)
    
    def _wakeup(self, future: Future):
        """被 Future 唤醒"""
        try:
            result = future.result()
            self._step()
        except Exception as exc:
            self._step(exc)


class SimpleEventLoop:
    """简单的事件循环实现"""
    
    def __init__(self):
        self._ready = deque()  # 准备执行的回调
        self._scheduled = []  # 定时回调
        self._stopping = False
        self._running = False
    
    def call_soon(self, callback: Callable, *args):
        """尽快调用回调"""
        self._ready.append((callback, args))
    
    def call_later(self, delay: float, callback: Callable, *args):
        """延迟调用"""
        when = time.time() + delay
        self._scheduled.append((when, callback, args))
        self._scheduled.sort()  # 按时间排序
    
    def create_task(self, coro: Coroutine) -> Task:
        """创建任务"""
        return Task(coro, self)
    
    def create_future(self) -> Future:
        """创建 Future"""
        return Future()
    
    def run_until_complete(self, coro: Coroutine) -> Any:
        """运行直到协程完成"""
        task = self.create_task(coro)
        return self.run_until_future(task)
    
    def run_until_future(self, future: Future) -> Any:
        """运行直到 Future 完成"""
        if self._running:
            raise RuntimeError("Event loop is already running")
        
        self._running = True
        
        try:
            while not future.done():
                self._run_once()
            return future.result()
        finally:
            self._running = False
    
    def _run_once(self):
        """运行一次迭代"""
        # 处理定时回调
        now = time.time()
        while self._scheduled and self._scheduled[0][0] <= now:
            _, callback, args = self._scheduled.pop(0)
            self.call_soon(callback, *args)
        
        # 处理准备好的回调
        num_ready = len(self._ready)
        for _ in range(num_ready):
            callback, args = self._ready.popleft()
            try:
                callback(*args)
            except Exception as e:
                print(f"Exception in callback: {e}")
        
        # 如果没有任务，短暂休眠
        if not self._ready and self._scheduled:
            delay = self._scheduled[0][0] - time.time()
            if delay > 0:
                time.sleep(min(delay, 0.01))
        elif not self._ready:
            time.sleep(0.001)
    
    def stop(self):
        """停止事件循环"""
        self._stopping = True


# 辅助函数：模拟 asyncio.sleep
def sleep(delay: float, loop: SimpleEventLoop) -> Future:
    """模拟异步睡眠"""
    future = loop.create_future()
    loop.call_later(delay, future.set_result, None)
    return future


# 演示程序
async def hello_world(loop: SimpleEventLoop):
    """Hello World 协程"""
    print("Hello")
    await sleep(1, loop)
    print("World")
    return "Done"


async def count_to(n: int, loop: SimpleEventLoop):
    """计数协程"""
    for i in range(1, n + 1):
        print(f"Count: {i}")
        await sleep(0.5, loop)
    return f"Counted to {n}"


async def concurrent_tasks(loop: SimpleEventLoop):
    """并发任务演示"""
    print("\n=== 并发任务演示 ===")
    
    # 创建多个任务
    task1 = loop.create_task(count_to(3, loop))
    task2 = loop.create_task(count_to(3, loop))
    
    # 等待任务完成
    result1 = await task1
    result2 = await task2
    
    print(f"Task 1: {result1}")
    print(f"Task 2: {result2}")
    
    return "All tasks completed"


def demo_simple_event_loop():
    """演示简单事件循环"""
    print("=" * 50)
    print("自定义事件循环演示")
    print("=" * 50)
    
    loop = SimpleEventLoop()
    
    # 示例1: 基本使用
    print("\n示例1: Hello World")
    result = loop.run_until_complete(hello_world(loop))
    print(f"Result: {result}")
    
    # 示例2: 并发任务
    print("\n示例2: 并发任务")
    result = loop.run_until_complete(concurrent_tasks(loop))
    print(f"Result: {result}")
    
    print("\n" + "=" * 50)
    print("演示完成！")
    print("=" * 50)


def explain_implementation():
    """解释实现原理"""
    print("\n" + "=" * 50)
    print("实现原理说明")
    print("=" * 50)
    print("""
这个简化的 asyncio 实现包含以下核心组件:

1. Future 类:
   - 代表一个未来的结果
   - 可以设置结果或异常
   - 支持添加完成回调

2. Task 类:
   - 继承自 Future
   - 封装协程的执行
   - 使用 send() 和 throw() 驱动协程
   - 处理 StopIteration 表示完成

3. EventLoop 类:
   - 维护准备队列 (_ready) 和定时队列 (_scheduled)
   - call_soon: 添加到准备队列
   - call_later: 添加到定时队列
   - _run_once: 执行一次事件循环迭代
   - run_until_complete: 运行直到任务完成

4. 协程执行流程:
   a. 创建 Task，将协程的第一步加入准备队列
   b. 事件循环从准备队列取出回调执行
   c. 执行协程的一步 (coro.send(None))
   d. 如果遇到 await Future，注册回调等待
   e. Future 完成时，唤醒协程继续执行
   f. 重复直到协程完成 (StopIteration)

5. 关键技术:
   - 生成器协议 (send/throw/StopIteration)
   - 事件驱动架构
   - 回调和 Future 模式
   - 时间调度

这个实现省略了:
   - I/O 多路复用 (select/epoll)
   - 信号处理
   - 子进程管理
   - 完整的异常处理
   - 性能优化

但它展示了 asyncio 的核心原理！
    """)


if __name__ == "__main__":
    # 运行演示
    demo_simple_event_loop()
    
    # 解释原理
    explain_implementation()

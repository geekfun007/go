"""
Python Asyncio - 事件循环
Event Loop Basics
"""

import asyncio
import time
from typing import Callable


# 示例1: 获取和使用事件循环
async def example1_get_event_loop():
    """示例1: 获取事件循环"""
    print("\n=== 示例1: 获取事件循环 ===")
    
    loop = asyncio.get_event_loop()
    print(f"事件循环: {loop}")
    print(f"是否运行中: {loop.is_running()}")
    print(f"是否关闭: {loop.is_closed()}")


# 示例2: 在事件循环中调度回调
def callback_function(name: str, loop):
    """回调函数"""
    print(f"回调 {name} 执行于 {time.time():.2f}")


async def example2_schedule_callbacks():
    """示例2: 调度回调"""
    print("\n=== 示例2: 调度回调 ===")
    
    loop = asyncio.get_event_loop()
    
    # call_soon: 下一次循环迭代时调用
    loop.call_soon(callback_function, "soon", loop)
    
    # call_later: 延迟调用
    loop.call_later(1, callback_function, "later-1s", loop)
    loop.call_later(2, callback_function, "later-2s", loop)
    
    # call_at: 在特定时间调用
    loop.call_at(loop.time() + 1.5, callback_function, "at-1.5s", loop)
    
    print("回调已调度")
    await asyncio.sleep(3)


# 示例3: 在事件循环中运行协程
async def example3_run_coroutine_threadsafe():
    """示例3: 线程安全地运行协程"""
    print("\n=== 示例3: 在事件循环中运行协程 ===")
    
    async def my_coroutine(name: str):
        print(f"协程 {name} 开始")
        await asyncio.sleep(1)
        print(f"协程 {name} 结束")
        return f"{name} 结果"
    
    loop = asyncio.get_event_loop()
    
    # 创建任务
    task = loop.create_task(my_coroutine("任务1"))
    result = await task
    print(f"结果: {result}")


# 示例4: 事件循环的未来对象 (Future)
async def example4_futures():
    """示例4: Future 对象"""
    print("\n=== 示例4: Future 对象 ===")
    
    loop = asyncio.get_event_loop()
    
    # 创建一个 Future 对象
    future = loop.create_future()
    
    async def set_future_result():
        """设置 future 的结果"""
        await asyncio.sleep(1)
        print("设置 future 结果")
        future.set_result("Future 完成")
    
    # 启动任务
    asyncio.create_task(set_future_result())
    
    # 等待 future 完成
    print("等待 future...")
    result = await future
    print(f"Future 结果: {result}")


# 示例5: 自定义事件循环策略
class CustomEventLoopPolicy(asyncio.DefaultEventLoopPolicy):
    """自定义事件循环策略"""
    
    def new_event_loop(self):
        """创建新的事件循环"""
        print("创建自定义事件循环")
        return super().new_event_loop()


async def example5_custom_policy():
    """示例5: 自定义事件循环策略"""
    print("\n=== 示例5: 自定义事件循环策略 ===")
    
    # 注意：这只是演示，实际不修改策略
    print(f"当前策略: {type(asyncio.get_event_loop_policy()).__name__}")


# 示例6: 事件循环异常处理
def exception_handler(loop, context):
    """自定义异常处理器"""
    print(f"捕获异常: {context['message']}")
    if 'exception' in context:
        print(f"异常类型: {type(context['exception']).__name__}")


async def example6_exception_handler():
    """示例6: 异常处理"""
    print("\n=== 示例6: 事件循环异常处理 ===")
    
    loop = asyncio.get_event_loop()
    
    # 设置异常处理器
    loop.set_exception_handler(exception_handler)
    
    async def buggy_coroutine():
        """有问题的协程"""
        await asyncio.sleep(0.1)
        raise ValueError("故意的错误")
    
    # 创建任务但不等待（异常会被异常处理器捕获）
    task = asyncio.create_task(buggy_coroutine())
    
    try:
        await task
    except ValueError:
        print("在 await 时捕获异常")


# 示例7: 运行直到完成
async def example7_run_until_complete():
    """示例7: 运行直到完成"""
    print("\n=== 示例7: 运行直到完成 ===")
    
    async def simple_task(n: int):
        await asyncio.sleep(0.5)
        return n * 2
    
    # asyncio.run() 内部就是使用 run_until_complete()
    print("使用 asyncio.run() 类似于:")
    print("  loop = asyncio.new_event_loop()")
    print("  try:")
    print("      loop.run_until_complete(coroutine)")
    print("  finally:")
    print("      loop.close()")


# 示例8: 事件循环调试
async def example8_debug_mode():
    """示例8: 调试模式"""
    print("\n=== 示例8: 调试模式 ===")
    
    loop = asyncio.get_event_loop()
    
    # 检查调试模式
    print(f"调试模式: {loop.get_debug()}")
    
    # 启用调试模式会:
    # 1. 记录未等待的协程
    # 2. 记录耗时的回调
    # 3. 提供更详细的错误信息
    print("\n可以通过以下方式启用调试:")
    print("  export PYTHONASYNCIODEBUG=1")
    print("  或 loop.set_debug(True)")


# 示例9: 事件循环中的同步代码
async def example9_run_in_executor():
    """示例9: 在执行器中运行同步代码"""
    print("\n=== 示例9: 运行同步代码 ===")
    
    import concurrent.futures
    
    def blocking_io():
        """阻塞式 I/O"""
        print("开始阻塞操作...")
        time.sleep(2)
        print("阻塞操作完成")
        return "阻塞操作结果"
    
    loop = asyncio.get_event_loop()
    
    # 在默认执行器中运行
    result = await loop.run_in_executor(None, blocking_io)
    print(f"结果: {result}")
    
    # 在自定义执行器中运行
    with concurrent.futures.ThreadPoolExecutor() as executor:
        result = await loop.run_in_executor(executor, blocking_io)
        print(f"结果: {result}")


# 示例10: 事件循环生命周期回调
async def example10_lifecycle_callbacks():
    """示例10: 生命周期回调"""
    print("\n=== 示例10: 生命周期回调 ===")
    
    async def startup():
        print("应用启动")
    
    async def shutdown():
        print("应用关闭")
    
    async def work():
        print("执行工作...")
        await asyncio.sleep(1)
        print("工作完成")
    
    try:
        await startup()
        await work()
    finally:
        await shutdown()


async def main():
    """主函数"""
    print("=" * 50)
    print("Python Asyncio - 事件循环")
    print("=" * 50)
    
    await example1_get_event_loop()
    await example2_schedule_callbacks()
    await example3_run_coroutine_threadsafe()
    await example4_futures()
    await example5_custom_policy()
    await example6_exception_handler()
    await example7_run_until_complete()
    await example8_debug_mode()
    await example9_run_in_executor()
    await example10_lifecycle_callbacks()
    
    print("\n" + "=" * 50)
    print("所有示例完成！")
    print("=" * 50)


if __name__ == "__main__":
    asyncio.run(main())

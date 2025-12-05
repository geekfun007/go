"""
Python Asyncio 基础 - 协程
Basic Coroutines in Python Asyncio
"""

import asyncio
import time


# 示例1: 基本的协程
async def simple_coroutine():
    """简单的协程函数"""
    print("协程开始执行")
    await asyncio.sleep(1)
    print("协程执行完成")
    return "结果"


async def example1_basic_coroutine():
    """示例1: 基本协程"""
    print("\n=== 示例1: 基本协程 ===")
    result = await simple_coroutine()
    print(f"返回值: {result}")


# 示例2: 多个协程并发执行
async def task_with_delay(name: str, delay: float):
    """带延迟的任务"""
    print(f"任务 {name} 开始")
    await asyncio.sleep(delay)
    print(f"任务 {name} 完成 (延迟 {delay}s)")
    return f"{name} 结果"


async def example2_concurrent_coroutines():
    """示例2: 并发执行多个协程"""
    print("\n=== 示例2: 并发执行协程 ===")
    
    start = time.time()
    
    # 并发执行
    results = await asyncio.gather(
        task_with_delay("A", 2),
        task_with_delay("B", 1),
        task_with_delay("C", 1.5),
    )
    
    elapsed = time.time() - start
    print(f"结果: {results}")
    print(f"总耗时: {elapsed:.2f}s (并发执行)")


# 示例3: 顺序执行 vs 并发执行
async def example3_sequential_vs_concurrent():
    """示例3: 顺序执行 vs 并发执行"""
    print("\n=== 示例3: 顺序 vs 并发 ===")
    
    # 顺序执行
    print("顺序执行:")
    start = time.time()
    await task_with_delay("1", 1)
    await task_with_delay("2", 1)
    await task_with_delay("3", 1)
    print(f"顺序执行耗时: {time.time() - start:.2f}s")
    
    # 并发执行
    print("\n并发执行:")
    start = time.time()
    await asyncio.gather(
        task_with_delay("1", 1),
        task_with_delay("2", 1),
        task_with_delay("3", 1),
    )
    print(f"并发执行耗时: {time.time() - start:.2f}s")


# 示例4: 创建任务 (Task)
async def example4_create_tasks():
    """示例4: 创建任务"""
    print("\n=== 示例4: 创建任务 ===")
    
    # 创建任务
    task1 = asyncio.create_task(task_with_delay("Task1", 2))
    task2 = asyncio.create_task(task_with_delay("Task2", 1))
    
    print("任务已创建，开始等待...")
    
    # 等待任务完成
    result1 = await task1
    result2 = await task2
    
    print(f"任务1结果: {result1}")
    print(f"任务2结果: {result2}")


# 示例5: 协程超时
async def slow_operation():
    """慢速操作"""
    print("慢速操作开始...")
    await asyncio.sleep(5)
    return "操作完成"


async def example5_timeout():
    """示例5: 超时处理"""
    print("\n=== 示例5: 超时处理 ===")
    
    try:
        result = await asyncio.wait_for(slow_operation(), timeout=2.0)
        print(f"结果: {result}")
    except asyncio.TimeoutError:
        print("操作超时！")


# 示例6: 协程异常处理
async def faulty_coroutine(will_fail: bool):
    """可能失败的协程"""
    await asyncio.sleep(1)
    if will_fail:
        raise ValueError("协程执行失败")
    return "成功"


async def example6_exception_handling():
    """示例6: 异常处理"""
    print("\n=== 示例6: 异常处理 ===")
    
    try:
        result = await faulty_coroutine(True)
        print(f"结果: {result}")
    except ValueError as e:
        print(f"捕获异常: {e}")
    
    # gather 中的异常处理
    results = await asyncio.gather(
        faulty_coroutine(False),
        faulty_coroutine(True),
        return_exceptions=True  # 返回异常而不是抛出
    )
    
    for i, result in enumerate(results):
        if isinstance(result, Exception):
            print(f"任务 {i} 失败: {result}")
        else:
            print(f"任务 {i} 成功: {result}")


# 示例7: 等待第一个完成
async def example7_wait_first():
    """示例7: 等待第一个完成"""
    print("\n=== 示例7: 等待第一个完成 ===")
    
    tasks = [
        asyncio.create_task(task_with_delay("快", 1)),
        asyncio.create_task(task_with_delay("中", 2)),
        asyncio.create_task(task_with_delay("慢", 3)),
    ]
    
    done, pending = await asyncio.wait(
        tasks,
        return_when=asyncio.FIRST_COMPLETED
    )
    
    print(f"完成的任务数: {len(done)}")
    print(f"待完成的任务数: {len(pending)}")
    
    # 取消待完成的任务
    for task in pending:
        task.cancel()
    
    # 获取完成任务的结果
    for task in done:
        print(f"结果: {task.result()}")


# 示例8: 协程链式调用
async def fetch_data():
    """获取数据"""
    await asyncio.sleep(1)
    return {"id": 1, "name": "数据"}


async def process_data(data):
    """处理数据"""
    await asyncio.sleep(1)
    return {**data, "processed": True}


async def save_data(data):
    """保存数据"""
    await asyncio.sleep(0.5)
    print(f"保存数据: {data}")
    return "保存成功"


async def example8_chaining():
    """示例8: 协程链式调用"""
    print("\n=== 示例8: 协程链式调用 ===")
    
    data = await fetch_data()
    print(f"获取数据: {data}")
    
    processed = await process_data(data)
    print(f"处理数据: {processed}")
    
    result = await save_data(processed)
    print(f"结果: {result}")


# 示例9: 使用 asyncio.as_completed
async def example9_as_completed():
    """示例9: 按完成顺序处理"""
    print("\n=== 示例9: 按完成顺序处理 ===")
    
    tasks = [
        task_with_delay("任务1", 3),
        task_with_delay("任务2", 1),
        task_with_delay("任务3", 2),
    ]
    
    for coro in asyncio.as_completed(tasks):
        result = await coro
        print(f"收到结果: {result}")


async def main():
    """主函数"""
    print("=" * 50)
    print("Python Asyncio 基础 - 协程示例")
    print("=" * 50)
    
    await example1_basic_coroutine()
    await example2_concurrent_coroutines()
    await example3_sequential_vs_concurrent()
    await example4_create_tasks()
    await example5_timeout()
    await example6_exception_handling()
    await example7_wait_first()
    await example8_chaining()
    await example9_as_completed()
    
    print("\n" + "=" * 50)
    print("所有示例完成！")
    print("=" * 50)


if __name__ == "__main__":
    asyncio.run(main())

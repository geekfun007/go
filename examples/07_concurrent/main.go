// 示例 7：并发安全迭代器
// 演示如何为并发数据结构实现安全的迭代器
package main

import (
	"fmt"
	"sync"
	"time"
)

// ========== 并发安全 Map ==========

// SafeMap 线程安全的 map
type SafeMap[K comparable, V any] struct {
	mu sync.RWMutex
	m  map[K]V
}

// NewSafeMap 创建安全 map
func NewSafeMap[K comparable, V any]() *SafeMap[K, V] {
	return &SafeMap[K, V]{m: make(map[K]V)}
}

// Set 设置键值
func (s *SafeMap[K, V]) Set(key K, value V) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[key] = value
}

// Get 获取值
func (s *SafeMap[K, V]) Get(key K) (V, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.m[key]
	return v, ok
}

// Delete 删除键
func (s *SafeMap[K, V]) Delete(key K) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.m, key)
}

// All 返回快照迭代器（安全）
// 使用快照方式避免在遍历时持有锁
func (s *SafeMap[K, V]) All() func(yield func(K, V) bool) {
	return func(yield func(K, V) bool) {
		// 创建快照
		s.mu.RLock()
		snapshot := make([]struct {
			k K
			v V
		}, 0, len(s.m))
		for k, v := range s.m {
			snapshot = append(snapshot, struct {
				k K
				v V
			}{k, v})
		}
		s.mu.RUnlock()

		// 遍历快照（不持有锁）
		for _, item := range snapshot {
			if !yield(item.k, item.v) {
				return
			}
		}
	}
}

// Keys 返回所有键的迭代器
func (s *SafeMap[K, V]) Keys() func(yield func(K) bool) {
	return func(yield func(K) bool) {
		s.mu.RLock()
		keys := make([]K, 0, len(s.m))
		for k := range s.m {
			keys = append(keys, k)
		}
		s.mu.RUnlock()

		for _, k := range keys {
			if !yield(k) {
				return
			}
		}
	}
}

// Values 返回所有值的迭代器
func (s *SafeMap[K, V]) Values() func(yield func(V) bool) {
	return func(yield func(V) bool) {
		s.mu.RLock()
		values := make([]V, 0, len(s.m))
		for _, v := range s.m {
			values = append(values, v)
		}
		s.mu.RUnlock()

		for _, v := range values {
			if !yield(v) {
				return
			}
		}
	}
}

// ========== 并发安全切片 ==========

// SafeSlice 线程安全的切片
type SafeSlice[T any] struct {
	mu    sync.RWMutex
	items []T
}

// NewSafeSlice 创建安全切片
func NewSafeSlice[T any]() *SafeSlice[T] {
	return &SafeSlice[T]{items: make([]T, 0)}
}

// Append 追加元素
func (s *SafeSlice[T]) Append(item T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = append(s.items, item)
}

// Len 返回长度
func (s *SafeSlice[T]) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.items)
}

// All 返回快照迭代器
func (s *SafeSlice[T]) All() func(yield func(int, T) bool) {
	return func(yield func(int, T) bool) {
		s.mu.RLock()
		snapshot := make([]T, len(s.items))
		copy(snapshot, s.items)
		s.mu.RUnlock()

		for i, v := range snapshot {
			if !yield(i, v) {
				return
			}
		}
	}
}

// ========== 带锁迭代器（高级用法） ==========

// LockedIterator 在迭代期间持有锁的迭代器
// 注意：消费者不应在回调中进行耗时操作
func (s *SafeMap[K, V]) LockedAll() func(yield func(K, V) bool) {
	return func(yield func(K, V) bool) {
		s.mu.RLock()
		defer s.mu.RUnlock()

		for k, v := range s.m {
			if !yield(k, v) {
				return
			}
		}
	}
}

func main() {
	fmt.Println("=== 并发安全迭代器示例 ===")
	fmt.Println()

	// 示例 1: SafeMap 基本使用
	fmt.Println("1. SafeMap 基本使用:")
	m := NewSafeMap[string, int]()
	m.Set("apple", 1)
	m.Set("banana", 2)
	m.Set("cherry", 3)

	for k, v := range m.All() {
		fmt.Printf("  %s: %d\n", k, v)
	}
	fmt.Println()

	// 示例 2: 并发安全演示
	fmt.Println("2. 并发安全演示（一边迭代一边修改）:")
	m2 := NewSafeMap[int, string]()
	for i := 0; i < 5; i++ {
		m2.Set(i, fmt.Sprintf("value-%d", i))
	}

	var wg sync.WaitGroup

	// 启动一个 goroutine 不断修改 map
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			m2.Set(100+i, fmt.Sprintf("new-%d", i))
			time.Sleep(10 * time.Millisecond)
		}
	}()

	// 主 goroutine 遍历（使用快照，不会受影响）
	fmt.Println("  开始遍历（快照方式）:")
	for k, v := range m2.All() {
		fmt.Printf("    %d: %s\n", k, v)
		time.Sleep(20 * time.Millisecond)
	}

	wg.Wait()
	fmt.Println("  遍历完成，当前 map 大小:", func() int {
		count := 0
		for _ = range m2.Keys() {
			count++
		}
		return count
	}())
	fmt.Println()

	// 示例 3: SafeSlice
	fmt.Println("3. SafeSlice 使用:")
	slice := NewSafeSlice[string]()
	slice.Append("first")
	slice.Append("second")
	slice.Append("third")

	for i, v := range slice.All() {
		fmt.Printf("  [%d] %s\n", i, v)
	}
	fmt.Println()

	// 示例 4: Keys 和 Values 迭代器
	fmt.Println("4. Keys 和 Values 迭代器:")
	config := NewSafeMap[string, string]()
	config.Set("host", "localhost")
	config.Set("port", "8080")
	config.Set("protocol", "https")

	fmt.Print("  Keys: ")
	for k := range config.Keys() {
		fmt.Printf("%s ", k)
	}
	fmt.Println()

	fmt.Print("  Values: ")
	for v := range config.Values() {
		fmt.Printf("%s ", v)
	}
	fmt.Println()
	fmt.Println()

	// 示例 5: 迭代中提前退出
	fmt.Println("5. 迭代中提前退出:")
	largeMap := NewSafeMap[int, int]()
	for i := 0; i < 100; i++ {
		largeMap.Set(i, i*i)
	}

	count := 0
	for k, v := range largeMap.All() {
		fmt.Printf("  %d: %d\n", k, v)
		count++
		if count >= 5 {
			fmt.Println("  ... (只显示前 5 个)")
			break
		}
	}
	fmt.Println()

	// 示例 6: 带锁迭代（需谨慎使用）
	fmt.Println("6. 带锁迭代 (LockedAll):")
	fmt.Println("  注意：在迭代期间持有读锁，不要进行耗时操作")
	smallMap := NewSafeMap[string, int]()
	smallMap.Set("a", 1)
	smallMap.Set("b", 2)
	smallMap.Set("c", 3)

	for k, v := range smallMap.LockedAll() {
		fmt.Printf("  %s: %d\n", k, v)
	}
}

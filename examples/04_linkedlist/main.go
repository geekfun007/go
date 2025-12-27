// 示例 4：链表遍历
// 演示如何为自定义数据结构实现迭代器
package main

import "fmt"

// Node 链表节点
type Node[T any] struct {
	Value T
	Next  *Node[T]
}

// LinkedList 单向链表
type LinkedList[T any] struct {
	Head *Node[T]
	size int
}

// NewLinkedList 创建新链表
func NewLinkedList[T any]() *LinkedList[T] {
	return &LinkedList[T]{}
}

// Append 追加元素
func (l *LinkedList[T]) Append(value T) {
	newNode := &Node[T]{Value: value}
	if l.Head == nil {
		l.Head = newNode
	} else {
		current := l.Head
		for current.Next != nil {
			current = current.Next
		}
		current.Next = newNode
	}
	l.size++
}

// Prepend 在头部插入
func (l *LinkedList[T]) Prepend(value T) {
	newNode := &Node[T]{Value: value, Next: l.Head}
	l.Head = newNode
	l.size++
}

// Size 返回链表大小
func (l *LinkedList[T]) Size() int {
	return l.size
}

// All 返回遍历所有元素的迭代器
func (l *LinkedList[T]) All() func(yield func(T) bool) {
	return func(yield func(T) bool) {
		for node := l.Head; node != nil; node = node.Next {
			if !yield(node.Value) {
				return
			}
		}
	}
}

// AllWithIndex 返回带索引的迭代器
func (l *LinkedList[T]) AllWithIndex() func(yield func(int, T) bool) {
	return func(yield func(int, T) bool) {
		i := 0
		for node := l.Head; node != nil; node = node.Next {
			if !yield(i, node.Value) {
				return
			}
			i++
		}
	}
}

// Reversed 返回反向迭代器（需要先收集所有元素）
func (l *LinkedList[T]) Reversed() func(yield func(T) bool) {
	return func(yield func(T) bool) {
		// 先收集所有元素到切片
		var values []T
		for node := l.Head; node != nil; node = node.Next {
			values = append(values, node.Value)
		}
		// 反向遍历
		for i := len(values) - 1; i >= 0; i-- {
			if !yield(values[i]) {
				return
			}
		}
	}
}

// ============ 双向链表示例 ============

// DNode 双向链表节点
type DNode[T any] struct {
	Value T
	Prev  *DNode[T]
	Next  *DNode[T]
}

// DoublyLinkedList 双向链表
type DoublyLinkedList[T any] struct {
	Head *DNode[T]
	Tail *DNode[T]
	size int
}

// NewDoublyLinkedList 创建双向链表
func NewDoublyLinkedList[T any]() *DoublyLinkedList[T] {
	return &DoublyLinkedList[T]{}
}

// Append 追加元素
func (l *DoublyLinkedList[T]) Append(value T) {
	newNode := &DNode[T]{Value: value, Prev: l.Tail}
	if l.Tail != nil {
		l.Tail.Next = newNode
	} else {
		l.Head = newNode
	}
	l.Tail = newNode
	l.size++
}

// Forward 正向迭代
func (l *DoublyLinkedList[T]) Forward() func(yield func(T) bool) {
	return func(yield func(T) bool) {
		for node := l.Head; node != nil; node = node.Next {
			if !yield(node.Value) {
				return
			}
		}
	}
}

// Backward 反向迭代（双向链表可以高效反向遍历）
func (l *DoublyLinkedList[T]) Backward() func(yield func(T) bool) {
	return func(yield func(T) bool) {
		for node := l.Tail; node != nil; node = node.Prev {
			if !yield(node.Value) {
				return
			}
		}
	}
}

func main() {
	fmt.Println("=== 链表迭代器示例 ===")
	fmt.Println()

	// 创建单向链表
	list := NewLinkedList[int]()
	list.Append(1)
	list.Append(2)
	list.Append(3)
	list.Append(4)
	list.Append(5)

	// 示例 1: 基本遍历
	fmt.Println("1. 单向链表 - 基本遍历:")
	for v := range list.All() {
		fmt.Printf("  %d\n", v)
	}
	fmt.Println()

	// 示例 2: 带索引遍历
	fmt.Println("2. 单向链表 - 带索引遍历:")
	for i, v := range list.AllWithIndex() {
		fmt.Printf("  [%d] = %d\n", i, v)
	}
	fmt.Println()

	// 示例 3: 反向遍历
	fmt.Println("3. 单向链表 - 反向遍历:")
	for v := range list.Reversed() {
		fmt.Printf("  %d\n", v)
	}
	fmt.Println()

	// 示例 4: 提前退出
	fmt.Println("4. 单向链表 - 找到 3 就停止:")
	for v := range list.All() {
		fmt.Printf("  处理: %d\n", v)
		if v == 3 {
			fmt.Println("  找到 3，停止遍历")
			break
		}
	}
	fmt.Println()

	// 示例 5: 双向链表
	fmt.Println("5. 双向链表:")
	dlist := NewDoublyLinkedList[string]()
	dlist.Append("A")
	dlist.Append("B")
	dlist.Append("C")
	dlist.Append("D")

	fmt.Println("  正向遍历:")
	for v := range dlist.Forward() {
		fmt.Printf("    %s\n", v)
	}

	fmt.Println("  反向遍历:")
	for v := range dlist.Backward() {
		fmt.Printf("    %s\n", v)
	}
}

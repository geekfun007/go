// 示例 5：二叉树遍历
// 演示如何为树结构实现不同的遍历迭代器
package main

import "fmt"

// TreeNode 二叉树节点
type TreeNode[T any] struct {
	Value T
	Left  *TreeNode[T]
	Right *TreeNode[T]
}

// NewNode 创建新节点
func NewNode[T any](value T) *TreeNode[T] {
	return &TreeNode[T]{Value: value}
}

// InOrder 中序遍历迭代器 (左 -> 根 -> 右)
// 对于 BST，这会按升序返回所有元素
func (t *TreeNode[T]) InOrder() func(yield func(T) bool) {
	return func(yield func(T) bool) {
		var traverse func(*TreeNode[T]) bool
		traverse = func(node *TreeNode[T]) bool {
			if node == nil {
				return true
			}
			if !traverse(node.Left) {
				return false
			}
			if !yield(node.Value) {
				return false
			}
			return traverse(node.Right)
		}
		traverse(t)
	}
}

// PreOrder 前序遍历迭代器 (根 -> 左 -> 右)
func (t *TreeNode[T]) PreOrder() func(yield func(T) bool) {
	return func(yield func(T) bool) {
		var traverse func(*TreeNode[T]) bool
		traverse = func(node *TreeNode[T]) bool {
			if node == nil {
				return true
			}
			if !yield(node.Value) {
				return false
			}
			if !traverse(node.Left) {
				return false
			}
			return traverse(node.Right)
		}
		traverse(t)
	}
}

// PostOrder 后序遍历迭代器 (左 -> 右 -> 根)
func (t *TreeNode[T]) PostOrder() func(yield func(T) bool) {
	return func(yield func(T) bool) {
		var traverse func(*TreeNode[T]) bool
		traverse = func(node *TreeNode[T]) bool {
			if node == nil {
				return true
			}
			if !traverse(node.Left) {
				return false
			}
			if !traverse(node.Right) {
				return false
			}
			return yield(node.Value)
		}
		traverse(t)
	}
}

// LevelOrder 层序遍历迭代器 (BFS)
func (t *TreeNode[T]) LevelOrder() func(yield func(T) bool) {
	return func(yield func(T) bool) {
		if t == nil {
			return
		}
		queue := []*TreeNode[T]{t}
		for len(queue) > 0 {
			node := queue[0]
			queue = queue[1:]
			if !yield(node.Value) {
				return
			}
			if node.Left != nil {
				queue = append(queue, node.Left)
			}
			if node.Right != nil {
				queue = append(queue, node.Right)
			}
		}
	}
}

// LevelOrderWithDepth 带层级信息的层序遍历
func (t *TreeNode[T]) LevelOrderWithDepth() func(yield func(int, T) bool) {
	return func(yield func(int, T) bool) {
		if t == nil {
			return
		}
		type nodeWithDepth struct {
			node  *TreeNode[T]
			depth int
		}
		queue := []nodeWithDepth{{t, 0}}
		for len(queue) > 0 {
			item := queue[0]
			queue = queue[1:]
			if !yield(item.depth, item.node.Value) {
				return
			}
			if item.node.Left != nil {
				queue = append(queue, nodeWithDepth{item.node.Left, item.depth + 1})
			}
			if item.node.Right != nil {
				queue = append(queue, nodeWithDepth{item.node.Right, item.depth + 1})
			}
		}
	}
}

// Leaves 只遍历叶子节点
func (t *TreeNode[T]) Leaves() func(yield func(T) bool) {
	return func(yield func(T) bool) {
		var traverse func(*TreeNode[T]) bool
		traverse = func(node *TreeNode[T]) bool {
			if node == nil {
				return true
			}
			// 叶子节点：没有子节点
			if node.Left == nil && node.Right == nil {
				return yield(node.Value)
			}
			if !traverse(node.Left) {
				return false
			}
			return traverse(node.Right)
		}
		traverse(t)
	}
}

func main() {
	fmt.Println("=== 二叉树遍历迭代器示例 ===")
	fmt.Println()

	// 构建二叉树
	//         4
	//        / \
	//       2   6
	//      / \ / \
	//     1  3 5  7
	tree := &TreeNode[int]{
		Value: 4,
		Left: &TreeNode[int]{
			Value: 2,
			Left:  &TreeNode[int]{Value: 1},
			Right: &TreeNode[int]{Value: 3},
		},
		Right: &TreeNode[int]{
			Value: 6,
			Left:  &TreeNode[int]{Value: 5},
			Right: &TreeNode[int]{Value: 7},
		},
	}

	fmt.Println("树结构:")
	fmt.Println("        4")
	fmt.Println("       / \\")
	fmt.Println("      2   6")
	fmt.Println("     / \\ / \\")
	fmt.Println("    1  3 5  7")
	fmt.Println()

	// 示例 1: 中序遍历
	fmt.Print("1. 中序遍历 (左-根-右): ")
	for v := range tree.InOrder() {
		fmt.Printf("%d ", v)
	}
	fmt.Println()

	// 示例 2: 前序遍历
	fmt.Print("2. 前序遍历 (根-左-右): ")
	for v := range tree.PreOrder() {
		fmt.Printf("%d ", v)
	}
	fmt.Println()

	// 示例 3: 后序遍历
	fmt.Print("3. 后序遍历 (左-右-根): ")
	for v := range tree.PostOrder() {
		fmt.Printf("%d ", v)
	}
	fmt.Println()

	// 示例 4: 层序遍历
	fmt.Print("4. 层序遍历 (BFS): ")
	for v := range tree.LevelOrder() {
		fmt.Printf("%d ", v)
	}
	fmt.Println()

	// 示例 5: 带层级的遍历
	fmt.Println("5. 层序遍历 (带层级):")
	currentLevel := -1
	for level, v := range tree.LevelOrderWithDepth() {
		if level != currentLevel {
			if currentLevel != -1 {
				fmt.Println()
			}
			fmt.Printf("  Level %d: ", level)
			currentLevel = level
		}
		fmt.Printf("%d ", v)
	}
	fmt.Println()
	fmt.Println()

	// 示例 6: 叶子节点
	fmt.Print("6. 叶子节点: ")
	for v := range tree.Leaves() {
		fmt.Printf("%d ", v)
	}
	fmt.Println()
	fmt.Println()

	// 示例 7: 提前终止
	fmt.Println("7. 中序遍历找到 5 就停止:")
	for v := range tree.InOrder() {
		fmt.Printf("  访问: %d\n", v)
		if v == 5 {
			fmt.Println("  找到 5，停止")
			break
		}
	}
	fmt.Println()

	// 示例 8: 收集遍历结果
	fmt.Println("8. 收集中序遍历结果到切片:")
	var result []int
	for v := range tree.InOrder() {
		result = append(result, v)
	}
	fmt.Printf("  结果: %v\n", result)
}

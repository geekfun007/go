# 递归与尾递归 / Recursion & Tail Recursion

```go
package main

import "fmt"

// 基本递归 - 阶乘 / Basic recursion - factorial
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n-1)
}

// 尾递归版本 / Tail recursive version
// 注意: Go 编译器不保证尾调用优化
// Note: Go compiler doesn't guarantee tail call optimization
func factorialTail(n, acc int) int {
    if n <= 1 {
        return acc
    }
    return factorialTail(n-1, n*acc)
}

// 斐波那契数列 - 普通递归 / Fibonacci - normal recursion
func fibonacci(n int) int {
    if n <= 1 {
        return n
    }
    return fibonacci(n-1) + fibonacci(n-2)
}

// 斐波那契 - 记忆化递归 / Fibonacci - memoized
func fibonacciMemo(n int, memo map[int]int) int {
    if n <= 1 {
        return n
    }
    if v, ok := memo[n]; ok {
        return v
    }
    memo[n] = fibonacciMemo(n-1, memo) + fibonacciMemo(n-2, memo)
    return memo[n]
}

// 斐波那契 - 迭代版本 (推荐) / Fibonacci - iterative (recommended)
func fibonacciIter(n int) int {
    if n <= 1 {
        return n
    }
    a, b := 0, 1
    for i := 2; i <= n; i++ {
        a, b = b, a+b
    }
    return b
}

// 树遍历递归示例 / Tree traversal recursion example
type TreeNode struct {
    Value int
    Left  *TreeNode
    Right *TreeNode
}

func inorderTraversal(node *TreeNode, result *[]int) {
    if node == nil {
        return
    }
    inorderTraversal(node.Left, result)
    *result = append(*result, node.Value)
    inorderTraversal(node.Right, result)
}

// 快速排序 / Quick sort
func quickSort(arr []int) []int {
    if len(arr) <= 1 {
        return arr
    }
    
    pivot := arr[0]
    var left, right []int
    
    for _, v := range arr[1:] {
        if v < pivot {
            left = append(left, v)
        } else {
            right = append(right, v)
        }
    }
    
    result := append(quickSort(left), pivot)
    return append(result, quickSort(right)...)
}

func main() {
    // 阶乘 / Factorial
    fmt.Println("factorial(5):", factorial(5))
    fmt.Println("factorialTail(5, 1):", factorialTail(5, 1))
    
    // 斐波那契 / Fibonacci
    fmt.Println("fibonacci(10):", fibonacci(10))
    
    memo := make(map[int]int)
    fmt.Println("fibonacciMemo(10):", fibonacciMemo(10, memo))
    fmt.Println("fibonacciIter(10):", fibonacciIter(10))
    
    // 树遍历 / Tree traversal
    tree := &TreeNode{
        Value: 4,
        Left: &TreeNode{
            Value: 2,
            Left:  &TreeNode{Value: 1},
            Right: &TreeNode{Value: 3},
        },
        Right: &TreeNode{
            Value: 6,
            Left:  &TreeNode{Value: 5},
            Right: &TreeNode{Value: 7},
        },
    }
    var result []int
    inorderTraversal(tree, &result)
    fmt.Println("Inorder:", result)
    
    // 快速排序 / Quick sort
    arr := []int{64, 34, 25, 12, 22, 11, 90}
    fmt.Println("QuickSort:", quickSort(arr))
}
```

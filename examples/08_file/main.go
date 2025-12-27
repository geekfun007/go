// 示例 8：文件迭代器
// 演示如何为文件操作创建迭代器
package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Lines 返回文件行迭代器
func Lines(filename string) func(yield func(int, string) bool) {
	return func(yield func(int, string) bool) {
		file, err := os.Open(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "无法打开文件 %s: %v\n", filename, err)
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		lineNum := 1
		for scanner.Scan() {
			if !yield(lineNum, scanner.Text()) {
				return
			}
			lineNum++
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "读取文件错误: %v\n", err)
		}
	}
}

// LinesOnly 只返回行内容的迭代器
func LinesOnly(filename string) func(yield func(string) bool) {
	return func(yield func(string) bool) {
		file, err := os.Open(filename)
		if err != nil {
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			if !yield(scanner.Text()) {
				return
			}
		}
	}
}

// NonEmptyLines 返回非空行迭代器
func NonEmptyLines(filename string) func(yield func(string) bool) {
	return func(yield func(string) bool) {
		for line := range LinesOnly(filename) {
			trimmed := strings.TrimSpace(line)
			if trimmed != "" {
				if !yield(trimmed) {
					return
				}
			}
		}
	}
}

// Words 返回文件中所有单词的迭代器
func Words(filename string) func(yield func(string) bool) {
	return func(yield func(string) bool) {
		for line := range LinesOnly(filename) {
			words := strings.Fields(line)
			for _, word := range words {
				if !yield(word) {
					return
				}
			}
		}
	}
}

// DirEntries 返回目录条目迭代器
func DirEntries(path string) func(yield func(os.DirEntry) bool) {
	return func(yield func(os.DirEntry) bool) {
		entries, err := os.ReadDir(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "无法读取目录 %s: %v\n", path, err)
			return
		}

		for _, entry := range entries {
			if !yield(entry) {
				return
			}
		}
	}
}

// WalkFiles 递归遍历目录中的文件
func WalkFiles(root string) func(yield func(string, os.FileInfo) bool) {
	return func(yield func(string, os.FileInfo) bool) {
		filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil // 跳过错误
			}
			if !info.IsDir() {
				if !yield(path, info) {
					return filepath.SkipAll
				}
			}
			return nil
		})
	}
}

// FilesWithExt 返回指定扩展名的文件迭代器
func FilesWithExt(root string, ext string) func(yield func(string) bool) {
	return func(yield func(string) bool) {
		for path, info := range WalkFiles(root) {
			if strings.HasSuffix(strings.ToLower(info.Name()), ext) {
				if !yield(path) {
					return
				}
			}
		}
	}
}

func main() {
	fmt.Println("=== 文件迭代器示例 ===")
	fmt.Println()

	// 创建示例文件
	testFile := "test_data.txt"
	content := `Go Range Iterator Example
这是第二行
Third line with some words

第五行是中文

Line 7: The last line
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		fmt.Println("创建测试文件失败:", err)
		return
	}
	defer os.Remove(testFile)

	// 示例 1: 遍历文件行
	fmt.Println("1. 遍历文件行 (带行号):")
	for lineNum, line := range Lines(testFile) {
		fmt.Printf("  %3d | %s\n", lineNum, line)
	}
	fmt.Println()

	// 示例 2: 只遍历非空行
	fmt.Println("2. 只遍历非空行:")
	for line := range NonEmptyLines(testFile) {
		fmt.Printf("  > %s\n", line)
	}
	fmt.Println()

	// 示例 3: 遍历单词
	fmt.Println("3. 遍历所有单词 (前 10 个):")
	count := 0
	for word := range Words(testFile) {
		fmt.Printf("  '%s'\n", word)
		count++
		if count >= 10 {
			fmt.Println("  ...")
			break
		}
	}
	fmt.Println()

	// 示例 4: 统计行数
	fmt.Println("4. 统计信息:")
	lineCount := 0
	for _ = range LinesOnly(testFile) {
		lineCount++
	}
	fmt.Printf("  总行数: %d\n", lineCount)

	wordCount := 0
	for _ = range Words(testFile) {
		wordCount++
	}
	fmt.Printf("  总单词数: %d\n", wordCount)
	fmt.Println()

	// 示例 5: 遍历当前目录
	fmt.Println("5. 当前目录条目:")
	for entry := range DirEntries(".") {
		typeStr := "文件"
		if entry.IsDir() {
			typeStr = "目录"
		}
		fmt.Printf("  [%s] %s\n", typeStr, entry.Name())
	}
	fmt.Println()

	// 示例 6: 查找特定文件
	fmt.Println("6. 查找当前目录下的 .go 文件:")
	for path := range FilesWithExt(".", ".go") {
		fmt.Printf("  %s\n", path)
	}
	fmt.Println()

	// 示例 7: 搜索包含特定内容的行
	fmt.Println("7. 搜索包含 'line' 的行 (忽略大小写):")
	for lineNum, line := range Lines(testFile) {
		if strings.Contains(strings.ToLower(line), "line") {
			fmt.Printf("  行 %d: %s\n", lineNum, line)
		}
	}
	fmt.Println()

	// 示例 8: 组合使用 - 统计每行字符数
	fmt.Println("8. 每行字符数统计:")
	for lineNum, line := range Lines(testFile) {
		fmt.Printf("  行 %d: %d 字符\n", lineNum, len(line))
	}
}

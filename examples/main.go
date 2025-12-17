package main

import (
	"fmt"
	"os"
	"os/exec"
)

// main.go - 运行所有示例的主程序

func main() {
	examples := []struct {
		name string
		file string
		desc string
	}{
		{"基础操作", "01_basic_operations.go", "Slice 的声明、初始化、访问和修改"},
		{"追加和复制", "02_append_copy.go", "append 和 copy 函数的详细用法"},
		{"切片操作", "03_slicing.go", "切片表达式和子切片操作"},
		{"增删改查", "04_crud_operations.go", "Slice 的 CRUD 操作实现"},
		{"高级技巧", "05_advanced_techniques.go", "反转、去重、过滤、映射等高级操作"},
		{"性能优化", "06_performance.go", "性能优化技巧和对比"},
		{"数据处理", "07_data_processing.go", "实际数据处理场景"},
		{"内存管理", "08_memory_management.go", "内存管理和优化"},
	}

	fmt.Println("╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║         Go Slice 操作方法详解 - 实战示例集合                  ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")
	fmt.Println()

	if len(os.Args) > 1 {
		// 运行指定的示例
		runExample(os.Args[1])
		return
	}

	// 显示菜单
	fmt.Println("可用的示例:")
	for i, ex := range examples {
		fmt.Printf("  %d. %s\n", i+1, ex.name)
		fmt.Printf("     %s\n", ex.desc)
		fmt.Println()
	}

	fmt.Println("运行方式:")
	fmt.Println("  go run main.go [示例编号]")
	fmt.Println("  例如: go run main.go 1")
	fmt.Println()
	fmt.Println("或者直接运行单个示例:")
	for i, ex := range examples {
		fmt.Printf("  go run %s   # %d. %s\n", ex.file, i+1, ex.name)
	}
	fmt.Println()

	fmt.Print("请选择要运行的示例 (1-8, 0=全部运行): ")
	var choice int
	fmt.Scanln(&choice)

	if choice == 0 {
		// 运行所有示例
		for i, ex := range examples {
			fmt.Printf("\n\n")
			fmt.Println("═══════════════════════════════════════════════════════════════")
			fmt.Printf("示例 %d: %s\n", i+1, ex.name)
			fmt.Println("═══════════════════════════════════════════════════════════════")
			runExample(ex.file)
			fmt.Println("\n按 Enter 继续...")
			fmt.Scanln()
		}
	} else if choice >= 1 && choice <= len(examples) {
		runExample(examples[choice-1].file)
	} else {
		fmt.Println("无效的选择！")
	}
}

func runExample(filename string) {
	cmd := exec.Command("go", "run", filename)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = "."

	if err := cmd.Run(); err != nil {
		fmt.Printf("运行示例失败: %v\n", err)
	}
}

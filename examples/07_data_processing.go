package main

import (
	"fmt"
	"sort"
	"strings"
)

// 07_data_processing.go - Slice 实际数据处理场景

type Student struct {
	ID    int
	Name  string
	Score int
	Grade string
}

type Product struct {
	ID       int
	Name     string
	Price    float64
	Category string
}

func main() {
	fmt.Println("=== Slice 实际数据处理场景 ===\n")

	// 1. 学生成绩处理
	processStudentScores()

	// 2. 商品数据分析
	processProductData()

	// 3. 日志数据处理
	processLogData()

	// 4. 数据统计
	demonstrateStatistics()
}

func processStudentScores() {
	fmt.Println("--- 1. 学生成绩处理 ---")

	students := []Student{
		{1, "Alice", 85, ""},
		{2, "Bob", 92, ""},
		{3, "Charlie", 78, ""},
		{4, "David", 95, ""},
		{5, "Eve", 88, ""},
		{6, "Frank", 73, ""},
		{7, "Grace", 90, ""},
		{8, "Henry", 67, ""},
	}

	fmt.Println("原始数据:")
	printStudents(students)

	// 1.1 计算等级
	for i := range students {
		students[i].Grade = calculateGrade(students[i].Score)
	}

	fmt.Println("\n添加等级后:")
	printStudents(students)

	// 1.2 按成绩排序
	sort.Slice(students, func(i, j int) bool {
		return students[i].Score > students[j].Score
	})

	fmt.Println("\n按成绩排序（降序）:")
	printStudents(students)

	// 1.3 筛选优秀学生（成绩 >= 90）
	excellent := filterStudents(students, func(s Student) bool {
		return s.Score >= 90
	})

	fmt.Println("\n优秀学生（成绩 >= 90）:")
	printStudents(excellent)

	// 1.4 计算平均分
	totalScore := 0
	for _, s := range students {
		totalScore += s.Score
	}
	avgScore := float64(totalScore) / float64(len(students))
	fmt.Printf("\n平均分: %.2f\n", avgScore)

	// 1.5 按等级分组
	gradeGroups := groupStudentsByGrade(students)
	fmt.Println("\n按等级分组:")
	for grade, group := range gradeGroups {
		fmt.Printf("  %s 级 (%d 人): ", grade, len(group))
		names := make([]string, len(group))
		for i, s := range group {
			names[i] = s.Name
		}
		fmt.Printf("%s\n", strings.Join(names, ", "))
	}

	fmt.Println()
}

func processProductData() {
	fmt.Println("--- 2. 商品数据分析 ---")

	products := []Product{
		{1, "iPhone 15", 5999.0, "电子产品"},
		{2, "MacBook Pro", 12999.0, "电子产品"},
		{3, "AirPods", 1299.0, "电子产品"},
		{4, "Nike 鞋", 699.0, "服装"},
		{5, "Adidas 外套", 899.0, "服装"},
		{6, "ThinkPad", 6999.0, "电子产品"},
		{7, "李宁运动裤", 399.0, "服装"},
		{8, "iPad", 3299.0, "电子产品"},
	}

	fmt.Println("商品列表:")
	printProducts(products)

	// 2.1 按价格排序
	sort.Slice(products, func(i, j int) bool {
		return products[i].Price < products[j].Price
	})

	fmt.Println("\n按价格排序（升序）:")
	printProducts(products)

	// 2.2 按类别分组
	categoryGroups := groupProductsByCategory(products)
	fmt.Println("\n按类别分组:")
	for category, items := range categoryGroups {
		fmt.Printf("\n  %s (%d 件):\n", category, len(items))
		for _, p := range items {
			fmt.Printf("    - %s: ¥%.2f\n", p.Name, p.Price)
		}
	}

	// 2.3 价格区间统计
	fmt.Println("\n价格区间统计:")
	priceRanges := []struct {
		Name  string
		Min   float64
		Max   float64
		Count int
	}{
		{"低价 (<1000)", 0, 1000, 0},
		{"中价 (1000-5000)", 1000, 5000, 0},
		{"高价 (5000-10000)", 5000, 10000, 0},
		{"奢侈 (>10000)", 10000, 999999, 0},
	}

	for _, p := range products {
		for i := range priceRanges {
			if p.Price >= priceRanges[i].Min && p.Price < priceRanges[i].Max {
				priceRanges[i].Count++
			}
		}
	}

	for _, r := range priceRanges {
		fmt.Printf("  %s: %d 件\n", r.Name, r.Count)
	}

	// 2.4 计算各类别平均价格
	fmt.Println("\n各类别平均价格:")
	for category, items := range categoryGroups {
		total := 0.0
		for _, p := range items {
			total += p.Price
		}
		avg := total / float64(len(items))
		fmt.Printf("  %s: ¥%.2f\n", category, avg)
	}

	// 2.5 找出最贵和最便宜的商品
	minProduct := products[0]
	maxProduct := products[0]
	for _, p := range products {
		if p.Price < minProduct.Price {
			minProduct = p
		}
		if p.Price > maxProduct.Price {
			maxProduct = p
		}
	}

	fmt.Printf("\n最便宜: %s (¥%.2f)\n", minProduct.Name, minProduct.Price)
	fmt.Printf("最贵: %s (¥%.2f)\n", maxProduct.Name, maxProduct.Price)

	fmt.Println()
}

func processLogData() {
	fmt.Println("--- 3. 日志数据处理 ---")

	logs := []string{
		"2024-01-15 10:23:45 [INFO] Application started",
		"2024-01-15 10:24:12 [ERROR] Database connection failed",
		"2024-01-15 10:24:15 [WARN] Retrying connection...",
		"2024-01-15 10:24:18 [INFO] Database connected",
		"2024-01-15 10:25:03 [ERROR] API request timeout",
		"2024-01-15 10:26:30 [INFO] Request processed successfully",
		"2024-01-15 10:27:45 [ERROR] File not found",
		"2024-01-15 10:28:12 [INFO] Cache cleared",
	}

	fmt.Println("所有日志:")
	for _, log := range logs {
		fmt.Printf("  %s\n", log)
	}

	// 3.1 筛选错误日志
	errorLogs := filterLogs(logs, func(log string) bool {
		return strings.Contains(log, "[ERROR]")
	})

	fmt.Println("\nERROR 日志:")
	for _, log := range errorLogs {
		fmt.Printf("  %s\n", log)
	}

	// 3.2 按日志级别分类
	logsByLevel := make(map[string][]string)
	levels := []string{"INFO", "WARN", "ERROR"}

	for _, level := range levels {
		logsByLevel[level] = filterLogs(logs, func(log string) bool {
			return strings.Contains(log, "["+level+"]")
		})
	}

	fmt.Println("\n按级别统计:")
	for _, level := range levels {
		fmt.Printf("  [%s]: %d 条\n", level, len(logsByLevel[level]))
	}

	// 3.3 提取错误消息
	fmt.Println("\n错误消息列表:")
	for _, log := range errorLogs {
		parts := strings.SplitN(log, "[ERROR] ", 2)
		if len(parts) == 2 {
			fmt.Printf("  - %s\n", parts[1])
		}
	}

	fmt.Println()
}

func demonstrateStatistics() {
	fmt.Println("--- 4. 数据统计 ---")

	data := []float64{23.5, 18.2, 45.6, 33.1, 28.9, 41.3, 19.7, 36.4, 25.8, 39.2}

	fmt.Printf("数据集: %v\n", data)

	// 4.1 基本统计
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	mean := sum / float64(len(data))

	min, max := data[0], data[0]
	for _, v := range data {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}

	fmt.Printf("\n基本统计:\n")
	fmt.Printf("  数量: %d\n", len(data))
	fmt.Printf("  总和: %.2f\n", sum)
	fmt.Printf("  平均值: %.2f\n", mean)
	fmt.Printf("  最小值: %.2f\n", min)
	fmt.Printf("  最大值: %.2f\n", max)
	fmt.Printf("  范围: %.2f\n", max-min)

	// 4.2 中位数
	sortedData := make([]float64, len(data))
	copy(sortedData, data)
	sort.Float64s(sortedData)

	var median float64
	n := len(sortedData)
	if n%2 == 0 {
		median = (sortedData[n/2-1] + sortedData[n/2]) / 2
	} else {
		median = sortedData[n/2]
	}

	fmt.Printf("  中位数: %.2f\n", median)

	// 4.3 方差
	variance := 0.0
	for _, v := range data {
		diff := v - mean
		variance += diff * diff
	}
	variance /= float64(len(data))

	fmt.Printf("  方差: %.2f\n", variance)

	// 4.4 分位数
	q1 := sortedData[len(sortedData)/4]
	q3 := sortedData[3*len(sortedData)/4]

	fmt.Printf("  第一四分位数 (Q1): %.2f\n", q1)
	fmt.Printf("  第三四分位数 (Q3): %.2f\n", q3)
	fmt.Printf("  四分位距 (IQR): %.2f\n", q3-q1)

	fmt.Println()
}

// === 辅助函数 ===

func calculateGrade(score int) string {
	switch {
	case score >= 90:
		return "A"
	case score >= 80:
		return "B"
	case score >= 70:
		return "C"
	case score >= 60:
		return "D"
	default:
		return "F"
	}
}

func filterStudents(students []Student, fn func(Student) bool) []Student {
	result := make([]Student, 0)
	for _, s := range students {
		if fn(s) {
			result = append(result, s)
		}
	}
	return result
}

func groupStudentsByGrade(students []Student) map[string][]Student {
	result := make(map[string][]Student)
	for _, s := range students {
		result[s.Grade] = append(result[s.Grade], s)
	}
	return result
}

func groupProductsByCategory(products []Product) map[string][]Product {
	result := make(map[string][]Product)
	for _, p := range products {
		result[p.Category] = append(result[p.Category], p)
	}
	return result
}

func filterLogs(logs []string, fn func(string) bool) []string {
	result := make([]string, 0)
	for _, log := range logs {
		if fn(log) {
			result = append(result, log)
		}
	}
	return result
}

func printStudents(students []Student) {
	for _, s := range students {
		if s.Grade != "" {
			fmt.Printf("  [%d] %s: %d 分 (等级: %s)\n", s.ID, s.Name, s.Score, s.Grade)
		} else {
			fmt.Printf("  [%d] %s: %d 分\n", s.ID, s.Name, s.Score)
		}
	}
}

func printProducts(products []Product) {
	for _, p := range products {
		fmt.Printf("  [%d] %s - ¥%.2f (%s)\n", p.ID, p.Name, p.Price, p.Category)
	}
}

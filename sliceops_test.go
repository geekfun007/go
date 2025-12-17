package sliceops

import (
	"testing"
)

// sliceops_test.go - Slice 操作函数的单元测试和基准测试

// === 单元测试 ===

func TestInsertAt(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		index    int
		value    int
		expected []int
	}{
		{"在开头插入", []int{2, 3, 4}, 0, 1, []int{1, 2, 3, 4}},
		{"在中间插入", []int{1, 3, 4}, 1, 2, []int{1, 2, 3, 4}},
		{"在末尾插入", []int{1, 2, 3}, 3, 4, []int{1, 2, 3, 4}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := InsertAt(tt.input, tt.index, tt.value)
			if !sliceEqual(result, tt.expected) {
				t.Errorf("got %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestRemoveAt(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		index    int
		expected []int
	}{
		{"删除第一个", []int{1, 2, 3, 4}, 0, []int{2, 3, 4}},
		{"删除中间", []int{1, 2, 3, 4}, 2, []int{1, 2, 4}},
		{"删除最后", []int{1, 2, 3, 4}, 3, []int{1, 2, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RemoveAt(tt.input, tt.index)
			if !sliceEqual(result, tt.expected) {
				t.Errorf("got %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestUniqueOrdered(t *testing.T) {
	input := []int{1, 2, 2, 3, 3, 3, 4, 5, 5, 1}
	expected := []int{1, 2, 3, 4, 5}
	result := UniqueOrdered(input)

	if !sliceEqual(result, expected) {
		t.Errorf("got %v, want %v", result, expected)
	}
}

func TestFilter(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	expected := []int{2, 4, 6, 8, 10}
	result := Filter(input, func(n int) bool { return n%2 == 0 })

	if !sliceEqual(result, expected) {
		t.Errorf("got %v, want %v", result, expected)
	}
}

func TestMap(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	expected := []int{1, 4, 9, 16, 25}
	result := Map(input, func(n int) int { return n * n })

	if !sliceEqual(result, expected) {
		t.Errorf("got %v, want %v", result, expected)
	}
}

func TestReduce(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	expected := 15
	result := Reduce(input, 0, func(acc, n int) int { return acc + n })

	if result != expected {
		t.Errorf("got %d, want %d", result, expected)
	}
}

func TestReverse(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	expected := []int{5, 4, 3, 2, 1}
	Reverse(input)

	if !sliceEqual(input, expected) {
		t.Errorf("got %v, want %v", input, expected)
	}
}

func TestContains(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	
	if !Contains(s, 3) {
		t.Error("should contain 3")
	}
	
	if Contains(s, 10) {
		t.Error("should not contain 10")
	}
}

func TestIndexOf(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	
	if idx := IndexOf(s, 3); idx != 2 {
		t.Errorf("expected index 2, got %d", idx)
	}
	
	if idx := IndexOf(s, 10); idx != -1 {
		t.Errorf("expected index -1, got %d", idx)
	}
}

func TestFindMinMax(t *testing.T) {
	s := []int{3, 7, 2, 9, 1, 5, 8}
	min, max := FindMinMax(s)
	
	if min != 1 {
		t.Errorf("expected min 1, got %d", min)
	}
	
	if max != 9 {
		t.Errorf("expected max 9, got %d", max)
	}
}

// === 基准测试 ===

func BenchmarkAppendWithoutPrealloc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := []int{}
		for j := 0; j < 1000; j++ {
			s = append(s, j)
		}
	}
}

func BenchmarkAppendWithPrealloc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := make([]int, 0, 1000)
		for j := 0; j < 1000; j++ {
			s = append(s, j)
		}
	}
}

func BenchmarkAppendWithPreallocAndAssign(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := make([]int, 1000)
		for j := 0; j < 1000; j++ {
			s[j] = j
		}
	}
}

func BenchmarkCopySmallSlice(b *testing.B) {
	src := make([]int, 10)
	for i := 0; i < b.N; i++ {
		dst := make([]int, len(src))
		copy(dst, src)
	}
}

func BenchmarkCopyLargeSlice(b *testing.B) {
	src := make([]int, 10000)
	for i := 0; i < b.N; i++ {
		dst := make([]int, len(src))
		copy(dst, src)
	}
}

func BenchmarkFilterNewSlice(b *testing.B) {
	data := make([]int, 1000)
	for i := range data {
		data[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Filter(data, func(x int) bool { return x%2 == 0 })
	}
}

func BenchmarkFilterInPlace(b *testing.B) {
	for i := 0; i < b.N; i++ {
		data := make([]int, 1000)
		for j := range data {
			data[j] = j
		}
		_ = FilterInPlace(data, func(x int) bool { return x%2 == 0 })
	}
}

func BenchmarkRemoveAtOrdered(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		_ = RemoveAt(s, 5)
	}
}

func BenchmarkRemoveAtFast(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		_ = RemoveAtFast(s, 5)
	}
}

func BenchmarkUniqueOrdered(b *testing.B) {
	data := make([]int, 1000)
	for i := range data {
		data[i] = i % 100 // 有重复
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = UniqueOrdered(data)
	}
}

func BenchmarkUniqueFast(b *testing.B) {
	data := make([]int, 1000)
	for i := range data {
		data[i] = i % 100 // 有重复
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = UniqueFast(data)
	}
}

func BenchmarkReverse(b *testing.B) {
	data := make([]int, 1000)
	for i := range data {
		data[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 每次都需要复制，因为 reverse 会修改原数组
		s := make([]int, len(data))
		copy(s, data)
		Reverse(s)
	}
}

// === 辅助函数 ===

func sliceEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

package main

import (
	"fmt"
	"sync"
	"time"
)

// 实战示例: 并发数据处理
// Practical Example: Concurrent Data Processing

// 示例1: 批量数据处理
func batchProcessing() {
	fmt.Println("\n=== 示例1: 批量数据处理 ===")
	
	// 模拟大量数据
	data := make([]int, 100)
	for i := 0; i < 100; i++ {
		data[i] = i + 1
	}
	
	const batchSize = 10
	const numWorkers = 5
	
	batches := make(chan []int, numWorkers)
	results := make(chan int, len(data))
	var wg sync.WaitGroup
	
	// 启动 workers
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for batch := range batches {
				sum := 0
				for _, val := range batch {
					sum += val
				}
				fmt.Printf("Worker %d 处理批次，和为: %d\n", id, sum)
				results <- sum
			}
		}(w)
	}
	
	// 分批发送数据
	go func() {
		for i := 0; i < len(data); i += batchSize {
			end := i + batchSize
			if end > len(data) {
				end = len(data)
			}
			batches <- data[i:end]
		}
		close(batches)
	}()
	
	// 等待所有 worker 完成
	go func() {
		wg.Wait()
		close(results)
	}()
	
	// 汇总结果
	total := 0
	for result := range results {
		total += result
	}
	fmt.Printf("总和: %d\n", total)
}

// 示例2: 并行 Map-Reduce
func mapReduceExample() {
	fmt.Println("\n=== 示例2: Map-Reduce 模式 ===")
	
	// 输入数据
	documents := []string{
		"hello world",
		"hello golang",
		"world of go",
		"go go go",
	}
	
	// Map 阶段: 计算每个文档的词频
	type WordCount map[string]int
	
	mapResults := make(chan WordCount, len(documents))
	var mapWg sync.WaitGroup
	
	for _, doc := range documents {
		mapWg.Add(1)
		go func(text string) {
			defer mapWg.Done()
			
			wc := make(WordCount)
			words := []string{}
			word := ""
			for _, char := range text + " " {
				if char == ' ' {
					if word != "" {
						words = append(words, word)
						word = ""
					}
				} else {
					word += string(char)
				}
			}
			
			for _, w := range words {
				wc[w]++
			}
			
			mapResults <- wc
		}(doc)
	}
	
	go func() {
		mapWg.Wait()
		close(mapResults)
	}()
	
	// Reduce 阶段: 合并所有结果
	finalCount := make(WordCount)
	for wc := range mapResults {
		for word, count := range wc {
			finalCount[word] += count
		}
	}
	
	fmt.Println("词频统计:")
	for word, count := range finalCount {
		fmt.Printf("  %s: %d\n", word, count)
	}
}

// 示例3: 数据流处理管道
func streamProcessingPipeline() {
	fmt.Println("\n=== 示例3: 数据流处理管道 ===")
	
	// 阶段1: 数据生成器
	generate := func(nums ...int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for _, n := range nums {
				out <- n
			}
		}()
		return out
	}
	
	// 阶段2: 过滤器
	filter := func(in <-chan int, predicate func(int) bool) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for n := range in {
				if predicate(n) {
					out <- n
				}
			}
		}()
		return out
	}
	
	// 阶段3: 转换器
	transform := func(in <-chan int, fn func(int) int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for n := range in {
				out <- fn(n)
			}
		}()
		return out
	}
	
	// 构建管道
	numbers := generate(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	evens := filter(numbers, func(n int) bool { return n%2 == 0 })
	squared := transform(evens, func(n int) int { return n * n })
	
	// 消费结果
	fmt.Println("偶数的平方:")
	for result := range squared {
		fmt.Printf("  %d\n", result)
	}
}

// 示例4: 并发缓存
type Cache struct {
	mu    sync.RWMutex
	data  map[string]interface{}
	ttl   time.Duration
	timer map[string]*time.Timer
}

func NewCache(ttl time.Duration) *Cache {
	return &Cache{
		data:  make(map[string]interface{}),
		ttl:   ttl,
		timer: make(map[string]*time.Timer),
	}
}

func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	// 取消旧的定时器
	if timer, exists := c.timer[key]; exists {
		timer.Stop()
	}
	
	c.data[key] = value
	
	// 设置新的过期定时器
	c.timer[key] = time.AfterFunc(c.ttl, func() {
		c.Delete(key)
	})
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	value, exists := c.data[key]
	return value, exists
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	delete(c.data, key)
	if timer, exists := c.timer[key]; exists {
		timer.Stop()
		delete(c.timer, key)
	}
}

func cacheExample() {
	fmt.Println("\n=== 示例4: 并发缓存 ===")
	
	cache := NewCache(2 * time.Second)
	
	// 并发写入
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", id)
			cache.Set(key, id*100)
			fmt.Printf("设置 %s = %d\n", key, id*100)
		}(i)
	}
	
	wg.Wait()
	time.Sleep(500 * time.Millisecond)
	
	// 并发读取
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", id)
			if val, ok := cache.Get(key); ok {
				fmt.Printf("获取 %s = %v\n", key, val)
			}
		}(i)
	}
	
	wg.Wait()
	
	// 等待过期
	fmt.Println("等待缓存过期...")
	time.Sleep(3 * time.Second)
	
	if val, ok := cache.Get("key0"); ok {
		fmt.Printf("key0 = %v\n", val)
	} else {
		fmt.Println("key0 已过期")
	}
}

// 示例5: 并发下载器
func concurrentDownloader() {
	fmt.Println("\n=== 示例5: 并发下载器 ===")
	
	urls := []string{
		"http://example.com/file1.txt",
		"http://example.com/file2.txt",
		"http://example.com/file3.txt",
		"http://example.com/file4.txt",
		"http://example.com/file5.txt",
	}
	
	type Result struct {
		URL   string
		Size  int
		Error error
	}
	
	download := func(url string) Result {
		// 模拟下载
		time.Sleep(time.Duration(len(url)%3+1) * time.Second)
		return Result{
			URL:  url,
			Size: len(url) * 1024,
		}
	}
	
	results := make(chan Result, len(urls))
	var wg sync.WaitGroup
	
	// 限制并发数为 3
	semaphore := make(chan struct{}, 3)
	
	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			
			semaphore <- struct{}{}        // 获取信号量
			defer func() { <-semaphore }() // 释放信号量
			
			fmt.Printf("开始下载: %s\n", u)
			result := download(u)
			results <- result
			fmt.Printf("完成下载: %s (%d bytes)\n", result.URL, result.Size)
		}(url)
	}
	
	go func() {
		wg.Wait()
		close(results)
	}()
	
	// 统计结果
	totalSize := 0
	for result := range results {
		if result.Error == nil {
			totalSize += result.Size
		}
	}
	
	fmt.Printf("总下载大小: %d bytes\n", totalSize)
}

// 示例6: 生产者-消费者模式
func producerConsumerPattern() {
	fmt.Println("\n=== 示例6: 生产者-消费者模式 ===")
	
	const bufferSize = 5
	queue := make(chan int, bufferSize)
	var wg sync.WaitGroup
	
	// 生产者
	producer := func(id int, count int) {
		defer wg.Done()
		for i := 0; i < count; i++ {
			item := id*100 + i
			fmt.Printf("生产者 %d 生产: %d\n", id, item)
			queue <- item
			time.Sleep(100 * time.Millisecond)
		}
	}
	
	// 消费者
	consumer := func(id int) {
		defer wg.Done()
		for item := range queue {
			fmt.Printf("消费者 %d 消费: %d\n", id, item)
			time.Sleep(150 * time.Millisecond)
		}
	}
	
	// 启动 2 个生产者
	wg.Add(2)
	go producer(1, 5)
	go producer(2, 5)
	
	// 启动 3 个消费者
	wg.Add(3)
	go consumer(1)
	go consumer(2)
	go consumer(3)
	
	// 等待生产者完成
	time.Sleep(2 * time.Second)
	close(queue)
	
	wg.Wait()
	fmt.Println("生产消费完成")
}

func main() {
	fmt.Println("Go 异步编程 - 数据处理实战")
	fmt.Println("============================")
	
	batchProcessing()
	mapReduceExample()
	streamProcessingPipeline()
	cacheExample()
	concurrentDownloader()
	producerConsumerPattern()
	
	fmt.Println("\n所有示例完成！")
}

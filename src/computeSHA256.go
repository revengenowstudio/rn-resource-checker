package main

import (
	"crypto/sha256"
	"fmt"
	"io"

	"os"
	"sync"
)

// 定义一个结构体，用于传递文件路径及其哈希结果
type fileResult struct {
	filename string
	hash     []byte
	err      error
}

// 计算单个文件的 SHA-256，使用 io.SectionReader 进行分段读取
func computeSHA256ForFileWithSectionReader(filename string, numWorkers int) ([]byte, error) {

	// 打开文件
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// 获取文件大小
	info, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}
	fileSize := info.Size()

	// 计算每个 sha256Worker 处理的文件段大小
	segmentSize := fileSize / int64(numWorkers)
	if fileSize%int64(numWorkers) != 0 {
		segmentSize++
	}

	// 创建一个新的 SHA-256 哈希对象
	hasher := sha256.New()

	// 创建一个 WaitGroup，用于等待所有 goroutine 完成
	var wg sync.WaitGroup

	// 启动多个 sha256Worker goroutine，负责读取文件段并更新哈希
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			// 计算当前 sha256Worker 处理的文件段的起始和结束位置
			start := int64(i) * segmentSize
			end := start + segmentSize
			if end > fileSize {
				end = fileSize
			}

			// 创建一个 SectionReader，用于读取文件段
			section := io.NewSectionReader(file, start, end-start)

			// 使用 io.Copy 将文件段内容写入哈希对象
			_, err := io.Copy(hasher, section)
			if err != nil {
				fmt.Printf("Error reading file segment: %v\n", err)
				return
			}
		}(i)
	}

	// 等待所有 sha256Worker goroutine 完成
	wg.Wait()

	// 获取最终的哈希值
	return hasher.Sum(nil), nil
}

// 计算单个文件的 SHA-256，使用 io.SectionReader 进行分段读取
func computeSHA256ForFile(filename string, resultChan chan<- fileResult,wg *sync.WaitGroup) {
	defer wg.Done()
	file, err := os.Open(filename)
	if err != nil {
		resultChan <- fileResult{
			filename: filename,
			hash:     nil,
			err:      fmt.Errorf("failed to open file: %w", err),
		}
	}
	defer file.Close()

	hasher := sha256.New()
	buffer := make([]byte, 4096) // 使用4KB的缓冲区

	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			resultChan <- fileResult{
				filename: filename,
				hash:     nil,
				err:      fmt.Errorf("failed to read file: %w", err),
			}

		}

		hasher.Write(buffer[:n])
	}
	hash := hasher.Sum(nil)
	resultChan <- fileResult{
		filename: filename,
		hash:     hash,
		err:      nil,
	}

}

// sha256Worker 函数，负责从任务队列中获取文件并计算其 SHA-256
func sha256Worker(taskChan <-chan string, resultChan chan<- fileResult, wg *sync.WaitGroup, numWorkers int) {
	defer wg.Done()

	for filename := range taskChan {
		hash, err := computeSHA256ForFileWithSectionReader(filename, numWorkers)
		resultChan <- fileResult{filename: filename, hash: hash, err: err}
	}
}

// 并发计算多个文件的 SHA-256
func ComputeSHA256ForMultipleFiles(filenames []string, numWorkers int) ([]fileResult, error) {

	results := make([]fileResult, 0, len(filenames))
	resultChan := make(chan fileResult, len(filenames))
	var wg sync.WaitGroup
	for _, filename := range filenames {
		wg.Add(1)
		go computeSHA256ForFile(filename, resultChan,&wg)
		// hash, err := computeSHA256ForFile(filename)
	}
	go func() {
		wg.Wait()
		close(resultChan)
	}()
	for result := range resultChan {
		fmt.Printf("File: %s, Hash: %x\n", result.filename, result.hash)
		results = append(results, result)
	}
	return results, nil
}

// 收集所有文件的哈希结果

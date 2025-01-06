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
func computeSHA256ForFile(filename string, resultChan chan<- fileResult, wg *sync.WaitGroup) {
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

// 并发计算多个文件的 SHA-256
func ComputeSHA256ForMultipleFiles(filenames []string, numWorkers int) ([]fileResult, error) {

	results := make([]fileResult, 0, len(filenames))
	resultChan := make(chan fileResult, len(filenames))
	var wg sync.WaitGroup
	for _, filename := range filenames {
		wg.Add(1)
		go computeSHA256ForFile(filename, resultChan, &wg)
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

package utils

import (
	"crypto/sha256"
	"fmt"
	"io"
	"rn-resource-checker/src/log"

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
		resultChan <- fileResult{filename: filename, hash: nil, err: fmt.Errorf("open: %w", err)}
		return // 关键
	}
	defer file.Close()

	hash := sha256.New()
	buf := make([]byte, 4*1024)

	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			resultChan <- fileResult{filename: filename, hash: nil, err: fmt.Errorf("read: %w", err)}
			return // 关键
		}
		if n > 0 {
			hash.Write(buf[:n])
		}
		if err == io.EOF {
			break
		}
	}

	result := hash.Sum(nil)

	resultChan <- fileResult{filename: filename, hash: result, err: nil}
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
		log.Info(fmt.Sprintf("File: %s, Hash: %x", result.filename, result.hash))
		results = append(results, result)
	}
	return results, nil
}

// 收集所有文件的哈希结果

package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// 定义一个结构体，用于传递文件路径
type filePath struct {
	path string
	err  error
}

// walkDirWorker 函数，负责从任务队列中获取目录并遍历其下的文件
func walkDirWorker(taskChan <-chan string, resultChan chan<- filePath, wg *sync.WaitGroup, whiteListSuffix *[]string) {
	defer wg.Done()

	for dir := range taskChan {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				resultChan <- filePath{path: path, err: err}
				return nil // 继续遍历其他文件
			}
			// if *whiteListSuffix != nil {
			// 	for _, suffix := range *whiteListSuffix {
			// 		if filepath.Ext(path) == suffix {
			// 			resultChan <- filePath{path: path, err: nil}
			// 			return nil
			// 		}
			// 	}
			// }
			if !info.IsDir() {
				if *whiteListSuffix != nil {
					for _, suffix := range *whiteListSuffix {
						if filepath.Ext(path) == "."+suffix {
							resultChan <- filePath{path: path, err: nil}
							return nil
						}
					}
				} else {
					resultChan <- filePath{path: path, err: nil}
				}

			}
			return nil
		})

		if err != nil {
			resultChan <- filePath{path: dir, err: err}
		}
	}
}

// 并发获取文件夹及其子文件夹中的文件名和路径
func WalkDirectoriesConcurrently(
	rootDirs []string,
	numWorkers int,
	whiteListSuffix *[]string) ([]string, []error) {
	// 创建任务队列和结果队列
	taskChan := make(chan string, len(rootDirs))
	resultChan := make(chan filePath, 1000) // 缓冲区大小可以根据需要调整

	// 创建一个 WaitGroup，用于等待所有 walkDirWorker goroutine 完成
	var wg sync.WaitGroup

	// 启动多个 walkDirWorker goroutine
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go walkDirWorker(taskChan, resultChan, &wg, whiteListSuffix)
	}

	// 将所有根目录发送到任务队列
	for _, dir := range rootDirs {
		taskChan <- dir
	}
	close(taskChan)

	// 等待所有 walkDirWorker goroutine 完成
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// 收集所有文件的路径和错误
	var files []string
	var errors []error

	for result := range resultChan {
		if result.err != nil {
			errors = append(errors, fmt.Errorf("error processing %s: %v", result.path, result.err))
		} else {
			files = append(files, result.path)
		}
	}

	return files, errors
}

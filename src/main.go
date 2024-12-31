package main

import (
	"fmt"
	"time"
)

func main() {

	// 记录程序开始时间
	startTime := time.Now()

	// 指定要遍历的根目录列表
	targetDirs := []string{"./"}

	whiteListSuffix := []string{"mix", "exe", "dll"}

	// 指定并发 worker 的数量
	numWorkers := 8

	// 并发获取文件夹及其子文件夹中的文件名和路径
	files, errs := WalkDirectoriesConcurrently(targetDirs, numWorkers, &whiteListSuffix)

	// 打印遇到的错误
	if len(errs) > 0 {
		fmt.Println("Encountered errors:")
		for _, err := range errs {
			fmt.Println(err)
		}
	}

	results, err := ComputeSHA256ForMultipleFiles(files, numWorkers)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	versionCodeContent, err := ReadFileToString("Versioncode")
	if err != nil {
		fmt.Printf("load Version file failed : %v\n", err)
		versionCodeContent = ""
	}
	now := time.Now()
	// ISO 8601 格式 "YYYY-MM-DDTHH:mm:ssZ"
	iso8601Time := now.Format("20060102-150405")
	OutputResult("result."+iso8601Time+".txt", results, versionCodeContent)

	// 记录程序结束时间
	endTime := time.Now()
	totalDuration := endTime.Sub(startTime)
	fmt.Printf("Total elapsed time: %v\n", totalDuration)
}

package src

import (
	"fmt"
	"time"

	"rn-resource-checker/src/log"
	"rn-resource-checker/src/utils"
)

const SleepSecondAtJobDone = 5

func DoHashJob(targetDirs []string, whiteListSuffix []string) error {

	// 记录程序开始时间
	startTime := time.Now()

	// 指定要遍历的根目录列表

	switch len(targetDirs) {
	case 0:
		targetDirs = []string{"./"}
	case 1:
		targetDirs = utils.SplitArgString(targetDirs[0])
	}

	switch len(whiteListSuffix) {
	case 0:
		whiteListSuffix = []string{"mix", "exe", "dll", "ext"}
	case 1:
		whiteListSuffix = utils.SplitArgString(whiteListSuffix[0])
	}

	log.Info(fmt.Sprintf("Target path: %q", targetDirs))
	log.Info(fmt.Sprintf("Match file suffix list: %q", whiteListSuffix))

	// 指定并发 worker 的数量
	numWorkers := 8

	// 并发获取文件夹及其子文件夹中的文件名和路径
	files, errs := utils.WalkDirectoriesConcurrently(targetDirs, numWorkers, &whiteListSuffix)

	// 打印遇到的错误
	if len(errs) > 0 {
		log.Warn("Encountered errors:")
		for _, err := range errs {
			log.Error(err)
		}
	}

	results, err := utils.ComputeSHA256ForMultipleFiles(files, numWorkers)
	if err != nil {
		log.Error(fmt.Sprintf("Error: %v", err))
		return err
	}

	versionCodeContent, err := utils.ReadFileToString("Versioncode")
	if err != nil {
		log.Warn(fmt.Sprintf("load Version file failed : %v", err))
		versionCodeContent = ""
	}
	now := time.Now()
	// ISO 8601 格式 "YYYY-MM-DDTHH:mm:ssZ"
	iso8601Time := now.Format("20060102-150405")
	utils.OutputResult("result."+iso8601Time+".txt", results, versionCodeContent)

	// 记录程序结束时间
	endTime := time.Now()
	totalDuration := endTime.Sub(startTime)
	log.Info(fmt.Sprintf("Total elapsed time: %v", totalDuration))

	fmt.Printf("\nWaiting for %d seconds before exiting...", SleepSecondAtJobDone)
	time.Sleep(SleepSecondAtJobDone * time.Second)

	return nil
}

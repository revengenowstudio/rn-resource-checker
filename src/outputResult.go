package main

import (
	"fmt"
	"os"
	"sort"
)

func OutputResult(outputFileName string, fileResults []fileResult, versionCodeContent string) {
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		panic(err)
	}
	sort.Slice(fileResults, func(i, j int) bool {
		return fileResults[i].filename < fileResults[j].filename
	})
	outputFile.WriteString(versionCodeContent + "\n")
	for _, fileResult := range fileResults {
		outputFile.WriteString(fileResult.filename + " , hash : " + fmt.Sprintf("%x", fileResult.hash) + "\n")
	}
	outputFile.WriteString(fmt.Sprintf("Total file numbers : %d", len(fileResults)) + "\n")

}

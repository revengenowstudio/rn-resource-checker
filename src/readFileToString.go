package main

import (
	"bufio"
	"fmt"
	"os"
)

func ReadFileToString(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var contentBuilder []byte
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		contentBuilder = append(contentBuilder, line...)
		contentBuilder = append(contentBuilder, '\n')
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("failed to scan file: %w", err)
	}

	content := string(contentBuilder)
	return content, nil
}

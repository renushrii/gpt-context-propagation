package main

import (
	"bufio"
	"os"
	"path/filepath"
)

func listFiles(directory string) ([]string, error) {
	var files []string

	err := filepath.Walk(directory, visit(&files))
	if err != nil {
		return nil, err
	}

	return files, nil
}

func reWrite(filePath, response string) {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = file.Write([]byte(response))
	if err != nil {
		panic(err)
	}
}

func visit(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".go" {
			*files = append(*files, path)
		}
		return nil
	}
}

func ReadCode(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var content string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		content += scanner.Text() + "\n"
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return content, nil
}

package internal

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type FileProcessor struct {
	word string
}

func NewFileProcessor(word string) *FileProcessor {
	return &FileProcessor{
		word: word,
	}
}

func (f *FileProcessor) ProcessDir(jobs chan<- string, dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Mode().IsRegular() {
			jobs <- path
		}

		return nil
	})
}

func (f *FileProcessor) CountInstances(filepath string) (int, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return 0, fmt.Errorf("failed to open %s: %w", filepath, err)
	}
	defer file.Close()

	counter := 0
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		word := scanner.Text()
		if strings.EqualFold(word, f.word) {
			counter++
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("scanner error for %s: %w", filepath, err)
	}

	return counter, nil
}

package main

import (
	"fmt"
	"log"
	"os"
	"parallel-file-processor/internal"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("go run main.go <directory> [word]")
	}

	directory := os.Args[1]
	var word string
	if len(os.Args) > 2 {
		word = os.Args[2]
	}

	f := internal.NewFileProcessor(word)
	p := internal.NewWorkerPool(10)

	p.Run(f.CountInstances)

	go func() {
		defer p.Close()
		err := f.ProcessDir(p.Jobs, directory)
		if err != nil {
			log.Fatalf("error processing directory: %v", err)
		}
	}()

	totalCount := p.CollectResults()

	fmt.Printf("Total occurrences of '%s': %d\n", word, totalCount)
}

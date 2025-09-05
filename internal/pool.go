package internal

import (
	"fmt"
	"log"
	"sync"
)

type Pool struct {
	workers int
	Jobs    chan string
	results chan int
	wg      sync.WaitGroup
}

func NewWorkerPool(workers int) *Pool {
	return &Pool{
		Jobs:    make(chan string, 100),
		results: make(chan int, 100),
		workers: workers,
	}
}

func (p *Pool) Run(countFunc func(string) (int, error)) {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go func(workerID int) {
			defer p.wg.Done()
			for job := range p.Jobs {
				wordCount, err := countFunc(job)
				if err != nil {
					log.Printf("Worker %d failed to process %s: %v", workerID, job, err)
					p.results <- 0
				} else {
					fmt.Printf("Worker %d processed %s: %d occurrences\n", workerID, job, wordCount)
					p.results <- wordCount
				}
			}
		}(i)
	}
}

func (p *Pool) Close() {
	close(p.Jobs)
}

func (p *Pool) Wait() {
	p.wg.Wait()
}

func (p *Pool) CollectResults() int {
	go func() {
		p.Wait()
		close(p.results)
	}()

	totalCount := 0
	for count := range p.results {
		totalCount += count
	}

	return totalCount
}

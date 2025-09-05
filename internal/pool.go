package internal

import (
	"fmt"
	"log"
	"sync"
)

type Pool struct {
	workers   int
	jobs      chan string
	results   chan int
	wg        sync.WaitGroup
	totalJobs int
}

func NewWorkerPool(workers int) *Pool {
	jobs := make(chan string, 100)
	results := make(chan int, 100)
	return &Pool{
		jobs:      jobs,
		results:   results,
		workers:   workers,
		wg:        sync.WaitGroup{},
		totalJobs: 100,
	}
}

func (p *Pool) Run(count func(string) (int, error)) {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for job := range p.jobs {
				wordCount, err := count(job)
				if err != nil {
					log.Fatalf("failed to read %s because of error %v\n", job, err)
				} else {
					p.results <- wordCount
					fmt.Println()
				}
			}
		}()
	}
}

func (p *Pool) StartJob(filePath string) {
	p.jobs <- filePath
}

func (p *Pool) Wait() {
	p.wg.Wait()
}

func (p *Pool) Close() {
	close(p.jobs)
}

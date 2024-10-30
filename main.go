package main

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

func addWorker(workers *map[int]context.CancelFunc, id int, jobs <-chan string) {
	ctx, cancel := context.WithCancel(context.Background())
	(*workers)[id] = cancel
	go worker(ctx, id, jobs)
}
func deleteWorker(workers *map[int]context.CancelFunc, id int) {
	cancelFunc, ok := (*workers)[id]
	if ok {
		cancelFunc()
		delete(*workers, id)
	}
}
func worker(ctx context.Context, id int, jobs <-chan string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("worker %d deleted\n", id)
			return
		case job, ok := <-jobs:
			if !ok {
				return
			}
			fmt.Println(job, "started by", id)
			time.Sleep(time.Second)
			fmt.Println(job, "processed by", id)
		}
	}
}

func main() {
	startWorkers := 4
	jobs := make(chan string, startWorkers)
	workers := make(map[int]context.CancelFunc)
	for i := 0; i < startWorkers; i++ {
		jobs <- "job" + strconv.Itoa(i)
	}
	for i := 0; i < startWorkers; i++ {
		addWorker(&workers, i, jobs)
	}
	deleteWorker(&workers, 0)
	time.Sleep(time.Second * 3)
	close(jobs)
}

package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	tasks := make(chan int, 5)
	var wg sync.WaitGroup

	for i := 1; i < 4; i++ {
		wg.Add(1)
		go worker(tasks, &wg, i)
	}

	for i := 1; i < 10; i++ {
		tasks <- i
	}

	close(tasks)
	wg.Wait()
}

func worker(tasks chan int, wg *sync.WaitGroup, worker int) {
	defer wg.Done()
	for val := range tasks {
		fmt.Printf("Worker %d starting task %d\n", worker, val)
		time.Sleep(time.Second)
		fmt.Printf("Worker %d finished task %d\n", worker, val)
	}
}

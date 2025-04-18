package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type Tasks struct {
	mu    sync.Mutex
	count int64
	tasks []int
}

func (t *Tasks) GetTasks() []int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.tasks
}

func (t *Tasks) AddTask(task int, wg *sync.WaitGroup) {
	defer wg.Done()
	t.mu.Lock()
	t.tasks = append(t.tasks, task)
	t.mu.Unlock()
	atomic.AddInt64(&t.count, 1)
}

func main() {
	var wg sync.WaitGroup
	t := Tasks{tasks: make([]int, 0)}
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go t.AddTask(i, &wg)
	}
	wg.Wait()
	fmt.Println("Tasks:", t.GetTasks())
	fmt.Println("Count of tasks:", t.count)
}

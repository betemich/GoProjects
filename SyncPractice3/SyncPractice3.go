package main

import (
	"fmt"
	"sync"
)

type RequestCounter struct {
	counter_map sync.Map
}

func (rc *RequestCounter) Increment(url string) {
	c, _ := rc.counter_map.LoadOrStore(url, 0)
	rc.counter_map.Store(url, c.(int)+1)
}

func (rc *RequestCounter) Get(url string) (int, bool) {
	c, ok := rc.counter_map.Load(url)
	if ok {
		return c.(int), true
	} else {
		return 0, false
	}
}

func (rc *RequestCounter) Clear() {
	rc.counter_map = sync.Map{}
}

func CreateRequestCounter() *RequestCounter {
	return &RequestCounter{}
}

func CreateWaitGroup() *sync.WaitGroup {
	return &sync.WaitGroup{}
}

func main() {
	wg := CreateWaitGroup()
	rc := CreateRequestCounter()
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			rc.Increment(fmt.Sprintf("https://%d.ru", i/2))
		}()
	}
	wg.Wait()
	for i := 0; i < 30; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			url := fmt.Sprintf("https://%d.ru", i)
			if val, ok := rc.Get(url); ok {
				fmt.Printf("%s count: %d\n", url, val)
			} else {
				fmt.Printf("%s is not found\n", url)
			}
		}()
	}
	wg.Wait()
	fmt.Println("Incrementing finished")
	rc.Clear()
	rc.counter_map.Range(func(key, value interface{}) bool {
		fmt.Printf("key: %v; value: %v\n", key, value)
		return true
	})
}

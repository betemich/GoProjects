package main

import (
	"fmt"
	"sync"
)

type SafeCache struct {
	rmu   sync.RWMutex
	cache map[string]interface{}
}

func (sc *SafeCache) Set(key string, value interface{}) {
	sc.rmu.Lock()
	sc.cache[key] = value
	sc.rmu.Unlock()
}

func (sc *SafeCache) Get(key string) (interface{}, bool) {
	sc.rmu.RLock()
	defer sc.rmu.RUnlock()
	if sc.cache[key] == nil {
		return nil, false
	}
	return sc.cache[key], true
}

func (sc *SafeCache) Delete(key string) {
	sc.rmu.Lock()
	delete(sc.cache, key)
	sc.rmu.Unlock()
}

func (sc *SafeCache) Clear() {
	sc.rmu.Lock()
	sc.cache = make(map[string]interface{})
	sc.rmu.Unlock()
}

func CreateCache() *SafeCache {
	return &SafeCache{cache: make(map[string]interface{})}
}

func CreateWaitGroup() *sync.WaitGroup {
	var wg sync.WaitGroup
	return &wg
}

func main() {
	wg := CreateWaitGroup()
	sc := CreateCache()

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sc.Set("key1", i)
		}()
	}
	wg.Wait()

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if val, ok := sc.Get("key1"); ok {
				fmt.Printf("key1: %v\n", val)
			} else {
				fmt.Print("Key1 not found\n")
			}
		}()
	}
	wg.Wait()
}

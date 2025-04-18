package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

func UrlProcessing(url string, ctx context.Context) {
	tm_ctx, timeout := context.WithTimeout(ctx, 1*time.Second)
	defer timeout()
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println(err)
		return
	}
	request = request.WithContext(tm_ctx)
	client := &http.Client{}
	start := time.Now()
	result, err := client.Do(request)
	response_time := time.Since(start)
	if err != nil {
		select {
		case <-tm_ctx.Done():
			log.Printf("%s: too long request\n", url)
		default:
			log.Printf("%s: %v\n", url, err)
		}
		return
	}
	fmt.Printf("url: %s, status code: %d, response time: %v\n", url, result.StatusCode, response_time)
}

func main() {
	var wg sync.WaitGroup
	urls := os.Args[1:]
	ctx := context.Background()
	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			UrlProcessing(url, ctx)
		}(url)
	}
	wg.Wait()

}

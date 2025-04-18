package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

//cancel context

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		err := CancelRequest(ctx)
		if err != nil {
			cancel()
		}
	}()
	DoRequest(ctx, "https://ya.ru")
}

func CancelRequest(ctx context.Context) error {
	time.Sleep(100 * time.Millisecond)
	return fmt.Errorf("fail request")
}

func DoRequest(ctx context.Context, requestStr string) {
	req, _ := http.NewRequest(http.MethodGet, requestStr, nil)
	req = req.WithContext(ctx)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("request too long")
		return
	}

	select {
	case <-time.After(500 * time.Millisecond):
		fmt.Printf("request done: status %d\n", res.StatusCode)
	case <-ctx.Done():
		fmt.Println("request too long")
	}

}

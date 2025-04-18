package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	client := http.Client{
		Timeout: 500 * time.Millisecond,
	}
	start := time.Now()
	resp, err := client.Get("https://almaz.comfortkino.ru/")
	response_time := time.Since(start)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	fmt.Printf("Response time: %v", response_time)

	io.Copy(os.Stdout, resp.Body)

}

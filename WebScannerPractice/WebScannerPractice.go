package main

import (
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

type URL struct {
	response_time time.Duration
	status_code   int16
	err           error
}

func (u *URL) Error() string {
	return fmt.Sprintf("%v", u.err)
}

type Result struct {
	result_map        sync.Map
	count             int64
	success           int64
	sum_response_time time.Duration
}

func (r *Result) ProcessURL(ch chan string) {
	for url := range ch {
		atomic.AddInt64(&r.count, 1)
		url_struct := &URL{}
		client := &http.Client{
			Timeout: 30 * time.Second,
		}
		start := time.Now()
		res, err := client.Get(url)
		url_struct.response_time = time.Since(start)
		r.sum_response_time += url_struct.response_time
		if err != nil {
			url_struct.err = err
			r.result_map.Store(url, url_struct)
			return
		}
		if res.StatusCode != http.StatusOK {
			url_struct.err = fmt.Errorf("bad status: %s", res.Status)
			r.result_map.Store(url, url_struct)
			return
		}
		atomic.AddInt64(&r.success, 1)
		url_struct.status_code = int16(res.StatusCode)
		r.result_map.Store(url, url_struct)
	}
}

func (r *Result) PrintStatistic() {
	r.result_map.Range(func(key, value interface{}) bool {
		fmt.Printf("%s: ", key)
		response_time := value.(*URL).response_time
		status_code := value.(*URL).status_code
		err := value.(*URL).err
		fmt.Printf("response time %v, ", response_time)
		if status_code != 0 {
			fmt.Printf("status code %d, ", status_code)
		}
		if err != nil {
			fmt.Print(value.(*URL).Error())
		}
		fmt.Print("\n")
		return true
	})
	fmt.Println("Total:", r.count)
	fmt.Println("Succesful:", r.success)
	fmt.Println("Average response time:", r.sum_response_time/time.Duration(r.count))
}

func main() {
	r := Result{}
	ch := make(chan string, 10)
	var wg sync.WaitGroup
	urls := [...]string{
		"https://almaz.comfortkino.ru/",
		"https://habr.com/ru/articles/759584/",
		"https://studlk.susu.ru/Account/Login?ReturnUrl=%2fru%2fSchedule",
		"https://www.deepl.com/ru/translator",
		"https://ya.ru",
		"https://metanit.com/go/tutorial/7.2.php",
		"https://www.google.com",
		"https://www.github.com",
		"https://www.nonexistentwebsite123.com",
		"https://www.stackoverflow.com",
	}

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r.ProcessURL(ch)
		}()
	}

	for _, val := range urls {
		ch <- val
	}
	close(ch)
	wg.Wait()
	r.PrintStatistic()
}

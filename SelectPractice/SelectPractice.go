package main

import (
	"fmt"
	"time"
)

func main() {
	IP := make(chan int, 10)
	DNS := make(chan int, 10)
	exit := make(chan int)
	go SendIP(IP)
	go SendDNS(DNS)
	go Stop(exit)

	select {
	case <-IP:
		fmt.Println("Reading IP")

	case <-DNS:
		fmt.Println("Reading DNS")

	case <-exit:
		fmt.Println("EXIT")
		return
	}

}

func SendIP(IP chan int) {
	time.Sleep(2 * time.Second)
	IP <- 500
}

func SendDNS(DNS chan int) {
	time.Sleep(time.Second)
	DNS <- 600
}

func Stop(exit chan int) {
	time.Sleep(4 * time.Second)
	exit <- 0
}

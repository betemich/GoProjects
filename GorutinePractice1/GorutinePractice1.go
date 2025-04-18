package main

import (
	"fmt"
	"math/rand"
)

func main() {
	ch := make(chan uint32)
	var INT_VAL uint32 = 0
	go SendCode(ch)
	for i := range ch {
		INT_VAL = i
		fmt.Println(INT_VAL)
	}
}

func SendCode(ch chan uint32) {
	for range 10 {
		a := rand.Uint32()
		ch <- a
	}
	close(ch)
}

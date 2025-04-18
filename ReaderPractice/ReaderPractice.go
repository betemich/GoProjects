package main

import (
	"fmt"
	"io"
	"strings"
)

func main() {
	str := "Hello, world!"
	rdr := strings.NewReader(str)
	buf := make([]byte, 5)
	for {
		n, err := io.ReadFull(rdr, buf)
		if err == io.EOF {
			fmt.Println("End of file")
			return
		}

		if err == nil {
			fmt.Printf("%s\n", buf[:n])
		}
	}
}

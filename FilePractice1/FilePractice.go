package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// Чтение файла до делимитера с помощью bufio.ReadString
func main() {
	file, err := os.Open("file1.txt")
	if err != nil {
	}
	defer file.Close()
	scanner := bufio.NewReader(file)
	for {
		str, err := scanner.ReadString('\n')
		fmt.Print(str)
		if err == io.EOF {
			return
		}
	}
}

package main

import (
	"log"
	"os"
)

// Построчная запись в файл c дефолтным io.Writer
func main() {
	file, err := os.Create("file1.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	file.WriteString("XXX\n")
	file.WriteString("Tototo")
}

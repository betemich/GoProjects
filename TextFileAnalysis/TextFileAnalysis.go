package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func Recovering() {
	if err := recover(); err != nil {
		log.Println(err)
	}
}

func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type FileInfo struct {
	number_of_lines   uint
	number_of_words   uint
	number_of_symbols uint
}

func WriteStatToFile(file_name string, file_info *FileInfo) {
	file_name = strings.TrimSuffix(filepath.Base(file_name), ".txt")
	result_file, err := os.Create(fmt.Sprintf("%s_stats.txt", file_name))
	CheckErr(err)
	defer result_file.Close()
	writer := bufio.NewWriter(result_file)
	_, err = writer.WriteString(fmt.Sprintf("Lines: %d\n", file_info.number_of_lines))
	CheckErr(err)
	_, err = writer.WriteString(fmt.Sprintf("Words: %d\n", file_info.number_of_words))
	CheckErr(err)
	_, err = writer.WriteString(fmt.Sprintf("Characters: %d\n", file_info.number_of_symbols))
	CheckErr(err)
	err = writer.Flush()
	CheckErr(err)
}

func AnalyzeFile() {
	defer Recovering()
	file_name := os.Args[1]
	file, err := os.Open(file_name)
	CheckErr(err)
	defer file.Close()
	file_info := FileInfo{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		str := scanner.Text()
		words := strings.Fields(str)
		for _, val := range words {
			file_info.number_of_symbols += uint(len(val)) + 1
			file_info.number_of_words += 1
		}
		file_info.number_of_lines++
	}

	WriteStatToFile(file_name, &file_info)
	log.Println("File has analized")
}

func main() {
	AnalyzeFile()
}

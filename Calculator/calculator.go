package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

var symbols [4]string
var a int
var b int

func contains(el string) bool {
	for _, l := range symbols {
		if el == l {
			return true
		}
	}
	return false
}

func isCorrect(elements []string) bool {
	if len(elements) != 3 {
		return false
	}
	if _, err := strconv.Atoi(elements[0]); err != nil {
		return false
	}
	if _, err := strconv.Atoi(elements[2]); err != nil {
		return false
	}
	if !contains(elements[1]) {
		return false
	}
	return true
}

func calculate(elements []string) {
	a, _ = strconv.Atoi(elements[0])
	b, _ = strconv.Atoi(elements[2])
	switch elements[1] {
	case "+":
		fmt.Printf("%d + %d = %d\n", a, b, a+b)
	case "-":
		fmt.Printf("%d - %d = %d\n", a, b, a-b)
	case "*":
		fmt.Printf("%d * %d = %d\n", a, b, a*b)
	case "/":
		if b != 0 {
			fmt.Printf("%d / %d = %d\n", a, b, a/b)
		} else {
			fmt.Println("Zero division error")
		}
	}
}

func run(example string, wg *sync.WaitGroup) {
	defer wg.Done()
	if example == "Exit" {
		return
	}
	elements := strings.Fields(example) //Разбиение строки на части, содержащие один или несколько пробелов
	if isCorrect(elements) {
		calculate(elements)
	} else {
		fmt.Println("Invalid Input")
	}
}

func main() {
	symbols[0] = "+"
	symbols[1] = "-"
	symbols[2] = "*"
	symbols[3] = "/"
	fmt.Println("Enter the expressions(for example 5 + 5) ")
	fmt.Println("If you want to exit print Exit")
	var wg sync.WaitGroup
	for {
		example, _ := bufio.NewReader(os.Stdin).ReadString('\n') //Создаем объект Reader и передаем ему на вход стандартный ввод и читаем все до \n
		fmt.Println(example)
		if example == "Exit\n" {
			fmt.Println("aa")
			wg.Wait()
			return
		}
		wg.Add(1)
		go run(example, &wg)
	}
}

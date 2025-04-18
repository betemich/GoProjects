package main

import (
	"errors"
	"fmt"
	"log"
)

type AppError struct {
	message string
	err     error
}

func (ae *AppError) Error() string {
	return ae.message
}

func main() {
	CauseDefaultPanic()
	fmt.Println("Message after first panic")
	CauseAppErrorPanic()
	fmt.Println("Message after second panic")
}

func Recovering() {
	var appErr *AppError
	if err := recover(); err != nil {
		switch err := err.(type) {
		case error:
			if errors.As(err, &appErr) {
				log.Println("app error panic")
			} else {
				log.Println("some error panic")
			}
		default:
			log.Println("default panic")
		}
	}
}

func CauseDefaultPanic() {
	defer Recovering()
	panic("default panic")
}

func CauseAppErrorPanic() {
	defer Recovering()
	panic(&AppError{
		message: "app error",
		err:     nil,
	})
}

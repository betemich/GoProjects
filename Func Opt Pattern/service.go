package main

import (
	"errors"
	"fmt"
	"log"
)

type Service struct {
	addr string
	port int
}

type Options struct {
	port *int
}

func NewService(addr string, opts ...Option) (*Service, error) {
	options := Options{}
	for _, opt := range opts {
		err := opt(&options)
		if err != nil {
			return nil, err
		}
	}

	port := 0
	if options.port != nil {
		port = *options.port
	}

	service := Service{
		addr: addr,
		port: port,
	}

	return &service, nil
}

type Option func(options *Options) error

func WithPort(port int) Option {
	return func(options *Options) error {
		if port < 0 {
			return errors.New("port must be positive")
		}
		options.port = &port
		return nil
	}
}

func main() {
	service, err := NewService("localhost", WithPort(8080))
	if err != nil {
		log.Printf("error: %v\n", err)
	}
	fmt.Println(service)
}

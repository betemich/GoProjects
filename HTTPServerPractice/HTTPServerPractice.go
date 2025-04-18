package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/favicon.ico" {
		return
	}

	switch r.Method {
	case "GET":
		w.Write([]byte("GET"))
	case "POST":
		w.Write([]byte("POST"))
	default:
		w.Write([]byte("DEFAULT"))
	}
	log.Printf("Home request for path %s succesfully processed", r.URL.Path)
}

func Recovering(w *http.ResponseWriter) {
	if err := recover(); err != nil {
		http.Error(*w, "Internal server error", http.StatusInternalServerError)
		log.Printf("Panic: %v", err)
	}
}

func DivideHandler(w http.ResponseWriter, r *http.Request) {
	defer Recovering(&w)
	a, err := strconv.Atoi(r.URL.Query().Get("a"))
	if err != nil {
		http.Error(w, "Wrong value", http.StatusBadRequest)
		return
	}
	b, err := strconv.Atoi(r.URL.Query().Get("b"))
	if err != nil {
		http.Error(w, "Wrong value", http.StatusBadRequest)
		return
	}
	res := a / b
	w.Write([]byte(fmt.Sprintf("Result: %d", res)))
	log.Printf("Divide request for %s with method %s succesfully processed", r.URL.Path, r.Method)
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("%v", http.StatusOK)))
	log.Printf("Status request for %s succesfully processed", r.URL.Path)
}

func main() {
	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/divide", DivideHandler)
	http.HandleFunc("/status", StatusHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Error starting server")
		return
	}

}

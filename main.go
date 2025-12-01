package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	mux := createHttpMux()
	addr := "localhost:8080"
	fmt.Println("Starting local server at", addr)
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Server stopped")
}

func createHttpMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleDefault)
	return mux
}

func handleDefault(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello, HTTP!\n")
}

package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	setupServer()
	launchServer("localhost:8080")
}

func launchServer(addr string) {
	fmt.Println("Started local server", addr)
	fmt.Println("To access the server")
	fmt.Println("1. Open your web browser")
	fmt.Println("2. Type \"localhost:8080\" in the address bar")

	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Println(err)
	}

	fmt.Println("Server stopped")
}

func setupServer() {
	http.HandleFunc("/", getIndex)
}

func getIndex(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello, HTTP!\n")
}

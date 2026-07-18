package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

func main() {
	var addr = "localhost:8080"
	if len(os.Args) == 2 {
		if os.Args[1] == "Host@Render" {
			addr = "0.0.0.0:10000"
		}
	}
	setupServer()
	launchServer(addr)
}

func launchServer(addr string) {
	fmt.Println("Server started @", addr)
	if strings.HasPrefix(addr, "localhost:") {
		fmt.Println()
		fmt.Println("To access the server")
		fmt.Println("1. Open your web browser")
		fmt.Printf("2. Type \"%s\" in the address bar\n", addr)
	}

	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Println(err)
	}
}

func setupServer() {
	http.HandleFunc("GET /", defaultActivityPresenter)
	http.HandleFunc("POST /", defaultActivityPresenter)
	http.HandleFunc("GET /ping", pingResponder)
	http.HandleFunc("GET /{activity}", activityRequester)
	fs := http.FileServer(http.Dir("style"))
	http.Handle("GET /style.css", fs)
}

func defaultActivityPresenter(w http.ResponseWriter, r *http.Request) {
	activity := getDefaultActivity()
	activity.Respond(w, r)
}

func activityRequester(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("activity")
	activity := selectActivityById(id)
	activity.Respond(w, r)
}

// Gets called by 3rd-party keep-alive service!
func pingResponder(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Ping'd")
}

package main

import (
	"fmt"
	"net/http"
	"os"
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
	fmt.Println("Started the server", addr)
	fmt.Println()
	fmt.Println("To access the server")
	fmt.Println("1. Open your web browser")
	fmt.Printf("2. Type \"%s\" in the address bar\n", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Println(err)
	}
}

func setupServer() {
	http.HandleFunc("GET /", activityPresenter)
	http.HandleFunc("POST /", activityPresenter)
	http.HandleFunc("GET /{activity}", activityRequester)
	fs := http.FileServer(http.Dir("style"))
	http.Handle("GET /style.css", fs)
}

func activityPresenter(w http.ResponseWriter, r *http.Request) {
	activity := getActivity()
	activity.Respond(w, r)
}

func activityRequester(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("activity")
	selectActivityById(id)
	activityPresenter(w, r)
}

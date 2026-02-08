package main

import (
	"fmt"
	"net/http"
)

func main() {
	startServing()
}

func startServing() {
	setupServer()
	launchServer("localhost:8080")
}

func launchServer(addr string) {
	fmt.Println("Started local server", addr)
	fmt.Println()
	fmt.Println("To access the server")
	fmt.Println("1. Open your web browser")
	fmt.Println("2. Type \"localhost:8080\" in the address bar")

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

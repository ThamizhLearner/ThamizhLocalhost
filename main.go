package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

// Log GET/POST requests...
var chatterOn bool = false

func main() {
	var addr = "localhost:8080"
	if len(os.Args) == 2 {
		if os.Args[1] == "Host@Render" {
			port := os.Getenv("PORT")
			fmt.Printf("%s requested using port %s\n", os.Args[1], port)
			addr = "0.0.0.0:10000"

			// Chatter mode!
			chatterOn = os.Getenv("Chatter") == "1"
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
	http.HandleFunc("POST /{activity}", activityPresenter)
	http.HandleFunc("GET /ping", pingResponder)
	http.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	http.HandleFunc("GET /{activity}", activityPresenter)
	fs := http.FileServer(http.Dir("style"))
	http.Handle("GET /style.css", fs)
}

func defaultActivityPresenter(w http.ResponseWriter, r *http.Request) {
	if chatterOn {
		fmt.Println("GET Default")
	}
	activity := getDefaultActivity()
	activity.Respond(w, r)
}

func activityPresenter(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("activity")

	if chatterOn {
		if r.Method == http.MethodPost {
			if id == "SylWord" || id == "Decomposition" {
				query := r.FormValue("inpStr")
				fmt.Printf("POST %s; %s\n", id, query)
			} else {
				fmt.Printf("POST %s\n", id)
			}
		} else {
			fmt.Printf("GET %s\n", id)
		}
	}

	activity := selectActivityById(id)
	activity.Respond(w, r)
}

// Gets called by 3rd-party keep-alive service!
func pingResponder(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Ping'd")
}

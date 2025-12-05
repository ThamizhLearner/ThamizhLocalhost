package main

import (
	"fmt"
	"html/template"
	"net/http"

	script "github.com/ThamizhLearner/Thamizh"
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
	fs := http.FileServer(http.Dir("./"))
	http.Handle("/style.css", fs)
}

var tmpl = template.Must(template.ParseFiles("index.tmpl"))

func getIndex(w http.ResponseWriter, r *http.Request) {
	post := r.Method == http.MethodPost

	seed := struct {
		InpStr string
		SylStr string
	}{"", ""}

	if post {
		seed.InpStr = r.FormValue("inpStr")
		str, ok := script.Decode(seed.InpStr)
		if ok {
			seed.SylStr = str.SyllabifiedUStr("-")
		}
	}

	tmpl.Execute(w, seed)
}

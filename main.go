package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	script "github.com/ThamizhLearner/Thamizh"
)

func main() {
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

	fmt.Println("Server stopped")
}

func setupServer() {
	http.HandleFunc("/", getIndex)
	fs := http.FileServer(http.Dir("style"))
	http.Handle("/style.css", fs)
}

var tmpl = template.Must(template.ParseFiles("tmpls/index.tmpl", "tmpls/sylWord.tmpl"))

func getIndex(w http.ResponseWriter, r *http.Request) {
	post := r.Method == http.MethodPost

	seed := struct {
		InpStr      string
		LetterCount int
		SylStr      string
		SylCount    int
		Graph       string
	}{"", 0, "", 0, ""}

	if post {
		seed.InpStr = r.FormValue("inpStr")
		str, ok := script.Decode(seed.InpStr)
		if ok {
			seed.LetterCount = str.Len()
			seed.SylStr, seed.SylCount = str.SyllabifiedUStr("-")
			seed.Graph = createSylGraph(seed.InpStr, strings.Split(seed.SylStr, "-"))
		}
	}

	tmpl.Execute(w, seed)
}

func createSylGraph(w string, syls []string) string {
	sb := strings.Builder{}
	sb.WriteString("graph TB\n")
	for idx, syl := range syls {
		if idx == 0 {
			sb.WriteString(fmt.Sprintf("N[%v] --> ", w))
		} else {
			sb.WriteString("N --> ")
		}
		sb.WriteString(fmt.Sprintf("N%v[%v]\n", idx, syl))
	}
	return sb.String()
}

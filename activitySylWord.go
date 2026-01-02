package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	script "github.com/ThamizhLearner/Thamizh"
)

type stylWordActivity struct{}

func (a stylWordActivity) GetID() string   { return "StylWord" }
func (a stylWordActivity) GetDesc() string { return "English style syllabification" }
func (a stylWordActivity) Respond(w http.ResponseWriter, r *http.Request) {
	post := r.Method == http.MethodPost // GET or POST response

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

	var tmpl = template.Must(template.ParseFiles("tmpls/index.tmpl", "tmpls/sylWord.tmpl"))
	tmpl.Execute(w, seed)
}

// Syllable derivation graph
func createSylGraph(w string, syls []string) string {
	// Graph driven by Mermaid (https://mermaid.js.org/)
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

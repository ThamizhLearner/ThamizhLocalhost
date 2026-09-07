package main

import (
	"fmt"
	"html/template"
	"math/rand/v2"
	"net/http"
	"strings"

	script "github.com/ThamizhLearner/Thamizh"
)

type sylWordActivity struct{}

func (a sylWordActivity) GetID() string   { return "SylWord" }
func (a sylWordActivity) GetDesc() string { return "English style syllabification" }
func (a sylWordActivity) Respond(w http.ResponseWriter, r *http.Request) {
	post := r.Method == http.MethodPost // GET or POST response

	seed := struct {
		InpStr      string
		LetterCount int
		SylStr      string
		SylStrEg    string
		SylCount    int
		Graph       string
	}{"", 0, "", "", 0, ""}

	if post {
		seed.InpStr = strings.TrimSpace(r.FormValue("inpStr"))
		str, ok := script.LetterSeqFrom(seed.InpStr)
		if ok {
			seed.LetterCount = str.Len()
			seed.SylStr, seed.SylCount = script.SyllabifiedUStr(str, "-")
			seed.Graph = createSylGraph(seed.InpStr, strings.Split(seed.SylStr, "-"))
		}
	} else {
		// Basically, there is no input from the user! And we simply want to show the expected result.

		strs := []string{"தொல்காப்பியம்", "திருக்குறள்", "நூலகம்", "மரபியல்", "தமிழிலக்கணம்", "கல்விக்கழகம்",
			"பெயர்கள்", "யாப்பிலக்கணம்", "உரையமைப்பும்", "உள்ளடக்கம்", "வாசிக்க", "ஆசிரியர்",
			"வடிவம்", "விடுதலை", "வேற்றுமை", "புணர்ச்சி", "படிக்கும்போது", "வரையறைகளை"}
		seed.InpStr = strs[rand.IntN(len(strs))]
		str, ok := script.LetterSeqFrom(seed.InpStr)
		if ok {
			seed.LetterCount = str.Len()
			seed.SylStrEg, seed.SylCount = script.SyllabifiedUStr(str, "-")
			seed.Graph = createSylGraph(seed.InpStr, strings.Split(seed.SylStrEg, "-"))
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
			sb.WriteString(fmt.Sprintf("N(%v) --> ", w))
		} else {
			sb.WriteString("N --> ")
		}
		sb.WriteString(fmt.Sprintf("N%v(%v)\n", idx, syl))
	}
	return sb.String()
}

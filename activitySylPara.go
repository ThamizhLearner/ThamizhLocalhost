package main

import (
	"html/template"
	"net/http"
	"strings"

	script "github.com/ThamizhLearner/Thamizh"
)

type sylParaActivity struct{}

func (a sylParaActivity) GetID() string   { return "SylPara" }
func (a sylParaActivity) GetDesc() string { return "English style syllabification" }

func (a sylParaActivity) Respond(w http.ResponseWriter, r *http.Request) {
	post := r.Method == http.MethodPost // GET or POST response

	seed := struct {
		InpStr string
		SylStr string
	}{"", ""}

	if post {
		seed.InpStr = r.FormValue("inpStr")
		seed.InpStr = strings.TrimSpace(seed.InpStr)
		strs := strings.Split(seed.InpStr, " ")
		for i, s := range strs {
			ustr := strings.Trim(s, ".,'-:\"") // Trim away punctuation marks
			str, ok := script.Decode(ustr)
			if ok {
				strs[i], _ = str.SyllabifiedUStr("-") // Replace with syllabified version
			}
		}
		seed.SylStr = strings.Join(strs, " | ")
	}

	var tmpl = template.Must(template.ParseFiles("tmpls/index.tmpl", "tmpls/sylPara.tmpl"))
	tmpl.Execute(w, seed)
}

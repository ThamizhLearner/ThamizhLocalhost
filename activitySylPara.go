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
		// Separate by space!
		// Attempt to decode each continuation... and be done!
		strs := strings.Split(seed.InpStr, " ")
		for i := 0; i < len(strs); i++ {
			ustr := strings.Trim(strs[i], "., \n\t")
			str, ok := script.Decode(ustr)
			if ok {
				strs[i], _ = str.SyllabifiedUStr("-")
			}
		}
		seed.SylStr = strings.Join(strs, " | ")
	}

	var tmpl = template.Must(template.ParseFiles("tmpls/index.tmpl", "tmpls/sylPara.tmpl"))
	tmpl.Execute(w, seed)
}

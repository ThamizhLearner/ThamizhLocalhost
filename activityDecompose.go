package main

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/ThamizhLearner/ThamizhLocalhost/decomposition"
)

type decompActivity struct{}

func (a decompActivity) GetID() string   { return "Decomposition" }
func (a decompActivity) GetDesc() string { return "Syntax analysis" }
func (a decompActivity) Respond(w http.ResponseWriter, r *http.Request) {
	post := r.Method == http.MethodPost // GET or POST response

	seed := struct {
		InpStr  string
		ResStrs []string
	}{"", nil}

	if post {
		seed.InpStr = strings.TrimSpace(r.FormValue("inpStr"))
		seed.ResStrs = append(seed.ResStrs, "---------------------")
		for _, strs := range decomposition.DecomposeWord(seed.InpStr) {
			seed.ResStrs = append(seed.ResStrs, strs...)
			seed.ResStrs = append(seed.ResStrs, "---------------------")
		}
	}

	var tmpl = template.Must(template.ParseFiles("tmpls/index.tmpl", "tmpls/decomp.tmpl"))
	tmpl.Execute(w, seed)
}

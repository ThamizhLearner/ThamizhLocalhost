package main

import (
	"html/template"
	"net/http"
)

type miscActivity struct{}

func (a miscActivity) GetID() string   { return "Misc" }
func (a miscActivity) GetDesc() string { return "Uncategorized cache" }
func (a miscActivity) Respond(w http.ResponseWriter, r *http.Request) {
	seed := struct {
		InfoTable SimpleTable
		VerbGraph string
	}{finalTable2(), createVerbGraph()}

	var tmpl = template.Must(template.ParseFiles("tmpls/index.tmpl", "tmpls/misc.tmpl"))
	tmpl.Execute(w, seed)
}

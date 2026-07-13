package main

import (
	"bufio"
	"html/template"
	"net/http"
	"os"
)

type kuralActivity struct{}

func (a kuralActivity) GetID() string   { return "Kural" }
func (a kuralActivity) GetDesc() string { return "Thirukkural" }
func (a kuralActivity) Respond(w http.ResponseWriter, r *http.Request) {
	seed := struct {
		Kurals []kural
	}{LoadAllKurals("kural/kurals.txt")}

	var tmpl = template.Must(template.ParseFiles("tmpls/index.tmpl", "tmpls/kural.tmpl"))
	tmpl.Execute(w, seed)
}

type kural struct {
	Num         int
	Line, Line2 string
}

func LoadAllKurals(fn string) []kural {
	file, err := os.Open(fn)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var kurals []kural
	num := 1
	for scanner.Scan() {
		line := scanner.Text()
		if !scanner.Scan() {
			panic("Unexpected EOF")
		}
		kurals = append(kurals, kural{num, line, scanner.Text()})
		num++
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return kurals
}

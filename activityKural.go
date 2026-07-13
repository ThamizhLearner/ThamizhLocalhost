package main

import (
	"bufio"
	"html/template"
	"net/http"
	"os"
	"strings"

	script "github.com/ThamizhLearner/Thamizh"
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
	Syl, Syl2   string
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
		line2 := scanner.Text()
		syl, syl2 := SyllabifyLong(line), SyllabifyLong(line2)
		kurals = append(kurals, kural{num, line, line2, syl, syl2})
		num++
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return kurals
}

func SyllabifyLong(str string) string {
	strs := strings.Split(str, " ")
	for i, ustr := range strs {
		str, ok := script.LetterSeqFrom(ustr)
		if !ok {
			continue
		}
		strs[i], _ = script.SyllabifiedUStr(str, "-")
	}
	return strings.Join(strs, " ")
}

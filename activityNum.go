package main

import (
	"html/template"
	"net/http"
	"strconv"

	script "github.com/ThamizhLearner/Thamizh"
)

type numActivity struct{}

func (a numActivity) GetID() string   { return "Num" }
func (a numActivity) GetDesc() string { return "Thamizh numbers" }
func (a numActivity) Respond(w http.ResponseWriter, r *http.Request) {
	seed := struct {
		NumTable SimpleTable
	}{numListing()}

	var tmpl = template.Must(template.ParseFiles("tmpls/index.tmpl", "tmpls/numbers.tmpl"))
	tmpl.Execute(w, seed)
}

func numListing() SimpleTable {
	t := SimpleTable{
		Title:       "Number listing",
		ColInfoList: []ColInfo{{"Number", 1}, {"Thamizh", 1}, {"Syllabified", 1}},
		Cells:       make([][]string, 20),
	}
	for i := range 20 {
		v, ok := getCardinalNumStr(i)
		if !ok {
			panic("Coding error")
		}
		row := make([]string, 3)
		row[0] = strconv.Itoa(i)
		row[1] = v
		row[2], _ = script.SyllabifiedUStr(script.MustLetterSeqFrom(v), "-")
		t.Cells[i] = row
	}
	return t
}

// Range (single digit): [0, 9]
// Range (double digits): [10, 99]
// 	Range (double digits): [10, 19]
// 	Range (double digits): [20, 99]
//		Range (tens): {10, 90, 10}
// Range (triple digits): [100, 999]
// 	Range (hundreds): {100, 900, 100}
// Range (quadruple digits): [1000, 9999]
// 	Range (thousands): {1000, 9000, 1000}

func getCardinalNumStr(num int) (string, bool) {
	if num < 0 {
		return "", false
	}
	if num < 10 {
		return unit[num], true
	}
	if num == 10 {
		return tens[0], true
	}
	if num < 20 {
		return unit2[num-11], true
	}
	return "", false
}

func getOrdinalNumStr(ord int) (string, bool) {
	if ord < 1 {
		return "", false
	}
	panic("")
}

var unit []string = []string{
	"சுழியம்",
	"ஒன்று",
	"இரண்டு",
	"மூன்று",
	"நான்கு",
	"ஐந்து",
	"ஆறு",
	"ஏழு",
	"எட்டு",
	"ஒன்பது",
}

var unit2 []string = []string{
	"பதினொன்று",
	"பன்னிரண்டு",
	"பதின்மூன்று",
	"பதினான்கு",
	"பதினைந்து",
	"பதினாறு",
	"பதினேழு",
	"பதினெட்டு",
	"பத்தொன்பது",
}

var tens []string = []string{
	"பத்து",    // பதினொன்று
	"இருபது",   // இருபத்து ஒன்று
	"முப்பது",  // முப்பத்து ஒன்று
	"நாற்பது",  // நாற்பத்து ஒன்று
	"ஐம்பது",   // ஐம்பத்து ஒன்று
	"அறுபது",   // அறுபத்து ஒன்று
	"எழுபது",   // எழுபத்து ஒன்று
	"எண்பது",   // எண்பத்து ஒன்று
	"தொன்னூறு", // தொன்னூற்று ஒன்று
}

var hundreds []string = []string{
	"நூறு",
	"இருனூற்று",
	"முன்னூற்று",
	"",
	"",
	"",
	"",
	"",
	"தொள்ளாயிரம்",
}

var thousands []string = []string{
	"ஆயிரம்", // ஆயிரத்து
	"இரண்டாயிரம்",
	"மூன்றாயிரம்",
	"",
	"",
	"",
	"",
	"",
	"",
}

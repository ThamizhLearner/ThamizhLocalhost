package main

import (
	"html/template"
	"net/http"

	script "github.com/ThamizhLearner/Thamizh"
)

type miscActivity struct{}

func (a miscActivity) GetID() string   { return "Misc" }
func (a miscActivity) GetDesc() string { return "Uncategorized cache" }
func (a miscActivity) Respond(w http.ResponseWriter, r *http.Request) {
	seed := struct {
		InfoTable   SimpleTable
		NerTable    SimpleTable
		NiraiTable  SimpleTable
		RhythmTable SimpleTable
		CVTable     SimpleTable
		VerbGraph   string
	}{finalTable2(), createNerTable(), createNiraiTable(), createRhythmTable(), compositeLetters(), createVerbGraph()}

	var tmpl = template.Must(template.ParseFiles("tmpls/index.tmpl", "tmpls/misc.tmpl"))
	tmpl.Execute(w, seed)
}

func compositeLetters() SimpleTable {
	table := SimpleTable{
		Title: "Compound letter matrix",
		Rows:  19, Columns: 13,
		ColInfoList: nil, // Note: Characteristic of 2D table!
		Cells:       make([][]string, 19),
	}

	consonantLetters := [18]script.Letter{
		script.MustLetterFrom("க்"), script.MustLetterFrom("ங்"), script.MustLetterFrom("ச்"),
		script.MustLetterFrom("ஞ்"), script.MustLetterFrom("ட்"), script.MustLetterFrom("ண்"),
		script.MustLetterFrom("த்"), script.MustLetterFrom("ந்"), script.MustLetterFrom("ப்"),
		script.MustLetterFrom("ம்"), script.MustLetterFrom("ய்"), script.MustLetterFrom("ர்"),
		script.MustLetterFrom("ல்"), script.MustLetterFrom("வ்"), script.MustLetterFrom("ழ்"),
		script.MustLetterFrom("ள்"), script.MustLetterFrom("ற்"), script.MustLetterFrom("ன்"),
	}
	vowelLetters := [12]script.Letter{
		script.MustLetterFrom("அ"), script.MustLetterFrom("ஆ"), script.MustLetterFrom("இ"),
		script.MustLetterFrom("ஈ"), script.MustLetterFrom("உ"), script.MustLetterFrom("ஊ"),
		script.MustLetterFrom("எ"), script.MustLetterFrom("ஏ"), script.MustLetterFrom("ஐ"),
		script.MustLetterFrom("ஒ"), script.MustLetterFrom("ஓ"), script.MustLetterFrom("ஔ"),
	}

	// Initialize the table matrix
	for r := range table.Rows {
		table.Cells[r] = make([]string, table.Columns)
	}

	for r := range table.Rows {
		row := table.Cells[r]
		for c := range table.Columns {
			if r == 0 && c == 0 {
				continue
			} else if r == 0 {
				// V zone - Row-wise
				row[c] = vowelLetters[c-1].String()
			} else if c == 0 {
				// C zone - Column-wise
				row[0] = consonantLetters[r-1].String()
			} else {
				// CV zone
				cv := consonantLetters[r-1].MustJoinCV(vowelLetters[c-1])
				row[c] = cv.String()
			}
		}
	}
	return table
}

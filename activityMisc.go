package main

import (
	"html/template"
	"maps"
	"net/http"
	"slices"
	"sort"

	script "github.com/ThamizhLearner/Thamizh"
	kural2 "github.com/ThamizhLearner/ThamizhLocalhost/kural"
)

type miscActivity struct{}

func (a miscActivity) GetID() string   { return "Misc" }
func (a miscActivity) GetDesc() string { return "Uncategorized cache" }
func (a miscActivity) Respond(w http.ResponseWriter, r *http.Request) {
	seed := struct {
		InfoTable        SimpleTable
		NerTable         SimpleTable
		NiraiTable       SimpleTable
		Class1RhymeTable SimpleTable
		Class2RhymeTable SimpleTable
		Class3RhymeTable SimpleTable
		Class4RhymeTable SimpleTable
		CVTable          SimpleTable
		VerbGraph        string
	}{finalTable2(), createNerTable(), createNiraiTable(),
		class1RhymeTable(), class2RhymeTable(), class3RhymeTable(), class4RhymeTable(),
		compositeLetters(), createVerbGraph()}

	var tmpl = template.Must(template.ParseFiles("tmpls/index.tmpl", "tmpls/misc.tmpl"))
	tmpl.Execute(w, seed)
}

func class1RhymeTable() SimpleTable {
	return rhymeTable("ஓரசைச்சீர் (ஓர் + அசை + சீர்)", kural2.GetSortedClass1RhymeMap())
}

func class2RhymeTable() SimpleTable {
	return rhymeTable("ஈரசைச்சீர் (ஈர் + அசை + சீர்)", kural2.GetSortedClass2RhymeMap())
}

func class3RhymeTable() SimpleTable {
	return rhymeTable("மூவசைச்சீர் (மூ + அசை + சீர்)", kural2.GetSortedClass3RhymeMap())
}

func class4RhymeTable() SimpleTable {
	return rhymeTable("நாலசைச்சீர் (நால் + அசை + சீர்)", kural2.GetSortedClass4RhymeMap())
}

func rhymeTable(title string, ceerMap map[string]kural2.Ceer) SimpleTable {
	var t = SimpleTable{
		Title:       title,
		ColInfoList: []ColInfo{{"வாய்ப்பாடு", 1}, {"Syllabified", 1}, {"Structure", 1}, {"அசை spans", 1}},
		Cells:       make([][]string, len(ceerMap)),
	}
	keys := slices.Collect(maps.Keys(ceerMap))

	// Sort the keys (which are encoded as sequence of indices into [நேர், நிரை, நேர்பு, நிரைபு])
	sort.Slice(keys, func(i, j int) bool {
		a, b := reversed(keys[i]), reversed(keys[j]) // Reversed to match the way அசை sequences are ordered!
		la, lb := len(a), len(b)
		if la == lb {
			return a < b
		}
		return la < lb
	})

	for r, k := range keys {
		row := make([]string, 4)
		t.Cells[r] = row
		rhythm := ceerMap[k]
		row[0] = rhythm.UStr
		row[1], _ = script.SyllabifiedUStr(script.MustLetterSeqFrom(rhythm.UStr), "-")
		row[2] = kural2.AcaaiSeq2UStr(rhythm.Captures, false)
		row[3] = kural2.AcaaiFragSeq2UStr(rhythm.Captures)
	}

	return t
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

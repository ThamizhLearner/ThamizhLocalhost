package main

import (
	"bufio"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"

	script "github.com/ThamizhLearner/Thamizh"
	kural2 "github.com/ThamizhLearner/ThamizhLocalhost/kural"
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

// Kural display source
type kural struct {
	Num         int
	Line, Line2 string
	Syl, Syl2   string
	RhythmTable SimpleTable
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
		line := strings.Trim(scanner.Text(), " ")
		if !scanner.Scan() {
			panic("Unexpected EOF")
		}
		line2 := strings.Trim(scanner.Text(), " ")
		syl, syl2 := SyllabifyLong(line), SyllabifyLong(line2)
		kuralStr := line + " " + line2
		kurals = append(kurals, kural{num, line, line2, syl, syl2, createThirukkuralRythmTable(kuralStr)})
		num++

		// Just handle first 10 kurals!
		if num == 11 {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return kurals
}

// Syllabifies spaced string!
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

var rhythmBaseMap = kural2.GetRhythmBaseMap()

func createThirukkuralRythmTable(k string) SimpleTable {
	var t = SimpleTable{
		Title:       "அலகிட்டு வாய்ப்பாடு",
		ColInfoList: []ColInfo{{"சீர்", 1}, {"அசை", 1}, {"வாய்ப்பாடு", 1}},
		Cells:       make([][]string, 7),
	}
	// fmt.Println(k)
	strs := strings.Split(k, " ")
	if len(strs) != 7 {
		fmt.Println(k)
		panic("Expected 7 சீர்கள்")
	}
	for i, str := range strs {
		if i == 7 {
			break
		}
		row := make([]string, 3)
		t.Cells[i] = row
		ok := !strings.Contains(str, "ஃ")
		if !ok {
			row[0] = str
			row[1] = "Failed"
			row[2] = "Failed"
			continue
		}
		ls := script.MustLetterSeqFrom(str)
		captures := kural2.CaptureRhythm(ls, i == 6)
		key := kural2.CreateKey(captures)
		rhythm := rhythmBaseMap[key]
		row[0] = str
		row[1] = kural2.GetRhythmBeats(captures)
		row[2] = rhythm.Value
	}

	return t
}

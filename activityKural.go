package main

import (
	"bufio"
	"html/template"
	"maps"
	"net/http"
	"os"
	"slices"
	"sort"
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
		syl, syl2 := SyllabifySpaced(line), SyllabifySpaced(line2)
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
func SyllabifySpaced(str string) string {
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

// Gets given குறள் text's அலகிட்டு வாய்ப்பாடு table
func createThirukkuralRythmTable(k string) SimpleTable {
	var t = SimpleTable{
		Title:       "அலகிட்டு வாய்ப்பாடு",
		ColInfoList: []ColInfo{{"சீர்", 2}, {"அசை", 1}, {"வாய்ப்பாடு", 2}},
		Cells:       make([][]string, 7),
	}
	strs := strings.Split(k, " ")
	for i, str := range strs {
		row := make([]string, 5)
		t.Cells[i] = row
		row[0] = str
		if strings.Contains(str, "ஃ") {
			row[1] = "Failed"
			row[2] = "Failed"
			row[3] = "Failed"
			row[4] = "Failed"
			continue
		}
		ls := script.MustLetterSeqFrom(str)
		captures := kural2.CaptureRhythm(ls, i == 6)
		key := kural2.CreateKey(captures)
		baseRhythm := rhythmBaseMap[key]
		row[1] = kural2.RhythmBreakup(captures)
		row[2] = kural2.GetRhythmBeats(captures, false)
		row[3] = baseRhythm.UStr
		row[4] = kural2.RhythmBreakup(baseRhythm.Captures)
	}

	return t
}

// நேரசை table
func createNerTable() SimpleTable {
	var t = SimpleTable{
		Title:       "நேரசை (நேர் + அசை) definition",
		ColInfoList: []ColInfo{{"Pattern", 1}, {"Example", 1}, {"Insight", 1}, {"Remarks", 1}},
		Cells:       make([][]string, 4),
	}
	strs := [4]string{
		"தனிக்குறில்|ப|ப (குறில்)|நேரசை(1)",
		"குறில் + ஒற்று|பல்|பல் (குறில் + ஒற்று)|நேரசை(2)",
		"தனிநெடில்|கா|கா (நெடில்)|நேரசை(3)",
		"நெடில் + ஒற்று|கால்|கால் (நெடில் + ஒற்று)|நேரசை(4)",
	}
	for r := range 4 {
		row := make([]string, 4)
		t.Cells[r] = row
		count := copy(row, strings.Split(strs[r], "|"))
		if count != 4 {
			panic("Coding error")
		}
	}
	return t
}

// நிரையசை table
func createNiraiTable() SimpleTable {
	var t = SimpleTable{
		Title:       "நிரையசை (நிரை + அசை) definition",
		ColInfoList: []ColInfo{{"Pattern", 1}, {"Example", 1}, {"Insight", 1}, {"Remarks", 1}},
		Cells:       make([][]string, 4),
	}
	strs := [4]string{
		"தனிக்குறில் + தனிக்குறில்|அணி (அ / ணி)|அ (குறில்) + ணி (குறில்)|நேரசை(1) + நேரசை(1)",
		"தனிக்குறில் + (குறில் + ஒற்று)|அவர் (அ / வர்)|அ (குறில்) + வர் (குறில் + ஒற்று)|நேரசை(1) + நேரசை(2)",
		"தனிக்குறில் + தனிநெடில்|தொழா (தொ / ழா)|தொ (குறில்) + ழா (நெடில்)|நேரசை(1) + நேரசை(3)",
		"தனிக்குறில் + (நெடில் + ஒற்று)|கினான் (கி / னான்)|கி (குறில்) + னான் (நெடில் + ஒற்று)|நேரசை(1) + நேரசை(4)",
	}
	for r := range 4 {
		row := make([]string, 4)
		t.Cells[r] = row
		count := copy(row, strings.Split(strs[r], "|"))
		if count != 4 {
			panic("Coding error")
		}
	}
	return t
}

// வாய்ப்பாடு table
func createRhythmTable() SimpleTable {
	var t = SimpleTable{
		Title:       "வாய்ப்பாடு definitions (list)",
		ColInfoList: []ColInfo{{"வாய்ப்பாடு", 1}, {"Syllabified", 1}, {"Structure", 1}, {"அசை spans", 1}},
		Cells:       make([][]string, len(rhythmBaseMap)),
	}
	keys := slices.Collect(maps.Keys(rhythmBaseMap))

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
		rhythm := rhythmBaseMap[k]
		row[0] = rhythm.UStr
		row[1], _ = script.SyllabifiedUStr(script.MustLetterSeqFrom(rhythm.UStr), "-")
		row[2] = kural2.GetRhythmBeats(rhythm.Captures, false)
		row[3] = kural2.RhythmBreakup(rhythm.Captures)
	}

	return t
}

// reverses the given string
func reversed(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		// swap the letters of the string,
		// like first with last and so on.
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

package main

import (
	"fmt"
)

// Tense composition expression template

// Classified verb tense template generator types
var infixesList [][3]string = [][3]string{
	{"த்த்", "க்", "ப்ப்"}, // Type A
	{"ந்த்", "", "வ்"},     // Type B
	{"இன்", "", "வ்"},      // Type C
	{"ந்த்", "க்", "ப்ப்"}, // Type D
}

var tenseDivs = [3]string{"இறந்த காலம்", "நிகழ் காலம்", "எதிர் காலம்"}
var coreTenseSfxs = [10]string{
	"ஏன்",
	"ஓம்",
	"ஆய்",
	"ஈர்கள்",
	"ஆன்",
	"ஆள்",
	"ஆர்",
	"ஆர்கள்",
	"(todo)",
	"(todo)",
}

// Tense form decomposition illustration (for Type A tense composition)
//
// Columns (5): Pronoun, Past, Present, Present2, Future
func decompTableA(vroot string) SimpleTable {
	return decompTable(vroot, 0)
}

// Tense form decomposition illustration (for Type B tense composition)
//
// Columns (5): Pronoun, Past, Present, Present2, Future
func decompTableB(vroot string) SimpleTable {
	return decompTable(vroot, 1)
}

// Tense form decomposition illustration
//
// Columns (5): Pronoun, Past, Present, Present2, Future
func decompTable(vroot string, infixesIdx int) SimpleTable {
	infixes := infixesList[infixesIdx]
	var coreTenseInserts = [4]string{"", "கிற்", "கின்ற்", ""}
	table := SimpleTable{
		Title: fmt.Sprintf("Tense pattern format (Type: %c)", "ABCD"[infixesIdx]),
		Rows:  10, Columns: 4,
		ColInfoList: []ColInfo{{"Pronoun", 1}, {tenseDivs[0], 1}, {tenseDivs[1], 2}, {tenseDivs[2], 1}},
		Cells:       make([][]string, 10),
	}
	for r := range 10 {
		pronounInfo := gPronounInfoList[r]
		row := make([]string, 5)
		row[0] = pronounInfo.Names[0]

		if r >= 8 { // TODO Handle tense for அஃறிணை
			for c := range 4 {
				row[1+c] = "(todo)"
			}
			table.Cells[r] = row
			continue
		}

		for c := range 4 {
			tenseInsert := ""
			switch c {
			case 0:
				tenseInsert = infixes[0]
			case 3:
				tenseInsert = infixes[2]
			default:
				tenseInsert = infixes[1]
				if len(tenseInsert) == 0 {
					tenseInsert = coreTenseInserts[c]
				} else {
					tenseInsert = tenseInsert + "+" + coreTenseInserts[c]
				}
			}
			row[1+c] = fmt.Sprintf("(%s)%s|%s", vroot, tenseInsert, coreTenseSfxs[r])
		}
		table.Cells[r] = row
	}
	return table
}

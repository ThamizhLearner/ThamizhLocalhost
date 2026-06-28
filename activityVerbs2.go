package main

import (
	"fmt"
	"strings"
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
	var coreTenseInserts = [4]string{"", "கிற", "கின்ற", ""}
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

// Verb classification graph
func createVerbGraph() string {
	// Graph driven by Mermaid (https://mermaid.js.org/)
	sb := strings.Builder{}
	sb.WriteString("graph LR\n")
	terms := []string{
		"A(வினைச்சொல்)",
		"B(\"வினைமுற்று\n(வினை + முற்று)\")",
		"C(\"வினையெச்சம்\n(வினை + எச்சம்)\")",
		"D(\"பெயரெச்சம்\n(பெயர் + எச்சம்)\")",
		"E(தெரிநிலை வினைமுற்று)", "F(குறிப்பு வினைமுற்று)",
		"G(தெரிநிலை வினையெச்சம்)", "H(குறிப்பு வினையெச்சம்)",
		"I(தெரிநிலை பெயரெச்சம்)", "J(குறிப்புப் பெயரெச்சம்)",
	}
	// Define the nodes
	for _, t := range terms {
		sb.WriteString(t)
		sb.WriteString("\n")
	}
	// Connect the nodes
	links := []string{
		"A --> B", "A --> C", "A --> D",
		"B --> E", "B --> F",
		"C --> G", "C --> H",
		"D --> I", "D --> J",
	}
	for _, l := range links {
		sb.WriteString(l)
		sb.WriteString("\n")
	}

	return sb.String()
}

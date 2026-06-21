package main

import (
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strings"

	script "github.com/ThamizhLearner/Thamizh"
)

type verbsActivity struct{}

func (a verbsActivity) GetID() string   { return `Verbs` }
func (a verbsActivity) GetDesc() string { return `Verb forms` }
func (a verbsActivity) Respond(w http.ResponseWriter, r *http.Request) {
	// Need the list of verbs
	// Determine the verb to present. [Default: Pick first one from the list, else user selected verb]

	if len(verbItems) == 0 {
		fmt.Fprintf(w, "Verb DB is empty")
		return
	}

	reqVerb := strings.TrimSpace(r.FormValue(`verb`))
	verbIdx := -1
	if len(reqVerb) == 0 {
		verbIdx = 0 // Pick the first one from the verb list.
	} else {
		// Verify the given verb is in the list
		for i, v := range verbItems {
			if v.Name == reqVerb {
				verbIdx = i
				break
			}
		}
		if verbIdx == -1 {
			fmt.Fprint(w, "Specified verb not found")
			return
		}
	}

	var verbUrls []Hyperlink
	for _, v := range verbItems {
		verbUrls = append(verbUrls, Hyperlink{
			Name: fmt.Sprintf("%s (%s)", v.Root, v.Name),
			Url:  template.URL(`/verbs?verb=` + v.Name),
		})
	}

	verbItem := verbItems[verbIdx]
	tenseForms := outTenseFormATable(script.MustLetterSeqFrom(verbItem.Root))
	seed := struct {
		Verb           string
		VerbalNoun     string
		VerbLinks      []Hyperlink
		RefTable       SimpleTable
		PhraseTable    SimpleTable
		SylPhraseTable SimpleTable
		Table          SimpleTable
		DecompTable    SimpleTable
	}{
		verbItem.Root, verbItem.Name, verbUrls, tenseQRefTable(tenseForms),
		tenseQPhraseTable(tenseForms, false), tenseQPhraseTable(tenseForms, true),
		finalTable2(), decompTable(verbItem.Root),
	}

	tmpl, err := template.ParseFiles(`tmpls/index.tmpl`, `tmpls/verb.tmpl`)
	if err != nil {
		panic(err)
	}
	tmpl.Execute(w, seed)
}

// Simple rendered table layout.
//
// Supports column-header merge (span extension)
type SimpleTable struct {
	Title         string
	Rows, Columns int
	ColInfoList   []ColInfo
	Cells         [][]string
}

type ColInfo struct {
	Header string
	Span   int
}

// Verb picker list item
type verbInfo struct {
	Root string
	Name string
}

var verbItems []verbInfo = getVerbItems()

// Get list of available verbs to choose from!
func getVerbItems() []verbInfo {
	verbItems := []verbInfo{
		{"அசை", "அசைத்தல்"},
		{"கடி", "கடித்தல்"},
		{"அடை", "அடைத்தல்"},
		{"அடி", "அடித்தல்"},
		{"அமை", "அமைத்தல்"},
		{"அவிழ்", "அவிழ்த்தல்"},
		{"அழி", "அழித்தல்"},
		{"அறு", "அறுத்தல்"},
		{"இடி", "இடித்தல்"},
		{"இழு", "இழுத்தல்"},
		{"இனி", "இனித்தல்"},
		{"உடை", "உடைத்தல்"},
		{"உதை", "உதைத்தல்"},
		{"எடு", "எடுத்தல்"},
		{"எரி", "எரித்தல்"},
		{"ஒடி", "ஒடித்தல்"},
		{"ஒழி", "ஒழித்தல்"},
		{"கரி", "கரித்தல்"},
		{"கரை", "கரைத்தல்"},
		{"கலை", "கலைத்தல்"},
		{"கவனி", "கவனித்தல்"},
		{"கவிழ்", "கவிழ்த்தல்"},
		{"கழி", "கழித்தல்"},
		{"கிடை", "கிடைத்தல்"},
		{"கிழி", "கிழித்தல்"},
	}

	// Sort!
	sort.Slice(verbItems, func(i, j int) bool {
		return verbItems[i].Name < verbItems[j].Name
	})

	return verbItems
}

// Columns (5): Pronoun, Past, Present, Present2, Future
func tenseQRefTable(tenseTable [8][4]string) SimpleTable {
	table := SimpleTable{
		Title: "Verbal tense forms",
		Rows:  8, Columns: 9,
		ColInfoList: []ColInfo{{"Pronoun", 1}, {"இறந்த காலம்", 1}, {"நிகழ் காலம்", 2}, {"எதிர் காலம்", 1}},
		Cells:       make([][]string, 10),
	}
	for r := range 10 {
		pronounInfo := gPronounInfoList[r]
		row := make([]string, 5)
		row[0] = pronounInfo.Names[0]

		if r >= 8 {
			for c := range 4 {
				row[1+c] = "(todo)"
			}
			table.Cells[r] = row
			continue
		}

		for c := range 4 {
			row[1+c] = tenseTable[r][c]
		}
		table.Cells[r] = row
	}
	return table
}

// Columns (4): Pronoun  Past, Pronoun  Present, Pronoun Present2, Pronoun Future
func tenseQPhraseTable(tenseTable [8][4]string, iSyllabified bool) SimpleTable {
	table := SimpleTable{
		Title: "Compact view",
		Rows:  8, Columns: 4,
		ColInfoList: []ColInfo{{"இறந்த காலம்", 1}, {"நிகழ் காலம்", 2}, {"எதிர் காலம்", 1}},
		Cells:       make([][]string, 8),
	}
	if iSyllabified {
		table.Title = "Syllabified (pronunciation hint) view"
	}
	for r := range 8 {
		pronounInfo := gPronounInfoList[r]
		row := make([]string, 4)
		pronoun := pronounInfo.Names[0]
		if iSyllabified {
			pronoun, _ = script.SyllabifiedUStr(script.MustLetterSeqFrom(pronoun), "|")
			pronoun = fmt.Sprintf("(%s)", pronoun)
		}

		for c := range 4 {
			cell := tenseTable[r][c]
			if iSyllabified {
				cell, _ = script.SyllabifiedUStr(script.MustLetterSeqFrom(cell), "|")
				cell = fmt.Sprintf("(%s)", cell)
			}
			row[c] = fmt.Sprintf("%s %s", pronoun, cell)
		}
		table.Cells[r] = row
	}
	return table
}

// Columns (5): Pronoun, Idam, Paal, Thinai, En.
func finalTable2() SimpleTable {
	table := SimpleTable{
		Title: "Classification (திணை | பால் | எண் | இடம்)",
		Rows:  8, Columns: 5,
		ColInfoList: []ColInfo{{"Pronouns", 1}, {"திணை", 1}, {"பால்", 1}, {"எண்", 1}, {"இடம்", 1}},
		Cells:       make([][]string, 10),
	}
	for r := range 10 {
		pronounInfo := gPronounInfoList[r]
		row := make([]string, 5)
		row[0] = strings.Join(pronounInfo.Names, ", ")
		row[1] = pronounInfo.Thinai
		row[2] = pronounInfo.Paal
		row[3] = pronounInfo.En
		row[4] = pronounInfo.Idam
		table.Cells[r] = row
	}
	return table
}

// Columns (5): Pronoun, Past, Present, Present2, Future
func decompTable(vroot string) SimpleTable {
	var suffixGrid = [8][4]string{
		{`த்த் # ஏன்`, `க் + கிற் # ஏன்`, `க் + கின்ற் # ஏன்`, `ப்ப் # ஏன்`},
		{`த்த் # ஓம்`, `க் + கிற் # ஓம்`, `க் + கின்ற் # ஓம்`, `ப்ப் # ஓம்`},
		{`த்த் # ஆய்`, `க் + கிற் # ஆய்`, `க் + கின்ற் # ஆய்`, `ப்ப் # ஆய்`},
		{`த்த் # ஈர்கள்`, `க் + கிற் # ஈர்கள்`, `க் + கின்ற் # ஈர்கள்`, `ப்ப் # ஈர்கள்`},
		{`த்த் # ஆன்`, `க் + கிற் # ஆன்`, `க் + கின்ற் # ஆன்`, `ப்ப் # ஆன்`},
		{`த்த் # ஆள்`, `க் + கிற் # ஆள்`, `க் + கின்ற் # ஆள்`, `ப்ப் # ஆள்`},
		{`த்த் # ஆர்`, `க் + கிற் # ஆர்`, `க் + கின்ற் # ஆர்`, `ப்ப் # ஆர்`},
		{`த்த் # ஆர்கள்`, `க் + கிற் # ஆர்கள்`, `க் + கின்ற் # ஆர்கள்`, `ப்ப் # ஆர்கள்`},
	}

	table := SimpleTable{
		Title: "Tense pattern format (Type: (todo))",
		Rows:  8, Columns: 9,
		ColInfoList: []ColInfo{{"Pronoun", 1}, {"இறந்த காலம்", 1}, {"நிகழ் காலம்", 2}, {"எதிர் காலம்", 1}},
		Cells:       make([][]string, 10),
	}
	for r := range 10 {
		pronounInfo := gPronounInfoList[r]
		row := make([]string, 5)
		row[0] = pronounInfo.Names[0]

		if r >= 8 {
			for c := range 4 {
				row[1+c] = "(todo)"
			}
			table.Cells[r] = row
			continue
		}

		for c := range 4 {
			row[1+c] = strings.Replace(fmt.Sprintf("(%s)%s", vroot, suffixGrid[r][c]), " ", "", -1)
			row[1+c] = strings.Replace(row[1+c], "#", "|", -1)
		}
		table.Cells[r] = row
	}
	return table
}

// Final tense forms table, for the given root verb.
func outTenseFormATable(vroot script.LetterSeq) [8][4]string {
	suffixes := TenseFormASuffixes
	var table [8][4]string
	for r := range 8 {
		for c := range 4 {
			suffix := suffixes[r][c]
			if suffix.First().IsV() {
				seq, ok := script.VSuffixAppended(vroot, suffix)
				if !ok {
					panic("Coding error")
				}
				table[r][c] = seq.String()
				continue
			}
			table[r][c] = script.SuffixAppended(vroot, suffix).String()
		}
	}
	return table
}

func createTenseFormASuffixes() [8][4]script.LetterSeq {
	return createTenseSuffixes(script.MustLetterSeqFrom("த்த்"), script.MustLetterSeqFrom("க்"), script.MustLetterSeqFrom("ப்ப்"))
}

// Constructs final suffix mask.
//
// Form expansion format: ◌த்த்!, ◌க்◌, ◌ப்ப்!
func createTenseSuffixes(past, present, future script.LetterSeq) [8][4]script.LetterSeq {
	var table [8][4]script.LetterSeq
	tenses := [4]script.LetterSeq{past, present, present, future} // Note: 2 versions of present tense!
	for r := range 8 {
		for c := range 4 {
			suffix := gCoreTenseSuffixes[r][c]
			if suffix.First().IsV() {
				seq, ok := script.VSuffixAppended(tenses[c], suffix)
				if !ok {
					panic("Coding error")
				}
				table[r][c] = seq
				continue
			}
			table[r][c] = script.SuffixAppended(tenses[c], suffix)
		}
	}
	return table
}

type Hyperlink struct {
	Name string
	Url  template.URL
}

// Internal (mutable!)
var gCoreTenseSuffixes [8][4]script.LetterSeq = createCoreTenseSuffixes()
var TenseFormASuffixes [8][4]script.LetterSeq = createTenseFormASuffixes()

// Core suffix mask for final tense form compositions.
func createCoreTenseSuffixes() [8][4]script.LetterSeq {
	var ustrTenseForms = [8][4]string{
		{`ஏன்`, `கிறேன்`, `கின்றேன்`, `ஏன்`},
		{`ஓம்`, `கிறோம்`, `கின்றோம்`, `ஓம்`},
		{`ஆய்`, `கிறாய்`, `கின்றாய்`, `ஆய்`},
		{`ஈர்கள்`, `கிறீர்கள்`, `கின்றீர்கள்`, `ஈர்கள்`},
		{`ஆன்`, `கிறான்`, `கின்றான்`, `ஆன்`},
		{`ஆள்`, `கிறாள்`, `கின்றாள்`, `ஆள்`},
		{`ஆர்`, `கிறார்`, `கின்றார்`, `ஆர்`},
		{`ஆர்கள்`, `கிறார்கள்`, `கின்றார்கள்`, `ஆர்கள்`},
	}
	var tenseForms [8][4]script.LetterSeq
	for rIdx := range 8 {
		for cIdx := range 4 {
			tenseForms[rIdx][cIdx] = script.MustLetterSeqFrom(ustrTenseForms[rIdx][cIdx])
		}
	}
	return tenseForms
}

// திணை பால் எண் இடம்
// இடம்: தன்மை, முன்னிலை, படர்க்கை
// எண்: ஒருமை பன்மை
// திணை: உயர்திணை அஃறிணை
// பால்: ஆண்பால், பெண்பால், பலர்பால், ஒன்றன்பால், பலவின்பால்
var gPronounInfoList = []pronounInfo{
	{Id: "1s", Names: []string{"நான்"}, Idam: "தன்மை", Paal: "ஆண்பால், பெண்பால்", Thinai: "உயர்திணை", En: "ஒருமை"},
	{Id: "1p", Names: []string{"நாங்கள்", "நாம்"}, Idam: "தன்மை", Paal: "பலர்பால்", Thinai: "உயர்திணை", En: "பன்மை"},
	{Id: "2s", Names: []string{"நீ", "நீம்"}, Idam: "முன்னிலை", Paal: "ஆண்பால், பெண்பால்", Thinai: "உயர்திணை", En: "ஒருமை"},
	{Id: "2p", Names: []string{"நீங்கள்"}, Idam: "முன்னிலை", Paal: "பலர்பால்", Thinai: "உயர்திணை", En: "பன்மை"},
	{Id: "3sm", Names: []string{"அவன்", "இவன்", "எவன்"}, Idam: "படர்க்கை", Paal: "ஆண்பால்", Thinai: "உயர்திணை", En: "ஒருமை"},
	{Id: "3sf", Names: []string{"அவள்", "இவள்", "எவள்"}, Idam: "படர்க்கை", Paal: "பெண்பால்", Thinai: "உயர்திணை", En: "ஒருமை"},
	{Id: "3sh", Names: []string{"அவர்", "இவர்", "எவர்"}, Idam: "படர்க்கை", Paal: "ஆண்பால், பெண்பால்", Thinai: "உயர்திணை", En: "ஒருமை"},
	{Id: "3ph", Names: []string{"அவர்கள்", "இவர்கள்", "எவர்கள்"}, Idam: "படர்க்கை", Paal: "பலர்பால்", Thinai: "உயர்திணை", En: "பன்மை"},
	{Id: "3sn", Names: []string{"அது", "இது", "எது"}, Idam: "படர்க்கை", Paal: "ஒன்றன்பால்", Thinai: "அஃறிணை", En: "ஒருமை"},
	{Id: "3pn", Names: []string{"அவை", "இவை", "எவை"}, Idam: "படர்க்கை", Paal: "பலவின்பால்", Thinai: "அஃறிணை", En: "பன்மை"},
}

type pronounInfo struct {
	Id     string
	Names  []string
	Idam   string
	Paal   string
	Thinai string
	En     string
}

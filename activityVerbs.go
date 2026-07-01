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
		reqVSeq, ok := script.LetterSeqFrom(reqVerb)
		if !ok {
			fmt.Fprint(w, "Not a valid verb")
			return
		}
		// Verify the given verb is in the list
		for i, v := range verbItems {
			if v.Name.Equals(reqVSeq) {
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
		vname := v.Name.String()
		verbUrls = append(verbUrls, Hyperlink{
			Name: fmt.Sprintf("%s (%s)", v.Root, vname),
			// Name: fmt.Sprintf("%s (%s) - %c", v.Root, vname, v.Type),
			Url: template.URL(`/verbs?verb=` + vname),
		})
	}

	verbItem := verbItems[verbIdx]
	var tenseForms [8][4]string
	var decompTable SimpleTable
	switch verbItem.Type {
	case 'A':
		tenseForms = outTenseTable(verbItem.Root, gTenseFormASuffixes)
		decompTable = decompTableA(verbItem.Root.String())
	case 'B':
		tenseForms = outTenseTable(verbItem.Root, gTenseFormBSuffixes)
		decompTable = decompTableB(verbItem.Root.String())
	case 'C':
		tenseForms = outTenseTable(verbItem.Root, gTenseFormCSuffixes)
		decompTable = decompTableC(verbItem.Root.String())
	case 'D':
		tenseForms = outTenseTable(verbItem.Root, gTenseFormDSuffixes)
		decompTable = decompTableD(verbItem.Root.String())
	default:
		panic("Coding error: Unsupported verb tense type")
	}

	seed := struct {
		Verb           script.LetterSeq
		VerbalNoun     script.LetterSeq
		VerbLinks      []Hyperlink
		RefTable       SimpleTable
		PhraseTable    SimpleTable
		SylPhraseTable SimpleTable
		DecompTable    SimpleTable
	}{
		verbItem.Root, verbItem.Name, verbUrls, tenseQRefTable(tenseForms),
		tenseQPhraseTable(tenseForms, false), tenseQPhraseTable(tenseForms, true),
		decompTable,
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
	Root script.LetterSeq
	Name script.LetterSeq
	Type rune // Tense composition type
}

var verbItems []verbInfo = getVerbItems()

// Type A - அடைத்தல் (அடை + த்தல்) ["த்த்", "க்", "ப்ப்"]
var verbRootsA []string = []string{
	"அசை",
	"கடி",
	"அடை",
	"அடி",
	"அமை",
	"அவிழ்",
	"அழி",
	"அறு",
	"இடி",
	"இழு",
	"இனி",
	"உடை",
	"உதை",
	"எடு",
	"எரி",
	"ஒடி",
	"ஒழி",
	"கரி",
	"கரை",
	"கலை",
	"கவனி",
	"கவிழ்",
	"கழி",
	"கிடை",
	"கிழி",
}

// Type B - அடைதல் (அடை + தல்) ["ந்த்", "", "வ்"]
var verbRootsB []string = []string{
	"சாய்",
	"குனி",
	"குளிர்",
	"அசை",
	"அடை",
	"அமை",
	"அழி",
	"அறி",
	"அறை",
	"எரி",
}

// Type C -தோன்றுதல் (தோன்று + தல்) ["இன்", "", "வ்"]
var verbRootsC []string = []string{
	"தோன்று",
	"தேடு",
}

// Type D - திறத்தல் (திற + த்தல்) ["ந்த்", "க்", "ப்ப்"]
var verbRootsD []string = []string{
	"திற",
	"சும",
}

// Get list of available verbs to choose from!
func getVerbItems() []verbInfo {
	var verbInfoItems []verbInfo

	fix := script.MustLetterSeqFrom("த்தல்")
	for _, vroot := range verbRootsA {
		r := script.MustLetterSeqFrom(vroot)
		n := r.Appended(fix)
		verbInfoItems = append(verbInfoItems, verbInfo{Root: r, Name: n, Type: 'A'})
	}

	fix = script.MustLetterSeqFrom("தல்")
	for _, vroot := range verbRootsB {
		r := script.MustLetterSeqFrom(vroot)
		n := r.Appended(fix)
		verbInfoItems = append(verbInfoItems, verbInfo{Root: r, Name: n, Type: 'B'})
	}

	fix = script.MustLetterSeqFrom("தல்")
	for _, vroot := range verbRootsC {
		r := script.MustLetterSeqFrom(vroot)
		n := r.Appended(fix)
		verbInfoItems = append(verbInfoItems, verbInfo{Root: r, Name: n, Type: 'C'})
	}

	fix = script.MustLetterSeqFrom("த்தல்")
	for _, vroot := range verbRootsD {
		r := script.MustLetterSeqFrom(vroot)
		n := r.Appended(fix)
		verbInfoItems = append(verbInfoItems, verbInfo{Root: r, Name: n, Type: 'D'})
	}

	// Sort!
	sort.Slice(verbInfoItems, func(i, j int) bool {
		a, b := verbInfoItems[i], verbInfoItems[j]
		if a.Root.String() == b.Root.String() {
			return a.Name.String() < b.Name.String()
		}
		return a.Root.String() < b.Root.String()
	})

	return verbInfoItems
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
		table.Title = "Syllabified (quick pronunciation hint) view"
	}
	for r := range 8 {
		pronounInfo := gPronounInfoList[r]
		row := make([]string, 4)
		pronoun := pronounInfo.Names[0]
		if iSyllabified {
			pronoun, _ = script.SyllabifiedUStr(script.MustLetterSeqFrom(pronoun), "-")
			pronoun = fmt.Sprintf("%s | ", pronoun)
		}

		for c := range 4 {
			cell := tenseTable[r][c]
			if iSyllabified {
				cell, _ = script.SyllabifiedUStr(script.MustLetterSeqFrom(cell), "-")
				cell = fmt.Sprintf("%s", cell)
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
		row[0] = strings.Join(pronounInfo.Names, " | ")
		row[1] = pronounInfo.Thinai
		row[2] = pronounInfo.Paal
		row[3] = pronounInfo.En
		row[4] = pronounInfo.Idam
		table.Cells[r] = row
	}
	return table
}

// Final tense forms table, for the given root verb.
func outTenseTable(vroot script.LetterSeq, suffixTable [8][4]script.LetterSeq) [8][4]string {
	var table [8][4]string
	for r := range 8 {
		for c := range 4 {
			suffix := suffixTable[r][c]
			if suffix.First().IsV() {
				seq, ok := script.VSuffixAppended(vroot, suffix)
				if !ok {
					seq, ok = script.VSuffixVSubst(vroot, suffix)
					if !ok {
						panic("Coding error")
					}
				}
				table[r][c] = seq.String()
				continue
			}
			table[r][c] = script.SuffixAppended(vroot, suffix).String()
		}
	}
	return table
}

func createTenseFormASuffixes() [8][4]script.LetterSeq { // Type A {"த்த்", "க்", "ப்ப்"}
	return createTenseSuffixes(script.MustLetterSeqFrom("த்த்"), script.MustLetterSeqFrom("க்"), script.MustLetterSeqFrom("ப்ப்"))
}

func createTenseFormBSuffixes() [8][4]script.LetterSeq { // Type B {"ந்த்", "", "வ்"}
	return createTenseSuffixes(script.MustLetterSeqFrom("ந்த்"), script.LetterSeq{}, script.MustLetterSeqFrom("வ்"))
}

func createTenseFormCSuffixes() [8][4]script.LetterSeq { // Type C {"இன்", "", "வ்"}
	return createTenseSuffixes(script.MustLetterSeqFrom("இன்"), script.LetterSeq{}, script.MustLetterSeqFrom("வ்"))
}

func createTenseFormDSuffixes() [8][4]script.LetterSeq { // Type D {"ந்த்", "க்", "ப்ப்"}
	return createTenseSuffixes(script.MustLetterSeqFrom("ந்த்"), script.MustLetterSeqFrom("க்"), script.MustLetterSeqFrom("ப்ப்"))
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
var gTenseFormASuffixes [8][4]script.LetterSeq = createTenseFormASuffixes()
var gTenseFormBSuffixes [8][4]script.LetterSeq = createTenseFormBSuffixes()
var gTenseFormCSuffixes [8][4]script.LetterSeq = createTenseFormCSuffixes()
var gTenseFormDSuffixes [8][4]script.LetterSeq = createTenseFormDSuffixes()

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
	{Id: "1s", Names: []string{"நான்"}, Idam: "தன்மை", Paal: "ஆண்பால் | பெண்பால்", Thinai: "உயர்திணை", En: "ஒருமை"},
	{Id: "1p", Names: []string{"நாங்கள்", "நாம்"}, Idam: "தன்மை", Paal: "பலர்பால்", Thinai: "உயர்திணை", En: "பன்மை"},
	{Id: "2s", Names: []string{"நீ", "நீம்"}, Idam: "முன்னிலை", Paal: "ஆண்பால் | பெண்பால்", Thinai: "உயர்திணை", En: "ஒருமை"},
	{Id: "2p", Names: []string{"நீங்கள்"}, Idam: "முன்னிலை", Paal: "பலர்பால்", Thinai: "உயர்திணை", En: "பன்மை"},
	{Id: "3sm", Names: []string{"அவன்", "இவன்", "எவன்"}, Idam: "படர்க்கை", Paal: "ஆண்பால்", Thinai: "உயர்திணை", En: "ஒருமை"},
	{Id: "3sf", Names: []string{"அவள்", "இவள்", "எவள்"}, Idam: "படர்க்கை", Paal: "பெண்பால்", Thinai: "உயர்திணை", En: "ஒருமை"},
	{Id: "3sh", Names: []string{"அவர்", "இவர்", "எவர்"}, Idam: "படர்க்கை", Paal: "ஆண்பால் | பெண்பால்", Thinai: "உயர்திணை", En: "ஒருமை"},
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

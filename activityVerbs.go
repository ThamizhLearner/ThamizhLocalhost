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
			fmt.Fprint(w, "Specified verb found")
			return
		}
	}

	var verbUrls []Hyperlink
	for _, v := range verbItems {
		verbUrls = append(verbUrls, Hyperlink{
			Name: fmt.Sprintf("%s (%s)", v.Root, v.Name),
			Url:  `/verbs?verb=` + v.Name,
		})
	}

	verbItem := verbItems[verbIdx]
	seed := struct {
		Verb      string
		VerbLinks []Hyperlink
		Tenses    [8][3]string
	}{verbItem.Root, verbUrls, outTenseFormATable(script.MustLetterSeqFrom(verbItem.Root))}

	// https://www.w3schools.com/html/html_layout.asp

	var tmpl = template.Must(template.ParseFiles(`tmpls/index.tmpl`, `tmpls/verb.tmpl`))
	tmpl.Execute(w, seed)

}

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

func outTenseFormATable(vroot script.LetterSeq) [8][3]string {
	suffixes := TenseFormASuffixes
	var table [8][3]string
	for r := range 8 {
		for c := range 3 {
			suffix := suffixes[r][c]
			if suffix.First().IsV() {
				seq, _ := script.VSuffixAppended(vroot, suffix)
				table[r][c] = seq.String()
				continue
			}
			table[r][c] = script.SuffixAppended(vroot, suffix).String()
		}
	}
	return table
}

func createTenseFormASuffixes() [8][3]script.LetterSeq {
	return createTenseSuffixes(script.MustLetterSeqFrom("த்த்"), script.MustLetterSeqFrom("க்"), script.MustLetterSeqFrom("ப்ப்"))
}

// Form expansion format: ◌த்த்!, ◌க்◌, ◌ப்ப்!
func createTenseSuffixes(past, present, future script.LetterSeq) [8][3]script.LetterSeq {
	var table [8][3]script.LetterSeq
	tenses := [3]script.LetterSeq{past, present, future}
	for r := range 8 {
		for c := range 3 {
			suffix := coreTenseSuffixes[r][c]
			if suffix.First().IsV() {
				seq, _ := script.VSuffixAppended(tenses[c], suffix)
				table[r][c] = seq
				continue
			}
			table[r][c] = script.SuffixAppended(tenses[c], suffix)
		}
	}
	return table
}

type VerbFormDump struct {
	Verb      string // Verb being considered
	TenseRows [][]string
	AllVerbs  []Hyperlink // All the verbs to choose from, for the next iteration.
}

type Hyperlink struct {
	Name string
	Url  string
}

type Text struct {
	Value        string
	SyllabicForm string
}

// Internal (mutable!)
var coreTenseSuffixes [8][3]script.LetterSeq = createCoreTenseSuffixes()
var TenseFormASuffixes [8][3]script.LetterSeq = createTenseFormASuffixes()

func createCoreTenseSuffixes() [8][3]script.LetterSeq {
	var ustrTenseForms = [8][3]string{
		{`ஏன்`, `கிறேன்`, `ஏன்`},
		{`ஓம்`, `கிறோம்`, `ஓம்`},
		{`ஆய்`, `கிறாய்`, `ஆய்`},
		{`ஈர்கள்`, `கிறீர்கள்`, `ஈர்கள்`},
		{`ஆன்`, `கிறான்`, `ஆன்`},
		{`ஆள்`, `கிறாள்`, `ஆள்`},
		{`ஆர்`, `கிறார்`, `ஆர்`},
		{`ஆர்கள்`, `கிறார்கள்`, `ஆர்கள்`},
	}
	var tenseForms [8][3]script.LetterSeq
	for rIdx := range 8 {
		for cIdx := range 3 {
			tenseForms[rIdx][cIdx] = script.MustLetterSeqFrom(ustrTenseForms[rIdx][cIdx])
		}
	}
	return tenseForms
}

package kural2

import (
	"fmt"
	"strings"

	script "github.com/ThamizhLearner/Thamizh"
)

// Formula/pattern matching! Formula is a sequence of நேர் and நிரை.
// Pattern encoding: [நேர், நிரை, நேர்பு, நிரைபு]

// Rhyme
// தேமா => நேர், நேர்
// decompose தேமா => நேர்(தே), நேர்(மா)

// Note: We cannot handle 'ஃ' symbol yet. [Need to formulate a suitable hack]

// சீர் - seq of [நேர், நிரை, நேர்பு, நிரைபு]
type Ceer struct {
	UStr     string
	Captures []AcaaiFrag
}

// Accai-seq as formatted string
func AcaaiSeq2UStr(captures []AcaaiFrag, withCatpure bool) string {
	var strs []string
	for _, capture := range captures {
		str := acaaiNames[capture.acaaiId]
		if withCatpure { // Inject the captured fragment
			str += fmt.Sprintf(" {%s}", capture.ufrag)
		}
		strs = append(strs, str)
	}
	return strings.Join(strs, " | ")
}

// Acaai(Captured frag)-seq as formatted string
func AcaaiFragSeq2UStr(captures []AcaaiFrag) string {
	var strs []string
	for _, capture := range captures {
		strs = append(strs, capture.ufrag)
	}
	return strings.Join(strs, "/")
}

var rhymeClass1Names = []string{
	"நாள்", "மலர்",

	"காசு", "பிறப்பு",
}
var rhymeClass2Names = []string{
	"தேமா", "புளிமா", "கருவிளம்", "கூவிளம்",
}
var rhymeClass3Names = []string{
	"தேமாங்காய்", "புளிமாங்காய்", "கருவிளங்காய்", "கூவிளங்காய்",
	"தேமாங்கனி", "புளிமாங்கனி", "கருவிளங்கனி", "கூவிளங்கனி",
}
var rhymeClass4Names = []string{
	"தேமாந்தண்பூ", "தேமாந்தண்ணிழல்", "தேமாநறும்பூ", "தேமாநறுநிழல்",
	"புளிமாந்தண்பூ", "புளிமாந்தண்ணிழல்", "புளிமாநறும்பூ", "புளிமாநறுநிழல்",
	"கூவிளந்தண்பூ", "கூவிளந்தண்ணிழல்", "கூவிளநறும்பூ", "கூவிளநறுநிழல்",
	"கருவிளந்தண்பூ", "கருவிளந்தண்ணிழல்", "கருவிளநறும்பூ", "கருவிளநறுநிழல்",
}

func GetSortedClass1RhymeMap() map[string]Ceer { return createRhymeMap(rhymeClass1Names) }
func GetSortedClass2RhymeMap() map[string]Ceer { return createRhymeMap(rhymeClass2Names) }
func GetSortedClass3RhymeMap() map[string]Ceer { return createRhymeMap(rhymeClass3Names) }
func GetSortedClass4RhymeMap() map[string]Ceer { return createRhymeMap(rhymeClass4Names) }
func GetRhythmBaseMap() map[string]Ceer {
	var strs []string
	strs = append(strs, rhymeClass1Names...)
	strs = append(strs, rhymeClass2Names...)
	strs = append(strs, rhymeClass3Names...)
	strs = append(strs, rhymeClass4Names...)

	return createRhymeMap(strs)
}
func createRhymeMap(names []string) map[string]Ceer {
	dict := make(map[string]Ceer)
	for _, str := range names {
		s := script.MustLetterSeqFrom(str)
		captures := ToAcaaiCaptures(s, true) // Note: Reduction == true, is fine here!
		dict[CreateRhymeKey(captures)] = Ceer{UStr: str, Captures: captures}
	}
	return dict
}

// Rhythm pattern string (for rhythm matching)
func CreateRhymeKey(captures []AcaaiFrag) string {
	var sb strings.Builder
	for _, capture := range captures {
		sb.WriteString(string(capture.acaaiId))
	}
	return sb.String()
}

var acaaiNames = []string{"நேர்", "நிரை", "நேர்பு", "நிரைபு"}

// அசை fragment capture
type AcaaiFrag struct {
	ufrag   string
	acaaiId uint8 // 0: நேர், 1: நிரை, 2: நேர்பு, 3: நிரைபு
}

// Gets slice of indices into [நேர், நிரை], corresponding to the given word
// Reduce to "நேர்பு" | "நிரைபு" form, as applicable
func ToAcaaiCaptures(s script.LetterSeq, reduced bool) []AcaaiFrag {
	syls := script.Syllables(s) // Each syllable is simply a நேர், which may be upto 2 letters long.
	pending := false
	var captures []AcaaiFrag
	for i, syl := range syls {
		if pending {
			// Form a நிரை
			captures = append(captures, AcaaiFrag{ufrag: syls[i-1].String() + syls[i].String(), acaaiId: 1})
			pending = false
			continue
		}
		// See if நிரை could be formed...
		if syl.Len() == 1 && syl.First().IsShortVocal() {
			pending = true
			continue
		}
		// Form a நேர்
		captures = append(captures, AcaaiFrag{ufrag: syls[i].String(), acaaiId: 0})
	}
	if pending { // Unconsumed pending == நேர்
		captures = append(captures, AcaaiFrag{ufrag: syls[len(syls)-1].String(), acaaiId: 0})
	}

	// Optionally, attempt reducing to single "நேர்பு" | "நிரைபு" form.
	if reduced && len(captures) == 2 && captures[1].acaaiId == 0 {
		syl := syls[len(syls)-1]
		if syl.Len() == 1 {
			// The last letter better be CV letter! [Unless syllabification is broken!]
			_, v := syl.Nth(0).MustSplitCV()
			if v.Is(உ) {
				var codeIdx uint8 = 2
				if captures[0].acaaiId == 1 {
					codeIdx = 3
				}
				return []AcaaiFrag{{ufrag: s.String(), acaaiId: codeIdx}}
			}
		}
	}

	return captures
}

// Vowel Letter உ
var உ = script.MustLetterFrom("உ")

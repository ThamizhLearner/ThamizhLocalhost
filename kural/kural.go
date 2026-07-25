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

// Rhythm of beats of [நேர், நிரை, நேர்பு, நிரைபு]
type rhythm struct {
	UStr     string
	Captures []captured
}

// Gets "நிரை, நேர், நேர்" captured beat sequence.
func GetRhythmBeats(captures []captured, withCatpure bool) string {
	var strs []string
	for _, capture := range captures {
		str := captureCodes[capture.codeIdx]
		if withCatpure {
			str += fmt.Sprintf(" {%s}", capture.ufrag)
		}
		strs = append(strs, str)
	}
	return strings.Join(strs, " | ")
}

func RhythmBreakup(captures []captured) string {
	var strs []string
	for _, capture := range captures {
		strs = append(strs, capture.ufrag)
	}
	return strings.Join(strs, "/")
}

func GetRhythmBaseMap() map[string]rhythm {
	strs := []string{
		"நாள்", "மலர்",

		"காசு", "பிறப்பு",

		"தேமா", "புளிமா", "கருவிளம்", "கூவிளம்",

		"தேமாங்காய்", "புளிமாங்காய்", "கருவிளங்காய்", "கூவிளங்காய்",
		"தேமாங்கனி", "புளிமாங்கனி", "கருவிளங்கனி", "கூவிளங்கனி",

		"தேமாந்தண்பூ", "தேமாந்தண்ணிழல்", "தேமாநறும்பூ", "தேமாநறுநிழல்",
		"புளிமாந்தண்பூ", "புளிமாந்தண்ணிழல்", "புளிமாநறும்பூ", "புளிமாநறுநிழல்",
		"கூவிளந்தண்பூ", "கூவிளந்தண்ணிழல்", "கூவிளநறும்பூ", "கூவிளநறுநிழல்",
		"கருவிளந்தண்பூ", "கருவிளந்தண்ணிழல்", "கருவிளநறும்பூ", "கருவிளநறுநிழல்",
	}

	dict := make(map[string]rhythm)
	for _, str := range strs {
		s := script.MustLetterSeqFrom(str)
		captures := CaptureRhythm(s, true) // Note: Reduction == true, is fine here!
		dict[CreateKey(captures)] = rhythm{UStr: str, Captures: captures}
	}
	return dict
}

func CreateKey(captures []captured) string {
	var sb strings.Builder
	for _, capture := range captures {
		sb.WriteString(string(capture.codeIdx))
	}
	return sb.String()
}

var captureCodes = []string{"நேர்", "நிரை", "நேர்பு", "நிரைபு"}

// சீர் decomposition capture
type captured struct {
	ufrag   string
	codeIdx uint8
}

// Gets slice of indices into [நேர், நிரை], corresponding to the given word
// Reduce to "நேர்பு" | "நிரைபு" form, as applicable
func CaptureRhythm(s script.LetterSeq, reduced bool) []captured {
	syls := script.Syllables(s) // Each syllable is simply a நேர், which may be upto 2 letters long.
	pending := false
	var captures []captured
	for i, syl := range syls {
		if pending {
			// Form a நிரை
			captures = append(captures, captured{ufrag: syls[i-1].String() + syls[i].String(), codeIdx: 1})
			pending = false
			continue
		}
		// See if நிரை could be formed...
		if syl.Len() == 1 && syl.First().IsShortVocal() {
			pending = true
			continue
		}
		// Form a நேர்
		captures = append(captures, captured{ufrag: syls[i].String(), codeIdx: 0})
	}
	if pending { // Unconsumed pending == நேர்
		captures = append(captures, captured{ufrag: syls[len(syls)-1].String(), codeIdx: 0})
	}

	// Optionally, attempt reducing to single "நேர்பு" | "நிரைபு" form.
	if reduced && len(captures) == 2 && captures[1].codeIdx == 0 {
		syl := syls[len(syls)-1]
		if syl.Len() == 1 {
			// The last letter better be CV letter! [Unless syllabification is broken!]
			_, v := syl.Nth(0).MustSplitCV()
			if v.Is(உ) {
				var codeIdx uint8 = 2
				if captures[0].codeIdx == 1 {
					codeIdx = 3
				}
				return []captured{{ufrag: s.String(), codeIdx: codeIdx}}
			}
		}
	}

	return captures
}

// Vowel Letter உ
var உ = script.MustLetterFrom("உ")

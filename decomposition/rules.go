package decomposition

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	script "github.com/ThamizhLearner/Thamizh"
)

type SuffixTrimmer interface { // Suffix trimming rule [Technically, how it works is not important in itself]
	GetSuffix() script.String
	// Trim the suffix and return the remnant/s
	Trim(str script.String) []script.String // Note: More than one pathways may come into existence!!!
}

var Trimmers []SuffixTrimmer = LoadDecompositionRules("decomposition/rules.txt")

// "த்த்"
var str_தத = script.MustDecode("த்த்")

// "ம்"
var str_ம = script.MustDecode("ம்")

var str_empty = script.MustDecode("")

func createSuffixTrimmer(s string) SuffixTrimmer {
	subs := strings.Split(s, " ") // Look out for first Space separator

	if strings.Contains(subs[0], "|") { // Compound-form code
		subs := strings.Split(subs[0], "|")
		var substRules []SubstRule
		for _, s := range subs[1:] {
			idx := strings.Index(s, ":")
			if idx == -1 { // got match-trim 'n replace pair
				sr := SubstRule{matchTrimStr: script.MustDecode(strings.TrimSpace(s)), subsStr: str_empty}
				substRules = append(substRules, sr)
			} else {
				sr := SubstRule{
					matchTrimStr: script.MustDecode(strings.TrimSpace(s[:idx])),
					subsStr:      script.MustDecode(strings.TrimSpace(s[idx+1:])),
				}
				substRules = append(substRules, sr)
			}
		}
		return SuffixTrimRule{name: script.MustDecode(subs[0]), substRules: substRules}
	}

	sfx := script.MustDecode(subs[0]) // (Eg. இல், அம்)
	// For V-Suffix, add trim rules for ய் and வ் forms
	if sfx.FirstLetter().IsV() {
		substRules := []SubstRule{
			// (Eg. கோயில் = கோ + இல்)
			{matchTrimStr: script.MustDecode("ய்").Appended(sfx), subsStr: str_empty},
			// (Eg. கோவில் = கோ + இல்)
			{matchTrimStr: script.MustDecode("வ்").Appended(sfx), subsStr: str_empty},
			// ம் restoration form த்த் (Eg. மாற்றத்தை = மாற்றம் + ஐ)
			{matchTrimStr: str_தத.Appended(sfx), subsStr: str_ம},
			// (Eg. மாற்றம் = மாற் + அம்) | Applies to Strong Consonants (க், ச், ட், த், ப், ற்.)
			{matchTrimStr: script.MustDecode("க்க்").Appended(sfx), subsStr: script.MustDecode("க்")},
			{matchTrimStr: script.MustDecode("ச்ச்").Appended(sfx), subsStr: script.MustDecode("ச்")},
			{matchTrimStr: script.MustDecode("ட்ட்").Appended(sfx), subsStr: script.MustDecode("ட்")},
			{matchTrimStr: script.MustDecode("த்த்").Appended(sfx), subsStr: script.MustDecode("த்")},
			{matchTrimStr: script.MustDecode("ப்ப்").Appended(sfx), subsStr: script.MustDecode("ப்")},
			{matchTrimStr: script.MustDecode("ற்ற்").Appended(sfx), subsStr: script.MustDecode("ற்")},
		}

		return SuffixTrimRule{name: sfx, substRules: substRules}
	}

	// For CV-Suffix, check if we have Strong consonant (Eg. கள், தல்)
	if sfx.FirstLetter().IsCV() && sfx.FirstLetter().IsStrongVocal() {
		c, _ := sfx.FirstLetter().SplitCV()
		substRules := []SubstRule{
			{matchTrimStr: script.MustDecode(c.String()).Appended(sfx), subsStr: str_empty},
			{matchTrimStr: sfx, subsStr: str_empty},
			// {matchTrimStr: script.MustDecode(c.String()).Appended(sfx), subsStr: script.MustDecode("உ")},
		}
		return SuffixTrimRule{name: sfx, substRules: substRules}
	}

	return SuffixTrimRule{name: sfx, substRules: []SubstRule{
		{matchTrimStr: sfx, subsStr: str_empty},
	}}
}

func Assert(assertion bool, msg string) {
	if !assertion {
		panic(msg)
	}
}

// Suffix decomposition rule
type SuffixTrimRule struct {
	name       script.String
	substRules []SubstRule
}

func (rule SuffixTrimRule) GetSuffix() script.String { return rule.name }

// Trim out the suffix and return all the possible trim remnants
func (rule SuffixTrimRule) Trim(str script.String) []script.String {
	var strs []script.String
	// Try (raw) matching one options
	for _, sfxRule := range rule.substRules { // No letter-split trim performed here
		res, ok := str.TailTrimmedRaw(sfxRule.matchTrimStr)
		if !ok {
			continue
		}
		// Substitute if there is any substitution specified
		if sfxRule.subsStr.Len() != 0 {
			res = res.AppendedRaw(sfxRule.subsStr)
		}
		strs = append(strs, res)

		break // One rule match only allowed
	}

	// Additional handling for V-Suffixes (requires letter-split trimming) (Eg. இல், அம்)
	sfx := rule.name
	if sfx.FirstLetter().IsV() { // Letter-split trim performed here
		rem, ok := str.TailTrimmed(sfx)
		if !ok {
			return nil
		}
		strs = append(strs, rem)
	}
	return strs
}

func (rule SuffixTrimRule) String() string {
	if len(rule.substRules) == 0 {
		return rule.name.String()
	}
	sb := strings.Builder{}
	sb.WriteString(rule.name.String())
	for _, r := range rule.substRules {
		fmt.Fprint(&sb, " | ", r.matchTrimStr.String())
		if r.subsStr.Len() != 0 {
			fmt.Fprint(&sb, " : ", r.subsStr.String())
		}
	}
	return sb.String()
}

type SubstRule struct { // match 'n substitute
	matchTrimStr script.String
	subsStr      script.String
}

func LoadDecompositionRules(fname string) []SuffixTrimmer {
	// Load and parse the suffix definition file
	file, err := os.Open(fname)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var trimmers []SuffixTrimmer
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := strings.TrimSpace(scanner.Text())
		if len(s) == 0 || s[0] == ';' { // Skip "empty" lines
			continue
		}
		trimmers = append(trimmers, createSuffixTrimmer(s))
	}
	if err = scanner.Err(); err != nil {
		panic(err)
	}

	// dump trimmers
	// for idx, t := range trimmers {
	// 	fmt.Println(idx, t)
	// }

	return trimmers
}

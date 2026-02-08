package decomposition

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	script "github.com/ThamizhLearner/Thamizh"
)

type SuffixTrimmer interface { // Suffix trimming rule [Technically, how it works is not important in itself]
	GetSuffix() string
	// Trim the suffix and return the remnant
	Trim(str script.String) []script.String // Note: More than one pathways may come into existence!!!
}

var Trimmers []SuffixTrimmer

func createSuffixTrimmer(s string) SuffixTrimmer {
	// கள்|க்கள்|ங்கள்:ம்|ற்கள்:ல்|ட்கள்:ள்|கள் [remnant=noun=singular] [suffix=noun=plural] (பன்மை விகுதி)
	// த்தல் [remnant=verb=?] [suffix=noun=verb] (? விகுதி)
	// தல் [remnant=verb=base] [suffix=noun=verb] (தொழிற்பெயர் விகுதி)
	// அல் [remnant=verb=stem] [suffix=noun=verb] (தொழிற்பெயர் விகுதி)

	panic("")
}

type SuffixSubstitutionDef struct { // Match 'n trim 'n replace
	name              string
	substitutionRules []SubstitutionRule
}

func (sfx SuffixSubstitutionDef) GetSuffix() string { return sfx.name }

// Trim the suffix and return the remnant
func (sfx SuffixSubstitutionDef) Trim(str script.String) []script.String {
	for _, sfxRule := range sfx.substitutionRules {
		res, ok := str.TailTrimmedRaw(sfxRule.matchTrimStr)
		if !ok {
			continue
		}
		if sfxRule.subsStr.Len() != 0 {
			res = res.AppendedRaw(sfxRule.subsStr)
		}
		return []script.String{res}
	}
	return nil
}

type SubstitutionRule struct { // match 'n substitute
	matchTrimStr script.String
	subsStr      script.String
}

func NewSuffixSubstitutionRule(matchStr string, subssubstitutionStr string) SubstitutionRule {
	return SubstitutionRule{
		matchTrimStr: script.MustDecode(matchStr),
		subsStr:      script.MustDecode(subssubstitutionStr),
	}
}

func NewSingularizationDef() SuffixSubstitutionDef {
	// கள்|க்கள்|ங்கள்:ம்|ற்கள்:ல்|ட்கள்:ள்|கள் [remnant=noun=singular] [suffix=noun=plural] (பன்மை விகுதி)
	substitutions := []SubstitutionRule{
		NewSuffixSubstitutionRule("க்கள்", ""),
		NewSuffixSubstitutionRule("ங்கள்", "ம்"),
		NewSuffixSubstitutionRule("ற்கள்", "ல்"),
		NewSuffixSubstitutionRule("ட்கள்", "ள்"),
		NewSuffixSubstitutionRule("கள்", ""),
	}
	return SuffixSubstitutionDef{name: "கள்", substitutionRules: substitutions}
}

// கள்|க்கள்|ங்கள்:ம்|ற்கள்:ல்|ட்கள்:ள்|கள் [remnant=noun=singular] [suffix=noun=plural] (பன்மை விகுதி)
func createOverriddenTrimmerQuick() SuffixTrimmer {
	panic("")
}

func loadDecompositionRules(fname string) []SuffixTrimmer {
	// Load and parse the suffix definition file
	file, err := os.Open(fname)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := strings.TrimSpace(scanner.Text())
		if len(s) == 0 || s[0] == ';' {
			continue
		}
		fmt.Println(s)
		createSuffixTrimmer(s)
	}
	if err = scanner.Err(); err != nil {
		fmt.Println(err)
	}

	panic("")
}

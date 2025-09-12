package pchecker

import (
	"pchecker/ds"
)

var (
	pd *ProfanityDetector
)

type ProfanityReplacementFunc func(match []rune) []rune

func getSafeTrie(m map[string]bool) *ds.SafeTrie[rune] {
	result := ds.NewSafeTrie[rune](len(m))
	for word := range m {
		result.Insert([]rune(word))
	}
	return result
}

// Censor takes in a string (word or sentence) and tries to censor all profanities found.
//
// Uses the default ProfanityDetector
func Censor(s string, f ProfanityReplacementFunc) string {
	if pd == nil {
		pd = NewDefaultProfanityDetector()
	}
	return pd.censor(s, f)
}

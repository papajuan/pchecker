package pchecker

var (
	pd *ProfanityDetector
)

type ReplacementFunc func(match []rune) string

func getSafeTrie(m map[string]bool) *SafeTrie[rune] {
	result := NewSafeTrie[rune](len(m))
	for word := range m {
		result.Insert([]rune(word))
	}
	return result
}

// Censor takes in a string (word or sentence) and tries to censor all profanities found.
//
// Uses the default ProfanityDetector
func Censor(s string, f ReplacementFunc) string {
	if pd == nil {
		pd = NewDefaultProfanityDetector()
	}
	return pd.censor(s, f)
}

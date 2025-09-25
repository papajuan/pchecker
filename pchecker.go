package pchecker

type ReplacementFunc func(match []rune) string

func getSafeTrie(m map[string]bool) *SafeTrie[rune] {
	result := NewSafeTrie[rune](len(m))
	for word := range m {
		result.Insert([]rune(word))
	}
	return result
}

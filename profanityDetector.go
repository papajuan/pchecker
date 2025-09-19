package pchecker

import (
	"unicode"
	"unicode/utf8"
)

/**
 * @author  papajuan
 * @date    9/12/2025
 **/

// ProfanityDetector contains the dictionaries as well as the configuration
// for determining how profanity detection is handled
type ProfanityDetector struct {
	profanities           *SafeTrie[rune]
	falsePositives        *SafeTrie[rune]
	falseNegatives        *SafeTrie[rune]
	characterReplacements map[rune]rune
}

func NewProfanityDetector() *ProfanityDetector {
	return &ProfanityDetector{}
}

// NewDefaultProfanityDetector creates a new ProfanityDetector with the default settings
func NewDefaultProfanityDetector() *ProfanityDetector {
	return NewProfanityDetector().
		WithDefaultFalsePositives().
		WithDefaultFalseNegatives().
		WithDefaultProfanities().
		WithDefaultCharacterReplacements()
}

func (pd *ProfanityDetector) WithProfanities(profanities map[string]bool) *ProfanityDetector {
	pd.profanities = getSafeTrie(profanities).WithComparator(unicode.ToLower)
	return pd
}

func (pd *ProfanityDetector) WithDefaultProfanities() *ProfanityDetector {
	pd.profanities = getSafeTrie(DefaultProfanities).WithComparator(unicode.ToLower)
	return pd
}

func (pd *ProfanityDetector) WithFalsePositives(falsePositives map[string]bool) *ProfanityDetector {
	pd.falsePositives = getSafeTrie(falsePositives).WithComparator(unicode.ToLower)
	return pd
}

func (pd *ProfanityDetector) WithDefaultFalsePositives() *ProfanityDetector {
	pd.falsePositives = getSafeTrie(DefaultFalsePositives).WithComparator(unicode.ToLower)
	return pd
}

func (pd *ProfanityDetector) WithFalseNegatives(falseNegatives map[string]bool) *ProfanityDetector {
	pd.falseNegatives = getSafeTrie(falseNegatives).WithComparator(unicode.ToLower)
	return pd
}

func (pd *ProfanityDetector) WithDefaultFalseNegatives() *ProfanityDetector {
	pd.falseNegatives = getSafeTrie(DefaultFalseNegatives).WithComparator(unicode.ToLower)
	return pd
}

func (pd *ProfanityDetector) WithCharacterReplacements(replacements map[rune]rune) *ProfanityDetector {
	pd.characterReplacements = replacements
	return pd
}

func (pd *ProfanityDetector) WithDefaultCharacterReplacements() *ProfanityDetector {
	pd.characterReplacements = DefaultCharacterReplacements
	return pd
}

func (pd *ProfanityDetector) censor(input string, f ReplacementFunc) string {
	tb := newTokenBuffer(utf8.RuneCountInString(input), f)
	defer tb.close()
	var curr *node[rune] = nil // rolling pointer in the profanities trie
	for _, r := range input {
		if unicode.IsPunct(r) || unicode.IsSpace(r) {
			tb.flush(pd.falsePositives)
			tb.result.WriteRune(r)
			curr = nil
			continue
		}
		normRune := pd.getCharReplacement(r)
		tb.buff = append(tb.buff, normRune)
		// Advance the trie
		if curr == nil {
			curr = pd.profanities.root
		}
		if next, ok := curr.children[normRune]; ok {
			curr = next
			if curr.isEnd {
				tb.badToken = true
			}
		} else {
			// Dead end: restart from root using this rune
			if next, ok = pd.profanities.root.children[normRune]; ok {
				curr = next
				if curr.isEnd {
					tb.badToken = true
				}
			} else {
				curr = nil
			}
		}
	}
	// flush the last token
	tb.flush(pd.falsePositives)
	return tb.String()
}

func (pd *ProfanityDetector) getCharReplacement(original rune) rune {
	if replacement, found := pd.characterReplacements[unicode.ToLower(original)]; found {
		return replacement
	}
	return original
}

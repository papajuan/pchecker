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

func (pd *ProfanityDetector) Profanities() *SafeTrie[rune] {
	return pd.profanities
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

func (pd *ProfanityDetector) Censor(input string, f ReplacementFunc) string {
	tb := newTokenBuffer(utf8.RuneCountInString(input), f)
	defer tb.close()
	// start from root so we never lose the trie
	curr := pd.profanities.root
	for _, r := range input {
		if (r != '@' && r != '_' && unicode.IsPunct(r)) || unicode.IsSpace(r) {
			tb.flush(pd.falsePositives, pd.falseNegatives)
			tb.result.WriteRune(r)
			curr = pd.profanities.root
			continue
		}
		normRune := pd.getCharReplacement(r)
		tb.buff = append(tb.buff, r)
		// try to advance current branch
		if next, ok := curr.children[unicode.ToLower(normRune)]; ok {
			curr = next
			if curr.isEnd {
				tb.badToken = true
			}
			continue
		}
		// dead end: try to match this rune from root
		if next, ok := pd.profanities.root.children[unicode.ToLower(normRune)]; ok {
			curr = next
			if curr.isEnd {
				tb.badToken = true
			}
			// additionally: maybe there's a profanity that starts earlier in tb.buff
			// (e.g. "bigblack" blocked "breast" inside "bigbreasts")
			if !tb.badToken && pd.containsProfanityInBuffer(tb.buff) {
				tb.badToken = true
			}
			continue
		}
		// couldn't match current rune even from root — reset to root
		curr = pd.profanities.root
		// but still: check the whole buffer for any substring match
		if !tb.badToken && pd.containsProfanityInBuffer(tb.buff) {
			tb.badToken = true
		}
	}
	// flush the last token
	tb.flush(pd.falsePositives, pd.falseNegatives)
	return tb.String()
}

// containsProfanityInBuffer checks if any substring of buf matches a word in profanities trie.
// It uses the same normalization pd.getCharReplacement and unicode.ToLower as lookups.
func (pd *ProfanityDetector) containsProfanityInBuffer(buf []rune) bool {
	n := len(buf)
	for start := 0; start < n; start++ {
		curr := pd.profanities.root
		for i := start; i < n; i++ {
			norm := unicode.ToLower(pd.getCharReplacement(buf[i]))
			if child, ok := curr.children[norm]; !ok {
				break
			} else {
				curr = child
				if curr.isEnd {
					return true
				}
			}
		}
	}
	return false
}

func (pd *ProfanityDetector) getCharReplacement(original rune) rune {
	if replacement, found := pd.characterReplacements[unicode.ToLower(original)]; found {
		return replacement
	}
	return original
}

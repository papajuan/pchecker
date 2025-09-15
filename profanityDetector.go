package pchecker

import (
	"sync"
	"unicode"
	"unicode/utf8"

	"github.com/papajuan/pchecker/ds"
)

/**
 * @author  papajuan
 * @date    9/12/2025
 **/

// ProfanityDetector contains the dictionaries as well as the configuration
// for determining how profanity detection is handled
type ProfanityDetector struct {
	profanities           *ds.SafeTrie[rune]
	falsePositives        *ds.SafeTrie[rune]
	falseNegatives        *ds.SafeTrie[rune]
	characterReplacements map[rune]rune
	runeBufPool           sync.Pool
}

func NewProfanityDetector() *ProfanityDetector {
	return &ProfanityDetector{
		runeBufPool: sync.Pool{New: func() any { return make([]rune, 0, 32) }},
	}
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
	pd.profanities = getSafeTrie(profanities)
	return pd
}

func (pd *ProfanityDetector) WithDefaultProfanities() *ProfanityDetector {
	pd.profanities = getSafeTrie(DefaultProfanities)
	return pd
}

func (pd *ProfanityDetector) WithFalsePositives(falsePositives map[string]bool) *ProfanityDetector {
	pd.falsePositives = getSafeTrie(falsePositives)
	return pd
}

func (pd *ProfanityDetector) WithDefaultFalsePositives() *ProfanityDetector {
	pd.falsePositives = getSafeTrie(DefaultFalsePositives)
	return pd
}

func (pd *ProfanityDetector) WithFalseNegatives(falseNegatives map[string]bool) *ProfanityDetector {
	pd.falseNegatives = getSafeTrie(falseNegatives)
	return pd
}

func (pd *ProfanityDetector) WithDefaultFalseNegatives() *ProfanityDetector {
	pd.falseNegatives = getSafeTrie(DefaultFalseNegatives)
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

func (pd *ProfanityDetector) censor(input string, replacementFunc ProfanityReplacementFunc) string {
	var (
		possibleBadWord               = pd.runeBufPool.Get().([]rune)[:0]
		result                        = make([]rune, 0, utf8.RuneCountInString(input))
		isBadWord, hasBadPrefix, keep bool
	)
	defer pd.runeBufPool.Put(possibleBadWord[:0])
	inputLength := len(input)
	for i, originalChar := range input {
		isLast := i == inputLength-1
		if unicode.IsPunct(originalChar) || unicode.IsSpace(originalChar) {
			if possibleBadWord != nil && len(possibleBadWord) > 0 {
				if isBadWord && !pd.falsePositives.IsPrefixInTrie(possibleBadWord) {
					result = append(result, replacementFunc(possibleBadWord)...)
				} else {
					result = append(result, possibleBadWord...)
				}
			}
			result = append(result, originalChar)
			possibleBadWord = possibleBadWord[:0]
			isBadWord = false
			keep = false
		} else if isLast {
			if hasBadPrefix && !isBadWord {
				possibleBadWord = append(possibleBadWord, originalChar)
				_, isBadWord = pd.profanities.StartsWith(possibleBadWord)
				if isBadWord && !pd.falsePositives.IsPrefixInTrie(possibleBadWord) {
					result = append(result, replacementFunc(possibleBadWord)...)
				} else {
					result = append(result, possibleBadWord...)
				}
			} else if isBadWord {
				possibleBadWord = append(possibleBadWord, originalChar)
				if pd.falsePositives.IsPrefixInTrie(possibleBadWord) {
					result = append(result, possibleBadWord...)
				} else {
					result = append(result, replacementFunc(possibleBadWord)...)
				}
				isBadWord = false
			} else {
				if possibleBadWord != nil && len(possibleBadWord) > 0 {
					result = append(result, possibleBadWord...)
				} else {
					result = append(result, originalChar)
				}
			}
		} else {
			if isBadWord {
				possibleBadWord = append(possibleBadWord, originalChar)
			} else {
				possibleBadWord = append(possibleBadWord, pd.getCharReplacement(originalChar))
				hasBadPrefix, isBadWord = pd.profanities.StartsWith(possibleBadWord)
				if !hasBadPrefix {
					result = append(result, possibleBadWord...)
					possibleBadWord = possibleBadWord[:0]
					keep = true
				} else if keep {
					result = append(result, originalChar)
					possibleBadWord = possibleBadWord[:0]
					isBadWord = false
				}
			}
		}
	}
	return string(result)
}

func (pd *ProfanityDetector) getCharReplacement(original rune) rune {
	if replacement, found := pd.characterReplacements[unicode.ToLower(original)]; found {
		return replacement
	}
	return original
}

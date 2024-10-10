package pchecker

import (
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"strings"
	"unicode"
)

var (
	defaultProfanityDetector *ProfanityDetector
)

// ProfanityDetector contains the dictionaries as well as the configuration
// for determining how profanity detection is handled
type ProfanityDetector struct {
	profanities           *SafeSet[string]
	falsePositives        *SafeSet[string]
	falseNegatives        *SafeSet[string]
	characterReplacements map[rune]rune
}

// NewProfanityDetector creates a new ProfanityDetector
func NewProfanityDetector() *ProfanityDetector {
	return &ProfanityDetector{
		profanities:           NewSafeSet(DefaultProfanities),
		falsePositives:        NewSafeSet(DefaultFalsePositives),
		falseNegatives:        NewSafeSet(DefaultFalseNegatives),
		characterReplacements: DefaultCharacterReplacements,
	}
}

func (g *ProfanityDetector) WithCustomDicts(profanities, falsePositives, falseNegatives map[string]bool) *ProfanityDetector {
	g.profanities.PutAll(profanities)
	g.falsePositives.PutAll(falsePositives)
	g.falseNegatives.PutAll(falseNegatives)
	return g
}

func (g *ProfanityDetector) indexToRune(s string, index int) int {
	count := 0
	for i := range s {
		if i == index {
			break
		}
		if i < index {
			count++
		}
	}
	return count
}

func (g *ProfanityDetector) Censor(s string) string {
	censored := []rune(s)
	var originalIndexes []int
	s, originalIndexes = g.sanitize(s, true)
	runeWordLength := 0

	g.checkProfanity(&s, &originalIndexes, &censored, g.falseNegatives, &runeWordLength)
	g.removeFalsePositives(&s, &originalIndexes, &runeWordLength)
	g.checkProfanity(&s, &originalIndexes, &censored, g.profanities, &runeWordLength)

	return string(censored)
}

func (g *ProfanityDetector) checkProfanity(s *string, originalIndexes *[]int, censored *[]rune,
	wordMap *SafeSet[string], runeWordLength *int) {
	wordMap.Range(func(word string) {
		currentIndex := 0
		*runeWordLength = len([]rune(word))
		for currentIndex != -1 {
			if foundIndex := strings.Index((*s)[currentIndex:], word); foundIndex != -1 {
				for i := 0; i < *runeWordLength; i++ {
					runeIndex := g.indexToRune(*s, currentIndex+foundIndex) + i
					if runeIndex < len(*originalIndexes) {
						(*censored)[(*originalIndexes)[runeIndex]] = '*'
					}
				}
				currentIndex += foundIndex + len([]byte(word))
			} else {
				break
			}
		}
	})
}

func (g *ProfanityDetector) removeFalsePositives(s *string, originalIndexes *[]int, runeWordLength *int) {
	g.falsePositives.Range(func(word string) {
		currentIndex := 0
		*runeWordLength = len([]rune(word))
		for currentIndex != -1 {
			if foundIndex := strings.Index((*s)[currentIndex:], word); foundIndex != -1 {
				foundRuneIndex := g.indexToRune(*s, foundIndex)
				*originalIndexes = append((*originalIndexes)[:foundRuneIndex], (*originalIndexes)[foundRuneIndex+*runeWordLength:]...)
				currentIndex += foundIndex + len([]byte(word))
			} else {
				break
			}
		}
		*s = strings.Replace(*s, word, "", -1)
	})
}

func (g *ProfanityDetector) sanitize(s string, rememberOriginalIndexes bool) (string, []int) {
	s = strings.ToLower(s)
	if !rememberOriginalIndexes {
		s = strings.ReplaceAll(s, "()", "o")
	}
	sb := strings.Builder{}
	for _, char := range s {
		if replacement, found := g.characterReplacements[char]; found {
			sb.WriteRune(replacement)
			continue
		}
		sb.WriteRune(char)
	}
	s = sb.String()
	s = g.removeAccents(s)
	var originalIndexes []int
	if rememberOriginalIndexes {
		for i, c := range []rune(s) {
			if c != ' ' {
				originalIndexes = append(originalIndexes, i)
			}
		}
	}
	s = strings.Replace(s, " ", "", -1)
	return s, originalIndexes
}

// removeAccents strips all accents from characters.
// Only called if ProfanityDetector.removeAccents is set to true
func (_ *ProfanityDetector) removeAccents(s string) string {
	removeAccentsTransformer := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	for _, character := range s {
		// If there's a character outside the range of supported runes, there might be some accented words
		if character < ' ' || character > '~' {
			s, _, _ = transform.String(removeAccentsTransformer, s)
			break
		}
	}
	return s
}

// Censor takes in a string (word or sentence) and tries to censor all profanities found.
//
// Uses the default ProfanityDetector
func Censor(s string) string {
	if defaultProfanityDetector == nil {
		defaultProfanityDetector = NewProfanityDetector()
	}
	return defaultProfanityDetector.Censor(s)
}

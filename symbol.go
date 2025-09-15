package pchecker

/**
 * @author  papajuan
 * @date    10/14/2024
 **/

type symbol struct {
	bad   bool
	index int
	char  rune
}

func newSymbol(index int, char rune) *symbol {
	return &symbol{
		index: index,
		char:  char,
	}
}

package pchecker

import (
	"strings"
	"sync"
)

/**
 * @author  papajuan
 * @date    9/17/2025
 **/

type tokenBuffer struct {
	buffPool sync.Pool
	buff     []rune
	badToken bool
	result   strings.Builder
	f        ReplacementFunc
}

func newTokenBuffer(resultLen int, f ReplacementFunc) *tokenBuffer {
	result := &tokenBuffer{
		f: f,
	}
	result.result.Grow(resultLen)
	result.buffPool = sync.Pool{New: func() any { return make([]rune, 0, 16) }}
	result.buff = result.buffPool.Get().([]rune)[:0]
	return result
}

func (tb *tokenBuffer) flush(falsePositives *SafeTrie[rune]) {
	if len(tb.buff) > 0 {
		if tb.badToken && !falsePositives.IsPrefixInTrie(tb.buff) {
			tb.result.WriteString(tb.f(tb.buff))
		} else {
			for _, r := range tb.buff {
				tb.result.WriteRune(r)
			}
		}
		tb.buff = tb.buff[:0]
		tb.badToken = false
	}
}

func (tb *tokenBuffer) String() string {
	return tb.result.String()
}

func (tb *tokenBuffer) close() {
	tb.buffPool.Put(tb.buff[:0])
}

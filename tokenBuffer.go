package pchecker

import (
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
	result   []rune
	f        ReplacementFunc
}

func newTokenBuffer(resultLen int, f ReplacementFunc) *tokenBuffer {
	result := &tokenBuffer{
		result: make([]rune, 0, resultLen),
		f:      f,
	}
	result.buffPool = sync.Pool{New: func() any { return make([]rune, 0, 32) }}
	result.buff = result.buffPool.Get().([]rune)[:0]
	return result
}

func (tb *tokenBuffer) flush(falsePositives *SafeTrie[rune]) {
	if len(tb.buff) > 0 {
		if tb.badToken && !falsePositives.IsPrefixInTrie(tb.buff) {
			tb.result = append(tb.result, tb.f(tb.buff)...)
		} else {
			tb.result = append(tb.result, tb.buff...)
		}
		tb.buff = tb.buff[:0]
		tb.badToken = false
	}
}

func (tb *tokenBuffer) String() string {
	return string(tb.result)
}

func (tb *tokenBuffer) close() {
	tb.buffPool.Put(tb.buff[:0])
}

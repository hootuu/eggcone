package token

import (
	"sync/atomic"
)

type Allocator struct {
	size   uint64
	tokens []Token
	seq    uint64
}

func New(size uint64) *Allocator {
	a := &Allocator{
		size:   size,
		tokens: make([]Token, size),
		seq:    0,
	}
	a.init()
	return a
}

func (a *Allocator) init() {
	size := len(a.tokens)
	for i := 0; i < size; i++ {
		a.tokens[i] = newToken()
	}
}

func (a *Allocator) Alloc() Token {
	n := atomic.AddUint64(&a.seq, 1)
	idx := n % a.size
	return a.tokens[idx]
}

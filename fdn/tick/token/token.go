package token

import (
	"fmt"
	"github.com/hootuu/gelato/errors"
	"github.com/hootuu/gelato/sys"
	"sync/atomic"
)

type Token string

func (t Token) S() string {
	return string(t)
}

func (t Token) Verify() *errors.Error {
	if len(t) == 0 {
		return errors.Verify("require Token")
	}
	return nil
}

var gTokenLocalSeq uint64 = 0

func newToken() Token {
	n := atomic.AddUint64(&gTokenLocalSeq, 1)
	str := fmt.Sprintf("%s_%04d", sys.ServerID, n)
	return Token(str)
}

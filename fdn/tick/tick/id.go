package tick

import (
	"fmt"
	"github.com/hootuu/gelato/sys"
	"sync/atomic"
)

type ID string

const (
	NilID ID = ""
)

func (id ID) S() string {
	return string(id)
}

var gTickIDLocalSeq uint64 = 0

func newID() ID {
	n := atomic.AddUint64(&gTickIDLocalSeq, 1)
	str := fmt.Sprintf("%s_%04d", sys.ServerID, n)
	return ID(str)
}

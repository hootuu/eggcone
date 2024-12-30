package schedule

import (
	"errors"
	"github.com/hootuu/gelato/idx"
)

type ID string

const (
	NilID ID = ""
)

func (id ID) String() string {
	return string(id)
}

func (id ID) Verify() error {
	if id == NilID {
		return errors.New("require ID")
	}
	return nil
}

func NewID() ID {
	return ID(idx.New())
}

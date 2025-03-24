package eggmq

import (
	"github.com/hootuu/gelato/errors"
)

type Listener interface {
	Topic() string
	Handle(msg *Message) *errors.Error
}

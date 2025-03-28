package eggmq

import (
	"github.com/hootuu/gelato/errors"
)

type Listener func(msg *Message) *errors.Error

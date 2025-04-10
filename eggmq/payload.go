package eggmq

import "github.com/hootuu/gelato/errors"

type Payload interface {
	Of(str string) *errors.Error
	To() string
}

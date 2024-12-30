package schedule

import (
	"github.com/hootuu/gelato/errors"
)

type Cron string

func (c Cron) S() string {
	return string(c)
}

func (c Cron) Verify() *errors.Error {
	if len(c) == 0 {
		return errors.Verify("require Cron")
	}
	return nil
}

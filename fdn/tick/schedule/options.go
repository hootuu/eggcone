package schedule

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type FailAct int

const (
	END   FailAct = 1
	RETRY FailAct = 2
)

type Options struct {
	FailAct   FailAct `bson:"fail_act" json:"fail_act"`
	RetryTime int     `bson:"retry_time,omitempty" json:"retry_time,omitempty"`
}

func NewDefaultOptions() *Options {
	return &Options{
		FailAct:   RETRY,
		RetryTime: 3,
	}
}

func (opt *Options) Scan(value interface{}) error {
	if value == nil {
		opt = nil
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, opt)
	default:
		return errors.New("invalid type for Dict")
	}
}

func (opt Options) Value() (driver.Value, error) {
	return json.Marshal(opt)
}

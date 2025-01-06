package def

import (
	"database/sql/driver"
	"encoding/json"
	nerrors "errors"
	"github.com/hootuu/gelato/errors"
	"github.com/hootuu/gelato/idx"
)

type JobID string

func NewJobID() JobID {
	return JobID(idx.New())
}

type Job struct {
	ID      JobID    `bson:"_id,omitempty" json:"id,omitempty"`
	Topic   string   `bson:"topic" json:"topic"`
	Payload *Payload `bson:"payload" json:"payload"`
	Tag     []string `bson:"tag,omitempty" json:"tag,omitempty"`
}

func (j *Job) Verify() *errors.Error {
	if len(j.ID) == 0 {
		return errors.Verify("require job.ID")
	}
	if len(j.Topic) == 0 {
		return errors.Verify("require job.topic")
	}
	if j.Payload == nil {
		return errors.Verify("require job.payload")
	}
	return nil
}

func (j *Job) Scan(value interface{}) error {
	if value == nil {
		j = nil
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, j)
	default:
		return nerrors.New("invalid type for Dict")
	}
}

func (j Job) Value() (driver.Value, error) {
	return json.Marshal(j)
}

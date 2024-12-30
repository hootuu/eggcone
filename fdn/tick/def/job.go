package def

import (
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

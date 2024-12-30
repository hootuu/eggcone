package schedule

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

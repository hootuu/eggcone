package tick

import (
	"github.com/hootuu/eggcone/fdn/tick/def"
	"github.com/hootuu/eggcone/fdn/tick/schedule"
	"github.com/hootuu/gelato/errors"
)

type ScheduleIO struct {
	Title     string        `json:"title"`
	Code      string        `json:"code"`
	OutID     string        `json:"out_id"`
	Cron      schedule.Cron `json:"cron"`
	Topic     string        `json:"topic"`
	Payload   *def.Payload  `json:"payload"`
	Signature string        `json:"signature"`
}

func (req *ScheduleIO) Verify() *errors.Error {
	if err := req.Cron.Verify(); err != nil {
		return err
	}
	if len(req.Code) == 0 {
		return errors.Verify("require code")
	}
	if len(req.OutID) == 0 {
		return errors.Verify("require out_id")
	}
	if len(req.Topic) == 0 {
		return errors.Verify("require topic")
	}
	if len(req.Signature) == 0 {
		return errors.Verify("require signature")
	}
	if req.Payload == nil {
		return errors.Verify("require Payload")
	}
	return nil
}

// doRegisterSchedule plan schedule
func doRegisterSchedule(req *ScheduleIO) (schedule.ID, *errors.Error) {
	if err := req.Verify(); err != nil {
		return schedule.NilID, err
	}

	srcM, err := schedule.LoadByCodeAndOutID(req.Code, req.OutID)
	if err != nil {
		return schedule.NilID, err
	}
	update := false
	if srcM != nil {
		if srcM.Signature == req.Signature {
			return srcM.ID, nil
		}
		update = true
	}
	destM := &schedule.Schedule{
		Token:   gTokenAllocator.Alloc(),
		Title:   req.Title,
		Code:    req.Code,
		OutID:   req.OutID,
		Cron:    req.Cron,
		Options: schedule.NewDefaultOptions(),
		Job: &def.Job{
			ID:      def.NewJobID(),
			Topic:   req.Topic,
			Payload: req.Payload,
			Tag:     nil,
		},
		Signature: req.Signature,
		Available: true,
	}

	if update {
		err = schedule.Update(srcM, destM)
		if err != nil {
			return srcM.ID, err
		}
		return srcM.ID, nil
	}

	id, err := schedule.New(destM)
	if err != nil {
		return schedule.NilID, err
	}
	return id, nil
}

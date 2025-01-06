package tick

import (
	"github.com/hootuu/eggcone/fdn/tick/dbx"
	"github.com/hootuu/eggcone/fdn/tick/schedule"
	"github.com/hootuu/gelato/errors"
	"github.com/hootuu/gelato/sys"
	"time"
)

func NewTick() (ID, *errors.Error) {
	id := newID()
	m := &Tick{
		ID:               id,
		Server:           sys.ServerID,
		LstHeartbeatTime: time.Now(),
		Living:           true,
		Version:          0,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	dbErr := dbx.DB().Create(m).Error
	if dbErr != nil {
		return NilID, errors.System("db error", dbErr)
	}

	return id, nil
}

func DealRecord(scheduleID schedule.ID, result bool, ltn string, ctx any) error {
	m := &Record{
		ScheduleID: scheduleID,
		Result:     result,
		DealTime:   time.Now(),
		Listener:   ltn,
		Ctx:        ctx,
		UpdatedAt:  time.Now(),
	}

	dbErr := dbx.DB().Create(m).Error
	if dbErr != nil {
		return errors.System("db error", dbErr)
	}
	return nil
}

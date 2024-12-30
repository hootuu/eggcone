package tick

import (
	"context"
	"github.com/hootuu/eggcone/fdn/tick/dbx"
	"github.com/hootuu/eggcone/fdn/tick/schedule"
	"github.com/hootuu/eggcone/fdn/tick/tick"
	"github.com/hootuu/gelato/errors"
	"github.com/hootuu/gelato/sys"
	"gorm.io/gorm"
)

type Daemon struct {
}

func NewDaemon(db *gorm.DB) *Daemon {
	d := &Daemon{}
	dbInit(db)
	return d
}

func (d *Daemon) Code() string {
	return "EGG_TICK"
}

func (d *Daemon) Startup() *errors.Error {
	doTimerInit()
	doStart()
	return nil
}

func (d *Daemon) Shutdown(_ context.Context) *errors.Error {
	doStop()
	return nil
}

func RegisterSchedule(req *ScheduleIO) (schedule.ID, *errors.Error) {
	return doRegisterSchedule(req)
}

func RegisterListener(ltn Listener) {
	doRegisterListener(ltn)
}

func dbInit(db *gorm.DB) {
	dbx.Init(db)
	err := dbx.DB().AutoMigrate(&schedule.Schedule{}, &tick.Bind{}, &tick.Tick{}, &tick.Record{})
	if err != nil {
		sys.Exit(errors.System("auto migrate for tick failed"))
	}
}

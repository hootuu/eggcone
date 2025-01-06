package tick

import (
	"fmt"
	"github.com/hootuu/gelato/configure"
	"github.com/hootuu/gelato/logger"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

var gCrontab = newCrontab()

const (
	CronExpHeartbeat    = "*/16 * * * * ?"
	CronExpSyncSchedule = "*/18 * * * * ?"
	CronExpWatchDied    = "*/32 * * * * ?"
)

func newCrontab() *cron.Cron {
	return cron.New(cron.WithSeconds(), cron.WithLogger(&cronLogger{}))
}

type cronLogger struct{}

var gShowCronDetails = configure.GetBool("tick.cron.details.show", false)

func (c *cronLogger) Info(msg string, keysAndValues ...interface{}) {
	if !gShowCronDetails {
		return
	}
	if keysAndValues == nil {
		logger.Logger.Info(msg)
		return
	}
	var d []zap.Field
	for idx, kv := range keysAndValues {
		d = append(d, zap.Any(fmt.Sprintf("%d", idx), kv))
	}
	logger.Logger.Debug(msg, d...)
}

func (c *cronLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	if keysAndValues == nil {
		logger.Logger.Error(msg, zap.Error(err))
		return
	}
	d := []zap.Field{zap.Error(err)}
	for idx, kv := range keysAndValues {
		d = append(d, zap.Any(fmt.Sprintf("%d", idx), kv))
	}
	logger.Logger.Error(msg, d...)
}

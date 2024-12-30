package tick

import (
	"github.com/hootuu/eggcone/fdn/tick/def"
	"github.com/hootuu/eggcone/fdn/tick/schedule"
	"github.com/hootuu/eggcone/fdn/tick/tick"
	"github.com/hootuu/eggcone/fdn/tick/token"
	"github.com/hootuu/gelato/errors"
	"github.com/hootuu/gelato/logger"
	"github.com/hootuu/gelato/sys"
	"go.uber.org/zap"
	"math"
)

const (
	tickCount     = 1
	tickTokenSize = 1
	tokenCount    = tickCount * tickTokenSize
)

var gTokenAllocator = token.New(tokenCount)
var gTicks [tickCount]*Tick
var gListener []Listener

func doTimerInit() {
	sys.Info("Tick Init [ ", tickCount, " ]")
	gTicks = [tickCount]*Tick{}
	for i := 0; i < tickCount; i++ {
		t, err := newTick()
		if err != nil {
			logger.Logger.Error("[init_fail]new Tick failed", zap.Error(err))
			return
		}
		for j := 0; j < tickTokenSize; j++ {
			err = t.bind(gTokenAllocator.Alloc())
			if err != nil {
				logger.Logger.Error("[init_fail]bind token failed", zap.Error(err))
				sys.Exit(err)
				return
			}
		}
		gTicks[i] = t
	}
	sys.Success("Tick Init [ ", tickCount, " ] [OK]")
}

func doRegisterListener(ltn Listener) {
	gListener = append(gListener, ltn)
	logger.Logger.Info("tick.listener register", zap.String("name", ltn.GetName()))
}

func dealJob(scheduleID schedule.ID, job *def.Job) {
	for _, ltn := range gListener {
		if !ltn.Match(job) {
			continue
		}
		ctx, err := ltn.Deal(job)
		result := true
		if err != nil {
			logger.Logger.Error("tick.deal.job failed",
				zap.String("listener", ltn.GetName()),
				zap.Any("job", job),
				zap.Error(err))
			result = false
		}

		err = tick.DealRecord(scheduleID, result, ltn.GetName(), ctx)
		if err != nil {
			logger.Logger.Error("tick.deal.job record failed",
				zap.String("listener", ltn.GetName()),
				zap.Any("job", job),
				zap.Error(err))
		}
	}
}

func doStart() {
	ticksStart()
	gCrontab.Start()
	regHeartbeat()
	regLoadSyncSchedule()
	regWatchDied()
}

func doStop() {
	ticksStop()
	gCrontab.Stop()
}

func ticksStart() {
	for _, t := range gTicks {
		t.start()
	}
}

func ticksStop() {
	for _, t := range gTicks {
		t.stop()
	}
}

func regHeartbeat() {
	_, err := gCrontab.AddFunc(CronExpHeartbeat, func() {
		for _, t := range gTicks {
			t.heartbeat()
		}
	})
	if err != nil {
		logger.Logger.Error("regHeartbeat failed",
			zap.Error(err),
			zap.String("cron", CronExpHeartbeat))
		return
	}
}

func regLoadSyncSchedule() {
	_, err := gCrontab.AddFunc(CronExpSyncSchedule, func() {
		for _, t := range gTicks {
			t.sync()
		}
	})
	if err != nil {
		logger.Logger.Error("regLoadSyncSchedule failed",
			zap.Error(err),
			zap.String("cron", CronExpSyncSchedule))
		return
	}
}

func regWatchDied() {
	_, err := gCrontab.AddFunc(CronExpWatchDied, func() {
		var arr []*tick.Tick
		var err error
		var lstSeqIdx int64 = -1
		for {
			arr, lstSeqIdx, err = tick.LoadWillDied(lstSeqIdx)
			if err != nil {
				logger.Logger.Error("load will died tick failed[ignore, wait next]",
					zap.Error(err))
				return
			}
			if len(arr) == 0 {
				return
			}
			for _, m := range arr {
				dealDiedTick(m)
			}
		}
	})
	if err != nil {
		logger.Logger.Error("regWatchDied failed",
			zap.Error(err),
			zap.String("cron", CronExpWatchDied))
		return
	}
}

func dealDiedTick(m *tick.Tick) {

	unbindTokens, err := tick.Die(m)
	if err != nil {
		logger.Logger.Error("die tick failed",
			zap.Error(err),
			zap.String("tick.ID", m.ID.S()))
		return
	}

	if len(unbindTokens) == 0 {
		return
	}

	for _, t := range unbindTokens {
		innerErr := reAllocTick(t)
		if innerErr != nil {
			logger.Logger.Error("reAllocTick failed[ignore]", zap.Error(innerErr))
			continue
		}
	}
}

func reAllocTick(token token.Token) *errors.Error {
	var least = math.MaxInt
	var cur *Tick = nil

	for _, t := range gTicks {
		size := len(t.Tokens)
		if size < least {
			least = size
			cur = t
		}
	}

	if cur == nil {
		return errors.System("no valid tick")
	}

	err := cur.bind(token)
	if err != nil {
		return err
	}

	return nil
}

package tick

import (
	"github.com/hootuu/eggcone/fdn/tick/def"
	"github.com/hootuu/eggcone/fdn/tick/schedule"
	"github.com/hootuu/eggcone/fdn/tick/tick"
	"github.com/hootuu/eggcone/fdn/tick/token"
	"github.com/hootuu/gelato/errors"
	"github.com/hootuu/gelato/logger"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"sync"
)

type Schedule struct {
	ID      schedule.ID
	Token   token.Token
	Cron    schedule.Cron
	Options *schedule.Options
	Job     *def.Job
	Version uint64

	cronID cron.EntryID
}

func scheduleOf(sM *schedule.Schedule) *Schedule {
	return &Schedule{
		ID:      sM.ID,
		Token:   sM.Token,
		Cron:    sM.Cron,
		Options: sM.Options,
		Job:     sM.Job,
		Version: sM.Version,
	}
}

func (s *Schedule) update(r *schedule.Schedule) {
	s.Cron = r.Cron
	s.Options = r.Options
	s.Job = r.Job
	s.Version = r.Version
}

func (s *Schedule) execute() {
	job := &def.Job{
		ID:      def.NewJobID(),
		Topic:   s.Job.Topic,
		Payload: s.Job.Payload,
		Tag:     s.Job.Tag,
	}
	dealJob(s.ID, job)
}

type Tick struct {
	ID     tick.ID
	Tokens []token.Token

	scheduleDict map[token.Token]*Schedule

	mutex sync.Mutex

	crontab *cron.Cron
}

func newTick() (*Tick, *errors.Error) {
	mID, err := tick.NewTick()
	if err != nil {
		return nil, err
	}

	t := &Tick{
		ID:           mID,
		Tokens:       []token.Token{},
		scheduleDict: make(map[token.Token]*Schedule),
		crontab:      newCrontab(),
	}

	return t, nil
}

func (t *Tick) start() {
	t.sync()
	t.crontab.Start()
}

func (t *Tick) stop() {
	t.crontab.Stop()
}

func (t *Tick) bind(token token.Token) *errors.Error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	err := tick.BindToken(t.ID, token)
	if err != nil {
		return err
	}
	t.Tokens = append(t.Tokens, token)
	return nil
}

func (t *Tick) heartbeat() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	err := tick.Heartbeat(t.ID)
	if err != nil {
		logger.Logger.Error("tick heartbeat failed", zap.String("id", t.ID.S()))
	}
}

func (t *Tick) sync() {
	var sArr []*schedule.Schedule
	var lstSeqIdx int64 = -1
	var err *errors.Error
	for {
		sArr, lstSeqIdx, err = schedule.LoadManyByTokens(t.Tokens, lstSeqIdx)
		if err != nil {
			logger.Logger.Error("schedule.LoadManyByTokens failed", zap.Error(err))
			return
		}
		if len(sArr) == 0 {
			return
		}

		for _, sM := range sArr {
			t.syncSchedule(sM)
		}
	}
}

func (t *Tick) syncSchedule(sM *schedule.Schedule) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	sLocal, exists := t.scheduleDict[sM.Token]
	if !exists {
		if !sM.Available {
			return
		}
		t.put(sM)
		return
	}

	if sLocal.Version == sM.Version {
		return
	}

	if !sM.Available {
		t.remove(sLocal)
		return
	}

	t.update(sLocal, sM)
}

func (t *Tick) put(sM *schedule.Schedule) {
	sLocal := scheduleOf(sM)
	cronID, err := t.crontab.AddFunc(sLocal.Cron.S(), sLocal.execute)
	if err != nil {
		logger.Logger.Error("register cron failed", zap.Error(err))
		return
	}
	sLocal.cronID = cronID
	t.scheduleDict[sM.Token] = sLocal
}

func (t *Tick) remove(sLocal *Schedule) {
	t.crontab.Remove(sLocal.cronID)
	delete(t.scheduleDict, sLocal.Token)
}

func (t *Tick) update(sLocal *Schedule, sM *schedule.Schedule) {
	t.crontab.Remove(sLocal.cronID)
	sLocal.update(sM)
	cronID, err := t.crontab.AddFunc(sLocal.Cron.S(), sLocal.execute)
	logger.Logger.Info("tick.set cron", zap.Int("cronID", int(cronID)))
	if err != nil {
		logger.Logger.Error("register cron failed", zap.Error(err))
		return
	}
	sLocal.cronID = cronID
}

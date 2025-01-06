package schedule

import (
	"github.com/hootuu/eggcone/fdn/tick/def"
	"github.com/hootuu/eggcone/fdn/tick/token"
	"github.com/hootuu/gelato/errors"
	"gorm.io/gorm"
	"time"
)

type Schedule struct {
	ID        ID          `gorm:"column:id;not null;size:32"`
	Token     token.Token `gorm:"column:token;not null;size:32"`
	Title     string      `gorm:"column:title;not null;size:64"`
	Code      string      `gorm:"column:code;not null;size:32"`
	OutID     string      `gorm:"column:out_id;not null;size:64"`
	Cron      Cron        `gorm:"column:cron;not null;size:128"`
	Options   *Options    `gorm:"column:options;type:json"`
	Job       *def.Job    `gorm:"column:job;type:json"`
	Available bool        `gorm:"column:available"`
	Signature string      `gorm:"column:signature;not null;size:128"`
	Version   uint64      `gorm:"column:version"`
	SeqIdx    int64       `gorm:"column:seq_idx"`

	AutoID    int64          `gorm:"column:auto_id;uniqueIndex;autoIncrement"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (m *Schedule) TableName() string {
	return "egg_fdn_tick_schedule"
}

func (m *Schedule) newVerify() *errors.Error {
	if err := m.Token.Verify(); err != nil {
		return err
	}
	if len(m.Title) == 0 {
		return errors.Verify("require title")
	}
	if len(m.Code) == 0 {
		return errors.Verify("require code")
	}
	if err := m.Cron.Verify(); err != nil {
		return err
	}
	if err := m.Job.Verify(); err != nil {
		return err
	}
	if len(m.Signature) == 0 {
		return errors.Verify("require signature")
	}
	return nil
}

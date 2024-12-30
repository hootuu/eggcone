package once

import (
	"github.com/hootuu/eggcone/database/pgx"
	"github.com/hootuu/gelato/errors"
	"github.com/hootuu/gelato/logger"
	"go.uber.org/zap"
	"time"
)

type Status int

const (
	EXECUTING Status = 1
	FAILED    Status = -1
	SUCCESS   Status = 8
)

type Once struct {
	Code        string    `gorm:"column:code;uniqueIndex;not null;size:64"`
	DoServerID  string    `gorm:"column:do_serv_id;index;not null;size:128"`
	DoStatus    Status    `gorm:"column:do_status"`
	DoStartTime time.Time `gorm:"column:do_start_time;index;not null"`
	DoEndTime   time.Time `gorm:"column:do_end_time;index;not null"`

	Version uint64 `gorm:"column:version;default:0"`
}

func (m *Once) TableName() string {
	return "egg_fdn_once"
}

func Get(code string) (*Once, *errors.Error) {
	return pgx.PgGet[Once](DB(), "code = ?", code)
}

func MustGet(code string) (*Once, *errors.Error) {
	m, err := Get(code)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.Verify("no such code: " + code)
	}
	return m, nil
}

func Create(m *Once) *errors.Error {
	tx := DB()
	dbErr := tx.Create(m).Error
	if dbErr != nil {
		logger.Logger.Error("sync.once.Create Failed", zap.Error(dbErr))
		return errors.System("db error", dbErr)
	}
	return nil
}

func SetEnd(m *Once, status Status) *errors.Error {
	dbErr := DB().Model(&Once{}).
		Where("code = ? AND version = ?", m.Code, m.Version).
		Update("do_status", status).
		Update("do_end_time", time.Now()).
		Update("version", m.Version+1).Error
	if dbErr != nil {
		logger.Logger.Error("sync.once.SetEnd Failed", zap.Error(dbErr))
		return errors.System("db error", dbErr)
	}
	return nil
}

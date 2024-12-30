package schedule

import (
	"github.com/hootuu/eggcone/fdn/tick/dbx"
	"github.com/hootuu/gelato/errors"
	"github.com/hootuu/gelato/logger"
	"go.uber.org/zap"
	"time"
)

func New(m *Schedule) (ID, *errors.Error) {
	if err := m.Token.Verify(); err != nil {
		return NilID, err
	}
	if err := m.newVerify(); err != nil {
		return NilID, err
	}

	return newModel(m)
}

func newModel(m *Schedule) (ID, *errors.Error) {
	m.ID = NewID()
	m.Available = true
	if m.Options == nil {
		m.Options = NewDefaultOptions()
	}
	m.Version = 0
	dbErr := dbx.DB().Create(m).Error
	if dbErr != nil {
		logger.Logger.Error("db error", zap.Error(dbErr),
			zap.Any("model", m))
		return NilID, errors.System("db error", dbErr)
	}
	return m.ID, nil
}

func Update(srcM *Schedule, destM *Schedule) *errors.Error {
	up := make(map[string]interface{})
	up["title"] = destM.Title
	up["cron"] = destM.Cron
	up["options"] = destM.Options
	up["job"] = destM.Job
	up["available"] = destM.Available
	up["version"] = srcM.Version + 1
	up["modified_at"] = time.Now()

	dbErr := dbx.DB().Model(srcM).Updates(up).Error
	if dbErr != nil {
		logger.Logger.Error("db field", zap.Error(dbErr),
			zap.Any("srcM", srcM),
			zap.Any("destM", destM))
		return errors.System("db error", dbErr)
	}
	return nil
}

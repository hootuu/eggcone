package tick

import (
	"github.com/hootuu/eggcone/fdn/tick/dbx"
	"github.com/hootuu/eggcone/fdn/tick/token"
	"github.com/hootuu/gelato/errors"
	"gorm.io/gorm"
	"time"
)

func Die(m *Tick) ([]token.Token, *errors.Error) {
	tx := dbx.DB().Begin()
	innerErr := doTxDie(m, tx)
	if innerErr != nil {
		tx.Rollback()
		return nil, innerErr
	}
	innerErr = doTxUnbindTokens(m.ID, tx)
	if innerErr != nil {
		tx.Rollback()
		return nil, innerErr
	}
	tx.Commit()
	return loadBindTokens(m.ID)
}

func doTxDie(m *Tick, tx *gorm.DB) *errors.Error {
	dbErr := tx.Model(&Tick{}).Where("id = ? AND version = ?", m.ID, m.Version).
		Update("living", false).
		Update("updated_at", time.Now()).
		Update("version", m.Version+1).Error
	if dbErr != nil {
		return errors.System("db error", dbErr)
	}
	return nil
}

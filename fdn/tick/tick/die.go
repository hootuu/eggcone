package tick

import (
	"github.com/hootuu/eggcone/fdn/tick/dbx"
	"github.com/hootuu/eggcone/fdn/tick/token"
	"github.com/hootuu/gelato/errors"
	"time"
)

func Die(m *Tick) ([]token.Token, *errors.Error) {
	tx := dbx.DB()
	tx.Begin()
	innerErr := doTxDie(m)
	if innerErr != nil {
		tx.Rollback()
		return nil, innerErr
	}
	innerErr = doTxUnbindTokens(m.ID)
	if innerErr != nil {
		tx.Rollback()
		return nil, innerErr
	}
	tx.Commit()
	return loadBindTokens(m.ID)
}

func doTxDie(m *Tick) *errors.Error {
	tx := dbx.DB()
	dbErr := tx.Model(&Tick{}).Where("id = ? AND version = ?", m.ID, m.Version).
		Update("living", false).
		Update("modified_at", time.Now()).
		Update("version", m.Version+1).Error
	if dbErr != nil {
		return errors.System("db error", dbErr)
	}
	return nil
}

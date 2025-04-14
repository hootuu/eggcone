package pgx

import (
	"github.com/hootuu/gelato/errors"
	"github.com/hootuu/gelato/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Tx = gorm.DB

type TxCtx[M any] struct {
	Tx    *Tx
	Model *M
}

func NewTxCtx[M any](tx *Tx, m *M) *TxCtx[M] {
	return &TxCtx[M]{
		Tx:    tx,
		Model: m,
	}
}

func Transaction(db *gorm.DB, fn func(tx *Tx) *errors.Error) *errors.Error {
	nErr := db.Transaction(func(tx *Tx) error {
		err := fn(tx)
		if err != nil {
			logger.Error.Error("db tx failed", zap.Error(err))
			return err.Native()
		}
		return nil
	})
	if nErr != nil {
		logger.Error.Error("db tx failed", zap.Error(nErr))
		return errors.System("db.Tx err", nErr)
	}
	return nil
}

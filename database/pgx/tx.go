package pgx

import (
	"github.com/hootuu/gelato/errors"
	"gorm.io/gorm"
)

func Tx(db *gorm.DB, fn func(tx *gorm.DB) *errors.Error) *errors.Error {
	nErr := db.Transaction(func(tx *gorm.DB) error {
		err := fn(tx)
		if err != nil {
			return err.Native()
		}
		return nil
	})
	if nErr != nil {
		return errors.System("db.Tx err", nErr)
	}
	return nil
}

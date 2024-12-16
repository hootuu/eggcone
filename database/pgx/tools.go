package pgx

import (
	"fmt"
	"github.com/hootuu/gelato/errors"
	"github.com/hootuu/gelato/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func PgExists(db *gorm.DB, model interface{}, query interface{}, args ...interface{}) (bool, *errors.Error) {
	var exists bool
	err := db.Model(model).
		Select("1").
		Where(query, args...).
		Limit(1).
		Find(&exists).
		Error
	if err != nil {
		return false, errors.System("db error", err)
	}
	return exists, nil
}

func PgGet(db *gorm.DB, model interface{}, cond ...interface{}) *errors.Error {
	tx := db.First(model, cond...)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil
		}
		logger.Logger.Error("db err", zap.Any("cond", cond), zap.Error(tx.Error))
		return errors.System("db err", tx.Error)
	}
	return nil
}

func PgLoad(db *gorm.DB, model interface{}, cond ...interface{}) *errors.Error {
	tx := db.First(model, cond...)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return errors.E("no_such", fmt.Sprintf("no such data %v", cond))
		}
		logger.Logger.Error("db err", zap.Any("cond", cond), zap.Error(tx.Error))
		return errors.System("db err", tx.Error)
	}
	return nil
}

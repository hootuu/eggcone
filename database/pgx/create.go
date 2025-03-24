package pgx

import (
	"github.com/hootuu/gelato/errors"
	"github.com/hootuu/gelato/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func Create[T any](db *gorm.DB, model *T) *errors.Error {
	tx := db.Create(model)
	if tx.Error != nil {
		logger.Logger.Error("db create error", zap.Any("model", model), zap.Error(tx.Error))
		return errors.System("db create error", tx.Error)
	}
	return nil
}

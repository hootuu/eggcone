package pgx

import (
	"github.com/hootuu/gelato/errors"
	"github.com/hootuu/gelato/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func Update[T any](
	db *gorm.DB,
	values map[string]interface{},
	query interface{},
	args ...interface{},
) *errors.Error {
	var model T
	tx := db.Model(&model).Where(query, args...).Updates(values)
	if tx.Error != nil {
		logger.Logger.Error("db error", zap.Any("query", query), zap.Any("args", args), zap.Error(tx.Error))
		return errors.System("db err", tx.Error)
	}
	return nil
}

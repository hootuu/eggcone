package pgx

import (
	"github.com/hootuu/gelato/errors"
	"github.com/hootuu/gelato/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func Save[T any](
	db *gorm.DB,
	model *T,
	rewriteColumns ...string,
) *errors.Error {
	clauseCond := clause.OnConflict{}
	if len(rewriteColumns) == 0 {
		clauseCond.UpdateAll = true
	} else {
		clauseCond.DoUpdates = clause.AssignmentColumns(rewriteColumns)
	}
	tx := db.Clauses(clauseCond).Create(model)
	if tx.Error != nil {
		logger.Logger.Error("db error", zap.Any("model", model), zap.Error(tx.Error))
		return errors.System("db err", tx.Error)
	}
	return nil
}

func SaveMulti[T any](
	db *gorm.DB,
	models []*T,
	rewriteColumns ...string,
) *errors.Error {
	clauseCond := clause.OnConflict{}
	if len(rewriteColumns) == 0 {
		clauseCond.UpdateAll = true
	} else {
		clauseCond.DoUpdates = clause.AssignmentColumns(rewriteColumns)
	}
	tx := db.Clauses(clauseCond).Create(&models)
	if tx.Error != nil {
		logger.Logger.Error("db error", zap.Any("models", models), zap.Error(tx.Error))
		return errors.System("db err", tx.Error)
	}
	return nil
}

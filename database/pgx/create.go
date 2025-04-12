package pgx

import (
	"github.com/avast/retry-go"
	"github.com/hootuu/gelato/configure"
	"github.com/hootuu/gelato/errors"
	"github.com/hootuu/gelato/logger"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

func Create[T any](db *gorm.DB, model *T) *errors.Error {
	_ = retry.Do(func() error {
		tx := db.Create(model)
		if tx.Error != nil {
			logger.Error.Error("db create error", zap.Any("model", model), zap.Error(tx.Error))
			return errors.System("db create error", tx.Error)
		}
		return nil
	},
		retry.Attempts(cast.ToUint(configure.GetInt("db.act.retry.attempts", 3))),
		retry.Delay(configure.GetDuration("db.act.retry.delay", 500*time.Millisecond)),
	)

	return nil
}

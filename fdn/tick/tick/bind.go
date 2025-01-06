package tick

import (
	"github.com/hootuu/eggcone/database/pgx"
	"github.com/hootuu/eggcone/fdn/tick/dbx"
	"github.com/hootuu/eggcone/fdn/tick/token"
	"github.com/hootuu/gelato/errors"
	"github.com/hootuu/gelato/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

func BindToken(tickID ID, token token.Token) *errors.Error {
	m := &Bind{
		Token:     token,
		TickID:    tickID,
		BindTime:  time.Now(),
		Available: true,
		Version:   0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	tx := dbx.DB()
	dbErr := tx.Create(m).Error
	if dbErr != nil {
		logger.Logger.Error("db failed",
			zap.Error(dbErr),
			zap.Any("model", m))
		return errors.System("db error", dbErr)
	}

	return nil
}

func doTxUnbindTokens(tickID ID, tx *gorm.DB) *errors.Error {
	dbErr := tx.Model(&Bind{}).Where("tick_id", tickID).
		Update("available", false).
		Update("updated_at", time.Now()).
		Update("version", 1).Error //todo
	if dbErr != nil {
		return errors.System("db error", dbErr)
	}

	return nil
}

func loadBindTokens(tickID ID) ([]token.Token, *errors.Error) {
	arr, dbErr := pgx.PgFind[Bind](dbx.DB(), "tick_id = ?", tickID)
	if dbErr != nil {
		return nil, errors.System("db error", dbErr)
	}
	if len(arr) == 0 {
		return []token.Token{}, nil
	}
	var tokenArr []token.Token
	for _, m := range arr {
		tokenArr = append(tokenArr, m.Token)
	}

	return tokenArr, nil
}

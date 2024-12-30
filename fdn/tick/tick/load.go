package tick

import (
	"github.com/hootuu/eggcone/database/pgx"
	"github.com/hootuu/eggcone/fdn/tick/dbx"
	"github.com/hootuu/gelato/errors"
	"time"
)

// TODO wait implement
//func LoadByTokenAndLiving(token token.Token, living boolean.Int) (*Model, error) {
//	var arr []*Model
//	sql := `SELECT *
//			FROM flywheel_timer_tick
//			WHERE token=? AND living=? LIMIT 1`
//
//	err := sqlx.Conn().Select(&arr, sql, token, living)
//
//	if err != nil {
//		log.Logger.Error("db failed",
//			zap.Error(err),
//			zap.String("sql", sql),
//			zap.String("token", token.To()),
//			zap.Bool("living", living.IsTrue()))
//		return nil, err
//	}
//
//	if len(arr) == 0 {
//		return nil, nil
//	}
//
//	return arr[0], nil
//}

func LoadWillDied(lstSeqIdx int64) ([]*Tick, int64, error) {
	checkTime := time.Now().Add(time.Minute * -1)
	mArr, dbErr := pgx.PgFind[Tick](dbx.DB(), "seq_idx > ? AND lst_heartbeat_time < ? AND living = ?",
		lstSeqIdx, checkTime, true)
	if dbErr != nil {
		return nil, 0, errors.System("db error", dbErr)
	}
	if len(mArr) == 0 {
		return nil, lstSeqIdx, nil
	}
	newLstSeqIdx := lstSeqIdx
	for _, m := range mArr {
		if m.SeqIdx > newLstSeqIdx {
			newLstSeqIdx = m.SeqIdx
		}
	}
	return mArr, newLstSeqIdx, nil
}

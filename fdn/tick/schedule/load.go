package schedule

import (
	"github.com/hootuu/eggcone/database/pgx"
	"github.com/hootuu/eggcone/fdn/tick/dbx"
	"github.com/hootuu/eggcone/fdn/tick/token"
	"github.com/hootuu/gelato/errors"
	"github.com/hootuu/gelato/io/pagination"
	"github.com/hootuu/gelato/logger"
	"go.uber.org/zap"
)

func LoadByCodeAndOutID(code string, outID string) (*Schedule, *errors.Error) {
	tx := dbx.DB()
	mds, _, err := pgx.PgPageFind[Schedule](tx, pagination.PageALL(), "code = ? AND out_id = ?", code, outID)
	if err != nil {
		return nil, err
	}
	if len(*mds) == 0 {
		return nil, nil
	}
	return (*mds)[0], nil
}

func LoadManyByTokens(tokens []token.Token, lstSeqIdx int64) ([]*Schedule, int64, *errors.Error) {
	tx := dbx.DB()
	mds, _, err := pgx.PgPageFind[Schedule](tx, pagination.PageALL(), "token IN ? AND seq_idx > ?", tokens, lstSeqIdx)
	if err != nil {
		logger.Logger.Error("loadManyByTokens Failed", zap.Any("tokens", tokens), zap.Int64("lstSeqIdx", lstSeqIdx), zap.Error(err))
		return nil, 0, err
	}
	arr := *mds
	if len(arr) == 0 {
		return nil, lstSeqIdx, nil
	}
	newLstSeqIdx := lstSeqIdx
	for _, m := range arr {
		if m.SeqIdx > newLstSeqIdx {
			newLstSeqIdx = m.SeqIdx
		}
	}
	return arr, newLstSeqIdx, nil
}

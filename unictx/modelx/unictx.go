package modelx

import (
	"fmt"
	"github.com/hootuu/eggcone/database/pgx"
	"github.com/hootuu/gelato/crtpto/md5x"
	"github.com/hootuu/gelato/errors"
	"github.com/hootuu/gelato/logger"
	"go.uber.org/zap"
)

type UniCtx struct {
	Basic
	UniKey string `gorm:"column:uni_key;primaryKey;not null;size:100"`
	Value  string `gorm:"column:value;index;not null;size:1000"`
	Domain string `gorm:"column:domain;not null;size:100"`
	KeyStr string `gorm:"column:key_str;not null;size:100"`
}

func (model *UniCtx) TableName() string {
	return "egg_unictx"
}

func Set(domain string, key string, value string) *errors.Error {
	uniKey := genUniKey(domain, key)
	uniCtxM, err := pgx.Get[UniCtx](PgDB(), "uni_key = ?", uniKey)
	if err != nil {
		logger.GetLogger(eggDbName).Error("get unictx failed", zap.Error(err))
		return err
	}
	if uniCtxM == nil {
		uniCtxM = &UniCtx{
			UniKey: uniKey,
			Value:  value,
			Domain: domain,
			KeyStr: key,
		}
		err = pgx.Create[UniCtx](PgDB(), uniCtxM)
		if err != nil {
			logger.GetLogger(eggDbName).Error("create unictx failed", zap.Error(err))
			return err
		}
		return nil
	}
	if uniCtxM.Value == value {
		return nil
	}
	err = pgx.Update[UniCtx](PgDB(), map[string]interface{}{
		"value": value,
	}, "uni_key = ?", uniKey)
	if err != nil {
		logger.GetLogger(eggDbName).Error("update unictx failed", zap.Error(err))
		return err
	}
	return nil
}

func Get(domain string, key string, dfVal string) (string, *errors.Error) {
	uniKey := genUniKey(domain, key)
	uniCtxM, err := pgx.Get[UniCtx](PgDB(), "uni_key = ?", uniKey)
	if err != nil {
		logger.GetLogger(eggDbName).Error("get unictx failed", zap.Error(err))
		return "", err
	}
	if uniCtxM == nil {
		return dfVal, nil
	}
	return uniCtxM.Value, nil
}

func genUniKey(domain string, key string) string {
	return md5x.MD5(fmt.Sprintf("%s.%s", domain, key))
}

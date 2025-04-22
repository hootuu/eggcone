package custodian

import (
	"github.com/hootuu/eggcone/database/pgx"
	"github.com/hootuu/eggcone/eggdbx"
	"github.com/hootuu/eggcone/eggdbx/basic"
	"github.com/hootuu/gelato/errors"
	"github.com/hootuu/gelato/logger"
	"github.com/hootuu/gelato/sys"
	"go.uber.org/zap"
)

type Model struct {
	basic.Basic
	Idx        string `gorm:"column:idx;primaryKey;not null;size:32"`
	PrivateKey []byte `gorm:"column:private_key"`
	Usage      uint64 `gorm:"column:usage"`
	Available  bool   `gorm:"column:available"`
}

func (m *Model) TableName() string {
	return "eggcone_security_custodian"
}

func multiCreate(arr []*Model) *errors.Error {
	err := pgx.MultiCreate[Model](eggdbx.EggPgDB(), arr)
	if err != nil {
		logger.Error.Error("custodian.multiCreate err:", zap.Error(err))
		return nil
	}
	return nil
}

func updateUsage(idx string, usage uint64) *errors.Error {
	mut := make(map[string]interface{})
	mut["usage"] = usage
	err := pgx.Update[Model](eggdbx.EggPgDB(), mut, "idx = ?", idx)
	if err != nil {
		logger.Error.Error("custodian.UpdateUsage err:", zap.Error(err))
		return err
	}
	return nil
}

func loadAvailable(limit int) ([]*Model, *errors.Error) {
	var arr []*Model
	tx := eggdbx.EggPgDB().Model(&Model{}).Where("available = ?", true).Limit(limit).Find(&arr)
	if tx.Error != nil {
		logger.Error.Error("custodian.loadAvailable err", zap.Error(tx.Error))
		return nil, errors.System("loadAvailable error", tx.Error)
	}
	return arr, nil
}

func init() {
	err := eggdbx.EggPgDB().AutoMigrate(&Model{})
	if err != nil {
		sys.Exit(errors.System("init custodian.table Custodian failed", err))
	}
}

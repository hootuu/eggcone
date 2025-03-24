package eggmq

import (
	"github.com/hootuu/eggcone/eggdbx"
	"github.com/hootuu/eggcone/eggmq/modelx"
	"github.com/hootuu/gelato/errors"
	"github.com/hootuu/gelato/sys"
)

func init() {
	err := eggdbx.EggPgDB().AutoMigrate(&modelx.MessageM{})
	if err != nil {
		sys.Exit(errors.System("init eggmq.table message failed", err))
	}
}

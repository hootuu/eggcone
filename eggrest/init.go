package eggrest

import (
	"github.com/hootuu/eggcone/eggdbx"
	"github.com/hootuu/gelato/errors"
	"github.com/hootuu/gelato/sys"
)

func init() {
	err := eggdbx.EggPgDB().AutoMigrate(&GuardM{})
	if err != nil {
		sys.Exit(errors.System("init eggrest.table failed", err))
	}
}

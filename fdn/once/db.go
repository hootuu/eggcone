package once

import (
	"github.com/hootuu/gelato/errors"
	"github.com/hootuu/gelato/sys"
	"gorm.io/gorm"
)

var gDB *gorm.DB

func DB() *gorm.DB {
	if gDB == nil {
		sys.Exit(errors.System("once: must set db first"))
	}
	return gDB
}

func Init(db *gorm.DB) {
	gDB = db
	err := gDB.AutoMigrate(&Once{})
	if err != nil {
		sys.Exit(errors.System("auto migrate for Once failed", err))
		return
	}
}

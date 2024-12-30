package dbx

import (
	"github.com/hootuu/gelato/errors"
	"github.com/hootuu/gelato/sys"
	"gorm.io/gorm"
)

var gDB *gorm.DB

func DB() *gorm.DB {
	if gDB == nil {
		sys.Exit(errors.System("must init db first"))
	}
	return gDB
}

func Init(db *gorm.DB) {
	gDB = db
}

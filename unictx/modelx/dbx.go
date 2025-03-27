package modelx

import (
	"github.com/hootuu/eggcone/database/pgx"
	"github.com/hootuu/gelato/configure"
	"github.com/hootuu/gelato/errors"
	"github.com/hootuu/gelato/sys"
	"gorm.io/gorm"
)

const eggDbName = "eggcone_db"

func init() {
	if !pgx.DatabaseExist(eggDbName) {
		dns := configure.GetString(
			"dbx.eggcone.dns",
			"host=localhost dbname=eggcone_db port=5432 sslmode=disable",
		)
		pgx.Register(eggDbName, dns, &gorm.Config{}, []interface{}{})
	}
	err := PgDB().AutoMigrate(&UniCtx{})
	if err != nil {
		sys.Exit(errors.System("AutoMigrate UniCtx Failed"))
	}
}

func DB() *pgx.Database {
	return pgx.GetDatabase(eggDbName)
}

func PgDB() *gorm.DB {
	return DB().DB()
}

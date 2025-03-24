package eggdbx

import (
	"github.com/hootuu/eggcone/database/pgx"
	"github.com/hootuu/gelato/configure"
	"gorm.io/gorm"
)

func init() {
	dns := configure.GetString(
		"dbx.eggcone.dns",
		"host=localhost dbname=eggcone_db port=5432 sslmode=disable",
	)
	pgx.Register(eggDbName, dns, &gorm.Config{}, []interface{}{})
}

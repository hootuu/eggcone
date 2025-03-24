package eggdbx

import (
	"github.com/hootuu/eggcone/database/pgx"
	"github.com/hootuu/gelato/logger"
	"gorm.io/gorm"
)

const eggDbName = "eggcone_db"

var Logger = logger.GetLogger(eggDbName)

func Egg() *pgx.Database {
	return pgx.GetDatabase(eggDbName)
}

func EggPgDB() *gorm.DB {
	return Egg().DB()
}

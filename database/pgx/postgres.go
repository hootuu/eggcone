package pgx

import (
	"github.com/hootuu/gelato/errors"
	"github.com/hootuu/gelato/sys"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"sync"
)

type Database struct {
	name   string
	dns    string
	config *gorm.Config
	db     *gorm.DB
	once   sync.Once
}

func newDB(name string, dns string, config *gorm.Config) *Database {
	return &Database{
		name:   name,
		dns:    dns,
		config: config,
	}
}

func (db *Database) DB() *gorm.DB {
	db.once.Do(func() {
		sys.Info("# Connecting to db [", db.name, "] ... #")
		var err error
		db.db, err = gorm.Open(postgres.Open(db.dns), db.config)
		if err != nil {
			sys.Exit(errors.System("db error: "+db.dns, err))
		}
		sys.Success("# Connecting to db [", db.name, "] OK #")
	})
	//if sys.RunMode.IsRd() {
	//	return db.db.Debug()
	//}
	return db.db
}

var gPostgresDbDict map[string]*Database
var gPostgresDbMutex sync.Mutex

func Register(name string, dns string, config *gorm.Config, models []interface{}) {
	gPostgresDbMutex.Lock()
	defer gPostgresDbMutex.Unlock()
	if gPostgresDbDict == nil {
		gPostgresDbDict = make(map[string]*Database)
	}
	if _, ok := gPostgresDbDict[name]; ok {
		return
	}
	postgresDB := newDB(name, dns, config)
	err := postgresDB.DB().AutoMigrate(models...)
	if err != nil {
		sys.Exit(errors.System("DB AutoMigrate Error: "+name, err))
	}
	gPostgresDbDict[name] = postgresDB
}

func DatabaseExist(name string) bool {
	gPostgresDbMutex.Lock()
	defer gPostgresDbMutex.Unlock()
	if gPostgresDbDict == nil {
		return false
	}
	_, ok := gPostgresDbDict[name]
	return ok
}

func GetDatabase(name string) *Database {
	gPostgresDbMutex.Lock()
	defer gPostgresDbMutex.Unlock()
	if gPostgresDbDict == nil {
		gPostgresDbDict = make(map[string]*Database)
	}
	db, ok := gPostgresDbDict[name]
	if !ok {
		sys.Exit(errors.System("DB does not exist: " + name))
	}
	return db
}

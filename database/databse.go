package database

import (
	"encoding/json"
	"errors"
	"github.com/glebarez/sqlite"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DBDriverNotSupported = errors.New("database driver not supported")
)

type Database struct {
	*gorm.DB
}

func New(cfg *Config) (database *Database, err error) {
	var driver gorm.Dialector
	logrus.Debugf("database config: %+v", cfg)
	switch cfg.Driver {
	case DBDriverSQLite:
		driver = sqlite.Open(cfg.DSN)
	case DBDriverMySQL:
		driver = mysql.Open(cfg.DSN)
	case DBDriverPostgres:
		driver = postgres.Open(cfg.DSN)
	default:
		return nil, DBDriverNotSupported
	}

	var dbLogLevel = logger.Error
	if cfg.Debug {
		dbLogLevel = logger.Info
	}

	client, err := gorm.Open(driver, &gorm.Config{
		Logger: logger.Default.LogMode(dbLogLevel),
	})
	if err != nil {
		return nil, err
	}

	db, err := client.DB()
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(cfg.MaxIdleConn)
	db.SetMaxOpenConns(cfg.MaxOpenConn)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdle)
	db.SetConnMaxLifetime(cfg.ConnMaxLift)

	dbStat, _ := json.Marshal(db.Stats())
	logrus.Debugf("database stats: %s", dbStat)
	return &Database{client}, nil
}

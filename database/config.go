package database

import (
	"time"
)

const (
	DBDriverPostgres DBDriver = "postgres"
	DBDriverMySQL    DBDriver = "mysql"
	DBDriverSQLite   DBDriver = "sqlite"
)

type DBDriver string

type Config struct {
	Driver DBDriver
	DSN    string

	MaxIdleConn int           `default:"10"`
	MaxOpenConn int           `default:"40"`
	ConnMaxLift time.Duration `default:"0s"`
	ConnMaxIdle time.Duration `default:"0s"`

	Debug bool `default:"false"`
}

package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	driver "gorm.io/driver/clickhouse"
	"gorm.io/gorm"
)

// DbClickHouse clickhouse db
var gClickhouseOrm *gorm.DB

func InitClickHouse(dsn string) {
	opt, err := clickhouse.ParseDSN(dsn)
	if err != nil {
		panic(fmt.Errorf("dsn:%s, err:%s", dsn, err))
	}
	opt.Settings = clickhouse.Settings{
		"max_execution_time": 60,
	}
	opt.DialTimeout = 5 * time.Second
	opt.Compression = &clickhouse.Compression{
		Method: clickhouse.CompressionLZ4,
	}
	db := clickhouse.OpenDB(opt)
	db.SetMaxIdleConns(8)
	db.SetMaxOpenConns(32)
	db.SetConnMaxLifetime(15 * time.Minute)

	initClickHouseOrm(db)
}

func initClickHouseOrm(dbConn *sql.DB) {
	db, err := gorm.Open(driver.New(driver.Config{Conn: dbConn}))
	if err != nil {
		panic("failed to connect database")
	}
	db = db.Debug()
	gClickhouseOrm = db
}

func GetORM() *gorm.DB {
	return gClickhouseOrm
}

func Exec(query string, args ...interface{}) (sql.Result, error) {
	db, err := gClickhouseOrm.DB()
	if err != nil {
		return nil, err
	}
	return db.Exec(query, args...)
}

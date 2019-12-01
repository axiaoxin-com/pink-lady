package database

import (
	"pink-lady/app/logging"

	"github.com/jinzhu/gorm"

	// need by gorm
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// UTDBFile
const UTDBFile = "/tmp/pink-lady-ut.db"

// UTDB 单元测试使用的Sqlite3 DB
func UTDB() *gorm.DB {
	db, _ := NewSQLite3Instance(UTDBFile, true, 10, 10, 10)
	logger, _ := logging.NewLogger("debug", "console", []string{"stderr"}, nil, true, true)
	db.SetLogger(NewLogger(logger))
	return db
}

// MySQL 根据instance名称返回MySQL实例
func MySQL(instance string) *gorm.DB {
	return InstanceMap["mysql"][instance]
}

// SQLite3 根据instance名称返回sqlite3实例
func SQLite3(instance string) *gorm.DB {
	return InstanceMap["sqlite3"][instance]
}

// Postgres 根据实例名称返回pg实例
func Postgres(instance string) *gorm.DB {
	return InstanceMap["postgres"][instance]
}

// MsSQL 根据实例名称返回sqlserver实例
func MsSQL(instance string) *gorm.DB {
	return InstanceMap["mssql"][instance]
}

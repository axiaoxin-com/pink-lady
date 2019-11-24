package db

import (
	"github.com/jinzhu/gorm"

	// need by gorm
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// MySQL 根据instance名称返回MySQL实例
func MySQL(instance string) *gorm.DB {
	return InstanceMap["mysql"][instance]
}

// SQLite3 根据instance名称返回sqlite3实例
func SQLite3(instance string) *gorm.DB {
	db := InstanceMap["sqlite3"][instance]
	if db == nil && instance == "testing" {
		db, _ = NewSQLite3Instance("/tmp/pink-lady-testing.db", true, 10, 10, 10)
	}
	return db
}

// Postgres 根据实例名称返回pg实例
func Postgres(instance string) *gorm.DB {
	return InstanceMap["postgres"][instance]
}

// MsSQL 根据实例名称返回sqlserver实例
func MsSQL(instance string) *gorm.DB {
	return InstanceMap["mssql"][instance]
}

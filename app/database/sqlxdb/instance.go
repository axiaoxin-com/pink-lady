package sqlxdb

import (
	"github.com/jmoiron/sqlx"

	// need by gorm
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// UTDBFile sqlite3单测数据库文件位置
const UTDBFile = "/tmp/pink-lady-ut.db"

// UTDB 单元测试使用的Sqlite3 DB
func UTDB() *sqlx.DB {
	db, _ := NewSQLite3Instance(UTDBFile, true, 10, 10, 10)
	return db
}

// MySQL 根据instance名称返回MySQL实例
func MySQL(instance string) *sqlx.DB {
	return InstanceMap["mysql"][instance]
}

// SQLite3 根据instance名称返回sqlite3实例
func SQLite3(instance string) *sqlx.DB {
	return InstanceMap["sqlite3"][instance]
}

// Postgres 根据实例名称返回pg实例
func Postgres(instance string) *sqlx.DB {
	return InstanceMap["postgres"][instance]
}

// MsSQL 根据实例名称返回sqlserver实例
func MsSQL(instance string) *sqlx.DB {
	return InstanceMap["mssql"][instance]
}

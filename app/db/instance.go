package db

import (
	"github.com/jinzhu/gorm"

	// need by gorm
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func MySQL(instance string) *gorm.DB {
	return DBInstanceMap["mysql"][instance]
}

func SQLite3(instance string) *gorm.DB {
	db := DBInstanceMap["sqlite3"][instance]
	if db == nil && instance == "testing" {
		db, _ = NewSQLite3Instance("/tmp/pink-lady-testing.db", true, 10, 10, 10)
	}
	return db
}

func Postgres(instance string) *gorm.DB {
	return DBInstanceMap["postgres"][instance]
}

func MsSQL(instance string) *gorm.DB {
	return DBInstanceMap["mssql"][instance]
}

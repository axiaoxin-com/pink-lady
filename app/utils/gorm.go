package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/sirupsen/logrus"
)

var DB *gorm.DB

type GormLogger struct{}

func (*GormLogger) Print(v ...interface{}) {
	if v[0] == "sql" {
		logrus.WithFields(logrus.Fields{"module": "gorm", "type": "sql"}).Print(v[3])
	}
	if v[0] == "log" {
		logrus.WithFields(logrus.Fields{"module": "gorm", "type": "log"}).Print(v[2])
	}
}

// InitGormDB init the grom DB
// engine is database type, such as mysql, sqlite3, etc.
// addr is database's address
// name is database dbname
// usename is database username
// password is dabase password
// maxIdleConns sets the maximum number of connections in the idle connection pool
// maxOpenConns sets the maximum number of open connections to the database.
// connMaxLifeMinutes sets the maximum amount of time(minutes) a connection may be reused
// logMode show detailed log
func InitGormDB(engine, addr, name, username, password string, maxIdleConns, maxOpenConns, connMaxLifeMinutes int, logMode bool) {

	var dsn string
	var err error

	switch strings.ToLower(engine) {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", username, password, addr, name)
	case "postgres":
		addr_ := strings.Split(addr, ":")
		dsn = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s", addr_[0], addr_[1], username, name, password)
	case "sqlite3":
		dsn = name
	case "mssql":
		dsn = fmt.Sprintf("sqlserver://%s:%s@%s?database=%s", username, password, addr, name)
	}
	// var scope: use `:=` will declare a new local variable named DB, raise compile error of var not be used
	DB, err = gorm.Open(engine, dsn)
	if err != nil {
		logrus.Fatal(err)
	}
	DB.SetLogger(&GormLogger{})
	DB.LogMode(logMode)
	if engine == "mysql" {
		DB.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8 auto_increment=1")
	}
	DB.DB().SetMaxIdleConns(maxIdleConns)
	DB.DB().SetMaxOpenConns(maxOpenConns)
	DB.DB().SetConnMaxLifetime(time.Duration(connMaxLifeMinutes) * time.Minute)
}

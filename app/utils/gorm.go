package utils

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/sirupsen/logrus"
)

var DB *gorm.DB

func InitGormDB(engine, addr, name, username, password string) {

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
}

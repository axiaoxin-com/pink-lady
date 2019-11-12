package utils

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	// need by gorm
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/sirupsen/logrus"
)

// DB is *gorm.DB, use it to do orm
var DB *gorm.DB

// GormLogger custom gorm logger
type GormLogger struct{}

var sqlRegexp = regexp.MustCompile(`(\$\d+)|\?`)

// Print define how to log
func (*GormLogger) Print(values ...interface{}) {
	if len(values) > 1 {
		level := values[0]
		source := values[1]
		entry := Logger.WithField("source", source)
		if level == "sql" {
			duration := values[2]
			// sql
			var formattedValues []interface{}
			for _, value := range values[4].([]interface{}) {
				indirectValue := reflect.Indirect(reflect.ValueOf(value))
				if indirectValue.IsValid() {
					value = indirectValue.Interface()
					if t, ok := value.(time.Time); ok {
						formattedValues = append(formattedValues, fmt.Sprintf("'%v'", t.Format(time.RFC3339)))
					} else if b, ok := value.([]byte); ok {
						formattedValues = append(formattedValues, fmt.Sprintf("'%v'", string(b)))
					} else if r, ok := value.(driver.Valuer); ok {
						if value, err := r.Value(); err == nil && value != nil {
							formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
						} else {
							formattedValues = append(formattedValues, "NULL")
						}
					} else {
						formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
					}
				} else {
					formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
				}
			}
			entry.WithField("took", duration).Debug(fmt.Sprintf(sqlRegexp.ReplaceAllString(values[3].(string), "%v"), formattedValues...))
		} else {
			entry.Error(values[2:]...)
		}
	} else {
		Logger.Error(values...)
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
func InitGormDB(engine, addr, name, username, password string, maxIdleConns, maxOpenConns, connMaxLifeMinutes int, logMode bool) error {

	var dsn string
	var err error

	switch strings.ToLower(engine) {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", username, password, addr, name)
	case "postgres":
		addrSlice := strings.Split(addr, ":")
		// if pq: SSL is not enabled on the server, set sslmode=disable
		dsn = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s", addrSlice[0], addrSlice[1], username, name, password)
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
	return errors.Wrap(err, "init gormdb error")
}

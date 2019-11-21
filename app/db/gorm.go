package db

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/axiaoxin/pink-lady/app/utils"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/prometheus/common/log"
	"github.com/spf13/viper"

	// need by gorm
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// DBInstanceMapType {"mysql": {"default": client}, "sqlite3": {"default": client}}
type DBInstanceMapType map[string]map[string]*gorm.DB

// DBInstanceMap is *gorm.DB instance group by db engine
var DBInstanceMap = make(DBInstanceMapType)

func (dbimt DBInstanceMapType) Close() {
	for _, ins := range dbimt {
		for _, i := range ins {
			i.Close()
		}
	}
}

// GormLogger custom gorm logger
type GormLogger struct{}

var sqlRegexp = regexp.MustCompile(`(\$\d+)|\?`)

// Print define how to log
func (*GormLogger) Print(values ...interface{}) {
	if len(values) > 1 {
		level := values[0]
		source := values[1]
		entry := utils.Logger.WithField("source", source)
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
		utils.Logger.Error(values...)
	}
}

// InitGorm init the DBInstanceMap
func InitGorm() error {
	var err error
	databaseMap := viper.GetStringMap("database")
	for engine, dbList := range databaseMap {
		switch strings.ToLower(engine) {
		case "mysql":
			if DBInstanceMap["mysql"] == nil {
				DBInstanceMap["mysql"] = make(map[string]*gorm.DB)
			}
			for _, dbItemItf := range dbList.([]interface{}) {
				dbItem := dbItemItf.(map[string]interface{})
				db, err := NewMySQLInstance(
					dbItem["host"].(string),
					int(dbItem["port"].(int64)),
					dbItem["username"].(string),
					dbItem["password"].(string),
					dbItem["dbname"].(string),
					dbItem["logMode"].(bool),
					int(dbItem["maxIdleConns"].(int64)),
					int(dbItem["maxOpenConns"].(int64)),
					int(dbItem["connMaxLifeMinutes"].(int64)),
				)
				if err != nil {
					log.Error("NewMySQLInstance error:", err)
				} else {
					DBInstanceMap["mysql"][dbItem["instance"].(string)] = db
				}
			}
		case "sqlite3":
			if DBInstanceMap["sqlite3"] == nil {
				DBInstanceMap["sqlite3"] = make(map[string]*gorm.DB)
			}
			for _, dbItemItf := range dbList.([]interface{}) {
				dbItem := dbItemItf.(map[string]interface{})
				db, err := NewSQLite3Instance(
					dbItem["dbname"].(string),
					dbItem["logMode"].(bool),
					int(dbItem["maxIdleConns"].(int64)),
					int(dbItem["maxOpenConns"].(int64)),
					int(dbItem["connMaxLifeMinutes"].(int64)),
				)
				if err != nil {
					log.Error("NewSQLite3Instance error:", err)
				} else {
					DBInstanceMap["sqlite3"][dbItem["instance"].(string)] = db
				}
			}
		case "postgres":
			if DBInstanceMap["postgres"] == nil {
				DBInstanceMap["postgres"] = make(map[string]*gorm.DB)
			}
			for _, dbItemItf := range dbList.([]interface{}) {
				dbItem := dbItemItf.(map[string]interface{})
				db, err := NewPostgresInstance(
					dbItem["host"].(string),
					int(dbItem["port"].(int64)),
					dbItem["username"].(string),
					dbItem["password"].(string),
					dbItem["dbname"].(string),
					dbItem["sslmode"].(string),
					dbItem["logMode"].(bool),
					int(dbItem["maxIdleConns"].(int64)),
					int(dbItem["maxOpenConns"].(int64)),
					int(dbItem["connMaxLifeMinutes"].(int64)),
				)
				if err != nil {
					log.Error("NewPostgresInstance error:", err)
				} else {
					DBInstanceMap["postgres"][dbItem["instance"].(string)] = db
				}
			}
		case "mssql":
			if DBInstanceMap["mssql"] == nil {
				DBInstanceMap["mssql"] = make(map[string]*gorm.DB)
			}
			for _, dbItemItf := range dbList.([]interface{}) {
				dbItem := dbItemItf.(map[string]interface{})
				db, err := NewMsSQLInstance(
					dbItem["host"].(string),
					int(dbItem["port"].(int64)),
					dbItem["username"].(string),
					dbItem["password"].(string),
					dbItem["dbname"].(string),
					dbItem["logMode"].(bool),
					int(dbItem["maxIdleConns"].(int64)),
					int(dbItem["maxOpenConns"].(int64)),
					int(dbItem["connMaxLifeMinutes"].(int64)),
				)
				if err != nil {
					log.Error("NewMsSQLInstance error:", err)
				} else {
					DBInstanceMap["mssql"][dbItem["instance"].(string)] = db
				}
			}
		}
	}
	return err
}

// NewSQLite3Instance return gorm sqlite3 instance
// dbname is dbfile path
// logMode show detailed log
// maxIdleConns sets the maximum number of connections in the idle connection pool
// maxOpenConns sets the maximum number of open connections to the database.
// connMaxLifeMinutes sets the maximum amount of time(minutes) a connection may be reused
func NewSQLite3Instance(dbname string, logMode bool, maxIdleConns, maxOpenConns, connMaxLifeMinutes int) (*gorm.DB, error) {
	db, err := gorm.Open("sqlite3", dbname)
	if err != nil {
		return nil, errors.Wrap(err, "NewSQLite3Instance gorm open error")
	}
	db.SetLogger(&GormLogger{})
	db.LogMode(logMode)
	db.DB().SetMaxIdleConns(maxIdleConns)                                       // 设置连接池中的最大闲置连接数
	db.DB().SetMaxOpenConns(maxOpenConns)                                       // 设置数据库的最大连接数量
	db.DB().SetConnMaxLifetime(time.Duration(connMaxLifeMinutes) * time.Minute) // 设置连接的最大可复用时间
	return db, nil
}

// NewMySQLInstance return gorm mysql instance
// host is database's host
// port is database's port
// dbname is database's dbname
// usename is database's username
// password is database's password
// logMode show detailed log
// maxIdleConns sets the maximum number of connections in the idle connection pool
// maxOpenConns sets the maximum number of open connections to the database.
// connMaxLifeMinutes sets the maximum amount of time(minutes) a connection may be reused
func NewMySQLInstance(host string, port int, username, password, dbname string, logMode bool, maxIdleConns, maxOpenConns, connMaxLifeMinutes int) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, dbname)
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		return nil, errors.Wrap(err, "NewMySQLInstance gorm open error")
	}
	db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8 auto_increment=1")
	db.SetLogger(&GormLogger{})
	db.LogMode(logMode)
	db.DB().SetMaxIdleConns(maxIdleConns)                                       // 设置连接池中的最大闲置连接数
	db.DB().SetMaxOpenConns(maxOpenConns)                                       // 设置数据库的最大连接数量
	db.DB().SetConnMaxLifetime(time.Duration(connMaxLifeMinutes) * time.Minute) // 设置连接的最大可复用时间
	return db, nil
}

// NewPostgresInstance return gorm postgresql instance
// host is database's host
// port is database's port
// dbname is database's dbname
// usename is database's username
// sslmode ssl is disable or not
// password is database's password
// logMode show detailed log
// maxIdleConns sets the maximum number of connections in the idle connection pool
// maxOpenConns sets the maximum number of open connections to the database.
// connMaxLifeMinutes sets the maximum amount of time(minutes) a connection may be reused
func NewPostgresInstance(host string, port int, username, password, dbname, sslmode string, logMode bool, maxIdleConns, maxOpenConns, connMaxLifeMinutes int) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s", host, port, username, dbname, password, sslmode)
	db, err := gorm.Open("postgres", dsn)
	if err != nil {
		return nil, errors.Wrap(err, "NewPostgresInstance gorm open error")
	}
	db.SetLogger(&GormLogger{})
	db.LogMode(logMode)
	db.DB().SetMaxIdleConns(maxIdleConns)                                       // 设置连接池中的最大闲置连接数
	db.DB().SetMaxOpenConns(maxOpenConns)                                       // 设置数据库的最大连接数量
	db.DB().SetConnMaxLifetime(time.Duration(connMaxLifeMinutes) * time.Minute) // 设置连接的最大可复用时间
	return db, nil
}

// NewMsSQLInstance return gorm sqlserver instance
// host is database's host
// port is database's port
// dbname is database's dbname
// usename is database's username
// password is database's password
// logMode show detailed log
// maxIdleConns sets the maximum number of connections in the idle connection pool
// maxOpenConns sets the maximum number of open connections to the database.
// connMaxLifeMinutes sets the maximum amount of time(minutes) a connection may be reused
func NewMsSQLInstance(host string, port int, username, password, dbname string, logMode bool, maxIdleConns, maxOpenConns, connMaxLifeMinutes int) (*gorm.DB, error) {
	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s", username, password, host, port, dbname)
	db, err := gorm.Open("mssql", dsn)
	if err != nil {
		return nil, errors.Wrap(err, "NewMsSQLInstance gorm open error")
	}
	db.SetLogger(&GormLogger{})
	db.LogMode(logMode)
	db.DB().SetMaxIdleConns(maxIdleConns)                                       // 设置连接池中的最大闲置连接数
	db.DB().SetMaxOpenConns(maxOpenConns)                                       // 设置数据库的最大连接数量
	db.DB().SetConnMaxLifetime(time.Duration(connMaxLifeMinutes) * time.Minute) // 设置连接的最大可复用时间
	return db, nil
}

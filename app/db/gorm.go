package db

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/axiaoxin/pink-lady/app/logging"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/spf13/viper"

	// need by gorm
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// InstanceMapType {"mysql": {"default": client}, "sqlite3": {"default": client}}
type InstanceMapType map[string]map[string]*gorm.DB

// InstanceMap is *gorm.DB instance group by db engine
var InstanceMap = make(InstanceMapType)

func (dbimt InstanceMapType) Close() {
	for _, ins := range dbimt {
		for _, i := range ins {
			i.Close()
		}
	}
}

var sqlRegexp = regexp.MustCompile(`(\$\d+)|\?`)

// InitGorm init the InstanceMap
func InitGorm() error {
	var err error
	var db *gorm.DB

	databaseMap := viper.GetStringMap("database")
	logger := NewLogger(logging.CloneLogger().Named("gorm"))
	for engine, dbList := range databaseMap {
		switch strings.ToLower(engine) {
		case "mysql":
			if InstanceMap["mysql"] == nil {
				InstanceMap["mysql"] = make(map[string]*gorm.DB)
			}
			for _, dbItemItf := range dbList.([]interface{}) {
				dbItem := dbItemItf.(map[string]interface{})
				db, err = NewMySQLInstance(
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
				if err == nil {
					db.SetLogger(logger)
					InstanceMap["mysql"][dbItem["instance"].(string)] = db
				}
			}
		case "sqlite3":
			if InstanceMap["sqlite3"] == nil {
				InstanceMap["sqlite3"] = make(map[string]*gorm.DB)
			}
			for _, dbItemItf := range dbList.([]interface{}) {
				dbItem := dbItemItf.(map[string]interface{})
				db, err = NewSQLite3Instance(
					dbItem["dbname"].(string),
					dbItem["logMode"].(bool),
					int(dbItem["maxIdleConns"].(int64)),
					int(dbItem["maxOpenConns"].(int64)),
					int(dbItem["connMaxLifeMinutes"].(int64)),
				)
				if err == nil {
					db.SetLogger(logger)
					InstanceMap["sqlite3"][dbItem["instance"].(string)] = db
				}
			}
		case "postgres":
			if InstanceMap["postgres"] == nil {
				InstanceMap["postgres"] = make(map[string]*gorm.DB)
			}
			for _, dbItemItf := range dbList.([]interface{}) {
				dbItem := dbItemItf.(map[string]interface{})
				db, err = NewPostgresInstance(
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
				if err == nil {
					db.SetLogger(logger)
					InstanceMap["postgres"][dbItem["instance"].(string)] = db
				}
			}
		case "mssql":
			if InstanceMap["mssql"] == nil {
				InstanceMap["mssql"] = make(map[string]*gorm.DB)
			}
			for _, dbItemItf := range dbList.([]interface{}) {
				dbItem := dbItemItf.(map[string]interface{})
				db, err = NewMsSQLInstance(
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
				if err == nil {
					db.SetLogger(logger)
					InstanceMap["mssql"][dbItem["instance"].(string)] = db
				}
			}
		}
	}
	if len(InstanceMap) == 0 {
		err = errors.New("db InstanceMap is empty, check your config file")
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
	db.LogMode(logMode)
	db.DB().SetMaxIdleConns(maxIdleConns)                                       // 设置连接池中的最大闲置连接数
	db.DB().SetMaxOpenConns(maxOpenConns)                                       // 设置数据库的最大连接数量
	db.DB().SetConnMaxLifetime(time.Duration(connMaxLifeMinutes) * time.Minute) // 设置连接的最大可复用时间
	return db, nil
}

package models

import (
	"context"
	"time"

	"github.com/axiaoxin-com/logging"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	// DefaultMySQLMaxIdleConns db默认最大空闲连接数
	DefaultMySQLMaxIdleConns = 2
	// DefaultMySQLMaxOpenConns db默认最大连接数
	DefaultMySQLMaxOpenConns = 10
	// DefaultMySQLConnMaxLifetime db连接默认可复用的最大时间
	DefaultMySQLConnMaxLifetime = time.Hour
)

// DBConfig doris相关的配置
type DBConfig struct {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	// "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	DSN             string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	*gorm.Config
}

// NewMySQLDB 创建gorm mysql db client
func NewMySQLDB(c DBConfig) (*gorm.DB, error) {
	// 设置config默认值
	if c.MaxIdleConns == 0 {
		c.MaxIdleConns = DefaultMySQLMaxIdleConns
	}
	if c.MaxOpenConns == 0 {
		c.MaxOpenConns = DefaultMySQLMaxOpenConns
	}
	if c.ConnMaxLifetime == 0 {
		c.ConnMaxLifetime = DefaultMySQLConnMaxLifetime
	}

	m := mysql.Open(c.DSN)
	if c.Config == nil {
		c.Config = &gorm.Config{}
	} else if c.Config.Logger == nil {
		c.Config.Logger = logging.NewGormLogger(zap.InfoLevel, zap.DebugLevel, viper.GetDuration("logging.access_logger.slow_threshold")*time.Millisecond)
	}

	gormdb, err := gorm.Open(m, c.Config)
	if err != nil {
		return nil, errors.Wrap(err, "gorm open error")
	}

	db, err := gormdb.DB()
	if err != nil {
		return nil, errors.Wrap(err, "gormdb get db error")
	}
	db.SetMaxIdleConns(c.MaxIdleConns)
	db.SetConnMaxLifetime(c.ConnMaxLifetime)
	db.SetMaxOpenConns(c.MaxOpenConns)

	if err := db.Ping(); err != nil {
		return nil, errors.Wrap(err, "gorm db ping error")
	}

	return gormdb, nil
}

// CheckMySQL 检查 mysql 服务状态
func CheckMySQL(ctx context.Context) map[string]string {
	// 检查 mysql
	result := map[string]string{
		"db": "ok",
	}
	db, err := DB.DB()
	if err != nil {
		result["db"] = err.Error()
	} else if err := db.Ping(); err != nil {
		result["db"] = err.Error()
	}
	return result
}

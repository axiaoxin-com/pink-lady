package services

import (
	"time"

	"github.com/axiaoxin-com/goutils"
	"github.com/axiaoxin-com/logging"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// DB 全局数据库对象
var DB *gorm.DB

// GormLogger 自定义 gorm logger
var GormLogger = logging.NewGormLogger(zap.InfoLevel, viper.GetDuration("logging.access_logger.slow_threshold")*time.Millisecond)

// GormMySQL 获取带有自定义 logger 的 gorm db 实例
func GormMySQL(which string) (*gorm.DB, error) {
	db, err := goutils.GormMySQL(which)
	if err != nil {
		return nil, err
	}
	db = db.Session(&gorm.Session{
		Logger: GormLogger,
	})
	return db, nil
}

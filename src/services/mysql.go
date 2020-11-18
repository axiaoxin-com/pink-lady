package services

import (
	"context"
	"time"

	"github.com/axiaoxin-com/goutils"
	"github.com/axiaoxin-com/logging"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// DB 获取带有ctx 和自定义 logger 的 gorm db 实例
func DB(ctx context.Context) *gorm.DB {
	env := viper.GetString("env")
	db, err := goutils.GormMySQL(env)
	if err != nil {
		panic(env + " get gorm mysql instance error:" + err.Error())
	}
	db = db.Session(&gorm.Session{
		Logger: logging.NewGormLogger(zap.InfoLevel, viper.GetDuration("logging.access_logger.slow_threshold")*time.Millisecond),
	})
	return db.WithContext(ctx)
}

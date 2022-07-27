// Package models 定义数据库 model
package models

import (
	"fmt"

	"github.com/axiaoxin-com/goutils"
	"github.com/axiaoxin-com/logging"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var (
	// DB mysql gorm db client
	DB *gorm.DB
	// Redis redis client
	Redis *redis.Client
)

// Init 初始化逻辑
func Init() {
	var err error
	// 初始化 gorm db
	dsn := viper.GetString(fmt.Sprintf("mysql.%s.dbname.dsn", viper.GetString("env")))
	if dsn != "" && DB == nil {
		DB, err = NewMySQLDB(DBConfig{
			DSN: dsn,
		})
		if err != nil {
			logging.Fatal(nil, "Init DB NewMySQLDB error:"+err.Error())
		} else {
			if viper.GetString("server.mode") == "debug" {
				DB = DB.Debug()
			}
			logging.Info(nil, "Init DB success")
		}
	}

	// 初始化 redis
	redisAddr := viper.GetString(fmt.Sprintf("redis.%s.addr", viper.GetString("env")))
	if redisAddr != "" && Redis == nil {
		Redis, err = goutils.RedisClient(viper.GetString("env"))
		if err != nil {
			logging.Fatal(nil, "Init Redis get client error:"+err.Error())
		}
		logging.Info(nil, "Init Redis success")
	}
}

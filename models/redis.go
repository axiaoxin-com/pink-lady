package models

import (
	"context"

	"github.com/axiaoxin-com/goutils"
	"github.com/spf13/viper"
)

// CheckRedis 检查 redis 服务状态
func CheckRedis(ctx context.Context) map[string]string {
	env := viper.GetString("env")
	envRedisStatus := "ok"
	if envRedis, err := goutils.RedisClient(env); err != nil {
		envRedisStatus = err.Error()
	} else if _, err := envRedis.Ping(context.TODO()).Result(); err != nil {
		envRedisStatus = err.Error()
	}
	return map[string]string{
		env: envRedisStatus,
	}
}

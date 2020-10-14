// 服务健康度检查

package services

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/axiaoxin-com/goutils"
	"github.com/axiaoxin-com/logging"
	"github.com/spf13/viper"
)

// CheckMySQL 检查 mysql 服务状态
func CheckMySQL(ctx context.Context) map[string]string {
	// 检查本地 mysql
	localhostMySQLStatus := "ok"
	if localhostMySQL, err := goutils.GormMySQL("localhost"); err != nil {
		localhostMySQLStatus = err.Error()
	} else if sqlDB, err := localhostMySQL.DB(); err != nil {
		localhostMySQLStatus = err.Error()
	} else if err := sqlDB.Ping(); err != nil {
		localhostMySQLStatus = err.Error()
	}
	return map[string]string{
		"localhost": localhostMySQLStatus,
	}
}

// CheckRedis 检查 redis 服务状态
func CheckRedis(ctx context.Context) map[string]string {
	localhostRedisStatus := "ok"
	if localhostRedis, err := goutils.RedisClient("localhost"); err != nil {
		localhostRedisStatus = err.Error()
	} else if _, err := localhostRedis.Ping(context.TODO()).Result(); err != nil {
		localhostRedisStatus = err.Error()
	}
	return map[string]string{
		"localhost": localhostRedisStatus,
	}

}

// CheckAtomicLevelServer 检查 logging 的 AtomicLevel server 是否正常
func CheckAtomicLevelServer(ctx context.Context) string {
	client := &http.Client{}
	url := "http://localhost" + viper.GetString("logging.atomic_level_server.addr") + viper.GetString("logging.atomic_level_server.path")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err.Error()
	}
	req.SetBasicAuth(viper.GetString("basic_auth.username"), viper.GetString("basic_auth.password"))
	req.Header.Set(string(logging.TraceIDKeyname), logging.CtxTraceID(ctx))
	rsp, err := client.Do(req)
	if err != nil {
		return err.Error()
	}
	lvl, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err.Error()
	}
	type levelJSON struct {
		Level string `json:"level"`
	}
	level := levelJSON{}
	if err := json.Unmarshal(lvl, &level); err != nil {
		return err.Error()
	}
	if level.Level == "" {
		return "atomiclevel server response invalid level json."
	}
	return "ok"
}

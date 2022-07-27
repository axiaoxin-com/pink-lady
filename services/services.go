// Package services 加载或初始化外部依赖服务
package services

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/axiaoxin-com/logging"
	"github.com/spf13/viper"
)

// Init 相关依赖服务的初始化或加载操作
func Init() {
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
